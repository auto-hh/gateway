package modules

import (
	"gateway/pkg/config"
	"strconv"
	"time"
)

const (
	SessionIdExpirationTime config.ConfigKey = "SESSION_ID_EXPIRATION_TIME"
	StateExpirationTime     config.ConfigKey = "STATE_EXPIRATION_TIME"

	DefaultSessionIdTimeout int = 20
	DefaultStateTimeout     int = 20
)

type TimeoutConfig struct {
	sessionIdExpirationTime time.Duration
	stateExpirationTime     time.Duration
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

func (timeoutConfig *TimeoutConfig) GetSessionIdExpirationTime() time.Duration {
	return timeoutConfig.sessionIdExpirationTime
}

func (timeoutConfig *TimeoutConfig) GetStateExpirationTime() time.Duration {
	return timeoutConfig.stateExpirationTime
}
