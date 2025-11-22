package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/models"
)

// UpdateUserRequest represents the update user payload
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// ChangePasswordRequest represents the change password payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// GetMe godoc
// @Summary Get current user profile
// @Description Returns the authenticated user's profile information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "User profile"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/users/me [get]
func GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateMe godoc
// @Summary Update current user profile
// @Description Updates the authenticated user's name and/or email
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserRequest true "Fields to update"
// @Success 200 {object} models.User "Updated user profile"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 409 {object} map[string]string "Email already in use"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/me [put]
func UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Find user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		// Normalize email
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Check if email is already taken by another user
		var existingUser models.User
		if err := db.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			http.Error(w, `{"error":"email already in use"}`, http.StatusConflict)
			return
		}

		user.Email = req.Email
	}

	// Save updates
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, `{"error":"failed to update user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ChangePassword godoc
// @Summary Change user password
// @Description Changes the authenticated user's password
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Current and new password"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} map[string]string "Invalid request body or missing required fields"
// @Failure 401 {object} map[string]string "Unauthorized - missing/invalid token or incorrect current password"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/change-password [post]
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, `{"error":"current_password and new_password are required"}`, http.StatusBadRequest)
		return
	}

	// Find user
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	// Verify current password
	if !user.CheckPassword(req.CurrentPassword) {
		http.Error(w, `{"error":"current password is incorrect"}`, http.StatusUnauthorized)
		return
	}

	// Set new password
	if err := user.SetPassword(req.NewPassword); err != nil {
		http.Error(w, `{"error":"failed to hash new password"}`, http.StatusInternalServerError)
		return
	}

	// Save updates
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, `{"error":"failed to update password"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "password changed successfully"})
}

// DeleteMe godoc
// @Summary Delete user account
// @Description Permanently deletes the authenticated user's account
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Account deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/me [delete]
func DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Delete user
	if err := db.DB.Delete(&models.User{}, userID).Error; err != nil {
		http.Error(w, `{"error":"failed to delete user"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "account deleted successfully"})
}
