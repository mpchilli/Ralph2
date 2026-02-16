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
