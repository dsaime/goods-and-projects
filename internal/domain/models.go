package domain

import (
	"time"
)

type Project struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type Good struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}

type GoodsRepository interface {
	List(filter GoodsFilter) ([]Good, error)
	Update(GoodForUpdate) (Good, error)
	Create(GoodForSave) (Good, error)
	InTransaction(func(txRepo GoodsRepository) error) error
}

type GoodForUpdate struct {
	ID          int
	Name        string
	Description string
	Priority    int
	Removed     bool
}

type GoodForSave struct {
	ID        int
	ProjectID int
	Name      string
}

type GoodsFilter struct {
	//WithRemoved         bool
	PriorityGreaterThan int
	Limit               int
	Offset              int
}
