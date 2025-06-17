package goodsCache

import (
	"github.com/dsaime/goods-and-projects/internal/domain"
)

// GoodsCache представляет собой интерфейс кэша товаров
type GoodsCache interface {
	// Get Возвращает товар, сохраненный по переданному ключу.
	// Если такой записи нет, вторым значением вернется false
	Get(Key) (domain.Good, bool)

	// Save Сохраняет товары в кэш.
	// Каждому товару будет соответствовать свой ключ
	Save(goods ...domain.Good)

	// Delete Удаляет(инвалидирует) запись по ключу
	Delete(...Key)
}

// Key представляет собой составной ключ товара
type Key interface {
	GetID() int
	GetProjectID() int
}
