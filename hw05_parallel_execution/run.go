package hw05parallelexecution

import (
	"errors"
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
	return nil
}
