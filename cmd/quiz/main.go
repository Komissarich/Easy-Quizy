package main

import (
	"awesomeProject2/internal/app"
	"context"
)

const (
	DEFAULT_CONFIG_PATH = "config.yml"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := app.New(DEFAULT_CONFIG_PATH)
	application.Run(ctx)
}
