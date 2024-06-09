package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/queue"
	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/sender"
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

	queue, err := queue.New(
		config.Rabbit.Host,
		config.Rabbit.Port,
		config.Rabbit.User,
		config.Rabbit.Password,
	)
	defer queue.Close()

	snd := sender.New(queue)

	interval, err := time.ParseDuration(config.Interval)
	if err != nil {
		slog.Error("failed to parse interval", "err", err)
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := snd.Schedule(ctx); err != nil {
				slog.Error("failed to send notification", "err", err)
				return
			}
		}

	}
}
