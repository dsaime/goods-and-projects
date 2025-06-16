package goodsEventStorage

import goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"

type GoodsEventStorage interface {
	Save(...goodsEvent.Event) error
}
