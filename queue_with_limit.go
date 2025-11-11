package collections

type QueueWithLimit[T any] interface {
	Enqueue(T) error
	Dequeue() (T, error)
	MaxSize() uint
	Size() uint
}
