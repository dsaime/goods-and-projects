package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	redisGoodsCache "github.com/dsaime/goods-and-projects/internal/adapter/redis_goods_cache"
	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

func Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Инициализация адаптеров
	adapterss, closeAdapters, err := initAdapters(redisGoodsCache.Config{})
	if err != nil {
		return err
	}
	defer closeAdapters()

	// Инициализация репозиториев
	repos, closeRepos, err := initPgsqlRepositories(pgsql.Config{}, adapterss)
	if err != nil {
		return err
	}
	defer closeRepos()

	// Инициализация сервисов
	ss := initServices(repos, adapterss)

	// Инициализация и Запуск http контроллера
	server := initHttpServer(ss)
	g.Go(func() error {
		return runHttpServer(ctx, server)
	})

	return g.Wait()
}
