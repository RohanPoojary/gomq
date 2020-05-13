package queue

import (
	"math/rand"
	"testing"
)

func BenchmarkPushWithoutPoll(b *testing.B) {

	// > go test -cpu 1,2,4 -bench  . -benchmem
	// goos: darwin
	// goarch: amd64
	// pkg: github.com/RohanPoojary/gomq/queue
	// BenchmarkPushWithoutPoll     	 3000000	       587 ns/op	      98 B/op	       0 allocs/op
	// BenchmarkPushWithoutPoll-2   	 3000000	       597 ns/op	      98 B/op	       0 allocs/op
	// BenchmarkPushWithoutPoll-4   	 2000000	       652 ns/op	      94 B/op	       0 allocs/op
	// PASS

	queue := NewQueue()
	defer queue.Close(0)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.Push(rand.Intn(10000))
		}
	})
}

func BenchmarkBothPushPollWithPrefilledQueue(b *testing.B) {

	// > go test -cpu 1,2,4 -bench  . -benchmem
	// goos: darwin
	// goarch: amd64
	// pkg: github.com/RohanPoojary/gomq/queue
	// BenchmarkBothPushPollWithPrefilledQueue     	 2000000	       597 ns/op	      91 B/op	       0 allocs/op
	// BenchmarkBothPushPollWithPrefilledQueue-2   	 2000000	       686 ns/op	      91 B/op	       0 allocs/op
	// BenchmarkBothPushPollWithPrefilledQueue-4   	 2000000	       682 ns/op	      91 B/op	       0 allocs/op
	// PASS

	queue := NewQueue()

	defer queue.Close(0)

	for i := 0; i < 10000; i++ {
		queue.Push(rand.Intn(10000))
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(10000)
			queue.Push(n)
			if n > 8000 {
				queue.Poll()
			}
		}
	})
}

func BenchmarkPollWithAsyncPublish(b *testing.B) {

	// > go test -cpu 1,2,4 -bench  . -benchmem
	// goos: darwin
	// goarch: amd64
	// pkg: github.com/RohanPoojary/gomq/queue
	// BenchmarkPollWithAsyncPublish     	 2000000	       808 ns/op	      94 B/op	       1 allocs/op
	// BenchmarkPollWithAsyncPublish-2   	 2000000	       975 ns/op	      36 B/op	       1 allocs/op
	// BenchmarkPollWithAsyncPublish-4   	 2000000	       943 ns/op	      31 B/op	       1 allocs/op
	// PASS

	queue := NewQueue()
	defer queue.Close(0)

	closeCh := make(chan bool)
	go func() {
		for  {
			select {
			case <-closeCh:
				return
			default:
				queue.Push(rand.Intn(100000))
			}
		}
	}()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			queue.Poll()
		}
	})

	closeCh<-true
}
