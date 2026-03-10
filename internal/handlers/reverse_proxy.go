package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/config/modules"
	"gateway/internal/repository"
	"gateway/internal/utils"
	"log/slog"
	"net/http"
	"net/http/httputil"
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

type ProxyHandler struct {
	backendHost   string
	frontendHost  string
	backendProxy  *httputil.ReverseProxy
	frontendProxy *httputil.ReverseProxy

	repo     repository.Repository
	hhConfig *modules.HHConfig
}

func NewProxyHandler(baseConfig *modules.BaseConfig) *ProxyHandler {
	return &ProxyHandler{
		backendHost:   baseConfig.GetBackendHost(),
		frontendHost:  baseConfig.GetFrontendHost(),
		backendProxy:  httputil.NewSingleHostReverseProxy(baseConfig.GetBackendUrl()),
		frontendProxy: httputil.NewSingleHostReverseProxy(baseConfig.GetFrontendUrl()),
	}
}

func (ph *ProxyHandler) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlers.ProxyHandler.Handler is ready to work")
	sessionId, err := r.Cookie(utils.CookieKeySessionId)

	if errors.Is(err, http.ErrNoCookie) { //проверка на наличие в куках
		slog.Info(fmt.Sprintf("handlers.ProxyHandler: Redirecting to authorization, %v", err))
		http.Redirect(w, r, "/oauth/begin/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	target := r.URL.Host
	slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: user with id %s exists, wants to see %s", sessionId, target))

	if target == ph.frontendHost {
		slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: succesfully redirecting to frontend"))
		ph.frontendProxy.ServeHTTP(w, r)
	} else if target == ph.backendHost {

		accessStatus, err := ph.CheckToken(r.Context(), sessionId.Value)

		if err != nil {
			http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		switch accessStatus {
		case Allow:
			slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: succesfully redirecting to backend"))
			ph.backendProxy.ServeHTTP(w, r)
		case Update:
			request, err := ph.BuildRefreshTokensRequest(r.Context(), sessionId.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
			}
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
				return
			}

			defer func() {
				if err := response.Body.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}()

			if response.StatusCode != http.StatusOK {
				http.Error(w, "handlers.ProxyHandler.Handler: something went wrong in authorization", http.StatusUnauthorized)
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

			err = ph.SetTokens(
				r.Context(),
				sessionId.Value,
				results.AccessToken,
				results.RefreshToken,
				time.Duration(results.ExpiresIn)*time.Second,
			)
			if err != nil {
				http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
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

			slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: succesfully redirecting to backend after authorization"))
			ph.backendProxy.ServeHTTP(w, r)
		default:
			http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", target), http.StatusUnauthorized)
			return

		}

	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: finished work"))
}

func (ph *ProxyHandler) CheckToken(ctx context.Context, sessionId string) (accessState AccessEnum, err error) {
	authTokenExists, err := ph.repo.Exists(ctx, fmt.Sprintf("access_token_%s", sessionId))

	if err != nil {
		accessState = -1
		return accessState, fmt.Errorf("handlers.CheckToken: %w", err)
	}

	if !authTokenExists {
		refreshTokenExists, err := ph.repo.Exists(ctx, fmt.Sprintf("refresh_token_%s", sessionId))
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

func (ph *ProxyHandler) BuildRefreshTokensRequest(ctx context.Context, sessionId string) (*http.Request, error) {
	refreshToken, err := ph.repo.Get(ctx, fmt.Sprintf("refresh_token_%s", sessionId))
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
	request.Header.Set("User-Agent", fmt.Sprintf("%s/%s (%s)", ph.hhConfig.GetAppName(), ph.hhConfig.GetAppVersion(), ph.hhConfig.GetDevContact()))

	return request, nil
}

func (ph *ProxyHandler) SetTokens(ctx context.Context, sessionId, accessToken, refreshToken string, expiresIn time.Duration) error {
	err := ph.repo.Set(ctx, fmt.Sprintf("access_token_%s", sessionId), accessToken, expiresIn)
	if err != nil {
		return fmt.Errorf("service.SetTokens: %w", err)
	}

	err = ph.repo.Set(ctx, fmt.Sprintf("refresh_token_%s", sessionId), refreshToken, expiresIn)
	if err != nil {
		return fmt.Errorf("service.SetTokens: %w", err)
	}
	return nil
}
