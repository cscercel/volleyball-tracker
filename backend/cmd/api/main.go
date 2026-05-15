package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"


	"github.com/cscercel/volleyball-tracker/internal/api"
	"github.com/cscercel/volleyball-tracker/internal/config"
	"github.com/cscercel/volleyball-tracker/internal/repository"

)

func main() {

	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	// Connect to Database
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Instantiate API
	a := api.New(pool, cfg)
	router := a.Routes()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Behold DnD server listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
