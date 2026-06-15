package handlers

import (
	"net/http"
	"strconv"

	"pizza-app/models"

	"github.com/gin-gonic/gin"
	"pizza-app/repositories"
	"log"
)

func GetMenuHandler(c *gin.Context) {
	pizzas, err := repositories.GetAllPizzas()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pizzas)
}

func CreatePizzaHandler(c *gin.Context) {

	var newPizza models.Pizza

	if err := c.ShouldBindJSON(&newPizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON",
		})
		return
	}

	if newPizza.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name cannot be empty",
		})
		return
	}

	if newPizza.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "price must be greater than 0",
		})
		return
	}

	newPizza.ID = models.NextPizzaID

	models.NextPizzaID++

	err := repositories.CreatePizza(newPizza)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid pizza id",
		})
		return
	}

	var pizza models.Pizza

	if err := c.ShouldBindJSON(&pizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json",
		})
		return
	}

	updatedPizza, err := repositories.UpdatePizza(id, pizza)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedPizza)
}

func DeletePizzaHandler(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid pizza ID",
		})
		return
	}
	pizza, err := repositories.DeletePizza(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pizza)
}

func GetPizzaOrder(c *gin.Context) {
 var req models.Order


 if err := c.ShouldBindJSON(&req); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
		"error" : "Invalid request body!",
	})
	return
 }
// log.Println("pizza_id:", req.PizzaId)
// log.Println("quantity:", req.Quantity)

 order, err := repositories.CreateOrder(req.PizzaId, req.Quantity)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)

}

func GetOrdersHandler(c *gin.Context) {
	orders, err := repositories.GetAllOrdersWithPizzaName()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "invalid id",
		})
		return
	}

	var req models.UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		
		return
	}

	err = repositories.UpdateOrderStatus(
		id,
		req.Status,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message" : "status updated",
	})
}

func GetDashboardStatsHandler(c *gin.Context) {
	stats, err := repositories.GetDashboardStats() 

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GetBestSellerHandler(c *gin.Context) {
	pizza, err := repositories.GetBestSellingPizza()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, pizza)
}

func SearchPizzaHandler(c *gin.Context) {
	query := c.Query("q")
	log.Println("Search query:", query)

	pizzas, err := repositories.SearchPizza(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pizzas)
}