package collections

import (
	"errors"
)

type entry[T any] struct {
	value T
	next  *entry[T]
}

// LinkedQueue a queue based on single-linked list. This implementation is not threadsafe.
type LinkedQueue[T any] struct {
	head *entry[T]
	tail *entry[T]
	size uint
}

func NewLinkedQueue[T any]() *LinkedQueue[T] {
	return &LinkedQueue[T]{}
}

func (q *LinkedQueue[T]) AddLast(value T) {
	e := &entry[T]{value: value}
	if q.tail != nil {
		q.tail.next = e
	} else {
		q.head = e
	}
	q.tail = e
	q.size++
}

func (q *LinkedQueue[T]) RemoveFirst() (T, error) {
	if q.head == nil {
		var t T
		return t, errors.New("the queue is empty")
	}
	e := q.head
	q.head = e.next
	if q.head == nil {
		q.tail = nil
	}
	q.size--
	return e.value, nil
}

func (q *LinkedQueue[T]) Size() uint {
	return q.size
}
