package batchy

import (
	"sync"
)

type batch struct {
	processor  Processor
	items      []interface{}
	waiting    chan struct{}
	done       chan struct{}
	mutex      *sync.Mutex
	processing bool
	resp       []error
}

func newBatch(p Processor) *batch {
	var b = &batch{
		processor: p,
		items:     []interface{}{},
		waiting:   make(chan struct{}),
		done:      make(chan struct{}),
		mutex:     &sync.Mutex{},
	}
	return b
}

func (b *batch) process() bool {
	b.mutex.Lock()
	if b.processing {
		return false
	}
	close(b.waiting)
	b.processing = true
	b.mutex.Unlock()
	go func() {
		b.resp = b.processor(b.items)
		close(b.done)
	}()
	return true
}
