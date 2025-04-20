package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort string      `yaml:"grpcPort" json:"grpcPort"`
	DB       DBConfig    `yaml:"db" json:"db"`
	Redis    RedisConfig `yaml:"redis" json:"redis"`
	JWT      JWTConfig   `yaml:"jwt" json:"jwt"`
}

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret" json:"secret"`
	TTL    time.Duration `yaml:"token_ttl" json:"token_ttl"`
}

type DBConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password"  json:"password"`
	Name     string `yaml:"name"  json:"name"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}
