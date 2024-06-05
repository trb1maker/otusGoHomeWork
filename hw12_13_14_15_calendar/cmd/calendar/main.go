package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/app"
	internalgrpc "github.com/trb1maker/otus_golang_home_work/hw12_13_14_15_calendar/internal/server/grpc"
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
		slog.Error("failed to load config", "err", err)
		return
	}

	var storage StorageConnectClose
	if config.Storage.Type == "postgres" {
		storage, err = sqlstorage.New(
			config.Storage.Postgres.Host,
			config.Storage.Postgres.Port,
			config.Storage.Postgres.Database,
			config.Storage.Postgres.User,
			config.Storage.Postgres.Password,
		)
		if err != nil {
			slog.Error("failed to connect to postgres", "err", err)
			return
		}
	} else {
		storage, err = memorystorage.New()
		if err != nil {
			slog.Error("failed to connect to memory storage", "err", err)
			return
		}
	}

	if err := storage.Connect(context.Background()); err != nil {
		slog.Error("failed to connect to storage", "err", err)
		return
	}
	defer storage.Close(context.Background())

	calendar := app.New(storage)

	httpServer := internalhttp.NewServer(calendar, config.Server.HTTP.Host, config.Server.HTTP.Port)
	grpcServer := internalgrpc.NewServer(calendar, config.Server.GRPC.Host, config.Server.GRPC.Port)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Stop(ctx); err != nil {
			slog.Error("failed to stop http server", "err", err)
		}

		grpcServer.Stop()
	}()

	slog.Info("calendar is running...")

	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			cancel()
			slog.Error("failed to start http server", "err", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			cancel()
			slog.Error("failed to start grpc server", "err", err)
			return
		}
	}()

	wg.Wait()
}

type StorageConnectClose interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	app.Storage
}
