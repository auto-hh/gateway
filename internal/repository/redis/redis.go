package redis

import (
	"context"
	"fmt"
	"gateway/config/modules"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	client *redis.Client
}

func NewRepo(ctx context.Context, repoConfig modules.RepoConfig) *Repo {

	client, err := NewClient(ctx, repoConfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("redis.NewRepo: %v", err))
	}

	return &Repo{
		client: client,
	}
}

func NewClient(
	ctx context.Context,
	repoConfig modules.RepoConfig,
) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:         repoConfig.GetAddr(),
		DB:           repoConfig.GetDbName(),
		Username:     repoConfig.GetUser(),
		Password:     repoConfig.GetPassword(),
		MaxRetries:   repoConfig.GetMaxRetries(),
		DialTimeout:  repoConfig.GetDialTimeout(),
		ReadTimeout:  repoConfig.GetReadTimeout(),
		WriteTimeout: repoConfig.GetWriteTimeout(),
	})

	if err := db.Ping(ctx).Err(); err != nil {
		fmt.Printf("redis.NewClient: %s", err.Error())
		return nil, err
	}

	return db, nil
}

func (r *Repo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Repo) Set(ctx context.Context, key string, value string, expirationTime time.Duration) error {
	return r.client.Set(ctx, key, value, expirationTime).Err()
}

func (r *Repo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
