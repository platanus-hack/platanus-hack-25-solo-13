package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"gorm.io/datatypes"
)

// GetQuestionTypes godoc
// @Summary Get all question types
// @Description Retrieve all registered question types from catalog
// @Tags Questions
// @Produce json
// @Success 200 {array} models.QuestionType
// @Failure 500 {object} map[string]interface{}
// @Router /api/question-types [get]
func GetQuestionTypes(w http.ResponseWriter, r *http.Request) {
	var questionTypes []models.QuestionType

	if err := db.DB.Where("activo = ?", true).Find(&questionTypes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionTypes)
}

// GetQuestions godoc
// @Summary Get all questions
// @Description Retrieve questions with optional filters
// @Tags Questions
// @Produce json
// @Param tipo query string false "Filter by question type"
// @Param tipo_uso query string false "Filter by usage type"
// @Param oa_bloom_objective_id query int false "Filter by OA Bloom objective"
// @Param activa query boolean false "Filter by active status"
// @Success 200 {array} models.Question
// @Failure 500 {object} map[string]interface{}
// @Router /api/questions [get]
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	var questions []models.Question
	query := db.DB.Preload("OABloomObjective").Preload("QuestionType")

	// Filters
	if tipo := r.URL.Query().Get("tipo"); tipo != "" {
		query = query.Where("tipo = ?", tipo)
	}
	if tipoUso := r.URL.Query().Get("tipo_uso"); tipoUso != "" {
		query = query.Where("tipo_uso = ?", tipoUso)
	}
	if oaBloomID := r.URL.Query().Get("oa_bloom_objective_id"); oaBloomID != "" {
		query = query.Where("oa_bloom_objective_id = ?", oaBloomID)
	}
	if activa := r.URL.Query().Get("activa"); activa != "" {
		query = query.Where("activa = ?", activa == "true")
	}

	if err := query.Find(&questions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// GetQuestion godoc
// @Summary Get question by ID
// @Description Retrieve a specific question WITHOUT validation_data (security)
// @Tags Questions
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/questions/{id} [get]
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var question models.Question

	if err := db.DB.Preload("OABloomObjective").Preload("QuestionType").First(&question, id).Error; err != nil {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	// Return question WITHOUT validation_data for security
	response := map[string]interface{}{
		"id":                   question.ID,
		"oa_bloom_objective_id": question.OABloomObjectiveID,
		"tipo":                 question.Tipo,
		"question_data":        question.QuestionData,
		"dificultad_relativa":  question.DificultadRelativa,
		"tags":                 question.Tags,
		"oa_bloom_objective":   question.OABloomObjective,
		"question_type":        question.QuestionType,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateQuestion godoc
// @Summary Create a new question
// @Description Create a new question in the bank (auth required)
// @Tags Questions
// @Accept json
// @Produce json
// @Param question body models.Question true "Question data"
// @Success 201 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/questions [post]
func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate question structure
	if err := question.Validate(); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Create(&question).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

// UpdateQuestion godoc
// @Summary Update a question
// @Description Update an existing question (auth required)
// @Tags Questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param question body models.Question true "Updated question data"
// @Success 200 {object} models.Question
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/questions/{id} [put]
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var question models.Question

	if err := db.DB.First(&question, id).Error; err != nil {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate question structure
	if err := question.Validate(); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Save(&question).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

// ValidateAnswerRequest represents the request to validate an answer
type ValidateAnswerRequest struct {
	UserAnswer datatypes.JSON `json:"user_answer"`
}

// ValidateAnswerResponse represents the validation response
type ValidateAnswerResponse struct {
	IsCorrect   bool           `json:"is_correct"`
	Score       float64        `json:"score"`
	Explanation string         `json:"explanation,omitempty"`
	CorrectAnswer interface{}  `json:"correct_answer,omitempty"`
}

// ValidateAnswer godoc
// @Summary Validate a user's answer
// @Description Validate an answer against question's validation_data
// @Tags Questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param answer body ValidateAnswerRequest true "User answer"
// @Success 200 {object} ValidateAnswerResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/questions/{id}/validate [post]
func ValidateAnswer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var question models.Question

	if err := db.DB.First(&question, id).Error; err != nil {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	var req ValidateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the answer
	isCorrect, score, err := question.ValidateAnswer(req.UserAnswer)
	if err != nil {
		// Some question types require manual validation
		if err.Error() == "requires manual or AI validation" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ValidateAnswerResponse{
				IsCorrect: false,
				Score:     0,
				Explanation: "This question requires manual or AI validation",
			})
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Increment usage count
	db.DB.Model(&question).Update("veces_usada", question.VecesUsada+1)

	// Get explanation if available
	var questionData map[string]interface{}
	json.Unmarshal(question.QuestionData, &questionData)
	explanation, _ := questionData["explicacion"].(string)

	response := ValidateAnswerResponse{
		IsCorrect:   isCorrect,
		Score:       score,
		Explanation: explanation,
	}

	// Include correct answer if incorrect
	if !isCorrect {
		var validationData map[string]interface{}
		json.Unmarshal(question.ValidationData, &validationData)
		response.CorrectAnswer = validationData["respuesta_correcta"]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
