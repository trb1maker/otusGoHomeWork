package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded          = errors.New("errors limit exceeded")
	ErrWorkersCountIsNegativeOrZero = errors.New("workers count is negative or zero")
)

type Task func() error

func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	if workersCount <= 0 {
		return ErrWorkersCountIsNegativeOrZero
	}
	if maxErrorsCount <= 0 {
		return ErrErrorsLimitExceeded
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var err error
	for range waitErrors(ctx, tasks, workersCount) {
		maxErrorsCount--
		if maxErrorsCount == 0 {
			cancel()
			err = ErrErrorsLimitExceeded
		}
	}
	return err
}

// waitErrors запускает набор воркеров и передает им задачи, если не поступит сигнал завершить работу,
// возвращает канал ошибок.
func waitErrors(ctx context.Context, tt []Task, workersCount int) <-chan error {
	out := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(workersCount)
	go func() {
		defer close(out)
		defer wg.Wait() // Жду когда воркеры завершат работу, чтобы они не попытались записать в закрытый канал

		// Запускаю воркеры
		tasks := make(chan Task)
		defer close(tasks)
		for i := 0; i < workersCount; i++ {
			go worker(ctx, tasks, out, wg)
		}

		// Направляю задачи на выполнение
		for _, t := range tt {
			select {
			case <-ctx.Done():
				return
			default:
				tasks <- t
			}
		}
	}()
	return out
}

// worker получает задачи из канала и выполняет их, если не поступил сигнал завершить работу.
func worker(ctx context.Context, tt <-chan Task, errs chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tt {
		select {
		case <-ctx.Done():
			return
		default:
			if err := t(); err != nil {
				errs <- err
			}
		}
	}
}
