package generator

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// slugify convierte un nombre en un slug para archivo
func slugify(text string) string {
	// Convertir a minúsculas y reemplazar espacios con guiones
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "á", "a")
	text = strings.ReplaceAll(text, "é", "e")
	text = strings.ReplaceAll(text, "í", "i")
	text = strings.ReplaceAll(text, "ó", "o")
	text = strings.ReplaceAll(text, "ú", "u")
	text = strings.ReplaceAll(text, "ñ", "n")

	// Eliminar caracteres no alfanuméricos excepto guiones
	var result strings.Builder
	for _, char := range text {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// GenerateAvatarImage genera una imagen de avatar usando DALL-E-3
func GenerateAvatarImage(avatar Avatar) (string, error) {
	timeoutStr := os.Getenv("OPENAI_TIMEOUT_SECONDS")
	timeout := 60
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

	// Construir prompt completo con lineamientos de estilo
	fullPrompt := buildFullPrompt(avatar)

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		// Generar imagen con DALL-E-3
		resp, err := openaiClient.CreateImage(ctx, openai.ImageRequest{
			Prompt:         fullPrompt,
			Model:          "dall-e-3",
			N:              1,
			Size:           "1024x1024",
			Quality:        "standard",
			ResponseFormat: "url",
		})

		if err != nil {
			lastError = err
			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				log.Printf("⚠ DALL-E error (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, waitTime)
				time.Sleep(waitTime)
				continue
			}
			return "", fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}

		if len(resp.Data) == 0 {
			lastError = fmt.Errorf("no image generated")
			continue
		}

		imageURL := resp.Data[0].URL

		// Descargar imagen y guardar localmente
		localPath, err := downloadImage(imageURL, avatar.Nombre)
		if err != nil {
			lastError = fmt.Errorf("failed to download image: %w", err)
			if attempt < maxRetries {
				waitTime := time.Duration(attempt*2) * time.Second
				log.Printf("⚠ Download error (attempt %d/%d): %v. Retrying in %v...", attempt, maxRetries, err, waitTime)
				time.Sleep(waitTime)
				continue
			}
			return "", lastError
		}

		log.Printf("✓ Generated avatar image for '%s' (tier %d⭐)", avatar.Nombre, avatar.Tier)
		return localPath, nil
	}

	return "", lastError
}

// buildFullPrompt construye el prompt completo con lineamientos de estilo
func buildFullPrompt(avatar Avatar) string {
	// Lineamientos base
	baseStyle := "friendly cartoon mascot character, educational theme, clean simple design, white background, professional quality, vibrant colors, appealing to middle school students (14-18 years old)"

	// Lineamientos por tier
	tierEnhancements := map[int]string{
		1: "simple design, basic shapes, primary colors, minimalist style",
		2: "more detailed, colorful accessories, expressive features",
		3: "detailed design, thematic accessories, rich color palette, dynamic pose",
		4: "highly detailed, magical elements, special effects (subtle glow, sparkles), heroic pose, unique design",
		5: "extremely detailed, prominent magical effects (glowing aura, energy effects), epic heroic pose, legendary appearance, premium quality, striking visual impact",
	}

	tierStyle := tierEnhancements[avatar.Tier]
	if tierStyle == "" {
		tierStyle = tierEnhancements[1]
	}

	// Combinar todo
	fullPrompt := fmt.Sprintf("%s. %s. %s", avatar.DallePrompt, tierStyle, baseStyle)

	return fullPrompt
}

// downloadImage descarga una imagen desde URL y la guarda localmente
func downloadImage(url, nombre string) (string, error) {
	// Crear directorio si no existe
	outputDir := "output/images"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generar nombre de archivo
	slug := slugify(nombre)
	filename := fmt.Sprintf("%s.png", slug)
	filepath := filepath.Join(outputDir, filename)

	// Descargar imagen
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Guardar imagen
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write image: %w", err)
	}

	// Retornar path relativo para la base de datos
	return fmt.Sprintf("/static/avatares/%s", filename), nil
}
