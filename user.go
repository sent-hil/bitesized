package bitesized

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// IsUserNew returns if an user has been seen by the library or not.
func (b *Bitesized) IsUserNew(user string) (bool, error) {
	userExists, err := redis.Bool(b.store.Do("HEXISTS", b.userListKey(), user))
	return !userExists, err
}

// EventUsers returns list of users who did a given event for given interval and time.
func (b *Bitesized) EventUsers(e string, t time.Time, i Interval) ([]string, error) {
	key := b.intervalkey(e, t, i)
	str, err := redis.String(b.store.Do("GET", key))
	if err != nil {
		return []string{}, err
	}

	idTobools := bitStringToBools(str)

	key = b.userIdListKey()
	args := []interface{}{key}

	for userIndex, userDidEvent := range idTobools {
		if userDidEvent {
			args = append(args, userIndex)
		}
	}

	return redis.Strings(b.store.Do("HMGET", args...))
}

// RemoveUser unsets all events did by an user.
func (b *Bitesized) RemoveUser(user string) error {
	eventkeys, err := redis.Strings(b.store.Do("KEYS", b.allEventsKey()))
	if err != nil {
		return err
	}

	offset, err := b.getOrSetUser(user)
	if err != nil {
		return err
	}

	b.store.Send("MULTI")

	for _, event := range eventkeys {
		b.store.Send("SETBIT", event, offset, Off)
	}

	_, err = b.store.Do("EXEC")

	return err
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
