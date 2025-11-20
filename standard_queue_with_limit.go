package collections

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/semaphore"
)

// StandardQueueWithLimit an implementation of QueueWithLimit built on a Queue.
type StandardQueueWithLimit[T any] struct {
	freeSlotsSemaphore *semaphore.Weighted
	fullSlotsSemaphore *semaphore.Weighted
	lock               *sync.Mutex
	maxSize            uint
	queue              Queue[T]
}

func NewLinkedQueueWithLimit[T any](maxSize uint) (*StandardQueueWithLimit[T], error) {
	return newQueueWithLimit(maxSize, NewLinkedQueue[T]())
}

func NewArrayQueueWithLimit[T any](maxSize uint) (*StandardQueueWithLimit[T], error) {
	return newQueueWithLimit(maxSize, NewArrayQueueWithInitialCapacity[T](maxSize))
}

func NewSimpleArrayQueueWithLimit[T any](maxSize uint) (*StandardQueueWithLimit[T], error) {
	return newQueueWithLimit(maxSize, NewSimpleArrayQueueWithInitialCapacity[T](maxSize))
}

func newQueueWithLimit[T any](maxSize uint, queue Queue[T]) (*StandardQueueWithLimit[T], error) {
	lock := new(sync.Mutex)
	freeSlotsSemaphore := semaphore.NewWeighted(int64(maxSize))
	if err := freeSlotsSemaphore.Acquire(context.Background(), int64(maxSize)); err != nil {
		return nil, err
	}
	return &StandardQueueWithLimit[T]{
		freeSlotsSemaphore: freeSlotsSemaphore,
		fullSlotsSemaphore: semaphore.NewWeighted(int64(maxSize)),
		lock:               lock,
		maxSize:            maxSize,
		queue:              queue,
	}, nil
}

func (q *StandardQueueWithLimit[T]) AddLast(ctx context.Context, value T) (err error) {
	if err = q.fullSlotsSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue.AddLast(value)
	q.freeSlotsSemaphore.Release(1)
	return nil
}

func (q *StandardQueueWithLimit[T]) MaxSize() uint {
	return q.maxSize
}

func (q *StandardQueueWithLimit[T]) RemoveFirst(ctx context.Context) (t T, err error) {
	if err = q.freeSlotsSemaphore.Acquire(ctx, 1); err != nil {
		return t, err
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	t, err = q.queue.RemoveFirst()
	if err != nil {
		return t, err
	}
	q.fullSlotsSemaphore.Release(1)
	return t, err
}

func (q *StandardQueueWithLimit[T]) Size() uint {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue.Size()
}

func (q *StandardQueueWithLimit[T]) TryAddLast(value T) (err error) {
	if !q.freeSlotsSemaphore.TryAcquire(1) {
		return errors.New("queue is full")
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue.AddLast(value)
	q.fullSlotsSemaphore.Release(1)
	return err
}

func (q *StandardQueueWithLimit[T]) TryRemoveFirst() (t T, err error) {
	if !q.freeSlotsSemaphore.TryAcquire(1) {
		return t, errors.New("queue is empty")
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	first, err := q.queue.RemoveFirst()
	if err != nil {
		return t, err
	}
	q.fullSlotsSemaphore.Release(1)
	return first, err
}
