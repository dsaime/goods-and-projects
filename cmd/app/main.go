package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dsaime/goods-and-projects/internal/app"
)

func main() {
	slog.Info("main: Starting")
	ctx, cancel := context.WithCancel(context.Background())
	go appRun(ctx)
	waitInterrupt(cancel)
}

func appRun(ctx context.Context) {
	err := app.Run(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("main: appRun:" + err.Error())
		os.Exit(1)
	}
	slog.Info("main: appRun: ")
	os.Exit(0)
}

func waitInterrupt(cancel context.CancelFunc) {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("main: waitInterrupt: Received signal " + (<-interrupt).String())
	cancel()
	slog.Info("main: waitInterrupt: Context canceled")
	time.Sleep(3 * time.Second)
}
