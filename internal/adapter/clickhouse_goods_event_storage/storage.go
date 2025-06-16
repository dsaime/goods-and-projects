package clickhouseGoodsEventStorage

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"

	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

type Config struct {
	DSN string
}

type Storage struct {
	db *sqlx.DB
}

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

func (s *Storage) Save(event goodsEvent.Event) error {
	if _, err := s.db.NamedExec(`
		INSERT INTO goods (Id, ProjectId, Name, Description, Priority, Removed, EventTime)
		VALUES (:Id, :ProjectId, :Name, :Description, :Priority, :Removed, :EventTime)
	`, toDB(event)); err != nil {
		return err
	}

	return nil
}

type dbEvent struct {
	ID          int       `json:"Id"`
	ProjectID   int       `json:"ProjectId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Priority    int       `json:"Priority"`
	Removed     bool      `json:"Removed"`
	EventTime   time.Time `json:"EventTime"`
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
