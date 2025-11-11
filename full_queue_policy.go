package collections

import (
	"errors"
	"sync"
)

type FullQueuePolicy[T any] interface {
	// EnsureCanAdd checks whether the queue is full. If it is not full, this function returns no error. If the queue
	// is full, the behaviour will depend on the implementation: one implementation might simply reject new entries
	// and return an error while the other might wait till the queue is not full again.
	EnsureCanAdd() error

	// ElementRemoved function called when the element is removed from the queue.
	ElementRemoved()
}

// RejectingPolicy will simply reject any new entry if the queue is full.
type RejectingPolicy[T any] struct {
	maxSize uint
	size    uint
}

func NewRejectingPolicy[T any](maxSize uint) *RejectingPolicy[T] {
	return &RejectingPolicy[T]{maxSize: maxSize}
}

func (p *RejectingPolicy[T]) EnsureCanAdd() error {
	if p.size == p.maxSize {
		return errors.New("queue is full")
	}
	p.size++
	return nil
}

func (p *RejectingPolicy[T]) ElementRemoved() {
	p.size--
}

// BlockingPolicy will wait undefinitely until the queue regains any capacity.
type BlockingPolicy[T any] struct {
	cond *sync.Cond

	maxSize uint
	size    uint
}

func NewBlockingPolicy[T any](maxSize uint) *BlockingPolicy[T] {
	return &BlockingPolicy[T]{
		cond:    sync.NewCond(&sync.Mutex{}),
		maxSize: maxSize,
		size:    0,
	}
}

func (p *BlockingPolicy[T]) EnsureCanAdd() error {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	if p.size == p.maxSize {
		p.cond.Wait()
	}
	p.size++
	return nil
}

func (p *BlockingPolicy[T]) ElementRemoved() {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	p.size--
	p.cond.Signal()
}
