package server

import (
	"fmt"
	"gateway/internal/handlers"
	"net/http"
	"net/url"
)

type Server struct {
	proxyHandler *handlers.ProxyHandler
	oauthHandler *handlers.OAuthHandler
}

func NewServer(backendUrl, frontendUrl *url.URL, backendHost, frontendHost string) *Server {
	return &Server{
		proxyHandler: handlers.NewProxyHandler(backendUrl, frontendUrl, backendHost, frontendHost),
	}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.proxyHandler.Handler)
	mux.HandleFunc("/oauth/begin/", s.oauthHandler.Begin)
	mux.HandleFunc("/oauth/complete/", s.oauthHandler.Complete)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		return fmt.Errorf("server.Start: %v", err)
	}
	return nil
}

func (proxy *Server) Stop() {}
