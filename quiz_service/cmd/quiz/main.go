package main

import (
	"context"

	"quizzes/internal/app"
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
