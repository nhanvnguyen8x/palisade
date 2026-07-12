package transport

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/workspace/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/workspace/domain"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		rest.BadRequest(c, "invalid organization id")
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	workspace, err := h.service.Create(c.Request.Context(), application.CreateCommand{
		OrganizationID: organizationID,
		Name:           req.Name,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidName) || errors.Is(err, domain.ErrOrganizationNotFound) {
			rest.BadRequest(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToWorkspaceResponse(workspace))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		rest.BadRequest(c, "invalid workspace id")
		return
	}

	workspace, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrWorkspaceNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToWorkspaceResponse(workspace))
}
