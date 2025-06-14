package pgsql

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(cacheKey) (domain.Good, bool)
	Save(goods ...domain.Good)
	Delete(...cacheKey)
}

type cacheKey = interface {
	GetID() int
	GetProjectID() int
}
