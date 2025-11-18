package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Daumas-hugo/ReSell/internal/companies"
	"github.com/Daumas-hugo/ReSell/internal/config"
	"github.com/Daumas-hugo/ReSell/internal/database"
	"github.com/Daumas-hugo/ReSell/internal/orders"
	"github.com/Daumas-hugo/ReSell/internal/products"
	"github.com/Daumas-hugo/ReSell/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Initialize services and handlers
	companiesService := companies.NewService(db.DB)
	companiesHandler := companies.NewHandler(companiesService, log)

	productsService := products.NewService(db.DB)
	productsHandler := products.NewHandler(productsService, log)

	ordersService := orders.NewService(db.DB)
	ordersHandler := orders.NewHandler(ordersService, log)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Companies
		r.Route("/companies", companiesHandler.RegisterRoutes)

		// Products
		r.Route("/products", productsHandler.RegisterRoutes)

		// Orders
		r.Route("/orders", ordersHandler.RegisterRoutes)

		// Customer portal
		r.Route("/me", func(r chi.Router) {
			// Get current user info (requires auth middleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// TODO: Implement after Keycloak auth
				w.WriteHeader(http.StatusNotImplemented)
			})

			// Get user's orders
			r.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
				// TODO: Implement after Keycloak auth
				w.WriteHeader(http.StatusNotImplemented)
			})

			// Get user's companies
			r.Get("/companies", func(w http.ResponseWriter, r *http.Request) {
				// TODO: Implement after Keycloak auth
				w.WriteHeader(http.StatusNotImplemented)
			})
		})
	})

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		log.Info().Msgf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
