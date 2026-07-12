package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/organizations/:organizationId/workspaces", handler.Create)
	group.GET("/workspaces/:workspaceId", handler.GetByID)
}
