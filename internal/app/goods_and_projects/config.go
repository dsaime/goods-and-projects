package app

import (
	natsGoodsEvent "github.com/dsaime/goods-and-projects/internal/adapter/nats_goods_event"
	redisGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/redis_goods_cache"
	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

type Config struct {
	Pgsql    pgsql.Config
	Nats     natsGoodsEvent.LoggerConfig
	Redis    redisGoodsCache.Config
	HttpAddr string
}
