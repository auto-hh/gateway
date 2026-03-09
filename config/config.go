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

type BaseConfig struct {
	backendUrl   *url.URL
	frontendUrl  *url.URL
	backendHost  string
	frontendHost string
}

type RepoConfig struct {
	user         string
	password     string
	addr         string
	maxRetries   int
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type TimeoutConfig struct {
	sessionIdExpirationTime time.Duration
	stateExpirationTime     time.Duration
}

type HHConfig struct {
	clientId     string
	clientSecret string
	appName      string
	appVersion   string
	redirectUri  string
	devContact   string
	rawUrl       string
}
type Config struct {
	BaseConfig
	TimeoutConfig
	HHConfig
	RepoConfig
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

func NewRepoConfig() *RepoConfig {
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
	return &RepoConfig{
		user:         user,
		password:     password,
		addr:         addr,
		maxRetries:   maxRetries,
		dialTimeout:  time.Duration(dialTimeout) * time.Second,
		readTimeout:  time.Duration(readTimeout) * time.Second,
		writeTimeout: time.Duration(writeTimeout) * time.Second,
	}
}

func NewTimeoutConfig() *TimeoutConfig {

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
	return &TimeoutConfig{
		sessionIdExpirationTime: sessionIdExpirationTime,
		stateExpirationTime:     stateExpirationTime,
	}
}

func NewHHConfig() *HHConfig {
	clientId := HHClientId.MustGet()
	clientSecret := HHClientSecret.MustGet()
	redirectUri := HHRedirectUri.MustGet()
	devContact := HHDevContact.Get(DefaultDevContact)
	rawUrl := HHRawUrl.MustGet()
	appName := HHAppName.Get(DefaultAppName)
	appVersion := HHAppVersion.Get(DefaultAppVersion)

	return &HHConfig{
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
		BaseConfig:    *NewBaseConfig(),
		TimeoutConfig: *NewTimeoutConfig(),
		RepoConfig:    *NewRepoConfig(),
		HHConfig:      *NewHHConfig(),
	}
}

func (c *Config) GetBaseConfig() BaseConfig {
	return c.BaseConfig
}

func (c *Config) GetTimeoutConfig() TimeoutConfig {
	return c.TimeoutConfig
}

func (c *Config) GetRepoConfig() RepoConfig {
	return c.RepoConfig
}

func (c *Config) GetHHConfig() HHConfig {
	return c.HHConfig
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

func (repoConfig *RepoConfig) GetUser() string {
	return repoConfig.user
}

func (repoConfig *RepoConfig) GetPassword() string {
	return repoConfig.password
}

func (repoConfig *RepoConfig) GetAddr() string {
	return repoConfig.addr
}

func (repoConfig *RepoConfig) GetMaxRetries() int {
	return repoConfig.maxRetries
}

func (repoConfig *RepoConfig) GetDialTimeout() time.Duration {
	return repoConfig.dialTimeout
}

func (repoConfig *RepoConfig) GetReadTimeout() time.Duration {
	return repoConfig.readTimeout
}

func (repoConfig *RepoConfig) GetWriteTimeout() time.Duration {
	return repoConfig.writeTimeout
}

func (timeoutConfig *TimeoutConfig) GetSessionIdExpirationTime() time.Duration {
	return timeoutConfig.sessionIdExpirationTime
}

func (timeoutConfig *TimeoutConfig) GetStateExpirationTime() time.Duration {
	return timeoutConfig.stateExpirationTime
}

func (hhConfig *HHConfig) GetAppName() string {
	return hhConfig.appName
}

func (hhConfig *HHConfig) GetAppVersion() string {
	return hhConfig.appVersion
}
func (hhConfig *HHConfig) GetRedirectUri() string {
	return hhConfig.redirectUri
}

func (hhConfig *HHConfig) GetDevContact() string {
	return hhConfig.devContact
}

func (hhConfig *HHConfig) GetRawUrl() string {
	return hhConfig.rawUrl
}

func (hhConfig *HHConfig) GetClientId() string {
	return hhConfig.clientId
}

func (hhConfig *HHConfig) GetClientSecret() string {
	return hhConfig.clientSecret
}
