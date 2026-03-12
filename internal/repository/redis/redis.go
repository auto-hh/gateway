package redis

import (
	"context"
	"fmt"
	"gateway/config/modules"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	client *redis.Client
}

func NewRepo(ctx context.Context, repoConfig *modules.RepoConfig) *Repo {
	client := redis.NewClient(&redis.Options{
		Addr:         repoConfig.GetAddr(),
		DB:           repoConfig.GetDbName(),
		Username:     repoConfig.GetUser(),
		Password:     repoConfig.GetPassword(),
		MaxRetries:   repoConfig.GetMaxRetries(),
		DialTimeout:  repoConfig.GetDialTimeout(),
		ReadTimeout:  repoConfig.GetReadTimeout(),
		WriteTimeout: repoConfig.GetWriteTimeout(),
	})

	return &Repo{
		client: client,
	}
}

func (r *Repo) Get(ctx context.Context, key string) (value string, err error) {
	value, err = r.client.Get(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("redis.Repo.Get: %v", err)
	}
	return
}

func (r *Repo) Set(ctx context.Context, key string, value string, expirationTime time.Duration) error {
	err := r.client.Set(ctx, key, value, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("redis.Repo.Set: %v", err)
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis.Repo.Del: %v", err)
	}
	return nil
}

func (r *Repo) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()

	if err != nil {
		return false, fmt.Errorf("redis.Exists: %w", err)
	}

	return result > 0, nil
}
