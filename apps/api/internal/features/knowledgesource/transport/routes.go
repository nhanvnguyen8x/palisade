package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/knowledge-bases/:knowledgeBaseId/sources", handler.Upload)
	group.GET("/knowledge-sources/:knowledgeSourceId", handler.GetByID)
}
