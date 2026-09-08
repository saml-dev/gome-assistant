package metrics

import "sync"

type broadcaster struct {
	mu          sync.Mutex
	subscribers map[chan struct{}]struct{}
	closed      bool
}

func newBroadcaster() *broadcaster {
	return &broadcaster{subscribers: make(map[chan struct{}]struct{})}
}

func (broadcaster *broadcaster) subscribe() (<-chan struct{}, func()) {
	changes := make(chan struct{}, 1)
	broadcaster.mu.Lock()
	if broadcaster.closed {
		close(changes)
		broadcaster.mu.Unlock()
		return changes, func() {}
	}
	broadcaster.subscribers[changes] = struct{}{}
	broadcaster.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			broadcaster.mu.Lock()
			if _, exists := broadcaster.subscribers[changes]; exists {
				delete(broadcaster.subscribers, changes)
				close(changes)
			}
			broadcaster.mu.Unlock()
		})
	}
	return changes, cancel
}

func (broadcaster *broadcaster) publish() {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	for subscriber := range broadcaster.subscribers {
		select {
		case subscriber <- struct{}{}:
		default:
		}
	}
}

func (broadcaster *broadcaster) close() {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	if broadcaster.closed {
		return
	}
	broadcaster.closed = true
	for subscriber := range broadcaster.subscribers {
		close(subscriber)
		delete(broadcaster.subscribers, subscriber)
	}
}
