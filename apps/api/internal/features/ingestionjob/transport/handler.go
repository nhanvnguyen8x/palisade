package transport

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/domain"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListPending(c *gin.Context) {
	limit := 10
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			rest.BadRequest(c, "invalid limit")
			return
		}

		limit = parsed
	}

	jobs, err := h.service.ListPending(c.Request.Context(), limit)
	if err != nil {
		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToIngestionJobListResponse(jobs))
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("ingestionJobId"))
	if err != nil {
		rest.BadRequest(c, "invalid ingestion job id")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	job, err := h.service.UpdateStatus(c.Request.Context(), application.UpdateStatusCommand{
		JobID:        id,
		Status:       domain.JobStatus(req.Status),
		ErrorMessage: req.ErrorMessage,
	})
	if err != nil {
		if errors.Is(err, domain.ErrIngestionJobNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		if errors.Is(err, domain.ErrInvalidStatus) {
			rest.BadRequest(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToIngestionJobResponse(job))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("ingestionJobId"))
	if err != nil {
		rest.BadRequest(c, "invalid ingestion job id")
		return
	}

	job, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrIngestionJobNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToIngestionJobResponse(job))
}
