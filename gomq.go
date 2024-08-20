package gomq

import (
	"time"

	"github.com/RohanPoojary/gomq/queue"
)

// Poller is the interface that wraps Poll function.
type Poller interface {

	// Poll reads the data from queue and returns the value.
	// This is blocking call until there is consumable data.
	// In which case, Poll will return the data and True.
	//
	// If the resource is closed, then Poll will return,
	// nil and False
	Poll() (interface{}, bool)
}

// Broker represents the Broker for interaction.
type Broker interface {

	// Publish publishes the `data` to the topic.
	// It returns the count of matched subscribers to which the `data` has been published.
	Publish(topic string, data interface{}) int

	// Subscribe creates a Poller which polls data
	// from matched topics.
	Subscribe(topic Matcher) Poller

	// Close closes the Broker and renders it read only.
	// Hence, all data pushed will be ignored.
	// All the open resources will be collected based on timeOut.
	//
	// If timeOut < 0, then resources will be closed once
	// all the queues are empty.
	// For any timeOut >= 0, all the resources will be force closed
	// after timeOut.
	Close(timeOut time.Duration)
}

// NewBroker creates a new broker for message exchange.
// A simple broker which synchronously publishes the data to all its matching subscribers.
func NewBroker() Broker {
	return &broker{
		brokerBase: brokerBase{
			queueMatchers: []queueMatcher{},
		},
	}
}

type broker struct {
	brokerBase
}

func (b *broker) Publish(topic string, data interface{}) int {
	b.RLock()
	defer b.RUnlock()

	count := 0
	for _, q := range b.queueMatchers {
		if q.matcher.MatchString(topic) {
			q.queue.Push(data)
			count += 1
		}
	}

	return count
}

func NewAsyncBroker() Broker {
	b := &asyncBroker{
		queue: queue.NewQueue(),
		brokerBase: brokerBase{
			queueMatchers: []queueMatcher{},
		},
	}

	go b.manage()

	return b
}

type asyncBroker struct {
	brokerBase
	queue queue.Queue
}

type asyncPayload struct {
	topic string
	data  interface{}
}

func (b *asyncBroker) Publish(topic string, data interface{}) int {
	b.RLock()
	defer b.RUnlock()

	minMatchCount := len(b.queueMatchers)

	if minMatchCount == 0 {
		return 0
	}

	b.queue.Push(asyncPayload{topic: topic, data: data})
	return minMatchCount
}

func (b *asyncBroker) manage() int {
	for {
		val, ok := b.queue.Poll()
		if !ok {
			return 0
		}

		payload := val.(asyncPayload)
		b.publish(payload.topic, payload.data)
	}
}

func (b *asyncBroker) publish(topic string, data interface{}) int {
	b.RLock()
	defer b.RUnlock()

	count := 0
	for _, q := range b.queueMatchers {
		if q.matcher.MatchString(topic) {
			q.queue.Push(data)
			count += 1
		}
	}

	return count
}

func (b *asyncBroker) Close(timeOut time.Duration) {

	b.Lock()
	defer b.Unlock()

	b.queue.Close(timeOut)
	b.brokerBase.unsafeClose(timeOut)
}
