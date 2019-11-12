# Disk Write Batching

This example illustrates 6x throughput for an HTTP server writing strings to local disk.

See [main.go](./main.go)

## 500 Concurrency
```
Unbatched:  2,309 req/s, stdev 236 ms, 0% failure
Batched:   13,829 req/s, stdev  12 ms, 0% failure
```

---

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
