package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/queue"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := loadConfig(configFile)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		return
	}

	storage, err := sqlstorage.New(
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
		config.Postgres.User,
		config.Postgres.Password,
	)
	if err != nil {
		slog.Error("failed to connect to postgres", "err", err)
		return
	}
	if err := storage.Connect(context.Background()); err != nil {
		slog.Error("failed to connect to storage", "err", err)
		return
	}
	defer storage.Close(context.Background())

	queue, err := queue.New(
		config.Rabbit.Host,
		config.Rabbit.Port,
		config.Rabbit.User,
		config.Rabbit.Password,
	)
	defer queue.Close()

	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		slog.Error("failed to parse interval", "err", err)
		return
	}

	if err := storage.SetInterval(interval); err != nil {
		slog.Error("failed to set interval", "err", err)
		return
	}

	sched := scheduler.New(storage, queue)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := sched.Clean(ctx); err != nil {
				slog.Error("failed clean events", "err", err)
				return
			}
			if err := sched.Notify(ctx); err != nil {
				slog.Error("failed send notifications", "err", err)
				return
			}
		}

	}
}
