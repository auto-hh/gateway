package handlers

import (
	"errors"
	"fmt"
	"gateway/config/modules"
	"gateway/internal/service/reverse_proxy"
	"gateway/internal/utils"
	"log/slog"
	"net/http"
	"net/http/httputil"
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
	target := r.Host
	slog.Info("handlers.ProxyHandler.Handler call", slog.String("host", target))

	switch target {
	case ph.frontendHost:
		slog.Info("handlers.ProxyHandler.Handler: succesfully redirecting to frontend")
		ph.frontendProxy.ServeHTTP(w, r)

	case ph.backendHost:
		sessionId, err := r.Cookie(utils.CookieKeySessionId)

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				slog.Info("unauthorized")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			} else {
				slog.Error("r.Cookie(utils.CookieKeySessionId)", slog.String("err", err.Error()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		accessStatus, err := ph.service.CheckToken(r.Context(), sessionId.Value)

		if err != nil {
			http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
			return
		}

		switch accessStatus {
		case reverse_proxy.Allow:
			slog.Info("handlers.ProxyHandler.Handler: succesfully redirecting to backend")
			ph.backendProxy.ServeHTTP(w, r)
		case reverse_proxy.Update:

			err = ph.service.DoTokenRequest(r.Context(), sessionId.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "auth",
				Value:    "true",
				Path:     "/",
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   SessionIdAgeTime,
			})

			slog.Info("handlers.ProxyHandler.Handler: succesfully redirecting to backend after authorization")
			ph.backendProxy.ServeHTTP(w, r)
		default:
			http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", target), http.StatusUnauthorized)
		}

	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
}
