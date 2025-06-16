package app

import (
	"context"
	"fmt"
	"log/slog"

	clickhouseGoodsEventStorage "github.com/dsaime/goods-and-projects/internal/adapter/clickhouse_goods_event_storage"
	natsGoodsEvent "github.com/dsaime/goods-and-projects/internal/adapter/nats_goods_event"
	goodsEvent "github.com/dsaime/goods-and-projects/internal/domain/goods_event"
)

func Run(ctx context.Context, config Config) error {
	// Инициализация хранилища событий
	storage, err := clickhouseGoodsEventStorage.Init(config.Clickhouse)
	if err != nil {
		return fmt.Errorf("clickhouseGoodsEventStorage.Init: %w", err)
	}

	// Инициализация подписчика событий
	listener, err := natsGoodsEvent.InitListener(config.Nats)
	if err != nil {
		return fmt.Errorf("natsGoodsEvent.InitListener: %w", err)
	}
	defer func() { _ = listener.Close() }()

	// Прослушивать и сохранять события
	return listener.Listen(ctx, func(event goodsEvent.Event) {
		slog.Info("новое событие", "event", fmt.Sprintf("%+v", event))
		if err := storage.Save(event); err != nil {
			slog.Error("listener handler: storage.Save: " + err.Error())
			return
		}
		slog.Info("событие сохранено", "event", fmt.Sprintf("%+v", event))
	})
}
