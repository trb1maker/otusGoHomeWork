package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Конфигурирую логгер
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// Парсинг аргументов командной строки и валидация необходимых аргументов
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		slog.Error("flag", "err", "not set address and port", "args", flag.Args())
		return
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// Запуск подключения
	c := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := c.Connect(); err != nil {
		slog.Error("Connect", "err", err)
		return
	}

	// Настройка реагирования приложения на сигналы пользователя / операционной системы
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, os.Interrupt, os.Kill)
	defer cancel()

	// Конкурентные чтение и запись
	go func() {
		if err := c.Receive(); err != nil {
			slog.Error("receive", "err", err)
		}
	}()

	go func() {
		if err := c.Send(); err != nil {
			slog.Error("send", "err", err)
		}
	}()

	<-ctx.Done()
}
