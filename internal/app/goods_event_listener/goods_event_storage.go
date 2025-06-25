package app

import goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"

// goodsEventStorage представляет собой интерфейс хранилища событий
type goodsEventStorage interface {
	// Save сохраняет пачку событий
	Save(...goodsEvent.Event) error
}
