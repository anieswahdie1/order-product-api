package models

import "time"

type OrderRequest struct {
	ProductID int    `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	BuyerID   string `json:"buyer_id" binding:"required"`
}

type Order struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	BuyerID   string    `json:"buyer_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type Transaction struct {
	ID          int       `json:"id"`
	MerchantID  int       `json:"merchant_id"`
	AmountCents int       `json:"amount_cents"`
	FeeCents    int       `json:"fee_cents"`
	Status      string    `json:"status"`
	PaidAt      time.Time `json:"paid_at"`
}

type Settlement struct {
	ID          int       `json:"id"`
	MerchantID  int       `json:"merchant_id"`
	Date        time.Time `json:"date"`
	GrossCents  int       `json:"gross_cents"`
	FeeCents    int       `json:"fee_cents"`
	NetCents    int       `json:"net_cents"`
	TxnCount    int       `json:"txn_count"`
	GeneratedAt time.Time `json:"generated_at"`
	UniqueRunID string    `json:"unique_run_id"`
}

type JobRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type Job struct {
	ID         string         `json:"job_id"`
	Type       string         `json:"type"`
	Status     string         `json:"status"`
	Progress   int            `json:"progress"`
	Processed  int            `json:"processed"`
	Total      int            `json:"total"`
	ResultPath string         `json:"download_url,omitempty"`
	Metadata   map[string]any `json:"metadata"`
}
