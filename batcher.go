package batchy

import (
	"sync"
	"time"
)

type Batcher interface {
	Add(interface{}) error
	Stop()
}

type Processor func(items []interface{}) []error

type batcher struct {
	processor Processor
	itemLimit int
	waitTime  time.Duration
	mutex     *sync.Mutex
	stopped   bool
	batch     *batch
}

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
	if len(ba.items) == b.itemLimit {
		go ba.process()
		b.batch = nil
	}
	b.mutex.Unlock()
	<-ba.done
	if len(ba.resp) > i {
		err = ba.resp[i]
	}
	return
}

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
