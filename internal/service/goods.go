package service

import (
	"errors"

	"github.com/dsaime/goods-and-projects/internal/domain"
)

type Goods struct {
	Repo domain.GoodsRepository
}

type GoodsIn struct {
	Limit  int
	Offset int
}

func (in GoodsIn) Validate() error {
	if in.Limit <= 0 {
		return errors.New("limit must be positive and greater than 0")
	}
	if in.Offset < 0 {
		return errors.New("offset must be positive")
	}

	return nil
}

type GoodsOut struct {
	Meta  GoodsOutMeta
	Goods []domain.Good
}

type GoodsOutMeta struct {
	Total   int
	Removed int
	Limit   int
	Offset  int
}

func (g *Goods) Goods(in GoodsIn) (GoodsOut, error) {
	if err := in.Validate(); err != nil {
		return GoodsOut{}, err
	}

	goods, err := g.Repo.List(domain.GoodsFilter{
		Limit:  in.Limit,
		Offset: in.Offset,
	})
	if err != nil {
		return GoodsOut{}, err
	}

	return GoodsOut{
		Meta: GoodsOutMeta{
			Total:   len(goods),
			Removed: countByRemoved(goods),
			Limit:   in.Limit,
			Offset:  in.Offset,
		},
		Goods: goods,
	}, err
}

func countByRemoved(goods []domain.Good) int {
	var count int
	for _, good := range goods {
		if good.Removed {
			count++
		}
	}

	return count
}

type CreateGoodIn struct {
	ProjectID int
	Name      string
}

func (in CreateGoodIn) Validate() error {
	if in.ProjectID <= 0 {
		return errors.New("projectID is required")
	}
	if in.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

type CreateGoodOut struct {
	CreatedGood domain.Good
}

func (g *Goods) CreateGood(in CreateGoodIn) (CreateGoodOut, error) {
	if err := in.Validate(); err != nil {
		return CreateGoodOut{}, err
	}

	var createdGood domain.Good
	err := g.Repo.InTransaction(func(txRepo domain.GoodsRepository) error {
		forCreate := domain.GoodForCreate{
			ID:        getNewGoodID(txRepo, in.ProjectID),
			ProjectID: in.ProjectID,
			Name:      in.Name,
		}
		var err error
		createdGood, err = g.Repo.Create(forCreate)
		return err
	})
	if err != nil {
		return CreateGoodOut{}, err
	}

	return CreateGoodOut{
		CreatedGood: createdGood,
	}, nil
}

func getNewGoodID(repo domain.GoodsRepository, projectID int) int {
	for randomID := range 10 {
		_, err := repo.Find(domain.GoodFilter{
			ID:        randomID,
			ProjectID: projectID,
		})
		if errors.Is(err, domain.ErrGoodNotFound) {
			return randomID
		}
	}

	return 0
}

type UpdateGoodIn struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
}

func (in UpdateGoodIn) Validate() error {
	if in.ID <= 0 {
		return errors.New("GetID is required")
	}
	if in.ProjectID <= 0 {
		return errors.New("projectID is required")
	}
	if in.Name == "" {
		return errors.New("name must not be empty")
	}

	return nil
}

type UpdateGoodOut struct {
	UpdatedGood domain.Good
}

func (g *Goods) UpdateGood(in UpdateGoodIn) (UpdateGoodOut, error) {
	if err := in.Validate(); err != nil {
		return UpdateGoodOut{}, err
	}

	var updatedGood domain.Good
	err := g.Repo.InTransaction(func(txRepo domain.GoodsRepository) error {
		oldGood, err := txRepo.Find(domain.GoodFilter{
			ID:        in.ID,
			ProjectID: in.ProjectID,
		})
		if err != nil {
			return err
		}
		forUpdate := domain.GoodForUpdate{
			ID:          oldGood.ID,
			ProjectID:   oldGood.ProjectID,
			Name:        in.Name,
			Description: in.Description,
			Priority:    oldGood.Priority,
			Removed:     oldGood.Removed,
		}
		updatedGood, err = txRepo.Update(forUpdate)
		return err
	})
	if err != nil {
		return UpdateGoodOut{}, err
	}

	return UpdateGoodOut{
		UpdatedGood: updatedGood,
	}, err
}

type DeleteGoodIn struct {
	ID        int
	ProjectID int
}

func (in DeleteGoodIn) Validate() error {
	if in.ID <= 0 {
		return errors.New("GetID is required")
	}
	if in.ProjectID <= 0 {
		return errors.New("projectID is required")
	}

	return nil
}

type DeleteGoodOut struct {
	DeletedGood DeleteGoodOutDeletedGood
}

type DeleteGoodOutDeletedGood struct {
	ID        int
	ProjectID int
	Removed   bool
}

func (g *Goods) DeleteGood(in DeleteGoodIn) (DeleteGoodOut, error) {
	if err := in.Validate(); err != nil {
		return DeleteGoodOut{}, err
	}
	var deletedGood DeleteGoodOutDeletedGood
	err := g.Repo.InTransaction(func(txRepo domain.GoodsRepository) error {
		oldGood, err := txRepo.Find(domain.GoodFilter{
			ID:        in.ID,
			ProjectID: in.ProjectID,
		})
		if err != nil {
			return err
		}

		if oldGood.Removed {
			return errors.New("good is already removed")
		}

		forUpdate := domain.GoodForUpdate{
			ID:          oldGood.ID,
			ProjectID:   oldGood.ProjectID,
			Name:        oldGood.Name,
			Description: oldGood.Description,
			Priority:    oldGood.Priority,
			Removed:     true,
		}
		updatedGood, err := txRepo.Update(forUpdate)
		deletedGood = DeleteGoodOutDeletedGood{
			ID:        updatedGood.ID,
			ProjectID: updatedGood.ProjectID,
			Removed:   true,
		}
		return err
	})
	if err != nil {
		return DeleteGoodOut{}, err
	}

	return DeleteGoodOut{
		DeletedGood: deletedGood,
	}, err
}

type ReprioritiizeGoodIn struct {
	ID          int
	ProjectID   int
	NewPriority int
}

func (in ReprioritiizeGoodIn) Validate() error {
	if in.ID <= 0 {
		return errors.New("GetID is required")
	}
	if in.ProjectID <= 0 {
		return errors.New("projectID is required")
	}
	if in.NewPriority <= 0 {
		return errors.New("newPriority is required")
	}

	return nil
}

type ReprioritiizeGoodOut struct {
	Priorities []ReprioritiizeGoodOutPriority
}

type ReprioritiizeGoodOutPriority struct {
	ID       int
	Priority int
}

func (g *Goods) ReprioritiizeGood(in ReprioritiizeGoodIn) (ReprioritiizeGoodOut, error) {
	if err := in.Validate(); err != nil {
		return ReprioritiizeGoodOut{}, err
	}

	var priorities []ReprioritiizeGoodOutPriority
	err := g.Repo.InTransaction(func(txRepo domain.GoodsRepository) error {
		goodForReprioritiize, err := txRepo.Find(domain.GoodFilter{
			ID:        in.ID,
			ProjectID: in.ProjectID,
		})
		if err != nil {
			return err
		}
		if goodForReprioritiize.Priority == in.NewPriority {
			return nil
		}

		goods, err := txRepo.List(domain.GoodsFilter{
			PriorityGreaterThan: in.NewPriority - 1,
		})
		if err != nil {
			return err
		}
		for _, good := range goods {
			if good.ID == in.ID && good.ProjectID == in.ProjectID {
				continue
			}
			updatedGood, err := updatePriority(txRepo, good, good.Priority+1)
			if err != nil {
				return err
			}
			addReprioritiizeElem(&priorities, updatedGood)
		}
		updatedGood, err := updatePriority(txRepo, goodForReprioritiize, in.NewPriority)
		if err != nil {
			return err
		}
		addReprioritiizeElem(&priorities, updatedGood)

		return err
	})
	if err != nil {
		return ReprioritiizeGoodOut{}, err
	}

	return ReprioritiizeGoodOut{
		Priorities: priorities,
	}, err
}

func updatePriority(txRepo domain.GoodsRepository, good domain.Good, newPriority int) (domain.Good, error) {
	forUpdate := domain.GoodForUpdate{
		ID:          good.ID,
		ProjectID:   good.ProjectID,
		Name:        good.Name,
		Description: good.Description,
		Priority:    newPriority,
		Removed:     good.Removed,
	}

	return txRepo.Update(forUpdate)
}

func addReprioritiizeElem(s *[]ReprioritiizeGoodOutPriority, good domain.Good) {
	*s = append(*s, ReprioritiizeGoodOutPriority{
		ID:       good.ID,
		Priority: good.Priority,
	})
}
