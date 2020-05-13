package queue

import (
	"fmt"
	"sync"
)

func Example_example1() {

	// The queue has no size limit. Hence, data can be pushed
	// even without any poller.
	queue := NewQueue()
	for i := 1; i <= 5; i++ {
		queue.Push(i)
	}

	queue.Close(-1)

	// Polls data from the queue. Once polled the data will be evicted from queue.
	for val, ok := queue.Poll(); ok; val, ok = queue.Poll() {
		fmt.Println(val)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func Example_example2() {

	queue := NewQueue()
	wg := sync.WaitGroup{}

	wg.Add(1)

	// Queue being thread safe. The push & poll can be done by any routines.
	go func() {
		defer wg.Done()

		for i := 1; i <= 100; i++ {
			queue.Push(i)
		}

		queue.Close(-1)
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for val, ok := queue.Poll(); ok; val, ok = queue.Poll() {
				fmt.Println("Value Obtained:", val)
			}
		}()
	}

	wg.Wait()
}

func ExampleQueue() {
	queue := NewQueue()

	for i := 1; i <= 3; i++ {
		queue.Push(i)
	}

	// The queue being similar to go channels should be closed once used.
	// Else Poll will be blocked indefinitely if queue gets empty.
	queue.Close(-1)

	fmt.Println(queue.Poll())
	fmt.Println(queue.Poll())
	fmt.Println(queue.Poll())
	fmt.Println(queue.Poll())
	fmt.Println(queue.Poll())

	// Output:
	// 1 true
	// 2 true
	// 3 true
	// <nil> false
	// <nil> false
}
