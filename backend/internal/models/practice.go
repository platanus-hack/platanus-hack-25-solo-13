package models

import (
	"time"

	"gorm.io/datatypes"
)

// PracticeSession represents a practice session for a specific OA
type PracticeSession struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	UserID              uint           `json:"user_id" gorm:"not null;index"`
	OAID                uint           `json:"oa_id" gorm:"not null;index"`
	OABloomObjectiveID  uint           `json:"oa_bloom_objective_id" gorm:"not null"` // Specific OA+Bloom level
	BloomLevelInicial   int            `json:"bloom_level_inicial"`                  // Starting level
	BloomLevelFinal     *int           `json:"bloom_level_final"`                    // Final level after practice
	NumeroPreguntas     int            `json:"numero_preguntas" gorm:"default:10"`   // Number of questions (default 10)
	PreguntasRespondidas int           `json:"preguntas_respondidas" gorm:"default:0"`
	PreguntasCorrectas  int            `json:"preguntas_correctas" gorm:"default:0"`
	Estado              string         `json:"estado" gorm:"type:varchar(20);default:'en_progreso';index"` // en_progreso, completado
	Estrategia          datatypes.JSON `json:"estrategia" gorm:"type:jsonb"`         // Adaptive strategy data
	Resultado           datatypes.JSON `json:"resultado" gorm:"type:jsonb"`          // Final results and analysis
	StartedAt           time.Time      `json:"started_at"`
	CompletedAt         *time.Time     `json:"completed_at"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`

	// Relations
	User             User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	OA               ObjetivoAprendizaje `json:"oa,omitempty" gorm:"foreignKey:OAID"`
	OABloomObjective OABloomObjective `json:"oa_bloom_objective,omitempty" gorm:"foreignKey:OABloomObjectiveID"`
	Answers          []PracticeAnswer `json:"answers,omitempty" gorm:"foreignKey:SessionID"`
}

// PracticeAnswer represents a single answer in a practice session
type PracticeAnswer struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	SessionID          uint           `json:"session_id" gorm:"not null;index"`
	QuestionID         uint           `json:"question_id" gorm:"not null"`
	BloomLevelID       uint           `json:"bloom_level_id" gorm:"not null"` // Difficulty level of question
	UserAnswer         datatypes.JSON `json:"user_answer" gorm:"type:jsonb"`
	IsCorrect          *bool          `json:"is_correct"`
	Score              *float64       `json:"score"`              // Partial credit (0-1)
	TiempoSegundos     *int           `json:"tiempo_segundos"`
	CreatedAt          time.Time      `json:"created_at"`

	// Relations
	Session  PracticeSession `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	Question Question        `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}

// PracticeStrategy represents the adaptive strategy used during practice
type PracticeStrategy struct {
	NivelBloomActual     int      `json:"nivel_bloom_actual"`
	AciertosConsecutivos int      `json:"aciertos_consecutivos"`
	FallosConsecutivos   int      `json:"fallos_consecutivos"`
	AciertosPorNivel     map[int]int `json:"aciertos_por_nivel"`
	FallosPorNivel       map[int]int `json:"fallos_por_nivel"`
	PatronRespuestas     []string `json:"patron_respuestas"` // "C" = correct, "I" = incorrect
}

func (PracticeSession) TableName() string {
	return "practice_sessions"
}

func (PracticeAnswer) TableName() string {
	return "practice_answers"
}
