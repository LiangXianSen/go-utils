package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInfluxDBConn(t *testing.T) {
	must := assert.New(t)
	_, err := NewInfluxDBConn(
		Address("http://127.0.0.1:8086"),
	)
	must.Nil(err)
}

func TestQueryAsMap(t *testing.T) {
	must := assert.New(t)
	cli, err := NewInfluxDBConn(
		Address("http://127.0.0.1:8086"),
	)
	must.Nil(err)

	const SQL = `SELECT * FROM system limit 2`

	res, err := cli.QueryAsMap(SQL, "impactor_raw", "")
	must.Nil(err)
	must.Equal(2, len(res))
	for _, v := range res {
		must.NotNil(v["time"])
	}
}
