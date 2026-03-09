package server

import (
	"fmt"
	"gateway/config/modules"
	"gateway/internal/handlers"
	"gateway/internal/service/oauth"
	"net/http"
)

type Server struct {
	proxyHandler *handlers.ProxyHandler
	oauthHandler *handlers.OAuthHandler
}

func NewServer(baseConfig *modules.BaseConfig, service *oauth.Service) *Server {
	return &Server{
		proxyHandler: handlers.NewProxyHandler(baseConfig),
		oauthHandler: handlers.NewOAuthHandler(service),
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.proxyHandler.Handler)
	mux.HandleFunc("/oauth/begin/", s.oauthHandler.Begin)
	mux.HandleFunc("/oauth/complete/", s.oauthHandler.Complete)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		return fmt.Errorf("server.Start: %v", err)
	}
	return nil
}

//func (proxy *Server) Stop() {} -> gracefull shutdown TBA
