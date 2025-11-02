package repositories

import (
	"context"
	"fmt"

	"github.com/anieswahdie1/order-product-api.git/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) CreateOrder(ctx context.Context, productID, quantity int, buyerID string) (*models.Order, error) {
	var order *models.Order

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&product, productID).Error; err != nil {
			return err
		}

		if product.Stock < quantity {
			return fmt.Errorf("OUT_OF_STOCK")
		}

		if err := tx.Model(&product).
			Update("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
			return err
		}

		order = &models.Order{
			ProductID: productID,
			Quantity:  quantity,
			BuyerID:   buyerID,
			Status:    "created",
		}

		return tx.Create(order).Error
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *DBRepository) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
