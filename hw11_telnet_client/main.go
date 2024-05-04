package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Конфигурирую логгер
	l := log.New(os.Stderr, "go-telnet", log.Ltime)

	// Парсинг аргументов командной строки и валидация необходимых аргументов
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) != 2 {
		l.Fatal("not set address and port")
		return
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// Запуск подключения
	c := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := c.Connect(); err != nil {
		l.Fatalf("connect: %v", err)
		return
	}

	// Настройка реагирования приложения на сигналы пользователя / операционной системы
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, os.Interrupt, os.Kill)
	defer cancel()

	// Конкурентные чтение и запись
	go func() {
		if err := c.Receive(); err != nil {
			l.Fatalf("receive: %v", err)
		}
	}()

	go func() {
		if err := c.Send(); err != nil {
			l.Fatalf("send: %v", err)
		}
	}()

	<-ctx.Done()
}
