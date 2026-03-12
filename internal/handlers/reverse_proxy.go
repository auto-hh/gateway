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

			err = ph.service.DoTokenRequest(r.Context(), sessionId.Value)
			if err != nil {
				http.Error(w, fmt.Sprintf("handlers.ProxyHandler.Handler: %v", err.Error()), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "authorized",
				Value:    "true",
				Path:     "/",
				HttpOnly: false,
				Secure:   false,
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
