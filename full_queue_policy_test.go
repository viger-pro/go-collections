package collections

import (
	"testing"
	"time"
)

func TestRejectingPolicy(t *testing.T) {
	var err error
	var n uint = 10
	p := NewRejectingPolicy[int](n)
	for i := 0; i < int(n); i++ {
		err = p.EnsureCanAdd()
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	err = p.EnsureCanAdd()
	if err == nil {
		t.Fatal("expected an error")
	}
	p.ElementRemoved()
	err = p.EnsureCanAdd()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestBlockingPolicy(t *testing.T) {
	var err error
	var n uint = 10
	p := NewBlockingPolicy[int](n)
	for i := 0; i < int(n); i++ {
		err = p.EnsureCanAdd()
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	resultChannel := make(chan error)
	go func() {
		err = p.EnsureCanAdd()
		resultChannel <- err
	}()

	time.Sleep(50 * time.Millisecond)
	p.ElementRemoved()
	err = p.EnsureCanAdd()
	if err != nil {
		t.Fatal(err.Error())
	}
}
