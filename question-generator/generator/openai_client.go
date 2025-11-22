package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

// InitOpenAI inicializa el cliente de OpenAI
func InitOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY no está configurada en el archivo .env")
	}
	openaiClient = openai.NewClient(apiKey)
	log.Println("✓ OpenAI client initialized")
}

// GenerateQuestion genera una pregunta usando OpenAI
func GenerateQuestion(questionType string, objective OABloomObjective, dificultad int) (map[string]interface{}, error) {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	timeoutStr := os.Getenv("OPENAI_TIMEOUT_SECONDS")
	timeout := 45
	if timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil {
			timeout = t
		}
	}

	maxRetriesStr := os.Getenv("OPENAI_MAX_RETRIES")
	maxRetries := 3
	if maxRetriesStr != "" {
		if r, err := strconv.Atoi(maxRetriesStr); err == nil {
			maxRetries = r
		}
	}

	prompt := BuildPrompt(questionType, objective, dificultad)

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		resp, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Eres un experto en diseño de evaluaciones educativas para el currículo chileno de enseñanza media.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   2000,
		})

		if err != nil {
			lastError = err
			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				log.Printf("⚠ OpenAI error (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, waitTime)
				time.Sleep(waitTime)
				continue
			}
			return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}

		if len(resp.Choices) == 0 {
			lastError = fmt.Errorf("no response from OpenAI")
			continue
		}

		content := resp.Choices[0].Message.Content

		// Parse JSON response
		var result map[string]interface{}
		err = json.Unmarshal([]byte(content), &result)
		if err != nil {
			lastError = fmt.Errorf("failed to parse OpenAI response as JSON: %w\nContent: %s", err, content)
			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				log.Printf("⚠ JSON parse error (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, waitTime)
				time.Sleep(waitTime)
				continue
			}
			return nil, lastError
		}

		// Validate structure
		if _, ok := result["question_data"]; !ok {
			lastError = fmt.Errorf("missing question_data in response")
			if attempt < maxRetries {
				log.Printf("⚠ Invalid structure (attempt %d/%d): missing question_data. Retrying...", attempt, maxRetries)
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastError
		}

		if _, ok := result["validation_data"]; !ok {
			lastError = fmt.Errorf("missing validation_data in response")
			if attempt < maxRetries {
				log.Printf("⚠ Invalid structure (attempt %d/%d): missing validation_data. Retrying...", attempt, maxRetries)
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastError
		}

		// Success
		log.Printf("✓ Generated %s question (difficulty %d) for OA-Bloom #%d", questionType, dificultad, objective.ID)
		return result, nil
	}

	return nil, lastError
}
