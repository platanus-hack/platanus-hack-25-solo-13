package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	case "drag_drop_matching":
		return q.validateDragDropMatchingAnswer(userAnswer)
	case "sequencing":
		return q.validateSequencingAnswer(userAnswer)
	case "fill_blanks":
		return q.validateFillBlanksAnswer(userAnswer)
	case "compare_contrast":
		return q.validateCompareContrastAnswer(userAnswer)
	case "open_ended":
		// Open ended requires manual or AI validation
		return false, 0, errors.New("requires manual or AI validation")
	case "concept_map":
		// Concept map requires manual or AI validation
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

// validateDragDropMatchingAnswer validates a drag-drop matching answer
func (q *Question) validateDragDropMatchingAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	// User answer format: { "matches": { "0": 0, "1": 1, "2": 2 } } or { "matches": { "term1": "def1", ... } }
	// Validation formats supported:
	//   - { "correct_matches": { "0": 0, "1": 1, "2": 2 } }
	//   - { "emparejamientos_correctos": { "term1": "def1", "term2": "def2" } }
	userMatches, ok := answer["matches"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("answer must contain 'matches' map")
	}

	// Support both 'correct_matches' and 'emparejamientos_correctos'
	var correctMatches map[string]interface{}
	if cm, ok := validation["correct_matches"].(map[string]interface{}); ok {
		correctMatches = cm
	} else if em, ok := validation["emparejamientos_correctos"].(map[string]interface{}); ok {
		correctMatches = em
	} else {
		return false, 0, errors.New("validation_data must contain 'correct_matches' or 'emparejamientos_correctos' map")
	}

	correctCount := 0
	totalCount := len(correctMatches)

	for termID, correctMatchID := range correctMatches {
		if userMatchID, exists := userMatches[termID]; exists {
			// Compare as strings for flexibility
			if fmt.Sprint(userMatchID) == fmt.Sprint(correctMatchID) {
				correctCount++
			}
		}
	}

	if totalCount == 0 {
		return false, 0, errors.New("no correct matches found in validation_data")
	}

	score := (float64(correctCount) / float64(totalCount)) * 100.0
	isCorrect := score >= 60.0

	return isCorrect, score, nil
}

// validateSequencingAnswer validates a sequencing answer
func (q *Question) validateSequencingAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	// User answer format: { "sequence": ["text1", "text2", "text3"] } or { "sequence": [0, 1, 2] }
	// Validation formats supported:
	//   - { "correct_sequence": [0, 1, 2, 3] }
	//   - { "orden_correcto": ["text1", "text2", "text3"] }
	userSequence, ok := answer["sequence"].([]interface{})
	if !ok {
		return false, 0, errors.New("answer must contain 'sequence' array")
	}

	// Support both 'correct_sequence' and 'orden_correcto'
	var correctSequence []interface{}
	if cs, ok := validation["correct_sequence"].([]interface{}); ok {
		correctSequence = cs
	} else if oc, ok := validation["orden_correcto"].([]interface{}); ok {
		correctSequence = oc
	} else {
		return false, 0, errors.New("validation_data must contain 'correct_sequence' or 'orden_correcto' array")
	}

	if len(userSequence) != len(correctSequence) {
		return false, 0, nil
	}

	correctCount := 0
	for i := range correctSequence {
		if fmt.Sprint(userSequence[i]) == fmt.Sprint(correctSequence[i]) {
			correctCount++
		}
	}

	score := (float64(correctCount) / float64(len(correctSequence))) * 100.0
	isCorrect := score == 100.0 // Sequencing must be perfect

	return isCorrect, score, nil
}

// validateFillBlanksAnswer validates a fill-in-the-blanks answer
func (q *Question) validateFillBlanksAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	// User answer format: { "blanks": { "1": "answer1", "2": "answer2" } }
	// Validation format: { "correct_blanks": { "1": "answer1", "2": "answer2" }, "case_sensitive": false }
	userBlanks, ok := answer["blanks"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("answer must contain 'blanks' map")
	}

	correctBlanks, ok := validation["correct_blanks"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("validation_data must contain 'correct_blanks' map")
	}

	caseSensitive := false
	if cs, ok := validation["case_sensitive"].(bool); ok {
		caseSensitive = cs
	}

	correctCount := 0
	totalCount := len(correctBlanks)

	for blankID, correctAnswer := range correctBlanks {
		if userAnswerVal, exists := userBlanks[blankID]; exists {
			userStr := strings.TrimSpace(fmt.Sprint(userAnswerVal))
			correctStr := strings.TrimSpace(fmt.Sprint(correctAnswer))

			if !caseSensitive {
				userStr = strings.ToLower(userStr)
				correctStr = strings.ToLower(correctStr)
			}

			if userStr == correctStr {
				correctCount++
			}
		}
	}

	score := (float64(correctCount) / float64(totalCount)) * 100.0
	isCorrect := score >= 60.0

	return isCorrect, score, nil
}

// validateCompareContrastAnswer validates a compare-contrast answer
func (q *Question) validateCompareContrastAnswer(userAnswer datatypes.JSON) (bool, float64, error) {
	var answer map[string]interface{}
	if err := json.Unmarshal(userAnswer, &answer); err != nil {
		return false, 0, err
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(q.ValidationData, &validation); err != nil {
		return false, 0, err
	}

	// User answer format: { "classifications": { "0": "A", "1": "B", "2": "both" } }
	// Validation format: { "correct_classifications": { "0": "A", "1": "B", "2": "both" } }
	userClassifications, ok := answer["classifications"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("answer must contain 'classifications' map")
	}

	correctClassifications, ok := validation["correct_classifications"].(map[string]interface{})
	if !ok {
		return false, 0, errors.New("validation_data must contain 'correct_classifications' map")
	}

	correctCount := 0
	totalCount := len(correctClassifications)

	for charID, correctColumn := range correctClassifications {
		if userColumn, exists := userClassifications[charID]; exists {
			if fmt.Sprint(userColumn) == fmt.Sprint(correctColumn) {
				correctCount++
			}
		}
	}

	score := (float64(correctCount) / float64(totalCount)) * 100.0
	isCorrect := score >= 60.0

	return isCorrect, score, nil
}
