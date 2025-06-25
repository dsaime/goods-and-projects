package domain

import "time"

// Good представляет собой сущность товара
type Good struct {
	// ID это числовой идентификатор товара.
	// Первая часть составного ключа товара.
	// После создания нельзя изменить.
	ID int

	// ProjectID это числовой идентификатор компании, которой он принадлежит.
	// Вторая часть составного ключа товара.
	// После создания нельзя изменить.
	ProjectID int

	// Name это строковое название товара, которое ни на что не влияет.
	// Можно изменить после создания.
	Name string

	// Description это текстовое описание товара, которое ни на что не влияет.
	// Можно изменить после создания.
	Description string

	// Priority это число, которое ни на что не влияет.
	// Можно изменять после создания.
	// Обновление приоритета товара требует пересчета приоритета у других товаров.
	// Если товар удален (Removed=true), то Priority не пересчитывается.
	Priority int

	// Removed это признак удаления товара.
	// После создания можно установить в true, но обратно в false нельзя.
	// Товары с Removed=true, нельзя изменять.
	Removed bool

	// CreatedAt это дата и время создания товара
	CreatedAt time.Time
}

// GetID возвращает идентификатор товара.
func (g Good) GetID() int {
	return g.ID
}

// GetProjectID возвращает идентификатор компании товара.
func (g Good) GetProjectID() int {
	return g.ID
}

// GoodsRepository представляет собой интерфейс доступа к товарам.
type GoodsRepository interface {
	// List возвращает список товаров с учетом фильтрации и информацию по выборке
	List(filter GoodsFilter) (GoodsListResult, error)

	// Find возвращает товар по составному ключу.
	// Возвращает ErrGoodNotFound, если такого товара не существует
	Find(filter GoodFilter) (Good, error)

	// Update обновляет определенные поля товара и возвращает обновленный товар
	Update(GoodForUpdate) (Good, error)

	// Create создает и возвращает новый товар
	Create(GoodForCreate) (Good, error)

	// InTransaction оборачивает выполнение переданной функции в транзакцию.
	// В случае ошибки — откатывает транзакцию, в случае успеха — фиксирует.
	InTransaction(func(txRepo GoodsRepository) error) error
}

// GoodsListResult представляет собой результат получения списка товаров с учетом фильтрации, сдвига и лимита.
type GoodsListResult struct {
	Total   int    // Количество товаров, без ограничений по сдвигу и лимиту.
	Removed int    // Количество удаленных товаров, без ограничений по сдвигу и лимиту.
	Goods   []Good // Товары, полученные, применив сдвиг, лимит и фильтрацию удаленных.
}

// GoodFilter представляет собой фильтр для поиска товара.
type GoodFilter struct {
	ID           int
	ProjectID    int
	AllowRemoved bool // Разрешать поиск по удаленным товарам
}

// GoodForUpdate представляет собой пару идентификаторов товара и новые значения его полей.
type GoodForUpdate struct {
	ID          int
	ProjectID   int
	Name        string
	Description string
	Priority    int
	Removed     bool
}

// GoodForCreate представляет собой идентификатор товара и значения полей,
// с которыми он будет создан (кроме тех, что по умолчанию устанавливаются СУБД(осуждаю такое поведение))
type GoodForCreate struct {
	ID        int
	ProjectID int
	Name      string
}

// GoodsFilter представляет собой фильтр списка товаров.
type GoodsFilter struct {
	PriorityGreaterThan int // Выбирать товары, у которых приоритет выше переданного
	Limit               int
	Offset              int
}
