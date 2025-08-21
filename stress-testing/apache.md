# Apache Benchmark
## GET Requests
### 1,000 concurrent users ; 100,000 requests

Results from Apache Bench. The GET requests were made on `/api/home` route (to get whole menu), which implements `database calls` to get the complete menu.

```bash
ab -n 100000 -c 1000 -H "Cookie: token_id= <token value>" http://localhost:8100/api/home
```
<br />

***Result:-***

```bash
ab -n 100000 -c 1000 http://localhost:8100/api/home
This is ApacheBench, Version 2.3 <$Revision: 1903618 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 10000 requests
Completed 20000 requests
Completed 30000 requests
Completed 40000 requests
Completed 50000 requests
Completed 60000 requests
Completed 70000 requests
Completed 80000 requests
Completed 90000 requests
Completed 100000 requests
Finished 100000 requests


Server Software:
Server Hostname:        localhost
Server Port:            8100

Document Path:          /api/home
Document Length:        1146 bytes

Concurrency Level:      1000
Time taken for tests:   3.138 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      125600000 bytes
HTML transferred:       114600000 bytes
Requests per second:    31871.49 [#/sec] (mean)
Time per request:       31.376 [ms] (mean)
Time per request:       0.031 [ms] (mean, across all concurrent requests)
Transfer rate:          39092.38 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   21   3.9     21      34
Processing:   -30   23   5.5     23      52
Waiting:        0   16   3.9     16      41
Total:          0   44   5.3     44      65

Percentage of the requests served within a certain time (ms)
  50%     44
  66%     45
  75%     47
  80%     47
  90%     48
  95%     50
  98%     51
  99%     52
 100%     65 (longest request)
 ```

## POST Requests
### 1,000 concurrent users ; 10,000 requests

The POST requests were made on `/home/addtocart` route, which implements `authorization middleware` and multiple `database calls`.
```bash
 ab -n 10000 -c 1000 -p postdata.json -T application/json -H "Cookie: token_id= <token value>; session= <session value after atleast one item is added to cart to create sessions.order_id>" http://localhost:8100/home/addToCart
```
<br />

`postdata.json` has body values:-
```JSON
{
    "item_id":"1",
    "qnty":1
}
```
<br />

***Result:-***

```bash
Benchmarking localhost (be patient)
Completed 1000 requests
Completed 2000 requests
Completed 3000 requests
Completed 4000 requests
Completed 5000 requests
Completed 6000 requests
Completed 7000 requests
Completed 8000 requests
Completed 9000 requests
Completed 10000 requests
Finished 10000 requests


Server Software:
Server Hostname:        localhost
Server Port:            8100

Document Path:          /home/addToCart
Document Length:        36 bytes

Concurrency Level:      1000
Time taken for tests:   1.419 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1440000 bytes
Total body sent:        5440000
HTML transferred:       360000 bytes
Requests per second:    7049.09 [#/sec] (mean)
Time per request:       141.862 [ms] (mean)
Time per request:       0.142 [ms] (mean, across all concurrent requests)
Transfer rate:          991.28 [Kbytes/sec] received
                        3744.83 kb/s sent
                        4736.11 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   5.9      0      23
Processing:    10  133  83.6    115     745
Waiting:        1  133  83.7    115     745
Total:         15  136  83.3    116     761

Percentage of the requests served within a certain time (ms)
  50%    116
  66%    151
  75%    177
  80%    195
  90%    248
  95%    295
  98%    365
  99%    407
 100%    761 (longest request)
 ```