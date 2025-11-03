package repositories

import (
	"context"
	"fmt"
	"time"

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
		result := tx.Exec(
			"UPDATE products SET stock = stock - ? WHERE id = ? AND stock >= ?",
			quantity, productID, quantity,
		)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("OUT_OF_STOCK")
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

func (r *DBRepository) CreateJob(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *DBRepository) UpdateJobProgress(ctx context.Context, jobID string, processed, total int) error {
	progress := 0
	if total > 0 {
		progress = (processed * 100) / total
	}

	return r.db.WithContext(ctx).Model(&models.Job{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"progress":   progress,
			"processed":  processed,
			"total":      total,
			"updated_at": time.Now(),
		}).Error
}

func (r *DBRepository) UpdateJobResult(ctx context.Context, jobID, resultPath string) error {
	return r.db.WithContext(ctx).Model(&models.Job{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":      "completed",
			"result_path": resultPath,
			"updated_at":  time.Now(),
		}).Error
}

func (r *DBRepository) GetJob(ctx context.Context, jobID string) (*models.Job, error) {
	var job models.Job
	if err := r.db.WithContext(ctx).First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *DBRepository) CancelJob(ctx context.Context, jobID string) error {
	return r.db.WithContext(ctx).Model(&models.Job{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"status":     "cancelled",
			"updated_at": time.Now(),
		}).Error
}

func (r *DBRepository) GetTotalTransactions(ctx context.Context, from, to time.Time) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("paid_at BETWEEN ? AND ? AND status = ?", from, to, "paid").
		Count(&count).Error
	return int(count), err
}

func (r *DBRepository) GetTransactionsBatch(ctx context.Context, from, to time.Time, offset, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.WithContext(ctx).
		Where("paid_at BETWEEN ? AND ? AND status = ?", from, to, "paid").
		Order("id").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (r *DBRepository) UpsertSettlement(ctx context.Context, settlement *models.Settlement) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "merchant_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"gross_cents", "fee_cents", "net_cents", "txn_count", "generated_at", "unique_run_id"}),
	}).Create(settlement).Error
}

// unting testing
func (r *DBRepository) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
