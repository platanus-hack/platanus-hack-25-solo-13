package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"github.com/platanus-hack-25/lumera_app/internal/utils"
)

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account and returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 200 {object} AuthResponse "Successfully registered user with JWT token"
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 409 {object} map[string]string "Email already registered"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Name == "" || req.Password == "" {
		http.Error(w, `{"error":"email, name, and password are required"}`, http.StatusBadRequest)
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user already exists
	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, `{"error":"email already registered"}`, http.StatusConflict)
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Create new user
	user := models.User{
		Email: req.Email,
		Name:  req.Name,
		Role:  req.Role,
	}

	// Hash password
	if err := user.SetPassword(req.Password); err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// Save to database
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, `{"error":"failed to create user"}`, http.StatusInternalServerError)
		return
	}

	// Initialize gamification data for new user
	gamification := models.UserGamification{
		UserID: user.ID,
		Level:  1,
		XP:     0,
		Coins:  100,
	}
	db.DB.Create(&gamification)

	// Initialize equipment record
	equipment := models.UserEquipment{
		UserID: user.ID,
	}
	db.DB.Create(&equipment)

	// Add default items to inventory (async to not block response)
	go func() {
		var defaultItems []models.CustomizationItem
		db.DB.Where("is_default = ?", true).Find(&defaultItems)

		for _, item := range defaultItems {
			inventory := models.UserInventory{
				UserID:          user.ID,
				ItemID:          item.ID,
				AcquisitionType: "default",
			}
			db.DB.Create(&inventory)
		}
	}()

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  &user,
	})
}

// Login godoc
// @Summary Authenticate a user
// @Description Authenticates user with email and password, returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} AuthResponse "Successfully authenticated with JWT token"
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"email and password are required"}`, http.StatusBadRequest)
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Find user by email
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  &user,
	})
}
