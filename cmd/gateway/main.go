package main

import (
	"context"
	"gateway/config"
	"gateway/internal/repository/redis"
	"gateway/internal/server"
	"gateway/internal/service/oauth"
	"log"
)

func main() {

	cfg := config.NewConfig()
	ctx := context.Background()
	repo := redis.NewRepo(ctx, cfg.GetRepoConfig())

	service := oauth.NewService(repo, cfg.GetHHConfig(), cfg.GetTimeoutConfig())
	srv := server.NewServer(cfg.GetBaseConfig(), service)

	err := srv.Start(cfg.GetBaseConfig().GetServerPort())
	if err != nil {
		log.Fatalf("server.Start failed: %v", err)
	}

}
