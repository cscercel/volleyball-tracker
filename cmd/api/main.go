package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/cscercel/volleyball-tracker/internal/config"
	"github.com/cscercel/volleyball-tracker/internal/db"
	"github.com/cscercel/volleyball-tracker/internal/handler"
	"github.com/cscercel/volleyball-tracker/internal/service"

	_ "github.com/cscercel/volleyball-tracker/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
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

	// Connect to Database
	conn, err := sql.Open("libsql", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open to database: %v", err)
	}
	defer conn.Close()

	// Verify connection
	if err := conn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Instantiate API
	queries := db.New(conn)

	// Player handlers
	playerService := service.NewPlayerService(queries)
	playerHandler := handler.NewPlayerHandler(playerService)

	// Match handlers
	matchService := service.NewMatchService(queries)
	matchHandler := handler.NewMatchHandler(matchService)

	// User handlers
	userService := service.NewUserService(queries, cfg.RegistrationCode, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService)

	// Routers
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Local middleware
	authMiddleware := handler.AuthenticateMiddleware(cfg.JWTSecret)

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		playerHandler.RegisterRoutes(r, authMiddleware)
		matchHandler.RegisterRoutes(r, authMiddleware)
		userHandler.RegisterRoutes(r, authMiddleware)
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
