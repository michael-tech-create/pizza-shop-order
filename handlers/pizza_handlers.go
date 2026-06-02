package handlers

import (
	"net/http"
	"strconv"

	"pizza-app/models"

	"github.com/gin-gonic/gin"
	"pizza-app/repositories"
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

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid pizza ID",
		})
		return
	}

	pizza, err := repositories.GetPizzaByID(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pizza)
	// return
}

func UpdatePizzaHandler(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid pizza ID",
		})
		return
	}

	var updatedPizza models.Pizza

	// Parse JSON body
	if err := c.ShouldBindJSON(&updatedPizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON",
		})
		return
	}

	// Validate input
	if updatedPizza.Name == "" || updatedPizza.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name cannot be empty and price must be greater than 0",
		})
		return
	}

	// Find and update pizza
	for i, pizza := range models.Pizzas {

		if pizza.ID == id {

			updatedPizza.ID = id

			models.Pizzas[i] = updatedPizza

			c.JSON(http.StatusOK, updatedPizza)

			return
		}
	}

	// Pizza not found
	c.JSON(http.StatusNotFound, gin.H{
		"error": "pizza not found",
	})
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


 order, err := repositories.CreateOrder(req.PizzaId, req.Quantity)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, order)

}