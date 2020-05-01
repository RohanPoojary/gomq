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

type broker struct {
	queues   []queue.Queue
	patterns []Matcher
	sync.RWMutex
}

func NewBroker() Broker {
	return &broker{
		queues:   make([]queue.Queue, 0),
		patterns: make([]Matcher, 0),
	}
}

func (b *broker) Publish(topic string, data interface{}) {
	b.RLock()
	defer b.RUnlock()

	// TODO: Cache the below information.
	for i, q := range b.queues {
		if b.patterns[i].MatchString(topic) {
			q.Push(data)
		}
	}
}

func (b *broker) Subscribe(matcher Matcher) Poller {

	b.Lock()
	defer b.Unlock()

	que := queue.NewQueue()
	b.queues = append(b.queues, que)
	b.patterns = append(b.patterns, matcher)

	return que
}

func (b *broker) Close(timeOut time.Duration) {
	b.Lock()
	defer b.Unlock()
	for i := 0; i < len(b.queues); i++ {
		q := b.queues[i]
		q.Close(timeOut)
	}

	b.queues = []queue.Queue{}
	b.patterns = []Matcher{}
}
