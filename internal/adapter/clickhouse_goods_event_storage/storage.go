package clickhouseGoodsEventStorage

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"

	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

// Config представляет собой конфигурацию адаптера
type Config struct {
	DSN string
}

// Storage это реализация хранилища событий на основе clickhouse
type Storage struct {
	db *sqlx.DB
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() error {
	return s.db.Close()
}

// Init инициализирует хранилище
func Init(cfg Config) (*Storage, error) {
	options, err := clickhouse.ParseDSN(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("clickhouse.ParseDSN: %w", err)
	}

	conn := sqlx.NewDb(clickhouse.OpenDB(options), "clickhouse")

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("sqlx.Ping: %w", err)
	}

	slog.Info("Successfully connected to ClickHouse")

	return &Storage{
		db: conn,
	}, nil
}

func (s *Storage) Save(events ...goodsEvent.Event) error {
	if _, err := s.db.NamedExec(`
		INSERT INTO goods (Id, ProjectId, Name, Description, Priority, Removed, EventTime)
		VALUES (:Id, :ProjectId, :Name, :Description, :Priority, :Removed, :EventTime)
	`, toDBs(events)); err != nil {
		return err
	}

	return nil
}

type dbEvent struct {
	ID          int       `db:"Id"`
	ProjectID   int       `db:"ProjectId"`
	Name        string    `db:"Name"`
	Description string    `db:"Description"`
	Priority    int       `db:"Priority"`
	Removed     bool      `db:"Removed"`
	EventTime   time.Time `db:"EventTime"`
}

func toDB(event goodsEvent.Event) dbEvent {
	return dbEvent{
		ID:          event.ID,
		ProjectID:   event.ProjectID,
		Name:        event.Name,
		Description: event.Description,
		Priority:    event.Priority,
		Removed:     event.Removed,
		EventTime:   event.EventTime,
	}
}

func toDBs(event []goodsEvent.Event) []dbEvent {
	dbEvents := make([]dbEvent, len(event))
	for i, e := range event {
		dbEvents[i] = toDB(e)
	}
	return dbEvents
}
