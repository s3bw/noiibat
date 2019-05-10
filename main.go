package main

import (
	"context"
	"log"

	"github.com/foxyblue/noiibat/config"
	noiibat "github.com/foxyblue/noiibat/pkg"
)

func main() {
	cfg := config.New()
	ctx := context.Background()

	app, err := noiibat.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	app.Start()
}
