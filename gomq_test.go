package gomq

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestBrokerFanoutPattern(t *testing.T) {
	broker := NewBroker()
	defer broker.Close(0)

	subscriberCount := 10
	topicName := "all-events"
	valueToAssert := int32(12345)

	wg := sync.WaitGroup{}
	valueReceiveCount := int32(0)
	for i := 0; i < subscriberCount; i++ {
		wg.Add(1)

		// Subscribing before spawning the go routine to ensure all subscribed routines have been added to wait group.
		sub := broker.Subscribe(ExactMatcher(topicName))

		go func() {
			defer wg.Done()
			v, _ := sub.Poll()
			if v.(int32) != valueToAssert {
				t.Errorf("Invalid Value: Expected: %d Obtained: %v", valueToAssert, v)
			} else {
				atomic.AddInt32(&valueReceiveCount, 1)
			}
		}()
	}

	count := broker.Publish(topicName, valueToAssert)
	wg.Wait()

	if count != subscriberCount {
		t.Errorf("Missed delivery to few subscribers: Expected: %d Obtained: %v", count, subscriberCount)
	}

	if atomic.LoadInt32(&valueReceiveCount) != int32(subscriberCount) {
		t.Errorf("Invalid Subscriber's Receive Count: Expected: %d Obtained: %v", subscriberCount, valueReceiveCount)
	}

}

func TestBrokerFanoutPatternForAsyncBroker(t *testing.T) {
	broker := NewAsyncBroker()
	defer broker.Close(0)

	subscriberCount := 10
	topicName := "all-events"
	valueToAssert := int32(12345)

	wg := sync.WaitGroup{}
	valueReceiveCount := int32(0)
	for i := 0; i < subscriberCount; i++ {
		wg.Add(1)

		// Subscribing before spawning the go routine to ensure all subscribed routines have been added to wait group.
		sub := broker.Subscribe(ExactMatcher(topicName))

		go func() {
			defer wg.Done()
			v, _ := sub.Poll()
			if v.(int32) != valueToAssert {
				t.Errorf("Invalid Value: Expected: %d Obtained: %v", valueToAssert, v)
			} else {
				atomic.AddInt32(&valueReceiveCount, 1)
			}
		}()
	}

	count := broker.Publish(topicName, valueToAssert)
	wg.Wait()

	if count != subscriberCount {
		t.Errorf("Missed delivery to few subscribers: Expected: %d Obtained: %v", count, subscriberCount)
	}

	if atomic.LoadInt32(&valueReceiveCount) != int32(subscriberCount) {
		t.Errorf("Invalid Subscriber's Receive Count: Expected: %d Obtained: %v", subscriberCount, valueReceiveCount)
	}

}

func TestBrokerDataIntegritySingleRoutine(t *testing.T) {
	broker := NewBroker()
	defer broker.Close(0)

	poller := broker.Subscribe(ExactMatcher("all"))

	maxCount := 100

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for expected := 0; expected < maxCount; expected++ {
			val, ok := poller.Poll()
			if !ok {
				t.Errorf("Poll on available value should be True Got False")
			}
			if expected != val.(int) {
				t.Errorf("Invalid Value: Expected: %d Obtained: %v", expected, val)
			}
		}
	}()

	for i := 0; i < maxCount; i++ {
		broker.Publish("all", i)
	}

	wg.Wait()
}

func TestBrokerDataIntegrityMultiRoutine(t *testing.T) {
	broker := NewBroker()
	defer broker.Close(0)

	poller := broker.Subscribe(ExactMatcher("all"))

	maxCount := 100

	routineCount := 10

	wg := sync.WaitGroup{}

	for i := 0; i < routineCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			last := int32(-1)
			for i := 0; i < maxCount/routineCount; i++ {
				val, ok := poller.Poll()

				if !ok {
					t.Errorf("Poll on available value should be True Got False")
				}
				if last > val.(int32) {
					t.Errorf("Invalid Value: Last Value: %d Obtained: %v", last, val)
				}

				last = val.(int32)
			}
		}()
	}

	for i := 0; i < maxCount; i++ {
		broker.Publish("all", int32(i))
	}

	wg.Wait()
}

func TestBrokerPollAfterClose(t *testing.T) {
	broker := NewBroker()
	sub := broker.Subscribe(ExactMatcher("all"))
	broker.Publish("all", "record-1")

	{
		val, ok := sub.Poll()
		if val != "record-1" {
			t.Errorf("Expeted Value: record-1, Obtained: %v\n", val)
		}
		if !ok {
			t.Error("Poll should be True")
		}
	}

	broker.Close(-1)

	{
		_, ok := sub.Poll()
		if ok {
			t.Error("Poll should be False")
		}
	}

	{
		_, ok := sub.Poll()
		if ok {
			t.Error("Poll should be False")
		}
	}
}

func TestBrokerDataIntegrityCloseBroker(t *testing.T) {

	b := NewBroker()

	topics := []string{"topic1", "topic2", "topic3"}

	wg := sync.WaitGroup{}
	for _, topic := range topics {
		topic := topic
		poller := b.Subscribe(ExactMatcher(topic))
		wg.Add(1)
		go func() {
			defer wg.Done()
			incr := 0
			for val, ok := poller.Poll(); ok; val, ok = poller.Poll() {
				expected := fmt.Sprintf("%s:%v", topic, incr)
				if expected != val.(string) {
					t.Errorf("Invalid Value: Expected: %v Obtained: %v", expected, val)
				}
				incr++
			}
		}()
	}

	maxCount := 100
	for _, topic := range topics {
		for i := 0; i < maxCount; i++ {
			b.Publish(topic, fmt.Sprintf("%s:%v", topic, i))
		}
	}

	b.Close(-1)
	wg.Wait()
}
