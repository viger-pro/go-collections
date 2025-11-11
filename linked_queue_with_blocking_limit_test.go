package collections

import "testing"

func TestLinkedQueueWithBlockingLimit_HappyPath(t *testing.T) {
	var n uint = 10
	var queue = NewLinkedQueueWithBlockingLimit[int](n)
	for i := 0; i < int(n); i++ {
		err := queue.AddLast(i)
		if err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < int(n); i++ {
		x, err := queue.RemoveFirst()
		if err != nil {
			t.Error(err)
		}
		if x != i {
			t.Fatalf("expected %d, got %d", i, x)
		}
	}
}

func TestLinkedQueueWithBlockingLimit_Blocking(t *testing.T) {
	var n uint = 100000
	var m uint = 100
	var enqueueResults = make(chan error, n)
	var queue = NewLinkedQueueWithBlockingLimit[int](m)
	for i := 0; i < int(n); i++ {
		go func() {
			enqueueResults <- queue.AddLast(i)
		}()
	}
	var dequeueResults = make(map[int]bool)
	for i := 0; i < int(n); i++ {
		x, err := queue.RemoveFirst()
		if err != nil {
			t.Fatalf("error while dequeueing %dth time: %v", i, err)
		}
		dequeueResults[x] = true
	}
	for i := 0; i < int(n); i++ {
		err := <-enqueueResults
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < int(n); i++ {
		_, ok := dequeueResults[i]
		if !ok {
			t.Fatalf("%d has not been inserted into the queue", i)
		}
	}
}

func TestLinkedQueueWithBlockingLimit_TryRemoveFirst(t *testing.T) {
	var n uint = 10
	var queue = NewLinkedQueueWithBlockingLimit[int](n)
	for i := 0; i < int(n); i++ {
		err := queue.AddLast(i)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < int(n); i++ {
		_, err := queue.RemoveFirst()
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < int(n); i++ {
		_, err := queue.TryRemoveFirst()
		if err == nil || err.Error() != "queue is empty" {
			t.Fatalf("expected queue is empty error, got %v", err)
		}
	}
}
