package server

import (
	"fmt"
	"gateway/config/modules"
	"gateway/internal/handlers"
	"gateway/internal/service/oauth"
	"gateway/internal/service/reverse_proxy"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	proxyHandler *handlers.ProxyHandler
	oauthHandler *handlers.OAuthHandler
}

func NewServer(baseConfig *modules.BaseConfig, serviceOauth *oauth.Service, serviceReverseProxy *reverse_proxy.Service) *Server {
	return &Server{
		proxyHandler: handlers.NewProxyHandler(baseConfig, serviceReverseProxy),
		oauthHandler: handlers.NewOAuthHandler(serviceOauth),
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.proxyHandler.Handler)
	mux.HandleFunc("/oauth/begin/", s.oauthHandler.Begin)
	mux.HandleFunc("/oauth/complete/", s.oauthHandler.Complete)
	slog.Info("start listening on 0.0.0.0:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), loggingMiddleware(mux))
	if err != nil {
		return fmt.Errorf("server.Start: %v", err)
	}
	return nil
}

func loggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("response", slog.String("ip", r.RemoteAddr), slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Duration("duration", time.Since(start)))
	}
}

//func (proxy *Server) Stop() {} -> gracefull shutdown TBA
