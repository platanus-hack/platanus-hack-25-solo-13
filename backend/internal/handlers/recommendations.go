package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/platanus-hack-25/lumera_app/internal/db"
	authmiddleware "github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/models"
)

// DailyRecommendation represents a personalized learning recommendation
type DailyRecommendation struct {
	OA                 models.ObjetivoAprendizaje `json:"oa"`
	OABloomObjective   models.OABloomObjective    `json:"oa_bloom_objective"`
	BloomLevel         models.BloomLevel          `json:"bloom_level"`
	RecommendationType string                     `json:"recommendation_type"` // "next_level", "practice_more", "new_topic"
	Reason             string                     `json:"reason"`
	Priority           int                        `json:"priority"` // 1-5, higher = more important
	EstimatedMinutes   int                        `json:"estimated_minutes"`
	NumeroPreguntas    int                        `json:"numero_preguntas"`
	XPReward           int                        `json:"xp_reward"`
	TokenReward        int                        `json:"token_reward"`
}

// GetDailyRecommendation godoc
// @Summary Get personalized daily learning recommendation
// @Description Get a smart recommendation based on diagnostic results, practice history, and learning progress
// @Tags Recommendations
// @Produce json
// @Success 200 {object} DailyRecommendation
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/recommendations/daily [get]
func GetDailyRecommendation(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user's profile to know their course/subjects
	var profile models.StudentProfile
	if err := db.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Get user's completed diagnostic sessions to know their Bloom levels
	var diagnosticSessions []models.DiagnosticSession
	db.DB.Where("user_id = ? AND estado = ?", userID, "completado").
		Order("completed_at DESC").
		Find(&diagnosticSessions)

	// Build map of materia_id -> average_bloom_level
	diagnosticLevels := make(map[uint]float64)
	for _, session := range diagnosticSessions {
		var estrategia map[string]interface{}
		if err := json.Unmarshal(session.Estrategia, &estrategia); err == nil {
			if avgBloom, ok := estrategia["average_bloom_level"].(float64); ok {
				diagnosticLevels[session.MateriaID] = avgBloom
			}
		}
	}

	// Get practice sessions to know what they've already practiced
	var practiceSessions []models.PracticeSession
	db.DB.Where("user_id = ? AND estado = ?", userID, "completado").
		Order("completed_at DESC").
		Limit(50).
		Find(&practiceSessions)

	// Build map of oa_id -> practice_count
	oaPracticeCounts := make(map[uint]int)
	for _, session := range practiceSessions {
		oaPracticeCounts[session.OAID]++
	}

	// Find recommendation based on diagnostic levels and practice history
	recommendation := findBestRecommendation(diagnosticLevels, oaPracticeCounts, userID)

	if recommendation == nil {
		// Fallback: random OA if no recommendation found
		recommendation = getFallbackRecommendation()
	}

	if recommendation == nil {
		http.Error(w, "No recommendations available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendation)
}

// findBestRecommendation finds the best OA to recommend based on user's learning state
func findBestRecommendation(
	diagnosticLevels map[uint]float64,
	oaPracticeCounts map[uint]int,
	userID uint,
) *DailyRecommendation {
	// Strategy:
	// 1. If user has diagnostics, recommend next level for a subject they've diagnosed
	// 2. Prefer OAs they haven't practiced much
	// 3. Prioritize by categoria and orden

	// Get available subjects (Lengua for now)
	var materia models.Materia
	if err := db.DB.Where("nombre = ?", "Lengua y Literatura").First(&materia).Error; err != nil {
		return nil
	}

	// Get user's bloom level for this subject
	userBloomLevel := 2.0 // Default to level 2 (Comprensión)
	if level, ok := diagnosticLevels[materia.ID]; ok {
		userBloomLevel = level
	}

	// Get OAs for this subject, ordered by categoria and orden
	var oas []models.ObjetivoAprendizaje
	db.DB.Where("materia_id = ? AND activo = ?", materia.ID, true).
		Order("orden ASC").
		Limit(50).
		Find(&oas)

	if len(oas) == 0 {
		return nil
	}

	// Score each OA and pick the best one
	type scoredOA struct {
		oa          models.ObjetivoAprendizaje
		score       float64
		practiceGap int // Days since last practice
	}

	var scoredOAs []scoredOA

	for _, oa := range oas {
		practiceCount := oaPracticeCounts[oa.ID]

		// Calculate score (higher = better)
		score := 100.0

		// Boost if not practiced much
		if practiceCount == 0 {
			score += 50.0
		} else if practiceCount < 3 {
			score += 30.0 - float64(practiceCount*10)
		} else {
			score -= float64(practiceCount * 5) // Penalize over-practiced OAs
		}

		// Boost if early in the curriculum (lower orden)
		if oa.Orden != nil && *oa.Orden <= 5 {
			score += 20.0
		}

		// Boost specific categories
		switch oa.Categoria {
		case "Lectura":
			score += 15.0
		case "Escritura":
			score += 10.0
		}

		scoredOAs = append(scoredOAs, scoredOA{
			oa:          oa,
			score:       score,
			practiceGap: practiceCount,
		})
	}

	// Sort by score (highest first)
	if len(scoredOAs) == 0 {
		return nil
	}

	bestScored := scoredOAs[0]
	for _, s := range scoredOAs {
		if s.score > bestScored.score {
			bestScored = s
		}
	}

	selectedOA := bestScored.oa

	// Get the appropriate Bloom objective for user's level
	targetBloomLevel := int(userBloomLevel)
	if targetBloomLevel < 1 {
		targetBloomLevel = 1
	}
	if targetBloomLevel > 6 {
		targetBloomLevel = 6
	}

	var oaBloomObjective models.OABloomObjective
	if err := db.DB.Where("oa_id = ? AND bloom_level_id = ?", selectedOA.ID, targetBloomLevel).
		First(&oaBloomObjective).Error; err != nil {
		// Fallback to level 2 if specific level not found
		db.DB.Where("oa_id = ? AND bloom_level_id = ?", selectedOA.ID, 2).
			First(&oaBloomObjective)
	}

	// Get Bloom level info
	var bloomLevel models.BloomLevel
	db.DB.First(&bloomLevel, oaBloomObjective.BloomLevelID)

	// Determine recommendation type and reason
	recType := "new_topic"
	reason := "Este objetivo es perfecto para comenzar tu día"

	if bestScored.practiceGap == 0 {
		recType = "new_topic"
		reason = "Nuevo objetivo recomendado para ti"
	} else if bestScored.practiceGap < 3 {
		recType = "practice_more"
		reason = "Continúa practicando para dominar este objetivo"
	}

	// Calculate rewards
	numPreguntas := 10
	xpReward := numPreguntas * 5         // 5 XP per question
	tokenReward := numPreguntas / 5      // 1 token per 5 questions
	estimatedMinutes := numPreguntas * 2 // 2 minutes per question

	return &DailyRecommendation{
		OA:                 selectedOA,
		OABloomObjective:   oaBloomObjective,
		BloomLevel:         bloomLevel,
		RecommendationType: recType,
		Reason:             reason,
		Priority:           5 - bestScored.practiceGap, // Higher priority if less practice
		EstimatedMinutes:   estimatedMinutes,
		NumeroPreguntas:    numPreguntas,
		XPReward:           xpReward,
		TokenReward:        tokenReward,
	}
}

// getFallbackRecommendation returns a default recommendation if no smart recommendation found
func getFallbackRecommendation() *DailyRecommendation {
	// Get a random active OA from Lengua
	var materia models.Materia
	if err := db.DB.Where("nombre = ?", "Lengua y Literatura").First(&materia).Error; err != nil {
		return nil
	}

	var count int64
	db.DB.Model(&models.ObjetivoAprendizaje{}).
		Where("materia_id = ? AND activo = ?", materia.ID, true).
		Count(&count)

	if count == 0 {
		return nil
	}

	// Get random offset
	rand.Seed(time.Now().UnixNano())
	offset := rand.Intn(int(count))

	var oa models.ObjetivoAprendizaje
	db.DB.Where("materia_id = ? AND activo = ?", materia.ID, true).
		Offset(offset).
		Limit(1).
		First(&oa)

	// Get bloom objective at level 2 (Comprensión)
	var oaBloomObjective models.OABloomObjective
	if err := db.DB.Where("oa_id = ? AND bloom_level_id = ?", oa.ID, 2).
		First(&oaBloomObjective).Error; err != nil {
		return nil
	}

	var bloomLevel models.BloomLevel
	db.DB.First(&bloomLevel, 2)

	return &DailyRecommendation{
		OA:                 oa,
		OABloomObjective:   oaBloomObjective,
		BloomLevel:         bloomLevel,
		RecommendationType: "new_topic",
		Reason:             "Recomendado para tu nivel actual",
		Priority:           3,
		EstimatedMinutes:   20,
		NumeroPreguntas:    10,
		XPReward:           50,
		TokenReward:        2,
	}
}
