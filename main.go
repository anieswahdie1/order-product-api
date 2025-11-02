package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anieswahdie1/order-product-api.git/internal/handlers"
	"github.com/anieswahdie1/order-product-api.git/internal/jobs"
	"github.com/anieswahdie1/order-product-api.git/internal/models"
	"github.com/anieswahdie1/order-product-api.git/internal/repositories"
	"github.com/anieswahdie1/order-product-api.git/internal/services"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "host=localhost user=postgres password=5396 dbname=order_product port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.Product{},
		&models.Order{},
		&models.Transaction{},
		&models.Settlement{},
		&models.Job{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully")

	repo := repositories.NewDBRepository(db)

	workers := 4
	if envWorkers := os.Getenv("WORKERS"); envWorkers != "" {
		fmt.Sscanf(envWorkers, "%d", &workers)
	}

	workerPool := jobs.NewWorkerPool(repo, workers)
	workerPool.Start()

	orderService := services.NewOrderService(repo)
	jobService := services.NewJobService(repo, workerPool)

	orderHandler := handlers.NewOrderHandler(orderService)
	jobHandler := handlers.NewJobHandler(jobService)

	r := gin.Default()

	r.POST("/orders", orderHandler.CreateOrder)
	r.GET("/orders/:id", orderHandler.GetOrder)

	r.POST("/jobs/settlement", jobHandler.CreateSettlementJob)
	r.GET("/jobs/:id", jobHandler.GetJob)
	r.POST("/jobs/:id/cancel", jobHandler.CancelJob)
	r.GET("/downloads/:id.csv", jobHandler.DownloadResult)

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))

}
