package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gateway/config/modules"
	"gateway/internal/repository"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	repo          repository.Repository
	hhConfig      *modules.HHConfig
	timeoutConfig *modules.TimeoutConfig
}

func NewService(repo repository.Repository, hhConfig *modules.HHConfig, timeoutConfig *modules.TimeoutConfig) *Service {
	return &Service{
		repo:          repo,
		hhConfig:      hhConfig,
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

func (s *Service) BuildTokenRequest(ctx context.Context, sessionId string, stateFromQuery, code string) (*http.Request, error) {
	stateFromRepo, err := s.repo.Get(ctx, fmt.Sprintf("state_%s", sessionId))
	if err != nil {
		return &http.Request{}, err
	}

	if stateFromQuery == "" || stateFromRepo == "" || stateFromQuery != stateFromRepo {
		return &http.Request{}, fmt.Errorf("service.BuildTokenRequest: invalid state")
	}

	params := url.Values{}
	params.Set("client_id", s.hhConfig.GetClientId())
	params.Set("client_secret", s.hhConfig.GetClientSecret())
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	params.Set("redirect_uri", "/")

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://hh.ru/oauth/token",
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
func generateState() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("oauth.generateState: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
