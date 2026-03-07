package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/internal/repository"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

const SessionIdAgeTime = int(24 * time.Hour / time.Second)

type AppData struct {
	clientId     string
	redirectUri  string
	clientSecret string
	name         string
	version      string
	devContact   string
}

func NewAppData(clientId, clientSecret, appName, appVersion, appDevContact, redirectUri string) *AppData {
	return &AppData{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
		name:         appName,
		version:      appVersion,
		devContact:   appDevContact,
	}
}

type OAuthHandler struct {
	ctx     context.Context
	repo    repository.Repository
	appData *AppData
	client  *http.Client
}

func NewOAuthHandler(ctx context.Context, appData *AppData, repo repository.Repository) *OAuthHandler {
	client := &http.Client{}
	return &OAuthHandler{
		ctx:     ctx,
		repo:    repo,
		appData: appData,
		client:  client,
	}
}

func (o *OAuthHandler) Begin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthBeginHandler is ready to work")

	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var sessionId string

	if sessionIdCookie == nil || sessionIdCookie.Value == "" {
		sessionId, err = generateSessionId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionId,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   SessionIdAgeTime,
		})

	} else {
		sessionId = sessionIdCookie.Value
	}

	state, err := generateState() // state новый для каждого запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//положить в redis
	err = o.repo.Set(o.ctx, fmt.Sprintf("state_%s", sessionId), state, time.Minute*5)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var redirectUrl *url.URL

	redirectUrl, err = url.Parse("https://hh.ru/oauth/authorize")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := redirectUrl.Query()
	params.Set("response_type", "code")
	params.Set("client_id", o.appData.clientId)
	params.Set("redirect_uri", o.appData.redirectUri)
	params.Set("state", state)
	redirectUrl.RawQuery = params.Encode()

	http.Redirect(w, r, redirectUrl.String(), http.StatusFound)
}

func (o *OAuthHandler) Complete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthEndHandler is ready to work")

	sessionIdCookie, err := r.Cookie("sessionId")

	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sessionId string

	if sessionIdCookie == nil || sessionIdCookie.Value == "" {
		http.Error(w, "oauth.Complete: invalid sessionId", http.StatusUnauthorized)
		return
	} else {
		sessionId = sessionIdCookie.Value
	}

	stateFromRepo, err := o.repo.Get(r.Context(), fmt.Sprintf("state_%s", sessionId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stateFromQuery := r.URL.Query().Get("state")
	if stateFromQuery == "" || stateFromRepo == "" || stateFromQuery != stateFromRepo {
		http.Error(w, "oauth.Complete: invalid sessionId", http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")

	params := url.Values{}
	params.Set("client_id", o.appData.clientId)
	params.Set("client_secret", o.appData.clientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	params.Set("redirect_uri", "/")

	request, err := http.NewRequestWithContext(
		o.ctx,
		http.MethodPost,
		"https://hh.ru/oauth/token",
		strings.NewReader(params.Encode()),
	)
	request.Header.Set("User-Agent", fmt.Sprintf("%s/%s (%s)", o.appData.name, o.appData.version, o.appData.devContact))
	response, err := o.client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "oauth.Complete: something went wrong in authorization", http.StatusUnauthorized)
		return
	}

	var results struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken := results.AccessToken
	refreshToken := results.RefreshToken
	expiresIn := results.ExpiresIn

	err = o.repo.Set(o.ctx, fmt.Sprintf("access_token_%s", sessionId), accessToken, time.Second*time.Duration(expiresIn))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = o.repo.Set(o.ctx, fmt.Sprintf("refresh_token_%s", sessionId), refreshToken, 9999999999999999)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "authorized",
		Value:    "true",
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   SessionIdAgeTime,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

// Функции генерации sessionId и state
func generateState() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("oauth.generateState: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func generateSessionId() (string, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("oauth.generateSessionId: %w", err)
	}
	sessionID := id.String()

	return sessionID, nil
}
