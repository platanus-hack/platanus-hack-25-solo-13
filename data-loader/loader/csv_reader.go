package loader

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// ReadCSV lee el archivo CSV y retorna un slice de CSVRecord
func ReadCSV(filePath string) ([]CSVRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo archivo CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Leer todas las filas
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error leyendo CSV: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("archivo CSV vacío")
	}

	// Parsear header para encontrar índices de columnas
	header := rows[0]
	indices := make(map[string]int)
	for i, col := range header {
		indices[strings.ToLower(strings.TrimSpace(col))] = i
	}

	// Validar que existan las columnas requeridas
	requiredCols := []string{"materia", "curso", "objetivo"}
	for _, col := range requiredCols {
		if _, exists := indices[col]; !exists {
			return nil, fmt.Errorf("columna requerida '%s' no encontrada en CSV", col)
		}
	}

	// Parsear filas de datos
	var records []CSVRecord
	for i, row := range rows[1:] { // Saltar header
		if len(row) < len(header) {
			return nil, fmt.Errorf("fila %d: número de columnas incorrecto", i+2)
		}

		record := CSVRecord{
			Materia:  strings.TrimSpace(row[indices["materia"]]),
			Curso:    strings.TrimSpace(row[indices["curso"]]),
			Objetivo: strings.TrimSpace(row[indices["objetivo"]]),
		}

		// Validar que no estén vacíos
		if record.Materia == "" || record.Curso == "" || record.Objetivo == "" {
			return nil, fmt.Errorf("fila %d: campos materia/curso/objetivo no pueden estar vacíos", i+2)
		}

		records = append(records, record)
	}

	return records, nil
}
