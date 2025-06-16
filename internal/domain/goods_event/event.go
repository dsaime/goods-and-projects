package goodsEvent

import "time"

type Event struct {
	ID          int       // Идентификатор
	ProjectID   int       // Идентификатор
	Name        string    //  название
	Description string    // Описание
	Priority    int       // Приоритет
	Removed     bool      // Статус удаления
	EventTime   time.Time // Дата и время
}
