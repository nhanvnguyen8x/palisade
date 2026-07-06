package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/features/health/application"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetHealth(c *gin.Context) {
	ctx := c.Request.Context()
	report, err := h.service.GetHealth(ctx)
	if err != nil {
		rest.InternalServerError(c)
		return
	}

	response := ToHealthResponse(report)
	rest.OK(c, response)
}
