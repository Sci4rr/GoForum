package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var router *Router

func initDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	
	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	return db
}

func main() {
	db := initDB()
	defer db.Close()

	router = InitializeRoutes(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}