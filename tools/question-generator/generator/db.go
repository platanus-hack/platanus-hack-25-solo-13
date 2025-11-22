package generator

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB conecta a PostgreSQL usando GORM
func ConnectDB() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✓ Database connected successfully")
	return nil
}

// GetOABloomObjectives obtiene todos los objetivos con información completa
func GetOABloomObjectives() ([]OABloomObjective, error) {
	var objectives []OABloomObjective

	// Query compleja para traer toda la metadata necesaria
	query := `
		SELECT DISTINCT ON (oab.id)
			oab.id,
			oab.oa_id,
			oab.bloom_level_id,
			bl.nivel as bloom_level_numero,
			bl.nombre as bloom_level_nombre,
			oab.objetivo_especifico,
			oab.indicadores_logro,
			oab.tipo_actividad_sugerida,
			oab.complejidad_estimada,
			oa.titulo as oa_titulo,
			oa.descripcion as oa_descripcion,
			m.nombre as materia_nombre,
			COALESCE(c.nombre, 'Sin curso asignado') as curso_nombre
		FROM oa_bloom_objectives oab
		INNER JOIN bloom_levels bl ON oab.bloom_level_id = bl.id
		INNER JOIN objetivos_aprendizaje oa ON oab.oa_id = oa.id
		INNER JOIN materias m ON oa.materia_id = m.id
		LEFT JOIN curso_materias cm ON cm.materia_id = m.id
		LEFT JOIN cursos c ON c.id = cm.curso_id
		ORDER BY oab.id, m.nombre, c.nombre, oa.id, bl.nivel
	`

	err := DB.Raw(query).Scan(&objectives).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch objectives: %w", err)
	}

	log.Printf("✓ Fetched %d OA-Bloom objectives from database", len(objectives))
	return objectives, nil
}

// InsertQuestions inserta preguntas en la base de datos en lotes
func InsertQuestions(questions []Question) error {
	if len(questions) == 0 {
		return nil
	}

	batchSize := 50
	for i := 0; i < len(questions); i += batchSize {
		end := i + batchSize
		if end > len(questions) {
			end = len(questions)
		}
		batch := questions[i:end]

		err := DB.Transaction(func(tx *gorm.DB) error {
			for _, q := range batch {
				// Insertar usando raw SQL para mejor control
				result := tx.Exec(`
					INSERT INTO questions (
						oa_bloom_objective_id,
						tipo,
						tipo_uso,
						question_data,
						validation_data,
						dificultad_relativa,
						tags,
						activa,
						created_at,
						updated_at
					) VALUES (?, ?, ?, ?, ?, ?, ?, true, NOW(), NOW())
				`, q.OABloomObjectiveID, q.Tipo, q.TipoUso, q.QuestionData, q.ValidationData, q.DificultadRelativa, q.Tags)

				if result.Error != nil {
					return fmt.Errorf("failed to insert question: %w", result.Error)
				}
			}
			return nil
		})

		if err != nil {
			return err
		}

		log.Printf("✓ Inserted batch %d-%d questions", i+1, end)
	}

	log.Printf("✓ Successfully inserted %d questions", len(questions))
	return nil
}

// CountExistingQuestions cuenta cuántas preguntas ya existen para un OA-Bloom objective
func CountExistingQuestions(oaBloomObjectiveID uint) (int64, error) {
	var count int64
	err := DB.Model(&Question{}).Where("oa_bloom_objective_id = ?", oaBloomObjectiveID).Count(&count).Error
	return count, err
}
