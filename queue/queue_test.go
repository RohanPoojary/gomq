package queue

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestQueueAsyncPoll(t *testing.T) {
	queue := NewQueue()

	wg := sync.WaitGroup{}

	wg.Add(1)
	lastVal := 0

	go func() {
		defer wg.Done()

		for v, ok := queue.Poll(); ok; v, ok = queue.Poll() {
			if v != lastVal {
				t.Errorf("Invalid Value: Last: %v, Current: %v\n", lastVal, v)
			}
			lastVal++
		}
	}()

	maxValue := 100
	for i := 0; i < maxValue; i++ {
		queue.Push(i)
	}

	queue.Close(-1)

	wg.Wait()

	if lastVal != 100 {
		t.Errorf("Invalid Last Value: %v\n", lastVal)
	}
}

func TestQueueAsyncPush(t *testing.T) {
	queue := NewQueue()

	go func() {
		maxValue := 100
		for i := 0; i < maxValue; i++ {
			queue.Push(i)
		}
		queue.Close(-1)
	}()

	lastVal := 0
	for v, ok := queue.Poll(); ok; v, ok = queue.Poll() {
		if v != lastVal {
			t.Errorf("Invalid Value: Last: %v, Current: %v\n", lastVal, v)
		}
		lastVal++
	}

	if lastVal != 100 {
		t.Errorf("Invalid Last Value: %v\n", lastVal)
	}
}

func TestQueuePushPollSequential(t *testing.T) {
	queue := NewQueue()

	maxValue := 100
	for i := 0; i < maxValue; i++ {
		queue.Push(i)
	}
	queue.Close(-1)

	lastVal := 0
	for v, ok := queue.Poll(); ok; v, ok = queue.Poll() {
		if v != lastVal {
			t.Errorf("Invalid Value: Last: %v, Current: %v\n", lastVal, v)
		}
		lastVal++
	}

	if lastVal != 100 {
		t.Errorf("Invalid Last Value: %v\n", lastVal)
	}
}

func TestSafeRoutineClose(t *testing.T) {
	expected := runtime.NumGoroutine()
	queue := NewQueue()

	for i := 0; i < 10; i++ {
		queue.Push(i)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		queue.Poll()
	}()

	queue.Close(0)

	wg.Wait()
	time.Sleep(time.Millisecond)

	if current := runtime.NumGoroutine(); current != expected {
		t.Errorf("Non Closed Routines. Start: %d, Now: %d", expected, current)
	}
}

func TestRoutineClose(t *testing.T) {

	expected := runtime.NumGoroutine()
	queue := NewQueue()

	// Store few elements into the queue
	for i := 0; i < 100; i++ {
		queue.Push(i)
	}

	workerExited := false
	go func() {
		for {
			_, ok := queue.Poll()
			if !ok {
				break
			}
		}

		workerExited = true
	}()

	queue.Close(100 * time.Millisecond)

	if current := runtime.NumGoroutine(); current != expected {
		t.Errorf("Non Closed Routines. Start: %d, Now: %d", expected, current)
	}

	if !workerExited {
		t.Errorf("Worker did not exit")
	}
}

func TestRoutineForceClose(t *testing.T) {

	queue := NewQueue()

	// Store few elements into the queue
	for i := 0; i < 100; i++ {
		queue.Push(i)
	}

	workerExited := false
	go func() {
		for {
			_, ok := queue.Poll()

			time.Sleep(time.Second) // Induce high latency
			if !ok {
				break
			}
		}

		workerExited = true
	}()

	start := time.Now()
	queue.Close(100 * time.Millisecond)

	if took := time.Since(start); took < 100*time.Millisecond {
		t.Errorf("Worker should take atleast 100ms to close, took: %v", took)
	}

	if workerExited {
		t.Errorf("Worker was not closed forcefully")
	}

}
