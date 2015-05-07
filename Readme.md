# bitesized

bitesized is a library that uses redis's bit operations to store and calculate analytics. It comes with a http server that can be used as an stand alone api.

## Motivation

It started when I saw a [blog post](http://blog.getspool.com/2011/11/29/fast-easy-realtime-metrics-using-redis-bitmaps/) about using redis bitmaps to store user retention data. It sounded pretty neat and simple, not to mention fun, to implement.

## Install

`go get github.com/sent-hil/bitesized`
