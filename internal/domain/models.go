package domain

import "time"

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
	Upsert(Good) error
}

type GoodsFilter struct {
	Limit  int
	Offset int
}
