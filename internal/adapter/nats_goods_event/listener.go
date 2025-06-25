package natsGoodsEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"

	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

// Listener реализует интерфейс обработчика поступающих событий
type Listener struct {
	conn *nats.Conn
}

// Listen вызывает handler для каждого нового события. Выход из функции по отмене контекста
func (l *Listener) Listen(ctx context.Context, handler goodsEvent.Handler) error {
	subscribe, err := l.conn.Subscribe(goodsEventSubject, func(msg *nats.Msg) {
		var event goodsEvent.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			slog.Error("natsGoodsEvent: Listener: Unmarshal: " + err.Error())
			return
		}
		handler(event)
	})
	if err != nil {
		return fmt.Errorf("nats.Subscribe: %w", err)
	}
	defer func() { _ = subscribe.Unsubscribe() }()

	<-ctx.Done()
	return ctx.Err()
}

// Close закрывает соединение с NATS
func (l *Listener) Close() error {
	l.conn.Close()
	return nil
}

// ListenerConfig представляет собой конфигурацию адаптера
type ListenerConfig struct {
	NatsURL string
}

// InitListener инициализирует обработчик событий
func InitListener(config ListenerConfig) (*Listener, error) {
	connect, err := nats.Connect(config.NatsURL)
	if err != nil {
		return nil, err
	}

	slog.Info("natsGoodsEvent: Listener: Successfully connected to NATS ")

	return &Listener{
		conn: connect,
	}, nil
}
