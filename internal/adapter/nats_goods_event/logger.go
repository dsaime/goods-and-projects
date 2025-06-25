package natsGoodsEvent

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

// Logger реализует интерфейс логгера событий
type Logger struct {
	conn *nats.Conn
}

// Close закрывает соединение с NATS
func (p *Logger) Close() error {
	p.conn.Close()
	return nil
}

// Log отправляет событие в NATS
func (p *Logger) Log(event goodsEvent.Event) {
	slog.Info("natsGoodsEvent: Logger.Log: new event for publishing",
		slog.String("event", fmt.Sprintf("%+v", event)))
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("natsGoodsEvent: Logger.Log: json.Marshal: "+err.Error(),
			slog.String("event", fmt.Sprintf("%+v", event)))
		return
	}

	if err = p.conn.Publish(goodsEventSubject, b); err != nil {
		slog.Error("natsGoodsEvent: Logger.Publish: nats.Publish: "+err.Error(),
			slog.String("event", fmt.Sprintf("%+v", event)))
		return
	}
	slog.Info("natsGoodsEvent: Logger.Publish: nats.Publish: published",
		slog.String("event", fmt.Sprintf("%+v", event)))
}

// LoggerConfig представляет собой конфигурацию адаптера
type LoggerConfig struct {
	NatsURL string
}

// InitLogger инициализирует логгер событий
func InitLogger(config LoggerConfig) (*Logger, error) {
	connect, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, err
	}

	slog.Info("natsGoodsEvent: Logger: Successfully connected to NATS ")

	return &Logger{
		conn: connect,
	}, nil
}
