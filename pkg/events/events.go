package events

import "fmt"

type Name string

type Listener[T any] func(T)

type Event[T any] struct {
	listeners map[string]Listener[T]
	name      Name
}

func NewEvent[T any](name Name) *Event[T] {
	return &Event[T]{
		listeners: make(map[string]Listener[T]),
		name:      name,
	}
}

func (e *Event[T]) Name() Name {
	return e.name
}

func (e *Event[T]) SetName(n Name) {
	e.name = n
}

func (e *Event[T]) Dispatch(data T) {
	for _, l := range e.listeners {
		l(data)
	}
}

func (e *Event[T]) Add(listenerName string, listener Listener[T]) error {
	if _, ok := e.listeners[listenerName]; ok {
		return fmt.Errorf("listener with name \"%s\" already exists", listenerName)
	}

	e.listeners[listenerName] = listener
	return nil
}

func (e *Event[T]) Remove(listenerName string) {
	for cn := range e.listeners {
		if cn == listenerName {
			delete(e.listeners, listenerName)
		}
	}
}
