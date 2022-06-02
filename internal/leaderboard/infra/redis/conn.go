package redis

import (
	"context"
	"leaderboard/config"

	goredis "github.com/go-redis/redis/v8"
)

// NewDial -
func NewDial(c config.Config) (*goredis.Client, error) {
	client, err := newDial(c.Redis)

	return client, err
}

func newDial(c config.Redis) (*goredis.Client, error) {
	ctx := context.Background()

	opts := goredis.Options{
		Addr:     c.Host,
		Password: c.Password,
		DB:       c.Database,
	}

	if c.MaxPoolSize > 0 {
		opts.PoolSize = int(c.MaxPoolSize)
	}

	if c.MaxRetries > 0 {
		opts.MaxRetries = int(c.MaxRetries)
	}

	if c.MinIdelConns > 0 {
		opts.MinIdleConns = int(c.MinIdelConns)
	}

	client := goredis.NewClient(&opts)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
