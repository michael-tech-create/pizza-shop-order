package database

import (
	"database/sql"

	"log"
	"fmt"
	"os"

	_ "github.com/lib/pq" 
	
)

var DB *sql.DB

func ConnectDataBase() {


	connectStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connectStr)


	if err != nil {
		log.Fatal("failed to connect db", err)
	}


	err = db.Ping()

	if err != nil {
		log.Fatal("Database is offline or unreachable" , err)
	}

	log.Println("database connected successfully")


	DB = db

}