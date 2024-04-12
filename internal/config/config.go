package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env          string          `yaml:"env" env-default:"local"`
	GinCfg       *HTTPConfig     `yaml:"server"`
	PostgresCfg  *PostgresConfig `yaml:"postgres"`
	RedisCfg     *RedisConfig    `yaml:"redis"`
	ServiceCfg   *ServiceConfig  `yaml:"service"`
	StartTimeout time.Duration   `yaml:"start_stop_timeout" env-default:"5s"`
}

type ServiceConfig struct {
	Timeout time.Duration `yaml:"request_timeout" env-default:"5s"`
}

type PostgresConfig struct {
	DSN     string        `yaml:"dsn" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type CacheConfig struct {
	LocalSize int           `yaml:"local_size" env-required:"true"`
	TTL       time.Duration `yaml:"ttl" env-default:"5m"`
}

type RedisConfig struct {
	Address string       `yaml:"address" env-required:"true"`
	Cache   *CacheConfig `yaml:"cache"`
}

type HTTPConfig struct {
	Address string `yaml:"address" env-default:":8080"`
}

func MustLoad() *Config {
	configPath := getConfigPath()
	cfg, err := Load(configPath)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}