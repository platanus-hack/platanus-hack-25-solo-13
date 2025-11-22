package services

import (
	"errors"
	"time"

	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/models"
)

// GamificationService handles XP, coins, levels, and streaks
type GamificationService struct {
	unlockService *UnlockService
}

// NewGamificationService creates a new gamification service instance
func NewGamificationService() *GamificationService {
	return &GamificationService{
		unlockService: NewUnlockService(),
	}
}

// AddXPResult represents the result of adding XP
type AddXPResult struct {
	NewXP       int  `json:"new_xp"`
	NewLevel    int  `json:"new_level"`
	LeveledUp   bool `json:"leveled_up"`
	OldLevel    int  `json:"old_level,omitempty"`
	ItemsUnlocked []models.CustomizationItem `json:"items_unlocked,omitempty"`
}

// AddXP adds XP to a user and checks for level up
func (s *GamificationService) AddXP(userID uint, amount int, reason string) (*AddXPResult, error) {
	var gamification models.UserGamification
	if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
		return nil, errors.New("gamification data not found")
	}

	oldLevel := gamification.Level

	// Add XP (trigger will auto-calculate new level)
	gamification.XP += amount
	if err := db.DB.Save(&gamification).Error; err != nil {
		return nil, err
	}

	// Re-fetch to get updated level (from trigger)
	db.DB.Where("user_id = ?", userID).First(&gamification)

	leveledUp := gamification.Level > oldLevel
	result := &AddXPResult{
		NewXP:     gamification.XP,
		NewLevel:  gamification.Level,
		LeveledUp: leveledUp,
		OldLevel:  oldLevel,
	}

	// If leveled up, check for level-based unlocks
	if leveledUp {
		unlockedItems, _ := s.unlockService.CheckAndUnlock(userID, UnlockEvent{
			Type: "level_up",
			Key:  formatLevelKey(gamification.Level),
			Data: map[string]interface{}{
				"level": gamification.Level,
			},
		})
		result.ItemsUnlocked = unlockedItems
	}

	return result, nil
}

// AddCoins adds coins to a user
func (s *GamificationService) AddCoins(userID uint, amount int, reason string) error {
	var gamification models.UserGamification
	if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
		return errors.New("gamification data not found")
	}

	gamification.Coins += amount
	return db.DB.Save(&gamification).Error
}

// DeductCoins deducts coins from a user (used internally by purchase)
func (s *GamificationService) DeductCoins(userID uint, amount int) error {
	var gamification models.UserGamification
	if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
		return errors.New("gamification data not found")
	}

	if gamification.Coins < amount {
		return errors.New("insufficient coins")
	}

	gamification.Coins -= amount
	return db.DB.Save(&gamification).Error
}

// UpdateStreakResult represents the result of updating streak
type UpdateStreakResult struct {
	CurrentStreak int  `json:"current_streak"`
	LongestStreak int  `json:"longest_streak"`
	IsNewRecord   bool `json:"is_new_record"`
	ItemsUnlocked []models.CustomizationItem `json:"items_unlocked,omitempty"`
}

// UpdateStreak updates the user's activity streak
func (s *GamificationService) UpdateStreak(userID uint) (*UpdateStreakResult, error) {
	var gamification models.UserGamification
	if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
		return nil, errors.New("gamification data not found")
	}

	today := time.Now().Truncate(24 * time.Hour)
	oldStreak := gamification.CurrentStreak

	// If last activity was today, no change
	if gamification.LastActivityDate != nil {
		lastActivity := gamification.LastActivityDate.Truncate(24 * time.Hour)
		if lastActivity.Equal(today) {
			return &UpdateStreakResult{
				CurrentStreak: gamification.CurrentStreak,
				LongestStreak: gamification.LongestStreak,
				IsNewRecord:   false,
			}, nil
		}

		yesterday := today.Add(-24 * time.Hour)
		if lastActivity.Equal(yesterday) {
			// Consecutive day
			gamification.CurrentStreak++
		} else {
			// Streak broken
			gamification.CurrentStreak = 1
		}
	} else {
		// First activity
		gamification.CurrentStreak = 1
	}

	now := time.Now()
	gamification.LastActivityDate = &now

	// Trigger will auto-update longest_streak
	if err := db.DB.Save(&gamification).Error; err != nil {
		return nil, err
	}

	// Re-fetch to get updated longest_streak
	db.DB.Where("user_id = ?", userID).First(&gamification)

	result := &UpdateStreakResult{
		CurrentStreak: gamification.CurrentStreak,
		LongestStreak: gamification.LongestStreak,
		IsNewRecord:   gamification.CurrentStreak > oldStreak && gamification.CurrentStreak == gamification.LongestStreak,
	}

	// Check for streak-based unlocks if streak increased
	if gamification.CurrentStreak > oldStreak {
		unlockedItems, _ := s.unlockService.CheckAndUnlock(userID, UnlockEvent{
			Type: "streak",
			Key:  formatStreakKey(gamification.CurrentStreak),
			Data: map[string]interface{}{
				"days": gamification.CurrentStreak,
			},
		})
		result.ItemsUnlocked = unlockedItems
	}

	return result, nil
}

// CalculateStreak calculates streak from activity history (expensive, use sparingly)
func (s *GamificationService) CalculateStreak(userID uint) (int, error) {
	// Get distinct activity dates from student_oa_history
	var dates []time.Time
	err := db.DB.Model(&models.StudentOAHistory{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Pluck("DATE(created_at)", &dates).Error
	if err != nil {
		return 0, err
	}

	if len(dates) == 0 {
		return 0, nil
	}

	// Count consecutive days from today backwards
	today := time.Now().Truncate(24 * time.Hour)
	streak := 0
	currentDate := today

	for _, date := range dates {
		dateOnly := date.Truncate(24 * time.Hour)
		if dateOnly.Equal(currentDate) {
			streak++
			currentDate = currentDate.Add(-24 * time.Hour)
		} else if dateOnly.Before(currentDate) {
			// Gap found, streak ends
			break
		}
	}

	return streak, nil
}

// GetGamificationStats retrieves user gamification stats
func (s *GamificationService) GetGamificationStats(userID uint) (*models.UserGamification, error) {
	var gamification models.UserGamification
	if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
		return nil, errors.New("gamification data not found")
	}
	return &gamification, nil
}

// GetLeaderboard retrieves top users by XP
func (s *GamificationService) GetLeaderboard(limit int) ([]models.UserGamification, error) {
	var leaderboard []models.UserGamification
	err := db.DB.Preload("User").
		Order("xp DESC").
		Limit(limit).
		Find(&leaderboard).Error
	return leaderboard, err
}

// Helper functions
func formatLevelKey(level int) string {
	return "level_" + formatInt(level)
}

func formatStreakKey(days int) string {
	return "streak_" + formatInt(days)
}

func formatInt(n int) string {
	return string(rune('0' + n/10)) + string(rune('0' + n%10))
}
