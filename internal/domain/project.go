package domain

import (
	"time"
)

// Project представляет собой сущность компании.
type Project struct {
	ID        int       // Идентификатор компании
	Name      string    // Название компании
	CreatedAt time.Time // Дата и время создания компании
}
