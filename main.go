package main

import (
	"pizza-app/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"pizza-app/database"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading my .env file")
	}

	database.ConnectDataBase()
	
	router := gin.Default()

	router.Use(cors.Default())

	routes.SetupRoutes(router)
	fmt.Println("server starting at http://localhost:8080")
	router.Run(":8080")
}