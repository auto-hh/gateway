package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway/config/modules"
	"gateway/internal/service/reverse_proxy"
	"gateway/internal/utils"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"
)

type ProxyHandler struct {
	backendHost   string
	frontendHost  string
	backendProxy  *httputil.ReverseProxy
	frontendProxy *httputil.ReverseProxy

	service *reverse_proxy.Service
}

func NewProxyHandler(baseConfig *modules.BaseConfig, service *reverse_proxy.Service) *ProxyHandler {
	return &ProxyHandler{
		backendHost:   baseConfig.GetBackendHost(),
		frontendHost:  baseConfig.GetFrontendHost(),
		backendProxy:  httputil.NewSingleHostReverseProxy(baseConfig.GetBackendUrl()),
		frontendProxy: httputil.NewSingleHostReverseProxy(baseConfig.GetFrontendUrl()),
		service:       service,
	}
}

func (ph *ProxyHandler) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlers.ProxyHandler.Handler is ready to work")
	sessionId, err := r.Cookie(utils.CookieKeySessionId)

	if errors.Is(err, http.ErrNoCookie) {
		slog.Info(fmt.Sprintf("handlers.ProxyHandler: Redirecting to authorization, %v", err))
		http.Redirect(w, r, "/oauth/begin/", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusBadRequest) //TODO: чзх что с обработкой ошибок?
		return
	}

	target := r.URL.Host
	slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: user with id %s exists, wants to see %s", sessionId, target))

	if target == ph.frontendHost {
		slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: succesfully redirecting to frontend"))
		ph.frontendProxy.ServeHTTP(w, r)
	} else if target == ph.backendHost {

		accessStatus, err := ph.service.CheckToken(r.Context(), sessionId.Value)

		if err != nil {
			http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		switch accessStatus {
		case reverse_proxy.Allow:
			slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: succesfully redirecting to backend"))
			ph.backendProxy.ServeHTTP(w, r)
		case reverse_proxy.Update:
			request, err := ph.service.BuildRefreshTokensRequest(r.Context(), sessionId.Value)
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

			err = ph.service.SetTokens(
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
