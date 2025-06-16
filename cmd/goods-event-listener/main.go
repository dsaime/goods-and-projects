package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"

	app "github.com/dsaime/goods-and-projects/internal/app/goods_event_listener"
)

func main() {
	slog.Info("main: starting")
	ctx, cancel := context.WithCancel(context.Background())
	go appRun(ctx)
	waitInterrupt(cancel)
}

// appRun запускает приложение и обрабатывает результат
func appRun(ctx context.Context) {
	err := initCliCommand().Run(ctx, os.Args)
	if errors.Is(err, context.Canceled) {
		slog.Info("main: appRun: exit by context canceled")
		os.Exit(0)
	} else if err != nil {
		slog.Error("main: appRun: " + err.Error())
		os.Exit(1)
	}
	slog.Info("main: appRun: finished")
	os.Exit(0)
}

// waitInterrupt отменяет контекст, когда в приложение поступает сигнал syscall.SIGINT или syscall.SIGTERM
func waitInterrupt(cancel context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("main: waitInterrupt: Received signal " + (<-interrupt).String())
	cancel()
	slog.Info("main: waitInterrupt: Context canceled")
	time.Sleep(3 * time.Second)
}

// initCliCommand создает команду, для разбора аргументов командной строки и запуска приложения
func initCliCommand() *cli.Command {
	var cfg app.Config
	return &cli.Command{
		Name: "goods-event-listener",
		Action: func(ctx context.Context, command *cli.Command) error {
			return app.Run(ctx, cfg)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "nats-url",
				Destination: &cfg.Nats.NatsURL,
				Usage:       "NATS connection string in format 'nats://user:password@host:port'",
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "clickhouse-dsn",
				Destination: &cfg.Clickhouse.DSN,
				Usage:       "ClickHouse connection string in format 'clickhouse://user:password@host:port/dbname?param1=value1&param2=value2'",
				Required:    true,
			},
		},
	}
}
