package config

import (
	"fmt"
	"gateway/pkg/config"
	"net/url"
	"strconv"
	"time"
)

const (
	BackendUrl       config.ConfigKey = "BACKEND_URL"
	FrontendUrl      config.ConfigKey = "FRONTEND_URL"
	RepoUser         config.ConfigKey = "REPO_USER"
	RepoPassword     config.ConfigKey = "REPO_PASSWORD"
	RepoAddr         config.ConfigKey = "REPO_ADDR"
	RepoMaxRetries   config.ConfigKey = "REPO_MAX_RETRIES"
	RepoDialTimeout  config.ConfigKey = "REPO_DIAL_TIMEOUT"
	RepoReadTimeout  config.ConfigKey = "REPO_READ_TIMEOUT"
	RepoWriteTimeout config.ConfigKey = "REPO_WRITE_TIMEOUT"

	DefaultMaxRetries   int = 3
	DefaultDialTimeout  int = 20
	DefaultReadTimeout  int = 20
	DefaultWriteTimeout int = 20
)

type Config struct {
	backendUrl       *url.URL
	frontendUrl      *url.URL
	backendHost      string
	frontendHost     string
	repoUser         string
	repoPassword     string
	repoAddr         string
	repoMaxRetries   int
	repoDialTimeout  time.Duration
	repoReadTimeout  time.Duration
	repoWriteTimeout time.Duration
}

func NewConfig() *Config {
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

	repoUser := RepoUser.MustGet()

	repoPassword := RepoPassword.MustGet()

	repoAddr := RepoAddr.MustGet()

	repoMaxRetries, err := strconv.Atoi(RepoMaxRetries.MustGet())
	if err != nil {
		repoMaxRetries = DefaultMaxRetries
	}

	repoDialTimeout, err := strconv.Atoi(RepoDialTimeout.MustGet())
	if err != nil {
		repoDialTimeout = DefaultDialTimeout
	}

	repoReadTimeout, err := strconv.Atoi(RepoReadTimeout.MustGet())
	if err != nil {
		repoReadTimeout = DefaultReadTimeout
	}

	repoWriteTimeout, err := strconv.Atoi(RepoWriteTimeout.MustGet())
	if err != nil {
		repoWriteTimeout = DefaultWriteTimeout
	}

	return &Config{
		backendUrl:       backendUrl,
		frontendUrl:      frontendUrl,
		backendHost:      backendHost,
		frontendHost:     frontendHost,
		repoUser:         repoUser,
		repoPassword:     repoPassword,
		repoAddr:         repoAddr,
		repoMaxRetries:   repoMaxRetries,
		repoDialTimeout:  time.Duration(repoDialTimeout) * time.Second,
		repoReadTimeout:  time.Duration(repoReadTimeout) * time.Second,
		repoWriteTimeout: time.Duration(repoWriteTimeout) * time.Second,
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

func (c *Config) RepoUser() string {
	return c.repoUser
}

func (c *Config) RepoPassword() string {
	return c.repoPassword
}

func (c *Config) RepoMaxRetries() int {
	return c.repoMaxRetries
}

func (c *Config) RepoDialTimeout() time.Duration {
	return c.repoDialTimeout
}

func (c *Config) RepoReadTimeout() time.Duration {
	return c.repoReadTimeout
}

func (c *Config) RepoWriteTimeout() time.Duration {
	return c.repoWriteTimeout
}
