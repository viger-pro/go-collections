package collections

import (
	"context"
	"errors"
)

type ChannelledQueueWithLimit[T any] struct {
	c chan T
}

func NewChannelledQueueWithLimit[T any](maxSize uint) *ChannelledQueueWithLimit[T] {
	return &ChannelledQueueWithLimit[T]{
		c: make(chan T, maxSize),
	}
}

func (q *ChannelledQueueWithLimit[T]) AddLast(ctx context.Context, t T) error {
	select {
	case q.c <- t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *ChannelledQueueWithLimit[T]) TryAddLast(t T) error {
	select {
	case q.c <- t:
	default:
		return errors.New("queue is full")
	}
	return nil
}

func (q *ChannelledQueueWithLimit[T]) RemoveFirst(ctx context.Context) (t T, err error) {
	select {
	case t = <-q.c:
		return t, nil
	case <-ctx.Done():
		return t, ctx.Err()
	}
}

func (q *ChannelledQueueWithLimit[T]) TryRemoveFirst() (t T, err error) {
	select {
	case t = <-q.c:
		return t, nil
	default:
		return t, errors.New("queue is empty")
	}
}

func (q *ChannelledQueueWithLimit[T]) MaxSize() uint {
	return uint(cap(q.c))
}

func (q *ChannelledQueueWithLimit[T]) Size() uint {
	return uint(len(q.c))
}
