package transport

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/platform/health/application"
)

type Handler struct {
	service application.Service
}

func NewHandler(service application.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetHealth(c *gin.Context) {
	var ctx = c.Request.Context()
	report, err := h.service.GetHealth(ctx)
	if err != nil {
		log.Println(err)
	}

	response := ToHealthResponse(report)
	http.OK(c, response)
}
