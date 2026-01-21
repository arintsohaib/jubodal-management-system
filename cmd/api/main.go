package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjdms/api/config"
	"github.com/bjdms/api/internal/activity"
	"github.com/bjdms/api/internal/analytics"
	"github.com/bjdms/api/internal/auth"
	"github.com/bjdms/api/internal/notification"
	"github.com/bjdms/api/internal/committee"
	"github.com/bjdms/api/internal/complaint"
	"github.com/bjdms/api/internal/finance"
	"github.com/bjdms/api/internal/join"
	"github.com/bjdms/api/internal/search"
	"github.com/bjdms/api/internal/database"
	internalMiddleware "github.com/bjdms/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisMgr, err := auth.NewRedisManager(cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize JWT manager
	jwtMgr := auth.NewJWTManager(
		cfg.JWTAccessSecret,
		cfg.JWTRefreshSecret,
		cfg.JWTAccessExpiry,
		cfg.JWTRefreshExpiry,
	)

	// Initialize repositories and services
	authRepo := auth.NewRepository(db.Pool)
	authService := auth.NewService(authRepo, redisMgr, jwtMgr, cfg)
	authHandler := auth.NewHandler(authService, redisMgr, jwtMgr)

	committeeRepo := committee.NewRepository(db.Pool)
	committeeService := committee.NewService(committeeRepo)
	committeeHandler := committee.NewHandler(committeeService)

	notificationService := notification.NewService(db.Pool, redisMgr.Client())
	notificationHandler := notification.NewHandler(notificationService)

	activityRepo := activity.NewRepository(db.Pool)
	activityService := activity.NewService(activityRepo, notificationService)
	activityHandler := activity.NewHandler(activityService)

	complaintRepo := complaint.NewRepository(db.Pool)
	complaintService := complaint.NewService(complaintRepo, notificationService)
	complaintHandler := complaint.NewHandler(complaintService)

	searchClient, err := search.NewClient(cfg.OpenSearchURL)
	if err == nil {
		searchClient.InitIndices(context.Background())
	}
	searchService := search.NewService(searchClient)
	searchHandler := search.NewHandler(searchService)

	financeRepo := finance.NewRepository(db.Pool)
	financeService := finance.NewService(financeRepo)
	financeHandler := finance.NewHandler(financeService)


	analyticsClient := analytics.NewClient()
	analyticsHandler := analytics.NewHandler(analyticsClient)

	joinRepo := join.NewRepository(db.Pool)
	joinService := join.NewService(joinRepo, authRepo, notificationService)
	joinHandler := join.NewHandler(joinService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Rate limiters
	loginLimiter := internalMiddleware.NewRateLimiter(cfg.RateLimitLogin, cfg.RateLimitWindow)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Health(ctx); err != nil {
			http.Error(w, `{"status":"unhealthy","database":"disconnected"}`, http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","database":"connected"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public Auth routes
		r.Group(func(r chi.Router) {
			r.Use(loginLimiter.Limit)
			r.Mount("/auth", authHandler.Routes())
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(internalMiddleware.AuthMiddleware(jwtMgr, redisMgr))
			
			// Committees & Jurisdictions
			r.Mount("/org", committeeHandler.Routes())

			// Activities & Tasks
			r.Mount("/activities", activityHandler.Routes())

			// Complaints Management (Internal)
			r.Group(func(r chi.Router) {
				r.Use(internalMiddleware.ABACJurisdictionMiddleware(committeeService, authRepo))
				r.Mount("/complaints", complaintHandler.Routes())
			})

			// Search (Global)
			r.Mount("/search", searchHandler.Routes())

			// Finance (Internal)
			r.Group(func(r chi.Router) {
				r.Use(internalMiddleware.ABACJurisdictionMiddleware(committeeService, authRepo))
				r.Mount("/finance", financeHandler.Routes())
				r.Mount("/join-requests", joinHandler.Routes())
				r.Mount("/analytics", analyticsHandler.Routes())
				r.Mount("/notifications", notificationHandler.Routes())
			})
			
			// r.Mount("/users", userHandler.Routes())
		})

		// Public Anonymous Complaints
		r.Mount("/public/complaints", complaintHandler.Routes())
		r.Mount("/public/join", joinHandler.PublicRoutes())
	})

	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("✓ Server starting on port %s", cfg.Port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)

	case <-shutdown:
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Fatalf("Could not stop server: %v", err)
			}
		}
	}

	log.Println("✓ Server stopped")
}
