package models

import (
	"time"

	"gorm.io/datatypes"
)

// DiagnosticSession represents an adaptive diagnostic session
type DiagnosticSession struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	UserID              uint           `json:"user_id" gorm:"not null"`
	MateriaID           uint           `json:"materia_id" gorm:"not null"`
	NumeroIntento       int            `json:"numero_intento" gorm:"default:1"`
	Estado              string         `json:"estado" gorm:"size:20;not null;default:en_progreso"`
	Estrategia          datatypes.JSON `json:"estrategia" gorm:"type:jsonb;not null"`
	PreguntasTotales    int            `json:"preguntas_totales" gorm:"default:0"`
	PreguntasCorrectas  int            `json:"preguntas_correctas" gorm:"default:0"`
	StartedAt           time.Time      `json:"started_at"`
	CompletedAt         *time.Time     `json:"completed_at"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Materia  Materia  `json:"materia,omitempty" gorm:"foreignKey:MateriaID"`
	Answers  []DiagnosticAnswer `json:"answers,omitempty" gorm:"foreignKey:SessionID"`
	Results  []DiagnosticResult `json:"results,omitempty" gorm:"foreignKey:SessionID"`
}

// TableName overrides the default table name
func (DiagnosticSession) TableName() string {
	return "diagnostic_sessions"
}

// DiagnosticAnswer represents an individual answer during a diagnostic session
type DiagnosticAnswer struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	SessionID            uint           `json:"session_id" gorm:"not null"`
	QuestionID           uint           `json:"question_id" gorm:"not null"`
	OABloomObjectiveID   uint           `json:"oa_bloom_objective_id" gorm:"not null"`
	BloomLevelID         uint           `json:"bloom_level_id" gorm:"not null"`
	UserAnswer           datatypes.JSON `json:"user_answer" gorm:"type:jsonb;not null"`
	IsCorrect            *bool          `json:"is_correct"`
	Score                *float64       `json:"score" gorm:"type:decimal(5,2)"`
	TiempoSegundos       *int           `json:"tiempo_segundos"`
	CreatedAt            time.Time      `json:"created_at"`

	// Relationships
	Session          DiagnosticSession `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	Question         Question          `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	OABloomObjective OABloomObjective  `json:"oa_bloom_objective,omitempty" gorm:"foreignKey:OABloomObjectiveID"`
}

// TableName overrides the default table name
func (DiagnosticAnswer) TableName() string {
	return "diagnostic_answers"
}

// DiagnosticResult represents consolidated results per OA after completing diagnostic
type DiagnosticResult struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	SessionID            uint      `json:"session_id" gorm:"not null"`
	OAID                 uint      `json:"oa_id" gorm:"not null"`
	NivelBloomDominado   int       `json:"nivel_bloom_dominado" gorm:"not null"`
	NivelBloomNombre     string    `json:"nivel_bloom_nombre" gorm:"size:20;not null"`
	PreguntasRespondidas int       `json:"preguntas_respondidas" gorm:"default:0"`
	PreguntasCorrectas   int       `json:"preguntas_correctas" gorm:"default:0"`
	PorcentajeAciertos   int       `json:"porcentaje_aciertos" gorm:"not null"`
	Recomendacion        string    `json:"recomendacion" gorm:"type:text"`
	CreatedAt            time.Time `json:"created_at"`

	// Relationships
	Session DiagnosticSession   `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	OA      ObjetivoAprendizaje `json:"oa,omitempty" gorm:"foreignKey:OAID"`
}

// TableName overrides the default table name
func (DiagnosticResult) TableName() string {
	return "diagnostic_results"
}

// AdaptiveStrategy represents the structure of the estrategia JSONB field
type AdaptiveStrategy struct {
	NivelBloomActual      int      `json:"nivel_bloom_actual"`
	OAsEvaluados          []uint   `json:"oas_evaluados"`
	AciertosConsecutivos  int      `json:"aciertos_consecutivos"`
	FallosConsecutivos    int      `json:"fallos_consecutivos"`
	PatronRespuestas      []string `json:"patron_respuestas"`
}
