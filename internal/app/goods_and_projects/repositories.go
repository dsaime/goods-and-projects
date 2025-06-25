package app

import (
	"fmt"
	"log/slog"

	pgsqlRepository "github.com/dsaime/goods-and-projects/internal/adapter/pgsql_repository"
	"github.com/dsaime/goods-and-projects/internal/domain"
	goodsCache "github.com/dsaime/goods-and-projects/internal/port/goods_cache"
)

type repositories struct {
	goods domain.GoodsRepository
}

type pgsqlDeps interface {
	GoodsCache() goodsCache.GoodsCache
}

func initPgsqlRepositories(config Config, deps pgsqlDeps) (*repositories, func(), error) {
	factory, err := pgsqlRepository.InitFactory(config.Pgsql)
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
