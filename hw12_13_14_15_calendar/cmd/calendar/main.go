package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/exp/slog"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := loadConfig(configFile)
	if err != nil {
		slog.Error("config", "err", err)
		return
	}

	var storage app.Storage
	if config.Storage.Type == "postgres" {
		storage = sqlstorage.New()
	} else {
		storage = memorystorage.New()
	}

	calendar := app.New(storage)
	server := internalhttp.NewServer(calendar)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			slog.Error("failed to stop http server", "err", err)
		}
	}()

	slog.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		cancel()
		slog.Error("failed to start http server", "err", err)
		return
	}
}
