package transport

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/application"
	"github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/domain"
	kbdomain "github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/domain"
	"github.com/nhanvnguyen8x/palisade/internal/shared/rest"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(c *gin.Context) {
	knowledgeBaseID, err := uuid.Parse(c.Param("knowledgeBaseId"))
	if err != nil {
		rest.BadRequest(c, "invalid knowledge base id")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		rest.BadRequest(c, "file is required")
		return
	}

	reader, err := file.Open()
	if err != nil {
		rest.InternalServerError(c)
		return
	}
	defer reader.Close()

	result, err := h.service.Upload(c.Request.Context(), application.UploadCommand{
		KnowledgeBaseID: knowledgeBaseID,
		FileName:        file.Filename,
		ContentType:     file.Header.Get("Content-Type"),
		Size:            file.Size,
		Reader:          reader,
	})
	if err != nil {
		if errors.Is(err, domain.ErrKnowledgeBaseNotFound) || errors.Is(err, kbdomain.ErrKnowledgeBaseNotFound) {
			rest.NotFound(c, "knowledge base not found")
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.Accepted(c, ToUploadResponse(result))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("knowledgeSourceId"))
	if err != nil {
		rest.BadRequest(c, "invalid knowledge source id")
		return
	}

	source, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrKnowledgeSourceNotFound) {
			rest.NotFound(c, err.Error())
			return
		}

		rest.InternalServerError(c)
		return
	}

	rest.OK(c, ToKnowledgeSourceResponse(source))
}

func ParseUploadRequest(c *gin.Context) (*UploadRequest, error) {
	knowledgeBaseID, err := uuid.Parse(c.Param("knowledgeBaseId"))
	if err != nil {
		return nil, fmt.Errorf("invalid knowledge base id")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("file is required")
	}

	return &UploadRequest{
		KnowledgeBaseID: knowledgeBaseID,
		FileName:        file.Filename,
		ContentType:     file.Header.Get("Content-Type"),
		Size:            file.Size,
	}, nil
}
