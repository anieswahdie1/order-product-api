package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderRequest struct {
	ProductID int    `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	BuyerID   string `json:"buyer_id" binding:"required"`
}

type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID int            `json:"product_id"`
	Quantity  int            `json:"quantity"`
	BuyerID   string         `json:"buyer_id"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Product struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Stock     int            `json:"stock"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Transaction struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MerchantID  int            `json:"merchant_id"`
	AmountCents int            `json:"amount_cents"`
	FeeCents    int            `json:"fee_cents"`
	Status      string         `json:"status"`
	PaidAt      time.Time      `json:"paid_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Settlement struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MerchantID  int            `json:"merchant_id"`
	Date        time.Time      `json:"date"`
	GrossCents  int            `json:"gross_cents"`
	FeeCents    int            `json:"fee_cents"`
	NetCents    int            `json:"net_cents"`
	TxnCount    int            `json:"txn_count"`
	GeneratedAt time.Time      `json:"generated_at"`
	UniqueRunID string         `json:"unique_run_id"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Job struct {
	ID         string         `json:"job_id" gorm:"primaryKey"`
	Type       string         `json:"type"`
	Status     string         `json:"status"`
	Progress   int            `json:"progress"`
	Processed  int            `json:"processed"`
	Total      int            `json:"total"`
	ResultPath string         `json:"result_path"`
	Metadata   string         `json:"metadata" gorm:"type:jsonb"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type JobRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}
