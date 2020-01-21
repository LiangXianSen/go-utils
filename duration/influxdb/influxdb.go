package influxdb

import (
	"time"

	"github.com/influxdata/influxdb1-client/v2"
)

type Client struct {
	client.Client
}

type Options struct {
	Addr     string
	Username string
	Password string
	Timeout  time.Duration
}

// NewInfluxDBConn returns a connection of influxdb by options which provides
func NewInfluxDBConn(opt ...Option) (*Client, error) {
	var (
		opts Options
		conf client.HTTPConfig
	)

	for _, o := range opt {
		o(&opts)
	}

	conf.Addr = opts.Addr
	conf.Username = opts.Username
	conf.Password = opts.Password
	conf.Timeout = opts.Timeout
	cli, err := client.NewHTTPClient(conf)
	if err != nil {
		return nil, err
	}

	if _, _, err := cli.Ping(defaultConnectTimeout); err != nil {
		return nil, err
	}

	return &Client{
		Client: cli,
	}, nil
}

// QueryAsMap requests influx responds as a map
func (c *Client) QueryAsMap(query string, database string, precision string) ([]map[string]interface{}, error) {
	var res []map[string]interface{}

	q := client.NewQuery(query, database, "")
	response, err := c.Query(q)
	if err != nil {
		return nil, err
	}

	if response.Error() != nil {
		return nil, response.Error()
	}

	for _, row := range response.Results[0].Series {
		for _, v := range row.Values {
			line := make(map[string]interface{})
			for i := range row.Columns {
				line[row.Columns[i]] = v[i]
			}
			res = append(res, line)
		}
	}

	return res, nil
}

const (
	defaultConnectTimeout = time.Second * 1
)

type Option func(*Options)

// Address sets address of influxdb
func Address(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

// Username sets username of influxdb
func Username(user string) Option {
	return func(o *Options) {
		o.Username = user
	}
}

// Password sets password of influxdb
func Password(pwd string) Option {
	return func(o *Options) {
		o.Password = pwd
	}
}
