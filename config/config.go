package config

import (
	"fmt"
	"gateway/pkg/config"
	"net/url"
	"strconv"
	"time"
)

const (
	BackendUrl  config.ConfigKey = "BACKEND_URL"
	FrontendUrl config.ConfigKey = "FRONTEND_URL"

	RepoUser         config.ConfigKey = "REPO_USER"
	RepoPassword     config.ConfigKey = "REPO_PASSWORD"
	RepoAddr         config.ConfigKey = "REPO_ADDR"
	RepoMaxRetries   config.ConfigKey = "REPO_MAX_RETRIES"
	RepoDialTimeout  config.ConfigKey = "REPO_DIAL_TIMEOUT"
	RepoReadTimeout  config.ConfigKey = "REPO_READ_TIMEOUT"
	RepoWriteTimeout config.ConfigKey = "REPO_WRITE_TIMEOUT"

	HHClientId     config.ConfigKey = "CLIENT_ID"
	HHClientSecret config.ConfigKey = "CLIENT_SECRET"
	HHAppName      config.ConfigKey = "APP_NAME"
	HHAppVersion   config.ConfigKey = "APP_VERSION"
	HHRedirectUri  config.ConfigKey = "REDIRECT_URI"
	HHDevContact   config.ConfigKey = "DEV_CONTACT"
	HHRawUrl       config.ConfigKey = "RAW_URL"

	SessionIdExpirationTime config.ConfigKey = "SESSION_ID_EXPIRATION_TIME"
	StateExpirationTime     config.ConfigKey = "STATE_EXPIRATION_TIME"

	DefaultMaxRetries   int = 3
	DefaultDialTimeout  int = 20
	DefaultReadTimeout  int = 20
	DefaultWriteTimeout int = 20

	DefaultSessionIdTimeout int = 20
	DefaultStateTimeout     int = 20

	DefaultAppName    string = "MyAPP"
	DefaultDevContact string = "dev@mail.ru"
	DefaultAppVersion string = "1.0.0"
)

type baseConfig struct {
	backendUrl   *url.URL
	frontendUrl  *url.URL
	backendHost  string
	frontendHost string
}

type repoConfig struct {
	user         string
	password     string
	addr         string
	maxRetries   int
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type timeoutConfig struct {
	sessionIdExpirationTime time.Duration
	stateExpirationTime     time.Duration
}

type hhConfig struct {
	clientId     string
	clientSecret string
	appName      string
	appVersion   string
	redirectUri  string
	devContact   string
	rawUrl       string
}
type Config struct {
	baseConfig
	timeoutConfig
	hhConfig
	repoConfig
}

func NewBaseConfig() *baseConfig {
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

	return &baseConfig{
		backendUrl:   backendUrl,
		frontendUrl:  frontendUrl,
		backendHost:  backendHost,
		frontendHost: frontendHost,
	}
}

func NewRepoConfig() *repoConfig {
	user := RepoUser.MustGet()

	password := RepoPassword.MustGet()

	addr := RepoAddr.MustGet()

	maxRetries, err := strconv.Atoi(RepoMaxRetries.MustGet())
	if err != nil {
		maxRetries = DefaultMaxRetries
	}

	dialTimeout, err := strconv.Atoi(RepoDialTimeout.MustGet())
	if err != nil {
		dialTimeout = DefaultDialTimeout
	}

	readTimeout, err := strconv.Atoi(RepoReadTimeout.MustGet())
	if err != nil {
		readTimeout = DefaultReadTimeout
	}

	writeTimeout, err := strconv.Atoi(RepoWriteTimeout.MustGet())
	if err != nil {
		writeTimeout = DefaultWriteTimeout
	}
	return &repoConfig{
		user:         user,
		password:     password,
		addr:         addr,
		maxRetries:   maxRetries,
		dialTimeout:  time.Duration(dialTimeout) * time.Second,
		readTimeout:  time.Duration(readTimeout) * time.Second,
		writeTimeout: time.Duration(writeTimeout) * time.Second,
	}
}

func NewTimeoutConfig() *timeoutConfig {

	var sessionIdExpirationTime time.Duration
	var stateExpirationTime time.Duration

	sessionIdExpirationTimeInt, err := strconv.Atoi(SessionIdExpirationTime.Get(strconv.Itoa(DefaultSessionIdTimeout)))
	if err != nil {
		sessionIdExpirationTime = time.Duration(DefaultSessionIdTimeout) * time.Minute
	} else {
		sessionIdExpirationTime = time.Duration(sessionIdExpirationTimeInt) * time.Minute
	}
	stateExpirationTimeInt, err := strconv.Atoi(StateExpirationTime.Get(strconv.Itoa(DefaultStateTimeout)))
	if err != nil {
		stateExpirationTime = time.Duration(DefaultStateTimeout) * time.Minute
	} else {
		sessionIdExpirationTime = time.Duration(stateExpirationTimeInt) * time.Minute
	}
	return &timeoutConfig{
		sessionIdExpirationTime: sessionIdExpirationTime,
		stateExpirationTime:     stateExpirationTime,
	}
}

func NewHHConfig() *hhConfig {
	clientId := HHClientId.MustGet()
	clientSecret := HHClientSecret.MustGet()
	redirectUri := HHRedirectUri.MustGet()
	devContact := HHDevContact.Get(DefaultDevContact)
	rawUrl := HHRawUrl.MustGet()
	appName := HHAppName.Get(DefaultAppName)
	appVersion := HHAppVersion.Get(DefaultAppVersion)

	return &hhConfig{
		clientId:     clientId,
		clientSecret: clientSecret,
		appName:      appName,
		appVersion:   appVersion,
		redirectUri:  redirectUri,
		devContact:   devContact,
		rawUrl:       rawUrl,
	}
}

func NewConfig() *Config {

	return &Config{
		baseConfig:    *NewBaseConfig(),
		timeoutConfig: *NewTimeoutConfig(),
		repoConfig:    *NewRepoConfig(),
		hhConfig:      *NewHHConfig(),
	}
}

func (c *Config) GetBaseConfig() baseConfig {
	return c.baseConfig
}

func (c *Config) GetTimeoutConfig() timeoutConfig {
	return c.timeoutConfig
}

func (c *Config) GetRepoConfig() repoConfig {
	return c.repoConfig
}

func (c *Config) GetHHConfig() hhConfig {
	return c.hhConfig
}

func (baseConfig *baseConfig) GetBackendUrl() *url.URL {
	return baseConfig.backendUrl
}

func (baseConfig *baseConfig) GetFrontendUrl() *url.URL {
	return baseConfig.frontendUrl
}

func (baseConfig *baseConfig) GetBackendHost() string {
	return baseConfig.backendHost
}

func (baseConfig *baseConfig) GetFrontendHost() string {
	return baseConfig.frontendHost
}

func (repoConfig *repoConfig) GetUser() string {
	return repoConfig.user
}

func (repoConfig *repoConfig) GetPassword() string {
	return repoConfig.password
}

func (repoConfig *repoConfig) GetAddr() string {
	return repoConfig.addr
}

func (repoConfig *repoConfig) GetMaxRetries() int {
	return repoConfig.maxRetries
}

func (repoConfig *repoConfig) GetDialTimeout() time.Duration {
	return repoConfig.dialTimeout
}

func (repoConfig *repoConfig) GetReadTimeout() time.Duration {
	return repoConfig.readTimeout
}

func (repoConfig *repoConfig) GetWriteTimeout() time.Duration {
	return repoConfig.writeTimeout
}

func (timeoutConfig *timeoutConfig) GetSessionIdExpirationTime() time.Duration {
	return timeoutConfig.sessionIdExpirationTime
}

func (timeoutConfig *timeoutConfig) GetStateExpirationTime() time.Duration {
	return timeoutConfig.stateExpirationTime
}

func (hhConfig *hhConfig) GetAppName() string {
	return hhConfig.appName
}

func (hhConfig *hhConfig) GetAppVersion() string {
	return hhConfig.appVersion
}
func (hhConfig *hhConfig) GetRedirectUri() string {
	return hhConfig.redirectUri
}

func (hhConfig *hhConfig) GetDevContact() string {
	return hhConfig.devContact
}

func (hhConfig *hhConfig) GetRawUrl() string {
	return hhConfig.rawUrl
}

func (hhConfig *hhConfig) GetClientId() string {
	return hhConfig.clientId
}

func (hhConfig *hhConfig) GetClientSecret() string {
	return hhConfig.clientSecret
}
