package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// ConnectionData representing all configuration values neede to establish connection with Redis.
type ConnectionData struct {
	Address  string
	Database string
	Password string
}

// Init creates new instance.
func Init(connData ConnectionData) (*redis.Client, error) {
	databaseIndex, err := strconv.Atoi(connData.Database)
	if err != nil {
		return nil, err
	}

	c := redis.NewClient(&redis.Options{
		Addr:     connData.Address,
		Password: connData.Password,
		DB:       databaseIndex,
	})

	// Test connection to the server
	_, err = c.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return c, nil
}
