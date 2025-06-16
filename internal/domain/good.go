package domain

import "time"

type Good struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
	Priority    int
	Removed     bool
	CreatedAt   time.Time
}

func (g Good) GetID() int {
	return g.ID
}

func (g Good) GetProjectID() int {
	return g.ID
}

type GoodsRepository interface {
	List(filter GoodsFilter) (GoodsListResult, error)
	Find(filter GoodFilter) (Good, error)
	Update(GoodForUpdate) (Good, error)
	Create(GoodForCreate) (Good, error)
	InTransaction(func(txRepo GoodsRepository) error) error
}

type GoodsListResult struct {
	Total   int
	Removed int
	Goods   []Good
}

type GoodFilter struct {
	ID          int
	ProjectID   int
	ShowRemoved bool
}

type GoodForUpdate struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
	Priority    int
	Removed     bool
}

type GoodForCreate struct {
	ID        int
	ProjectID int
	Name      string
}

type GoodsFilter struct {
	PriorityGreaterThan int
	Limit               int
	Offset              int
}
