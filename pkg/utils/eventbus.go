package utils

import (
	"sync"
)

type Event struct {
	Topic   string
	Payload interface{}
}

type EventBus struct {
	subscribers map[string][]chan Event
	mu          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Event),
	}
}

func (eb *EventBus) Subscribe(topic string) <-chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	ch := make(chan Event, 10)
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	return ch
}

func (eb *EventBus) Unsubscribe(ch <-chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	for topic, chans := range eb.subscribers {
		for i, subCh := range chans {
			if subCh == ch {
				// Remove channel from slice
				eb.subscribers[topic] = append(chans[:i], chans[i+1:]...)
				// Consider if we should close here. 
				// Closing might be safer for range loops on the consumer side.
				// However, since Subscribe returns <-chan, we need to cast to chan to close?
				// Actually, the stored type is `chan Event`, so we can close it.
				close(subCh)
				return
			}
		}
	}
}

func (eb *EventBus) Publish(topic string, payload interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	
	if chans, found := eb.subscribers[topic]; found {
		event := Event{Topic: topic, Payload: payload}
		for _, ch := range chans {
			// Non-blocking publish to avoid deadlocks
			select {
			case ch <- event:
			default:
				// Drop event if channel full
			}
		}
	}
}
