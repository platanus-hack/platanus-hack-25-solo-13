package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/handlers"
	authmiddleware "github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/services"
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

	// Initialize OpenAI for content generation
	if err := services.InitOpenAI(); err != nil {
		log.Printf("⚠ Warning: OpenAI initialization failed: %v", err)
		log.Println("Learning plan generation will not be available")
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

	// Educational routes - Cursos (Courses)
	r.Route("/api/cursos", func(r chi.Router) {
		r.Get("/", handlers.GetCursos)              // Public: List all courses
		r.Get("/{id}", handlers.GetCurso)           // Public: Get course by ID

		// Protected write operations
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware)
			r.Post("/", handlers.CreateCurso)       // Create course
			r.Put("/{id}", handlers.UpdateCurso)    // Update course
		})
	})

	// Educational routes - Materias (Subjects)
	r.Route("/api/materias", func(r chi.Router) {
		r.Get("/", handlers.GetMaterias)            // Public: List all subjects
		r.Get("/{id}", handlers.GetMateria)         // Public: Get subject by ID

		// Protected write operations
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware)
			r.Post("/", handlers.CreateMateria)     // Create subject
			r.Put("/{id}", handlers.UpdateMateria)  // Update subject
			r.Get("/{materia_id}/oa-progress", handlers.GetOAProgressByMateria) // Get OA progress for user
		})
	})

	// Educational routes - Curso-Materia assignments
	r.Route("/api/curso-materias", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Post("/", handlers.AssignMateriaToCurso)  // Assign subject to course
	})

	// Educational routes - Bloom Levels (read-only)
	r.Get("/api/bloom-levels", handlers.GetBloomLevels)

	// Educational routes - Learning Objectives
	r.Route("/api/objetivos-aprendizaje", func(r chi.Router) {
		r.Get("/", handlers.GetObjetivosAprendizaje)        // Public: List all OAs
		r.Get("/{id}", handlers.GetObjetivoAprendizaje)     // Public: Get OA by ID

		// Protected write operations
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware)
			r.Post("/", handlers.CreateObjetivoAprendizaje)    // Create OA with all 6 Bloom levels
			r.Put("/{id}", handlers.UpdateObjetivoAprendizaje) // Update OA
		})
	})

	// Progress tracking routes (all protected)
	r.Route("/api/progress", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Post("/", handlers.RegisterProgress)                  // Register progress
		r.Get("/{user_id}", handlers.GetStudentProgress)        // Get student progress
		r.Get("/{user_id}/history", handlers.GetProgressHistory) // Get progress history
	})

	// Question Types catalog (public)
	r.Get("/api/question-types", handlers.GetQuestionTypes)

	// Questions Bank
	r.Route("/api/questions", func(r chi.Router) {
		r.Get("/", handlers.GetQuestions)               // Public: List questions with filters
		r.Get("/{id}", handlers.GetQuestion)            // Public: Get question (without validation_data)
		r.Post("/{id}/validate", handlers.ValidateAnswer) // Public: Validate answer

		// Protected write operations
		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware)
			r.Post("/", handlers.CreateQuestion)        // Create question
			r.Put("/{id}", handlers.UpdateQuestion)     // Update question
		})
	})

	// Diagnostic System (all protected)
	r.Route("/api/diagnostic-sessions", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/", handlers.GetDiagnosticSessions)             // List user's diagnostic sessions with filters
		r.Post("/", handlers.StartDiagnostic)                  // Start diagnostic session
		r.Get("/{id}", handlers.GetSessionProgress)            // Get session progress
		r.Get("/{id}/next-question", handlers.GetNextQuestion) // Get next adaptive question
		r.Post("/{id}/answer", handlers.SubmitAnswer)          // Submit answer
		r.Post("/{id}/complete", handlers.CompleteDiagnostic)  // Complete diagnostic
		r.Get("/{id}/results", handlers.GetDiagnosticResults)  // Get diagnostic results
	})

	// Practice System (all protected)
	r.Route("/api/practice-sessions", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/", handlers.GetPracticeSessions)                      // List user's practice sessions with filters
		r.Post("/", handlers.StartPractice)                           // Start practice session
		r.Get("/{id}/next-question", handlers.GetPracticeNextQuestion) // Get next adaptive question
		r.Post("/{id}/answer", handlers.SubmitPracticeAnswer)          // Submit answer
		r.Post("/{id}/complete", handlers.CompletePracticeSession)     // Complete practice session
	})

	// Gamification System (all protected)
	r.Route("/api/gamification", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/stats", handlers.GetGamificationStats)     // Get user gamification stats
		r.Get("/leaderboard", handlers.GetLeaderboard)     // Get leaderboard
	})

	// Customization System (all protected)
	r.Route("/api/customization", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/catalog", handlers.GetCustomizationCatalog)       // Get catalog with ownership status
		r.Get("/inventory", handlers.GetInventory)                // Get user inventory
		r.Get("/equipment", handlers.GetEquipment)                // Get equipped items
		r.Post("/equip", handlers.EquipItem)                      // Equip item
		r.Post("/purchase", handlers.PurchaseItem)                // Purchase item with coins
		r.Get("/notifications", handlers.GetUnlockNotifications)  // Get unlock notifications
	})

	// Learning Plans System (all protected)
	r.Route("/api/learning-plans", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)

		// Generate endpoint with extended timeout (120s) for full content generation
		r.With(middleware.Timeout(120 * time.Second)).Post("/generate", handlers.GenerateLearningPlanHandler)

		r.Get("/{id}", handlers.GetLearningPlanByIDHandler)                                  // Get plan by ID
		r.Get("/by-oa/{oa_bloom_objective_id}", handlers.GetLearningPlanByOAHandler)        // Get plan by OA
		r.Post("/{plan_id}/components/{component_id}/generate-content", handlers.GenerateComponentContentHandler) // Generate component content (legacy, for individual components)

		// Completion tracking
		r.Post("/{id}/start", handlers.StartLearningPlanHandler)       // Mark plan as started
		r.Post("/{id}/complete", handlers.CompleteLearningPlanHandler) // Mark plan as completed
	})

	// Text-to-Speech System (all protected)
	r.Route("/api/tts", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Post("/generate", handlers.GenerateTTSHandler) // Generate audio from text using ElevenLabs
	})

	// Recommendations System (all protected)
	r.Route("/api/recommendations", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleware)
		r.Get("/daily", handlers.GetDailyRecommendation) // Get personalized daily recommendation
	})

	// Static file server for avatars
	workDir, _ := os.Getwd()
	filesDir := http.Dir(workDir + "/static")
	FileServer(r, "/static", filesDir)

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

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
