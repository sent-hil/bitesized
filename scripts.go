package bitesized

var getOrsetUserScript = `
if redis.call('HEXISTS', KEYS[1], KEYS[2]) == 1 then
  return redis.call('HGET', KEYS[1], KEYS[2])
else
  local id = redis.call('INCR', KEYS[3])
  redis.call('HSET', KEYS[1], KEYS[2], id)
  return id
end`
