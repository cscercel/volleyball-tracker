package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"github.com/cscercel/volleyball-tracker/internal/config"
	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/cscercel/volleyball-tracker/internal/handler"
	"github.com/cscercel/volleyball-tracker/internal/repository"
	"github.com/cscercel/volleyball-tracker/internal/service"

	_ "github.com/cscercel/volleyball-tracker/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Volleyball Tracker API
// @version         2.0.0
// @description     API for managing volleyball games played with friends.

// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by your JWT token
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
	queries := db.New(pool)

	// Player handlers
	playerService := service.NewPlayerService(queries)
	playerHandler := handler.NewPlayerHandler(playerService)

	// Match handlers
	matchService := service.NewMatchService(queries)
	matchHandler := handler.NewMatchHandler(matchService)

	// User handlers
	userService := service.NewUserService(queries, cfg.RegistrationCode, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService)

	// Page handler
	secureCookies := cfg.Env == "production"
	pageHandler := handler.NewPageHandler(
		userService, 
		playerService,
		matchService,
		cfg.JWTSecret, 
		secureCookies,
		)

	// Routers
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Local middleware
	authMiddleware := handler.AuthenticateMiddleware(cfg.JWTSecret)

	// JSON Routes
	r.Route("/api/v1", func(r chi.Router) {
		playerHandler.RegisterRoutes(r, authMiddleware)
		matchHandler.RegisterRoutes(r, authMiddleware)
		userHandler.RegisterRoutes(r, authMiddleware)
	})

	// HTMX Routes
	pageHandler.RegisterRoutes(r)

	// Small `Mandatory` test route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	// Docs
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Start server
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	log.Printf("server running on port: %v", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
