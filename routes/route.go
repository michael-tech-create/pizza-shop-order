package routes

import (
	"pizza-app/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// admin := router.Group("/")

	router.GET("/menu", handlers.GetMenuHandler)

	router.POST("/api/pizzas", handlers.CreatePizzaHandler)
	router.GET("/api/pizzas/:id", handlers.GetPizzaByIdHandler)
	router.PUT("/api/pizzas/:id", handlers.UpdatePizzaHandler)
	router.DELETE("/api/pizzas/:id", handlers.DeletePizzaHandler)

	router.POST("/api/pizzas/:id/images", handlers.UploadPizzaImageHandler)
	router.GET("/api/pizzas/:id/images", handlers.GetPizzaImagesHandler)
	router.DELETE("/api/images/:id", handlers.DeletePizzaImageHandler)

	router.POST("/api/orders", handlers.GetPizzaOrder)
	router.GET("/api/orders", handlers.GetOrdersHandler)
	router.PATCH("/api/orders/:id/status", handlers.UpdateOrderStatus)
	router.POST("/api/auth/login")

	router.GET("/api/admin/stats", handlers.GetDashboardStatsHandler)
	router.GET("/api/admin/best-seller", handlers.GetBestSellerHandler)
	router.GET("/api/admin/search", handlers.SearchPizzaHandler)
}