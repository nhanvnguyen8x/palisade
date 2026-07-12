package transport

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadRequest struct {
	WorkspaceID uuid.UUID
	FileName    string
	ContentType string
	Size        int64
}

func ParseUploadRequest(c *gin.Context) (*UploadRequest, error) {

	workspaceID, err := uuid.Parse(c.Param("workspaceId"))
	if err != nil {
		return nil, fmt.Errorf("invalid workspace id")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("file is required")
	}

	return &UploadRequest{
		WorkspaceID: workspaceID,
		FileName:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
	}, nil
}
