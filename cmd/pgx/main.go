package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
	// Connect to the database
	connStr := os.Getenv("DATABASE_CONNECTION_STRING")
	// Connect to Database
	dbPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()
}
