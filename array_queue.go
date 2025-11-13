package collections

import "errors"

type ArrayQueue[T any] struct {
	array []T
	head  int
	tail  int
	size  uint
}

func NewArrayQueue[T any]() *ArrayQueue[T] {
	return &ArrayQueue[T]{
		array: make([]T, 0),
	}
}

func NewArrayQueueWithInitialCapacity[T any](capacity uint) *ArrayQueue[T] {
	return &ArrayQueue[T]{
		array: make([]T, 0, capacity),
	}
}

func (q *ArrayQueue[T]) AddLast(t T) {
	if q.size == uint(len(q.array)) {
		q.increaseCapacity()
	}
	q.array[q.tail] = t
	q.tail = (q.tail + 1) % len(q.array)
	q.size++
}

func (q *ArrayQueue[T]) RemoveFirst() (T, error) {
	var zero T
	if q.size == 0 {
		return zero, errors.New("queue is empty")
	}
	x := q.array[q.head]
	q.array[q.head] = zero
	q.head = (q.head + 1) % len(q.array)
	q.size--
	return x, nil
}

func (q *ArrayQueue[T]) Size() uint {
	return q.size
}

func (q *ArrayQueue[T]) increaseCapacity() {
	length := len(q.array)
	var zero T
	if length == 0 {
		q.array = append(q.array, zero)
	} else {
		if q.head > 0 {
			for i := q.head; i < length; i++ {
				q.array = append(q.array, q.array[i])
			}
			for i := 0; i < q.head; i++ {
				q.array = append(q.array, q.array[i])
			}
			for i := 0; i < length; i++ {
				q.array[i] = q.array[length+i]
			}
		} else {
			for i := 0; i < length; i++ {
				q.array = append(q.array, zero)
			}
		}
	}
	q.head = 0
	q.tail = length
}
