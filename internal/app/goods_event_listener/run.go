package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	clickhouseGoodsEventStorage "github.com/dsaime/goods-and-projects/internal/adapter/clickhouse_goods_event_storage"
	natsGoodsEvent "github.com/dsaime/goods-and-projects/internal/adapter/nats_goods_event"
	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

// Run запускает приложение
func Run(ctx context.Context, config Config) error {
	// Инициализация хранилища событий
	storage, err := clickhouseGoodsEventStorage.Init(config.Clickhouse)
	if err != nil {
		return fmt.Errorf("clickhouseGoodsEventStorage.Init: %w", err)
	}
	defer func() { _ = storage.Close() }()

	// Инициализация подписчика событий
	listener, err := natsGoodsEvent.InitListener(config.Nats)
	if err != nil {
		return fmt.Errorf("natsGoodsEvent.InitListener: %w", err)
	}
	defer func() { _ = listener.Close() }()

	// Прослушивать и сохранять события
	return listener.Listen(ctx, eventHandler(storage))
}

// eventHandler возвращает обработчик событий.
// Приходящие события сохраняются в хранилище пачками либо по таймеру
func eventHandler(storage goodsEventStorage) func(event goodsEvent.Event) {
	var mu sync.Mutex
	const batchSize = 1000
	const flushInterval = 2 * time.Second // ms
	batch := make([]goodsEvent.Event, 0, batchSize)
	flushTimer := time.NewTimer(flushInterval)
	go func() {
		for {
			<-flushTimer.C
			mu.Lock()
			if len(batch) > 0 {
				saveBatch(storage, batch)
				batch = batch[:0]
			}
			flushTimer.Reset(flushInterval)
			mu.Unlock()
		}
	}()

	return func(event goodsEvent.Event) {
		mu.Lock()
		defer mu.Unlock()
		slog.Info("новое событие",
			slog.String("event", fmt.Sprintf("%+v", event)))

		batch = append(batch, event)
		if len(batch) >= batchSize {
			saveBatch(storage, batch)
			batch = batch[:0]
			flushTimer.Reset(flushInterval)
		} else {
			slog.Info("событие будет сохранено в пачке",
				slog.String("event", fmt.Sprintf("%+v", event)))
		}
	}
}

// saveBatch сохраняет пачку событий
func saveBatch(storage goodsEventStorage, batch []goodsEvent.Event) {
	if err := storage.Save(batch...); err != nil {
		slog.Error("listener handler: storage.Save: " + err.Error())
		return
	}
	slog.Info(fmt.Sprintf("пачка из %d элементов сохранена", len(batch)))
}
