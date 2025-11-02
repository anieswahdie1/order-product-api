package services

import (
	"context"

	"github.com/anieswahdie1/order-product-api.git/internal/models"
	"github.com/anieswahdie1/order-product-api.git/internal/repositories"
)

type OrderService struct {
	repo *repositories.DBRepository
}

func NewOrderService(repo *repositories.DBRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, req models.OrderRequest) (*models.Order, error) {
	return s.repo.CreateOrder(ctx, req.ProductID, req.Quantity, req.BuyerID)
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	return s.repo.GetOrder(ctx, id)
}
