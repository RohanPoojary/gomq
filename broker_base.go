package gomq

import (
	"sync"
	"time"

	"github.com/RohanPoojary/gomq/queue"
)

type queueMatcher struct {
	queue   queue.Queue
	matcher Matcher
}

type brokerBase struct {
	queueMatchers []queueMatcher
	sync.RWMutex
}

func (b *brokerBase) Subscribe(matcher Matcher) Poller {

	b.Lock()
	defer b.Unlock()

	que := queue.NewQueue()
	b.queueMatchers = append(b.queueMatchers, queueMatcher{queue: que, matcher: matcher})

	return que
}

func (b *brokerBase) Close(timeOut time.Duration) {
	b.Lock()
	defer b.Unlock()

	b.unsafeClose(timeOut)
}

func (b *brokerBase) unsafeClose(timeOut time.Duration) {

	wg := sync.WaitGroup{}

	for _, qm := range b.queueMatchers {
		qm := qm
		wg.Add(1)
		go func() {
			defer wg.Done()
			qm.queue.Close(timeOut)
		}()
	}

	wg.Wait() // Wait until all subscribers are closed.

	b.queueMatchers = []queueMatcher{}
}
