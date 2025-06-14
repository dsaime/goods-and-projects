package domain

import (
	"time"
)

type Project struct {
	ID        int
	Name      string
	CreatedAt time.Time
}
