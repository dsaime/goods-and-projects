package goodsEventStorage

import goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"

type GoodsEventStorage interface {
	Save(event goodsEvent.Event) error
}
