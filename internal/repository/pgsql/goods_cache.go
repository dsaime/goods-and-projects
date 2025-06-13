package pgsql

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(CacheKey) (domain.Good, bool)
	Save(goods ...domain.Good)
	Delete(...CacheKey)
}

type CacheKey interface {
	GetID() int
	GetProjectID() int
}
