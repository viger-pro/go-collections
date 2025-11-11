package collections

type Queue[T any] interface {
	AddLast(T)
	RemoveFirst() (T, error)
	Size() uint
}
