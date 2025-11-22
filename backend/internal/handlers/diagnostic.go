package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	authmiddleware "github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"github.com/platanus-hack-25/lumera_app/internal/services"
	"gorm.io/datatypes"
)

// StartDiagnosticRequest represents the request to start a diagnostic session
type StartDiagnosticRequest struct {
	MateriaID uint `json:"materia_id"`
}

// StartDiagnostic godoc
// @Summary Start a diagnostic session
// @Description Create a new diagnostic session for a student
// @Tags Diagnostic
// @Accept json
// @Produce json
// @Param diagnostic body StartDiagnosticRequest true "Diagnostic session data"
// @Success 201 {object} models.DiagnosticSession
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/diagnostic-sessions [post]
func StartDiagnostic(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req StartDiagnosticRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Count previous attempts
	var count int64
	db.DB.Model(&models.DiagnosticSession{}).
		Where("user_id = ? AND materia_id = ?", userID, req.MateriaID).
		Count(&count)

	// Initialize adaptive strategy
	defaultStrategy := map[string]interface{}{
		"nivel_bloom_actual":      2,
		"oas_evaluados":           []uint{},
		"aciertos_consecutivos":   0,
		"fallos_consecutivos":     0,
		"patron_respuestas":       []string{},
	}
	strategyJSON, _ := json.Marshal(defaultStrategy)

	session := models.DiagnosticSession{
		UserID:        userID,
		MateriaID:     req.MateriaID,
		NumeroIntento: int(count) + 1,
		Estado:        "en_progreso",
		Estrategia:    datatypes.JSON(strategyJSON),
		StartedAt:     time.Now(),
	}

	if err := db.DB.Create(&session).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// GetSessionProgress godoc
// @Summary Get diagnostic session progress
// @Description Retrieve diagnostic session details and progress
// @Tags Diagnostic
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} models.DiagnosticSession
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/diagnostic-sessions/{id} [get]
func GetSessionProgress(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var session models.DiagnosticSession

	if err := db.DB.Preload("Materia").Preload("Answers").First(&session, id).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// SubmitAnswerRequest represents the request to submit an answer
type SubmitAnswerRequest struct {
	QuestionID         uint           `json:"question_id"`
	UserAnswer         datatypes.JSON `json:"user_answer"`
	TiempoSegundos     *int           `json:"tiempo_segundos"`
}

// SubmitAnswer godoc
// @Summary Submit an answer during diagnostic
// @Description Submit an answer and get validation result
// @Tags Diagnostic
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Param answer body SubmitAnswerRequest true "Answer data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/diagnostic-sessions/{id}/answer [post]
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var session models.DiagnosticSession

	if err := db.DB.First(&session, sessionID).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	var req SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get question
	var question models.Question
	if err := db.DB.Preload("OABloomObjective").First(&question, req.QuestionID).Error; err != nil {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	// Validate answer
	isCorrect, score, err := question.ValidateAnswer(req.UserAnswer)
	if err != nil && err.Error() != "requires manual or AI validation" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save answer
	answer := models.DiagnosticAnswer{
		SessionID:          session.ID,
		QuestionID:         req.QuestionID,
		OABloomObjectiveID: question.OABloomObjectiveID,
		BloomLevelID:       question.OABloomObjective.BloomLevelID,
		UserAnswer:         req.UserAnswer,
		IsCorrect:          &isCorrect,
		Score:              &score,
		TiempoSegundos:     req.TiempoSegundos,
	}

	if err := db.DB.Create(&answer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update session stats
	session.PreguntasTotales++
	if isCorrect {
		session.PreguntasCorrectas++
	}
	db.DB.Save(&session)

	response := map[string]interface{}{
		"is_correct": isCorrect,
		"score":      score,
		"answer_id":  answer.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompleteDiagnostic godoc
// @Summary Complete a diagnostic session
// @Description Mark session as completed and generate results
// @Tags Diagnostic
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/diagnostic-sessions/{id}/complete [post]
func CompleteDiagnostic(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var session models.DiagnosticSession

	if err := db.DB.First(&session, sessionID).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	now := time.Now()
	session.Estado = "completado"
	session.CompletedAt = &now

	if err := db.DB.Save(&session).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Generate diagnostic_results based on answers
	// This would analyze answers and create DiagnosticResult records
	// which will trigger the auto-update of student_oa_progress

	// Gamification events (async)
	go func() {
		// Award XP for completing diagnostic
		gamificationService.AddXP(session.UserID, 100, "diagnostic_completed")

		// Bonus coins for good performance
		if session.PreguntasTotales > 0 {
			scorePercent := (session.PreguntasCorrectas * 100) / session.PreguntasTotales
			if scorePercent >= 80 {
				gamificationService.AddCoins(session.UserID, 50, "diagnostic_excellent")
			} else if scorePercent >= 60 {
				gamificationService.AddCoins(session.UserID, 25, "diagnostic_good")
			}

			// Check for diagnostic achievement unlocks
			unlockService.CheckAndUnlock(session.UserID, services.UnlockEvent{
				Type: "diagnostic_achievement",
				Key:  fmt.Sprintf("diagnostic_materia_%d_score_%d", session.MateriaID, scorePercent),
				Data: map[string]interface{}{
					"materia_id": session.MateriaID,
					"score":      scorePercent,
				},
			})
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Diagnostic completed successfully",
		"session": session,
	})
}

// GetDiagnosticResults godoc
// @Summary Get diagnostic results
// @Description Retrieve consolidated results for a completed diagnostic session
// @Tags Diagnostic
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {array} models.DiagnosticResult
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/diagnostic-sessions/{id}/results [get]
func GetDiagnosticResults(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var results []models.DiagnosticResult

	if err := db.DB.Where("session_id = ?", sessionID).
		Preload("OA").
		Find(&results).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
