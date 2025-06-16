package goodsEvent

import "context"

type Listener interface {
	Listen(context context.Context, handler Handler) error
}

type Handler func(event Event)
