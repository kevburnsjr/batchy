package batchy

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func processorFailsEven(items []interface{}) (resp []error) {
	resp = make([]error, len(items))
	for i, v := range items {
		if v.(int)%2 == 0 {
			resp[i] = errors.New("even")
		}
	}
	return
}

// Batch should process when full
func TestBatcherFullBatch(t *testing.T) {
	b := New(4, time.Second, processorFailsEven)
	var errors int64
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			err := b.Add(i)
			if err != nil {
				atomic.AddInt64(&errors, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if errors != 2 {
		t.Fatal("Expected 2 errors")
	}
}

// Batch should process after timeout
func TestBatcherTimeout(t *testing.T) {
	b := New(4, 10*time.Millisecond, processorFailsEven)
	var errors int64
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			err := b.Add(i)
			if err != nil {
				atomic.AddInt64(&errors, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if errors != 1 {
		t.Fatal("Expected 1 error")
	}
}

// Batch should immediately process after call to Stop
func TestBatcherStop(t *testing.T) {
	start := time.Now()
	b := New(4, time.Second, processorFailsEven)
	var errors int64
	wg := sync.WaitGroup{}
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			err := b.Add(i)
			if err != nil {
				atomic.AddInt64(&errors, 1)
			}
			wg.Done()
		}(i)
	}
	time.Sleep(10 * time.Millisecond)
	b.Stop()
	wg.Wait()
	if time.Now().Sub(start) >= time.Second {
		t.Fatal("Expected stop to work")
	}
	if errors != 3 {
		t.Fatal("Expected 3 errors")
	}
	b.Stop()
	err := b.Add(1)
	if err.Error() != ErrBatcherStopped.Error() {
		t.Fatal("Expected " + ErrBatcherStopped.Error())
	}
}

// Batch should Stop cleanly even with no outstanding batches
func TestBatcherStopEmpty(t *testing.T) {
	start := time.Now()
	b := New(4, time.Second, processorFailsEven)
	var errors int64
	wg := sync.WaitGroup{}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			err := b.Add(i)
			if err != nil {
				atomic.AddInt64(&errors, 1)
			}
			wg.Done()
		}(i)
	}
	time.Sleep(10 * time.Millisecond)
	b.Stop()
	wg.Wait()
	if time.Now().Sub(start) >= time.Second {
		t.Fatal("Expected stop to work")
	}
	if errors != 4 {
		t.Fatal("Expected 4 errors")
	}
}

// Batcher should handle thousands of jobs efficiently
func TestBatcherStress(t *testing.T) {
	b := New(1000, 10*time.Millisecond, processorFailsEven)
	var errors int64
	wg := sync.WaitGroup{}
	for i := 0; i < 100001; i++ {
		wg.Add(1)
		go func(i int) {
			err := b.Add(i)
			if err != nil {
				atomic.AddInt64(&errors, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if errors != 50001 {
		t.Fatal("Expected 50001 errors")
	}
}
