package main

import (
	"database/sql"
	"log"

	"github.com/anieswahdie1/order-product-api.git/internal/handlers"
	"github.com/anieswahdie1/order-product-api.git/internal/repositories"
	"github.com/anieswahdie1/order-product-api.git/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:5396@localhost:5432/order_product?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to database successfully")

	repo := repositories.NewDBRepository(db)

	orderService := services.NewOrderService(repo)

	orderHandler := handlers.NewOrderHandler(orderService)

	r := gin.Default()

	r.POST("/orders", orderHandler.CreateOrder)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))

}
