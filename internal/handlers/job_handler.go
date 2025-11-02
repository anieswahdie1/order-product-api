package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/anieswahdie1/order-product-api.git/internal/services"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

func (h *JobHandler) CreateSettlementJob(c *gin.Context) {
	var req struct {
		From string `json:"from" binding:"required"`
		To   string `json:"to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobService.CreateSettlementJob(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"job_id": job.ID,
		"status": job.Status,
	}

	c.JSON(http.StatusAccepted, response)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	jobID := c.Param("id")

	job, err := h.jobService.GetJob(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	response := gin.H{
		"job_id":    job.ID,
		"status":    job.Status,
		"progress":  job.Progress,
		"processed": job.Processed,
		"total":     job.Total,
	}

	if job.ResultPath != "" {
		response["download_url"] = "/downloads/" + jobID + ".csv"
	}

	c.JSON(http.StatusOK, response)
}

func (h *JobHandler) CancelJob(c *gin.Context) {
	jobID := c.Param("id")

	if err := h.jobService.CancelJob(c.Request.Context(), jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cancellation requested"})
}

func (h *JobHandler) DownloadResult(c *gin.Context) {
	jobID := c.Param("id")
	filePath := filepath.Join("/tmp/settlements", jobID+".csv")

	c.File(filePath)
}
