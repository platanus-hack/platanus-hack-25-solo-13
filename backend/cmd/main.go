package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/handlers"
	authmiddleware "github.com/platanus-hack-25/lumera_app/internal/middleware"
)

// @title Lumera API
// @version 1.0
// @description API del proyecto Lumera - Platanus Hack 25
// @description Sistema de gestión con autenticación JWT

// @contact.name Equipo Lumera
// @contact.email info@lumera.com

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Connect to database with retry logic
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if err := db.Connect(); err != nil {
			log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(time.Second * 2)
				continue
			}
			log.Fatal("Failed to connect to database after retries")
		}
		break
	}

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/api/health", handlers.HealthCheck)

	// Public auth routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", handlers.Register)
		r.Post("/login", handlers.Login)
	})

	// Protected user routes
	r.Route("/api/users", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/me", handlers.GetMe)
		r.Put("/me", handlers.UpdateMe)
		r.Post("/change-password", handlers.ChangePassword)
		r.Delete("/me", handlers.DeleteMe)
	})

	// Protected profile routes
	r.Route("/api/profiles", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/export", handlers.ExportProfiles)          // Export for ML
		r.Get("/{user_id}", handlers.GetProfile)           // Get profile by user_id
		r.Post("/", handlers.CreateProfile)                 // Create new profile
		r.Patch("/{user_id}", handlers.UpdateProfile)      // Update profile
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("✓ Server starting on http://localhost%s", addr)
	log.Printf("✓ Health check available at http://localhost%s/api/health", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
