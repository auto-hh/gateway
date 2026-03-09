package config

import "gateway/config/modules"

type Config struct {
	modules.BaseConfig
	modules.TimeoutConfig
	modules.HHConfig
	modules.RepoConfig
}

func NewConfig() *Config {
	return &Config{
		BaseConfig:    *modules.NewBaseConfig(),
		TimeoutConfig: *modules.NewTimeoutConfig(),
		HHConfig:      *modules.NewHHConfig(),
		RepoConfig:    *modules.NewRepoConfig(),
	}
}

func (c *Config) GetBaseConfig() modules.BaseConfig {
	return c.BaseConfig
}

func (c *Config) GetTimeoutConfig() modules.TimeoutConfig {
	return c.TimeoutConfig
}

func (c *Config) GetRepoConfig() modules.RepoConfig {
	return c.RepoConfig
}

func (c *Config) GetHHConfig() modules.HHConfig {
	return c.HHConfig
}
