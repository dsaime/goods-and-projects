package pgsql

import (
	"github.com/jmoiron/sqlx"

	"github.com/dsaime/goods-and-projects/internal/domain"

	_ "github.com/lib/pq"
)

type Factory struct {
	db *sqlx.DB
}

func InitFactory(cfg Config) (*Factory, error) {
	conn, err := sqlx.Connect("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &Factory{
		db: conn,
	}, nil
}

func (f *Factory) Close() error {
	return f.db.Close()
}

func (f *Factory) NewGoodsRepository(cache GoodsCache) domain.GoodsRepository {
	return &goodsRepository{
		db:         f.db,
		txBeginner: f.db,
		isTx:       false,
		cache:      cache,
	}
}
