package modules

import (
	"fmt"
	"gateway/pkg/config"
	"net/url"
)

const (
	BackendUrl  config.ConfigKey = "BACKEND_URL"
	FrontendUrl config.ConfigKey = "FRONTEND_URL"
)

type BaseConfig struct {
	backendUrl   *url.URL
	frontendUrl  *url.URL
	backendHost  string
	frontendHost string
}

func NewBaseConfig() *BaseConfig {
	backendUrl, err := url.Parse(BackendUrl.MustGet())
	if err != nil {
		panic(fmt.Sprintf("config.NewConfig: failed to parse BACKEND_URL: %v", err))
	}

	frontendUrl, err := url.Parse(FrontendUrl.MustGet())
	if err != nil {
		panic(fmt.Sprintf("config.NewConfig: failed to parse FRONTEND_URL: %v", err))
	}

	backendHost := BackendUrl.MustGet()

	frontendHost := FrontendUrl.MustGet()

	return &BaseConfig{
		backendUrl:   backendUrl,
		frontendUrl:  frontendUrl,
		backendHost:  backendHost,
		frontendHost: frontendHost,
	}
}

func (baseConfig *BaseConfig) GetBackendUrl() *url.URL {
	return baseConfig.backendUrl
}

func (baseConfig *BaseConfig) GetFrontendUrl() *url.URL {
	return baseConfig.frontendUrl
}

func (baseConfig *BaseConfig) GetBackendHost() string {
	return baseConfig.backendHost
}

func (baseConfig *BaseConfig) GetFrontendHost() string {
	return baseConfig.frontendHost
}
