package collections

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestLinkedQueueWithBlockingLimit_HappyPath(t *testing.T) {
	var n uint = 10
	queue, err := NewLinkedQueueWithBlockingLimit[int](n)
	if err != nil {
		t.Fatal(err)
	}
	var ctx = context.Background()
	for i := 0; i < int(n); i++ {
		err := queue.AddLast(ctx, i)
		if err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < int(n); i++ {
		x, err := queue.RemoveFirst(ctx)
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

	queue, err := NewLinkedQueueWithBlockingLimit[int](m)
	if err != nil {
		t.Fatal(err)
	}

	var ctx = context.Background()
	wg := &sync.WaitGroup{}
	for i := 0; i < int(n); i++ {
		wg.Go(func() {
			enqueueResults <- queue.AddLast(ctx, i)
		})
	}
	var dequeueResults = make(map[int]bool)
	for i := 0; i < int(n); i++ {
		x, err := queue.RemoveFirst(ctx)
		if err != nil {
			t.Fatalf("error while dequeueing %dth time: %v", i, err)
		}
		dequeueResults[x] = true
	}
	wg.Wait()
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
	queue, err := NewLinkedQueueWithBlockingLimit[int](n)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for i := 0; i < int(n); i++ {
		err := queue.AddLast(ctx, i)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < int(n); i++ {
		_, err := queue.RemoveFirst(ctx)
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

func TestLinkedQueueWithBlockingLimit_AddLast_WithTimeout(t *testing.T) {
	var n uint = 2
	queue, err := NewLinkedQueueWithBlockingLimit[int](n)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	for i := 0; i < int(n); i++ {
		err := queue.AddLast(ctx, i)
		if err != nil {
			t.Fatal(err)
		}
	}
	var noTimeout = false
	wg := sync.WaitGroup{}
	for i := 0; i < int(n); i++ {
		wg.Go(func() {
			timeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
			defer cancel()
			err := queue.AddLast(timeout, i)
			if err == nil {
				noTimeout = true
			}
		})
	}
	wg.Wait()
	if noTimeout {
		t.Fatalf("expected timeout while trying to add to a full queue")
	}

	wg = sync.WaitGroup{}
	wg.Add(int(n))
	results := make(chan error, n)
	timeout := 5 * time.Second
	for i := 0; i < int(n); i++ {
		go func(x int) {
			wg.Done()
			timeout, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			err = queue.AddLast(timeout, x)
			results <- err
		}(i)
	}
	wg.Wait()
	for i := 0; i < int(n); i++ {
		_, err := queue.RemoveFirst(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < int(n); i++ {
		err := <-results
		if err != nil {
			t.Fatalf("error when adding %dth element: %v", i, err)
		}
	}
}

func TestLinkedQueueWithBlockingLimit_Randomized(t *testing.T) {
	var limit uint = 100
	var n uint = 10_000
	queue, err := NewLinkedQueueWithBlockingLimit[uint](limit)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	var i uint

	type Result struct {
		value uint
		err   error
	}

	added := make([]bool, n)
	removed := make([]bool, n)

	addedChan := make(chan Result, n)
	removedChan := make(chan Result, n)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(int(n) * 2)
	for i = 0; i < n; i++ {
		go func() {
			value, err := queue.RemoveFirst(ctx)
			removedChan <- Result{value, err}
			waitGroup.Done()
		}()
		go func(x uint) {
			err := queue.AddLast(ctx, x)
			addedChan <- Result{x, err}
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()
	close(addedChan)
	close(removedChan)

	if queue.Size() != 0 {
		t.Fatalf("expected queue to be empty")
	}

	for addedResult := range addedChan {
		if addedResult.err != nil {
			t.Fatal(err)
		}
		added[addedResult.value] = true
	}
	for removedResult := range removedChan {
		if removedResult.err != nil {
			t.Fatal(err)
		}
		removed[removedResult.value] = true
	}

	notAdded := not(added)
	if len(notAdded) != 0 {
		t.Errorf("elements that should have been be added but have not: %v", notAdded)
		t.Fatalf("elements that should have been be added but have not: %d", len(notAdded))
	}
	notRemoved := not(removed)
	if len(notRemoved) != 0 {
		t.Errorf("elements that should have been be removed but have not: %v", notRemoved)
		t.Fatalf("elements that should have been be removed but have not: %d", len(notRemoved))
	}
}

func not(a []bool) []uint {
	result := make([]uint, 0, len(a))
	for i := 0; i < len(a); i++ {
		if !a[i] {
			result = append(result, uint(i))
		}
	}
	return result
}
