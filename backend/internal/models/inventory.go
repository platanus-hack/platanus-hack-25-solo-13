package models

import (
	"time"
)

// UserInventory represents items unlocked by a user
type UserInventory struct {
	UserID          uint      `json:"user_id" gorm:"primaryKey"`
	ItemID          uint      `json:"item_id" gorm:"primaryKey"`
	AcquiredAt      time.Time `json:"acquired_at"`
	AcquisitionType string    `json:"acquisition_type" gorm:"size:20;not null;default:'auto_unlock'"` // 'auto_unlock', 'purchase', 'gift', 'default'

	// Relationships
	User User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Item CustomizationItem `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// TableName overrides the default table name
func (UserInventory) TableName() string {
	return "user_inventory"
}

// UserEquipment represents currently equipped items for a user
type UserEquipment struct {
	UserID           uint      `json:"user_id" gorm:"primaryKey"`
	EquippedAvatarID *uint     `json:"equipped_avatar_id"`
	EquippedFrameID  *uint     `json:"equipped_frame_id"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relationships
	User           User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	EquippedAvatar *CustomizationItem `json:"equipped_avatar,omitempty" gorm:"foreignKey:EquippedAvatarID"`
	EquippedFrame  *CustomizationItem `json:"equipped_frame,omitempty" gorm:"foreignKey:EquippedFrameID"`
}

// TableName overrides the default table name
func (UserEquipment) TableName() string {
	return "user_equipment"
}

// UnlockNotification represents notifications for newly unlocked items
type UnlockNotification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ItemID    uint      `json:"item_id" gorm:"not null"`
	IsRead    bool      `json:"is_read" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Item CustomizationItem `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// TableName overrides the default table name
func (UnlockNotification) TableName() string {
	return "unlock_notifications"
}
