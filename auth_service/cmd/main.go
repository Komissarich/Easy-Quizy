package main

import (
	"context"
	"eazy-quizy-auth/internal/application"
)

const (
	DEFAULT_CONFIG_PATH = "config.yml"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := application.New(DEFAULT_CONFIG_PATH)
	app.Run(ctx)
}
