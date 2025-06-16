package goodsEvent

import "time"

// Event представляет собой событие изменения состояния товара
type Event struct {
	ID          int       // Идентификатор
	ProjectID   int       // Идентификатор
	Name        string    // название
	Description string    // Описание
	Priority    int       // Приоритет
	Removed     bool      // Признак удаления
	EventTime   time.Time // Дата и время произошедшего изменения
}
