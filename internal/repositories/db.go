package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/anieswahdie1/order-product-api.git/internal/models"
)

type DBRepository struct {
	db *sql.DB
}

func NewDBRepository(db *sql.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) CreateOrder(ctx context.Context, productID, quantity int, buyerID string) (*models.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var currentStock int
	err = tx.QueryRowContext(ctx,
		"SELECT stock FROM products WHERE id = $1 FOR UPDATE",
		productID).Scan(&currentStock)
	if err != nil {
		return nil, err
	}

	if currentStock < quantity {
		return nil, fmt.Errorf("OUT_OF_STOCK")
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE products SET stock = stock - $1 WHERE id = $2",
		quantity, productID)
	if err != nil {
		return nil, err
	}

	var order models.Order
	err = tx.QueryRowContext(ctx,
		`INSERT INTO orders (product_id, quantity, buyer_id, status) VALUES ($1, $2, $3, 'created') RETURNING id, product_id, quantity, buyer_id, status, created_at`,
		productID, quantity, buyerID).Scan(
		&order.ID, &order.ProductID, &order.Quantity, &order.BuyerID,
		&order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}
