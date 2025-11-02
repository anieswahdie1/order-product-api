package jobs

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anieswahdie1/order-product-api.git/internal/models"
	"github.com/anieswahdie1/order-product-api.git/internal/repositories"
)

type Job struct {
	ID     string
	Type   string
	From   time.Time
	To     time.Time
	Cancel context.CancelFunc
}

type WorkerPool struct {
	repo       *repositories.DBRepository
	jobQueue   chan Job
	workers    int
	activeJobs sync.Map
}

func NewWorkerPool(repo *repositories.DBRepository, workers int) *WorkerPool {
	return &WorkerPool{
		repo:     repo,
		jobQueue: make(chan Job, 100),
		workers:  workers,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)

	for job := range wp.jobQueue {
		log.Printf("Worker %d processing job %s", id, job.ID)
		wp.processSettlementJob(job)
	}
}

func (wp *WorkerPool) SubmitJob(job Job) {
	wp.activeJobs.Store(job.ID, job)
	wp.jobQueue <- job
}

func (wp *WorkerPool) CancelJob(jobID string) {
	if job, ok := wp.activeJobs.Load(jobID); ok {
		if j, ok := job.(Job); ok && j.Cancel != nil {
			j.Cancel()
		}
		wp.activeJobs.Delete(jobID)
	}
}

func (wp *WorkerPool) processSettlementJob(job Job) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if existing, ok := wp.activeJobs.Load(job.ID); ok {
		if j, ok := existing.(Job); ok {
			j.Cancel = cancel
			wp.activeJobs.Store(job.ID, j)
		}
	}

	csvPath := filepath.Join("/tmp/settlements", job.ID+".csv")
	if err := os.MkdirAll(filepath.Dir(csvPath), 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		return
	}

	file, err := os.Create(csvPath)
	if err != nil {
		log.Printf("Error creating CSV file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"merchant_id", "date", "gross", "fee", "net", "txn_count"})

	total, err := wp.repo.GetTotalTransactions(ctx, job.From, job.To)
	if err != nil {
		log.Printf("Error getting total transactions: %v", err)
		return
	}

	wp.repo.UpdateJobProgress(ctx, job.ID, 0, total)

	batchSize := 10000
	offset := 0
	processed := 0

	settlements := make(map[string]*models.Settlement)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Job %s cancelled", job.ID)
			wp.repo.CancelJob(ctx, job.ID)
			wp.activeJobs.Delete(job.ID)
			return
		default:
		}

		transactions, err := wp.repo.GetTransactionsBatch(ctx, job.From, job.To, offset, batchSize)
		if err != nil {
			log.Printf("Error getting transactions batch: %v", err)
			break
		}

		if len(transactions) == 0 {
			break
		}

		for _, txn := range transactions {
			dateKey := fmt.Sprintf("%d-%s", txn.MerchantID, txn.PaidAt.Format("2006-01-02"))

			if settlement, exists := settlements[dateKey]; exists {
				settlement.GrossCents += txn.AmountCents
				settlement.FeeCents += txn.FeeCents
				settlement.NetCents += (txn.AmountCents - txn.FeeCents)
				settlement.TxnCount++
			} else {
				settlements[dateKey] = &models.Settlement{
					MerchantID:  txn.MerchantID,
					Date:        time.Date(txn.PaidAt.Year(), txn.PaidAt.Month(), txn.PaidAt.Day(), 0, 0, 0, 0, time.UTC),
					GrossCents:  txn.AmountCents,
					FeeCents:    txn.FeeCents,
					NetCents:    txn.AmountCents - txn.FeeCents,
					TxnCount:    1,
					UniqueRunID: job.ID,
				}
			}
		}

		processed += len(transactions)
		offset += batchSize

		wp.repo.UpdateJobProgress(ctx, job.ID, processed, total)

		log.Printf("Job %s processed %d/%d transactions", job.ID, processed, total)
	}

	for _, settlement := range settlements {
		if err := wp.repo.UpsertSettlement(ctx, settlement); err != nil {
			log.Printf("Error upserting settlement: %v", err)
			continue
		}

		writer.Write([]string{
			fmt.Sprintf("%d", settlement.MerchantID),
			settlement.Date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", float64(settlement.GrossCents)/100),
			fmt.Sprintf("%.2f", float64(settlement.FeeCents)/100),
			fmt.Sprintf("%.2f", float64(settlement.NetCents)/100),
			fmt.Sprintf("%d", settlement.TxnCount),
		})
	}

	wp.repo.UpdateJobResult(ctx, job.ID, csvPath)
	wp.activeJobs.Delete(job.ID)

	log.Printf("Job %s completed", job.ID)
}
