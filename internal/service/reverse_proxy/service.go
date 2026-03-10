package reverse_proxy

import (
	"context"
	"encoding/json"
	"errors"
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
	repo     *redis.Repo
	hhConfig *modules.HHConfig
}

func NewService(repo *redis.Repo, hhConfig *modules.HHConfig) *Service {
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

func (s *Service) buildRefreshTokensRequest(ctx context.Context, sessionId string) (*http.Request, error) {
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

func (s *Service) setTokens(ctx context.Context, sessionId, accessToken, refreshToken string, expiresIn time.Duration) error {
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

func (s *Service) DoTokenRequest(ctx context.Context, sessionId string) (err error) {
	request, err := s.buildRefreshTokensRequest(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("service.DoTokenRequest: %w", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("service.DoTokenRequest: %w, status:%s", err, http.StatusText(response.StatusCode))
	}

	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			if err != nil {
				err = errors.Join(err, fmt.Errorf("service.DoTokenRequest: %w", closeErr))
			} else {
				err = closeErr
			}
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("service.DoTokenRequest: status %s", http.StatusText(response.StatusCode))
	}

	var results struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return fmt.Errorf("service.DoTokenRequest: %w", err)
	}

	err = s.setTokens(
		ctx,
		sessionId,
		results.AccessToken,
		results.RefreshToken,
		time.Duration(results.ExpiresIn)*time.Second,
	)
	if err != nil {
		return fmt.Errorf("service.DoTokenRequest: %w", err)
	}
	return nil
}
