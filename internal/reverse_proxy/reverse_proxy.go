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

func (proxy *ReverseProxy) Start() {
	http.HandleFunc("/", handlers.ProxyHandler)

	err := http.ListenAndServe(proxy.frontendHost, nil)
	if err != nil {
		panic(fmt.Errorf("reverse_proxy.Start: %v", err))
	}
}

func (proxy *ReverseProxy) Stop() {}
