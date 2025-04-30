package main

import (
	"context"
	"log"
	"os"

	"github.com/extizout/pokemon-api-server/config"
	"github.com/extizout/pokemon-api-server/server"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Please provide a path to the config file")
		}
		return os.Args[1]
	}())

	server.Start(ctx, &cfg)
}
