package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"testing"
)

func TestConnectToPostgres(t *testing.T) {
	connStr := "postgres://postgres:*******localhost:5432/players" //Add your secret password
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to PostgreSQL")
}
