
This package makes use of Go's empty interface `interface{}`. For this reason, it is best not to export
any `Batcher` directly from your package. Instead the batcher should be hidden behind an existing synchronous
interface.

Suppose you have the following code that writes bytes to a file:

```go
package repo

import (
	"io/ioutil"
)

type DataWriter interface {
	Write(data []byte) error
}

type dataWriter struct{}

func (r *dataWriter) Write(data []byte) error {
	return ioutil.WriteFile("test1", data, 0644)
}

func NewDataWriter() *dataWriter {
	return &dataWriter{}
}
```

You could create a batched version that satisfies the same interface:
```go
package repo

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
```

Now during dependency injection just replace

```go
dw := repo.NewDataWriter()
dw.Write([]byte("asdf"))
```

with

```go
dw := repo.NewDataWriterBatched()
dw.Write([]byte("asdf"))
```

and your code shouldn't need to know the difference because you've used interfaces to effectively hide the
implementation details (in this case, the use of batching).