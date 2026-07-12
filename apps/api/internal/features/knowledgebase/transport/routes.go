package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/workspaces/:workspaceId/knowledge-bases", handler.Create)
	group.GET("/workspaces/:workspaceId/knowledge-bases", handler.ListByWorkspace)
	group.GET("/knowledge-bases/:knowledgeBaseId", handler.GetByID)
}
