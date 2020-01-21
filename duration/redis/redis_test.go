package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedisConn(t *testing.T) {
	must := assert.New(t)

	cli, err := NewRedisConn(
		Address("127.0.0.1:6379"),
	)
	must.Nil(err)

	err = cli.Set("k1", "v1", 0).Err()
	must.Nil(err)

	val, err := cli.Get("k1").Result()
	must.Nil(err)
	must.NotEmpty(val)
}
