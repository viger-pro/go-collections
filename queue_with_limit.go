package collections

import (
	"context"
)

// QueueWithLimit threadsafe queue that might block while trying to add an element to a full queue or removing
// an element from an empty queue.
type QueueWithLimit[T any] interface {

	// AddLast adds element to the end of the queue. Blocks if the queue is full until there is a space for
	// a new element or until the given context is done. Returns an error if the element has not been added.
	AddLast(context.Context, T) error

	// TryAddLast adds element to the end of the queue. If the queue is full it returns immediately returning an error.
	TryAddLast(T) error

	// RemoveFirst removes first element from the queue. If the queue is empty if blocks until there is some element
	// added or until the given context is done. Return either (element, nil) or (zero, error) in case no element has
	// not removed.
	RemoveFirst(context.Context) (T, error)

	// TryRemoveFirst removes first element from the queue. If the queue is empty it returns an error immediately.
	TryRemoveFirst() (T, error)

	// MaxSize returns max number of elements this queue can store.
	MaxSize() uint

	// Size returns current number of elements this queue stores.
	Size() uint
}
