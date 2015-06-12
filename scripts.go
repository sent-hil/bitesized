package bitesized

// Defines lua scripts that are used by library.
var getOrSetUserScript = `
if redis.call('HEXISTS', KEYS[1], KEYS[2]) == 1 then
  return redis.call('HGET', KEYS[1], KEYS[2])
else
  local id = redis.call('INCR', KEYS[3])
  redis.call('HSET', KEYS[1], KEYS[2], id)
  redis.call('HSET', KEYS[4], id, KEYS[2])

  return id
end`
