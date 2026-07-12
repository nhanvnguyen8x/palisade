package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	chatapp "github.com/nhanvnguyen8x/palisade/internal/features/chat/application"
	chattransport "github.com/nhanvnguyen8x/palisade/internal/features/chat/transport"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhanvnguyen8x/palisade/internal/config"
	ingestionapp "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/application"
	ingestiontransport "github.com/nhanvnguyen8x/palisade/internal/features/ingestionjob/transport"
	kbapp "github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/application"
	kbtransport "github.com/nhanvnguyen8x/palisade/internal/features/knowledgebase/transport"
	ksapp "github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/application"
	kstransport "github.com/nhanvnguyen8x/palisade/internal/features/knowledgesource/transport"
	orgapp "github.com/nhanvnguyen8x/palisade/internal/features/organization/application"
	orgtransport "github.com/nhanvnguyen8x/palisade/internal/features/organization/transport"
	healthapp "github.com/nhanvnguyen8x/palisade/internal/features/health/application"
	healthtransport "github.com/nhanvnguyen8x/palisade/internal/features/health/transport"
	wsapp "github.com/nhanvnguyen8x/palisade/internal/features/workspace/application"
	wstransport "github.com/nhanvnguyen8x/palisade/internal/features/workspace/transport"
	embedding "github.com/nhanvnguyen8x/palisade/internal/platform/ai/embedding"
	"github.com/nhanvnguyen8x/palisade/internal/platform/database/postgres"
	eventlog "github.com/nhanvnguyen8x/palisade/internal/platform/event/log"
)

type Application struct {
	config *config.Config
	engine *gin.Engine
	pool   *pgxpool.Pool
}

func NewApplication(cfg *config.Config) (*Application, error) {
	ctx := context.Background()

	pool, err := NewDatabase(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	storage, err := NewObjectStorage(ctx, cfg.Storage)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("storage: %w", err)
	}

	orgRepo := postgres.NewOrganizationRepository(pool)
	wsRepo := postgres.NewWorkspaceRepository(pool)
	kbRepo := postgres.NewKnowledgeBaseRepository(pool)
	ksRepo := postgres.NewKnowledgeSourceRepository(pool)
	jobRepo := postgres.NewIngestionJobRepository(pool)
	publisher := eventlog.NewPublisher()

	orgService := orgapp.NewService(orgRepo)
	orgHandler := orgtransport.NewHandler(orgService)

	wsService := wsapp.NewService(wsRepo, orgRepo)
	wsHandler := wstransport.NewHandler(wsService)

	kbService := kbapp.NewService(kbRepo, wsRepo)
	kbHandler := kbtransport.NewHandler(kbService)

	ksService := ksapp.NewService(pool, kbRepo, ksRepo, storage, publisher)
	ksHandler := kstransport.NewHandler(ksService)

	jobService := ingestionapp.NewService(jobRepo, ksRepo)
	jobHandler := ingestiontransport.NewHandler(jobService)

	chatRepo := postgres.NewChatRepository(pool)
	chatService := chatapp.NewService(chatRepo, embedding.DefaultEmbeddingDim, "hash-embedding")
	chatHandler := chattransport.NewHandler(chatService)

	healthService := healthapp.NewService()
	healthHandler := healthtransport.NewHandler(healthService)

	engine := gin.New()
	engine.Use(gin.Recovery())

	RegisterRoutes(engine, Routes{
		Health:         healthHandler,
		Organization:   orgHandler,
		Workspace:      wsHandler,
		KnowledgeBase:  kbHandler,
		KnowledgeSource: ksHandler,
		IngestionJob:   jobHandler,
		Chat:           chatHandler,
	})

	return &Application{
		config: cfg,
		engine: engine,
		pool:   pool,
	}, nil
}

func (a *Application) Run() error {
	addr := fmt.Sprintf(":%d", a.config.HTTP.Port)
	log.Printf("server listening on %s", addr)
	return a.engine.Run(addr)
}

func (a *Application) Close() {
	if a.pool != nil {
		a.pool.Close()
	}
}
