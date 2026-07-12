package transport

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/features/chat/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/chat/domain"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	answer, err := h.service.Chat(c.Request.Context(), application.ChatRequest{
		KnowledgeBaseID: req.KnowledgeBaseID,
		Question:        req.Question,
		TopK:            req.TopK,
	})
	if err != nil {
		if errors.Is(err, domain.ErrKnowledgeBaseNotFound) || errors.Is(err, domain.ErrNoRelevantChunks) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToChatResponse(answer))
}

