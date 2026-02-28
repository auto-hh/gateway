package reverse_proxy

import (
	"fmt"
	"gateway/internal/handlers"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	backendHost   string
	frontendHost  string
	backendProxy  *httputil.ReverseProxy
	frontendProxy *httputil.ReverseProxy
}

func NewReverseProxy(backendUrl, frontendUrl *url.URL, backendHost, frontendHost string) *ReverseProxy {
	return &ReverseProxy{
		backendHost:   backendHost,
		frontendHost:  frontendHost,
		backendProxy:  httputil.NewSingleHostReverseProxy(backendUrl),
		frontendProxy: httputil.NewSingleHostReverseProxy(frontendUrl),
	}
}

func (proxy *ReverseProxy) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ProxyHandler)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		return fmt.Errorf("reverse_proxy.Start: %v", err)
	}
	return nil
}

func (proxy *ReverseProxy) Stop() {}
