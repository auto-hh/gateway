package main

import (
	"context"
	"gateway/config"
	"gateway/internal/repository/redis"
	"gateway/internal/server"
	"gateway/internal/service/oauth"
	"gateway/internal/service/reverse_proxy"
	"log"
)

func main() {
	cfg := config.NewConfig()
	ctx := context.Background()
	repo := redis.NewRepo(ctx, cfg.GetRepoConfig())

	serviceOauth := oauth.NewService(repo, cfg.GetHHConfig(), cfg.GetTimeoutConfig())
	serviceReverseProxy := reverse_proxy.NewService(repo, cfg.GetHHConfig())

	srv := server.NewServer(cfg.GetBaseConfig(), serviceOauth, serviceReverseProxy)

	err := srv.Start(cfg.GetBaseConfig().GetServerPort())
	if err != nil {
		log.Fatalf("server.Start failed: %v", err)
	}

}
