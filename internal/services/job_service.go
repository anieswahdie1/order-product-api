package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/anieswahdie1/order-product-api.git/internal/jobs"
	"github.com/anieswahdie1/order-product-api.git/internal/models"
	"github.com/anieswahdie1/order-product-api.git/internal/repositories"
	"github.com/google/uuid"
)

type JobService struct {
	repo       *repositories.DBRepository
	workerPool *jobs.WorkerPool
}

func NewJobService(repo *repositories.DBRepository, workerPool *jobs.WorkerPool) *JobService {
	return &JobService{
		repo:       repo,
		workerPool: workerPool,
	}
}

func (s *JobService) CreateSettlementJob(ctx context.Context, req models.JobRequest) (*models.Job, error) {
	from, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		return nil, err
	}

	to, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		return nil, err
	}

	jobID := "job_" + uuid.New().String()[:8]

	metadata, _ := json.Marshal(map[string]interface{}{
		"from": req.From,
		"to":   req.To,
	})

	job := &models.Job{
		ID:        jobID,
		Type:      "settlement",
		Status:    "QUEUED",
		Progress:  0,
		Processed: 0,
		Total:     0,
		Metadata:  string(metadata),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}

	s.workerPool.SubmitJob(jobs.Job{
		ID:   jobID,
		Type: "settlement",
		From: from,
		To:   to,
	})

	return job, nil
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (*models.Job, error) {
	job, err := s.repo.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *JobService) CancelJob(ctx context.Context, jobID string) error {
	s.workerPool.CancelJob(jobID)
	return s.repo.CancelJob(ctx, jobID)
}
