package database

import (
	"database/sql"
	"time"

	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDataBase() {


	connectStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
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

	// log.Println("database connected successfully")

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)


	DB = db

}