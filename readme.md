# 多种抢红包模式

## demo1

简单抽奖模式

## demo2

彩票抽奖，双色球

## demo3

中奖概率类型抽奖

## demo4

抢红包模式

## demo5

抢红包模式，channel优化


## 压测使用wrk
wrk -t10 -c10 -d5 http://127.0.0.1:8080/get?id=17657329

```_asciidoc_
➜  log wrk -t10 -c10 -d5 http://127.0.0.1:8080/get?id=17657329
Running 5s test @ http://127.0.0.1:8080/get?id=17657329
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   421.07us    1.22ms  34.77ms   99.32%
    Req/Sec     2.94k     1.48k    6.65k    82.35%
  149330 requests in 5.10s, 18.99MB read
Requests/sec:  29283.19
Transfer/sec:      3.72MB
```