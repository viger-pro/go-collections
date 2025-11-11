package collections

import (
	"testing"
)

func TestLinkedQueue(t *testing.T) {
	var queue = NewLinkedQueue[int]()
	var n uint = 10
	fillQueue(queue, n, t)
	emptyQueue(queue, n, t)
}

func fillQueue(queue *LinkedQueue[int], n uint, t *testing.T) {
	for i := 0; i < int(n); i++ {
		queue.AddLast(i)
		if queue.Size() != uint(i+1) {
			t.Fatalf("expected %d got %d", i+1, queue.Size())
		}
	}
}

func emptyQueue(queue *LinkedQueue[int], n uint, t *testing.T) {
	for i := 0; i < int(n); i++ {
		if queue.Size() != n-uint(i) {
			t.Fatalf("expected %d got %d", n, queue.Size())
		}
		x, err := queue.RemoveFirst()
		if err != nil {
			t.Fatal(err)
		}
		if x != i {
			t.Fatalf("expected %d got %d", i, x)
		}
	}
}
