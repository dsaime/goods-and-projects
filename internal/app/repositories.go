package app

import (
	"fmt"
	"log/slog"

	"github.com/dsaime/goods-and-projects/internal/domain"
	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

type repositories struct {
	goods domain.GoodsRepository
}

type pgsqlDeps interface {
	GoodsCache() pgsql.GoodsCache
}

func initPgsqlRepositories(config Config, deps pgsqlDeps) (*repositories, func(), error) {
	factory, err := pgsql.InitFactory(config.Pgsql)
	if err != nil {
		return nil, nil, fmt.Errorf("pgsql.InitFactory: %w", err)
	}

	rs := &repositories{
		goods: factory.NewGoodsRepository(deps.GoodsCache()),
	}

	return rs, func() {
		if err := factory.Close(); err != nil {
			slog.Error("initPgsqlRepositories: factory.Close: " + err.Error())
		}
	}, nil
}
