package config

import (
	"awesomeProject2/pkg/postgres"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
	Host     string          `yaml:"HOST" env:"HOST" env-default:"127.0.0.1"`
	GRPCPort int             `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
	HTTPPort int             `yaml:"HTTP_PORT" env:"HTTP_PORT" env-default:"8080"`
}

func getProjectRoot() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
}

func New() (*Config, error) {
	var cfg Config
	root := getProjectRoot()
	configPath := filepath.Join(root, "config", "config.yaml")
	if err := cleanenv.ReadConfig(configPath, &cfg); err == nil {
		return &cfg, nil
	} else if err2 := cleanenv.ReadEnv(&cfg); err2 == nil {
		return &cfg, nil
	} else {
		return nil, err
	}

}
