package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"github.com/platanus-hack-25/lumera_app/internal/services"
)

// ============================================================================
// CURSOS (Courses)
// ============================================================================

// GetCursos godoc
// @Summary Get all courses
// @Description Retrieve all courses with optional filters
// @Tags Educational
// @Produce json
// @Param activo query boolean false "Filter by active status"
// @Success 200 {array} models.Curso
// @Failure 500 {object} map[string]interface{}
// @Router /api/cursos [get]
func GetCursos(w http.ResponseWriter, r *http.Request) {
	var cursos []models.Curso
	query := db.DB

	// Filter by active status
	if activoParam := r.URL.Query().Get("activo"); activoParam != "" {
		if activo, err := strconv.ParseBool(activoParam); err == nil {
			query = query.Where("activo = ?", activo)
		}
	}

	if err := query.Order("codigo ASC").Find(&cursos).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cursos)
}

// GetCurso godoc
// @Summary Get course by ID
// @Description Retrieve a specific course with its materias
// @Tags Educational
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Curso
// @Failure 404 {object} map[string]interface{}
// @Router /api/cursos/{id} [get]
func GetCurso(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var curso models.Curso

	if err := db.DB.Preload("Materias").First(&curso, id).Error; err != nil {
		http.Error(w, "Curso not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(curso)
}

// CreateCurso godoc
// @Summary Create a new course
// @Description Create a new course
// @Tags Educational
// @Accept json
// @Produce json
// @Param curso body models.Curso true "Course data"
// @Success 201 {object} models.Curso
// @Failure 400 {object} map[string]interface{}
// @Router /api/cursos [post]
func CreateCurso(w http.ResponseWriter, r *http.Request) {
	var curso models.Curso
	if err := json.NewDecoder(r.Body).Decode(&curso); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Create(&curso).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(curso)
}

// UpdateCurso godoc
// @Summary Update a course
// @Description Update an existing course
// @Tags Educational
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param curso body models.Curso true "Updated course data"
// @Success 200 {object} models.Curso
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/cursos/{id} [put]
func UpdateCurso(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var curso models.Curso

	if err := db.DB.First(&curso, id).Error; err != nil {
		http.Error(w, "Curso not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&curso); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Save(&curso).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(curso)
}

// ============================================================================
// MATERIAS (Subjects)
// ============================================================================

// GetMaterias godoc
// @Summary Get all subjects
// @Description Retrieve all subjects with optional filters
// @Tags Educational
// @Produce json
// @Param activo query boolean false "Filter by active status"
// @Success 200 {array} models.Materia
// @Failure 500 {object} map[string]interface{}
// @Router /api/materias [get]
func GetMaterias(w http.ResponseWriter, r *http.Request) {
	var materias []models.Materia
	query := db.DB

	if activoParam := r.URL.Query().Get("activo"); activoParam != "" {
		if activo, err := strconv.ParseBool(activoParam); err == nil {
			query = query.Where("activo = ?", activo)
		}
	}

	if err := query.Order("nombre ASC").Find(&materias).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materias)
}

// GetMateria godoc
// @Summary Get subject by ID
// @Description Retrieve a specific subject with its learning objectives
// @Tags Educational
// @Produce json
// @Param id path int true "Subject ID"
// @Success 200 {object} models.Materia
// @Failure 404 {object} map[string]interface{}
// @Router /api/materias/{id} [get]
func GetMateria(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var materia models.Materia

	if err := db.DB.Preload("OAs").First(&materia, id).Error; err != nil {
		http.Error(w, "Materia not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materia)
}

// CreateMateria godoc
// @Summary Create a new subject
// @Description Create a new subject
// @Tags Educational
// @Accept json
// @Produce json
// @Param materia body models.Materia true "Subject data"
// @Success 201 {object} models.Materia
// @Failure 400 {object} map[string]interface{}
// @Router /api/materias [post]
func CreateMateria(w http.ResponseWriter, r *http.Request) {
	var materia models.Materia
	if err := json.NewDecoder(r.Body).Decode(&materia); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Create(&materia).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(materia)
}

// UpdateMateria godoc
// @Summary Update a subject
// @Description Update an existing subject
// @Tags Educational
// @Accept json
// @Produce json
// @Param id path int true "Subject ID"
// @Param materia body models.Materia true "Updated subject data"
// @Success 200 {object} models.Materia
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/materias/{id} [put]
func UpdateMateria(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var materia models.Materia

	if err := db.DB.First(&materia, id).Error; err != nil {
		http.Error(w, "Materia not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&materia); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Save(&materia).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(materia)
}

// AssignMateriaToCurso godoc
// @Summary Assign a subject to a course
// @Description Create a relationship between a course and a subject
// @Tags Educational
// @Accept json
// @Produce json
// @Param assignment body map[string]interface{} true "Assignment data (curso_id, materia_id, horas_semanales)"
// @Success 201 {object} models.CursoMateria
// @Failure 400 {object} map[string]interface{}
// @Router /api/curso-materias [post]
func AssignMateriaToCurso(w http.ResponseWriter, r *http.Request) {
	var cursoMateria models.CursoMateria
	if err := json.NewDecoder(r.Body).Decode(&cursoMateria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Create(&cursoMateria).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cursoMateria)
}

// ============================================================================
// BLOOM LEVELS (Read-only)
// ============================================================================

// GetBloomLevels godoc
// @Summary Get all Bloom taxonomy levels
// @Description Retrieve all Bloom levels (seeded data, read-only)
// @Tags Educational
// @Produce json
// @Success 200 {array} models.BloomLevel
// @Failure 500 {object} map[string]interface{}
// @Router /api/bloom-levels [get]
func GetBloomLevels(w http.ResponseWriter, r *http.Request) {
	var levels []models.BloomLevel

	if err := db.DB.Order("nivel ASC").Find(&levels).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(levels)
}

// ============================================================================
// OBJETIVOS DE APRENDIZAJE (Learning Objectives)
// ============================================================================

// CreateOARequest represents the request to create a complete OA with all Bloom levels
type CreateOARequest struct {
	MateriaID   uint                       `json:"materia_id"`
	Codigo      string                     `json:"codigo"`
	Titulo      string                     `json:"titulo"`
	Descripcion string                     `json:"descripcion"`
	Orden       *int                       `json:"orden"`
	BloomObjectives []CreateBloomObjective `json:"bloom_objectives"` // All 6 levels
}

type CreateBloomObjective struct {
	BloomLevelID           uint     `json:"bloom_level_id"`
	ObjetivoEspecifico     string   `json:"objetivo_especifico"`
	IndicadoresLogro       []string `json:"indicadores_logro"`
	TipoActividadSugerida  string   `json:"tipo_actividad_sugerida"`
	ComplejidadEstimada    *int     `json:"complejidad_estimada"`
}

// GetObjetivosAprendizaje godoc
// @Summary Get all learning objectives
// @Description Retrieve all learning objectives with optional filtering by subject
// @Tags Educational
// @Produce json
// @Param materia_id query int false "Filter by subject ID"
// @Param activo query boolean false "Filter by active status"
// @Success 200 {array} models.ObjetivoAprendizaje
// @Failure 500 {object} map[string]interface{}
// @Router /api/objetivos-aprendizaje [get]
func GetObjetivosAprendizaje(w http.ResponseWriter, r *http.Request) {
	var oas []models.ObjetivoAprendizaje
	query := db.DB.Preload("BloomObjectives.BloomLevel")

	if materiaID := r.URL.Query().Get("materia_id"); materiaID != "" {
		query = query.Where("materia_id = ?", materiaID)
	}

	if activoParam := r.URL.Query().Get("activo"); activoParam != "" {
		if activo, err := strconv.ParseBool(activoParam); err == nil {
			query = query.Where("activo = ?", activo)
		}
	}

	if err := query.Order("orden ASC, codigo ASC").Find(&oas).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oas)
}

// GetObjetivoAprendizaje godoc
// @Summary Get learning objective by ID
// @Description Retrieve a specific learning objective with all Bloom-level details
// @Tags Educational
// @Produce json
// @Param id path int true "Learning Objective ID"
// @Success 200 {object} models.ObjetivoAprendizaje
// @Failure 404 {object} map[string]interface{}
// @Router /api/objetivos-aprendizaje/{id} [get]
func GetObjetivoAprendizaje(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var oa models.ObjetivoAprendizaje

	if err := db.DB.Preload("BloomObjectives.BloomLevel").
		Preload("Materia").
		First(&oa, id).Error; err != nil {
		http.Error(w, "Objetivo de Aprendizaje not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oa)
}

// CreateObjetivoAprendizaje godoc
// @Summary Create a new learning objective with all Bloom levels
// @Description Create a complete OA with specific objectives for all 6 Bloom taxonomy levels
// @Tags Educational
// @Accept json
// @Produce json
// @Param oa body CreateOARequest true "Complete OA data including all 6 Bloom objectives"
// @Success 201 {object} models.ObjetivoAprendizaje
// @Failure 400 {object} map[string]interface{}
// @Router /api/objetivos-aprendizaje [post]
func CreateObjetivoAprendizaje(w http.ResponseWriter, r *http.Request) {
	var req CreateOARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that we have objectives for all 6 Bloom levels
	if len(req.BloomObjectives) != 6 {
		http.Error(w, "Must provide objectives for all 6 Bloom levels", http.StatusBadRequest)
		return
	}

	// Start transaction
	tx := db.DB.Begin()

	// Create base OA
	oa := models.ObjetivoAprendizaje{
		MateriaID:   req.MateriaID,
		Codigo:      req.Codigo,
		Titulo:      req.Titulo,
		Descripcion: req.Descripcion,
		Orden:       req.Orden,
		Activo:      true,
	}

	if err := tx.Create(&oa).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create all 6 Bloom objectives
	for _, bloomObj := range req.BloomObjectives {
		oaBloom := models.OABloomObjective{
			OAID:                   oa.ID,
			BloomLevelID:           bloomObj.BloomLevelID,
			ObjetivoEspecifico:     bloomObj.ObjetivoEspecifico,
			IndicadoresLogro:       models.StringArray(bloomObj.IndicadoresLogro),
			TipoActividadSugerida:  bloomObj.TipoActividadSugerida,
			ComplejidadEstimada:    bloomObj.ComplejidadEstimada,
		}

		if err := tx.Create(&oaBloom).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	tx.Commit()

	// Reload with all relationships
	db.DB.Preload("BloomObjectives.BloomLevel").
		Preload("Materia").
		First(&oa, oa.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(oa)
}

// UpdateObjetivoAprendizaje godoc
// @Summary Update a learning objective
// @Description Update an existing learning objective (base fields only, use separate endpoint for Bloom objectives)
// @Tags Educational
// @Accept json
// @Produce json
// @Param id path int true "Learning Objective ID"
// @Param oa body models.ObjetivoAprendizaje true "Updated OA data"
// @Success 200 {object} models.ObjetivoAprendizaje
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/objetivos-aprendizaje/{id} [put]
func UpdateObjetivoAprendizaje(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var oa models.ObjetivoAprendizaje

	if err := db.DB.First(&oa, id).Error; err != nil {
		http.Error(w, "Objetivo de Aprendizaje not found", http.StatusNotFound)
		return
	}

	var updates models.ObjetivoAprendizaje
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update only base fields
	oa.Titulo = updates.Titulo
	oa.Descripcion = updates.Descripcion
	oa.Orden = updates.Orden
	oa.Activo = updates.Activo

	if err := db.DB.Save(&oa).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Reload with relationships
	db.DB.Preload("BloomObjectives.BloomLevel").First(&oa, oa.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oa)
}

// ============================================================================
// STUDENT PROGRESS TRACKING
// ============================================================================

// RegisterProgressRequest represents a progress update
type RegisterProgressRequest struct {
	UserID             uint     `json:"user_id"`
	OABloomObjectiveID uint     `json:"oa_bloom_objective_id"`
	Estado             string   `json:"estado"` // no_iniciado, en_proceso, logrado, dominado
	PorcentajeLogro    int      `json:"porcentaje_logro"`
	TipoEvento         string   `json:"tipo_evento"` // evaluacion, practica, diagnostico, repaso
	PuntajeObtenido    *float64 `json:"puntaje_obtenido"`
	PuntajeMaximo      *float64 `json:"puntaje_maximo"`
	Notas              string   `json:"notas"`
}

// RegisterProgress godoc
// @Summary Register student progress on an OA-Bloom objective
// @Description Record a student's progress on a specific Bloom-level objective
// @Tags Progress
// @Accept json
// @Produce json
// @Param progress body RegisterProgressRequest true "Progress data"
// @Success 201 {object} models.StudentOAProgress
// @Failure 400 {object} map[string]interface{}
// @Router /api/progress [post]
func RegisterProgress(w http.ResponseWriter, r *http.Request) {
	var req RegisterProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate estado
	validEstados := map[string]bool{
		"no_iniciado": true,
		"en_proceso":  true,
		"logrado":     true,
		"dominado":    true,
	}
	if !validEstados[req.Estado] {
		http.Error(w, "Invalid estado value", http.StatusBadRequest)
		return
	}

	tx := db.DB.Begin()

	// Upsert current progress
	var progress models.StudentOAProgress
	result := tx.Where("user_id = ? AND oa_bloom_objective_id = ?",
		req.UserID, req.OABloomObjectiveID).First(&progress)

	now := db.DB.NowFunc()
	if result.Error != nil {
		// Create new record
		progress = models.StudentOAProgress{
			UserID:              req.UserID,
			OABloomObjectiveID:  req.OABloomObjectiveID,
			Estado:              req.Estado,
			PorcentajeLogro:     req.PorcentajeLogro,
			Intentos:            1,
			UltimaActividadFecha: &now,
			Notas:               req.Notas,
		}
		if err := tx.Create(&progress).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Update existing record
		progress.Estado = req.Estado
		progress.PorcentajeLogro = req.PorcentajeLogro
		progress.Intentos++
		progress.UltimaActividadFecha = &now
		progress.Notas = req.Notas
		if err := tx.Save(&progress).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Create history entry
	history := models.StudentOAHistory{
		UserID:             req.UserID,
		OABloomObjectiveID: req.OABloomObjectiveID,
		Estado:             req.Estado,
		PorcentajeLogro:    &req.PorcentajeLogro,
		TipoEvento:         req.TipoEvento,
		PuntajeObtenido:    req.PuntajeObtenido,
		PuntajeMaximo:      req.PuntajeMaximo,
		Notas:              req.Notas,
	}

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx.Commit()

	// Gamification and unlock events (after successful commit)
	go func() {
		// Update streak
		gamificationService.UpdateStreak(req.UserID)

		// Award XP based on estado
		xpAmount := 0
		switch req.Estado {
		case "logrado":
			xpAmount = 30
		case "dominado":
			xpAmount = 50
		}
		if xpAmount > 0 {
			gamificationService.AddXP(req.UserID, xpAmount, "oa_progress")
			gamificationService.AddCoins(req.UserID, xpAmount/5, "oa_progress")
		}

		// Check for unlocks if estado is "dominado" or "logrado"
		if req.Estado == "dominado" || req.Estado == "logrado" {
			// Get OA details for event
			var oaBloom models.OABloomObjective
			db.DB.Preload("OA").First(&oaBloom, req.OABloomObjectiveID)

			// Trigger unlock check for OA completion
			unlockService.CheckAndUnlock(req.UserID, services.UnlockEvent{
				Type: "oa_complete",
				Key:  fmt.Sprintf("oa_%d", oaBloom.OAID),
				Data: map[string]interface{}{
					"materia_id": oaBloom.OA.MateriaID,
				},
			})

			// Trigger unlock check for Bloom mastery
			unlockService.CheckAndUnlock(req.UserID, services.UnlockEvent{
				Type: "bloom_mastery",
				Key:  fmt.Sprintf("bloom_%d_oa_%d", oaBloom.BloomLevelID, oaBloom.OAID),
				Data: map[string]interface{}{
					"bloom_level": oaBloom.BloomLevelID,
				},
			})
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(progress)
}

// GetStudentProgress godoc
// @Summary Get student progress
// @Description Retrieve all progress records for a specific student
// @Tags Progress
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} models.StudentOAProgress
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/{user_id} [get]
func GetStudentProgress(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	var progress []models.StudentOAProgress

	if err := db.DB.Where("user_id = ?", userID).
		Preload("OABloomObjective.OA.Materia").
		Preload("OABloomObjective.BloomLevel").
		Order("updated_at DESC").
		Find(&progress).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GetProgressHistory godoc
// @Summary Get student progress history
// @Description Retrieve complete history of progress changes for a student
// @Tags Progress
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Limit number of records (default 100)"
// @Success 200 {array} models.StudentOAHistory
// @Failure 500 {object} map[string]interface{}
// @Router /api/progress/{user_id}/history [get]
func GetProgressHistory(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	limit := 100
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			limit = parsedLimit
		}
	}

	var history []models.StudentOAHistory
	if err := db.DB.Where("user_id = ?", userID).
		Preload("OABloomObjective.OA.Materia").
		Preload("OABloomObjective.BloomLevel").
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
