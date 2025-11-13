package collections

import (
	"context"
	"sync"
	"testing"
	"time"
)

type testCase struct {
	name  string
	queue QueueWithLimit[uint]
}

func createTests(queueSize uint) ([]testCase, error) {
	var result = make([]testCase, 0)
	var queue QueueWithLimit[uint]
	var err error

	if queue, err = NewLinkedQueueWithLimit[uint](queueSize); err != nil {
		return nil, err
	}
	result = append(result, testCase{
		name:  "standard queue on linked queue",
		queue: queue,
	})

	if queue, err = NewArrayQueueWithLimit[uint](queueSize); err != nil {
		return nil, err
	}
	result = append(result, testCase{
		name:  "standard queue on array queue",
		queue: queue,
	})

	if queue, err = NewSimpleArrayQueueWithLimit[uint](queueSize); err != nil {
		return nil, err
	}
	result = append(result, testCase{
		name:  "standard queue on simple array queue",
		queue: queue,
	})

	queue = NewChannelledQueueWithLimit[uint](queueSize)
	result = append(result, testCase{
		name:  "standard queue on channelled queue",
		queue: queue,
	})

	return result, nil
}

func TestQueueWithLimit_HappyPath(t *testing.T) {
	var queueSize uint = 10
	tests, err := createTests(queueSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ctx = context.Background()
			var i uint
			for i = 0; i < queueSize; i++ {
				err := test.queue.AddLast(ctx, i)
				if err != nil {
					t.Error(err)
				}
			}
			for i = 0; i < queueSize; i++ {
				x, err := test.queue.RemoveFirst(ctx)
				if err != nil {
					t.Error(err)
				}
				if x != i {
					t.Fatalf("expected %d, got %d", i, x)
				}
			}
		})
	}
}

func TestQueueWithLimit_Blocking(t *testing.T) {
	var steps uint = 100
	var queueSize uint = 10
	tests, err := createTests(queueSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ctx = context.Background()
			wg := &sync.WaitGroup{}
			var i uint
			var enqueueResults = make(chan error, steps)
			for i = 0; i < steps; i++ {
				x := i
				wg.Go(func() {
					enqueueResults <- test.queue.AddLast(ctx, x)
				})
			}
			var dequeueResults = make(map[uint]bool)
			for i = 0; i < steps; i++ {
				x, err := test.queue.RemoveFirst(ctx)
				if err != nil {
					t.Fatalf("error while dequeueing %dth time: %v", i, err)
				}
				dequeueResults[x] = true
			}
			wg.Wait()
			for i = 0; i < steps; i++ {
				err := <-enqueueResults
				if err != nil {
					t.Fatal(err)
				}
			}
			for i = 0; i < steps; i++ {
				_, ok := dequeueResults[i]
				if !ok {
					t.Fatalf("%d has not been inserted into the queue", i)
				}
			}
		})
	}
}

func TestQueueWithLimit_TryRemoveFirst(t *testing.T) {
	var queueSize uint = 10
	tests, err := createTests(queueSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var i uint
			for i = 0; i < queueSize; i++ {
				err := test.queue.AddLast(ctx, i)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i = 0; i < queueSize; i++ {
				_, err := test.queue.RemoveFirst(ctx)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i = 0; i < queueSize; i++ {
				_, err := test.queue.TryRemoveFirst()
				if err == nil || err.Error() != "queue is empty" {
					t.Fatalf("expected queue is empty error, got %v", err)
				}
			}
		})
	}
}

func TestQueueWithLimit_AddLast_WithTimeout(t *testing.T) {
	var queueSize uint = 2
	tests, err := createTests(queueSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var i uint
			for i = 0; i < queueSize; i++ {
				err := test.queue.AddLast(ctx, i)
				if err != nil {
					t.Fatal(err)
				}
			}
			var noTimeout = false
			wg := sync.WaitGroup{}
			for i = 0; i < queueSize; i++ {
				wg.Go(func() {
					timeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
					defer cancel()
					err := test.queue.AddLast(timeout, i)
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
			wg.Add(int(queueSize))
			results := make(chan error, queueSize)
			timeout := 5 * time.Second
			for i = 0; i < queueSize; i++ {
				go func(x uint) {
					wg.Done()
					timeout, cancel := context.WithTimeout(ctx, timeout)
					defer cancel()
					err := test.queue.AddLast(timeout, x)
					results <- err
				}(i)
			}
			wg.Wait()
			for i = 0; i < queueSize; i++ {
				_, err := test.queue.RemoveFirst(ctx)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i = 0; i < queueSize; i++ {
				err := <-results
				if err != nil {
					t.Fatalf("error when adding %dth element: %v", i, err)
				}
			}
		})
	}
}

func TestQueueWithLimit_Randomized(t *testing.T) {
	var queueSize uint = 100
	var steps uint = 10_000
	tests, err := createTests(queueSize)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctx := context.Background()
			var i uint

			type Result struct {
				value uint
				err   error
			}

			added := make([]bool, steps)
			removed := make([]bool, steps)

			addedChan := make(chan Result, steps)
			removedChan := make(chan Result, steps)
			waitGroup := sync.WaitGroup{}
			waitGroup.Add(int(steps) * 2)
			for i = 0; i < steps; i++ {
				go func() {
					value, err := test.queue.RemoveFirst(ctx)
					removedChan <- Result{value, err}
					waitGroup.Done()
				}()
				go func(x uint) {
					err := test.queue.AddLast(ctx, x)
					addedChan <- Result{x, err}
					waitGroup.Done()
				}(i)
			}
			waitGroup.Wait()
			close(addedChan)
			close(removedChan)

			if test.queue.Size() != 0 {
				t.Fatalf("expected queue to be empty")
			}

			for addedResult := range addedChan {
				if addedResult.err != nil {
					t.Fatal(addedResult.err)
				}
				added[addedResult.value] = true
			}
			for removedResult := range removedChan {
				if removedResult.err != nil {
					t.Fatal(removedResult.err)
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
		})
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

func BenchmarkQueues(b *testing.B) {
	tests, err := createTests(1000)
	if err != nil {
		b.Fatal(err)
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			ctx := context.Background()
			for i := 0; i < b.N; i++ {
				err := test.queue.AddLast(ctx, uint(i))
				if err != nil {
					b.Fatal(err)
				}
				x, err := test.queue.RemoveFirst(ctx)
				if err != nil {
					b.Fatal(err)
				}
				if x != uint(i) {
					b.Fatalf("expected %d, got %d", i, x)
				}
			}
		})
	}
}
