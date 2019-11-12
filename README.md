# Batchy

A nice little package with no dependencies for fan-in batching of highly concurrent workloads

[![GoDoc](https://godoc.org/github.com/kevburnsjr/batchy?status.svg)](https://godoc.org/github.com/kevburnsjr/batchy)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevburnsjr/batchy?2)](https://goreportcard.com/report/github.com/kevburnsjr/batchy)
[![Code Coverage](http://gocover.io/_badge/github.com/kevburnsjr/batchy?2)](http://gocover.io/github.com/kevburnsjr/batchy)

The throughput of APIs, web services and background workers can sometimes be improved by orders of magnitude
through the introduction of artificial latency in support of concurrent batching. When latency and batch size
are well tuned, the client may not even experience added latency in most cases. These efficiency improvements
can result in increased service stability and total system throughput while lowering infrastructure costs.

This is a general purpose library for concurrent batching of any sort of operation one might desire. It could
be used to batch SQL inserts, API calls, disk writes, queue messages, stream records, emails, etc. The batcher
hides asynchronous processing behind a syncronous interface.

## How to use it

```go
// 100 max batch size
// 100 milliseconds max batch wait time
var table1 = batchy.New(100, 100*time.Millisecond, func(items []interface{}) (errs []error) {
	q := fmt.Sprintf(`INSERT INTO table1 (data) VALUES %s`,
		strings.Trim(strings.Repeat(`(?),`, len(items)), ","))
	_, err := db.Exec(q, items...)
	if err != nil {
		errs = make([]error, len(items))
		for i := range errs {
			errs[i] = err
		}
	}
	return
})
```
```go
// Call to Add blocks calling go routine for up to 100ms + processing time.
// If batch is filled before wait time expires, blocking will be reduced.
// Wait time begins when the first item is added to a batch.
err := table1.Add("data")
```

## Examples

See examples below for more complete integrations

- [Disk Write Batching](_examples/disk)
6x throughput improvement

- [Database Write Batching](_examples/db)
3x - 15x throughput improvement plus reduced failure rate
