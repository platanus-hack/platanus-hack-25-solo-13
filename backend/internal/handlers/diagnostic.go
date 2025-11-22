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

// GetNextQuestion godoc
// @Summary Get next adaptive question
// @Description Get the next question based on adaptive strategy
// @Tags Diagnostic
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/diagnostic-sessions/{id}/next-question [get]
func GetNextQuestion(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var session models.DiagnosticSession

	if err := db.DB.First(&session, sessionID).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Parse current strategy
	var strategy models.AdaptiveStrategy
	if err := json.Unmarshal(session.Estrategia, &strategy); err != nil {
		http.Error(w, "Invalid strategy data", http.StatusInternalServerError)
		return
	}

	// If first question, initialize strategy
	if session.PreguntasTotales == 0 {
		// Get all OAs for this materia
		var oaBloomObjectives []models.OABloomObjective
		err := db.DB.Joins("JOIN objetivos_aprendizaje ON oa_bloom_objectives.oa_id = objetivos_aprendizaje.id").
			Where("objetivos_aprendizaje.materia_id = ?", session.MateriaID).
			Group("oa_bloom_objectives.oa_id").
			Select("MIN(oa_bloom_objectives.id) as id, oa_bloom_objectives.oa_id").
			Limit(5).
			Find(&oaBloomObjectives).Error

		if err != nil || len(oaBloomObjectives) == 0 {
			http.Error(w, "No OAs available for this materia", http.StatusNotFound)
			return
		}

		// Initialize strategy with 5 random OAs
		strategy.NivelBloomActual = 3 // Start at Bloom level 3
		strategy.OAsEvaluados = []uint{}

		// Store OA IDs to evaluate
		var oasToEvaluate []uint
		for _, oa := range oaBloomObjectives {
			oasToEvaluate = append(oasToEvaluate, oa.OAID)
		}

		// Add oas_a_evaluar to strategy
		strategyMap := map[string]interface{}{
			"nivel_bloom_actual":      strategy.NivelBloomActual,
			"oas_evaluados":           strategy.OAsEvaluados,
			"oas_a_evaluar":           oasToEvaluate,
			"aciertos_consecutivos":   0,
			"fallos_consecutivos":     0,
			"patron_respuestas":       []string{},
		}
		strategyJSON, _ := json.Marshal(strategyMap)
		session.Estrategia = datatypes.JSON(strategyJSON)

		// Re-unmarshal to get updated strategy
		json.Unmarshal(session.Estrategia, &strategy)
	}

	// Get OAs to evaluate
	var strategyMap map[string]interface{}
	json.Unmarshal(session.Estrategia, &strategyMap)
	oasToEvaluateInterface, _ := strategyMap["oas_a_evaluar"].([]interface{})
	var oasToEvaluate []uint
	for _, v := range oasToEvaluateInterface {
		if floatVal, ok := v.(float64); ok {
			oasToEvaluate = append(oasToEvaluate, uint(floatVal))
		}
	}

	// Check if we have OAs left to evaluate
	if len(strategy.OAsEvaluados) >= len(oasToEvaluate) {
		http.Error(w, "No more questions available", http.StatusNotFound)
		return
	}

	// Get next OA to evaluate
	var nextOAID uint
	for _, oaID := range oasToEvaluate {
		alreadyEvaluated := false
		for _, evaluatedID := range strategy.OAsEvaluados {
			if oaID == evaluatedID {
				alreadyEvaluated = true
				break
			}
		}
		if !alreadyEvaluated {
			nextOAID = oaID
			break
		}
	}

	if nextOAID == 0 {
		http.Error(w, "No more OAs to evaluate", http.StatusNotFound)
		return
	}

	// Find OABloomObjective for this OA at current Bloom level
	var oaBloomObjective models.OABloomObjective
	err := db.DB.Where("oa_id = ? AND bloom_level_id = ?", nextOAID, strategy.NivelBloomActual).
		First(&oaBloomObjective).Error

	if err != nil {
		// Bloom level not found, try closest one
		db.DB.Where("oa_id = ?", nextOAID).
			Order(fmt.Sprintf("ABS(bloom_level_id - %d)", strategy.NivelBloomActual)).
			First(&oaBloomObjective)
	}

	// Find available question for this OABloomObjective
	var question models.Question
	err = db.DB.Where("oa_bloom_objective_id = ? AND activa = ? AND (tipo_uso = ? OR tipo_uso = ?)",
		oaBloomObjective.ID, true, "diagnostico", "all").
		Order("RANDOM()").
		First(&question).Error

	if err != nil {
		http.Error(w, "No questions available for current criteria", http.StatusNotFound)
		return
	}

	// Mark OA as evaluated
	strategy.OAsEvaluados = append(strategy.OAsEvaluados, nextOAID)
	strategyMap["oas_evaluados"] = strategy.OAsEvaluados
	updatedStrategyJSON, _ := json.Marshal(strategyMap)
	session.Estrategia = datatypes.JSON(updatedStrategyJSON)
	db.DB.Save(&session)

	// Prepare response without validation_data
	response := map[string]interface{}{
		"id":                  question.ID,
		"oa_bloom_objective_id": question.OABloomObjectiveID,
		"tipo":                question.Tipo,
		"question_data":       question.QuestionData,
		"question_number":     session.PreguntasTotales + 1,
		"total_questions":     len(oasToEvaluate),
		"current_bloom_level": strategy.NivelBloomActual,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Update adaptive strategy
	var strategyMap map[string]interface{}
	json.Unmarshal(session.Estrategia, &strategyMap)

	currentBloomLevel, _ := strategyMap["nivel_bloom_actual"].(float64)

	if isCorrect {
		// Increase Bloom level (max 6)
		newLevel := int(currentBloomLevel) + 1
		if newLevel > 6 {
			newLevel = 6
		}
		strategyMap["nivel_bloom_actual"] = newLevel
	} else {
		// Decrease Bloom level (min 1)
		newLevel := int(currentBloomLevel) - 1
		if newLevel < 1 {
			newLevel = 1
		}
		strategyMap["nivel_bloom_actual"] = newLevel
	}

	// Update strategy in session
	updatedStrategyJSON, _ := json.Marshal(strategyMap)
	session.Estrategia = datatypes.JSON(updatedStrategyJSON)

	db.DB.Save(&session)

	response := map[string]interface{}{
		"is_correct": isCorrect,
		"score":      score,
		"answer_id":  answer.ID,
		"new_bloom_level": strategyMap["nivel_bloom_actual"],
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

	// Generate diagnostic_results based on answers
	var answers []models.DiagnosticAnswer
	db.DB.Where("session_id = ?", session.ID).
		Preload("OABloomObjective").
		Preload("OABloomObjective.OA").
		Preload("OABloomObjective.BloomLevel").
		Find(&answers)

	// Group answers by OA
	oaAnswers := make(map[uint][]models.DiagnosticAnswer)
	for _, answer := range answers {
		oaID := answer.OABloomObjective.OAID
		oaAnswers[oaID] = append(oaAnswers[oaID], answer)
	}

	// Create diagnostic results for each OA
	var bloomLevels []uint
	for oaID, oaAnswerList := range oaAnswers {
		var correct int
		var maxBloomLevel uint = 0
		var bloomLevelName string

		for _, answer := range oaAnswerList {
			if answer.IsCorrect != nil && *answer.IsCorrect {
				correct++
				// Track highest bloom level achieved for this OA
				if answer.OABloomObjective.BloomLevelID > maxBloomLevel {
					maxBloomLevel = answer.OABloomObjective.BloomLevelID
					bloomLevelName = answer.OABloomObjective.BloomLevel.Nombre
					bloomLevels = append(bloomLevels, maxBloomLevel)
				}
			}
		}

		total := len(oaAnswerList)
		percentage := 0
		if total > 0 {
			percentage = (correct * 100) / total
		}

		// Generate recommendation based on result
		recommendation := fmt.Sprintf("Completaste %d de %d preguntas correctamente.", correct, total)
		if percentage >= 80 {
			recommendation += " ¡Excelente dominio de este objetivo!"
		} else if percentage >= 60 {
			recommendation += " Buen nivel, con práctica adicional alcanzarás la maestría."
		} else {
			recommendation += " Te recomendamos reforzar este contenido con ejercicios adicionales."
		}

		// Create diagnostic result
		result := models.DiagnosticResult{
			SessionID:            session.ID,
			OAID:                 oaID,
			NivelBloomDominado:   int(maxBloomLevel),
			NivelBloomNombre:     bloomLevelName,
			PreguntasRespondidas: total,
			PreguntasCorrectas:   correct,
			PorcentajeAciertos:   percentage,
			Recomendacion:        recommendation,
		}

		if err := db.DB.Create(&result).Error; err != nil {
			fmt.Printf("Error creating diagnostic result: %v\n", err)
		}
	}

	// Calculate average Bloom level from correct answers
	averageBloomLevel := 0
	if len(bloomLevels) > 0 {
		var sum uint = 0
		for _, level := range bloomLevels {
			sum += level
		}
		averageBloomLevel = int(sum) / len(bloomLevels)
	}

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
		"message":              "Diagnostic completed successfully",
		"session":              session,
		"average_bloom_level":  averageBloomLevel,
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
