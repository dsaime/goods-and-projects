package pgsql

import (
	"time"

	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache взаимодействует с кэшем, читать сохранить
type GoodsCache interface {
	Get(id int) (domain.Good, error)
	Set(domain.Good, time.Duration) error
	Delete(id int) error
}
