package routes

import (
	"pizza-app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/menu", handlers.GetMenuHandler)
	router.POST("/pizzas", handlers.CreatePizzaHandler)
	router.GET("/api/pizzas/:id", handlers.GetPizzaByIdHandler)
	router.PUT("/api/pizzas/:id", handlers.UpdatePizzaHandler)
	router.DELETE("/api/pizzas/:id", handlers.DeletePizzaHandler)
	router.POST("api/orders", handlers.GetPizzaOrder)
	//
}