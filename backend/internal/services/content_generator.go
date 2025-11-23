package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/platanus-hack-25/lumera_app/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

// InitOpenAI inicializa el cliente de OpenAI
func InitOpenAI() error {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY no está configurada")
	}
	openaiClient = openai.NewClient(apiKey)
	log.Println("✓ OpenAI client initialized for content generation")
	return nil
}

// cleanMarkdownJSON removes markdown code block markers from JSON response
func cleanMarkdownJSON(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
	}
	if strings.HasSuffix(content, "```") {
		content = strings.TrimSuffix(content, "```")
	}
	return strings.TrimSpace(content)
}

// getEnvInt obtiene una variable de entorno como entero con valor por defecto
func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvString obtiene una variable de entorno como string con valor por defecto
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// LearningPlanStructure representa la estructura del plan retornado por OpenAI
type LearningPlanStructure struct {
	Titulo      string                    `json:"titulo"`
	Descripcion string                    `json:"descripcion"`
	Componentes []ComponentStructure       `json:"componentes"`
}

// ComponentStructure representa la estructura de un componente en el plan
type ComponentStructure struct {
	Tipo                string `json:"tipo"`
	ObjetivoEspecifico  string `json:"objetivo_especifico"`
	TiempoEstimadoMin   int    `json:"tiempo_estimado_minutos"`
}

// OAContext contiene el contexto educativo del OA para los prompts
type OAContext struct {
	MateriaNombre       string
	MateriaDescripcion  string
	CursoNombre         string
	OATitulo            string
	OADescripcion       string
	BloomLevelNombre    string
	BloomLevelNumero    int
	BloomDescripcion    string
	ObjetivoEspecifico  string
	IndicadoresLogro    []string

	// Student profile personalization
	InteresesPersonales  []string
	ProfesionSoñada      string
	FormatoPreferido     string
	TipoActividad        []string
	CanalPreferido       string
}

// GenerateLearningPlanStructure genera la estructura del plan de aprendizaje
func GenerateLearningPlanStructure(oaContext OAContext) (*LearningPlanStructure, error) {
	model := getEnvString("OPENAI_MODEL", "gpt-4o-mini")
	timeout := getEnvInt("OPENAI_TIMEOUT_SECONDS", 45)
	maxRetries := getEnvInt("OPENAI_MAX_RETRIES", 3)

	prompt := buildPlanStructurePrompt(oaContext)

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		resp, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Eres un experto en diseño instruccional para el currículo chileno de enseñanza media. Tu tarea es crear planes de aprendizaje personalizados y efectivos.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.8,
			MaxTokens:   1500,
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
		content = cleanMarkdownJSON(content)

		var result LearningPlanStructure
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

		// Validar estructura
		if result.Titulo == "" {
			lastError = fmt.Errorf("missing titulo in response")
			if attempt < maxRetries {
				log.Printf("⚠ Invalid structure (attempt %d/%d): missing titulo. Retrying...", attempt, maxRetries)
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastError
		}

		if len(result.Componentes) == 0 {
			lastError = fmt.Errorf("no components in learning plan")
			if attempt < maxRetries {
				log.Printf("⚠ Invalid structure (attempt %d/%d): no components. Retrying...", attempt, maxRetries)
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastError
		}

		// Validar que todos los componentes tengan tipo válido
		for i, comp := range result.Componentes {
			if !models.IsValidComponentType(comp.Tipo) {
				lastError = fmt.Errorf("invalid component type at index %d: %s", i, comp.Tipo)
				if attempt < maxRetries {
					log.Printf("⚠ Invalid component type (attempt %d/%d): %s. Retrying...", attempt, maxRetries, comp.Tipo)
					time.Sleep(2 * time.Second)
					continue
				}
				return nil, lastError
			}
		}

		log.Printf("✓ Generated learning plan structure: %s (%d components)", result.Titulo, len(result.Componentes))
		return &result, nil
	}

	return nil, lastError
}

// GenerateComponentContent genera el contenido (props) de un componente específico
func GenerateComponentContent(componentType string, oaContext OAContext, componentObjective string) (map[string]interface{}, error) {
	if !models.IsValidComponentType(componentType) {
		return nil, fmt.Errorf("invalid component type: %s", componentType)
	}

	model := getEnvString("OPENAI_MODEL", "gpt-4o-mini")
	timeout := getEnvInt("OPENAI_TIMEOUT_SECONDS", 60) // Más tiempo para contenido detallado
	maxRetries := getEnvInt("OPENAI_MAX_RETRIES", 3)

	prompt := buildComponentPrompt(componentType, oaContext, componentObjective)

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		resp, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Eres un experto en diseño de contenido educativo para el currículo chileno de enseñanza media. Creas contenido pedagógico claro, preciso y adaptado al nivel del estudiante.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   3000,
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
		content = cleanMarkdownJSON(content)

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

		// Validación básica según el tipo de componente
		if err := validateComponentContent(componentType, result); err != nil {
			lastError = err
			if attempt < maxRetries {
				log.Printf("⚠ Invalid content structure (attempt %d/%d): %v. Retrying...", attempt, maxRetries, err)
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, lastError
		}

		log.Printf("✓ Generated content for component type: %s", componentType)
		return result, nil
	}

	return nil, lastError
}

// validateComponentContent valida que el contenido de ExplainAndExploreSlide tenga los campos necesarios
func validateComponentContent(componentType string, content map[string]interface{}) error {
	// Validar que el componente sea ExplainAndExploreSlide
	if componentType != models.ComponentTipoExplainAndExplore {
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	// Campos requeridos para ExplainAndExploreSlide
	if _, ok := content["titulo"]; !ok {
		return fmt.Errorf("missing required field: titulo")
	}

	if _, ok := content["bloques"]; !ok {
		return fmt.Errorf("missing required field: bloques")
	}

	// Validar que bloques sea un array
	bloques, ok := content["bloques"].([]interface{})
	if !ok {
		return fmt.Errorf("bloques must be an array")
	}

	if len(bloques) == 0 {
		return fmt.Errorf("bloques array cannot be empty")
	}

	// Validar tipos de bloques permitidos
	validBlockTypes := map[string]bool{
		"texto":       true,
		"ejemplo":     true,
		"definicion":  true,
		"nota":        true,
		"ejercicio":   true,
		"resumen":     true,
		"comparacion": true,
	}

	for i, bloque := range bloques {
		bloqueMap, ok := bloque.(map[string]interface{})
		if !ok {
			return fmt.Errorf("bloque at index %d is not an object", i)
		}

		tipo, ok := bloqueMap["tipo"].(string)
		if !ok {
			return fmt.Errorf("bloque at index %d missing 'tipo' field", i)
		}

		if !validBlockTypes[tipo] {
			return fmt.Errorf("bloque at index %d has invalid tipo: %s", i, tipo)
		}
	}

	return nil
}
