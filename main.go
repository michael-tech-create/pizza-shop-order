package main

import (
	"fmt"
	"log"
	"pizza-app/routes"

	"github.com/gin-gonic/gin"

	"pizza-app/database"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading my .env file")
	}

	database.ConnectDataBase()

	router := gin.Default()

	router.Use(cors.Default())
	
	router.Static("/uploads", "./uploads")

	routes.SetupRoutes(router)
	fmt.Println("server starting at http://localhost:8080")
	router.Run(":8080")
}
