# bitesized [![Build Status](https://travis-ci.org/sent-hil/bitesized.svg?branch=master)](https://travis-ci.org/sent-hil/bitesized)

bitesized is a library that uses redis's bit operations to store and calculate analytics. It comes with a http server that can be used as an stand alone api.

## Motivation

It started when I saw a [blog post](http://blog.getspool.com/2011/11/29/fast-easy-realtime-metrics-using-redis-bitmaps/) about using redis bitmaps to store user retention data. It sounded pretty neat and simple, not to mention fun, to implement.

## Install

`go get github.com/sent-hil/bitesized`

## Usage

Initialize client:

```go
package main

import (
  "log"

  "github.com/garyburd/redigo/redis"
  "github.com/sent-hil/bitesized"
)

func main() {
  redisuri := "localhost:6379"
  client, err := bitesized.NewClient(redisuri)
}
```

Optionally, set intervals you want to track; by default these intervals are tracked: `hourly, daily, weekly and monthly`:

```go
client.Intervals = []Interval{
  bitesized.Hour, bitesized.Daily, bitesized.Week, bitesized.Month,
}
```

Optionally, set prefix to use for ALL keys; defaults to `bitesized`:

```go
client.KeyPrefix = "bitesized"
```

Track an event that an user did:

```go
err = client.TrackEvent("dodge rock", "indianajones", time.Now())
```

Get count of users who did an event on particular interval:

```go
count, err = client.CountEvent("dodge rock", time.Now(), bitesized.Hour)
```

Check if user did an event for particular interval:

```go
didEvent, err := client.DidEvent("dodge rock", "indianajones", time.Now(), bitesized.Hour)
```

Get a metric:

```go
from := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
till := time.Date(2015, time.January, 3, 0, 0, 0, 0, time.UTC)

rs, err := client.GetRetention("dodge rock", bitesized.Daily, from, till)
```

This returns a result like below. Result key is sorted asc by time:

```
{
    "2015-01-01 00:00": [30, 17, 60],
    "2015-01-02 00:00": [49, 24,  0],
    "2015-01-03 00:00": [67,  0,  0]
}
```
