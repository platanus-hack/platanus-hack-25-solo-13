package models

import (
	"time"

	"gorm.io/datatypes"
)

// LearningPlan represents a personalized learning plan generated for a user
type LearningPlan struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	UserID               uint      `json:"user_id" gorm:"not null"`
	OABloomObjectiveID   uint      `json:"oa_bloom_objective_id" gorm:"not null"`
	Titulo               string    `json:"titulo" gorm:"size:500;not null"`
	Descripcion          string    `json:"descripcion" gorm:"type:text"`
	TiempoEstimadoMin    int        `json:"tiempo_estimado_minutos" gorm:"column:tiempo_estimado_minutos;default:0"`
	Estado               string     `json:"estado" gorm:"size:50;not null;default:'generando'"`
	ErrorMensaje         string     `json:"error_mensaje,omitempty" gorm:"type:text"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	// Completion tracking fields
	Completado           bool       `json:"completado" gorm:"default:false;not null"`
	FechaInicio          *time.Time `json:"fecha_inicio,omitempty"`
	FechaCompletado      *time.Time `json:"fecha_completado,omitempty"`
	ProgresoActual       int        `json:"progreso_actual" gorm:"default:0;not null"`
	TotalSlides          int        `json:"total_slides" gorm:"default:0;not null"`

	// Relationships
	User             User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OABloomObjective OABloomObjective `json:"oa_bloom_objective,omitempty" gorm:"foreignKey:OABloomObjectiveID"`
	Components       []LearningPlanComponent `json:"components,omitempty" gorm:"foreignKey:LearningPlanID;constraint:OnDelete:CASCADE;"`
}

// TableName overrides the default table name
func (LearningPlan) TableName() string {
	return "learning_plans"
}

// LearningPlanComponent represents a single teaching component (slide) in a learning plan
type LearningPlanComponent struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	LearningPlanID      uint           `json:"learning_plan_id" gorm:"not null"`
	Orden               int            `json:"orden" gorm:"not null"`
	TipoComponente      string         `json:"tipo_componente" gorm:"size:100;not null"`
	ObjetivoEspecifico  string         `json:"objetivo_especifico" gorm:"type:text;not null"`
	TiempoEstimadoMin   int            `json:"tiempo_estimado_minutos" gorm:"column:tiempo_estimado_minutos;default:0"`
	Estado              string         `json:"estado" gorm:"size:50;not null;default:'pendiente'"`
	ContenidoProps      datatypes.JSON `json:"contenido_props,omitempty" gorm:"type:jsonb"`
	ErrorMensaje        string         `json:"error_mensaje,omitempty" gorm:"type:text"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`

	// Relationships
	LearningPlan LearningPlan `json:"learning_plan,omitempty" gorm:"foreignKey:LearningPlanID"`
}

// TableName overrides the default table name
func (LearningPlanComponent) TableName() string {
	return "learning_plan_components"
}

// Constants for LearningPlan states
const (
	LearningPlanEstadoGenerando = "generando"
	LearningPlanEstadoGenerado  = "generado"
	LearningPlanEstadoError     = "error"
)

// Constants for LearningPlanComponent states
const (
	ComponentEstadoPendiente  = "pendiente"
	ComponentEstadoGenerando  = "generando"
	ComponentEstadoGenerado   = "generado"
	ComponentEstadoError      = "error"
)

// Available teaching component type - flexible blocks-based component
const (
	ComponentTipoExplainAndExplore = "ExplainAndExploreSlide"
)

// AvailableComponentTypes returns the list of available teaching component types
func AvailableComponentTypes() []string {
	return []string{
		ComponentTipoExplainAndExplore,
	}
}

// IsValidComponentType checks if a component type is valid
func IsValidComponentType(tipo string) bool {
	for _, validType := range AvailableComponentTypes() {
		if tipo == validType {
			return true
		}
	}
	return false
}
