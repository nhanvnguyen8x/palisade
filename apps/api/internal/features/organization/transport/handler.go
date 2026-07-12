package transport

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/features/organization/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/organization/domain"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	org, err := h.service.Create(c.Request.Context(), application.CreateCommand{
		Name: req.Name,
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidName) {
			rest.BadRequest(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToOrganizationResponse(org))
}

func (h *Handler) GetByID(c *gin.Context) {
	org, err := h.service.GetByID(c.Request.Context(), c.Param("organizationId"))
	if err != nil {
		if errors.Is(err, domain.ErrOrganizationNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToOrganizationResponse(org))
}
