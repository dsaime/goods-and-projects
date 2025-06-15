package redisGoodsCache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dsaime/goods-and-projects/internal/domain"
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

type cacheKey = interface {
	GetID() int
	GetProjectID() int
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

	return &GoodsCache{
		expiration: cfg.Expiration,
		client:     client,
	}, nil
}

func (g *GoodsCache) Get(key cacheKey) (domain.Good, bool) {
	b, err := g.client.Get(context.TODO(), strKey(key)).Bytes()
	if err != nil {
		return domain.Good{}, false
	}

	var good domain.Good
	if json.Unmarshal(b, &good) != nil {
		return domain.Good{}, false
	}

	return good, true
}

func (g *GoodsCache) Save(goods ...domain.Good) {
	var wg sync.WaitGroup
	wg.Add(len(goods))
	for _, good := range goods {
		go func() {
			defer wg.Done()
			if b, err := json.Marshal(good); err == nil {
				g.client.Set(context.TODO(), strKey(good), b, g.expiration)
			}
		}()
	}
	wg.Wait()
}

func (g *GoodsCache) Delete(keys ...cacheKey) {
	kk := make([]string, len(keys))
	for i := range keys {
		kk[i] = strKey(keys[i])
	}

	g.client.Del(context.TODO(), kk...)
}

func strKey(key cacheKey) string {
	return fmt.Sprintf("goods:%d_%d", key.GetID(), key.GetProjectID())
}
