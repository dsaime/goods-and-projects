package redisGoodsCache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dsaime/goods-and-projects/internal/domain"
	goodsCache "github.com/dsaime/goods-and-projects/internal/port/goods_cache"
)

type GoodsCache struct {
	expiration time.Duration
	client     *redis.Client
	// keymutex.KeyMutex
}

func (g *GoodsCache) Close() error {
	return g.client.Close()
}

type Config struct {
	// RedisURL это строка в формате redis[s]://[[username][:password]@][host][:port][/db-number],
	// для подключения к redis.
	//
	// Пример: redis://alice:foobared@awesome.redis.server:6380
	RedisURL   string
	Expiration time.Duration
}

func Init(cfg Config) (*GoodsCache, error) {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err = client.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	slog.Info("Successfully connected to Redis")

	return &GoodsCache{
		expiration: cfg.Expiration,
		client:     client,
	}, nil
}

func (g *GoodsCache) Get(key goodsCache.Key) (domain.Good, bool) {
	slog.Debug("redisGoodsCache.Get",
		slog.Int("id", key.GetID()),
		slog.Int("projectID", key.GetProjectID()))
	b, err := g.client.Get(context.TODO(), strKey(key)).Bytes()
	if err != nil {
		slog.Error("redisGoodsCache.Get: "+err.Error(),
			slog.Int("id", key.GetID()),
			slog.Int("projectID", key.GetProjectID()))
		return domain.Good{}, false
	}

	var good domain.Good
	if json.Unmarshal(b, &good) != nil {
		return domain.Good{}, false
	}

	slog.Error("redisGoodsCache.Get: got "+fmt.Sprintf("%+v", good),
		slog.Int("id", key.GetID()),
		slog.Int("projectID", key.GetProjectID()),
	)
	return good, true
}

func (g *GoodsCache) Save(goods ...domain.Good) {
	var wg sync.WaitGroup
	wg.Add(len(goods))
	for _, good := range goods {
		go func() {
			defer wg.Done()
			if b, err := json.Marshal(good); err == nil {
				slog.Error("redisGoodsCache.Save: "+fmt.Sprintf("%+v", good),
					slog.Int("id", good.GetID()),
					slog.Int("projectID", good.GetProjectID()),
				)
				g.client.Set(context.TODO(), strKey(good), b, g.expiration)
			}
		}()
	}
	wg.Wait()
}

func (g *GoodsCache) Delete(keys ...goodsCache.Key) {
	kk := make([]string, len(keys))
	for i := range keys {
		kk[i] = strKey(keys[i])
		slog.Error("redisGoodsCache.Delete:",
			slog.Int("id", keys[i].GetID()),
			slog.Int("projectID", keys[i].GetProjectID()),
		)
	}

	g.client.Del(context.TODO(), kk...)
}

func strKey(key goodsCache.Key) string {
	return fmt.Sprintf("goods:%d_%d", key.GetID(), key.GetProjectID())
}
