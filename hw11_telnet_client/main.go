package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
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
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// Запуск подключения
	c := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := c.Connect(); err != nil {
		l.Fatalf("connect: %v", err)
	}

	// Настройка реагирования приложения на сигналы пользователя / операционной системы
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// Конкурентные чтение и запись
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := c.Close(); err != nil {
			l.Printf("Close: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := c.Receive(); err != nil {
			l.Printf("receive: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := c.Send(); err != nil {
			l.Printf("send: %v", err)
		}
	}()

	wg.Wait()
}
