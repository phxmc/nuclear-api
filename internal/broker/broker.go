package broker

import "sync"

type Broker[T any] struct {
	subscribers map[chan T]any
	mu          sync.RWMutex
}

func New[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: make(map[chan T]any),
	}
}

func (broker *Broker[T]) Subscribe() chan T {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	message := make(chan T, 128)
	broker.subscribers[message] = nil
	return message
}

func (broker *Broker[T]) Unsubscribe(message chan T) {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	if _, ok := broker.subscribers[message]; ok {
		delete(broker.subscribers, message)
		close(message)
	}
}

func (broker *Broker[T]) Publish(message T) {
	broker.mu.RLock()
	defer broker.mu.RUnlock()

	for subscriber := range broker.subscribers {
		select {
		case subscriber <- message:
		default:
		}
	}
}

func (broker *Broker[T]) Shutdown() {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	for subscriber := range broker.subscribers {
		delete(broker.subscribers, subscriber)
		close(subscriber)
	}
}
