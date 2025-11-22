package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/platanus-hack-25/lumera_app/question-generator/generator"
)

func main() {
	// Flags
	skipExisting := flag.Bool("skip-existing", true, "Skip objectives that already have questions")
	batchSize := flag.Int("batch-size", 10, "Number of objectives to process before saving to database")
	flag.Parse()

	log.Println("=== Question Generator for Lumera App ===")
	log.Printf("Skip existing: %v", *skipExisting)
	log.Printf("Batch size: %d objectives\n", *batchSize)

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

	// Initialize stats
	stats := &generator.Stats{
		TotalAttempts:   0,
		SuccessCount:    0,
		FailCount:       0,
		FailedQuestions: []generator.FailedQuestion{},
		TypeCounts:      make(map[string]int),
		StartTime:       time.Now(),
	}

	// Fetch all OA-Bloom objectives
	log.Println("\n📚 Fetching OA-Bloom objectives from database...")
	objectives, err := generator.GetOABloomObjectives()
	if err != nil {
		log.Fatalf("❌ Failed to fetch objectives: %v", err)
	}

	log.Printf("Found %d OA-Bloom objectives\n", len(objectives))
	log.Printf("Expected to generate ~%d questions (5 per objective)\n", len(objectives)*5)

	// Process objectives
	var allQuestions []generator.Question
	var allFailed []generator.FailedQuestion

	for i, objective := range objectives {
		log.Printf("\n[%d/%d] Processing OA-Bloom #%d: %s (Bloom: %s - Level %d)",
			i+1, len(objectives), objective.ID, objective.OATitulo,
			objective.BloomLevelNombre, objective.BloomLevelNumero)

		// Skip if already has questions
		if *skipExisting {
			count, err := generator.CountExistingQuestions(objective.ID)
			if err != nil {
				log.Printf("⚠ Error checking existing questions: %v", err)
			} else if count > 0 {
				log.Printf("⏭ Skipping (already has %d questions)", count)
				continue
			}
		}

		// Generate questions
		questions, failed := generator.GenerateQuestionsForObjective(objective, stats)
		allQuestions = append(allQuestions, questions...)
		allFailed = append(allFailed, failed...)

		// Save batch to database
		if len(allQuestions) >= *batchSize*5 {
			log.Printf("\n💾 Saving batch of %d questions to database...", len(allQuestions))
			if err := generator.InsertQuestions(allQuestions); err != nil {
				log.Printf("❌ Failed to insert questions: %v", err)
				// Save to failed
				for _, q := range allQuestions {
					allFailed = append(allFailed, generator.FailedQuestion{
						OABloomObjectiveID: q.OABloomObjectiveID,
						Tipo:               q.Tipo,
						Dificultad:         q.DificultadRelativa,
						Error:              "database insertion failed",
						Timestamp:          time.Now(),
					})
				}
			} else {
				log.Printf("✓ Batch saved successfully")
			}
			allQuestions = []generator.Question{} // Clear batch
		}

		// Progress update
		if (i+1)%10 == 0 {
			elapsed := time.Since(stats.StartTime)
			rate := float64(i+1) / elapsed.Seconds()
			remaining := float64(len(objectives)-(i+1)) / rate
			log.Printf("\n📊 Progress: %d/%d objectives (%.1f obj/sec, ETA: %.1f min)",
				i+1, len(objectives), rate, remaining/60)
		}
	}

	// Save remaining questions
	if len(allQuestions) > 0 {
		log.Printf("\n💾 Saving final batch of %d questions...", len(allQuestions))
		if err := generator.InsertQuestions(allQuestions); err != nil {
			log.Printf("❌ Failed to insert final batch: %v", err)
		} else {
			log.Printf("✓ Final batch saved successfully")
		}
	}

	stats.EndTime = time.Now()

	// Save failed questions to JSON
	if len(allFailed) > 0 {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("output/failed_questions_%s.json", timestamp)
		data, _ := json.MarshalIndent(allFailed, "", "  ")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			log.Printf("⚠ Failed to save failed questions: %v", err)
		} else {
			log.Printf("\n📝 Saved %d failed questions to %s", len(allFailed), filename)
		}
	}

	// Print final stats
	printStats(stats)
}

func printStats(stats *generator.Stats) {
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("\n" + repeat("=", 60))
	fmt.Println("📊 FINAL STATISTICS")
	fmt.Println(repeat("=", 60))
	fmt.Printf("Total Duration:     %v\n", duration.Round(time.Second))
	fmt.Printf("Total Attempts:     %d\n", stats.TotalAttempts)
	fmt.Printf("✓ Successful:       %d (%.1f%%)\n",
		stats.SuccessCount,
		float64(stats.SuccessCount)/float64(stats.TotalAttempts)*100)
	fmt.Printf("✗ Failed:           %d (%.1f%%)\n",
		stats.FailCount,
		float64(stats.FailCount)/float64(stats.TotalAttempts)*100)

	if len(stats.TypeCounts) > 0 {
		fmt.Println("\n📋 Questions by Type:")
		for qType, count := range stats.TypeCounts {
			fmt.Printf("  - %-20s: %d\n", qType, count)
		}
	}

	fmt.Println(repeat("=", 60))

	if stats.FailCount > 0 {
		fmt.Printf("\n⚠ Review failed_questions_*.json in output/ to retry failed generations\n")
	}
}

// Helper for string repeat
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}
