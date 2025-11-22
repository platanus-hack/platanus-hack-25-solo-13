package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// QuestionType represents a registered question type in the catalog
type QuestionType struct {
	Tipo           string         `json:"tipo" gorm:"primaryKey"`
	NombreDisplay  string         `json:"nombre_display" gorm:"size:100;not null"`
	Descripcion    string         `json:"descripcion" gorm:"type:text"`
	SchemaExample  datatypes.JSON `json:"schema_example,omitempty" gorm:"type:jsonb"`
	Activo         bool           `json:"activo" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
}

// TableName overrides the default table name
func (QuestionType) TableName() string {
	return "question_types"
}

// Question represents a flexible question in the question bank
type Question struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	OABloomObjectiveID   uint           `json:"oa_bloom_objective_id" gorm:"not null"`
	Tipo                 string         `json:"tipo" gorm:"size:30;not null"`
	TipoUso              string         `json:"tipo_uso" gorm:"size:20;not null;default:all"`
	QuestionData         datatypes.JSON `json:"question_data" gorm:"type:jsonb;not null"`
	ValidationData       datatypes.JSON `json:"validation_data" gorm:"type:jsonb;not null"`
	DificultadRelativa   int            `json:"dificultad_relativa" gorm:"default:3"`
	VecesUsada           int            `json:"veces_usada" gorm:"default:0"`
	Activa               bool           `json:"activa" gorm:"default:true"`
	Tags                 pq.StringArray `json:"tags" gorm:"type:text[]"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`

	// Relationships
	OABloomObjective OABloomObjective `json:"oa_bloom_objective,omitempty" gorm:"foreignKey:OABloomObjectiveID"`
	QuestionType     QuestionType     `json:"question_type,omitempty" gorm:"foreignKey:Tipo;references:Tipo"`
}

// TableName overrides the default table name
func (Question) TableName() string {
	return "questions"
}

// Validate validates the question structure based on its type
func (q *Question) Validate() error {
	switch q.Tipo {
	case "multiple_choice":
		return q.validateMultipleChoice()
	case "true_false":
		return q.validateTrueFalse()
	case "fill_blanks":
		return q.validateFillBlanks()
	case "drag_drop_matching":
		return q.validateDragDropMatching()
	case "sequencing":
		return q.validateSequencing()
	case "compare_contrast":
		return q.validateCompareContrast()
	case "open_ended":
		return q.validateOpenEnded()
	case "criteria_evaluation":
		return q.validateCriteriaEvaluation()
	case "concept_map":
		return q.validateConceptMap()
	default:
		return errors.New("unknown question type: " + q.Tipo)
	}
}

// validateMultipleChoice validates multiple choice question structure
func (q *Question) validateMultipleChoice() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	// Required fields in question_data
	if _, ok := questionData["pregunta"]; !ok {
		return errors.New("question_data must contain 'pregunta'")
	}
	if _, ok := questionData["opciones"]; !ok {
		return errors.New("question_data must contain 'opciones'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	// Required fields in validation_data
	if _, ok := validationData["respuesta_correcta"]; !ok {
		return errors.New("validation_data must contain 'respuesta_correcta'")
	}

	return nil
}

// validateTrueFalse validates true/false question structure
func (q *Question) validateTrueFalse() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["statement"]; !ok {
		return errors.New("question_data must contain 'statement'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["correct_answer"]; !ok {
		return errors.New("validation_data must contain 'correct_answer'")
	}

	return nil
}

// validateFillBlanks validates fill-in-the-blanks question structure
func (q *Question) validateFillBlanks() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["text"]; !ok {
		return errors.New("question_data must contain 'text'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["blanks"]; !ok {
		return errors.New("validation_data must contain 'blanks'")
	}

	return nil
}

// validateDragDropMatching validates drag-drop matching question structure
func (q *Question) validateDragDropMatching() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["pairs"]; !ok {
		return errors.New("question_data must contain 'pairs'")
	}

	return nil
}

// validateSequencing validates sequencing question structure
func (q *Question) validateSequencing() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["items"]; !ok {
		return errors.New("question_data must contain 'items'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["correct_order"]; !ok {
		return errors.New("validation_data must contain 'correct_order'")
	}

	return nil
}

// validateCompareContrast validates compare/contrast question structure
func (q *Question) validateCompareContrast() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["characteristics"]; !ok {
		return errors.New("question_data must contain 'characteristics'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["correct_columns"]; !ok {
		return errors.New("validation_data must contain 'correct_columns'")
	}

	return nil
}

// validateOpenEnded validates open-ended question structure
func (q *Question) validateOpenEnded() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["prompt"]; !ok {
		return errors.New("question_data must contain 'prompt'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["rubric"]; !ok {
		return errors.New("validation_data must contain 'rubric'")
	}

	return nil
}

// validateCriteriaEvaluation validates criteria evaluation question structure
func (q *Question) validateCriteriaEvaluation() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["criteria"]; !ok {
		return errors.New("question_data must contain 'criteria'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["expected_ratings"]; !ok {
		return errors.New("validation_data must contain 'expected_ratings'")
	}

	return nil
}

// validateConceptMap validates concept map question structure
func (q *Question) validateConceptMap() error {
	var questionData map[string]interface{}
	if err := json.Unmarshal(q.QuestionData, &questionData); err != nil {
		return errors.New("invalid question_data JSON")
	}

	if _, ok := questionData["required_concepts"]; !ok {
		return errors.New("question_data must contain 'required_concepts'")
	}

	var validationData map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validationData); err != nil {
		return errors.New("invalid validation_data JSON")
	}

	if _, ok := validationData["suggested_connections"]; !ok {
		return errors.New("validation_data must contain 'suggested_connections'")
	}

	return nil
}

// ValidateAnswer validates a user's answer against the validation_data
func (q *Question) ValidateAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	switch q.Tipo {
	case "multiple_choice":
		return q.validateMultipleChoiceAnswer(userAnswer)
	case "true_false":
		return q.validateTrueFalseAnswer(userAnswer)
	case "criteria_evaluation":
		return q.validateCriteriaEvaluationAnswer(userAnswer)
	case "open_ended":
		// Open ended requires manual or AI validation
		return false, 0, errors.New("requires manual or AI validation")
	// Add other types as needed
	default:
		return false, 0, errors.New("validation not implemented for this question type")
	}
}

// validateMultipleChoiceAnswer validates a multiple choice answer
func (q *Question) validateMultipleChoiceAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	userChoice, ok := answer["selected"]
	if !ok {
		return false, 0, errors.New("answer must contain 'selected' field")
	}

	correctAnswer := validation["respuesta_correcta"]
	isCorrect := userChoice == correctAnswer

	score := 0.0
	if isCorrect {
		score = 100.0
	}

	return isCorrect, score, nil
}

// validateTrueFalseAnswer validates a true/false answer
func (q *Question) validateTrueFalseAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	userChoice, ok := answer["answer"]
	if !ok {
		return false, 0, errors.New("answer must contain 'answer' field")
	}

	correctAnswer := validation["correct_answer"]
	isCorrect := userChoice == correctAnswer

	score := 0.0
	if isCorrect {
		score = 100.0
	}

	return isCorrect, score, nil
}

// validateCriteriaEvaluationAnswer validates a criteria evaluation answer with tolerance
func (q *Question) validateCriteriaEvaluationAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	ratings, ok := answer["ratings"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("answer must contain 'ratings' map")
	}

	expectedRatings, ok := validation["expected_ratings"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("validation_data must contain 'expected_ratings' map")
	}

	tolerance := 1.0
	if t, ok := validation["tolerance"].(float64); ok {
		tolerance = t
	}

	totalScore := 0.0
	criteriaCount := 0

	for criteriaID, expectedRating := range expectedRatings {
		userRating, ok := ratings[criteriaID].(float64)
		if !ok {
			continue
		}

		expected := expectedRating.(float64)
		difference := abs(userRating - expected)

		var criteriaScore float64
		if difference == 0 {
			criteriaScore = 100.0
		} else if difference <= tolerance {
			criteriaScore = 60.0
		} else if difference <= tolerance*2 {
			criteriaScore = 30.0
		} else {
			criteriaScore = 0.0
		}

		totalScore += criteriaScore
		criteriaCount++
	}

	if criteriaCount == 0 {
		return false, 0, errors.New("no valid criteria ratings found")
	}

	finalScore := totalScore / float64(criteriaCount)
	isCorrect := finalScore >= 60.0

	return isCorrect, finalScore, nil
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
