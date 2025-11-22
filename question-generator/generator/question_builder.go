package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

// QuestionTypeDistribution define qué tipos de pregunta generar para cada nivel de Bloom
type QuestionTypeDistribution struct {
	Types       []string
	Counts      []int
	Dificultad  []int // Dificultad relativa para cada tipo
}

// GetQuestionTypesForBloomLevel retorna la distribución de tipos de pregunta según nivel de Bloom
func GetQuestionTypesForBloomLevel(bloomLevel int) QuestionTypeDistribution {
	switch bloomLevel {
	case 1: // Recordar
		return QuestionTypeDistribution{
			Types:      []string{"multiple_choice", "multiple_choice", "multiple_choice", "true_false", "fill_blanks"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{1, 2, 3, 1, 2},
		}
	case 2: // Comprender
		return QuestionTypeDistribution{
			Types:      []string{"multiple_choice", "multiple_choice", "true_false", "drag_drop_matching", "sequencing"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{2, 3, 2, 2, 3},
		}
	case 3: // Aplicar
		return QuestionTypeDistribution{
			Types:      []string{"multiple_choice", "multiple_choice", "drag_drop_matching", "sequencing", "open_ended"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{3, 4, 3, 3, 3},
		}
	case 4: // Analizar
		return QuestionTypeDistribution{
			Types:      []string{"compare_contrast", "compare_contrast", "open_ended", "concept_map", "concept_map"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{3, 4, 4, 3, 4},
		}
	case 5: // Evaluar
		return QuestionTypeDistribution{
			Types:      []string{"criteria_evaluation", "criteria_evaluation", "criteria_evaluation", "open_ended", "open_ended"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{4, 5, 5, 4, 5},
		}
	case 6: // Crear
		return QuestionTypeDistribution{
			Types:      []string{"open_ended", "open_ended", "concept_map", "concept_map", "open_ended"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{4, 5, 4, 5, 5},
		}
	default:
		// Fallback genérico
		return QuestionTypeDistribution{
			Types:      []string{"multiple_choice", "true_false", "fill_blanks", "open_ended", "sequencing"},
			Counts:     []int{1, 1, 1, 1, 1},
			Dificultad: []int{2, 2, 2, 3, 3},
		}
	}
}

// GenerateQuestionsForObjective genera todas las preguntas para un OA-Bloom objective
func GenerateQuestionsForObjective(objective OABloomObjective, stats *Stats) ([]Question, []FailedQuestion) {
	distribution := GetQuestionTypesForBloomLevel(objective.BloomLevelNumero)

	var questions []Question
	var failed []FailedQuestion

	for i, questionType := range distribution.Types {
		dificultad := distribution.Dificultad[i]

		log.Printf("→ Generating %s (difficulty %d) for OA-Bloom #%d (%s - Bloom %d)",
			questionType, dificultad, objective.ID, objective.MateriaNombre, objective.BloomLevelNumero)

		stats.TotalAttempts++

		// Llamar a OpenAI
		result, err := GenerateQuestion(questionType, objective, dificultad)
		if err != nil {
			log.Printf("✗ Failed to generate question: %v", err)
			stats.AddFail(FailedQuestion{
				OABloomObjectiveID: objective.ID,
				Tipo:               questionType,
				Dificultad:         dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			failed = append(failed, FailedQuestion{
				OABloomObjectiveID: objective.ID,
				Tipo:               questionType,
				Dificultad:         dificultad,
				Error:              err.Error(),
				Timestamp:          time.Now(),
			})
			continue
		}

		// Construir Question struct
		question, err := buildQuestionStruct(objective.ID, questionType, result, dificultad)
		if err != nil {
			log.Printf("✗ Failed to build question struct: %v", err)
			stats.AddFail(FailedQuestion{
				OABloomObjectiveID: objective.ID,
				Tipo:               questionType,
				Dificultad:         dificultad,
				Error:              fmt.Sprintf("build error: %v", err),
				Timestamp:          time.Now(),
			})
			failed = append(failed, FailedQuestion{
				OABloomObjectiveID: objective.ID,
				Tipo:               questionType,
				Dificultad:         dificultad,
				Error:              fmt.Sprintf("build error: %v", err),
				Timestamp:          time.Now(),
			})
			continue
		}

		questions = append(questions, question)
		stats.AddSuccess(questionType)
		log.Printf("✓ Successfully created %s question (difficulty %d)", questionType, dificultad)
	}

	return questions, failed
}

// buildQuestionStruct convierte el resultado de OpenAI en un struct Question
func buildQuestionStruct(oaBloomObjectiveID uint, questionType string, result map[string]interface{}, dificultad int) (Question, error) {
	// Extract question_data
	questionDataMap, ok := result["question_data"].(map[string]interface{})
	if !ok {
		return Question{}, fmt.Errorf("invalid question_data format")
	}
	questionDataJSON, err := json.Marshal(questionDataMap)
	if err != nil {
		return Question{}, fmt.Errorf("failed to marshal question_data: %w", err)
	}

	// Extract validation_data
	validationDataMap, ok := result["validation_data"].(map[string]interface{})
	if !ok {
		return Question{}, fmt.Errorf("invalid validation_data format")
	}
	validationDataJSON, err := json.Marshal(validationDataMap)
	if err != nil {
		return Question{}, fmt.Errorf("failed to marshal validation_data: %w", err)
	}

	// Extract tags
	var tags pq.StringArray
	if tagsInterface, ok := result["tags"]; ok {
		if tagsList, ok := tagsInterface.([]interface{}); ok {
			for _, tag := range tagsList {
				if tagStr, ok := tag.(string); ok {
					tags = append(tags, tagStr)
				}
			}
		}
	}

	// Default tags if none provided
	if len(tags) == 0 {
		tags = pq.StringArray{"auto-generated"}
	}

	// Determine tipo_uso based on question type
	// Most questions can be used for diagnostico, practica, and evaluacion
	// But some specialized types are better for certain uses
	tipoUso := "all" // diagnostico, practica, evaluacion

	return Question{
		OABloomObjectiveID: oaBloomObjectiveID,
		Tipo:               questionType,
		TipoUso:            tipoUso,
		QuestionData:       questionDataJSON,
		ValidationData:     validationDataJSON,
		DificultadRelativa: dificultad,
		Tags:               tags,
	}, nil
}
