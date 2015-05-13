package bitesized

import "github.com/garyburd/redigo/redis"

func (b *Bitesized) IsUserNew(user string) (bool, error) {
	userExists, err := redis.Bool(b.store.Do("HEXISTS", b.userListKey(), user))
	return !userExists, err
}
