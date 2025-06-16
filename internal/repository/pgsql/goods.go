package pgsql

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/nullism/bqb"

	"github.com/dsaime/goods-and-projects/internal/domain"
	goodsCache "github.com/dsaime/goods-and-projects/internal/port/goods_cache"
)

type goodsRepository struct {
	db interface {
		NamedExec(query string, arg interface{}) (sql.Result, error)
		Exec(query string, args ...interface{}) (sql.Result, error)
		Select(dest interface{}, query string, args ...interface{}) error
		Get(dest interface{}, query string, args ...interface{}) error
	}
	txBeginner interface {
		Beginx() (*sqlx.Tx, error)
	}
	isTx  bool
	cache goodsCache.GoodsCache
}

func (r *goodsRepository) Find(filter domain.GoodFilter) (domain.Good, error) {
	if filter.ID == 0 || filter.ProjectID == 0 {
		return domain.Good{}, errors.New("ID или ProjectID не указан")
	}

	if good, ok := r.cache.Get(domain.Good{
		ID:        filter.ID,
		ProjectID: filter.ProjectID,
	}); ok {
		return good, nil
	}
	var good dbGood
	err := r.db.Get(&good, `
		SELECT * FROM goods
		WHERE id = $1 
		  AND project_id = $2 
		  AND ($3::BOOLEAN OR NOT removed)
	`, filter.ID, filter.ProjectID, filter.ShowRemoved)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Good{}, domain.ErrGoodNotFound
	} else if err != nil {
		return domain.Good{}, err
	}

	domainGood := toDomain(good)

	r.cache.Save(domainGood)

	return domainGood, nil
}

func (r *goodsRepository) Update(goodForUpdate domain.GoodForUpdate) (domain.Good, error) {
	if goodForUpdate.ID == 0 || goodForUpdate.ProjectID == 0 {
		return domain.Good{}, errors.New("ID или ProjectID не указан")
	}

	var good dbGood
	if err := r.db.Get(&good, `
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

	r.cache.Delete(toDomain(good))

	return toDomain(good), nil
}

func (r *goodsRepository) Create(goodForSave domain.GoodForCreate) (domain.Good, error) {
	if goodForSave.ID == 0 || goodForSave.ProjectID == 0 {
		return domain.Good{}, errors.New("ID или ProjectID не указан")
	}

	var good dbGood
	if err := r.db.Get(&good, `
		INSERT INTO goods (id, project_id, name)
		VALUES ($1, $2, $3)
		RETURNING *
	`, goodForSave.ID, goodForSave.ProjectID, goodForSave.Name); err != nil {
		return domain.Good{}, err
	}

	return toDomain(good), nil
}

func (r *goodsRepository) List(filter domain.GoodsFilter) (domain.GoodsListResult, error) {
	// Запросить количественные общие значения
	queryMeta, argsMeta, err := buildQueryListMeta(filter)
	if err != nil {
		return domain.GoodsListResult{}, err
	}
	var meta struct {
		Total   int `db:"total"`
		Removed int `db:"removed"`
	}
	if err = r.db.Get(&meta, queryMeta, argsMeta...); err != nil {
		return domain.GoodsListResult{}, err
	}

	// Запросить товары с учетом сдвига и лимита
	query, args, err := buildQueryList(filter, r.isTx)
	if err != nil {
		return domain.GoodsListResult{}, err
	}
	var goods []dbGood
	if err = r.db.Select(&goods, query, args...); err != nil {
		return domain.GoodsListResult{}, err
	}

	// Сохранить в кэш
	domainGoods := toDomains(goods)
	r.cache.Save(domainGoods...)

	return domain.GoodsListResult{
		Total:   meta.Total,
		Removed: meta.Removed,
		Goods:   domainGoods,
	}, nil
}

func buildQueryList(filter domain.GoodsFilter, forUpdate bool) (query string, args []any, err error) {
	selFrom := bqb.New("SELECT * FROM goods")

	where := bqb.New("WHERE NOT removed")
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
	if forUpdate {
		q = q.Space("FOR UPDATE")
	}

	return q.ToPgsql()
}

func buildQueryListMeta(filter domain.GoodsFilter) (query string, args []any, err error) {
	selFrom := bqb.New(`
		SELECT COUNT(1) AS total,
		       COUNT(1) FILTER ( WHERE removed ) AS removed
		FROM goods
	`)

	where := bqb.Optional("WHERE")
	if filter.PriorityGreaterThan > 0 {
		where = where.And("priority > ?", filter.PriorityGreaterThan)
	}

	q := bqb.New("? ?", selFrom, where)

	return q.ToPgsql()
}

func (r *goodsRepository) InTransaction(fn func(txRepo domain.GoodsRepository) error) error {
	// Начинаем транзакцию
	tx, err := r.txBeginner.Beginx()
	if err != nil {
		return err
	}

	// Создаем транзакционный репозиторий
	txRepo := &goodsRepository{
		db:         tx,
		txBeginner: r.txBeginner,
		isTx:       true,
		cache:      r.cache,
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
