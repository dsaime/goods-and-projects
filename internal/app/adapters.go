package app

import (
	dummyGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/dummy_goods_cache"
	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

type adapters struct {
	//discovery      adapter.ServiceDiscovery
	goodsCache pgsql.GoodsCache
}

func (a *adapters) GoodsCache() pgsql.GoodsCache {
	return a.goodsCache
}

//func (a *adapters) Discovery() adapter.ServiceDiscovery {
//	return a.discovery
//}
//

func initAdapters() *adapters {

	return &adapters{
		goodsCache: dummyGoodsCache.GoodsCache{},
	}
}
