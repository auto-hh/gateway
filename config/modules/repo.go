package modules

import (
	"gateway/pkg/config"
	"strconv"
	"time"
)

const (
	RepoUser         config.ConfigKey = "REPO_USER"
	RepoPassword     config.ConfigKey = "REPO_PASSWORD"
	RepoAddr         config.ConfigKey = "REPO_ADDR"
	RepoDbName       config.ConfigKey = "REPO_DB_NAME"
	RepoMaxRetries   config.ConfigKey = "REPO_MAX_RETRIES"
	RepoDialTimeout  config.ConfigKey = "REPO_DIAL_TIMEOUT"
	RepoReadTimeout  config.ConfigKey = "REPO_READ_TIMEOUT"
	RepoWriteTimeout config.ConfigKey = "REPO_WRITE_TIMEOUT"

	DefaultDbName       string = "1"
	DefaultMaxRetries   int    = 3
	DefaultDialTimeout  int    = 20
	DefaultReadTimeout  int    = 20
	DefaultWriteTimeout int    = 20
)

type RepoConfig struct {
	user         string
	password     string
	addr         string
	dbName       int
	maxRetries   int
	dialTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewRepoConfig() *RepoConfig {
	user := RepoUser.MustGet()

	password := RepoPassword.MustGet()

	addr := RepoAddr.MustGet()

	var dbName int
	dbName, err := strconv.Atoi(RepoDbName.Get(DefaultDbName))
	if err != nil {
		dbName = 1
	}

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
		dbName:       dbName,
		maxRetries:   maxRetries,
		dialTimeout:  time.Duration(dialTimeout) * time.Second,
		readTimeout:  time.Duration(readTimeout) * time.Second,
		writeTimeout: time.Duration(writeTimeout) * time.Second,
	}
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

func (repoConfig *RepoConfig) GetDbName() int {
	return repoConfig.dbName
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
