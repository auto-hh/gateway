package modules

import (
	"fmt"
	"gateway/pkg/config"
	"net/url"
	"strconv"
)

const (
	BackendUrl  config.ConfigKey = "BACKEND_URL"
	FrontendUrl config.ConfigKey = "FRONTEND_URL"
	BackendHost config.ConfigKey = "BACKEND_HOST"
	FrontendHost config.ConfigKey = "FRONTEND_HOST"
	ServerPort  config.ConfigKey = "SERVER_PORT"

	DefaultServerPort int = 8080
)

type BaseConfig struct {
	serverPort   int
	backendUrl   *url.URL
	frontendUrl  *url.URL
	backendHost  string
	frontendHost string
}

func NewBaseConfig() *BaseConfig {

	serverPortString := ServerPort.MustGet()

	serverPort, err := strconv.Atoi(serverPortString)
	if err != nil {
		serverPort = DefaultServerPort
	}

	backendUrl, err := url.Parse(BackendUrl.MustGet())
	if err != nil {
		panic(fmt.Sprintf("config.NewConfig: failed to parse BACKEND_URL: %v", err))
	}

	frontendUrl, err := url.Parse(FrontendUrl.MustGet())
	if err != nil {
		panic(fmt.Sprintf("config.NewConfig: failed to parse FRONTEND_URL: %v", err))
	}

	backendHost := BackendHost.MustGet()

	frontendHost := FrontendHost.MustGet()

	return &BaseConfig{
		serverPort:   serverPort,
		backendUrl:   backendUrl,
		frontendUrl:  frontendUrl,
		backendHost:  backendHost,
		frontendHost: frontendHost,
	}
}

func (baseConfig *BaseConfig) GetServerPort() int {
	return baseConfig.serverPort
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
