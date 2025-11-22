package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/platanus-hack-25/lumera_app/question-generator/generator"
)

func main() {
	log.Println("=== Retry Failed Questions ===")

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

	// Read failed questions file
	failedFile := "output/failed_questions_20251122_155646.json"
	data, err := os.ReadFile(failedFile)
	if err != nil {
		log.Fatalf("❌ Failed to read failed questions file: %v", err)
	}

	var failedQuestions []generator.FailedQuestion
	if err := json.Unmarshal(data, &failedQuestions); err != nil {
		log.Fatalf("❌ Failed to parse failed questions: %v", err)
	}

	log.Printf("Found %d failed questions to retry\n", len(failedQuestions))

	// Get unique OA-Bloom IDs and their question types
	type RetryItem struct {
		OAID       uint
		Tipo       string
		Dificultad int
	}
	retryMap := make(map[string]RetryItem)

	for _, failed := range failedQuestions {
		key := fmt.Sprintf("%d_%s_%d", failed.OABloomObjectiveID, failed.Tipo, failed.Dificultad)
		retryMap[key] = RetryItem{
			OAID:       failed.OABloomObjectiveID,
			Tipo:       failed.Tipo,
			Dificultad: failed.Dificultad,
		}
	}

	log.Printf("Will retry %d unique question attempts\n", len(retryMap))

	// Fetch all objectives
	objectives, err := generator.GetOABloomObjectives()
	if err != nil {
		log.Fatalf("❌ Failed to fetch objectives: %v", err)
	}

	// Create a map for quick lookup
	objMap := make(map[uint]generator.OABloomObjective)
	for _, obj := range objectives {
		objMap[obj.ID] = obj
	}

	// Stats
	stats := &generator.Stats{
		TotalAttempts:   0,
		SuccessCount:    0,
		FailCount:       0,
		FailedQuestions: []generator.FailedQuestion{},
		TypeCounts:      make(map[string]int),
		StartTime:       time.Now(),
	}

	var allQuestions []generator.Question
	var stillFailed []generator.FailedQuestion

	// Process each unique retry item
	i := 1
	for _, retry := range retryMap {
		objective, ok := objMap[retry.OAID]
		if !ok {
			log.Printf("[%d/%d] ❌ Objective #%d not found", i, len(retryMap), retry.OAID)
			stillFailed = append(stillFailed, generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              "objective not found",
				Timestamp:          time.Now(),
			})
			i++
			continue
		}

		log.Printf("\n[%d/%d] Retrying OA-Bloom #%d: %s (%s, difficulty %d)",
			i, len(retryMap), retry.OAID, objective.OATitulo, retry.Tipo, retry.Dificultad)

		stats.TotalAttempts++

		// Generate the specific question
		result, err := generator.GenerateQuestion(retry.Tipo, objective, retry.Dificultad)
		if err != nil {
			log.Printf("✗ Failed to generate question: %v", err)
			stats.AddFail(generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			stillFailed = append(stillFailed, generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			i++
			continue
		}

		// Build question struct using unexported function - we need to duplicate the logic
		questionDataMap, ok := result["question_data"].(map[string]interface{})
		if !ok {
			err := fmt.Errorf("invalid question_data format")
			log.Printf("✗ Failed to build question struct: %v", err)
			stats.AddFail(generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			stillFailed = append(stillFailed, generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			i++
			continue
		}

		questionDataJSON, err := json.Marshal(questionDataMap)
		if err != nil {
			log.Printf("✗ Failed to marshal question_data: %v", err)
			stats.AddFail(generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			stillFailed = append(stillFailed, generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			i++
			continue
		}

		validationDataMap, ok := result["validation_data"].(map[string]interface{})
		if !ok {
			err := fmt.Errorf("invalid validation_data format")
			log.Printf("✗ Failed to build question struct: %v", err)
			stats.AddFail(generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			stillFailed = append(stillFailed, generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			i++
			continue
		}

		validationDataJSON, err := json.Marshal(validationDataMap)
		if err != nil {
			log.Printf("✗ Failed to marshal validation_data: %v", err)
			stats.AddFail(generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			stillFailed = append(stillFailed, generator.FailedQuestion{
				OABloomObjectiveID: retry.OAID,
				Tipo:               retry.Tipo,
				Dificultad:         retry.Dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			i++
			continue
		}

		// Extract tags
		var tags []string
		if tagsInterface, ok := result["tags"]; ok {
			if tagsList, ok := tagsInterface.([]interface{}); ok {
				for _, tag := range tagsList {
					if tagStr, ok := tag.(string); ok {
						tags = append(tags, tagStr)
					}
				}
			}
		}

		if len(tags) == 0 {
			tags = []string{"auto-generated"}
		}

		question := generator.Question{
			OABloomObjectiveID: retry.OAID,
			Tipo:               retry.Tipo,
			TipoUso:            "all",
			QuestionData:       questionDataJSON,
			ValidationData:     validationDataJSON,
			DificultadRelativa: retry.Dificultad,
			Tags:               tags,
		}

		allQuestions = append(allQuestions, question)
		stats.AddSuccess(retry.Tipo)
		log.Printf("✓ Successfully generated %s question (difficulty %d)", retry.Tipo, retry.Dificultad)
		i++
	}

	// Save all questions to database
	if len(allQuestions) > 0 {
		log.Printf("\n💾 Saving %d questions to database...", len(allQuestions))
		if err := generator.InsertQuestions(allQuestions); err != nil {
			log.Printf("❌ Failed to insert questions: %v", err)
		} else {
			log.Printf("✓ Successfully saved %d questions", len(allQuestions))
		}
	}

	// Save still-failed questions
	if len(stillFailed) > 0 {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("output/still_failed_%s.json", timestamp)
		data, _ := json.MarshalIndent(stillFailed, "", "  ")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			log.Printf("⚠ Failed to save still-failed questions: %v", err)
		} else {
			log.Printf("\n📝 Saved %d still-failed questions to %s", len(stillFailed), filename)
		}
	}

	// Print stats
	stats.EndTime = time.Now()
	duration := stats.EndTime.Sub(stats.StartTime)

	fmt.Println("\n============================================================")
	fmt.Println("📊 RETRY STATISTICS")
	fmt.Println("============================================================")
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

	fmt.Println("============================================================")
}
