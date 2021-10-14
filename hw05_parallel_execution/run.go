package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var mu sync.Mutex
	taskCh := make(chan Task)
	doneCh := make(chan bool)
	errors := 0
	var result error

	var wg sync.WaitGroup
	wg.Add(n + 1)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-doneCh:
					return
				case t, ok := <-taskCh:
					if !ok {
						return
					}
					err := t()
					if err != nil {
						mu.Lock()
						errors++
						if errors == m {
							doneCh <- true
							result = ErrErrorsLimitExceeded
						}
						mu.Unlock()
					}
				}
			}
		}()
	}

	go func() {
		defer wg.Done()
		defer close(taskCh)
		defer close(doneCh)

		taskKey := 0
		for taskKey < len(tasks) {
			select {
			case taskCh <- tasks[taskKey]:
				taskKey++
			case <-doneCh:
				return
			}
		}
	}()

	wg.Wait()

	return result
}
