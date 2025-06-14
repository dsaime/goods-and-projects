package dummy_goods_cache

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

type GoodsCache struct{}

func (g GoodsCache) Get(key cacheKey) (domain.Good, bool) {
	return domain.Good{}, false
}

func (g GoodsCache) Save(goods ...domain.Good) {}

func (g GoodsCache) Delete(key ...cacheKey) {}

type cacheKey = interface {
	GetID() int
	GetProjectID() int
}
