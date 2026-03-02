package handlers

import (
	"fmt"
	"gateway/internal/utils"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	backendHost   string
	frontendHost  string
	backendProxy  *httputil.ReverseProxy
	frontendProxy *httputil.ReverseProxy
}

func NewProxyHandler(backendUrl, frontendUrl *url.URL, backendHost, frontendHost string) *ProxyHandler {
	return &ProxyHandler{
		backendHost:   backendHost,
		frontendHost:  frontendHost,
		backendProxy:  httputil.NewSingleHostReverseProxy(backendUrl),
		frontendProxy: httputil.NewSingleHostReverseProxy(frontendUrl),
	}
}

func (ph *ProxyHandler) Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlers.ProxyHandler.Handler is ready to work")
	userId, err := r.Cookie(utils.CookieKeyUserId)
	if err != nil {
		slog.Info(fmt.Sprintf("handlers.ProxyHandler: Redirecting to authorization, %v", err))
		http.Redirect(w, r, "/oauth/begin/", http.StatusSeeOther)
	}
	target := r.URL.Host

	slog.Info(fmt.Sprintf("handlers.ProxyHandler.Handler: user with id %s exists, wants to see %s", userId, target))
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
