package app

import (
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"

	natsGoodsEvent "github.com/dsaime/goods-and-projects/internal/adapter/nats_goods_event"
	redisGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/redis_goods_cache"
	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
	goodsCache "github.com/dsaime/goods-and-projects/internal/port/goods_cache"
)

type adapters struct {
	goodsCache       goodsCache.GoodsCache
	goodsEventLogger goodsEvent.Logger
}

func (a *adapters) GoodsCache() goodsCache.GoodsCache {
	return a.goodsCache
}

func initAdapters(config Config) (*adapters, func(), error) {
	cache, err := redisGoodsCache.Init(config.Redis)
	if err != nil {
		return nil, nil, fmt.Errorf("redisGoodsCache.Init: %w", err)
	}

	eventLogger, err := natsGoodsEvent.InitLogger(config.Nats)
	if err != nil {
		return nil, nil, fmt.Errorf("natsGoodsEvent.InitLogger: %w", err)
	}

	closer := func() {
		var eg errgroup.Group
		eg.Go(func() error { return cache.Close() })
		eg.Go(func() error { return eventLogger.Close() })
		if err = eg.Wait(); err != nil {
			slog.Error("initAdapters: close: " + err.Error())
		}
	}

	return &adapters{
		goodsCache:       cache,
		goodsEventLogger: eventLogger,
	}, closer, nil
}
