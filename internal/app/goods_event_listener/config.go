package app

import (
	"github.com/dsaime/goods-and-projects/internal/adapter/clickhouse_goods_event_storage"
	natsGoodsEvent "github.com/dsaime/goods-and-projects/internal/adapter/nats_goods_event"
)

type Config struct {
	Nats       natsGoodsEvent.ListenerConfig
	Clickhouse clickhouseGoodsEventStorage.Config
}
