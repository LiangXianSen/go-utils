// go-redis v6 for golang v.1.12 redis client
package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// NewRedisConn returns a connection of redis by options which provides
// Without options will use default values
func NewRedisConn(opt ...Option) (*redis.Client, error) {
	opts := newOptions()

	for _, o := range opt {
		o(opts)
	}

	client := redis.NewClient(opts)
	_, err := client.Ping().Result()

	return client, err
}

const (
	defaultPoolSize     = 10
	defaultDialTimeout  = time.Second * 3
	defaultReadTimeout  = time.Second * 2
	defaultWriteTimeout = time.Second * 3
)

// newOptions returns default Options
func newOptions() *redis.Options {
	return &redis.Options{
		PoolSize:     defaultPoolSize,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		DialTimeout:  defaultDialTimeout,
	}
}

type Option func(*redis.Options)

// Address sets address of redis
func Address(addr string) Option {
	return func(o *redis.Options) {
		o.Addr = addr
	}
}

// Password sets password of redis
func Password(pwd string) Option {
	return func(o *redis.Options) {
		o.Password = pwd
	}
}

// DB selects db of redis
func DB(db int) Option {
	return func(o *redis.Options) {
		o.DB = db
	}
}

// PoolSize sets pool size
func PoolSize(size int) Option {
	return func(o *redis.Options) {
		o.PoolSize = size
	}
}
