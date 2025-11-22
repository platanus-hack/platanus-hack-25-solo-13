package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"gorm.io/datatypes"
)

// CreateProfileRequest represents the request to create a student profile
type CreateProfileRequest struct {
	UserID      uint           `json:"user_id"`
	Edad        *int           `json:"edad"`
	CursoActual string         `json:"curso_actual"`
	ProfileData datatypes.JSON `json:"profile_data"`
}

// UpdateProfileRequest represents the request to update a profile
type UpdateProfileRequest struct {
	Edad        *int           `json:"edad,omitempty"`
	CursoActual string         `json:"curso_actual,omitempty"`
	ProfileData datatypes.JSON `json:"profile_data,omitempty"`
}

// GetProfile godoc
// @Summary Get student profile
// @Description Get the complete profile of a student by user ID
// @Tags Profiles
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Student profile"
// @Failure 404 {object} map[string]string "Profile not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/profiles/{user_id} [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid user_id"}`, http.StatusBadRequest)
		return
	}

	var profile models.StudentProfile
	if err := db.DB.Preload("User").Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, `{"error":"profile not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// CreateProfile godoc
// @Summary Create student profile
// @Description Create a new student profile with adaptive learning data
// @Tags Profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProfileRequest true "Profile data"
// @Success 201 {object} map[string]interface{} "Created profile"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 409 {object} map[string]string "Profile already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/profiles [post]
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Check if profile already exists
	var existing models.StudentProfile
	if err := db.DB.Where("user_id = ?", req.UserID).First(&existing).Error; err == nil {
		http.Error(w, `{"error":"profile already exists for this user"}`, http.StatusConflict)
		return
	}

	// Initialize default profile_data if empty
	if len(req.ProfileData) == 0 {
		req.ProfileData = datatypes.JSON([]byte("{}"))
	}

	profile := models.StudentProfile{
		UserID:      req.UserID,
		Edad:        req.Edad,
		CursoActual: req.CursoActual,
		ProfileData: req.ProfileData,
	}

	if err := db.DB.Create(&profile).Error; err != nil {
		http.Error(w, `{"error":"failed to create profile"}`, http.StatusInternalServerError)
		return
	}

	// Load user relation
	db.DB.Preload("User").First(&profile, profile.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile godoc
// @Summary Update student profile
// @Description Update student profile data (partial update)
// @Tags Profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Param request body UpdateProfileRequest true "Updated profile data"
// @Success 200 {object} map[string]interface{} "Updated profile"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Profile not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/profiles/{user_id} [patch]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid user_id"}`, http.StatusBadRequest)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Find profile
	var profile models.StudentProfile
	if err := db.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, `{"error":"profile not found"}`, http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Edad != nil {
		profile.Edad = req.Edad
	}
	if req.CursoActual != "" {
		profile.CursoActual = req.CursoActual
	}
	if len(req.ProfileData) > 0 {
		// Merge strategy: for now, replace completely
		// TODO: Implement deep merge if needed
		profile.ProfileData = req.ProfileData
	}

	// Save
	if err := db.DB.Save(&profile).Error; err != nil {
		http.Error(w, `{"error":"failed to update profile"}`, http.StatusInternalServerError)
		return
	}

	// Load user relation
	db.DB.Preload("User").First(&profile, profile.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// ExportProfiles godoc
// @Summary Export all profiles for ML
// @Description Export all student profiles in JSON format ready for machine learning
// @Tags Profiles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Export result with profiles array"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/profiles/export [get]
func ExportProfiles(w http.ResponseWriter, r *http.Request) {
	// Only allow admins or authenticated users with proper permissions
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// TODO: Add role check for admin

	var profiles []models.StudentProfile
	if err := db.DB.Find(&profiles).Error; err != nil {
		http.Error(w, `{"error":"failed to export profiles"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":    len(profiles),
		"profiles": profiles,
		"exported_by": userID,
	})
}
