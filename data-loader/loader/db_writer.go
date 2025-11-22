package loader

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBWriter maneja las escrituras a la base de datos
type DBWriter struct {
	db *gorm.DB
}

// NewDBWriter crea una nueva instancia del escritor de BD
func NewDBWriter(dsn string) (*DBWriter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("error conectando a BD: %w", err)
	}

	return &DBWriter{db: db}, nil
}

// Close cierra la conexión a la BD
func (w *DBWriter) Close() error {
	sqlDB, err := w.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// bloomNameToID mapea nombres de niveles de Bloom a sus IDs
var bloomNameToID = map[string]uint{
	"recordar":   1,
	"comprender": 2,
	"aplicar":    3,
	"analizar":   4,
	"evaluar":    5,
	"crear":      6,
}

// SaveOA guarda un OA completo con sus objetivos de Bloom en la BD
func (w *DBWriter) SaveOA(record CSVRecord, bloomObjectives OpenAIResponse) error {
	// Iniciar transacción
	return w.db.Transaction(func(tx *gorm.DB) error {
		// 1. Buscar ID de materia
		var materia Materia
		if err := tx.Where("LOWER(nombre) = LOWER(?)", record.Materia).First(&materia).Error; err != nil {
			return fmt.Errorf("materia '%s' no encontrada: %w", record.Materia, err)
		}

		// 2. Buscar ID de curso
		var curso Curso
		if err := tx.Where("LOWER(nombre) = LOWER(?)", record.Curso).First(&curso).Error; err != nil {
			return fmt.Errorf("curso '%s' no encontrado: %w", record.Curso, err)
		}

		// 3. Generar código automático único
		var count int64
		tx.Model(&ObjetivoAprendizaje{}).Where("materia_id = ?", materia.ID).Count(&count)
		codigo := fmt.Sprintf("OA-%d", count+1)

		// 4. Insertar Objetivo de Aprendizaje
		oa := ObjetivoAprendizaje{
			MateriaID:   materia.ID,
			Codigo:      codigo,
			Titulo:      truncate(record.Objetivo, 255),
			Descripcion: record.Objetivo,
			Orden:       int(count + 1),
			Activo:      true,
		}

		if err := tx.Create(&oa).Error; err != nil {
			return fmt.Errorf("error insertando OA: %w", err)
		}

		// 5. Insertar objetivos de Bloom
		for _, bloomObj := range bloomObjectives {
			bloomLevelID, ok := bloomNameToID[strings.ToLower(bloomObj.NivelBloom)]
			if !ok {
				// Si el nombre no coincide exactamente, intentar buscar en BD
				var bloomLevel BloomLevel
				if err := tx.Where("LOWER(nombre) = LOWER(?)", bloomObj.NivelBloom).First(&bloomLevel).Error; err != nil {
					return fmt.Errorf("nivel de Bloom '%s' no encontrado: %w", bloomObj.NivelBloom, err)
				}
				bloomLevelID = bloomLevel.ID
			}

			oaBloom := OABloomObjective{
				OAID:                    oa.ID,
				BloomLevelID:            bloomLevelID,
				ObjetivoEspecifico:      bloomObj.Objetivo,
				IndicadoresLogro:        bloomObj.IndicadoresLogro,
				TipoActividadSugerida:   bloomObj.TipoActividadSugerida,
				ComplejidadEstimada:     bloomObj.ComplejidadEstimada,
			}

			if err := tx.Create(&oaBloom).Error; err != nil {
				return fmt.Errorf("error insertando OA-Bloom nivel '%s': %w", bloomObj.NivelBloom, err)
			}
		}

		return nil
	})
}

// GetMaterias retorna todas las materias activas
func (w *DBWriter) GetMaterias() ([]Materia, error) {
	var materias []Materia
	if err := w.db.Where("activo = ?", true).Find(&materias).Error; err != nil {
		return nil, err
	}
	return materias, nil
}

// GetCursos retorna todos los cursos activos
func (w *DBWriter) GetCursos() ([]Curso, error) {
	var cursos []Curso
	if err := w.db.Where("activo = ?", true).Find(&cursos).Error; err != nil {
		return nil, err
	}
	return cursos, nil
}

// truncate recorta un string a una longitud máxima de caracteres (no bytes)
// Maneja correctamente caracteres UTF-8 multi-byte
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
