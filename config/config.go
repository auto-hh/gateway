package config

import "gateway/config/modules"

type Config struct {
	base    *modules.BaseConfig
	timeout *modules.TimeoutConfig
	hh      *modules.HHConfig
	repo    *modules.RepoConfig
}

func NewConfig() *Config {

	return &Config{
		base:    modules.NewBaseConfig(),
		timeout: modules.NewTimeoutConfig(),
		hh:      modules.NewHHConfig(),
		repo:    modules.NewRepoConfig(),
	}
}

func (c *Config) GetBaseConfig() *modules.BaseConfig {
	return c.base
}

func (c *Config) GetTimeoutConfig() *modules.TimeoutConfig {
	return c.timeout
}

func (c *Config) GetRepoConfig() *modules.RepoConfig {
	return c.repo
}

func (c *Config) GetHHConfig() *modules.HHConfig {
	return c.hh
}
