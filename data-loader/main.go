package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/platanus-hack-25/lumera_app/data-loader/loader"
)

func main() {
	// Flags
	inputFile := flag.String("input", "input/oas.csv", "Ruta al archivo CSV de entrada")
	flag.Parse()

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Validar variables requeridas
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ Error: OPENAI_API_KEY no está configurada en .env")
	}

	// Configuración OpenAI
	model := getEnvOrDefault("OPENAI_MODEL", "gpt-4o-mini")
	timeoutSeconds := getEnvOrDefaultInt("OPENAI_TIMEOUT_SECONDS", 30)
	maxRetries := getEnvOrDefaultInt("OPENAI_MAX_RETRIES", 3)

	// Configuración BD
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "admin")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := getEnvOrDefault("DB_NAME", "hackathon")
	dbSSLMode := getEnvOrDefault("DB_SSLMODE", "disable")

	if dbPassword == "" {
		log.Fatal("❌ Error: DB_PASSWORD no está configurada en .env")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Inicializar clientes
	log.Println("🔧 Inicializando clientes...")
	openaiClient := loader.NewOpenAIClient(apiKey, model, time.Duration(timeoutSeconds)*time.Second, maxRetries)
	dbWriter, err := loader.NewDBWriter(dsn)
	if err != nil {
		log.Fatalf("❌ Error conectando a BD: %v", err)
	}
	defer dbWriter.Close()

	// Leer CSV
	log.Printf("📂 Leyendo archivo CSV: %s\n", *inputFile)
	records, err := loader.ReadCSV(*inputFile)
	if err != nil {
		log.Fatalf("❌ Error leyendo CSV: %v", err)
	}
	log.Printf("✅ Se encontraron %d OAs para cargar\n\n", len(records))

	// Procesar cada OA
	var successCount, failCount int
	var failedOAs []loader.FailedOA

	for i, record := range records {
		log.Printf("[%d/%d] Procesando OA: %s (Materia: %s, Curso: %s)\n",
			i+1, len(records), truncateString(record.Objetivo, 60), record.Materia, record.Curso)

		// 1. Generar objetivos de Bloom con OpenAI
		log.Printf("  🤖 Generando niveles de Bloom con OpenAI...")
		bloomObjectives, err := openaiClient.GenerateBloomObjectives(record.Objetivo)
		if err != nil {
			log.Printf("  ❌ Error llamando a OpenAI: %v\n", err)
			failedOAs = append(failedOAs, loader.FailedOA{
				Materia:   record.Materia,
				Curso:     record.Curso,
				Objetivo:  record.Objetivo,
				Error:     err.Error(),
				Timestamp: time.Now(),
			})
			failCount++
			continue
		}
		log.Printf("  ✅ Generados %d niveles de Bloom\n", len(bloomObjectives))

		// 2. Guardar en BD
		log.Printf("  💾 Guardando en base de datos...")
		if err := dbWriter.SaveOA(record, bloomObjectives); err != nil {
			log.Printf("  ❌ Error guardando en BD: %v\n", err)
			failedOAs = append(failedOAs, loader.FailedOA{
				Materia:   record.Materia,
				Curso:     record.Curso,
				Objetivo:  record.Objetivo,
				Error:     fmt.Sprintf("Error BD: %v", err),
				Timestamp: time.Now(),
			})
			failCount++
			continue
		}
		log.Printf("  ✅ OA guardado exitosamente\n\n")
		successCount++
	}

	// Guardar OAs fallidos si existen
	if len(failedOAs) > 0 {
		outputDir := "output"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Printf("⚠️  Error creando carpeta output: %v\n", err)
		} else {
			outputFile := filepath.Join(outputDir, fmt.Sprintf("failed_oas_%s.json", time.Now().Format("20060102_150405")))
			data, _ := json.MarshalIndent(failedOAs, "", "  ")
			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				log.Printf("⚠️  Error guardando OAs fallidos: %v\n", err)
			} else {
				log.Printf("📄 OAs fallidos guardados en: %s\n", outputFile)
			}
		}
	}

	// Reporte final
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("📊 REPORTE FINAL")
	log.Println(strings.Repeat("=", 60))
	log.Printf("Total OAs procesados:  %d\n", len(records))
	log.Printf("✅ Éxitos:              %d (%.1f%%)\n", successCount, float64(successCount)/float64(len(records))*100)
	log.Printf("❌ Fallos:              %d (%.1f%%)\n", failCount, float64(failCount)/float64(len(records))*100)
	if failCount > 0 {
		log.Printf("\n💡 Revisa los OAs fallidos en: output/failed_oas_*.json\n")
	}
	log.Println(strings.Repeat("=", 60))
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
