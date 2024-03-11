package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/cicd-lectures/vehicle-server/app"
	"go.uber.org/zap"
)

func main() { os.Exit(run()) }

func run() int {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var cfg app.Config

	flag.StringVar(&cfg.DatabaseURL, "database-url", "", "URL of the database")
	flag.StringVar(&cfg.ListenAddress, "listen-address", ":8080", "Address to listen to")

	flag.Parse()

	logger := zap.Must(zap.NewDevelopment())

	app, err := app.New(ctx, cfg, logger)
	if err != nil {
		return 1
	}

	defer func() {
		_ = app.Close()
	}()

	if err := app.Run(ctx); err != nil {
		return 1
	}

	return 0
}
