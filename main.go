package main

import (
	"context"
	"log"

	"github.com/foxyblue/noiibat/config"
	noiibat "github.com/foxyblue/noiibat/pkg"
)

func main() {
	// cmd.Execute()
	cfg := config.New()
	ctx := context.Background()

	bat, err := noiibat.Spawn(ctx, cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	bat.ListenAndServe()
}
