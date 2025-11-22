package models

import (
	"time"

	"gorm.io/datatypes"
)

// StudentProfile represents a student's adaptive learning profile
type StudentProfile struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	Edad        *int           `json:"edad"`
	CursoActual string         `json:"curso_actual" gorm:"size:50"`
	ProfileData datatypes.JSON `json:"profile_data" gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ProfileHistory represents a snapshot of a student's profile at a point in time
type ProfileHistory struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Snapshot  datatypes.JSON `json:"snapshot" gorm:"type:jsonb;not null"`
	Evento    string         `json:"evento" gorm:"size:100"`
	CreatedAt time.Time      `json:"created_at"`

	// Relationship
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ProfileDataStructure defines the expected structure of profile_data JSON
// This is for documentation and type safety in code (not enforced by GORM)
type ProfileDataStructure struct {
	ConocimientoPrevio struct {
		Lectura     *NivelConocimiento `json:"lectura,omitempty"`
		Escritura   *NivelConocimiento `json:"escritura,omitempty"`
		Matematicas *NivelConocimiento `json:"matematicas,omitempty"`
	} `json:"conocimiento_previo,omitempty"`

	PerfilCognitivo struct {
		MemoriaTrabajo         int    `json:"memoria_trabajo,omitempty"`          // 1-10
		RazonamientoInductivo  int    `json:"razonamiento_inductivo,omitempty"`   // 1-10
		EstiloCognitivo        string `json:"estilo_cognitivo,omitempty"`         // "reflexivo", "impulsivo", "analitico"
		CargaCognitivaTolerada string `json:"carga_cognitiva_tolerada,omitempty"` // "bajo", "medio", "alto"
	} `json:"perfil_cognitivo,omitempty"`

	PreferenciasAprendizaje struct {
		FormatoPreferido     string   `json:"formato_preferido,omitempty"`      // "visual", "texto", "interactivo"
		TipoActividad        []string `json:"tipo_actividad,omitempty"`         // ["proyectos", "juegos", "lecturas"]
		CanalPreferido       string   `json:"canal_preferido,omitempty"`        // "video", "audio", "animaciones"
	} `json:"preferencias_aprendizaje,omitempty"`

	Motivacion struct {
		Intrinseca      int    `json:"intrinseca,omitempty"`        // 1-5
		Extrinseca      int    `json:"extrinseca,omitempty"`        // 1-5
		InteresActual   string `json:"interes_actual,omitempty"`    // "alto", "medio", "bajo"
		OrientacionMeta string `json:"orientacion_meta,omitempty"`  // "maestria", "desempeño", "evitar_fracaso"
	} `json:"motivacion,omitempty"`

	Autoeficacia struct {
		General             int `json:"general,omitempty"`               // 1-5
		ConfianzaResolutiva int `json:"confianza_resolutiva,omitempty"`  // 1-5
	} `json:"autoeficacia,omitempty"`

	Autonomia struct {
		Nivel            string   `json:"nivel,omitempty"`             // "bajo", "medio", "alto"
		GestionaTiempo   bool     `json:"gestiona_tiempo,omitempty"`
		Estrategias      []string `json:"estrategias,omitempty"`       // ["resumen", "mapas_conceptuales"]
	} `json:"autonomia,omitempty"`

	InteresesPersonales struct {
		Temas           []string `json:"temas,omitempty"`             // ["videojuegos", "musica", "tecnologia"]
		ProfesionSoñada string   `json:"profesion_soñada,omitempty"`
	} `json:"intereses_personales,omitempty"`

	UltimaActualizacion string `json:"ultima_actualizacion,omitempty"` // ISO 8601 timestamp
}

// NivelConocimiento represents the knowledge level for a specific topic
type NivelConocimiento struct {
	Nivel  int    `json:"nivel"`  // 0: sin nociones, 1: basico, 2: intermedio, 3: avanzado, 4: experto
	Fuente string `json:"fuente"` // "diagnostico_inicial", "inferido", "prueba", "docente"
}
