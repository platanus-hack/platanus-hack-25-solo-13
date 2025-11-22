package generator

import (
	"time"

	"github.com/lib/pq"
)

// OABloomObjective representa un objetivo de aprendizaje con nivel de Bloom
type OABloomObjective struct {
	ID                    uint
	OAID                  uint
	BloomLevelID          uint
	BloomLevelNumero      int
	BloomLevelNombre      string
	ObjetivoEspecifico    string
	IndicadoresLogro      pq.StringArray
	TipoActividadSugerida string
	ComplejidadEstimada   int
	OATitulo              string
	OADescripcion         string
	MateriaNombre         string
	CursoNombre           string
}

// Question representa una pregunta para insertar en BD
type Question struct {
	OABloomObjectiveID uint
	Tipo               string
	TipoUso            string
	QuestionData       []byte // JSON
	ValidationData     []byte // JSON
	DificultadRelativa int
	Tags               pq.StringArray
}

// FailedQuestion representa una pregunta que falló al generarse
type FailedQuestion struct {
	OABloomObjectiveID uint      `json:"oa_bloom_objective_id"`
	Tipo               string    `json:"tipo"`
	Dificultad         int       `json:"dificultad"`
	Error              string    `json:"error"`
	Timestamp          time.Time `json:"timestamp"`
}

// Stats representa estadísticas de generación
type Stats struct {
	TotalAttempts  int
	SuccessCount   int
	FailCount      int
	FailedQuestions []FailedQuestion
	TypeCounts     map[string]int
	StartTime      time.Time
	EndTime        time.Time
}

// AddSuccess incrementa contador de éxitos
func (s *Stats) AddSuccess(tipo string) {
	s.SuccessCount++
	s.TypeCounts[tipo]++
}

// AddFail agrega una pregunta fallida
func (s *Stats) AddFail(failed FailedQuestion) {
	s.FailCount++
	s.FailedQuestions = append(s.FailedQuestions, failed)
}
