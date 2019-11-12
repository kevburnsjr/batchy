package batchy

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Batcher should handle millions of jobs efficiently
func BenchmarkBatcher(b *testing.B) {
	for _, n := range []int{10, 20, 100, 1000} {
		b.Run(fmt.Sprintf("itemLimit_%d", n), func(b *testing.B) {
			benchmarkBatcher(b, n)
		})
	}
}

func benchmarkBatcher(b *testing.B, n int) {
	batch := New(n, 10*time.Millisecond, processorFailsEven)
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			batch.Add(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// Batcher should handle millions of jobs efficiently even with 100ms latency
func BenchmarkBatcher100ms(b *testing.B) {
	for _, n := range []int{10, 20, 100, 1000} {
		b.Run(fmt.Sprintf("itemLimit_%d", n), func(b *testing.B) {
			benchmarkBatcher100ms(b, n)
		})
	}
}

func benchmarkBatcher100ms(b *testing.B, n int) {
	batch := New(n, 10*time.Millisecond, func(items []interface{}) (resp []error) {
		time.Sleep(100 * time.Millisecond)
		return
	})
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			batch.Add(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
