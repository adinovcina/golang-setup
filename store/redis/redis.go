package redisstore

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	r "github.com/redis/go-redis/v9"
)

// RedisStore - New Redis RedisStore.
type RedisStore struct {
	redis *r.Client
	rs    *redsync.Redsync
}

func New(redis *r.Client) *RedisStore {
	pool := goredis.NewPool(redis)

	// Create an instance of redisync to be used to obtain a mutual exclusion lock.
	rs := redsync.New(pool)

	return &RedisStore{
		redis: redis,
		rs:    rs,
	}
}
