package app

import (
	natsGoodsEvent "github.com/dsaime/goods-and-projects/internal/adapter/nats_goods_event"
	pgsqlRepository "github.com/dsaime/goods-and-projects/internal/adapter/pgsql_repository"
	redisGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/redis_goods_cache"
)

type Config struct {
	Pgsql    pgsqlRepository.Config
	Nats     natsGoodsEvent.LoggerConfig
	Redis    redisGoodsCache.Config
	HttpAddr string
}
