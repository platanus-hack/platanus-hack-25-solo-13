# Question Generator for Lumera App

Generador automático de preguntas educativas usando OpenAI GPT-4o-mini para poblar el banco de preguntas del sistema de diagnóstico adaptativo.

## 📋 Características

- ✅ Genera **5 preguntas por cada OA-Bloom objective** (~730 preguntas total)
- ✅ Soporta **9 tipos de preguntas** diferentes
- ✅ **Mapeo automático** de tipo de pregunta según nivel de Bloom
- ✅ **Distribución de dificultad** (1-5) apropiada por nivel
- ✅ Inserción en **lotes** para eficiencia
- ✅ Manejo robusto de errores con **retry automático**
- ✅ Guarda preguntas fallidas en JSON para retry manual
- ✅ Estadísticas detalladas de generación

## 🗂️ Estructura del Proyecto

```
question-generator/
├── main.go                    # Entry point, orchestrator
├── go.mod                     # Go dependencies
├── .env.example               # Template de configuración
├── .gitignore
├── generator/
│   ├── types.go              # Structs (OABloomObjective, Question, Stats)
│   ├── db.go                 # Database queries (fetch objectives, insert questions)
│   ├── prompts.go            # System prompts para cada tipo de pregunta
│   ├── openai_client.go      # OpenAI API client con retry
│   └── question_builder.go   # Bloom level → question type mapping
└── output/
    └── failed_questions_*.json   # Preguntas que fallaron (para retry)
```

## ⚙️ Configuración

### 1. Variables de Entorno

Copia `.env.example` a `.env` y configura:

```bash
# OpenAI API Key (requerido)
OPENAI_API_KEY=sk-your-openai-api-key-here

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=your_secure_password
DB_NAME=hackathon
DB_SSLMODE=disable

# OpenAI Configuration (opcional)
OPENAI_MODEL=gpt-4o-mini
OPENAI_TIMEOUT_SECONDS=45
OPENAI_MAX_RETRIES=3
```

### 2. Dependencias

```bash
go mod download
```

## 🚀 Uso

### Ejecución Básica

```bash
# Desde el directorio raíz del proyecto
cd question-generator

# Generar preguntas para todos los OA-Bloom objectives
go run main.go
```

### Desde Docker (recomendado para consistencia de DB)

```bash
# Desde raíz del proyecto lumera_app
docker run --rm \
  --network lumera_app_default \
  -v $(pwd)/question-generator:/app \
  -w /app \
  -e DB_HOST=postgres \
  --env-file question-generator/.env \
  golang:1.23-alpine \
  sh -c "go run main.go"
```

### Opciones de CLI

```bash
# Skip objectives que ya tienen preguntas (default: true)
go run main.go -skip-existing=false

# Cambiar tamaño de batch para inserción (default: 10 objectives = ~50 questions)
go run main.go -batch-size=20

# Combinar opciones
go run main.go -skip-existing=false -batch-size=5
```

## 📊 Mapeo de Tipos de Pregunta por Nivel de Bloom

El sistema genera automáticamente la distribución apropiada de tipos de pregunta según el nivel cognitivo:

### Nivel 1: Recordar (5 preguntas)
- `multiple_choice` (3x) - Dificultad: 1, 2, 3
- `true_false` (1x) - Dificultad: 1
- `fill_blanks` (1x) - Dificultad: 2

### Nivel 2: Comprender (5 preguntas)
- `multiple_choice` (2x) - Dificultad: 2, 3
- `true_false` (1x) - Dificultad: 2
- `drag_drop_matching` (1x) - Dificultad: 2
- `sequencing` (1x) - Dificultad: 3

### Nivel 3: Aplicar (5 preguntas)
- `multiple_choice` (2x) - Dificultad: 3, 4
- `drag_drop_matching` (1x) - Dificultad: 3
- `sequencing` (1x) - Dificultad: 3
- `open_ended` (1x) - Dificultad: 3

### Nivel 4: Analizar (5 preguntas)
- `compare_contrast` (2x) - Dificultad: 3, 4
- `open_ended` (1x) - Dificultad: 4
- `concept_map` (2x) - Dificultad: 3, 4

### Nivel 5: Evaluar (5 preguntas)
- `criteria_evaluation` (3x) - Dificultad: 4, 5, 5
- `open_ended` (2x) - Dificultad: 4, 5

### Nivel 6: Crear (5 preguntas)
- `open_ended` (3x) - Dificultad: 4, 5, 5
- `concept_map` (2x) - Dificultad: 4, 5

## 🎯 Tipos de Pregunta Soportados

1. **multiple_choice** - Selección múltiple A/B/C/D
2. **true_false** - Verdadero/Falso
3. **fill_blanks** - Completar espacios en blanco
4. **drag_drop_matching** - Emparejar conceptos
5. **sequencing** - Ordenar secuencia
6. **compare_contrast** - Comparar y contrastar
7. **open_ended** - Respuesta abierta
8. **criteria_evaluation** - Evaluar según rúbrica
9. **concept_map** - Crear mapa conceptual

Cada tipo tiene un **prompt especializado** en `generator/prompts.go` que guía a OpenAI para generar preguntas apropiadas.

## 📝 Estructura de Datos

### Entrada: OA-Bloom Objective

El sistema lee de la base de datos:
- Materia (ej: "Lengua y Literatura")
- Curso (ej: "Primero Medio")
- OA completo (título y descripción)
- Nivel de Bloom (1-6)
- Objetivo específico del nivel
- Indicadores de logro
- Tipo de actividad sugerida
- Complejidad estimada

### Salida: Question

```json
{
  "oa_bloom_objective_id": 34,
  "tipo": "multiple_choice",
  "tipo_uso": "all",
  "question_data": {
    "pregunta": "...",
    "opciones": {"A": "...", "B": "...", "C": "...", "D": "..."},
    "explicacion": "..."
  },
  "validation_data": {
    "respuesta_correcta": "B"
  },
  "dificultad_relativa": 3,
  "tags": ["literatura", "análisis", "narrativa"],
  "activa": true
}
```

## 🔄 Proceso de Generación

1. **Fetch Objectives**: Lee todos los OA-Bloom objectives activos de la BD
2. **Check Existing**: Opcionalmente skip si ya tiene preguntas
3. **Map Question Types**: Determina los 5 tipos según nivel de Bloom
4. **Call OpenAI**: Genera cada pregunta con prompt especializado
5. **Validate Response**: Verifica estructura JSON correcta
6. **Build Question Struct**: Convierte respuesta a modelo de BD
7. **Batch Insert**: Guarda en lotes para eficiencia
8. **Handle Failures**: Guarda fallidas en JSON para retry

## 🛠️ Manejo de Errores

### Retry Automático
- 3 intentos por pregunta (configurable con `OPENAI_MAX_RETRIES`)
- Backoff exponencial (2s, 4s, 6s)
- Timeout de 45 segundos por request

### Failed Questions File

Si alguna pregunta falla después de los retries, se guarda en:
```
output/failed_questions_YYYYMMDD_HHMMSS.json
```

Formato:
```json
[
  {
    "oa_bloom_objective_id": 42,
    "tipo": "concept_map",
    "dificultad": 5,
    "error": "timeout after 45s",
    "timestamp": "2025-11-22T15:30:45Z"
  }
]
```

Para retry manual:
1. Revisa el archivo de fallidas
2. Aumenta timeout o retries en `.env`
3. Corre nuevamente con `-skip-existing=false` para ese objetivo específico

## 📈 Estadísticas de Generación

Al finalizar, el sistema muestra:

```
============================================================
📊 FINAL STATISTICS
============================================================
Total Duration:     45m 23s
Total Attempts:     730
✓ Successful:       715 (97.9%)
✗ Failed:           15 (2.1%)

📋 Questions by Type:
  - multiple_choice    : 256
  - true_false         : 48
  - fill_blanks        : 24
  - drag_drop_matching : 72
  - sequencing         : 72
  - compare_contrast   : 48
  - open_ended         : 120
  - criteria_evaluation: 72
  - concept_map        : 48
============================================================

⚠ Review failed_questions_*.json in output/ to retry failed generations
```

## 🎓 Ejemplo Completo de Pregunta Generada

### Input (OA-Bloom Objective #34)
- **Materia**: Lengua y Literatura
- **Curso**: Primero Medio
- **OA**: Analizar textos narrativos considerando narrador, personajes, ambiente
- **Bloom**: Nivel 4 - Analizar
- **Objetivo Específico**: Comparar diferentes tipos de narradores y su efecto en la historia
- **Dificultad**: 4

### Output (compare_contrast)
```json
{
  "question_data": {
    "instruccion": "Compara los siguientes tipos de narradores según los criterios dados",
    "conceptos": ["Narrador protagonista", "Narrador omnisciente"],
    "criterios": ["Conocimiento de la historia", "Perspectiva", "Impacto en el lector"],
    "explicacion": "El narrador protagonista cuenta desde su experiencia personal..."
  },
  "validation_data": {
    "tabla_correcta": {
      "Narrador protagonista": {
        "Conocimiento de la historia": "Limitado a su experiencia",
        "Perspectiva": "Primera persona, subjetiva",
        "Impacto en el lector": "Mayor empatía e identificación"
      },
      "Narrador omnisciente": {
        "Conocimiento de la historia": "Conoce todo, incluso pensamientos",
        "Perspectiva": "Tercera persona, objetiva",
        "Impacto en el lector": "Visión completa y panorámica"
      }
    }
  },
  "tags": ["narrativa", "tipos-narrador", "análisis-literario"]
}
```

## 🔍 Debugging

### Ver queries ejecutadas
```bash
# Los logs muestran cada SQL query
go run main.go 2>&1 | grep "INSERT INTO"
```

### Verificar preguntas generadas
```bash
# Contar preguntas en BD
docker compose exec postgres psql -U admin -d hackathon -c \
  "SELECT COUNT(*) FROM questions;"

# Ver distribución por tipo
docker compose exec postgres psql -U admin -d hackathon -c \
  "SELECT tipo, COUNT(*) FROM questions GROUP BY tipo ORDER BY COUNT DESC;"

# Ver preguntas de un OA-Bloom específico
docker compose exec postgres psql -U admin -d hackathon -c \
  "SELECT tipo, dificultad_relativa FROM questions WHERE oa_bloom_objective_id = 34;"
```

## 🚧 Troubleshooting

### Error: "OPENAI_API_KEY no está configurada"
- Verifica que `.env` existe y tiene `OPENAI_API_KEY=sk-...`
- Si usas Docker, verifica `--env-file question-generator/.env`

### Error: "failed to connect to database"
- En Docker: asegúrate que `DB_HOST=postgres` (no localhost)
- Verifica que el contenedor de postgres está corriendo: `docker compose ps`

### Error: "timeout after 45s"
- Aumenta `OPENAI_TIMEOUT_SECONDS` en `.env` (ej: 60)
- Verifica tu conexión a internet
- OpenAI puede estar lento, retry después

### JSON Parse Error
- OpenAI puede devolver texto explicativo antes del JSON
- El código automáticamente reintenta hasta `MAX_RETRIES`
- Si persiste, revisa el prompt en `generator/prompts.go`

## 📚 Referencias

- Sistema de diagnóstico: `backend/QUESTION_DIAGNOSTIC_SYSTEM.md`
- Migraciones de BD: `backend/migrations/000009-000012*.sql`
- Data loader (OAs): `data-loader/README.md`

## 🎉 Próximos Pasos

Una vez generadas las preguntas:

1. **Frontend Integration**: Conectar componentes de actividad a `/api/questions`
2. **Diagnostic Algorithm**: Implementar lógica adaptativa de selección
3. **AI Validation**: Para preguntas `open_ended` y `concept_map`
4. **Question Review UI**: Panel para docentes revisar/editar preguntas generadas

---

**Desarrollado para Platanus Hack 25 - Lumera App**
