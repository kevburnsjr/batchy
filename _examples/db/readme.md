# Database Write Batching

This example illustrates 3x - 15x throughput improvement plus reduced failure rate

See [main.go](./main.go)

## 100 Concurrency

```
Unbatched: 2,690 req/s, stdev 78 ms, 0% failure
Batched:   8,060 req/s, stdev  7 ms, 0% failure
```

## 500 Concurrency

```
Unbatched:   540 req/s, stdev 1378 ms, 18% failure
Batched:   8,326 req/s, stdev   17 ms,  0% failure
```

---

c100 unbatched
```
> ab -n10000 -c100 localhost:8080/unbatched?id=123

Concurrency Level:      100
Time taken for tests:   3.716 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Total transferred:      750000 bytes
HTML transferred:       0 bytes
Requests per second:    2690.71 [#/sec] (mean)
Time per request:       37.165 [ms] (mean)
Time per request:       0.372 [ms] (mean, across all concurrent requests)
Transfer rate:          197.07 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.9      1       6
Processing:     1   36  78.0     30    1449
Waiting:        1   36  78.0     30    1449
Total:          1   37  78.1     31    1452

Percentage of the requests served within a certain time (ms)
  50%     31
  66%     36
  75%     40
  80%     42
  90%     49
  95%     55
  98%     66
  99%     77
 100%   1452 (longest request)
```

c100 batched
```
> ab -n10000 -c100 localhost:8080/batched?id=123

Concurrency Level:      100
Time taken for tests:   1.241 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Total transferred:      750000 bytes
HTML transferred:       0 bytes
Requests per second:    8060.68 [#/sec] (mean)
Time per request:       12.406 [ms] (mean)
Time per request:       0.124 [ms] (mean, across all concurrent requests)
Transfer rate:          590.38 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.9      2       6
Processing:     1    9   7.3      8     108
Waiting:        1    8   7.3      7     108
Total:          2   11   7.3     10     110

Percentage of the requests served within a certain time (ms)
  50%     10
  66%     11
  75%     12
  80%     12
  90%     14
  95%     16
  98%     18
  99%     19
 100%    110 (longest request)
```

c500 unbatched
```
> ab -n10000 -c500 localhost:8080/unbatched?id=123

Concurrency Level:      500
Time taken for tests:   18.496 seconds
Complete requests:      10000
Failed requests:        1826
   (Connect: 0, Receive: 0, Length: 1826, Exceptions: 0)
Write errors:           0
Non-2xx responses:      1826
Total transferred:      981902 bytes
HTML transferred:       60258 bytes
Requests per second:    540.67 [#/sec] (mean)
Time per request:       924.776 [ms] (mean)
Time per request:       1.850 [ms] (mean, across all concurrent requests)
Transfer rate:          51.84 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    4  43.7      1    1005
Processing:     1  267 1375.8     40   16077
Waiting:        1  266 1375.9     39   16077
Total:          1  271 1378.9     41   17063

Percentage of the requests served within a certain time (ms)
  50%     41
  66%     54
  75%     69
  80%     81
  90%   1014
  95%   1047
  98%   1122
  99%   2633
 100%  17063 (longest request)
```

c500 batched
```
> ab -n10000 -c500 localhost:8080/batched?id=123

Concurrency Level:      500
Time taken for tests:   1.201 seconds
Complete requests:      10000
Failed requests:        0
Write errors:           0
Total transferred:      750000 bytes
HTML transferred:       0 bytes
Requests per second:    8326.15 [#/sec] (mean)
Time per request:       60.052 [ms] (mean)
Time per request:       0.120 [ms] (mean, across all concurrent requests)
Transfer rate:          609.83 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    6   5.0      4      20
Processing:     8   37  16.1     35     226
Waiting:        7   33  15.3     32     224
Total:          9   43  16.7     41     227

Percentage of the requests served within a certain time (ms)
  50%     41
  66%     48
  75%     53
  80%     55
  90%     64
  95%     71
  98%     82
  99%     96
 100%    227 (longest request)
```