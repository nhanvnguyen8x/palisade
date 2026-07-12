package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/workspaces/:workspaceId/documents", handler.Upload)
}
