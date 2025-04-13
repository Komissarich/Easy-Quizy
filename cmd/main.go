package main

import (
	"context"
	"eazy-quizy-auth/internal/application"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := application.New("config.yml")
	app.Run(ctx)
}
