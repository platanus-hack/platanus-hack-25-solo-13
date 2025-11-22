package loader

import (
	"time"

	"github.com/lib/pq"
)

// CSVRecord representa una fila del CSV de entrada
type CSVRecord struct {
	Materia  string `csv:"materia"`
	Curso    string `csv:"curso"`
	Objetivo string `csv:"objetivo"`
}

// BloomObjective representa un subobjetivo generado por OpenAI
type BloomObjective struct {
	NivelBloom              string   `json:"nivel_bloom"`
	Objetivo                string   `json:"objetivo"`
	IndicadoresLogro        []string `json:"indicadores_logro"`
	TipoActividadSugerida   string   `json:"tipo_actividad_sugerida"`
	ComplejidadEstimada     int      `json:"complejidad_estimada"`
}

// OpenAIResponse es la respuesta completa de OpenAI (array de BloomObjective)
type OpenAIResponse []BloomObjective

// FailedOA representa un OA que falló al procesarse
type FailedOA struct {
	Materia    string    `json:"materia"`
	Curso      string    `json:"curso"`
	Objetivo   string    `json:"objetivo"`
	Error      string    `json:"error"`
	Timestamp  time.Time `json:"timestamp"`
}

// --- Modelos de BD (copiados del backend) ---

// ObjetivoAprendizaje representa la tabla objetivos_aprendizaje
type ObjetivoAprendizaje struct {
	ID          uint      `gorm:"primaryKey"`
	MateriaID   uint      `gorm:"not null"`
	Codigo      string    `gorm:"size:20"`
	Titulo      string    `gorm:"size:255;not null"`
	Descripcion string    `gorm:"type:text"`
	Orden       int
	Activo      bool      `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName especifica el nombre de la tabla
func (ObjetivoAprendizaje) TableName() string {
	return "objetivos_aprendizaje"
}

// OABloomObjective representa la tabla oa_bloom_objectives
type OABloomObjective struct {
	ID                      uint           `gorm:"primaryKey"`
	OAID                    uint           `gorm:"column:oa_id;not null"`
	BloomLevelID            uint           `gorm:"not null"`
	ObjetivoEspecifico      string         `gorm:"type:text;not null"`
	IndicadoresLogro        pq.StringArray `gorm:"type:text[]"`
	TipoActividadSugerida   string         `gorm:"size:50"`
	ComplejidadEstimada     int            `gorm:"check:complejidad_estimada >= 1 AND complejidad_estimada <= 10"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// TableName especifica el nombre de la tabla
func (OABloomObjective) TableName() string {
	return "oa_bloom_objectives"
}

// Materia representa la tabla materias
type Materia struct {
	ID          uint   `gorm:"primaryKey"`
	Nombre      string `gorm:"size:100;not null"`
	Codigo      string `gorm:"size:20;unique"`
	Descripcion string `gorm:"type:text"`
	Color       string `gorm:"size:7"`
	Activo      bool   `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName especifica el nombre de la tabla
func (Materia) TableName() string {
	return "materias"
}

// Curso representa la tabla cursos
type Curso struct {
	ID             uint   `gorm:"primaryKey"`
	Nombre         string `gorm:"size:100;not null"`
	Codigo         string `gorm:"size:20;unique"`
	NivelEducativo string `gorm:"size:50"`
	Descripcion    string `gorm:"type:text"`
	Activo         bool   `gorm:"default:true"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName especifica el nombre de la tabla
func (Curso) TableName() string {
	return "cursos"
}

// BloomLevel representa la tabla bloom_levels
type BloomLevel struct {
	ID           uint           `gorm:"primaryKey"`
	Nivel        int            `gorm:"unique;not null"`
	Nombre       string         `gorm:"size:50;not null"`
	NombreEN     string         `gorm:"column:nombre_en;size:50"`
	Descripcion  string         `gorm:"type:text"`
	VerbosAccion pq.StringArray `gorm:"type:text[]"`
	Color        string         `gorm:"size:7"`
	CreatedAt    time.Time
}

// TableName especifica el nombre de la tabla
func (BloomLevel) TableName() string {
	return "bloom_levels"
}
