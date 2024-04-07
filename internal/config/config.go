package config

import "time"

type Config struct {
	GinConfig *HTTPConfig
}

type HTTPConfig struct {
	Address string
	Timeout time.Duration
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load() (*Config, error) {
	_ = "config.load"
	panic("implement me")
	return &Config{}, nil
}
