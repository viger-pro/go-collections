package collections

import (
	"errors"
	"sync"
	"time"
)

// LinkedQueueWithBlockingLimit an implementation of QueueWithLimit that will block when you try to enqueue an element
// if the queue is full. Some other thread must dequeue from the queue to unblock such enqueueing.
// It will also block when you try to dequeue from an empty queue. Some other thread must enqueue to such a queue to
// unblock dequeue operation.
// This implementation is threadsafe.
type LinkedQueueWithBlockingLimit[T any] struct {
	dequeueChan chan bool
	enqueueChan chan bool
	lock        *sync.Mutex
	maxSize     uint
	queue       *LinkedQueue[T]
}

func NewLinkedQueueWithBlockingLimit[T any](maxSize uint) *LinkedQueueWithBlockingLimit[T] {
	lock := new(sync.Mutex)
	return &LinkedQueueWithBlockingLimit[T]{
		dequeueChan: make(chan bool),
		enqueueChan: make(chan bool),
		lock:        lock,
		maxSize:     maxSize,
		queue:       NewLinkedQueue[T](),
	}
}

func (q *LinkedQueueWithBlockingLimit[T]) Enqueue(value T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	select {
	case <-q.enqueueChan:
		break
	default:
		break
	}

	if q.queue.Size() >= q.maxSize {
		q.lock.Unlock()
		<-q.enqueueChan
		q.lock.Lock()
	}
	q.queue.Enqueue(value)

	select {
	case q.dequeueChan <- true:
		break
	default:
		break
	}
	return nil
}

func (q *LinkedQueueWithBlockingLimit[T]) Dequeue() (t T, err error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	select {
	case <-q.dequeueChan:
		break
	default:
		break
	}

	if q.queue.Size() == 0 {
		q.lock.Unlock()
		<-q.dequeueChan
		q.lock.Lock()
	}
	t, err = q.queue.Dequeue()

	select {
	case q.enqueueChan <- true:
		break
	default:
		break
	}
	return t, err
}

func (q *LinkedQueueWithBlockingLimit[T]) MaxSize() uint {
	return q.maxSize
}

func (q *LinkedQueueWithBlockingLimit[T]) Size() uint {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.queue.Size()
}

// TryEnqueue attempts to insert an element into the queue and in case the queue is full,
// it returns an error immediately, without waiting.
func (q *LinkedQueueWithBlockingLimit[T]) TryEnqueue(value T) (err error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.queue.Size() >= q.maxSize {
		err = errors.New("queue is full")
	} else {
		q.queue.Enqueue(value)
		select {
		case q.dequeueChan <- true:
			break
		default:
			break
		}
	}
	return err
}

// TryEnqueueWithTimeout attempts to insert an element into the queue and in case the queue is full it waits up to
// the given timeout for the queue to be dequeued.
func (q *LinkedQueueWithBlockingLimit[T]) TryEnqueueWithTimeout(value T, timeout time.Duration) (err error) {
	return nil
}
