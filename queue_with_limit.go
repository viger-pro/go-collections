package collections

import "time"

type QueueWithLimit[T any] interface {
	AddLast(T) error
	TryAddLast(T) error
	TryAddLastWithTimeout(T, time.Duration) error

	RemoveFirst() (T, error)
	TryRemoveFirst() (T, error)

	MaxSize() uint
	Size() uint
}
