package pgsql

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(id int) (domain.Good, bool)
	Save(goods ...domain.Good)
	Delete(ids ...int)
}
