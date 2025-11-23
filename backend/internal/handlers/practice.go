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

// StartPracticeRequest represents the request to start a practice session
type StartPracticeRequest struct {
	OAID               uint `json:"oa_id"`
	OABloomObjectiveID uint `json:"oa_bloom_objective_id"`
	NumeroPreguntas    *int `json:"numero_preguntas"` // Optional, default 10
}

// StartPractice godoc
// @Summary Start a practice session for a specific OA
// @Description Create a new practice session focusing on one OA at a specific Bloom level
// @Tags Practice
// @Accept json
// @Produce json
// @Param practice body StartPracticeRequest true "Practice session data"
// @Success 201 {object} models.PracticeSession
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/practice-sessions [post]
func StartPractice(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req StartPracticeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get OA Bloom Objective to verify it exists
	var oaBloomObjective models.OABloomObjective
	if err := db.DB.First(&oaBloomObjective, req.OABloomObjectiveID).Error; err != nil {
		http.Error(w, "OA Bloom Objective not found", http.StatusNotFound)
		return
	}

	// Set default number of questions
	numPreguntas := 10
	if req.NumeroPreguntas != nil && *req.NumeroPreguntas > 0 {
		numPreguntas = *req.NumeroPreguntas
	}

	// Initialize strategy
	strategy := models.PracticeStrategy{
		NivelBloomActual:     int(oaBloomObjective.BloomLevelID),
		AciertosConsecutivos: 0,
		FallosConsecutivos:   0,
		AciertosPorNivel:     make(map[int]int),
		FallosPorNivel:       make(map[int]int),
		PatronRespuestas:     []string{},
	}
	strategyJSON, _ := json.Marshal(strategy)

	session := models.PracticeSession{
		UserID:             userID,
		OAID:               req.OAID,
		OABloomObjectiveID: req.OABloomObjectiveID,
		BloomLevelInicial:  int(oaBloomObjective.BloomLevelID),
		NumeroPreguntas:    numPreguntas,
		Estado:             "en_progreso",
		Estrategia:         datatypes.JSON(strategyJSON),
		StartedAt:          time.Now(),
	}

	if err := db.DB.Create(&session).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// GetPracticeSessions godoc
// @Summary Get user's practice sessions
// @Description Get list of practice sessions with optional filters
// @Tags Practice
// @Produce json
// @Param oa_id query int false "Filter by OA ID"
// @Param estado query string false "Filter by estado (en_progreso, completado)"
// @Success 200 {array} models.PracticeSession
// @Failure 500 {string} string "Server error"
// @Security BearerAuth
// @Router /api/practice-sessions [get]
func GetPracticeSessions(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := db.DB.Where("user_id = ?", userID).Order("started_at DESC")

	// Apply filters
	if oaID := r.URL.Query().Get("oa_id"); oaID != "" {
		query = query.Where("oa_id = ?", oaID)
	}

	if estado := r.URL.Query().Get("estado"); estado != "" {
		query = query.Where("estado = ?", estado)
	}

	var sessions []models.PracticeSession
	if err := query.Preload("OA").Find(&sessions).Error; err != nil {
		http.Error(w, "Error fetching sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// GetPracticeNextQuestion godoc
// @Summary Get next practice question
// @Description Get the next question for practice session with adaptive difficulty
// @Tags Practice
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/practice-sessions/{id}/next-question [get]
func GetPracticeNextQuestion(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var session models.PracticeSession

	if err := db.DB.First(&session, sessionID).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	// Check if session is complete
	if session.PreguntasRespondidas >= session.NumeroPreguntas {
		http.Error(w, "Session already completed", http.StatusBadRequest)
		return
	}

	// Parse current strategy
	var strategy models.PracticeStrategy
	if err := json.Unmarshal(session.Estrategia, &strategy); err != nil {
		http.Error(w, "Invalid strategy data", http.StatusInternalServerError)
		return
	}

	// Find question at current Bloom level for this OA
	// First try exact level, then try nearby levels
	var question models.Question
	var err error

	// Get OA ID from the OABloomObjective
	var oaBloomObjective models.OABloomObjective
	if err := db.DB.First(&oaBloomObjective, session.OABloomObjectiveID).Error; err != nil {
		http.Error(w, "OA Bloom Objective not found", http.StatusNotFound)
		return
	}

	// Try current level first
	err = db.DB.Joins("JOIN oa_bloom_objectives ON questions.oa_bloom_objective_id = oa_bloom_objectives.id").
		Where("oa_bloom_objectives.oa_id = ? AND oa_bloom_objectives.bloom_level_id = ? AND questions.activa = ? AND (questions.tipo_uso = ? OR questions.tipo_uso = ?)",
			oaBloomObjective.OAID, strategy.NivelBloomActual, true, "practica", "all").
		Order("RANDOM()").
		First(&question).Error

	if err != nil {
		// Try nearby levels (±1)
		err = db.DB.Joins("JOIN oa_bloom_objectives ON questions.oa_bloom_objective_id = oa_bloom_objectives.id").
			Where("oa_bloom_objectives.oa_id = ? AND oa_bloom_objectives.bloom_level_id IN (?, ?, ?) AND questions.activa = ? AND (questions.tipo_uso = ? OR questions.tipo_uso = ?)",
				oaBloomObjective.OAID,
				max(1, strategy.NivelBloomActual-1),
				strategy.NivelBloomActual,
				min(6, strategy.NivelBloomActual+1),
				true, "practica", "all").
			Order("RANDOM()").
			First(&question).Error
	}

	if err != nil {
		http.Error(w, "No questions available for current criteria", http.StatusNotFound)
		return
	}

	// Prepare response without validation_data
	response := map[string]interface{}{
		"id":                   question.ID,
		"oa_bloom_objective_id": question.OABloomObjectiveID,
		"tipo":                 question.Tipo,
		"question_data":        question.QuestionData,
		"question_number":      session.PreguntasRespondidas + 1,
		"total_questions":      session.NumeroPreguntas,
		"current_bloom_level":  strategy.NivelBloomActual,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SubmitPracticeAnswer godoc
// @Summary Submit an answer during practice
// @Description Submit an answer and get validation result with adaptive difficulty adjustment
// @Tags Practice
// @Accept json
// @Produce json
// @Param id path int true "Session ID"
// @Param answer body SubmitAnswerRequest true "Answer data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/practice-sessions/{id}/answer [post]
func SubmitPracticeAnswer(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var session models.PracticeSession

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
	answer := models.PracticeAnswer{
		SessionID:      session.ID,
		QuestionID:     req.QuestionID,
		BloomLevelID:   question.OABloomObjective.BloomLevelID,
		UserAnswer:     req.UserAnswer,
		IsCorrect:      &isCorrect,
		Score:          &score,
		TiempoSegundos: req.TiempoSegundos,
	}

	if err := db.DB.Create(&answer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update session stats
	session.PreguntasRespondidas++
	if isCorrect {
		session.PreguntasCorrectas++
	}

	// Update adaptive strategy
	var strategy models.PracticeStrategy
	json.Unmarshal(session.Estrategia, &strategy)

	questionBloomLevel := int(question.OABloomObjective.BloomLevelID)

	if isCorrect {
		strategy.AciertosConsecutivos++
		strategy.FallosConsecutivos = 0
		strategy.PatronRespuestas = append(strategy.PatronRespuestas, "C")

		// Track correct answers by level
		if strategy.AciertosPorNivel == nil {
			strategy.AciertosPorNivel = make(map[int]int)
		}
		strategy.AciertosPorNivel[questionBloomLevel]++

		// Increase difficulty after 2 consecutive correct answers at current level
		if strategy.AciertosConsecutivos >= 2 && strategy.NivelBloomActual < 6 {
			strategy.NivelBloomActual++
			strategy.AciertosConsecutivos = 0
		}
	} else {
		strategy.FallosConsecutivos++
		strategy.AciertosConsecutivos = 0
		strategy.PatronRespuestas = append(strategy.PatronRespuestas, "I")

		// Track incorrect answers by level
		if strategy.FallosPorNivel == nil {
			strategy.FallosPorNivel = make(map[int]int)
		}
		strategy.FallosPorNivel[questionBloomLevel]++

		// Decrease difficulty after 2 consecutive incorrect answers
		if strategy.FallosConsecutivos >= 2 && strategy.NivelBloomActual > 1 {
			strategy.NivelBloomActual--
			strategy.FallosConsecutivos = 0
		}
	}

	// Save updated strategy
	strategyJSON, _ := json.Marshal(strategy)
	session.Estrategia = datatypes.JSON(strategyJSON)
	db.DB.Save(&session)

	response := map[string]interface{}{
		"is_correct":           isCorrect,
		"score":                score,
		"new_bloom_level":      strategy.NivelBloomActual,
		"preguntas_respondidas": session.PreguntasRespondidas,
		"total_preguntas":      session.NumeroPreguntas,
		"is_complete":          session.PreguntasRespondidas >= session.NumeroPreguntas,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompletePracticeSession godoc
// @Summary Complete a practice session
// @Description Finalize practice session and calculate final Bloom level
// @Tags Practice
// @Produce json
// @Param id path int true "Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/practice-sessions/{id}/complete [post]
func CompletePracticeSession(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "id")
	var session models.PracticeSession

	if err := db.DB.Preload("Answers").Preload("OA").First(&session, sessionID).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	if session.Estado == "completado" {
		http.Error(w, "Session already completed", http.StatusBadRequest)
		return
	}

	// Parse strategy
	var strategy models.PracticeStrategy
	json.Unmarshal(session.Estrategia, &strategy)

	// Calculate final Bloom level based on performance
	finalBloomLevel := calculateFinalBloomLevel(
		session.BloomLevelInicial,
		strategy.AciertosPorNivel,
		strategy.FallosPorNivel,
		session.PreguntasCorrectas,
		session.PreguntasRespondidas,
	)

	now := time.Now()
	session.Estado = "completado"
	session.BloomLevelFinal = &finalBloomLevel
	session.CompletedAt = &now

	// Create result summary
	resultado := map[string]interface{}{
		"bloom_level_inicial": session.BloomLevelInicial,
		"bloom_level_final":   finalBloomLevel,
		"cambio_nivel":        finalBloomLevel - session.BloomLevelInicial,
		"porcentaje_aciertos": float64(session.PreguntasCorrectas) / float64(session.PreguntasRespondidas) * 100,
		"preguntas_totales":   session.PreguntasRespondidas,
		"preguntas_correctas": session.PreguntasCorrectas,
		"aciertos_por_nivel":  strategy.AciertosPorNivel,
		"fallos_por_nivel":    strategy.FallosPorNivel,
		"patron_respuestas":   strategy.PatronRespuestas,
	}
	resultadoJSON, _ := json.Marshal(resultado)
	session.Resultado = datatypes.JSON(resultadoJSON)

	if err := db.DB.Save(&session).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update user's progress for this OA-Bloom objective
	accuracy := float64(session.PreguntasCorrectas) / float64(session.PreguntasRespondidas) * 100
	if err := updateUserProgress(session.UserID, session.OABloomObjectiveID, finalBloomLevel, accuracy, session.PreguntasCorrectas, session.PreguntasRespondidas); err != nil {
		// Log error but don't fail the request
		http.Error(w, "Failed to update progress: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Award XP and Coins for completing practice
	xpEarned := session.PreguntasCorrectas * 5 // 5 XP per correct answer
	coinsEarned := session.PreguntasCorrectas / 5 // 1 coin per 5 correct answers

	gamificationService := services.NewGamificationService()

	// Add XP
	xpResult, xpErr := gamificationService.AddXP(session.UserID, xpEarned, "practice_session_complete")
	if xpErr != nil {
		// Log but don't fail
		println("Error adding XP:", xpErr.Error())
	}

	// Add Coins
	coinsErr := gamificationService.AddCoins(session.UserID, coinsEarned, "practice_session_complete")
	if coinsErr != nil {
		println("Error adding coins:", coinsErr.Error())
	}

	// Update Streak
	streakResult, streakErr := gamificationService.UpdateStreak(session.UserID)
	if streakErr != nil {
		println("Error updating streak:", streakErr.Error())
	}

	response := map[string]interface{}{
		"session":             session,
		"bloom_level_inicial": session.BloomLevelInicial,
		"bloom_level_final":   finalBloomLevel,
		"cambio_nivel":        finalBloomLevel - session.BloomLevelInicial,
		"resultado":           resultado,
		"rewards": map[string]interface{}{
			"xp_earned":      xpEarned,
			"coins_earned":   coinsEarned,
			"xp_result":      xpResult,
			"streak_result":  streakResult,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// calculateFinalBloomLevel determines the final Bloom level based on practice performance
func calculateFinalBloomLevel(
	initialLevel int,
	aciertosPorNivel map[int]int,
	fallosPorNivel map[int]int,
	totalCorrect int,
	totalAnswered int,
) int {
	// Base level on overall accuracy
	accuracy := float64(totalCorrect) / float64(totalAnswered)

	// Calculate weighted average of levels where user succeeded
	var weightedSum, totalWeight float64
	for level, aciertos := range aciertosPorNivel {
		fallos := 0
		if f, ok := fallosPorNivel[level]; ok {
			fallos = f
		}
		total := aciertos + fallos
		if total > 0 {
			levelAccuracy := float64(aciertos) / float64(total)
			weight := float64(total) * levelAccuracy
			weightedSum += float64(level) * weight
			totalWeight += weight
		}
	}

	var finalLevel int
	if totalWeight > 0 {
		finalLevel = int(weightedSum / totalWeight)
	} else {
		finalLevel = initialLevel
	}

	// Apply accuracy bonus/penalty
	if accuracy >= 0.8 && finalLevel < 6 {
		finalLevel++ // Increase if high accuracy
	} else if accuracy < 0.5 && finalLevel > 1 {
		finalLevel-- // Decrease if low accuracy
	}

	// Ensure level is in valid range
	if finalLevel < 1 {
		finalLevel = 1
	}
	if finalLevel > 6 {
		finalLevel = 6
	}

	return finalLevel
}

// updateUserProgress updates or creates student progress for an OA-Bloom objective
func updateUserProgress(userID uint, oaBloomObjectiveID uint, bloomLevel int, accuracy float64, correctas int, totales int) error {
	// Determine state based on accuracy
	var estado string
	if accuracy >= 80 {
		estado = "dominado"
	} else if accuracy >= 60 {
		estado = "logrado"
	} else if totales > 0 {
		estado = "en_proceso"
	} else {
		estado = "no_iniciado"
	}

	// Find or create progress record
	var progress models.StudentOAProgress
	result := db.DB.Where("user_id = ? AND oa_bloom_objective_id = ?", userID, oaBloomObjectiveID).First(&progress)

	now := time.Now()
	porcentajeLogro := int(accuracy)

	if result.Error != nil {
		// Create new progress record
		progress = models.StudentOAProgress{
			UserID:              userID,
			OABloomObjectiveID:  oaBloomObjectiveID,
			Estado:              estado,
			PorcentajeLogro:     porcentajeLogro,
			Intentos:            1,
			UltimaActividadFecha: &now,
		}
		if err := db.DB.Create(&progress).Error; err != nil {
			return err
		}
	} else {
		// Update existing record - only if performance improved or state changed
		shouldUpdate := false

		// Update if accuracy improved
		if porcentajeLogro > progress.PorcentajeLogro {
			progress.PorcentajeLogro = porcentajeLogro
			shouldUpdate = true
		}

		// Always update state if it changed
		if estado != progress.Estado {
			progress.Estado = estado
			shouldUpdate = true
		}

		// Always increment attempts and update last activity
		progress.Intentos++
		progress.UltimaActividadFecha = &now
		shouldUpdate = true

		if shouldUpdate {
			if err := db.DB.Save(&progress).Error; err != nil {
				return err
			}
		}
	}

	// Create history record
	puntajeObtenido := float64(correctas)
	puntajeMaximo := float64(totales)

	history := models.StudentOAHistory{
		UserID:             userID,
		OABloomObjectiveID: oaBloomObjectiveID,
		Estado:             estado,
		PorcentajeLogro:    &porcentajeLogro,
		TipoEvento:         "practica",
		PuntajeObtenido:    &puntajeObtenido,
		PuntajeMaximo:      &puntajeMaximo,
		Notas:              fmt.Sprintf("Bloom level %d - %d/%d correctas", bloomLevel, correctas, totales),
	}

	return db.DB.Create(&history).Error
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
