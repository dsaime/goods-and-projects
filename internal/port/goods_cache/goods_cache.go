package goodsCache

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(Key) (domain.Good, bool)
	Save(goods ...domain.Good)
	Delete(...Key)
}

type Key interface {
	GetID() int
	GetProjectID() int
}
