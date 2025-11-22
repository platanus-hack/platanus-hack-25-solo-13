package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

const systemPrompt = `Eres un especialista en diseño instruccional, currículo chileno de enseñanza media y aprendizaje adaptativo.

Tu tarea es tomar un Objetivo de Aprendizaje (OA) del currículo escolar y descomponerlo en subobjetivos alineados a la Taxonomía de Bloom: Recordar, Comprender, Aplicar, Analizar, Evaluar y Crear.

Para cada nivel de la taxonomía debes:

1. Generar un subobjetivo específico que represente una acción cognitiva realista que un estudiante de enseñanza media podría realizar para contribuir al cumplimiento del OA original.
2. Proponer entre 1 y 3 indicadores de logro para ese subobjetivo.
   - Cada indicador debe ser observable y evaluable.
   - Debe estar redactado en términos de lo que el/la estudiante es capaz de hacer.
3. Sugerir un tipo de actividad principal para trabajar ese subobjetivo, en el campo ` + "`tipo_actividad_sugerida`" + `.
   - Usa nombres cortos y funcionales, por ejemplo:
     - "video_explicativo"
     - "lectura_guiada"
     - "ejercicio_practico"
     - "juego_interactivo"
     - "debate_guiado"
     - "proyecto_corto"
     - "escritura_creativa"
4. Asignar una complejidad estimada en el campo ` + "`complejidad_estimada`" + `, como un número entero de 1 a 5:
   - 1 = muy baja
   - 2 = baja
   - 3 = media
   - 4 = alta
   - 5 = muy alta
   La complejidad debe ser coherente con el nivel de Bloom, el tipo de actividad y un contexto de 1° medio.

Entrega SIEMPRE tu respuesta en formato JSON con el siguiente esquema:

[
  {
    "nivel_bloom": "recordar",
    "objetivo": "Subobjetivo claro y observable relacionado con el nivel Recordar",
    "indicadores_logro": [
      "Indicador observable de logro para este subobjetivo",
      "Otro indicador, si corresponde"
    ],
    "tipo_actividad_sugerida": "nombre_corto_de_actividad",
    "complejidad_estimada": 1
  },
  {
    "nivel_bloom": "comprender",
    "objetivo": "Subobjetivo claro y observable relacionado con Comprender",
    "indicadores_logro": [
      "Indicador de logro asociado a Comprender"
    ],
    "tipo_actividad_sugerida": "nombre_corto_de_actividad",
    "complejidad_estimada": 2
  }
  // ... repetir para otros niveles que apliquen
]

Instrucciones adicionales:

- Si el OA no permite generar subobjetivos para todos los niveles de Bloom, omite los niveles que no apliquen.
- Usa verbos operativos apropiados a cada nivel, por ejemplo:
  - Recordar: identificar, listar, reconocer;
  - Comprender: explicar, resumir, describir;
  - Aplicar: usar, resolver, implementar;
  - Analizar: comparar, clasificar, descomponer;
  - Evaluar: justificar, argumentar, valorar;
  - Crear: diseñar, producir, elaborar.
- Adapta los subobjetivos, indicadores, tipo_actividad_sugerida y complejidad_estimada al contexto escolar de enseñanza media en Chile (por ejemplo 1° medio), usando lenguaje claro y didáctico.
- No incluyas explicaciones fuera del JSON. El output debe ser SOLO el JSON válido.

Input: un solo Objetivo de Aprendizaje (texto simple).

Output: lista JSON de subobjetivos con sus indicadores de logro, tipo_actividad_sugerida y complejidad_estimada, categorizados por nivel de la Taxonomía de Bloom.`

// OpenAIClient wrapper para el cliente de OpenAI
type OpenAIClient struct {
	client     *openai.Client
	model      string
	timeout    time.Duration
	maxRetries int
}

// NewOpenAIClient crea una nueva instancia del cliente OpenAI
func NewOpenAIClient(apiKey, model string, timeout time.Duration, maxRetries int) *OpenAIClient {
	return &OpenAIClient{
		client:     openai.NewClient(apiKey),
		model:      model,
		timeout:    timeout,
		maxRetries: maxRetries,
	}
}

// GenerateBloomObjectives genera los subobjetivos de Bloom para un OA dado
func (c *OpenAIClient) GenerateBloomObjectives(objetivo string) (OpenAIResponse, error) {
	var lastErr error

	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()

		resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: objetivo,
				},
			},
			Temperature: 0.7,
		})

		if err != nil {
			lastErr = fmt.Errorf("intento %d/%d falló: %w", attempt, c.maxRetries, err)
			if attempt < c.maxRetries {
				// Backoff exponencial
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
				continue
			}
			return nil, lastErr
		}

		if len(resp.Choices) == 0 {
			lastErr = fmt.Errorf("intento %d/%d: respuesta de OpenAI vacía", attempt, c.maxRetries)
			if attempt < c.maxRetries {
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
				continue
			}
			return nil, lastErr
		}

		content := resp.Choices[0].Message.Content

		// Parsear JSON
		var bloomObjectives OpenAIResponse
		if err := json.Unmarshal([]byte(content), &bloomObjectives); err != nil {
			lastErr = fmt.Errorf("intento %d/%d: error parseando JSON de OpenAI: %w\nContenido: %s", attempt, c.maxRetries, err, content)
			if attempt < c.maxRetries {
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Validar que tenga al menos un objetivo
		if len(bloomObjectives) == 0 {
			lastErr = fmt.Errorf("intento %d/%d: OpenAI retornó array vacío de objetivos", attempt, c.maxRetries)
			if attempt < c.maxRetries {
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
				continue
			}
			return nil, lastErr
		}

		// Validar complejidad (debe estar entre 1 y 5, pero la BD acepta 1-10)
		for i := range bloomObjectives {
			if bloomObjectives[i].ComplejidadEstimada < 1 {
				bloomObjectives[i].ComplejidadEstimada = 1
			}
			if bloomObjectives[i].ComplejidadEstimada > 10 {
				bloomObjectives[i].ComplejidadEstimada = 10
			}
		}

		return bloomObjectives, nil
	}

	return nil, lastErr
}
