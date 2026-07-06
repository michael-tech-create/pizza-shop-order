package routes

import (
	"pizza-app/handlers"
	"pizza-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {


	router.GET("/menu", handlers.GetMenuHandler)

	router.GET("/api/pizzas/:id", handlers.GetPizzaByIdHandler)
	router.GET("/api/pizzas/:id/images", handlers.GetPizzaImagesHandler)

	router.POST("/api/orders", handlers.GetPizzaOrder)
	router.GET("/api/admin/search", handlers.SearchPizzaHandler)

	router.POST("/api/auth/login", handlers.LoginHandler)


	admin := router.Group("/")
	admin.Use(middleware.RequireAuth())
	{
		admin.GET("/api/auth/verify", handlers.VerifySessionHandler)

		// Pizza CRUD (mutating routes — read routes above stay public so
		// the customer-facing menu page keeps working without a token)
		admin.POST("/api/pizzas", handlers.CreatePizzaHandler)
		admin.PUT("/api/pizzas/:id", handlers.UpdatePizzaHandler)
		admin.DELETE("/api/pizzas/:id", handlers.DeletePizzaHandler)

		// Pizza Images (mutating)
		admin.POST("/api/pizzas/:id/images", handlers.UploadPizzaImageHandler)
		admin.DELETE("/api/images/:id", handlers.DeletePizzaImageHandler)

		// Orders (admin views + mutates status)
		admin.GET("/api/orders", handlers.GetOrdersHandler)
		admin.PATCH("/api/orders/:id/status", handlers.UpdateOrderStatus)

		// Dashboard / stats / search
		admin.GET("/api/admin/stats", handlers.GetDashboardStatsHandler)
		admin.GET("/api/admin/best-seller", handlers.GetBestSellerHandler)
	}
}