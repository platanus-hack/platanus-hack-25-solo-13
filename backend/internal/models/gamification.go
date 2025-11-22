package models

import (
	"time"

	"gorm.io/datatypes"
)

// UserGamification stores gamification stats per user
type UserGamification struct {
	UserID           uint       `json:"user_id" gorm:"primaryKey"`
	Level            int        `json:"level" gorm:"not null;default:1"`
	XP               int        `json:"xp" gorm:"not null;default:0"`
	Coins            int        `json:"coins" gorm:"not null;default:100"`
	CurrentStreak    int        `json:"current_streak" gorm:"not null;default:0"`
	LongestStreak    int        `json:"longest_streak" gorm:"not null;default:0"`
	LastActivityDate *time.Time `json:"last_activity_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName overrides the default table name
func (UserGamification) TableName() string {
	return "user_gamification"
}

// CustomizationItem represents an item in the global catalog (avatars, frames, etc)
type CustomizationItem struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"size:100;not null"`
	Type            string         `json:"type" gorm:"size:20;not null"` // 'avatar' or 'frame'
	Rarity          string         `json:"rarity" gorm:"size:20;not null"` // 'common', 'rare', 'epic', 'legendary'
	ImageURL        string         `json:"image_url" gorm:"type:text;not null"`
	Description     string         `json:"description" gorm:"type:text"`
	UnlockCondition datatypes.JSON `json:"unlock_condition" gorm:"type:jsonb;not null"` // For display/reference
	BaseCoinsCost   int            `json:"base_coins_cost" gorm:"default:0"`
	IsDefault       bool           `json:"is_default" gorm:"default:false"`
	CreatedAt       time.Time      `json:"created_at"`

	// Relationships
	Triggers []ItemUnlockTrigger `json:"triggers,omitempty" gorm:"foreignKey:ItemID"`
}

// TableName overrides the default table name
func (CustomizationItem) TableName() string {
	return "customization_items"
}

// ItemUnlockTrigger represents the event index mapping events to items
type ItemUnlockTrigger struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ItemID         uint           `json:"item_id" gorm:"not null"`
	TriggerType    string         `json:"trigger_type" gorm:"size:50;not null"` // 'oa_complete', 'bloom_mastery', 'streak', etc
	TriggerKey     string         `json:"trigger_key" gorm:"size:100;not null"` // 'oa_5', 'bloom_4_oa_10', 'streak_7', etc
	AdditionalData datatypes.JSON `json:"additional_data" gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time      `json:"created_at"`

	// Relationships
	Item CustomizationItem `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// TableName overrides the default table name
func (ItemUnlockTrigger) TableName() string {
	return "item_unlock_triggers"
}

// UnlockCondition represents the structure of unlock_condition JSON
type UnlockCondition struct {
	Type   string      `json:"type"` // 'default', 'oa_complete', 'bloom_mastery', 'streak', 'coins', 'level', 'multiple'
	Value  interface{} `json:"value,omitempty"`
	OAID   uint        `json:"oa_id,omitempty"`
	MateriaID uint     `json:"materia_id,omitempty"`
	MateriaName string `json:"materia,omitempty"`
	BloomLevel int     `json:"bloom_level,omitempty"`
	Count      int     `json:"count,omitempty"`
	Days       int     `json:"days,omitempty"`
	Amount     int     `json:"amount,omitempty"`
	MinScore   int     `json:"min_score,omitempty"`
	Conditions []UnlockCondition `json:"conditions,omitempty"` // For 'multiple' type
}

// AdditionalTriggerData represents the structure of additional_data in triggers
type AdditionalTriggerData struct {
	MateriaID        uint   `json:"materia_id,omitempty"`
	MateriaNombre    string `json:"materia_nombre,omitempty"`
	RequiresCount    int    `json:"requires_count,omitempty"`
	BloomLevel       int    `json:"bloom_level,omitempty"`
	BloomName        string `json:"bloom_name,omitempty"`
	Days             int    `json:"days,omitempty"`
	Amount           int    `json:"amount,omitempty"`
	Count            int    `json:"count,omitempty"`
	MinScore         int    `json:"min_score,omitempty"`
	Level            int    `json:"level,omitempty"`
	AlsoRequiresOACount int `json:"also_requires_oa_count,omitempty"`
}
