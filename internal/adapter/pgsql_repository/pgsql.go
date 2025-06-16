package pgsqlRepository

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/dsaime/goods-and-projects/internal/domain"
	goodsCache "github.com/dsaime/goods-and-projects/internal/port/goods_cache"

	_ "github.com/lib/pq"
)

type Factory struct {
	db *sqlx.DB
}

type Config struct {
	DSN string
}

func InitFactory(cfg Config) (*Factory, error) {
	conn, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("sqlx.Ping: %w", err)
	}

	slog.Info("Successfully connected to PostgreSQL")

	return &Factory{
		db: conn,
	}, nil
}

func (f *Factory) Close() error {
	return f.db.Close()
}

func (f *Factory) NewGoodsRepository(cache goodsCache.GoodsCache) domain.GoodsRepository {
	return &goodsRepository{
		db:         f.db,
		txBeginner: f.db,
		isTx:       false,
		cache:      cache,
	}
}
