package transport

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/domain"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		rest.BadRequest(c, "invalid workspace id")
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	kb, err := h.service.Create(c.Request.Context(), application.CreateCommand{
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidName) || errors.Is(err, domain.ErrWorkspaceNotFound) {
			rest.BadRequest(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToKnowledgeBaseResponse(kb))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("knowledgeBaseId"))
	if err != nil {
		rest.BadRequest(c, "invalid knowledge base id")
		return
	}

	kb, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrKnowledgeBaseNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToKnowledgeBaseResponse(kb))
}

func (h *Handler) ListByWorkspace(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		rest.BadRequest(c, "invalid workspace id")
		return
	}

	bases, err := h.service.ListByWorkspace(c.Request.Context(), workspaceID)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToKnowledgeBaseListResponse(bases))
}
