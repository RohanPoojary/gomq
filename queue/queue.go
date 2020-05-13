package queue

import (
	"sync"
	"time"
)

// Queue provides thread safe queue functions.
type Queue interface {

	// Push pushes value to the queue.
	Push(value interface{})

	// Poll pops the top most element of queue.
	// If not data is present in the queue, then ok will be false.
	//
	// This is blocking call,
	// Hence, it will wait till a queue is non empty.
	// In case of closed queue, Ok will be false.
	Poll() (value interface{}, ok bool)

	// Close closes the queue for any write operations.
	//
	// For negative timeOut, resources will be closed once all the data are polled,
	// else the resources will be forcefully collected after timeOut.
	Close(timeOut time.Duration)
}

type queue struct {
	in    chan interface{}
	out   chan interface{}
	close chan bool
	once  sync.Once
}

// NewQueue creates a new thread safe queue.
func NewQueue() Queue {
	q := queue{
		in:    make(chan interface{}, 1),
		out:   make(chan interface{}, 1),
		close: make(chan bool, 1),
	}
	go q.manage()

	return &q
}

func (q *queue) manage() {
	queue := []interface{}{}
	defer close(q.out)

	for {
		if len(queue) == 0 {
			select {
			case <-q.close:
				return
			case v, ok := <-q.in:
				if !ok {
					return
				}
				queue = append(queue, v)
			}
		} else {
			select {
			case <-q.close:
				return
			case v, ok := <-q.in:
				if ok {
					queue = append(queue, v)
				}
			case q.out <- queue[0]:
				queue[0] = nil
				queue = queue[1:]
			}
		}
	}
}

func (q *queue) Push(value interface{}) {
	q.in <- value
}

func (q *queue) Poll() (interface{}, bool) {
	val, ok := <-q.out
	return val, ok
}

func (q *queue) forceClose() {
	q.close <- true
	close(q.close)
}

func (q *queue) Close(timeout time.Duration) {
	q.once.Do(func() {
		close(q.in)
		if timeout >= 0 {
			go func() {
				select {
				case <-time.After(timeout):
					q.forceClose()
				}
			}()
		}
	})
}
