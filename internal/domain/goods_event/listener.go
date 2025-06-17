package goodsEvent

import "context"

// Listener представляет собой интерфейс обработки поступающих событий об обновлении состояния товаров
type Listener interface {
	Listen(context context.Context, handler Handler) error
}

// Handler обработчик событий об обновлении состояния товаров
type Handler func(event Event)
