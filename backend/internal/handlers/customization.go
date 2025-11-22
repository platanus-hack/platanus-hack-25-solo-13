package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/platanus-hack-25/lumera_app/internal/db"
	authmiddleware "github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"github.com/platanus-hack-25/lumera_app/internal/services"
)

var unlockService = services.NewUnlockService()

// ItemWithStatus represents an item with ownership status
type ItemWithStatus struct {
	models.CustomizationItem
	IsOwned      bool   `json:"is_owned"`
	CanPurchase  bool   `json:"can_purchase"`
	IsEquipped   bool   `json:"is_equipped"`
	Status       string `json:"status"` // 'owned', 'locked', 'can_purchase'
}

// GetCustomizationCatalog godoc
// @Summary Get customization catalog
// @Description Retrieve all customization items with user ownership status
// @Tags Customization
// @Produce json
// @Success 200 {array} ItemWithStatus
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customization/catalog [get]
func GetCustomizationCatalog(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all items
	var items []models.CustomizationItem
	if err := db.DB.Find(&items).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user's inventory
	var inventoryItemIDs []uint
	db.DB.Model(&models.UserInventory{}).
		Where("user_id = ?", userID).
		Pluck("item_id", &inventoryItemIDs)

	// Get user's equipment
	var equipment models.UserEquipment
	db.DB.Where("user_id = ?", userID).First(&equipment)

	// Get user's coins
	var gamification models.UserGamification
	db.DB.Where("user_id = ?", userID).First(&gamification)

	// Build response with status
	itemsWithStatus := make([]ItemWithStatus, len(items))
	for i, item := range items {
		isOwned := contains(inventoryItemIDs, item.ID)
		canPurchase := !isOwned && item.BaseCoinsCost > 0 && gamification.Coins >= item.BaseCoinsCost

		var isEquipped bool
		if item.Type == "avatar" && equipment.EquippedAvatarID != nil {
			isEquipped = *equipment.EquippedAvatarID == item.ID
		} else if item.Type == "frame" && equipment.EquippedFrameID != nil {
			isEquipped = *equipment.EquippedFrameID == item.ID
		}

		status := "locked"
		if isOwned {
			status = "owned"
		} else if canPurchase {
			status = "can_purchase"
		}

		itemsWithStatus[i] = ItemWithStatus{
			CustomizationItem: item,
			IsOwned:           isOwned,
			CanPurchase:       canPurchase,
			IsEquipped:        isEquipped,
			Status:            status,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itemsWithStatus)
}

// GetInventory godoc
// @Summary Get user inventory
// @Description Retrieve all items the user has unlocked
// @Tags Customization
// @Produce json
// @Success 200 {array} models.UserInventory
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customization/inventory [get]
func GetInventory(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var inventory []models.UserInventory
	if err := db.DB.Preload("Item").Where("user_id = ?", userID).Find(&inventory).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

// GetEquipment godoc
// @Summary Get equipped items
// @Description Retrieve currently equipped avatar and frame
// @Tags Customization
// @Produce json
// @Success 200 {object} models.UserEquipment
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customization/equipment [get]
func GetEquipment(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var equipment models.UserEquipment
	if err := db.DB.Preload("EquippedAvatar").Preload("EquippedFrame").
		Where("user_id = ?", userID).First(&equipment).Error; err != nil {
		http.Error(w, "Equipment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

// EquipItemRequest represents the request to equip an item
type EquipItemRequest struct {
	ItemID uint   `json:"item_id"`
	Slot   string `json:"slot"` // 'avatar' or 'frame'
}

// EquipItem godoc
// @Summary Equip an item
// @Description Equip an avatar or frame
// @Tags Customization
// @Accept json
// @Produce json
// @Param request body EquipItemRequest true "Equip request"
// @Success 200 {object} models.UserEquipment
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customization/equip [post]
func EquipItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req EquipItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate slot
	if req.Slot != "avatar" && req.Slot != "frame" {
		http.Error(w, "Invalid slot. Must be 'avatar' or 'frame'", http.StatusBadRequest)
		return
	}

	// Check if user owns the item
	var inventory models.UserInventory
	if err := db.DB.Where("user_id = ? AND item_id = ?", userID, req.ItemID).First(&inventory).Error; err != nil {
		http.Error(w, "You don't own this item", http.StatusBadRequest)
		return
	}

	// Get the item to validate type
	var item models.CustomizationItem
	if err := db.DB.First(&item, req.ItemID).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Validate item type matches slot
	if item.Type != req.Slot {
		http.Error(w, "Item type doesn't match slot", http.StatusBadRequest)
		return
	}

	// Get or create equipment
	var equipment models.UserEquipment
	err := db.DB.Where("user_id = ?", userID).First(&equipment).Error
	if err != nil {
		// Create equipment if doesn't exist
		equipment = models.UserEquipment{UserID: userID}
	}

	// Equip item
	if req.Slot == "avatar" {
		equipment.EquippedAvatarID = &req.ItemID
	} else {
		equipment.EquippedFrameID = &req.ItemID
	}

	if err := db.DB.Save(&equipment).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reload with preloads
	db.DB.Preload("EquippedAvatar").Preload("EquippedFrame").
		Where("user_id = ?", userID).First(&equipment)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(equipment)
}

// PurchaseItemRequest represents the request to purchase an item
type PurchaseItemRequest struct {
	ItemID uint `json:"item_id"`
}

// PurchaseItem godoc
// @Summary Purchase an item with coins
// @Description Purchase a customization item using coins
// @Tags Customization
// @Accept json
// @Produce json
// @Param request body PurchaseItemRequest true "Purchase request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customization/purchase [post]
func PurchaseItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req PurchaseItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Purchase item (handles coin deduction and unlock)
	if err := unlockService.PurchaseItem(userID, req.ItemID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated gamification data
	var gamification models.UserGamification
	db.DB.Where("user_id = ?", userID).First(&gamification)

	// Get the purchased item
	var item models.CustomizationItem
	db.DB.First(&item, req.ItemID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Item purchased successfully",
		"item":          item,
		"remaining_coins": gamification.Coins,
	})
}

// GetUnlockNotifications godoc
// @Summary Get unlock notifications
// @Description Retrieve unread unlock notifications
// @Tags Customization
// @Produce json
// @Success 200 {array} models.UnlockNotification
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/customization/notifications [get]
func GetUnlockNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var notifications []models.UnlockNotification
	if err := db.DB.Preload("Item").
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Mark as read
	db.DB.Model(&models.UnlockNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// Helper function (duplicate from unlock_service, could be moved to utils)
func contains(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
