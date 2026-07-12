package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhanvnguyen8x/palisade/internal/config"
	"github.com/nhanvnguyen8x/palisade/internal/platform/database/postgres"
)

func NewDatabase(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	return postgres.NewPool(ctx, cfg)
}
