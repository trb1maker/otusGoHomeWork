package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("invalid workersCount", func(t *testing.T) {
		err := Run([]Task{}, -1, 12)
		require.Error(t, err)

		err = Run([]Task{}, 0, 12)
		require.Error(t, err)
	})

	t.Run("invalid maxErrorsCount", func(t *testing.T) {
		err := Run([]Task{}, 1, -1)
		require.Error(t, err)

		err = Run([]Task{}, 1, 0)
		require.Error(t, err)
	})

	t.Run("empty task array", func(t *testing.T) {
		err := Run([]Task{}, 5, 1)
		require.NoError(t, err)
	})
}

// Тест не работает: я ожидаю, что успешной будет проверка runtime.NumGoroutine() >= workersCount,
// но запущены оказываются только 3 горутины. Вероятно это текущий поток, require.Eventually и Run.
// Ниже будет пример теста, который дает верный результат. Вероятно я не учитываю чего-то.

// TODO: вернуться к этому тесту позже.
// func TestEventually(t *testing.T) {
// 	tasksCount := 50
// 	tasks := make([]Task, 0, tasksCount)

// 	var runTasksCount int32

// 	for i := 0; i < tasksCount; i++ {
// 		tasks = append(tasks, func() error {
// 			atomic.AddInt32(&runTasksCount, 1)
// 			return nil
// 		})
// 	}

// 	workersCount := 5
// 	maxErrorsCount := 1

// 	err := Run(tasks, workersCount, maxErrorsCount)
// 	require.Eventually(t, func() bool {
// 		return runtime.NumGoroutine() >= 3
// 	}, time.Millisecond*500, time.Millisecond*50)
// 	require.NoError(t, err)
// 	require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
// }

// func TestRequireEventually(t *testing.T) {
// 	numWorkers := 5
// 	wg := &sync.WaitGroup{}
// 	wg.Add(numWorkers)

// 	for i := 0; i < numWorkers; i++ {
// 		go func() {
// 			defer wg.Done()
// 			time.Sleep(time.Millisecond * 500)
// 		}()
// 	}
// 	require.Eventually(t, func() bool {
// 		return runtime.NumGoroutine() >= numWorkers
// 	}, time.Millisecond*500, time.Millisecond*50)
// }
