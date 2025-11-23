package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/platanus-hack-25/lumera_app/internal/db"
	"github.com/platanus-hack-25/lumera_app/internal/middleware"
	"github.com/platanus-hack-25/lumera_app/internal/models"
	"github.com/platanus-hack-25/lumera_app/internal/services"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// GenerateLearningPlanRequest es el payload para generar un plan
type GenerateLearningPlanRequest struct {
	OABloomObjectiveID uint `json:"oa_bloom_objective_id"`
}

// GenerateLearningPlanHandler genera un nuevo plan de aprendizaje para el usuario
// POST /api/learning-plans/generate
func GenerateLearningPlanHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req GenerateLearningPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Verificar si ya existe un plan para este usuario y OA
	var existingPlan models.LearningPlan
	err := db.DB.Preload("Components").
		Where("user_id = ? AND oa_bloom_objective_id = ?", userID, req.OABloomObjectiveID).
		First(&existingPlan).Error

	if err == nil {
		// Plan already exists, return it
		log.Printf("✓ Learning plan already exists for user %d and OA %d", userID, req.OABloomObjectiveID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingPlan)
		return
	}

	// Obtener datos del OA para el contexto
	var oaBloomObjective models.OABloomObjective
	if err := db.DB.Preload("OA").Preload("OA.Materia").Preload("BloomLevel").
		First(&oaBloomObjective, req.OABloomObjectiveID).Error; err != nil {
		log.Printf("Error loading OA: %v", err)
		http.Error(w, `{"error":"objective not found"}`, http.StatusNotFound)
		return
	}

	// Fetch student profile for personalization
	var profile models.StudentProfile
	var interesesPersonales []string
	var profesionSoñada string
	var formatoPreferido string
	var tipoActividad []string
	var canalPreferido string

	if err := db.DB.Where("user_id = ?", userID).First(&profile).Error; err == nil {
		// Parse profile_data JSONB field
		var profileData models.ProfileDataStructure
		if err := profile.ProfileData.Scan(&profileData); err == nil {
			interesesPersonales = profileData.InteresesPersonales.Temas
			profesionSoñada = profileData.InteresesPersonales.ProfesionSoñada
			formatoPreferido = profileData.PreferenciasAprendizaje.FormatoPreferido
			tipoActividad = profileData.PreferenciasAprendizaje.TipoActividad
			canalPreferido = profileData.PreferenciasAprendizaje.CanalPreferido
		}
	}

	// Construir contexto para OpenAI
	oaContext := services.OAContext{
		MateriaNombre:      oaBloomObjective.OA.Materia.Nombre,
		MateriaDescripcion: oaBloomObjective.OA.Materia.Descripcion,
		CursoNombre:        oaBloomObjective.OA.Materia.Nombre, // TODO: Get actual curso
		OATitulo:           oaBloomObjective.OA.Titulo,
		OADescripcion:      oaBloomObjective.OA.Descripcion,
		BloomLevelNombre:   oaBloomObjective.BloomLevel.Nombre,
		BloomLevelNumero:   oaBloomObjective.BloomLevel.Nivel,
		BloomDescripcion:   oaBloomObjective.BloomLevel.Descripcion,
		ObjetivoEspecifico: oaBloomObjective.ObjetivoEspecifico,
		IndicadoresLogro:   oaBloomObjective.IndicadoresLogro,
		// Student personalization
		InteresesPersonales: interesesPersonales,
		ProfesionSoñada:     profesionSoñada,
		FormatoPreferido:    formatoPreferido,
		TipoActividad:       tipoActividad,
		CanalPreferido:      canalPreferido,
	}

	// Generar estructura del plan con OpenAI
	planStructure, err := services.GenerateLearningPlanStructure(oaContext)
	if err != nil {
		log.Printf("Error generating plan structure: %v", err)
		http.Error(w, `{"error":"failed to generate learning plan"}`, http.StatusInternalServerError)
		return
	}

	// Calcular tiempo total estimado
	totalTime := 0
	for _, comp := range planStructure.Componentes {
		totalTime += comp.TiempoEstimadoMin
	}

	// Crear el plan en la base de datos
	plan := models.LearningPlan{
		UserID:             userID,
		OABloomObjectiveID: req.OABloomObjectiveID,
		Titulo:             planStructure.Titulo,
		Descripcion:        planStructure.Descripcion,
		TiempoEstimadoMin:  totalTime,
		Estado:             models.LearningPlanEstadoGenerando,
	}

	if err := db.DB.Create(&plan).Error; err != nil {
		log.Printf("Error creating plan: %v", err)
		http.Error(w, `{"error":"failed to create plan"}`, http.StatusInternalServerError)
		return
	}

	// Crear componentes y generar su contenido inmediatamente
	var components []models.LearningPlanComponent
	totalComponents := len(planStructure.Componentes)

	for i, compStruct := range planStructure.Componentes {
		component := models.LearningPlanComponent{
			LearningPlanID:     plan.ID,
			Orden:              i + 1,
			TipoComponente:     compStruct.Tipo,
			ObjetivoEspecifico: compStruct.ObjetivoEspecifico,
			TiempoEstimadoMin:  compStruct.TiempoEstimadoMin,
			Estado:             models.ComponentEstadoGenerando,
		}

		if err := db.DB.Create(&component).Error; err != nil {
			log.Printf("Error creating component: %v", err)
			continue
		}

		// Generar contenido del componente inmediatamente
		log.Printf("⏳ Generating content for component %d/%d (%s)...", i+1, totalComponents, compStruct.Tipo)

		content, err := services.GenerateComponentContent(
			component.TipoComponente,
			oaContext,
			component.ObjetivoEspecifico,
		)

		if err != nil {
			log.Printf("❌ Error generating content for component %d: %v", i+1, err)
			component.Estado = models.ComponentEstadoError
			component.ErrorMensaje = err.Error()
			db.DB.Save(&component)
			components = append(components, component)
			continue
		}

		// Guardar contenido generado
		contentJSON, _ := json.Marshal(content)
		component.ContenidoProps = datatypes.JSON(contentJSON)
		component.Estado = models.ComponentEstadoGenerado
		component.ErrorMensaje = ""
		db.DB.Save(&component)

		log.Printf("✅ Content generated for component %d/%d", i+1, totalComponents)
		components = append(components, component)
	}

	// Actualizar estado del plan a generado (todo está listo)
	plan.Estado = models.LearningPlanEstadoGenerado
	plan.Components = components
	plan.TotalSlides = len(components) // Set total slides for completion tracking
	db.DB.Save(&plan)

	log.Printf("✅ Learning plan fully generated: %s (%d/%d components successful)", plan.Titulo, len(components), totalComponents)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GetLearningPlanByIDHandler obtiene un plan por ID
// GET /api/learning-plans/{id}
func GetLearningPlanByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	planIDStr := chi.URLParam(r, "id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid plan ID"}`, http.StatusBadRequest)
		return
	}

	var plan models.LearningPlan
	if err := db.DB.Preload("Components").
		Where("id = ? AND user_id = ?", planID, userID).
		First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"plan not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GetLearningPlanByOAHandler obtiene el plan del usuario para un OA específico
// GET /api/learning-plans/by-oa/{oa_bloom_objective_id}
func GetLearningPlanByOAHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	oaIDStr := chi.URLParam(r, "oa_bloom_objective_id")
	oaID, err := strconv.ParseUint(oaIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid OA ID"}`, http.StatusBadRequest)
		return
	}

	var plan models.LearningPlan
	if err := db.DB.Preload("Components").
		Where("user_id = ? AND oa_bloom_objective_id = ?", userID, oaID).
		First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"plan not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// GenerateComponentContentHandler genera el contenido de un componente específico
// POST /api/learning-plans/{plan_id}/components/{component_id}/generate-content
func GenerateComponentContentHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	planIDStr := chi.URLParam(r, "plan_id")
	componentIDStr := chi.URLParam(r, "component_id")

	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid plan ID"}`, http.StatusBadRequest)
		return
	}

	componentID, err := strconv.ParseUint(componentIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid component ID"}`, http.StatusBadRequest)
		return
	}

	// Verificar que el plan pertenece al usuario
	var plan models.LearningPlan
	if err := db.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		http.Error(w, `{"error":"plan not found"}`, http.StatusNotFound)
		return
	}

	// Obtener el componente
	var component models.LearningPlanComponent
	if err := db.DB.Where("id = ? AND learning_plan_id = ?", componentID, planID).
		First(&component).Error; err != nil {
		http.Error(w, `{"error":"component not found"}`, http.StatusNotFound)
		return
	}

	// Si ya está generado, devolverlo
	if component.Estado == models.ComponentEstadoGenerado && component.ContenidoProps != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(component)
		return
	}

	// Marcar como generando
	component.Estado = models.ComponentEstadoGenerando
	db.DB.Save(&component)

	// Obtener datos del OA para el contexto
	var oaBloomObjective models.OABloomObjective
	if err := db.DB.Preload("OA").Preload("OA.Materia").Preload("BloomLevel").
		First(&oaBloomObjective, plan.OABloomObjectiveID).Error; err != nil {
		component.Estado = models.ComponentEstadoError
		component.ErrorMensaje = "Failed to load OA data"
		db.DB.Save(&component)
		http.Error(w, `{"error":"objective not found"}`, http.StatusInternalServerError)
		return
	}

	// Fetch student profile for personalization
	var profile models.StudentProfile
	var interesesPersonales []string
	var profesionSoñada string
	var formatoPreferido string
	var tipoActividad []string
	var canalPreferido string

	if err := db.DB.Where("user_id = ?", userID).First(&profile).Error; err == nil {
		// Parse profile_data JSONB field
		var profileData models.ProfileDataStructure
		if err := profile.ProfileData.Scan(&profileData); err == nil {
			interesesPersonales = profileData.InteresesPersonales.Temas
			profesionSoñada = profileData.InteresesPersonales.ProfesionSoñada
			formatoPreferido = profileData.PreferenciasAprendizaje.FormatoPreferido
			tipoActividad = profileData.PreferenciasAprendizaje.TipoActividad
			canalPreferido = profileData.PreferenciasAprendizaje.CanalPreferido
		}
	}

	// Construir contexto para OpenAI
	oaContext := services.OAContext{
		MateriaNombre:      oaBloomObjective.OA.Materia.Nombre,
		MateriaDescripcion: oaBloomObjective.OA.Materia.Descripcion,
		CursoNombre:        oaBloomObjective.OA.Materia.Nombre,
		OATitulo:           oaBloomObjective.OA.Titulo,
		OADescripcion:      oaBloomObjective.OA.Descripcion,
		BloomLevelNombre:   oaBloomObjective.BloomLevel.Nombre,
		BloomLevelNumero:   oaBloomObjective.BloomLevel.Nivel,
		BloomDescripcion:   oaBloomObjective.BloomLevel.Descripcion,
		ObjetivoEspecifico: oaBloomObjective.ObjetivoEspecifico,
		IndicadoresLogro:   oaBloomObjective.IndicadoresLogro,
		// Student personalization
		InteresesPersonales: interesesPersonales,
		ProfesionSoñada:     profesionSoñada,
		FormatoPreferido:    formatoPreferido,
		TipoActividad:       tipoActividad,
		CanalPreferido:      canalPreferido,
	}

	// Generar contenido
	content, err := services.GenerateComponentContent(
		component.TipoComponente,
		oaContext,
		component.ObjetivoEspecifico,
	)

	if err != nil {
		log.Printf("Error generating component content: %v", err)
		component.Estado = models.ComponentEstadoError
		component.ErrorMensaje = err.Error()
		db.DB.Save(&component)
		http.Error(w, `{"error":"failed to generate content"}`, http.StatusInternalServerError)
		return
	}

	// Guardar contenido
	contentJSON, _ := json.Marshal(content)
	component.ContenidoProps = datatypes.JSON(contentJSON)
	component.Estado = models.ComponentEstadoGenerado
	component.ErrorMensaje = ""
	db.DB.Save(&component)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(component)
}

// StartLearningPlanHandler marks a learning plan as started
// POST /api/learning-plans/{id}/start
func StartLearningPlanHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	planIDStr := chi.URLParam(r, "id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid plan ID"}`, http.StatusBadRequest)
		return
	}

	var plan models.LearningPlan
	if err := db.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"plan not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Only set fecha_inicio if not already set
	if plan.FechaInicio == nil {
		now := time.Now()
		plan.FechaInicio = &now
		if err := db.DB.Save(&plan).Error; err != nil {
			http.Error(w, `{"error":"failed to update plan"}`, http.StatusInternalServerError)
			return
		}
		log.Printf("Learning plan %d marked as started by user %d", planID, userID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}

// CompleteLearningPlanHandler marks a learning plan as completed
// POST /api/learning-plans/{id}/complete
func CompleteLearningPlanHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	planIDStr := chi.URLParam(r, "id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		http.Error(w, `{"error":"invalid plan ID"}`, http.StatusBadRequest)
		return
	}

	var plan models.LearningPlan
	if err := db.DB.Where("id = ? AND user_id = ?", planID, userID).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"error":"plan not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
		}
		return
	}

	// Mark as completed
	now := time.Now()
	plan.Completado = true
	plan.FechaCompletado = &now
	plan.ProgresoActual = plan.TotalSlides // All slides completed

	// Set fecha_inicio if not already set (edge case)
	if plan.FechaInicio == nil {
		plan.FechaInicio = &now
	}

	if err := db.DB.Save(&plan).Error; err != nil {
		http.Error(w, `{"error":"failed to update plan"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Learning plan %d marked as completed by user %d", planID, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plan)
}
