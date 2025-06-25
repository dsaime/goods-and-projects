package goodsEvent

// Logger представляет собой интерфейс отправки логов об обновлении состояния товаров
type Logger interface {
	Log(event Event)
}
