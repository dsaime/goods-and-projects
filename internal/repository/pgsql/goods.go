package pgsql

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/nullism/bqb"

	"github.com/dsaime/goods-and-projects/internal/domain"
)

type GoodsRepository struct {
	db interface {
		NamedExec(query string, arg interface{}) (sql.Result, error)
		Exec(query string, args ...interface{}) (sql.Result, error)
		Select(dest interface{}, query string, args ...interface{}) error
	}
	txBeginner interface {
		BeginTx() (*sqlx.Tx, error)
	}
	isTx bool
}

func (r *GoodsRepository) Update(goodForUpdate domain.GoodForUpdate) (domain.Good, error) {
	if goodForUpdate.ID == 0 {
		return domain.Good{}, errors.New("ID не указан")
	}

	var good domain.Good
	if err := r.db.Select(&good, `
		UPDATE goods
		    SET name = $2, 
		        description = $3,
		        priority = $4,
		        removed = $5
		WHERE id = $1
		RETURNING goods.*
	`, goodForUpdate.ID,
		goodForUpdate.Name,
		goodForUpdate.Description,
		goodForUpdate.Priority,
		goodForUpdate.Removed,
	); err != nil {
		return domain.Good{}, err
	}

	return good, nil
}

func (r *GoodsRepository) Create(goodForSave domain.GoodForSave) (domain.Good, error) {
	if goodForSave.ID == 0 || goodForSave.ProjectID == 0 {
		return domain.Good{}, errors.New("ID или ProjectID не указан")
	}

	var good domain.Good
	if err := r.db.Select(&good, `
		INSERT INTO goods (id, project_id, name)
		VALUES ($1, $2, $3)
		RETURNING goods.*
	`, goodForSave.ID, goodForSave.ProjectID, goodForSave.Name); err != nil {
		return domain.Good{}, err
	}

	return good, nil
}

func (r *GoodsRepository) List(filter domain.GoodsFilter) ([]domain.Good, error) {
	query, args, err := buildListQuery(filter)
	if err != nil {
		return nil, err
	}

	var goods []domain.Good
	if err = r.db.Select(&goods, query, args...); err != nil {
		return nil, err
	}

	return goods, nil
}

func buildListQuery(filter domain.GoodsFilter) (string, []any, error) {
	selFrom := bqb.New("SELECT * FROM goods")

	where := bqb.Optional("WHERE")
	if filter.PriorityGreaterThan > 0 {
		where = where.And("priority > ?", filter.PriorityGreaterThan)
	}

	q := bqb.New("? ?", selFrom, where)
	if filter.Offset > 0 {
		q = q.Space("OFFSET ?", filter.Offset)
	}
	if filter.Limit > 0 {
		q = q.Space("LIMIT ?", filter.Limit)
	}

	return q.ToPgsql()
}

func (r *GoodsRepository) InTransaction(fn func(txRepo domain.GoodsRepository) error) error {
	// Начинаем транзакцию
	tx, err := r.txBeginner.BeginTx()
	if err != nil {
		return err
	}

	// Создаем транзакционный репозиторий
	txRepo := &GoodsRepository{
		db:         tx,
		txBeginner: r.txBeginner,
		isTx:       true,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // пробрасываем panic дальше
		}
	}()

	// Выполняем callback
	err = fn(txRepo)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Коммитим, если не было ошибок
	return tx.Commit()
}
