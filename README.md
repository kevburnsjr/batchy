# Batchy

A nice little library for concurrent batching

[![GoDoc](https://godoc.org/github.com/kevburnsjr/batchy?status.svg)](https://godoc.org/github.com/kevburnsjr/batchy)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevburnsjr/batchy?2)](https://goreportcard.com/report/github.com/kevburnsjr/batchy)
[![Code Coverage](http://gocover.io/_badge/github.com/kevburnsjr/batchy?2)](http://gocover.io/github.com/kevburnsjr/batchy)

The throughput of APIs, web services and background workers can sometimes be improved by orders of magnitude
through the introduction of artificial latency in support of concurrent batching. When latency and batch size
are well tuned, the client may not even experience added latency in most cases. These efficiency improvements
can improve stability and scalability while lowering server costs.

This is a general purpose library for concurrent batching of any sort of operation one might desire.

The [example](examples/example.go) below illustrates an HTTP server writing strings to a file.

```go
package main

import (
	"net/http"
	"time"
	"os"
	"strings"
	"errors"

	"github.com/kevburnsjr/batchy"
)

type unbatchedHandler struct {}

func (h unbatchedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := appendIdToFile(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type batchedHandler struct {
	batcher batchy.Batcher
}

func (h batchedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.batcher.Add(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	http.Handle("/unbatched", unbatchedHandler{})
	http.Handle("/batched", batchedHandler{
		// Max batch size 100
		// Max wait time 100 milliseconds
		batchy.New(100, 100 * time.Millisecond, func(items []interface{}) (errs []error) {
			var ok bool
			errs = make([]error, len(items))
			var ids = make([]string, len(items))
			for i, v := range items {
				ids[i], ok = v.(string)
				if !ok {
					errs[i] = errors.New("Wrong data type")
				}
			}
			err := appendIdsToFile(ids)
			if err != nil {
				for i := range items {
					errs[i] = err
				}
			}
			return
		}),
	})
	http.ListenAndServe(":8080", nil)
}

func appendIdsToFile(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	f, err := os.OpenFile("items", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(strings.Join(ids, "\n") + "\n"))
	if err != nil {
		return
	}
	return
}

func appendIdToFile(id string) (err error) {
	f, err := os.OpenFile("items", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(id + "\n"))
	if err != nil {
		return
	}
	return
}
```

Unbatched (2,309 req/s)

```
> ab -k -n10000 -c500 localhost:8080/unbatched?id=123

Concurrency Level:      500
Time taken for tests:   4.329 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Keep-Alive requests:    10000
Total transferred:      990000 bytes
HTML transferred:       0 bytes
Requests per second:    2309.75 [#/sec] (mean)
Time per request:       216.473 [ms] (mean)
Time per request:       0.433 [ms] (mean, across all concurrent requests)
Transfer rate:          223.31 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   3.2      0      17
Processing:     1  213 100.8    233     584
Waiting:        1  213 100.8    233     583
Total:          1  214 100.8    236     586

Percentage of the requests served within a certain time (ms)
  50%    236
  66%    253
  75%    262
  80%    281
  90%    363
  95%    386
  98%    421
  99%    457
 100%    586 (longest request)
```

Batched (13,829 req/s)

```
> ab -k -n10000 -c500 localhost:8080/batched?id=123

Concurrency Level:      500
Time taken for tests:   0.723 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Keep-Alive requests:    10000
Total transferred:      990000 bytes
HTML transferred:       0 bytes
Requests per second:    13829.44 [#/sec] (mean)
Time per request:       36.155 [ms] (mean)
Time per request:       0.072 [ms] (mean, across all concurrent requests)
Transfer rate:          1337.03 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   3.1      0      17
Processing:     2   26  11.4     25     120
Waiting:        2   26  11.3     25     104
Total:          2   27  11.7     26     120

Percentage of the requests served within a certain time (ms)
  50%     26
  66%     29
  75%     32
  80%     33
  90%     38
  95%     44
  98%     57
  99%     94
 100%    120 (longest request)
```
