package collections

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/semaphore"
)

// LinkedQueueWithBlockingLimit an implementation of QueueWithLimit built on single-linked list.
type LinkedQueueWithBlockingLimit[T any] struct {
	freeSlotsSemaphore *semaphore.Weighted
	fullSlotsSemaphore *semaphore.Weighted
	lock               *sync.Mutex
	maxSize            uint
	queue              *LinkedQueue[T]
}

func NewLinkedQueueWithBlockingLimit[T any](maxSize uint) (*LinkedQueueWithBlockingLimit[T], error) {
	lock := new(sync.Mutex)
	freeSlotsSemaphore := semaphore.NewWeighted(int64(maxSize))
	if err := freeSlotsSemaphore.Acquire(context.Background(), int64(maxSize)); err != nil {
		return nil, err
	}
	return &LinkedQueueWithBlockingLimit[T]{
		freeSlotsSemaphore: freeSlotsSemaphore,
		fullSlotsSemaphore: semaphore.NewWeighted(int64(maxSize)),
		lock:               lock,
		maxSize:            maxSize,
		queue:              NewLinkedQueue[T](),
	}, nil
}

func (q *LinkedQueueWithBlockingLimit[T]) AddLast(ctx context.Context, value T) (err error) {
	if err = q.fullSlotsSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue.AddLast(value)
	q.freeSlotsSemaphore.Release(1)
	return nil
}

func (q *LinkedQueueWithBlockingLimit[T]) MaxSize() uint {
	return q.maxSize
}

func (q *LinkedQueueWithBlockingLimit[T]) RemoveFirst(ctx context.Context) (t T, err error) {
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

func (q *LinkedQueueWithBlockingLimit[T]) Size() uint {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue.Size()
}

func (q *LinkedQueueWithBlockingLimit[T]) TryAddLast(value T) (err error) {
	if !q.freeSlotsSemaphore.TryAcquire(1) {
		return errors.New("queue is full")
	}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.queue.AddLast(value)
	q.fullSlotsSemaphore.Release(1)
	return err
}

func (q *LinkedQueueWithBlockingLimit[T]) TryRemoveFirst() (t T, err error) {
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
