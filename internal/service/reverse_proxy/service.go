package reverse_proxy

import (
	"context"
	"fmt"
	"gateway/config/modules"
	"gateway/internal/repository/redis"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessEnum int

const (
	Update AccessEnum = iota
	Allow
	Deny
)

type Service struct {
	repo     redis.Repo
	hhConfig modules.HHConfig
}

func NewService(repo redis.Repo, hhConfig modules.HHConfig) *Service {
	return &Service{repo: repo, hhConfig: hhConfig}
}

func (s *Service) CheckToken(ctx context.Context, sessionId string) (accessState AccessEnum, err error) {
	authTokenExists, err := s.repo.Exists(ctx, fmt.Sprintf("access_token_%s", sessionId))

	if err != nil {
		accessState = -1
		return accessState, fmt.Errorf("handlers.CheckToken: %w", err)
	}

	if !authTokenExists {
		refreshTokenExists, err := s.repo.Exists(ctx, fmt.Sprintf("refresh_token_%s", sessionId))
		if err != nil {
			accessState = -1
			return accessState, fmt.Errorf("handlers.CheckToken: %w", err)
		}

		if !refreshTokenExists {
			return Deny, nil
		}

		return Update, nil
	}

	return Allow, nil
}

func (s *Service) BuildRefreshTokensRequest(ctx context.Context, sessionId string) (*http.Request, error) {
	refreshToken, err := s.repo.Get(ctx, fmt.Sprintf("refresh_token_%s", sessionId))
	if err != nil || refreshToken == "" {
		return nil, fmt.Errorf("handlers.BuildTokenRequest: %w", err)
	}

	params := url.Values{}

	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.hh.ru/oauth/token",
		strings.NewReader(params.Encode()),
	)
	request.Header.Set("User-Agent", fmt.Sprintf("%s/%s (%s)", s.hhConfig.GetAppName(), s.hhConfig.GetAppVersion(), s.hhConfig.GetDevContact()))

	return request, nil
}

func (s *Service) SetTokens(ctx context.Context, sessionId, accessToken, refreshToken string, expiresIn time.Duration) error {
	err := s.repo.Set(ctx, fmt.Sprintf("access_token_%s", sessionId), accessToken, expiresIn)
	if err != nil {
		return fmt.Errorf("service.SetTokens: %w", err)
	}

	err = s.repo.Set(ctx, fmt.Sprintf("refresh_token_%s", sessionId), refreshToken, expiresIn)
	if err != nil {
		return fmt.Errorf("service.SetTokens: %w", err)
	}
	return nil
}
