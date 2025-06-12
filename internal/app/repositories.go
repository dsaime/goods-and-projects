package app

import (
	"fmt"
	"log/slog"

	"github.com/dsaime/goods-and-projects/internal/repository/pgsql"
)

type repositories struct {
	chats    chatt.Repository
	users    userr.Repository
	sessions sessionn.Repository
}

func initPgsqlRepositories(config pgsql.Config) (*repositories, func(), error) {
	factory, err := pgsql.InitFactory(config)
	if err != nil {
		return nil, func() {}, fmt.Errorf("pgsql.InitFactory: %w", err)
	}

	rs := &repositories{
		chats:    factory.NewChattRepository(),
		users:    factory.NewUserrRepository(),
		sessions: factory.NewSessionnRepository(),
	}

	return rs, func() {
		if err := factory.Close(); err != nil {
			slog.Error("initPgsqlRepositories: factory.Close: " + err.Error())
		}
	}, nil
}
