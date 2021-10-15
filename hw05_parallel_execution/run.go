package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type producer struct {
	taskCh *chan Task
	doneCh *chan bool
}

type consumer struct {
	taskCh      *chan Task
	doneCh      *chan bool
	resultCh    *chan error
	result      error
	mu          sync.Mutex
	errorsCount int
	errorsLimit int
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task)
	doneCh := make(chan bool)
	resultCh := make(chan error, n)
	defer close(resultCh)

	go newProducer(&taskCh, &doneCh).produce(tasks)
	go newConsumer(&taskCh, &doneCh, &resultCh, m).consume(n)

	return <-resultCh
}

func newProducer(taskCh *chan Task, doneCh *chan bool) *producer {
	return &producer{
		taskCh: taskCh,
		doneCh: doneCh,
	}
}

func (p *producer) produce(tasks []Task) {
	defer close(*p.doneCh)
	defer close(*p.taskCh)

	taskKey := 0
	for taskKey < len(tasks) {
		select {
		case *p.taskCh <- tasks[taskKey]:
			taskKey++
		case <-*p.doneCh:
			return
		}
	}
}

func newConsumer(taskCh *chan Task, doneCh *chan bool, resultCh *chan error, errorsLimit int) *consumer {
	return &consumer{
		taskCh:      taskCh,
		doneCh:      doneCh,
		resultCh:    resultCh,
		errorsLimit: errorsLimit,
	}
}

func (c *consumer) consume(n int) {
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range *c.taskCh {
				c.doTask(t)
			}
		}()
	}

	wg.Wait()
	*c.resultCh <- c.result
}

func (c *consumer) doTask(t Task) {
	err := t()
	if c.errorsLimit > 0 && err != nil {
		c.mu.Lock()
		c.errorsCount++
		if c.errorsCount == c.errorsLimit {
			*c.doneCh <- true
			c.result = ErrErrorsLimitExceeded
		}
		c.mu.Unlock()
	}
}
