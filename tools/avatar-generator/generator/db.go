package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB conecta a PostgreSQL usando GORM
func ConnectDB() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✓ Database connected successfully")
	return nil
}

// InsertAvatar inserta un avatar y su trigger de unlock en la base de datos
func InsertAvatar(avatar Avatar) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Preparar unlock_condition como JSONB
		unlockCondition := map[string]interface{}{}
		if avatar.UnlockTrigger.TriggerType != "" && avatar.UnlockTrigger.TriggerType != "default" {
			unlockCondition["type"] = avatar.UnlockTrigger.TriggerType
			unlockCondition["display_text"] = avatar.UnlockTrigger.DisplayText
			if avatar.UnlockTrigger.TriggerKey != "" {
				unlockCondition["trigger_key"] = avatar.UnlockTrigger.TriggerKey
			}
			if avatar.UnlockTrigger.ExtraData != nil {
				for k, v := range avatar.UnlockTrigger.ExtraData {
					unlockCondition[k] = v
				}
			}
		}

		unlockConditionJSON, err := json.Marshal(unlockCondition)
		if err != nil {
			return fmt.Errorf("failed to marshal unlock_condition: %w", err)
		}

		// Insertar en customization_items
		var itemID int
		err = tx.Raw(`
			INSERT INTO customization_items (
				name,
				type,
				rarity,
				tier,
				image_url,
				description,
				unlock_condition,
				base_coins_cost,
				is_default,
				created_at
			) VALUES (?, 'avatar', ?, ?, ?, ?, ?::jsonb, ?, ?, NOW())
			RETURNING id
		`, avatar.Nombre, avatar.Rarity, avatar.Tier, avatar.ImageURL, avatar.Descripcion,
			string(unlockConditionJSON), avatar.PrecioPuntos, avatar.IsDefault).Scan(&itemID).Error

		if err != nil {
			return fmt.Errorf("failed to insert avatar: %w", err)
		}

		// Insertar trigger de unlock si no es default
		if avatar.UnlockTrigger.TriggerType != "" && avatar.UnlockTrigger.TriggerType != "default" {
			additionalData := map[string]interface{}{}
			if avatar.UnlockTrigger.ExtraData != nil {
				additionalData = avatar.UnlockTrigger.ExtraData
			}

			additionalDataJSON, err := json.Marshal(additionalData)
			if err != nil {
				return fmt.Errorf("failed to marshal additional_data: %w", err)
			}

			triggerKey := avatar.UnlockTrigger.TriggerKey
			if triggerKey == "" {
				triggerKey = fmt.Sprintf("%s_%d", avatar.UnlockTrigger.TriggerType, itemID)
			}

			err = tx.Exec(`
				INSERT INTO item_unlock_triggers (
					item_id,
					trigger_type,
					trigger_key,
					additional_data,
					created_at
				) VALUES (?, ?, ?, ?::jsonb, NOW())
			`, itemID, avatar.UnlockTrigger.TriggerType, triggerKey, string(additionalDataJSON)).Error

			if err != nil {
				return fmt.Errorf("failed to insert unlock trigger: %w", err)
			}
		}

		log.Printf("✓ Inserted avatar '%s' (ID: %d, tier: %d⭐)", avatar.Nombre, itemID, avatar.Tier)
		return nil
	})
}

// InsertAvatarsBatch inserta múltiples avatares en lotes
func InsertAvatarsBatch(avatars []Avatar) error {
	if len(avatars) == 0 {
		return nil
	}

	batchSize := 10
	for i := 0; i < len(avatars); i += batchSize {
		end := i + batchSize
		if end > len(avatars) {
			end = len(avatars)
		}
		batch := avatars[i:end]

		for _, avatar := range batch {
			if err := InsertAvatar(avatar); err != nil {
				return fmt.Errorf("failed to insert avatar '%s': %w", avatar.Nombre, err)
			}
		}

		log.Printf("✓ Inserted batch %d-%d avatars", i+1, end)
	}

	log.Printf("✓ Successfully inserted %d avatars", len(avatars))
	return nil
}
