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

// waitErrors запускает taskChannel и набор worker, возвращает канал ошибок.
func waitErrors(ctx context.Context, tt []Task, workersCount int) <-chan error {
	out := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(workersCount)
	go func() {
		defer close(out)
		tasks := taskChannel(ctx, tt)
		for i := 0; i < workersCount; i++ {
			go worker(ctx, tasks, out, wg)
		}
		wg.Wait()
	}()
	return out
}

// taskChannel отправляет задачи в канал пока не получит сигнал завершить работу или пока вся работа не будет завершена.
func taskChannel(ctx context.Context, tt []Task) <-chan Task {
	out := make(chan Task)
	go func() {
		defer close(out)
		for _, t := range tt {
			select {
			case <-ctx.Done():
				return
			default:
				out <- t
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
