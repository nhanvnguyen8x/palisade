package transport

import "github.com/gin-gonic/gin"

func RegisterRoutes(group *gin.RouterGroup, handler *Handler) {
	group.GET("/ingestion-jobs/pending", handler.ListPending)
	group.GET("/ingestion-jobs/:ingestionJobId", handler.GetByID)
	group.PATCH("/ingestion-jobs/:ingestionJobId", handler.UpdateStatus)
}
