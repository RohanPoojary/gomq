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
