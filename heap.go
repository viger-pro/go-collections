package collections

import (
	"cmp"
	"errors"
)

var ErrEmptyHeap = errors.New("heap is empty")

type Heap[T any] struct {
	array   []T
	size    int
	compare func(T, T) int
}

func NewHeap[T cmp.Ordered](initialCapacity int) *Heap[T] {
	return NewHeapWithCompare(initialCapacity, func(t1, t2 T) int {
		return cmp.Compare(t1, t2)
	})
}

func NewHeapWithCompare[T any](initialCapacity int, compare func(t1, t2 T) int) *Heap[T] {
	return &Heap[T]{
		array:   make([]T, 0, initialCapacity),
		size:    0,
		compare: compare,
	}
}

func (heap *Heap[T]) IsEmpty() bool {
	return heap.size == 0
}

func (heap *Heap[T]) Size() int {
	return heap.size
}

func (heap *Heap[T]) Add(element T) {
	if heap.size < len(heap.array) {
		heap.array[heap.size] = element
	} else {
		heap.array = append(heap.array, element)
	}
	heap.siftUp(heap.array, heap.size)
	heap.size += 1
}

func (heap *Heap[T]) GetFirst() (t T, err error) {
	if heap.size == 0 {
		return t, ErrEmptyHeap
	}
	return heap.array[0], nil
}

func (heap *Heap[T]) Remove() (t T, err error) {
	if heap.size == 0 {
		return t, ErrEmptyHeap
	}
	element := heap.array[0]
	heap.size -= 1
	swap(heap.array, 0, heap.size)
	heap.siftDown(heap.array, 0, heap.size-1)
	return element, nil
}

func (heap *Heap[T]) siftUp(array []T, index int) {
	for index > 0 {
		parentIndex := (index - 1) / 2
		if heap.compare(array[parentIndex], array[index]) > 0 {
			swap(array, index, parentIndex)
		} else {
			break
		}
		index = parentIndex
	}
}

func (heap *Heap[T]) siftDown(array []T, i int, lastIndex int) {
	for {
		j := i*2 + 1
		if j > lastIndex {
			break
		}
		k := j + 1
		var minIndex int
		if k > lastIndex || heap.compare(array[j], array[k]) < 0 {
			minIndex = j
		} else {
			minIndex = k
		}
		if heap.compare(array[i], array[minIndex]) < 0 {
			break
		} else {
			swap(array, i, minIndex)
			i = minIndex
		}
	}
}

func swap[K any](array []K, index int, index2 int) {
	array[index], array[index2] = array[index2], array[index]
}
