package bitesized

import "github.com/garyburd/redigo/redis"

func (b *Bitesized) IsUserNew(user string) (bool, error) {
	userExists, err := redis.Bool(b.store.Do("HEXISTS", b.userListKey(), user))
	return !userExists, err
}

func (b *Bitesized) getOrSetUser(user string) (int, error) {
	user = dasherize(user)

	script := redis.NewScript(4, getOrSetUserScript)
	raw, err := script.Do(
		b.store, b.userListKey(), user, b.userCounterKey(), b.userIdListKey(),
	)

	return redis.Int(raw, err)
}

func (b *Bitesized) getUserById(id int) (string, error) {
	key := b.key(UserIdListKey)
	return redis.String(b.store.Do("HGET", key, id))
}
