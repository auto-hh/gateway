package config

import (
	"fmt"
	"gateway/pkg/config"
	"net/url"
)

const BACKEND_URL config.ConfigKey = "BACKEND_URL"
const FRONTEND_URL config.ConfigKey = "FRONTEND_URL"

type Config struct {
	backendUrl   *url.URL
	frontendUrl  *url.URL
	backendHost  string
	frontendHost string
}

func NewConfig() *Config {
	backendUrl, err := url.Parse(BACKEND_URL.MustGet())
	if err != nil {
		panic(fmt.Sprintf("config.NewConfig: failed to parse BACKEND_URL: %v", err))
	}
	frontendUrl, err := url.Parse(FRONTEND_URL.MustGet())
	if err != nil {
		panic(fmt.Sprintf("config.NewConfig: failed to parse FRONTEND_URL: %v", err))
	}
	backendHost := BACKEND_URL.MustGet()
	frontendHost := FRONTEND_URL.MustGet()

	return &Config{
		backendUrl:   backendUrl,
		frontendUrl:  frontendUrl,
		backendHost:  backendHost,
		frontendHost: frontendHost,
	}
}

func (c *Config) BackendUrl() *url.URL {
	return c.backendUrl
}

func (c *Config) FrontendUrl() *url.URL {
	return c.frontendUrl
}

func (c *Config) FrontendHost() string {
	return c.frontendHost
}

func (c *Config) BackendHost() string {
	return c.backendHost
}
