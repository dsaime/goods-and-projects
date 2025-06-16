package app

import (
	redisGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/redis_goods_cache"
	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

type Config struct {
	Pgsql    pgsql.Config
	Redis    redisGoodsCache.Config
	HttpAddr string
}

//
//func Parse(args []string) (Config, error) {
//
//}
