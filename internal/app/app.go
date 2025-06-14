package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

func Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Инициализация адаптеров
	adapterss := initAdapters()

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
