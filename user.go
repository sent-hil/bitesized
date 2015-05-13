package bitesized

import "github.com/garyburd/redigo/redis"

func (b *Bitesized) IsUserNew(user string) (bool, error) {
	userExists, err := redis.Bool(b.store.Do("HEXISTS", b.userListKey(), user))
	return !userExists, err
}

func (b *Bitesized) getOrSetUser(user string) (int, error) {
	user = dasherize(user)

	script := redis.NewScript(3, getOrSetUserScript)
	raw, err := script.Do(b.store, b.userListKey(), user, b.userCounterKey())

	return redis.Int(raw, err)
}
