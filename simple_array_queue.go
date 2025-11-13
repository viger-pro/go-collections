package collections

import "errors"

type SimpleArrayQueue[T any] struct {
	queue []T
}

func NewSimpleArrayQueue[T any]() *SimpleArrayQueue[T] {
	return &SimpleArrayQueue[T]{
		queue: make([]T, 0),
	}
}

func NewSimpleArrayQueueWithInitialCapacity[T any](capacity uint) *SimpleArrayQueue[T] {
	return &SimpleArrayQueue[T]{
		queue: make([]T, 0, capacity),
	}
}

func (q *SimpleArrayQueue[T]) AddLast(t T) {
	q.queue = append(q.queue, t)
}

func (q *SimpleArrayQueue[T]) RemoveFirst() (T, error) {
	var zero T
	if len(q.queue) == 0 {
		return zero, errors.New("queue is empty")
	}
	var result T = q.queue[0]
	q.queue = q.queue[1:]
	return result, nil
}

func (q *SimpleArrayQueue[T]) Size() uint {
	return uint(len(q.queue))
}
