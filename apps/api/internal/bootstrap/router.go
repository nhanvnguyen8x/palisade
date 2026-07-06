package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/nhanvnguyen8x/palisade/internal/platform/health/transport"
)

func RegisterRoutes(group *gin.RouterGroup, handler *transport.Handler) {
	group.Any("/health", handler.GetHealth)
}
