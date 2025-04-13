package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort string    `yaml:"grpcPort" json:"grpcPort"`
	DB       DBConfig  `yaml:"db" json:"db"`
	JWT      JWTConfig `yaml:"jwt" json:"jwt"`
}

type JWTConfig struct {
	Secret string `yaml:"secret" json:"secret"`
	TTL    string `yaml:"token_ttl" json:"token_ttl"`
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
