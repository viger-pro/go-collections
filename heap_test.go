package collections

import (
	"cmp"
	"math/rand"
	"sort"
	"testing"
)

type heapTestCase[T any] struct {
	name  string
	heap  *Heap[int]
	input []T
}

func createStandardComparisonHeapTestCases() []heapTestCase[int] {
	return createHeapTestCases(func(initialCapacity int) *Heap[int] {
		return NewHeap[int](initialCapacity)
	})
}

func createCustomComparisonHeapTestCases() []heapTestCase[int] {
	return createHeapTestCases(func(initialCapacity int) *Heap[int] {
		return NewHeapWithCompare[int](initialCapacity, func(x, y int) int {
			return -cmp.Compare(x, y)
		})
	})
}

func createHeapTestCases(createHeapFunc func(initialCapacity int) *Heap[int]) []heapTestCase[int] {
	return []heapTestCase[int]{
		{
			name:  "0 elements ordered",
			heap:  createHeapFunc(0),
			input: []int{},
		},
		{
			name:  "1 elements ordered",
			heap:  createHeapFunc(1),
			input: []int{0},
		},
		{
			name:  "10 elements ordered",
			heap:  createHeapFunc(10),
			input: orderedIntArray(10),
		},
		{
			name:  "100 elements ordered",
			heap:  createHeapFunc(10),
			input: orderedIntArray(100),
		},
		{
			name:  "1000 elements ordered",
			heap:  createHeapFunc(10),
			input: orderedIntArray(1000),
		},
		{
			name:  "0 elements random",
			heap:  createHeapFunc(0),
			input: []int{},
		},
		{
			name:  "1 elements random",
			heap:  createHeapFunc(1),
			input: []int{0},
		},
		{
			name:  "10 elements random",
			heap:  createHeapFunc(10),
			input: randomIntArray(10),
		},
		{
			name:  "100 elements random",
			heap:  createHeapFunc(10),
			input: randomIntArray(100),
		},
		{
			name:  "1000 elements random",
			heap:  createHeapFunc(10),
			input: randomIntArray(1000),
		},
	}
}

func TestHeapWithStandardComparison_HappyPath(t *testing.T) {
	heapTestCases := createStandardComparisonHeapTestCases()
	for _, heapTestCase := range heapTestCases {
		t.Run(heapTestCase.name, func(t *testing.T) {
			testHappyPath(heapTestCase.heap, heapTestCase.input, sorted(heapTestCase.input), t)
		})
	}
}

func TestHeapWithCustomComparison_HappyPath(t *testing.T) {
	heapTestCases := createCustomComparisonHeapTestCases()
	for _, heapTestCase := range heapTestCases {
		t.Run(heapTestCase.name, func(t *testing.T) {
			testHappyPath(heapTestCase.heap, heapTestCase.input, reverse(sorted(heapTestCase.input)), t)
		})
	}
}

func testHappyPath(heap *Heap[int], input []int, expected []int, t *testing.T) {
	if !heap.IsEmpty() {
		t.Fatalf("heap is not empty")
	}
	if heap.Size() != 0 {
		t.Fatalf("heap size is not zero")
	}

	var n = len(input)
	for i := 0; i < n; i++ {
		heap.Add(input[i])
		if heap.Size() != i+1 {
			t.Fatalf("expected heap size to be %d, got %d, input: %v", i+1, heap.Size(), input)
		}
	}
	for i := 0; i < n; i++ {
		x, err := heap.Remove()
		if err != nil {
			t.Fatalf("error while removing from heap: %v, input: %v", err, input)
		}
		if x != expected[i] {
			t.Fatalf("expected removed element to be %d, got %d, input: %v", expected[i], x, input)
		}
		if heap.Size() != n-i-1 {
			t.Fatalf("expected heap size to be %d, got %d, input: %v", n-i, heap.Size(), input)
		}
	}
	if !heap.IsEmpty() {
		t.Fatalf("heap is not empty, input: %v", input)
	}
	if heap.Size() != 0 {
		t.Fatalf("heap size is not zero, input: %v", input)
	}
	x, err := heap.Remove()
	if err == nil {
		t.Fatalf("expected error, got %d, input: %v", x, input)
	}
}

func TestHeapWithStandardComparison_Duplicated(t *testing.T) {
	var n = 100
	var array1 = randomIntArray(n)
	var array2 = randomIntArray(n)

	var heap = NewHeap[int](1)
	for i := 0; i < n; i++ {
		heap.Add(array1[i])
		heap.Add(array2[i])
	}

	for i := 0; i < n; i++ {
		x, err1 := heap.Remove()
		y, err2 := heap.Remove()
		if err1 != nil || err2 != nil {
			t.Fatalf("error while removing from heap: %v, %v, \ninput1: %v, \ninput2: %v",
				err1, err2, array1, array2)
		}
		if x != i {
			t.Fatalf("expected removed element to be %d, got %d, \ninput1: %v, \ninput2: %v",
				i, x, array1, array2)
		}
		if y != i {
			t.Fatalf("expected removed element to be %d, got %d, \ninput1: %v, \ninput2: %v",
				i, y, array1, array2)
		}
	}
}

func orderedIntArray(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = i
	}
	return result
}

func randomIntArray(n int) []int {
	result := orderedIntArray(n)
	for i := n - 1; i > 0; i-- {
		swap(result, i, rand.Intn(i))
	}
	return result
}

func reverse(ints []int) []int {
	left := 0
	right := len(ints) - 1
	for left < right {
		swap(ints, left, right)
		left++
		right--
	}
	return ints
}

func sorted(input []int) []int {
	result := make([]int, len(input))
	copy(result, input)
	sort.Ints(result)
	return result
}
