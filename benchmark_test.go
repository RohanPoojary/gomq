package gomq

import (
	"fmt"
	// "fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
)

func benchmarkPublishNConsumer(b *testing.B, n int) {

	broker := NewBroker()
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
			broker.Publish("test" + id, rand.Intn(1000))
		}
	})
}


func BenchmarkPublish(b *testing.B) {
	// BenchmarkPublish/Consumers=1-4         	 1000000	      1178 ns/op
	// BenchmarkPublish/Consumers=10-4        	  300000	      4504 ns/op
	// BenchmarkPublish/Consumers=30-4        	  200000	     11523 ns/op
	// BenchmarkPublish/Consumers=50-4        	  100000	     19384 ns/op

	consumerTopicRatios := []int{1, 10, 30, 50}

	for _, consumerCount := range consumerTopicRatios {
		name := fmt.Sprintf("Consumers=%d", consumerCount)
		b.Run(name, func(b *testing.B) {
			benchmarkPublishNConsumer(b, consumerCount)
		})
	}
}

func benchmarkConsumeNPublisher(b *testing.B, n int) {

	broker := NewBroker()
	defer broker.Close(0)

	stop := make(chan bool, 0)

	sub := broker.Subscribe(regexp.MustCompile(`test\.*`))

	for i := 0; i < n; i++ {
		i := i
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					broker.Publish("test:" + strconv.Itoa(i), rand.Intn(1000))
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

	stop <-true
	close(stop)
}

func benchmarkMConsumerNPublisher(b *testing.B, m, n int) {

	broker := NewBroker()
	defer broker.Close(0)

	stop := make(chan bool, 0)
	subs := make([]Poller, m)
	for i := 0; i < m; i++ {
		subs[i] = broker.Subscribe(regexp.MustCompile(`test\.*`))
	}

	for i := 0; i < n; i++ {
		i := i
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					broker.Publish("test:" + strconv.Itoa(i), rand.Intn(1000))
				}
			}
		}()
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		sub := subs[rand.Intn(m)]
		for pb.Next() {
			sub.Poll()
		}
	})

	stop <-true
	close(stop)
}


func BenchmarkConsume(b *testing.B) {
	// BenchmarkConsume/Publisher=1-4         	 1000000	      1599 ns/op
	// BenchmarkConsume/Publisher=10-4        	 1000000	      1788 ns/op
	// BenchmarkConsume/Publisher=30-4        	 1000000	      2182 ns/op
	// BenchmarkConsume/Publisher=50-4        	 1000000	      2330 ns/op

	publisherCounts := []int{1, 10, 30, 50}

	for _, publisherCount := range publisherCounts {
		name := fmt.Sprintf("Publisher=%d", publisherCount)
		b.Run(name, func(b *testing.B) {
			benchmarkConsumeNPublisher(b, publisherCount)
		})
	}
}


func BenchmarkMultiConsume(b *testing.B) {
	// BenchmarkMultiConsume/Publisher=50/Consumer=10-4         	 1000000	      1134 ns/op
	// BenchmarkMultiConsume/Publisher=50/Consumer=30-4         	  500000	      3568 ns/op
	// BenchmarkMultiConsume/Publisher=50/Consumer=50-4         	  500000	      4416 ns/op

	consumerCounts := []int{10, 30, 50}
	publisherCount := 50
	for _, consumerCount := range consumerCounts {
		name := fmt.Sprintf("Publisher=%d/Consumer=%d", publisherCount, consumerCount)
		b.Run(name, func(b *testing.B) {
			benchmarkMConsumerNPublisher(b, consumerCount, publisherCount)
		})
	}
}
