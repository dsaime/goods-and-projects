package pgsql

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(id int) (domain.Good, error)
	Set(domain.Good) error
	Delete(id int) error
}
