# bitesized [![Build Status](https://travis-ci.org/sent-hil/bitesized.svg?branch=master)](https://travis-ci.org/sent-hil/bitesized)

bitesized is a library that uses redis's bit operations to store and calculate analytics. It comes with a http server that can be used as an stand alone api.

## Motivation

It started when I saw a [blog post](http://blog.getspool.com/2011/11/29/fast-easy-realtime-metrics-using-redis-bitmaps/) about using redis bitmaps to store user retention data. It sounded pretty neat and simple, not to mention fun, to implement.

## Install

`go get github.com/sent-hil/bitesized`

## Usage

Initialize client with Redis connection and defaults:

```go
package main

import (
  "log"

  "github.com/garyburd/redigo/redis"
  "github.com/sent-hil/bitesized"
)

func main() {
  // initialize client
  redisuri := "localhost:6379"
  client, err := bitesized.NewClient(rs)
  if err != nil {
    log.Fatal(err)
  }
}
```

Optionally set timings you want to track; by default these are created:
`hourly, daily, weekly and monthly`

```go
client.Intervals = []Interval{
  bitesized.Hour, bitesized.Daily, bitesized.Week, bitesized.Month,
}
```

Optionally set prefix to use for ALL keys; defaults to `bitesized`

```go
client.KeyPrefix = "bitesized"
```

```go
// initialize new event with name, username & optional timestamp;
// last arg is optional, defaults to `time.Now()`
event := bitesized.NewEvent("dodge rock", "indianajones", nil)

// track an event, ie. save to redis
err = client.Track(event)
if err != nil {
  return err
}
```

Get a metric:

```go
from := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
till := time.Date(2015, time.January, 3, 0, 0, 0, 0, time.UTC)

rs, err := bitesized.GetMetricRetention(event, bitesized.Daily, from, till)
if err != nil {
  return err
}
```

This returns a result like below. Result key is sorted desc by time:

```
{
    "2015-01-03 00:00": [30, 17, 60],
    "2015-01-02 00:00": [49, 24,  0],
    "2015-01-01 00:00": [67,  0,  0]
}
```

Get list of events seen by library:

```go
rs, err := bitesized.GetMetricEvents("dodge*")
if err != nil {
  return err
}
```

This returns a result like below.

```
["dodge rock", "dodge snakes", "dodge nazis"]
```
