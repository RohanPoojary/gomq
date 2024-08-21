package gomq

import (
	"fmt"
	// "fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
)

func benchmarkPublishNPoller(b *testing.B, creator func() Broker, n int) {

	broker := creator()
	defer broker.Close(0)

	reg := regexp.MustCompile(`test\.*`)

	for i := 0; i < n; i++ {
		go func() {
			sub := broker.Subscribe(reg)
			for _, ok := sub.Poll(); ok; _, ok = sub.Poll() {
			}
		}()
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := strconv.Itoa(100)
			broker.Publish("test"+id, rand.Intn(1000))
		}
	})
}

func benchmarkPollNPublisher(b *testing.B, creator func() Broker, n int) {

	broker := creator()
	defer broker.Close(0)

	stop := make(chan bool)

	sub := broker.Subscribe(regexp.MustCompile(`test\.*`))

	for i := 0; i < n; i++ {
		i := i
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					broker.Publish("test:"+strconv.Itoa(i), rand.Intn(1000))
				}
			}
		}()
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sub.Poll()
		}
	})

	stop <- true
	close(stop)
}

func BenchmarkPublish(b *testing.B) {

	pollerCounts := []int{1, 10, 30, 50}

	for _, pollerCount := range pollerCounts {
		name := fmt.Sprintf("SyncBroker/Pollers=%d", pollerCount)
		b.Run(name, func(b *testing.B) {
			benchmarkPublishNPoller(b, NewBroker, pollerCount)
		})
	}

	for _, pollerCount := range pollerCounts {
		name := fmt.Sprintf("AsyncBroker/Pollers=%d", pollerCount)
		b.Run(name, func(b *testing.B) {
			benchmarkPublishNPoller(b, NewAsyncBroker, pollerCount)
		})
	}

}

func BenchmarkPoll(b *testing.B) {

	publisherCounts := []int{1, 10, 30, 50}

	for _, publisherCount := range publisherCounts {
		name := fmt.Sprintf("SyncBroker/Publisher=%d", publisherCount)
		b.Run(name, func(b *testing.B) {
			benchmarkPollNPublisher(b, NewBroker, publisherCount)
		})
	}

	for _, publisherCount := range publisherCounts {
		name := fmt.Sprintf("AsyncBroker/Publisher=%d", publisherCount)
		b.Run(name, func(b *testing.B) {
			benchmarkPollNPublisher(b, NewAsyncBroker, publisherCount)
		})
	}
}
