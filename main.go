package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/order_product?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	r := gin.Default()

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))

}
