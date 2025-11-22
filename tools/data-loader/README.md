# Data Loader - Generación Automática de Objetivos de Bloom

Sistema de carga masiva de Objetivos de Aprendizaje (OAs) con generación automática de niveles de Bloom usando OpenAI GPT.

## Descripción

Este módulo lee un archivo CSV con Objetivos de Aprendizaje del currículo escolar chileno y:

1. **Genera automáticamente** los 6 niveles de la Taxonomía de Bloom (Recordar, Comprender, Aplicar, Analizar, Evaluar, Crear) usando OpenAI
2. **Inserta en BD** los OAs completos con todos sus subobjetivos, indicadores de logro, tipo de actividad sugerida y complejidad estimada
3. **Maneja errores** guardando OAs fallidos en archivos JSON para procesamiento manual posterior

## Estructura de Archivos

```
data-loader/
├── .env.example              # Template de configuración
├── .env                      # Tu configuración (crear desde .env.example)
├── .gitignore
├── go.mod                    # Dependencias Go
├── main.go                   # Entry point CLI
├── loader/
│   ├── csv_reader.go         # Lee y parsea CSV
│   ├── openai_client.go      # Cliente OpenAI con system prompt
│   ├── db_writer.go          # Inserta OAs en PostgreSQL
│   └── types.go              # Structs compartidos
├── input/
│   └── oas_example.csv       # CSV de ejemplo
└── output/
    └── failed_oas_*.json     # OAs que fallaron (generado automáticamente)
```

## Requisitos Previos

1. **Base de datos configurada:**
   - PostgreSQL corriendo (via Docker Compose o local)
   - Migraciones aplicadas (`make migrate-up` desde `/backend`)
   - Tablas `cursos`, `materias`, `bloom_levels` con datos
   - Tabla `bloom_levels` debe tener los 6 niveles pre-cargados (ver migración 000006)

2. **API Key de OpenAI:**
   - Cuenta OpenAI con créditos disponibles
   - API key con permisos de Chat Completion

3. **Go 1.23+** instalado

## Instalación

### 1. Configurar Variables de Entorno

```bash
cd data-loader
cp .env.example .env
```

Editar `.env` con tus credenciales:

```bash
# REQUERIDO: Tu API Key de OpenAI
OPENAI_API_KEY=sk-proj-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

# Base de datos (ajustar según tu setup)
DB_HOST=localhost        # o 'postgres' si corres desde Docker
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=tu_password
DB_NAME=hackathon
DB_SSLMODE=disable

# Configuración OpenAI (opcional)
OPENAI_MODEL=gpt-4o-mini              # Modelo a usar
OPENAI_TIMEOUT_SECONDS=30             # Timeout por request
OPENAI_MAX_RETRIES=3                  # Reintentos ante errores
```

### 2. Instalar Dependencias

```bash
go mod download
```

## Uso

### Formato del CSV de Entrada

El CSV debe tener 3 columnas: `materia`, `curso`, `objetivo`

**Ejemplo (`input/oas.csv`):**

```csv
materia,curso,objetivo
Lenguaje y Comunicación,1° Medio,Comprender la estructura narrativa de textos literarios identificando personajes principales y secundarios
Matemática,1° Medio,Aplicar propiedades de las funciones lineales y cuadráticas para resolver problemas
Historia,2° Medio,Analizar el proceso de independencia de Chile considerando múltiples causas
```

**Importante:**
- Los nombres de `materia` y `curso` deben coincidir exactamente con los existentes en la BD
- El campo `objetivo` es el texto completo del OA que se enviará a OpenAI

### Ejecutar Carga

```bash
# Usar archivo por defecto (input/oas.csv)
go run main.go

# Especificar archivo custom
go run main.go --input=input/mi_archivo.csv
```

### Ejemplo de Salida

```
🔧 Inicializando clientes...
📂 Leyendo archivo CSV: input/oas.csv
✅ Se encontraron 5 OAs para cargar

[1/5] Procesando OA: Comprender la estructura narrativa de textos lite... (Materia: Lenguaje y Comunicación, Curso: 1° Medio)
  🤖 Generando niveles de Bloom con OpenAI...
  ✅ Generados 6 niveles de Bloom
  💾 Guardando en base de datos...
  ✅ OA guardado exitosamente

[2/5] Procesando OA: Aplicar propiedades de las funciones lineales y ... (Materia: Matemática, Curso: 1° Medio)
  🤖 Generando niveles de Bloom con OpenAI...
  ✅ Generados 6 niveles de Bloom
  💾 Guardando en base de datos...
  ✅ OA guardado exitosamente

...

============================================================
📊 REPORTE FINAL
============================================================
Total OAs procesados:  5
✅ Éxitos:              5 (100.0%)
❌ Fallos:              0 (0.0%)
============================================================
```

## Manejo de Errores

Si un OA falla (por error de OpenAI o BD), se guarda en `output/failed_oas_YYYYMMDD_HHMMSS.json`:

```json
[
  {
    "materia": "Matemática",
    "curso": "1° Medio",
    "objetivo": "Resolver ecuaciones cuadráticas...",
    "error": "intento 3/3 falló: context deadline exceeded",
    "timestamp": "2025-01-22T14:30:00Z"
  }
]
```

Puedes:
1. Revisar el error
2. Corregir el CSV si es problema de datos
3. Reejecutar solo los OAs fallidos creando un nuevo CSV

## System Prompt de OpenAI

El prompt usado está en `loader/openai_client.go` y genera:

- **6 niveles de Bloom:** Recordar, Comprender, Aplicar, Analizar, Evaluar, Crear
- **Subobjetivo específico** para cada nivel
- **1-3 indicadores de logro** observables
- **Tipo de actividad sugerida:** (ej: `lectura_guiada`, `ejercicio_practico`, `debate_guiado`)
- **Complejidad estimada:** 1-5 (muy baja a muy alta)

**Ejemplo de respuesta de OpenAI:**

```json
[
  {
    "nivel_bloom": "recordar",
    "objetivo": "Identificar los elementos narrativos de un texto literario",
    "indicadores_logro": [
      "Nombra personajes principales y secundarios",
      "Señala el conflicto central de la historia"
    ],
    "tipo_actividad_sugerida": "lectura_guiada",
    "complejidad_estimada": 2
  },
  {
    "nivel_bloom": "comprender",
    "objetivo": "Explicar cómo el conflicto genera tensión narrativa",
    "indicadores_logro": [
      "Describe el desarrollo del conflicto",
      "Relaciona acciones de personajes con el conflicto"
    ],
    "tipo_actividad_sugerida": "analisis_texto",
    "complejidad_estimada": 3
  }
  // ... niveles 3-6
]
```

## Costos de OpenAI

**Estimación con GPT-4o-mini (modelo recomendado):**

- **Input:** ~800 tokens por OA (system prompt + objetivo)
- **Output:** ~400 tokens por OA (6 niveles JSON)
- **Total:** ~1200 tokens/OA ≈ $0.0002 USD/OA

**Para 100 OAs:** ~$0.02 USD (~$20 pesos chilenos)

Si prefieres usar `gpt-4o` (más potente pero 60x más caro), cambia `OPENAI_MODEL=gpt-4o` en `.env`.

## Validaciones

El script valida:

1. **CSV:**
   - Columnas `materia`, `curso`, `objetivo` presentes
   - Campos no vacíos

2. **Base de Datos:**
   - Materia existe y está activa
   - Curso existe y está activo

3. **OpenAI:**
   - Respuesta es JSON válido
   - Array no vacío
   - Complejidad entre 1-10 (se ajusta automáticamente si está fuera de rango)

## Troubleshooting

### Error: "materia 'X' no encontrada"

**Causa:** La materia no existe en la BD o el nombre no coincide exactamente.

**Solución:**
```bash
# Ver materias disponibles
psql -h localhost -U admin -d hackathon -c "SELECT id, nombre FROM materias WHERE activo = true;"

# Ajustar nombre en CSV para que coincida exactamente
```

### Error: "OPENAI_API_KEY no está configurada"

**Causa:** El archivo `.env` no existe o no tiene la variable.

**Solución:**
```bash
cp .env.example .env
# Editar .env y agregar tu API key
```

### Error: "context deadline exceeded"

**Causa:** OpenAI está tardando más de 30 segundos (timeout).

**Solución:**
```bash
# Aumentar timeout en .env
OPENAI_TIMEOUT_SECONDS=60
```

### Error: "error parseando JSON de OpenAI"

**Causa:** OpenAI retornó texto que no es JSON válido (raro, pero puede pasar).

**Solución:**
- El OA se guardará en `output/failed_oas_*.json`
- Revisar el objetivo original (puede ser muy complejo o ambiguo)
- Simplificar el texto del objetivo y reintentar

## Scripts de Utilidad

### Verificar Datos Previos

Antes de cargar OAs, verifica que tengas cursos y materias:

```bash
# Listar cursos
psql -h localhost -U admin -d hackathon -c "SELECT * FROM cursos;"

# Listar materias
psql -h localhost -U admin -d hackathon -c "SELECT * FROM materias;"

# Verificar niveles de Bloom (deben ser 6)
psql -h localhost -U admin -d hackathon -c "SELECT nivel, nombre FROM bloom_levels ORDER BY nivel;"
```

### Limpiar OAs de Prueba

Si quieres eliminar OAs cargados durante pruebas:

```bash
# ⚠️ CUIDADO: Esto elimina TODOS los OAs
psql -h localhost -U admin -d hackathon -c "
DELETE FROM oa_bloom_objectives;
DELETE FROM objetivos_aprendizaje;
"
```

## Arquitectura Técnica

### Flujo de Datos

```
CSV → ReadCSV() → CSVRecord[]
  ↓
  Para cada CSVRecord:
    ↓
    OpenAIClient.GenerateBloomObjectives(objetivo)
    ↓
    OpenAI API (GPT-4o-mini) → JSON con 6 niveles
    ↓
    DBWriter.SaveOA(record, bloomObjectives)
    ↓
    Transacción PostgreSQL:
      - INSERT objetivos_aprendizaje
      - INSERT 6 × oa_bloom_objectives
    ↓
  Si falla: guardar en failed_oas.json
    ↓
Reporte final
```

### Manejo de Transacciones

Cada OA se inserta en una **transacción atómica**:

- Si falla insertar el OA base → rollback completo
- Si falla insertar un nivel de Bloom → rollback completo
- Solo si todo tiene éxito → commit

Esto garantiza que nunca tendrás un OA sin sus niveles de Bloom.

### Reintentos y Backoff

Para errores de red/OpenAI:

1. **Intento 1:** Inmediato
2. **Intento 2:** Espera 1 segundo (backoff 1²)
3. **Intento 3:** Espera 4 segundos (backoff 2²)

Después de 3 intentos → guardar en `failed_oas.json`

## Desarrollo

### Agregar Nueva Validación

Editar `loader/openai_client.go`:

```go
// Validar que tenga al menos un objetivo
if len(bloomObjectives) == 0 {
    return nil, fmt.Errorf("OpenAI retornó array vacío")
}

// NUEVA VALIDACIÓN
for _, obj := range bloomObjectives {
    if len(obj.IndicadoresLogro) == 0 {
        return nil, fmt.Errorf("objetivo '%s' no tiene indicadores", obj.NivelBloom)
    }
}
```

### Modificar System Prompt

El system prompt está hardcoded en `loader/openai_client.go` (constante `systemPrompt`).

Para cambiar el comportamiento de OpenAI, editar esa constante.

### Agregar Nuevo Campo al CSV

1. Agregar campo a `CSVRecord` en `loader/types.go`
2. Actualizar parsing en `loader/csv_reader.go`
3. Usar nuevo campo en `loader/db_writer.go`

## Licencia

Parte del proyecto Lumera App (Platanus Hack 25).
