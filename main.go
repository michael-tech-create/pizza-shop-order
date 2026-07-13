package main

import (
	"fmt"
	"log"
	"pizza-app/routes"

	"github.com/gin-gonic/gin"

	"pizza-app/database"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"time"
)

func main() {

	// In local dev, .env supplies config. In Docker/production, env vars
	// are injected directly by the container runtime (docker run -e,
	// --env-file, or your host's secret manager) and no .env file exists
	// — that's expected, not an error, so we only log and continue.
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	database.ConnectDataBase()
	router := gin.Default()

router.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:5500",
        "http://127.0.0.1:5500",
    },
    AllowMethods: []string{
        "GET",
        "POST",
        "PUT",
        "PATCH",
        "DELETE",
        "OPTIONS",
    },
    AllowHeaders: []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
    },
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
}))
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.Static("/public", "./frontend/public")
	router.Static("/uploads", "./uploads")

	router.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("./frontend/login.html")
	})

	router.GET("/admin", func(c *gin.Context) {
		c.File("./frontend/admin.html")
	})

	router.GET("/activity", func(c *gin.Context) {
		c.File("./frontend/activity.html")
	})

	router.GET("/pizza", func(c *gin.Context) {
		c.File("./frontend/pizza.html")
	})

	router.GET("/checkout", func(c *gin.Context) {
		c.File("./frontend/checkout.html")
	})

	routes.SetupRoutes(router)

	fmt.Println("server starting at http://localhost:8080")

	router.Run(":8080")
}
