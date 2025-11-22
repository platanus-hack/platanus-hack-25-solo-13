package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/models"
)

// UnlockEvent represents an event that can trigger item unlocks
type UnlockEvent struct {
	Type string                 // Event type: 'oa_complete', 'bloom_mastery', 'streak', 'level_up', 'diagnostic_achievement'
	Key  string                 // Event key for trigger lookup: 'oa_5', 'bloom_4_oa_10', 'streak_7', 'level_10'
	Data map[string]interface{} // Additional event data for validation
}

// UnlockService handles item unlock logic
type UnlockService struct{}

// NewUnlockService creates a new unlock service instance
func NewUnlockService() *UnlockService {
	return &UnlockService{}
}

// CheckAndUnlock is the main method that checks triggers and unlocks items
// This implements the Approach B2 (Push with event index)
func (s *UnlockService) CheckAndUnlock(userID uint, event UnlockEvent) ([]models.CustomizationItem, error) {
	// 1. Find triggers related to this specific event
	var triggerRecords []models.ItemUnlockTrigger
	err := db.DB.Preload("Item").
		Where("trigger_type = ? AND trigger_key = ?", event.Type, event.Key).
		Find(&triggerRecords).Error
	if err != nil {
		return nil, err
	}

	// If no items are associated with this event, return early
	if len(triggerRecords) == 0 {
		return []models.CustomizationItem{}, nil
	}

	// 2. Get item IDs from triggers
	itemIDs := make([]uint, len(triggerRecords))
	for i, tr := range triggerRecords {
		itemIDs[i] = tr.ItemID
	}

	// 3. Check which items user already has
	var unlockedItemIDs []uint
	db.DB.Model(&models.UserInventory{}).
		Where("user_id = ? AND item_id IN (?)", userID, itemIDs).
		Pluck("item_id", &unlockedItemIDs)

	// 4. Filter to only locked items
	lockedTriggers := []models.ItemUnlockTrigger{}
	for _, tr := range triggerRecords {
		if !contains(unlockedItemIDs, tr.ItemID) {
			lockedTriggers = append(lockedTriggers, tr)
		}
	}

	if len(lockedTriggers) == 0 {
		return []models.CustomizationItem{}, nil
	}

	// 5. Validate additional conditions and unlock
	newlyUnlocked := []models.CustomizationItem{}
	for _, trigger := range lockedTriggers {
		if s.ValidateAdditionalConditions(userID, trigger, event) {
			err := s.UnlockItem(userID, trigger.ItemID, "auto_unlock")
			if err == nil {
				newlyUnlocked = append(newlyUnlocked, trigger.Item)
			}
		}
	}

	return newlyUnlocked, nil
}

// ValidateAdditionalConditions validates trigger-specific conditions
func (s *UnlockService) ValidateAdditionalConditions(userID uint, trigger models.ItemUnlockTrigger, event UnlockEvent) bool {
	// Parse additional_data
	var additionalData models.AdditionalTriggerData
	if err := json.Unmarshal(trigger.AdditionalData, &additionalData); err != nil {
		return false // Invalid JSON, don't unlock
	}

	switch trigger.TriggerType {
	case "default":
		// Default items are always unlocked (should already be in inventory via seed)
		return true

	case "oa_complete":
		// Validate materia if specified
		if additionalData.MateriaID > 0 {
			if materiaID, ok := event.Data["materia_id"].(uint); ok {
				if materiaID != additionalData.MateriaID {
					return false
				}
			}
		}

		// If requires count, check how many OAs of this materia user has completed
		if additionalData.RequiresCount > 0 {
			var count int64
			db.DB.Model(&models.StudentOAProgress{}).
				Joins("JOIN oa_bloom_objectives ON oa_bloom_objectives.id = student_oa_progress.oa_bloom_objective_id").
				Joins("JOIN objetivos_aprendizaje ON objetivos_aprendizaje.id = oa_bloom_objectives.oa_id").
				Where("student_oa_progress.user_id = ? AND objetivos_aprendizaje.materia_id = ? AND student_oa_progress.estado IN ('dominado', 'logrado')",
					userID, additionalData.MateriaID).
				Count(&count)

			return int(count) >= additionalData.RequiresCount
		}

		return true // Single OA completion, already validated by trigger match

	case "bloom_mastery":
		// Check if user has required number of OAs dominated at this Bloom level
		if additionalData.RequiresCount > 0 && additionalData.BloomLevel > 0 {
			var count int64
			db.DB.Model(&models.StudentOAProgress{}).
				Joins("JOIN oa_bloom_objectives ON oa_bloom_objectives.id = student_oa_progress.oa_bloom_objective_id").
				Where("student_oa_progress.user_id = ? AND oa_bloom_objectives.bloom_level_id = ? AND student_oa_progress.estado IN ('dominado', 'logrado')",
					userID, additionalData.BloomLevel).
				Count(&count)

			return int(count) >= additionalData.RequiresCount
		}
		return false

	case "streak":
		// Check current streak
		if additionalData.Days > 0 {
			var gamification models.UserGamification
			if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
				return false
			}
			return gamification.CurrentStreak >= additionalData.Days
		}
		return false

	case "level_up":
		// Check user level
		if additionalData.Level > 0 {
			var gamification models.UserGamification
			if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
				return false
			}

			if gamification.Level < additionalData.Level {
				return false
			}

			// Check additional OA count requirement if specified
			if additionalData.AlsoRequiresOACount > 0 {
				var count int64
				db.DB.Model(&models.StudentOAProgress{}).
					Where("user_id = ? AND estado IN ('dominado', 'logrado')", userID).
					Count(&count)

				return int(count) >= additionalData.AlsoRequiresOACount
			}

			return true
		}
		return false

	case "diagnostic_achievement":
		// Check completed diagnostics with min score
		if additionalData.Count > 0 && additionalData.MinScore > 0 {
			var count int64
			db.DB.Model(&models.DiagnosticSession{}).
				Where("user_id = ? AND estado = 'completado' AND (preguntas_correctas * 100 / NULLIF(preguntas_totales, 0)) >= ?",
					userID, additionalData.MinScore).
				Count(&count)

			return int(count) >= additionalData.Count
		}
		return false

	case "coins":
		// Coin items are not auto-unlocked, they must be purchased
		return false

	default:
		return false
	}
}

// UnlockItem unlocks an item for a user
func (s *UnlockService) UnlockItem(userID uint, itemID uint, acquisitionType string) error {
	// Check if item exists
	var item models.CustomizationItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	// Check if user already has it
	var existing models.UserInventory
	err := db.DB.Where("user_id = ? AND item_id = ?", userID, itemID).First(&existing).Error
	if err == nil {
		return errors.New("user already has this item")
	}

	// Insert into inventory
	inventory := models.UserInventory{
		UserID:          userID,
		ItemID:          itemID,
		AcquisitionType: acquisitionType,
	}
	if err := db.DB.Create(&inventory).Error; err != nil {
		return err
	}

	// Create notification
	notification := models.UnlockNotification{
		UserID: userID,
		ItemID: itemID,
		IsRead: false,
	}
	db.DB.Create(&notification)

	return nil
}

// PurchaseItem purchases an item with coins
func (s *UnlockService) PurchaseItem(userID uint, itemID uint) error {
	// Get item
	var item models.CustomizationItem
	if err := db.DB.First(&item, itemID).Error; err != nil {
		return errors.New("item not found")
	}

	// Check if item is purchasable
	if item.BaseCoinsCost <= 0 {
		return errors.New("this item cannot be purchased with coins")
	}

	// Get user gamification
	var gamification models.UserGamification
	if err := db.DB.Where("user_id = ?", userID).First(&gamification).Error; err != nil {
		return errors.New("gamification data not found")
	}

	// Check if user has enough coins
	if gamification.Coins < item.BaseCoinsCost {
		return fmt.Errorf("insufficient coins. Need %d, have %d", item.BaseCoinsCost, gamification.Coins)
	}

	// Deduct coins
	gamification.Coins -= item.BaseCoinsCost
	if err := db.DB.Save(&gamification).Error; err != nil {
		return err
	}

	// Unlock item
	return s.UnlockItem(userID, itemID, "purchase")
}

// Helper function to check if slice contains a value
func contains(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
