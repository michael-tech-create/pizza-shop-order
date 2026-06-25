package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"pizza-app/models"
	"pizza-app/repositories"

	"github.com/gin-gonic/gin"
)

func GetMenuHandler(c *gin.Context) {
	pizzas, err := repositories.GetAllPizzas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pizzas)
}

func CreatePizzaHandler(c *gin.Context) {
	var newPizza models.Pizza

	if err := c.ShouldBindJSON(&newPizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if newPizza.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
		return
	}

	if newPizza.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
		return
	}

	if err := repositories.CreatePizza(newPizza); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "pizza created successfully"})
}

func GetPizzaByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pizza ID"})
		return
	}

	pizza, err := repositories.GetPizzaByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	images, err := repositories.GetPizzaImages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pizza":  pizza,
		"images": images,
	})
}

func UpdatePizzaHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pizza id"})
		return
	}

	var pizza models.Pizza
	if err := c.ShouldBindJSON(&pizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	updatedPizza, err := repositories.UpdatePizza(id, pizza)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPizza)
}

func DeletePizzaHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pizza ID"})
		return
	}

	pizza, err := repositories.DeletePizza(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pizza)
}

func GetPizzaOrder(c *gin.Context) {
	var req models.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.PizzaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pizza_id must be a positive integer"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be greater than 0"})
		return
	}

	order, err := repositories.CreateOrder(req.PizzaID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrdersHandler(c *gin.Context) {
	orders, err := repositories.GetAllOrdersWithPizzaName()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allowed := map[string]bool{
		"pending":   true,
		"preparing": true,
		"delivered": true,
		"cancelled": true,
	}
	if !allowed[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status must be one of: pending, preparing, delivered, cancelled",
		})
		return
	}

	if err := repositories.UpdateOrderStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func GetDashboardStatsHandler(c *gin.Context) {
	stats, err := repositories.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GetBestSellerHandler(c *gin.Context) {
	pizza, err := repositories.GetBestSellingPizza()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pizza)
}

func SearchPizzaHandler(c *gin.Context) {
	query := c.Query("q")
	log.Println("Search query:", query)

	pizzas, err := repositories.SearchPizza(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pizzas)
}

func UploadPizzaImageHandler(c *gin.Context) {
	pizzaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pizza id"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read files"})
		return
	}

	files := form.File["images"]
	log.Println("Files Received:", len(files))
	var uploaded []string

	for _, file := range files {
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		path := filepath.Join("uploads", filename)

		if err := c.SaveUploadedFile(file, path); err != nil {
			log.Println(err)
			continue
		}

		imageURL := "/uploads/" + filename
		log.Println("Writing To DB:", imageURL)

		if err := repositories.SavePizzaImage(pizzaID, imageURL); err != nil {
			log.Println("DB Write Error:", err)
			continue
		}
		uploaded = append(uploaded, imageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"pizza_id": pizzaID,
		"images":   uploaded,
	})
}

func GetPizzaImagesHandler(c *gin.Context) {
	pizzaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pizza id"})
		return
	}

	images, err := repositories.GetPizzaImages(pizzaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pizza_id": pizzaID,
		"images":   images,
	})
}

func DeletePizzaImageHandler(c *gin.Context) {
	imageID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"deleted_image": imageID,
	})
}