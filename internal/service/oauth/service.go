package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gateway/config"
	"gateway/internal/repository"
	"net/url"
	"time"
)

type Service struct {
	repo          repository.Repository
	hhConfig      config.HHConfig
	timeoutConfig config.TimeoutConfig
}

func NewService(repo repository.Repository, hhconfig config.HHConfig, timeoutConfig config.TimeoutConfig) *Service {
	return &Service{
		repo:          repo,
		hhConfig:      hhconfig,
		timeoutConfig: timeoutConfig,
	}
}

func (s *Service) SetValueInRepository(ctx context.Context, key, value string, explicitIn time.Duration) error {
	err := s.repo.Set(ctx, key, value, explicitIn)
	return err
}

func (s *Service) BuildCodeRequest(ctx context.Context, sessionId string) (*url.URL, error) {
	state, err := generateState() // state новый для каждого запроса
	if err != nil {
		return nil, err
	}

	if err = s.repo.Set(ctx, fmt.Sprintf("state_%s", sessionId), state, s.timeoutConfig.GetStateExpirationTime()); err != nil {
		return nil, err
	}

	redirectUrl, err := url.Parse(s.hhConfig.GetRawUrl())
	if err != nil {
		return nil, fmt.Errorf("oauth.BuildRequest: %w", err)
	}

	params := redirectUrl.Query()
	params.Set("response_type", "code")
	params.Set("client_id", s.hhConfig.GetClientId())
	params.Set("redirect_uri", s.hhConfig.GetRedirectUri())
	params.Set("state", state)
	redirectUrl.RawQuery = params.Encode()

	return redirectUrl, nil
}

func generateState() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("oauth.generateState: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
