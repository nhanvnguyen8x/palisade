package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/features/document/application"
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

func (h *Handler) Upload(
	c *gin.Context,
) {

	req, err := ParseUploadRequest(c)
	if err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		rest.BadRequest(c, err.Error())
		return
	}

	reader, err := file.Open()
	if err != nil {
		rest.InternalServerError(c)
		return
	}

	defer reader.Close()

	result, err := h.service.Upload(
		c.Request.Context(),
		application.UploadCommand{

			WorkspaceID: req.WorkspaceID,

			FileName: req.FileName,

			ContentType: req.ContentType,

			Size: req.Size,

			Reader: reader,
		},
	)

	if err != nil {
		rest.InternalServerError(c)
		return
	}

	rest.Accepted(
		c,
		ToUploadResponse(result),
	)
}
