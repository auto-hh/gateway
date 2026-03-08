package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gateway/internal/repository"
	"net/url"
	"time"
)

type Service struct {
	repo                repository.Repository
	clientId            string
	redirectUri         string
	rawUrl              string
	clientSecret        string
	name                string
	version             string
	devContact          string
	StateExpirationTime time.Duration
}

func NewService(repo repository.Repository, clientId, clientSecret, rawUrl, name, version, devContact, redirectUri string, stateExpirationTime time.Duration) *Service {
	return &Service{
		repo:                repo,
		clientId:            clientId,
		clientSecret:        clientSecret,
		redirectUri:         redirectUri,
		rawUrl:              rawUrl,
		name:                name,
		version:             version,
		devContact:          devContact,
		StateExpirationTime: stateExpirationTime,
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

	if err = s.repo.Set(ctx, fmt.Sprintf("state_%s", sessionId), state, s.StateExpirationTime); err != nil {
		return nil, err
	}

	redirectUrl, err := url.Parse(s.rawUrl)
	if err != nil {
		return nil, fmt.Errorf("oauth.BuildRequest: %w", err)
	}

	params := redirectUrl.Query()
	params.Set("response_type", "code")
	params.Set("client_id", s.clientId)
	params.Set("redirect_uri", s.redirectUri)
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
