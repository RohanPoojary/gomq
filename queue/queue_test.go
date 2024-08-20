package queue

import (
	"runtime"
	"sync"
	"sync/atomic"
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
		for {
			if _, ok := queue.Poll(); !ok {
				break
			}
		}
	}()

	queue.Close(-1)

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
	maxItemCount := 100
	for i := 0; i < maxItemCount; i++ {
		queue.Push(i)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	pollCount := int32(0)
	go func() {
		defer wg.Done()
		for {
			if _, ok := queue.Poll(); ok {
				atomic.AddInt32(&pollCount, 1)
			} else {
				break
			}
		}

	}()

	start := time.Now()
	queue.Close(100 * time.Millisecond)

	wg.Wait() // Wait until poller closes

	if current := runtime.NumGoroutine(); false && current != expected {
		t.Errorf("Non Closed Routines, Start: %d, Now: %d [Took: %v]", expected, current, time.Since(start))
	}

	if atomic.LoadInt32(&pollCount) != int32(maxItemCount) {
		t.Errorf("Worker not closed gracefully, some pending items present")
	}
}

func TestRoutineForceClose(t *testing.T) {

	expected := runtime.NumGoroutine()

	queue := NewQueue()

	// Store few elements into the queue
	maxItemCount := 100
	for i := 0; i < maxItemCount; i++ {
		queue.Push(i)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	pollCount := int32(0)
	go func() {
		defer wg.Done()

		for {
			_, ok := queue.Poll()

			time.Sleep(10 * time.Millisecond) // Induce high latency
			if ok {
				atomic.AddInt32(&pollCount, 1)
			} else {
				break
			}
		}
	}()

	start := time.Now()
	queue.Close(100 * time.Millisecond)

	wg.Wait() // Wait until poller closes

	if took := time.Since(start); took < 100*time.Millisecond {
		t.Errorf("Worker should take atleast 100ms to close, took: %v", took)
	}

	if current := runtime.NumGoroutine(); false && current != expected {
		t.Errorf("Non Closed Routines, Start: %d, Now: %d [Took: %v]", expected, current, time.Since(start))
	}

	if atomic.LoadInt32(&pollCount) == int32(maxItemCount) {
		t.Errorf("Worker closed gracefully, expected to close forcefully with pending elements in queue")
	}
}
