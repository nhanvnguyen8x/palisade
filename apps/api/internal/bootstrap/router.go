package bootstrap

import (
	"github.com/gin-gonic/gin"
	healthtransport "github.com/nhanvnguyen8x/palisade/internal/features/health/transport"
	ingestiontransport "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/transport"
	chattransport "github.com/nhanvnguyen8x/palisade/internal/features/chat/transport"
	kbtransport "github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/transport"
	kstransport "github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/transport"
	orgtransport "github.com/nhanvnguyen8x/palisade/internal/features/organization/transport"
	wstransport "github.com/nhanvnguyen8x/palisade/internal/features/workspace/transport"
)

type Routes struct {
	Health          *healthtransport.Handler
	Organization    *orgtransport.Handler
	Workspace       *wstransport.Handler
	KnowledgeBase   *kbtransport.Handler
	KnowledgeSource *kstransport.Handler
	IngestionJob    *ingestiontransport.Handler
	Chat            *chattransport.Handler
}

func RegisterRoutes(engine *gin.Engine, routes Routes) {
	engine.Any("/health", routes.Health.GetHealth)

	api := engine.Group("/api/v1")
	orgtransport.RegisterRoutes(api, routes.Organization)
	wstransport.RegisterRoutes(api, routes.Workspace)
	kbtransport.RegisterRoutes(api, routes.KnowledgeBase)
	kstransport.RegisterRoutes(api, routes.KnowledgeSource)
	ingestiontransport.RegisterRoutes(api, routes.IngestionJob)
	chattransport.RegisterRoutes(api, routes.Chat)
}
