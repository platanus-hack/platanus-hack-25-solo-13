package generator

import "time"

// Avatar representa la configuración de un avatar a generar
type Avatar struct {
	Nombre        string        `json:"nombre"`
	Descripcion   string        `json:"descripcion"`
	Categoria     string        `json:"categoria"` // inicio, logro, compra
	Tier          int           `json:"tier"`      // 1-5 estrellas
	Rarity        string        `json:"rarity"`    // common, rare, epic, legendary
	PrecioPuntos  int           `json:"precio_puntos"`
	IsDefault     bool          `json:"is_default"`
	UnlockTrigger UnlockTrigger `json:"unlock_trigger"`
	DallePrompt   string        `json:"dalle_prompt"`
	ImageURL      string        `json:"image_url,omitempty"` // Se llena después de generar
}

// UnlockTrigger define cómo se desbloquea un avatar
type UnlockTrigger struct {
	TriggerType string                 `json:"trigger_type"` // default, oa_complete, bloom_mastery, streak, level, coins
	TriggerKey  string                 `json:"trigger_key,omitempty"`
	DisplayText string                 `json:"display_text,omitempty"`
	ExtraData   map[string]interface{} `json:"extra_data,omitempty"`
}

// GenerationResult representa el resultado de generar un avatar
type GenerationResult struct {
	Avatar    Avatar
	Success   bool
	Error     string
	Timestamp time.Time
}

// Stats mantiene estadísticas de generación
type Stats struct {
	TotalAttempts int
	SuccessCount  int
	FailCount     int
	Failed        []GenerationResult
	StartTime     time.Time
	EndTime       time.Time
}

// AddSuccess incrementa contador de éxitos
func (s *Stats) AddSuccess() {
	s.SuccessCount++
}

// AddFail registra un fallo
func (s *Stats) AddFail(result GenerationResult) {
	s.FailCount++
	s.Failed = append(s.Failed, result)
}
