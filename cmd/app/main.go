package main

import (
	"context"
	"log"

	"lab/internal/config"
	"lab/internal/infra"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return
	}

	app, err := infra.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer app.Stop(ctx)

	app.Run(ctx)
}
