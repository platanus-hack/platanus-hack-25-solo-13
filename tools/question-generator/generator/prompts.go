package generator

import "fmt"

// BuildPrompt construye el prompt completo para OpenAI basado en el tipo de pregunta
func BuildPrompt(questionType string, objective OABloomObjective, dificultad int) string {
	baseContext := fmt.Sprintf(`
CONTEXTO EDUCATIVO:
- Materia: %s
- Curso: %s
- Objetivo de Aprendizaje (OA): %s
- Nivel de Bloom: %s (Nivel %d)
- Objetivo Específico: %s
- Indicadores de Logro: %v
- Tipo de Actividad Sugerida: %s
- Complejidad Estimada: %d/10
- Dificultad de la Pregunta: %d/5

`,
		objective.MateriaNombre,
		objective.CursoNombre,
		objective.OATitulo,
		objective.BloomLevelNombre,
		objective.BloomLevelNumero,
		objective.ObjetivoEspecifico,
		objective.IndicadoresLogro,
		objective.TipoActividadSugerida,
		objective.ComplejidadEstimada,
		dificultad,
	)

	switch questionType {
	case "multiple_choice":
		return baseContext + promptMultipleChoice
	case "true_false":
		return baseContext + promptTrueFalse
	case "fill_blanks":
		return baseContext + promptFillBlanks
	case "drag_drop_matching":
		return baseContext + promptDragDropMatching
	case "sequencing":
		return baseContext + promptSequencing
	case "compare_contrast":
		return baseContext + promptCompareContrast
	case "open_ended":
		return baseContext + promptOpenEnded
	case "criteria_evaluation":
		return baseContext + promptCriteriaEvaluation
	case "concept_map":
		return baseContext + promptConceptMap
	default:
		return baseContext + "Genera una pregunta apropiada para este nivel de Bloom."
	}
}

const promptMultipleChoice = `
TAREA: Generar una pregunta de selección múltiple (4 opciones: A, B, C, D)

INSTRUCCIONES:
1. La pregunta debe evaluar el objetivo específico del nivel de Bloom indicado
2. Las 4 opciones deben ser plausibles y relacionadas con el contenido
3. Solo una opción debe ser correcta
4. Los distractores deben representar errores conceptuales comunes
5. Incluye una explicación clara de por qué la respuesta es correcta

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "pregunta": "Texto de la pregunta aquí",
    "opciones": {
      "A": "Primera opción",
      "B": "Segunda opción",
      "C": "Tercera opción",
      "D": "Cuarta opción"
    },
    "explicacion": "Explicación de por qué la respuesta correcta es correcta y análisis de los distractores"
  },
  "validation_data": {
    "respuesta_correcta": "B"
  },
  "tags": ["tag1", "tag2", "tag3"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptTrueFalse = `
TAREA: Generar una pregunta de Verdadero/Falso

INSTRUCCIONES:
1. La afirmación debe ser clara y específica
2. Evita ambigüedades o trucos de lenguaje
3. La afirmación debe evaluar comprensión genuina del concepto
4. Incluye explicación de por qué es verdadero o falso

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "afirmacion": "La afirmación a evaluar",
    "explicacion": "Explicación detallada de por qué es verdadero o falso"
  },
  "validation_data": {
    "es_verdadero": true
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptFillBlanks = `
TAREA: Generar una pregunta de completar espacios en blanco

INSTRUCCIONES:
1. El texto debe tener 2-4 espacios en blanco marcados como [BLANK_1], [BLANK_2], etc.
2. Las palabras faltantes deben ser conceptos clave del objetivo de aprendizaje
3. El contexto debe dar suficientes pistas sin ser obvio
4. Proporciona las respuestas correctas para cada espacio

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "texto": "El texto con espacios [BLANK_1] que el estudiante debe [BLANK_2]",
    "pistas": ["Pista opcional para BLANK_1", "Pista para BLANK_2"],
    "explicacion": "Explicación del concepto completo"
  },
  "validation_data": {
    "respuestas_correctas": {
      "BLANK_1": ["respuesta1", "sinónimo1"],
      "BLANK_2": ["respuesta2", "variante2"]
    }
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptDragDropMatching = `
TAREA: Generar una pregunta de emparejar conceptos (drag & drop)

INSTRUCCIONES:
1. Crea 4-6 pares de elementos relacionados (concepto-definición, causa-efecto, término-ejemplo)
2. Los elementos deben estar mezclados para que el estudiante los empareje
3. Las relaciones deben ser claras y directas
4. Incluye distractores opcionales (elementos que no tienen pareja)

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "instruccion": "Empareja cada concepto con su definición/ejemplo correspondiente",
    "columna_izquierda": ["Elemento 1", "Elemento 2", "Elemento 3", "Elemento 4"],
    "columna_derecha": ["Match A", "Match B", "Match C", "Match D"],
    "explicacion": "Explicación de cada emparejamiento correcto"
  },
  "validation_data": {
    "emparejamientos_correctos": {
      "Elemento 1": "Match B",
      "Elemento 2": "Match D",
      "Elemento 3": "Match A",
      "Elemento 4": "Match C"
    }
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptSequencing = `
TAREA: Generar una pregunta de ordenar secuencia de pasos/eventos

INSTRUCCIONES:
1. Crea una lista de 4-6 pasos o eventos que deben ordenarse
2. La secuencia debe tener un orden lógico claro (temporal, causal, procedimental)
3. Los pasos deben estar desordenados inicialmente
4. Incluye explicación del orden correcto

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "instruccion": "Ordena los siguientes pasos/eventos en la secuencia correcta",
    "elementos_desordenados": ["Paso C", "Paso A", "Paso D", "Paso B"],
    "explicacion": "Explicación de por qué este es el orden correcto"
  },
  "validation_data": {
    "orden_correcto": ["Paso A", "Paso B", "Paso C", "Paso D"]
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptCompareContrast = `
TAREA: Generar una pregunta de comparar y contrastar conceptos

INSTRUCCIONES:
1. Selecciona 2-3 conceptos relacionados del objetivo de aprendizaje
2. Define 3-5 características o criterios para comparar
3. El estudiante debe identificar similitudes y diferencias
4. Incluye tabla de comparación correcta

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "instruccion": "Compara y contrasta los siguientes conceptos según los criterios dados",
    "conceptos": ["Concepto A", "Concepto B"],
    "criterios": ["Criterio 1", "Criterio 2", "Criterio 3"],
    "explicacion": "Explicación de las diferencias y similitudes clave"
  },
  "validation_data": {
    "tabla_correcta": {
      "Concepto A": {
        "Criterio 1": "Valor A1",
        "Criterio 2": "Valor A2",
        "Criterio 3": "Valor A3"
      },
      "Concepto B": {
        "Criterio 1": "Valor B1",
        "Criterio 2": "Valor B2",
        "Criterio 3": "Valor B3"
      }
    }
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptOpenEnded = `
TAREA: Generar una pregunta de respuesta abierta

INSTRUCCIONES:
1. La pregunta debe requerir pensamiento crítico, análisis o creación
2. Debe ser apropiada para el nivel de Bloom (Aplicar, Analizar, Evaluar, Crear)
3. Proporciona criterios de evaluación claros
4. Incluye ejemplo de respuesta modelo (no se mostrará al estudiante)

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "pregunta": "Pregunta abierta que requiere pensamiento profundo",
    "criterios_evaluacion": [
      "Criterio 1: Descripción de qué se evalúa (peso: 30%)",
      "Criterio 2: Descripción de qué se evalúa (peso: 40%)",
      "Criterio 3: Descripción de qué se evalúa (peso: 30%)"
    ],
    "contexto_adicional": "Información o escenario si es necesario",
    "explicacion": "Guía de evaluación para el docente"
  },
  "validation_data": {
    "respuesta_modelo": "Ejemplo de respuesta completa y bien estructurada",
    "puntos_clave": ["Punto clave 1 que debe mencionar", "Punto clave 2", "Punto clave 3"],
    "requiere_revision_humana": true
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptCriteriaEvaluation = `
TAREA: Generar una pregunta de evaluación según criterios/rúbrica

INSTRUCCIONES:
1. Presenta un caso, producto o situación para evaluar
2. Proporciona 3-5 criterios de evaluación con niveles de desempeño
3. El estudiante debe evaluar y justificar su calificación
4. Apropiado para nivel Bloom: Evaluar

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "caso_a_evaluar": "Descripción del caso, producto o situación",
    "rubrica": {
      "Criterio 1": {
        "Excelente": "Descripción nivel excelente",
        "Satisfactorio": "Descripción nivel satisfactorio",
        "Necesita mejorar": "Descripción nivel bajo"
      },
      "Criterio 2": {
        "Excelente": "Descripción nivel excelente",
        "Satisfactorio": "Descripción nivel satisfactorio",
        "Necesita mejorar": "Descripción nivel bajo"
      }
    },
    "instruccion": "Evalúa el caso según la rúbrica y justifica tu calificación",
    "explicacion": "Guía de evaluación esperada"
  },
  "validation_data": {
    "evaluacion_modelo": {
      "Criterio 1": "Satisfactorio",
      "Criterio 2": "Excelente"
    },
    "justificacion_modelo": "Explicación de la evaluación correcta",
    "requiere_revision_humana": true
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`

const promptConceptMap = `
TAREA: Generar una pregunta de crear mapa conceptual

INSTRUCCIONES:
1. Proporciona un conjunto de conceptos clave relacionados
2. El estudiante debe organizar y conectar los conceptos
3. Define las relaciones esperadas entre conceptos
4. Apropiado para niveles Bloom: Analizar, Crear

FORMATO DE RESPUESTA (JSON):
{
  "question_data": {
    "instruccion": "Crea un mapa conceptual que muestre las relaciones entre los siguientes conceptos",
    "conceptos": ["Concepto 1", "Concepto 2", "Concepto 3", "Concepto 4", "Concepto 5"],
    "pregunta_guia": "¿Cómo se relacionan estos conceptos en el contexto de [tema]?",
    "explicacion": "Descripción de las relaciones esperadas"
  },
  "validation_data": {
    "relaciones_esperadas": [
      {"origen": "Concepto 1", "relacion": "causa", "destino": "Concepto 2"},
      {"origen": "Concepto 2", "relacion": "incluye", "destino": "Concepto 3"},
      {"origen": "Concepto 3", "relacion": "se opone a", "destino": "Concepto 4"}
    ],
    "conceptos_centrales": ["Concepto 1", "Concepto 2"],
    "requiere_revision_humana": true
  },
  "tags": ["tag1", "tag2"]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.
`
