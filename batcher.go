package batchy

import (
	"sync"
	"time"
)

// Batcher can add an item, returning the corresponding error
type Batcher interface {
	// Add adds an item to the current batch
	Add(interface{}) error

	// Stop stops the batcher
	Stop()
}

// Processor is a function that accepts items and returns a corresponding array of errors
type Processor func(items []interface{}) []error

type batcher struct {
	processor Processor
	itemLimit int
	waitTime  time.Duration
	mutex     *sync.Mutex
	stopped   bool
	batch     *batch
}

// New returns a new batcher
// - itemLimit indicates the maximum number of items per batch
// - waitTime indicates the amount of time to wait before processing a non-full batch
// - processor is the processing function to call for the batch
func New(itemLimit int, waitTime time.Duration, processor Processor) *batcher {
	return &batcher{
		processor: processor,
		itemLimit: itemLimit,
		waitTime:  waitTime,
		mutex:     &sync.Mutex{},
	}
}

func (b *batcher) newBatch() {
	var ba = newBatch(b.processor)
	go func(ba *batch) {
		select {
		case <-time.After(b.waitTime):
			b.mutex.Lock()
			if ba.process() {
				b.batch = nil
			}
			b.mutex.Unlock()
		case <-ba.waiting:
		}
	}(ba)
	b.batch = ba
}

// Add adds an item to the current batch
func (b *batcher) Add(item interface{}) (err error) {
	b.mutex.Lock()
	if b.stopped {
		return ErrBatcherStopped
	}
	if b.batch == nil {
		b.newBatch()
	}
	var ba = b.batch
	var i = len(ba.items)
	ba.items = append(ba.items, item)
	if len(ba.items) == b.itemLimit && ba.process() {
		b.batch = nil
	}
	b.mutex.Unlock()
	<-ba.done
	if len(ba.resp) > i {
		err = ba.resp[i]
	}
	return
}

// Stop stops the batcher
func (b *batcher) Stop() {
	b.mutex.Lock()
	b.stopped = true
	var ba = b.batch
	if ba == nil {
		return
	}
	ba.process()
	b.mutex.Unlock()
	<-ba.done
}
