package app

import goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"

type goodsEventStorage interface {
	Save(...goodsEvent.Event) error
}
