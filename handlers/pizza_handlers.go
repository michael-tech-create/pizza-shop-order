package handlers

import (
	"net/http"
	"strconv"

	"pizza-app/models"

	"github.com/gin-gonic/gin"
)

func GetMenuHandler(c *gin.Context) {
	c.JSON(http.StatusOK, models.Pizzas)
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

	models.Pizzas = append(models.Pizzas, newPizza)

	c.JSON(http.StatusCreated, newPizza)
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

	for _, pizza := range models.Pizzas {

		if pizza.ID == id {
			c.JSON(http.StatusOK, pizza)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "pizza not found",
	})
	return
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

	for i, pizza := range models.Pizzas {

		if pizza.ID == id {
			models.Pizzas = append(models.Pizzas[:i], models.Pizzas[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "pizza deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "pizza not found",
	})
	return
}