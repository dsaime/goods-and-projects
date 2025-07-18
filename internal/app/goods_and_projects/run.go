package app

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Run запускает приложение
func Run(ctx context.Context, config Config) error {
	g, ctx := errgroup.WithContext(ctx)

	// Инициализация адаптеров
	adapterss, closeAdapters, err := initAdapters(config)
	if err != nil {
		return err
	}
	defer closeAdapters()

	// Инициализация репозиториев
	repos, closeRepos, err := initPgsqlRepositories(config, adapterss)
	if err != nil {
		return err
	}
	defer closeRepos()

	// Инициализация сервисов
	ss := initServices(repos, adapterss)

	// Инициализация и Запуск http контроллера
	server := initHttpServer(config, ss)
	g.Go(func() error {
		return runHttpServer(ctx, server)
	})

	return g.Wait()
}
