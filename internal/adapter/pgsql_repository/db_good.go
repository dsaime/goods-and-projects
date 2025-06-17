package pgsqlRepository

import (
	"time"

	"github.com/dsaime/goods-and-projects/internal/domain"
)

type dbGood struct {
	ID          int       `db:"id"`
	ProjectID   int       `db:"project_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	CreatedAt   time.Time `db:"created_at"`
}

func toDomain(good dbGood) domain.Good {
	return domain.Good{
		ID:          good.ID,
		ProjectID:   good.ProjectID,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		Removed:     good.Removed,
		CreatedAt:   good.CreatedAt,
	}
}

func toDomains(goods []dbGood) []domain.Good {
	goodsDomain := make([]domain.Good, len(goods))
	for i, good := range goods {
		goodsDomain[i] = toDomain(good)
	}
	return goodsDomain
}
