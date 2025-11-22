package models

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
)

// StringArray is a custom type for PostgreSQL text arrays
type StringArray []string

// Scan implements sql.Scanner interface for GORM
func (a *StringArray) Scan(value interface{}) error {
	return (*pq.StringArray)(a).Scan(value)
}

// Value implements driver.Valuer interface for GORM
func (a StringArray) Value() (driver.Value, error) {
	return pq.StringArray(a).Value()
}

// Curso represents an educational level (e.g., "1° medio")
type Curso struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Nombre          string    `json:"nombre" gorm:"size:100;not null"`
	Codigo          string    `json:"codigo" gorm:"size:20;uniqueIndex"`
	NivelEducativo  string    `json:"nivel_educativo" gorm:"size:50"`
	Descripcion     string    `json:"descripcion" gorm:"type:text"`
	Activo          bool      `json:"activo" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	Materias []Materia `json:"materias,omitempty" gorm:"many2many:curso_materias;"`
}

// TableName overrides the default table name
func (Curso) TableName() string {
	return "cursos"
}

// Materia represents a subject (e.g., "Lenguaje", "Matemáticas")
type Materia struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Nombre      string    `json:"nombre" gorm:"size:100;not null"`
	Codigo      string    `json:"codigo" gorm:"size:20;uniqueIndex"`
	Descripcion string    `json:"descripcion" gorm:"type:text"`
	Color       string    `json:"color" gorm:"size:7"` // Hex color for UI
	Activo      bool      `json:"activo" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Cursos []Curso `json:"cursos,omitempty" gorm:"many2many:curso_materias;"`
	OAs    []ObjetivoAprendizaje `json:"oas,omitempty" gorm:"foreignKey:MateriaID"`
}

// TableName overrides the default table name
func (Materia) TableName() string {
	return "materias"
}

// CursoMateria represents the many-to-many relationship
type CursoMateria struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CursoID        uint      `json:"curso_id" gorm:"not null"`
	MateriaID      uint      `json:"materia_id" gorm:"not null;uniqueIndex:idx_curso_materia"`
	HorasSemanales *int      `json:"horas_semanales"`
	CreatedAt      time.Time `json:"created_at"`

	// Relationships
	Curso   Curso   `json:"curso,omitempty" gorm:"foreignKey:CursoID"`
	Materia Materia `json:"materia,omitempty" gorm:"foreignKey:MateriaID"`
}

// TableName overrides the default table name
func (CursoMateria) TableName() string {
	return "curso_materias"
}

// BloomLevel represents a level in Bloom's Taxonomy
type BloomLevel struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	Nivel        int         `json:"nivel" gorm:"not null;uniqueIndex"`
	Nombre       string      `json:"nombre" gorm:"size:50;not null"`
	NombreEn     string      `json:"nombre_en" gorm:"size:50"`
	Descripcion  string      `json:"descripcion" gorm:"type:text"`
	VerbosAccion StringArray `json:"verbos_accion" gorm:"type:text[]"`
	Color        string      `json:"color" gorm:"size:7"`
	CreatedAt    time.Time   `json:"created_at"`
}

// TableName overrides the default table name
func (BloomLevel) TableName() string {
	return "bloom_levels"
}

// ObjetivoAprendizaje represents a base learning objective
type ObjetivoAprendizaje struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	MateriaID   uint      `json:"materia_id" gorm:"not null"`
	Codigo      string    `json:"codigo" gorm:"size:20;not null;uniqueIndex:idx_materia_codigo"`
	Titulo      string    `json:"titulo" gorm:"size:255;not null"`
	Descripcion string    `json:"descripcion" gorm:"type:text"`
	Orden       *int      `json:"orden"`
	Activo      bool      `json:"activo" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Materia          Materia            `json:"materia,omitempty" gorm:"foreignKey:MateriaID"`
	BloomObjectives  []OABloomObjective `json:"bloom_objectives,omitempty" gorm:"foreignKey:OAID"`
}

// TableName overrides the default table name
func (ObjetivoAprendizaje) TableName() string {
	return "objetivos_aprendizaje"
}

// OABloomObjective represents an OA broken down by Bloom level
type OABloomObjective struct {
	ID                      uint        `json:"id" gorm:"primaryKey"`
	OAID                    uint        `json:"oa_id" gorm:"not null;uniqueIndex:idx_oa_bloom"`
	BloomLevelID            uint        `json:"bloom_level_id" gorm:"not null;uniqueIndex:idx_oa_bloom"`
	ObjetivoEspecifico      string      `json:"objetivo_especifico" gorm:"type:text;not null"`
	IndicadoresLogro        StringArray `json:"indicadores_logro" gorm:"type:text[]"`
	TipoActividadSugerida   string      `json:"tipo_actividad_sugerida" gorm:"size:50"`
	ComplejidadEstimada     *int        `json:"complejidad_estimada"` // 1-10
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`

	// Relationships
	OA         ObjetivoAprendizaje `json:"oa,omitempty" gorm:"foreignKey:OAID"`
	BloomLevel BloomLevel          `json:"bloom_level,omitempty" gorm:"foreignKey:BloomLevelID"`
}

// TableName overrides the default table name
func (OABloomObjective) TableName() string {
	return "oa_bloom_objectives"
}

// StudentOAProgress represents current progress on an OA-Bloom objective
type StudentOAProgress struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	UserID              uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_oa_bloom"`
	OABloomObjectiveID  uint       `json:"oa_bloom_objective_id" gorm:"not null;uniqueIndex:idx_user_oa_bloom"`
	Estado              string     `json:"estado" gorm:"size:20;not null"` // no_iniciado, en_proceso, logrado, dominado
	PorcentajeLogro     int        `json:"porcentaje_logro" gorm:"default:0"`
	Intentos            int        `json:"intentos" gorm:"default:0"`
	UltimaActividadFecha *time.Time `json:"ultima_actividad_fecha"`
	Notas               string     `json:"notas" gorm:"type:text"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`

	// Relationships
	User            User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OABloomObjective OABloomObjective `json:"oa_bloom_objective,omitempty" gorm:"foreignKey:OABloomObjectiveID"`
}

// TableName overrides the default table name
func (StudentOAProgress) TableName() string {
	return "student_oa_progress"
}

// StudentOAHistory represents history of progress changes
type StudentOAHistory struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	UserID             uint      `json:"user_id" gorm:"not null"`
	OABloomObjectiveID uint      `json:"oa_bloom_objective_id" gorm:"not null"`
	Estado             string    `json:"estado" gorm:"size:20;not null"`
	PorcentajeLogro    *int      `json:"porcentaje_logro"`
	TipoEvento         string    `json:"tipo_evento" gorm:"size:50"` // evaluacion, practica, diagnostico, repaso
	PuntajeObtenido    *float64  `json:"puntaje_obtenido" gorm:"type:decimal(5,2)"`
	PuntajeMaximo      *float64  `json:"puntaje_maximo" gorm:"type:decimal(5,2)"`
	Notas              string    `json:"notas" gorm:"type:text"`
	CreatedAt          time.Time `json:"created_at"`

	// Relationships
	User            User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OABloomObjective OABloomObjective `json:"oa_bloom_objective,omitempty" gorm:"foreignKey:OABloomObjectiveID"`
}

// TableName overrides the default table name
func (StudentOAHistory) TableName() string {
	return "student_oa_history"
}
