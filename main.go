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
// 	"pizza-app/repositories"
// )
)

func main() {

err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDataBase()
	router := gin.Default()

router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5500", "http://127.0.0.1:5500"}, // wherever your frontend is served from
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
	router.Static("/uploads", "./uploads")

	routes.SetupRoutes(router)

	fmt.Println("server starting at http://localhost:8080")

	router.Run(":8080")
}
