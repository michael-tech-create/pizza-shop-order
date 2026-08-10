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

	// Public search — used by the customer-facing menu search box.
	// admin.GET("/api/admin/search", ...) below stays too, since the
	// admin dashboard's own search UI already calls that one.
	router.GET("/api/search", handlers.SearchPizzaHandler)

	router.POST("/api/orders", handlers.GetPizzaOrder)

	router.POST("/api/auth/login", handlers.LoginHandler)

	// Categories — public read so the menu filter can populate itself
	router.GET("/api/categories", handlers.GetCategoriesHandler)

	// Reviews — public read for the homepage; submitting one requires a
	// logged-in customer whose own delivered order it must reference.
	router.GET("/api/reviews", handlers.GetPublicReviewsHandler)

	// Customer accounts — public signup/login. Guest checkout at
	// /api/orders doesn't require any of this; these just enable the
	// optional "create an account" path for order history.
	router.POST("/api/customers/signup", handlers.CustomerSignupHandler)
	router.POST("/api/customers/login", handlers.CustomerLoginHandler)

	customer := router.Group("/")
	customer.Use(middleware.RequireCustomerAuth())
	{
		customer.GET("/api/customers/orders", handlers.GetMyOrdersHandler)
		customer.POST("/api/reviews", handlers.CreateReviewHandler)
	}

	// Payments — public routes. /initialize and /verify are called by the
	// customer-facing checkout flow (no admin session exists yet at that
	// point), and /webhook is called by Paystack's servers directly, which
	// can't attach a Bearer token. Security instead comes from: the amount
	// always being read from our own DB (never client input), the
	// reference being a server-generated lookup key, and the webhook
	// requiring a valid X-Paystack-Signature.
	router.POST("/api/payments/initialize", handlers.InitializePaymentHandler)
	router.GET("/api/payments/verify/:reference", handlers.VerifyPaymentHandler)
	router.POST("/api/payments/webhook", handlers.PaystackWebhookHandler)

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
		admin.GET("/api/admin/search", handlers.SearchPizzaHandler)

		admin.POST("/api/categories", handlers.CreateCategoryHandler)
		admin.PUT("/api/categories/:id", handlers.UpdateCategoryHandler)
		admin.DELETE("/api/categories/:id", handlers.DeleteCategoryHandler)
	}
}