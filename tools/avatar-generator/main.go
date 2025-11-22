package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/platanus-hack-25/lumera_app/avatar-generator/generator"
)

type AvatarConfig struct {
	Avatares []generator.Avatar `json:"avatares"`
}

func main() {
	log.Println("=== Avatar Generator con DALL-E ===")

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ No .env file found, using environment variables")
	}

	// Connect to database
	if err := generator.ConnectDB(); err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	// Initialize OpenAI
	generator.InitOpenAI()

	// Read avatar configuration
	configFile := "input/avatar_config.json"
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("❌ Failed to read config file: %v", err)
	}

	var config AvatarConfig
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("❌ Failed to parse config file: %v", err)
	}

	log.Printf("Found %d avatars to generate\n", len(config.Avatares))

	// Statistics
	stats := &generator.Stats{
		TotalAttempts: 0,
		SuccessCount:  0,
		FailCount:     0,
		Failed:        []generator.GenerationResult{},
		StartTime:     time.Now(),
	}

	// Category counters
	categoryCounts := map[string]int{
		"inicio":  0,
		"logro":   0,
		"compra":  0,
	}
	tierCounts := map[int]int{}

	var generatedAvatars []generator.Avatar

	// Process each avatar
	for i, avatar := range config.Avatares {
		log.Printf("\n[%d/%d] Generating '%s' (tier %d⭐, %s)",
			i+1, len(config.Avatares), avatar.Nombre, avatar.Tier, avatar.Categoria)

		stats.TotalAttempts++

		// Generate image with DALL-E
		imageURL, err := generator.GenerateAvatarImage(avatar)
		if err != nil {
			log.Printf("✗ Failed to generate image: %v", err)
			stats.AddFail(generator.GenerationResult{
				Avatar:    avatar,
				Success:   false,
				Error:     err.Error(),
				Timestamp: time.Now(),
			})
			continue
		}

		// Update avatar with image URL
		avatar.ImageURL = imageURL
		generatedAvatars = append(generatedAvatars, avatar)

		// Update stats
		stats.AddSuccess()
		categoryCounts[avatar.Categoria]++
		tierCounts[avatar.Tier]++

		log.Printf("✓ Successfully generated '%s'", avatar.Nombre)

		// Small delay to avoid rate limiting
		if i < len(config.Avatares)-1 {
			time.Sleep(2 * time.Second)
		}
	}

	// Insert all generated avatars into database
	if len(generatedAvatars) > 0 {
		log.Printf("\n💾 Inserting %d avatars into database...", len(generatedAvatars))
		if err := generator.InsertAvatarsBatch(generatedAvatars); err != nil {
			log.Printf("❌ Failed to insert avatars: %v", err)
		} else {
			log.Printf("✓ Successfully inserted %d avatars", len(generatedAvatars))
		}
	}

	// Save failed avatars to JSON
	if len(stats.Failed) > 0 {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("output/failed_avatars_%s.json", timestamp)
		data, _ := json.MarshalIndent(stats.Failed, "", "  ")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			log.Printf("⚠ Failed to save failed avatars: %v", err)
		} else {
			log.Printf("\n📝 Saved %d failed avatars to %s", len(stats.Failed), filename)
		}
	}

	// Print final statistics
	stats.EndTime = time.Now()
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("\n============================================================")
	fmt.Println("📊 AVATAR GENERATION STATISTICS")
	fmt.Println("============================================================")
	fmt.Printf("Total Duration:     %v\n", duration.Round(time.Second))
	fmt.Printf("Total Attempts:     %d\n", stats.TotalAttempts)
	fmt.Printf("✓ Successful:       %d (%.1f%%)\n",
		stats.SuccessCount,
		float64(stats.SuccessCount)/float64(stats.TotalAttempts)*100)
	fmt.Printf("✗ Failed:           %d (%.1f%%)\n",
		stats.FailCount,
		float64(stats.FailCount)/float64(stats.TotalAttempts)*100)

	fmt.Println("\n📋 Avatars by Category:")
	for cat, count := range categoryCounts {
		fmt.Printf("  - %-10s: %d\n", cat, count)
	}

	fmt.Println("\n⭐ Avatars by Tier:")
	for tier := 1; tier <= 5; tier++ {
		if count, ok := tierCounts[tier]; ok {
			stars := ""
			for i := 0; i < tier; i++ {
				stars += "⭐"
			}
			fmt.Printf("  - Tier %d %s: %d\n", tier, stars, count)
		}
	}

	fmt.Println("============================================================")

	if stats.FailCount > 0 {
		fmt.Printf("\n⚠ Some avatars failed to generate. Check output/failed_avatars_*.json for details\n")
	} else {
		fmt.Printf("\n🎉 All avatars generated successfully!\n")
	}
}
