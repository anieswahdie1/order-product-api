package services

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/anieswahdie1/order-product-api.git/internal/models"
	"github.com/anieswahdie1/order-product-api.git/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConcurrentOrders(t *testing.T) {
	dsn := "host=localhost user=postgres password=5396 dbname=order_product port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	// BERSIHKAN DATA LAMA + RESET STOCK
	db.Exec("DELETE FROM orders")
	db.Model(&models.Product{}).Where("id = ?", 1).Update("stock", 100)

	// Verify initial state
	var initialProduct models.Product
	db.First(&initialProduct, 1)
	t.Logf("Initial stock: %d", initialProduct.Stock)

	repo := repositories.NewDBRepository(db)
	service := NewOrderService(repo)

	// Test concurrent orders
	concurrentUsers := 500
	successfulOrders := 0
	failedOrders := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	ctx := context.Background()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			order, err := service.CreateOrder(ctx, models.OrderRequest{
				ProductID: 1,
				Quantity:  1,
				BuyerID:   fmt.Sprintf("user-%d", userID),
			})

			if err == nil && order != nil {
				mu.Lock()
				successfulOrders++
				mu.Unlock()
			} else {
				mu.Lock()
				failedOrders++
				mu.Unlock()
				if err.Error() != "OUT_OF_STOCK" {
					t.Logf("Unexpected error for user %d: %v", userID, err)
				}
			}
		}(i)
	}

	wg.Wait()

	var product models.Product
	db.First(&product, 1)

	var orderCount int64
	db.Model(&models.Order{}).Where("product_id = ?", 1).Count(&orderCount)

	t.Logf("=== RESULTS ===")
	t.Logf("Successful orders: %d", successfulOrders)
	t.Logf("Failed orders: %d", failedOrders)
	t.Logf("Final stock: %d", product.Stock)
	t.Logf("Total orders in database: %d", orderCount)

	if successfulOrders != 100 {
		t.Errorf("Expected exactly 100 successful orders, got %d", successfulOrders)
	}
	if product.Stock != 0 {
		t.Errorf("Expected final stock to be 0, got %d", product.Stock)
	}
	if orderCount != 100 {
		t.Errorf("Expected 100 orders in database, got %d", orderCount)
	}
}
