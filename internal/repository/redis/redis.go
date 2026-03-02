package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	client *redis.Client
}

func NewClient(ctx context.Context, addr string, dbName int, username, password string, maxRetries int, dialTimeout, readTimeout, writeTimeout time.Duration) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           dbName,
		Username:     username,
		Password:     password,
		MaxRetries:   maxRetries,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("redis.NewClient: %s", err.Error())
		return nil, err
	}

	return db, nil
}
