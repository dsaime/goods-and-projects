package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/dsaime/goods-and-projects/internal/domain"
	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

// Goods сервис, объединяющий случаи использования(юзкейсы) в контексте товаров
type Goods struct {
	Repo             domain.GoodsRepository
	GoodsEventLogger goodsEvent.Logger
}

// GoodsIn входящие параметры
type GoodsIn struct {
	Limit  int
	Offset int
}

// Validate валидирует значение отдельно каждого параметры
func (in GoodsIn) Validate() error {
	if in.Limit <= 0 {
		return ErrLimitMustBeGTZero
	}
	if in.Offset < 0 {
		return ErrOffsetMustBePositive
	}

	return nil
}

// GoodsOut результат запроса чатов
type GoodsOut struct {
	Total   int
	Removed int
	Limit   int
	Offset  int
	Goods   []domain.Good
}

// Goods возвращает список товаров, с учетом сдвига и лимита
func (g *Goods) Goods(in GoodsIn) (GoodsOut, error) {
	if err := in.Validate(); err != nil {
		return GoodsOut{}, err
	}

	goodsList, err := g.Repo.List(domain.GoodsFilter{
		Limit:  in.Limit,
		Offset: in.Offset,
	})
	if err != nil {
		return GoodsOut{}, err
	}

	return GoodsOut{
		Total:   goodsList.Total,
		Removed: goodsList.Removed,
		Limit:   in.Limit,
		Offset:  in.Offset,
		Goods:   goodsList.Goods,
	}, err
}

// CreateGoodIn описывает входящие параметры
type CreateGoodIn struct {
	ProjectID int
	Name      string
}

// Validate валидирует значение отдельно каждого параметры
func (in CreateGoodIn) Validate() error {
	if in.ProjectID <= 0 {
		return ErrPriorityMustBeGTZero
	}
	if in.Name == "" {
		return ErrNameMustNotBeEmpty
	}

	return nil
}

// CreateGoodOut результат запроса чатов
type CreateGoodOut struct {
	CreatedGood domain.Good
}

// CreateGood создает товар
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

	event := newEventFromGood(createdGood)
	g.GoodsEventLogger.Log(event)

	return CreateGoodOut{
		CreatedGood: createdGood,
	}, nil
}

// getNewGoodID отдает ID, с которым можно сохранять новый товар.
// Может вернуть 0, если не найдет свободный ID.
// Свой выбор основывает на рандоме.
func getNewGoodID(repo domain.GoodsRepository, projectID int) int {
	for range 10 {
		randomID := int(rand.Int31())
		_, err := repo.Find(domain.GoodFilter{
			ID:          randomID,
			ProjectID:   projectID,
			ShowRemoved: true,
		})
		if errors.Is(err, domain.ErrGoodNotFound) {
			return randomID
		}
	}

	return 0
}

// UpdateGoodIn описывает входящие параметры
type UpdateGoodIn struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
}

// Validate валидирует значение отдельно каждого параметры
func (in UpdateGoodIn) Validate() error {
	if in.ID <= 0 {
		return ErrIDMustBeGTZero
	}
	if in.ProjectID <= 0 {
		return ErrProjectIDMustBeGTZero
	}
	if in.Name == "" {
		return ErrNameMustNotBeEmpty
	}

	return nil
}

// UpdateGoodOut результат запроса чатов
type UpdateGoodOut struct {
	UpdatedGood domain.Good
}

// UpdateGood обновляет товар
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

	event := newEventFromGood(updatedGood)
	g.GoodsEventLogger.Log(event)

	return UpdateGoodOut{
		UpdatedGood: updatedGood,
	}, err
}

// DeleteGoodIn описывает входящие параметры
type DeleteGoodIn struct {
	ID        int
	ProjectID int
}

// Validate валидирует значение отдельно каждого параметры
func (in DeleteGoodIn) Validate() error {
	if in.ID <= 0 {
		return ErrIDMustBeGTZero
	}
	if in.ProjectID <= 0 {
		return ErrProjectIDMustBeGTZero
	}

	return nil
}

// DeleteGoodOut результат запроса чатов
type DeleteGoodOut struct {
	DeletedGood DeleteGoodOutDeletedGood
}

type DeleteGoodOutDeletedGood struct {
	ID        int
	ProjectID int
	Removed   bool
}

// DeleteGood удаляет товар
func (g *Goods) DeleteGood(in DeleteGoodIn) (DeleteGoodOut, error) {
	if err := in.Validate(); err != nil {
		return DeleteGoodOut{}, err
	}
	var deletedGood domain.Good
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
		deletedGood, err = txRepo.Update(forUpdate)
		return err
	})
	if err != nil {
		return DeleteGoodOut{}, err
	}

	event := newEventFromGood(deletedGood)
	g.GoodsEventLogger.Log(event)

	return DeleteGoodOut{
		DeletedGood: DeleteGoodOutDeletedGood{
			ID:        deletedGood.ID,
			ProjectID: deletedGood.ProjectID,
			Removed:   true,
		},
	}, err
}

// ReprioritiizeGoodIn описывает входящие параметры
type ReprioritiizeGoodIn struct {
	ID          int
	ProjectID   int
	NewPriority int
}

// Validate валидирует значение отдельно каждого параметры
func (in ReprioritiizeGoodIn) Validate() error {
	if in.ID <= 0 {
		return ErrIDMustBeGTZero
	}
	if in.ProjectID <= 0 {
		return ErrProjectIDMustBeGTZero
	}
	if in.NewPriority <= 0 {
		return ErrPriorityMustBeGTZero
	}

	return nil
}

// ReprioritiizeGoodOut результат запроса чатов
type ReprioritiizeGoodOut struct {
	Priorities []ReprioritiizeGoodOutPriority
}

type ReprioritiizeGoodOutPriority struct {
	ID       int
	Priority int
}

// ReprioritiizeGood изменяет позицию у товара и двигает приоритет товаров, приоритет которых равен или больше переданному
func (g *Goods) ReprioritiizeGood(in ReprioritiizeGoodIn) (ReprioritiizeGoodOut, error) {
	if err := in.Validate(); err != nil {
		return ReprioritiizeGoodOut{}, err
	}

	var updatedGoods []domain.Good
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

		goodsList, err := txRepo.List(domain.GoodsFilter{
			PriorityGreaterThan: in.NewPriority - 1,
		})
		if err != nil {
			return err
		}
		for _, good := range goodsList.Goods {
			if good.ID == in.ID && good.ProjectID == in.ProjectID {
				continue
			}
			updatedGood, err := updatePriority(txRepo, good, good.Priority+1)
			if err != nil {
				return err
			}
			updatedGoods = append(updatedGoods, updatedGood)
		}
		updatedGood, err := updatePriority(txRepo, goodForReprioritiize, in.NewPriority)
		if err != nil {
			return err
		}
		updatedGoods = append(updatedGoods, updatedGood)

		return err
	})
	if err != nil {
		return ReprioritiizeGoodOut{}, err
	}

	for _, updatedGood := range updatedGoods {
		event := newEventFromGood(updatedGood)
		g.GoodsEventLogger.Log(event)
	}

	return ReprioritiizeGoodOut{
		Priorities: prioritiesFrom(updatedGoods),
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

func prioritiesFrom(goods []domain.Good) []ReprioritiizeGoodOutPriority {
	var outPriorities []ReprioritiizeGoodOutPriority
	for _, good := range goods {
		outPriorities = append(outPriorities, ReprioritiizeGoodOutPriority{
			ID:       good.ID,
			Priority: good.Priority,
		})
	}
	return outPriorities
}

func newEventFromGood(good domain.Good) goodsEvent.Event {
	return goodsEvent.Event{
		ID:          good.ID,
		ProjectID:   good.ProjectID,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		Removed:     good.Removed,
		EventTime:   time.Now(),
	}
}
