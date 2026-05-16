package usecase

import (
	"context"
	"sync"
)

// Event represents an IPC event.
type Event struct {
	Topic   string
	Payload []byte
}

// EventHandler is a function that handles an Event.
type EventHandler func(ctx context.Context, event Event)

// EventBus is a non-blocking Pub/Sub router.
type EventBus interface {
	Subscribe(topic string, handler EventHandler)
	Publish(topic string, payload []byte)
}

// MemoryEventBus implements the EventBus interface in memory.
type MemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]EventHandler
}

// NewMemoryEventBus creates a new MemoryEventBus.
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe adds an EventHandler for the given topic.
func (b *MemoryEventBus) Subscribe(topic string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], handler)
}

// Publish executes all registered EventHandlers for the given topic asynchronously.
func (b *MemoryEventBus) Publish(topic string, payload []byte) {
	b.mu.RLock()
	handlers, ok := b.subscribers[topic]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, handler := range handlers {
		// Execute handlers asynchronously so publishing never blocks
		go handler(context.Background(), Event{Topic: topic, Payload: payload})
	}
}
