package natsGoodsEvent

import (
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"

	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

type Logger struct {
	conn *nats.Conn
}

func (p *Logger) Close() error {
	p.conn.Close()
	return nil
}

func (p *Logger) Log(event goodsEvent.Event) {
	b, err := json.Marshal(event)
	if err != nil {
		slog.Error("natsGoodsEvent: Logger: json.Marshal: " + err.Error())
		return
	}

	if err = p.conn.Publish(goodsEventSubject, b); err != nil {
		slog.Error("natsGoodsEvent: Logger: nats.Publish: " + err.Error())
	}
}

type LoggerConfig struct {
	NatsURL string
}

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
