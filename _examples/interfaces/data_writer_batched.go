package main

import (
	"io/ioutil"
	"time"

	"github.com/kevburnsjr/batchy"
)

type dataWriterBatched struct {
	batcher batchy.Batcher
}

func (r *dataWriterBatched) Write(data []byte) error {
	return r.batcher.Add(data)
}

func NewDataWriterBatched(maxItems int, maxWait time.Duration) *dataWriterBatched {
	return &dataWriterBatched{batchy.New(maxItems, maxWait, func(items []interface{}) (errs []error) {
		var data []byte
		for _, d := range items {
			data = append(data, d.([]byte)...)
		}
		err := ioutil.WriteFile("test2", data, 0644)
		if err != nil {
			errs = make([]error, len(items))
			for i := range errs {
				errs[i] = err
			}
		}
		return
	})}
}
