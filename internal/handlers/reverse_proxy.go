package handlers

import (
	"errors"
	"fmt"
	"gateway/config/modules"
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
		slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: succesfully redirecting to backend"))
		ph.backendProxy.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: finished work"))
}
