package main

// One-off CLI seed script for creating admin users.
// This is intentionally kept OUT of routes.go — admin self-signup is
// out of scope for a single-restaurant dashboard, so this is the only
// way to create/reset an admin account.
//
// Usage (run from your project root, next to your existing main.go / .env):
//
//   go run ./cmd/seed -username=admin -password=yourpassword
//
// Suggested file location: cmd/seed/main.go
// (Go allows multiple main packages in a repo as long as each lives in
// its own directory — this won't conflict with your existing main.go.)

import (
	"flag"
	"fmt"
	"log"

	"pizza-app/database"
	"pizza-app/repositories"

	"github.com/joho/godotenv"
)

func main() {
	username := flag.String("username", "", "admin username to create")
	password := flag.String("password", "", "admin password (plaintext, will be hashed)")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("both -username and -password are required, e.g. go run ./cmd/seed -username=admin -password=yourpassword")
	}

	if len(*password) < 8 {
		log.Fatal("password should be at least 8 characters")
	}

	// Load the same .env your main server uses, so DB_HOST/DB_USER/etc.
	// match. If your .env lives elsewhere relative to this file, adjust
	// the path here, e.g. godotenv.Load("../../.env").
	if err := godotenv.Load(); err != nil {
		log.Fatal("could not load .env file: ", err)
	}

	database.ConnectDataBase()

	// Guard against accidentally creating a duplicate. GetAdminByUsername
	// returns ErrInvalidCredentials when not found, so that specific error
	// is the "safe to proceed" case — anything else is a real DB problem.
	if _, err := repositories.GetAdminByUsername(*username); err == nil {
		log.Fatalf("admin '%s' already exists — pick a different username or delete the existing row first", *username)
	} else if err != repositories.ErrInvalidCredentials {
		log.Fatal("error checking for existing admin: ", err)
	}

	if err := repositories.CreateAdmin(*username, *password); err != nil {
		log.Fatal("failed to create admin: ", err)
	}

	fmt.Printf("✅ admin '%s' created successfully\n", *username)
}
