package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg          sync.WaitGroup
		err         error
		errorsCount int32
	)

	taskCh := make(chan Task)
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range taskCh {
				err := t()
				if err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			if m > 0 && atomic.LoadInt32(&errorsCount) >= int32(m) {
				err = ErrErrorsLimitExceeded
				break
			}
			taskCh <- task
		}
	}()

	wg.Wait()
	return err
}
