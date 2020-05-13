package gomq

import (
	"sync"
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
	Publish(topic string, data interface{})

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

type queueMatcher struct {
	queue queue.Queue
	matcher Matcher
}

type broker struct {
	queueMatchers []queueMatcher
	sync.RWMutex
}

// NewBroker creates a new broker for message exchange.
func NewBroker() Broker {
	return &broker{
		queueMatchers: []queueMatcher{},
	}
}

func (b *broker) Publish(topic string, data interface{}) {
	b.RLock()
	defer b.RUnlock()

	// TODO: Cache the below information.
	for _, q := range b.queueMatchers {
		if q.matcher.MatchString(topic) {
			q.queue.Push(data)
		}
	}
}

func (b *broker) Subscribe(matcher Matcher) Poller {

	b.Lock()
	defer b.Unlock()

	que := queue.NewQueue()
	b.queueMatchers = append(b.queueMatchers, queueMatcher{queue:que, matcher: matcher})

	return que
}

func (b *broker) Close(timeOut time.Duration) {
	b.Lock()
	defer b.Unlock()
	for _, qm := range b.queueMatchers {
		qm.queue.Close(timeOut)
	}

	b.queueMatchers = []queueMatcher{}
}
