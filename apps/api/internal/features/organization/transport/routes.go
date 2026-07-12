package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.POST("/organizations", handler.Create)
	group.GET("/organizations/:organizationId", handler.GetByID)
}
