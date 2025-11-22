package handlers

import (
	"encoding/json"
	"net/http"

	authmiddleware "github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/services"
)

var gamificationService = services.NewGamificationService()

// GamificationStatsResponse represents gamification stats with calculated fields
type GamificationStatsResponse struct {
	UserID           uint   `json:"user_id"`
	Level            int    `json:"level"`
	XP               int    `json:"xp"`
	XPForNextLevel   int    `json:"xp_for_next_level"`
	XPProgress       int    `json:"xp_progress"` // XP earned in current level
	Coins            int    `json:"coins"`
	CurrentStreak    int    `json:"current_streak"`
	LongestStreak    int    `json:"longest_streak"`
	LastActivityDate string `json:"last_activity_date,omitempty"`
}

// GetGamificationStats godoc
// @Summary Get gamification stats
// @Description Retrieve user's gamification stats (XP, level, coins, streaks)
// @Tags Gamification
// @Produce json
// @Success 200 {object} GamificationStatsResponse
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/gamification/stats [get]
func GetGamificationStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	gamification, err := gamificationService.GetGamificationStats(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Calculate XP for next level
	// Formula: level = floor(sqrt(xp / 100)) + 1
	// Reverse: xp_for_level = (level - 1)^2 * 100
	currentLevelXP := (gamification.Level - 1) * (gamification.Level - 1) * 100
	nextLevelXP := gamification.Level * gamification.Level * 100
	xpProgress := gamification.XP - currentLevelXP

	response := GamificationStatsResponse{
		UserID:         gamification.UserID,
		Level:          gamification.Level,
		XP:             gamification.XP,
		XPForNextLevel: nextLevelXP,
		XPProgress:     xpProgress,
		Coins:          gamification.Coins,
		CurrentStreak:  gamification.CurrentStreak,
		LongestStreak:  gamification.LongestStreak,
	}

	if gamification.LastActivityDate != nil {
		response.LastActivityDate = gamification.LastActivityDate.Format("2006-01-02")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LeaderboardEntry represents a leaderboard entry
type LeaderboardEntry struct {
	Rank   int    `json:"rank"`
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Level  int    `json:"level"`
	XP     int    `json:"xp"`
}

// LeaderboardResponse represents the leaderboard with user's position
type LeaderboardResponse struct {
	Leaderboard  []LeaderboardEntry `json:"leaderboard"`
	UserPosition *LeaderboardEntry  `json:"user_position,omitempty"`
}

// GetLeaderboard godoc
// @Summary Get leaderboard
// @Description Retrieve top users by XP
// @Tags Gamification
// @Produce json
// @Success 200 {object} LeaderboardResponse
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/gamification/leaderboard [get]
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get top 10
	topUsers, err := gamificationService.GetLeaderboard(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build leaderboard entries
	leaderboard := make([]LeaderboardEntry, len(topUsers))
	var userEntry *LeaderboardEntry

	for i, user := range topUsers {
		entry := LeaderboardEntry{
			Rank:   i + 1,
			UserID: user.UserID,
			Name:   user.User.Name,
			Level:  user.Level,
			XP:     user.XP,
		}
		leaderboard[i] = entry

		// Check if this is the current user
		if user.UserID == userID {
			userEntry = &entry
		}
	}

	// If user not in top 10, find their position
	if userEntry == nil {
		userGamification, err := gamificationService.GetGamificationStats(userID)
		if err == nil {
			// Count how many users have more XP
			var rank int64
			// This is a simplified rank calculation
			// In production, you'd want a more efficient query
			rank = int64(len(topUsers)) + 1 // Placeholder, user is outside top 10

			userEntry = &LeaderboardEntry{
				Rank:   int(rank),
				UserID: userID,
				Name:   "(You)", // Would need to fetch user name
				Level:  userGamification.Level,
				XP:     userGamification.XP,
			}
		}
	}

	response := LeaderboardResponse{
		Leaderboard:  leaderboard,
		UserPosition: userEntry,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
