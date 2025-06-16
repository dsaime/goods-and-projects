package app

import (
	"errors"
	"fmt"
	"log/slog"

	redisGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/redis_goods_cache"

	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

type adapters struct {
	goodsCache pgsql.GoodsCache
}

func (a *adapters) GoodsCache() pgsql.GoodsCache {
	return a.goodsCache
}

func initAdapters(config Config) (*adapters, func(), error) {
	goodsCache, err := redisGoodsCache.Init(config.Redis)
	if err != nil {
		return nil, nil, fmt.Errorf("redisGoodsCache.Init: %w", err)
	}

	closer := func() {
		if err := errors.Join(goodsCache.Close()); err != nil {
			slog.Error("initAdapters: close: " + err.Error())
		}
	}

	return &adapters{
		goodsCache: goodsCache,
	}, closer, nil
}
