# API de Planes de Aprendizaje - Sistema de Generación Dinámica

## 📋 Resumen

Sistema completo para generar planes de aprendizaje personalizados usando OpenAI (gpt-4o-mini). Los planes se componen de componentes TEACH que se generan de forma lazy (estructura primero, contenido después).

## 🗄️ Estructura de Datos

### Tablas

**`learning_plans`**
- Almacena el plan general para un usuario + OA específico
- Un usuario solo puede tener 1 plan por OA (constraint único)
- Estados: `generando` → `generado` → `error`

**`learning_plan_components`**
- Componentes individuales de cada plan (ordenados secuencialmente)
- Tipo válido: `ExplainAndExploreSlide` - Componente flexible basado en bloques tipados
- Campo `contenido_props` (JSONB) almacena el contenido generado por OpenAI en formato de bloques
- Estados: `pendiente` → `generando` → `generado` → `error`

## 🔌 Endpoints Disponibles

### 1. Generar Plan de Aprendizaje

**POST** `/api/learning-plans/generate`

**Auth:** Requiere JWT Bearer token

**Body:**
```json
{
  "oa_bloom_objective_id": 123
}
```

**Respuesta exitosa (201):**
```json
{
  "id": 1,
  "user_id": 7,
  "oa_bloom_objective_id": 123,
  "titulo": "Plan de Aprendizaje: Oraciones Compuestas",
  "descripcion": "Aprenderás a identificar y construir oraciones compuestas...",
  "tiempo_estimado_minutos": 45,
  "estado": "generado",
  "components": [
    {
      "id": 1,
      "learning_plan_id": 1,
      "orden": 1,
      "tipo_componente": "ExplainAndExploreSlide",
      "objetivo_especifico": "Comprender qué son las oraciones compuestas y sus tipos básicos",
      "tiempo_estimado_minutos": 15,
      "estado": "pendiente",
      "contenido_props": null
    },
    {
      "id": 2,
      "learning_plan_id": 1,
      "orden": 2,
      "tipo_componente": "ExplainAndExploreSlide",
      "objetivo_especifico": "Identificar oraciones compuestas en textos reales",
      "tiempo_estimado_minutos": 20,
      "estado": "pendiente",
      "contenido_props": null
    }
  ],
  "created_at": "2025-11-22T21:00:00Z",
  "updated_at": "2025-11-22T21:00:00Z"
}
```

**Comportamiento:**
- Si ya existe un plan para ese user_id + oa_bloom_objective_id, retorna el existente (no genera uno nuevo)
- OpenAI decide cuántos componentes ExplainAndExploreSlide usar y define el objetivo específico de cada uno
- Aplica SCAFFOLDING PEDAGÓGICO: los primeros componentes enseñan fundamentos, los últimos aumentan complejidad
- Los componentes se crean con estado `pendiente` (sin contenido/bloques aún)

---

### 2. Obtener Plan por ID

**GET** `/api/learning-plans/{id}`

**Auth:** Requiere JWT Bearer token

**Respuesta:** Mismo formato que el endpoint de generar

---

### 3. Obtener Plan por OA

**GET** `/api/learning-plans/by-oa/{oa_bloom_objective_id}`

**Auth:** Requiere JWT Bearer token

**Respuesta:** Mismo formato que el endpoint de generar

**Uso:** Para verificar si ya existe un plan antes de intentar generar uno nuevo

---

### 4. Generar Contenido de Componente

**POST** `/api/learning-plans/{plan_id}/components/{component_id}/generate-content`

**Auth:** Requiere JWT Bearer token

**No requiere body**

**Respuesta exitosa (200):**
```json
{
  "id": 1,
  "learning_plan_id": 1,
  "orden": 1,
  "tipo_componente": "ExplainAndExploreSlide",
  "objetivo_especifico": "Comprender qué son las oraciones compuestas y sus tipos básicos",
  "tiempo_estimado_minutos": 15,
  "estado": "generado",
  "contenido_props": {
    "titulo": "Fundamentos de las Oraciones Compuestas",
    "bloques": [
      {"tipo": "texto", "contenido": "Una oración compuesta es aquella que contiene dos o más verbos conjugados..."},
      {"tipo": "definicion", "termino": "Oración Compuesta", "texto": "Oración que tiene dos o más verbos conjugados unidos por conectores"},
      {"tipo": "ejemplo", "titulo": "Oración Coordinada", "contenido": "María estudia y Pedro trabaja", "analisis": "Dos oraciones independientes unidas por 'y'"},
      {"tipo": "nota", "estilo": "tip", "texto": "Recuerda: cada parte de la oración compuesta podría funcionar sola"},
      {"tipo": "resumen", "puntos": ["Las oraciones compuestas tienen 2+ verbos", "Se unen con conectores", "Cada parte mantiene sentido propio"]}
    ]
  },
  "created_at": "2025-11-22T21:00:00Z",
  "updated_at": "2025-11-22T21:00:15Z"
}
```

**Comportamiento:**
- Si el componente ya tiene contenido generado (bloques), lo retorna directamente
- Si no, consulta OpenAI con el prompt de bloques tipados
- OpenAI genera una secuencia flexible de bloques (texto, ejemplos, definiciones, etc.) adaptada al objetivo específico
- Los bloques pueden alternarse libremente para crear contenido pedagógico profundo

---

## 🎯 Flujo de Uso Recomendado

### Frontend: Generar y Mostrar Plan

```javascript
// 1. Verificar si ya existe un plan
const checkResponse = await fetch(`/api/learning-plans/by-oa/${oaId}`, {
  headers: { 'Authorization': `Bearer ${token}` }
});

let plan;
if (checkResponse.status === 404) {
  // 2. No existe, generar nuevo plan
  const generateResponse = await fetch('/api/learning-plans/generate', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ oa_bloom_objective_id: oaId })
  });
  plan = await generateResponse.json();
} else {
  // 3. Ya existe, usar el existente
  plan = await checkResponse.json();
}

// 4. Mostrar estructura del plan (títulos, objetivos, tiempos)
displayPlanOverview(plan);

// 5. Cuando el usuario quiera ver un componente específico
async function loadComponent(componentId) {
  const response = await fetch(
    `/api/learning-plans/${plan.id}/components/${componentId}/generate-content`,
    {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    }
  );
  const component = await response.json();

  // 6. Renderizar el componente según su tipo
  renderComponent(component.tipo_componente, component.contenido_props);
}
```

---

## 📝 Estructura de Bloques Tipados (ExplainAndExploreSlide)

El componente `ExplainAndExploreSlide` utiliza un sistema de bloques flexibles que permite alternar libremente entre diferentes tipos de contenido pedagógico.

### Estructura General
```typescript
{
  titulo: string;
  bloques: Array<Bloque>; // Secuencia flexible de bloques
}
```

### Tipos de Bloques Disponibles

#### 1. BLOQUE TEXTO
Explicaciones, introducciones, desarrollo de conceptos fundacionales.
```typescript
{
  tipo: "texto",
  contenido: string // Puede ser largo (varios párrafos)
}
```

#### 2. BLOQUE EJEMPLO
Ilustrar conceptos con casos concretos. Incluye análisis opcional.
```typescript
{
  tipo: "ejemplo",
  titulo: string,
  contenido: string,
  analisis?: string // Explicación de qué ilustra este ejemplo
}
```

#### 3. BLOQUE DEFINICIÓN
Términos clave y vocabulario técnico.
```typescript
{
  tipo: "definicion",
  termino: string,
  texto: string
}
```

#### 4. BLOQUE NOTA
Destacar información importante, tips, advertencias.
```typescript
{
  tipo: "nota",
  estilo: "info" | "warning" | "tip",
  texto: string
}
```

#### 5. BLOQUE EJERCICIO
Práctica guiada o actividades de aplicación.
```typescript
{
  tipo: "ejercicio",
  instruccion: string,
  ejemplo?: string // Ejemplo de cómo completar el ejercicio
}
```

#### 6. BLOQUE RESUMEN
Síntesis de conceptos clave al final de secciones.
```typescript
{
  tipo: "resumen",
  puntos: string[] // Array de puntos clave
}
```

#### 7. BLOQUE COMPARACIÓN
Tablas comparativas entre conceptos similares.
```typescript
{
  tipo: "comparacion",
  items: Array<{
    aspecto: string,
    opcion1: string,
    opcion2: string
  }>
}
```

### Ejemplo Completo de Contenido Generado

```json
{
  "titulo": "Fundamentos de las Oraciones Compuestas",
  "bloques": [
    {
      "tipo": "texto",
      "contenido": "Una oración compuesta es aquella que contiene dos o más verbos conjugados. A diferencia de las oraciones simples, las compuestas pueden expresar ideas más complejas al combinar múltiples proposiciones."
    },
    {
      "tipo": "definicion",
      "termino": "Oración Compuesta",
      "texto": "Oración que tiene dos o más verbos conjugados unidos por conectores o yuxtapuestas"
    },
    {
      "tipo": "ejemplo",
      "titulo": "Oración Coordinada Simple",
      "contenido": "María estudia medicina y Pedro trabaja en un hospital",
      "analisis": "Dos oraciones independientes ('María estudia medicina' y 'Pedro trabaja en un hospital') unidas por el conector 'y'"
    },
    {
      "tipo": "nota",
      "estilo": "tip",
      "texto": "Recuerda: cada parte de una oración coordinada podría funcionar como oración independiente"
    },
    {
      "tipo": "ejercicio",
      "instruccion": "Identifica los verbos en la siguiente oración: 'El sol brilla y los pájaros cantan'",
      "ejemplo": "Verbos: 'brilla' y 'cantan'"
    },
    {
      "tipo": "resumen",
      "puntos": [
        "Las oraciones compuestas tienen 2 o más verbos conjugados",
        "Se pueden unir con conectores (coordinadas) o sin ellos (yuxtapuestas)",
        "Cada parte mantiene sentido propio en oraciones coordinadas"
      ]
    }
  ]
}
```

### Ventajas del Sistema de Bloques

1. **Flexibilidad Pedagógica**: OpenAI decide qué bloques usar y en qué orden según el objetivo
2. **Profundidad Variable**: Los bloques "texto" pueden ser tan largos como necesario
3. **Scaffolding Natural**: Permite alternar teoría → ejemplo → teoría → ejemplo → práctica
4. **Adaptabilidad**: Se ajusta a niveles Bloom desde fundamentos (1-2) hasta análisis complejo (5-6)

---

## ⚙️ Configuración OpenAI

Las variables de entorno en `backend/.env`:

```bash
OPENAI_API_KEY=sk-tu-api-key-aqui
OPENAI_MODEL=gpt-4o-mini              # opcional, este es el default
OPENAI_TIMEOUT_SECONDS=60             # opcional
OPENAI_MAX_RETRIES=3                  # opcional
```

---

## 🚨 Manejo de Errores

### Plan con error
Si `learning_plans.estado = 'error'`, el campo `error_mensaje` contiene detalles.

### Componente con error
Si `learning_plan_components.estado = 'error'`, el campo `error_mensaje` contiene detalles.

### Errores comunes
- **401 Unauthorized**: Falta token JWT o es inválido
- **404 Not Found**: Plan o OA no existe
- **500 Internal Server Error**: Error de OpenAI o base de datos (revisar logs)

---

## 📊 Arquitectura Híbrida (Lazy Loading)

**Ventajas:**
1. **Respuesta rápida** (~2-5s): El endpoint `/generate` solo crea la estructura
2. **Menos costo**: Solo genera contenido cuando el usuario lo necesita
3. **Mejor UX**: El usuario ve el plan inmediatamente
4. **Cacheable**: El contenido generado se guarda en DB

**Flujo:**
1. POST `/generate` → OpenAI decide componentes (rápido)
2. GET `/plans/{id}` → Ver estructura del plan (instantáneo)
3. POST `/generate-content` → OpenAI genera contenido detallado (lento, ~5-15s)

---

## 🔍 Ejemplos de Uso con curl

```bash
# 1. Login (obtener token)
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}' \
  | jq -r '.token')

# 2. Generar plan
curl -X POST http://localhost:8080/api/learning-plans/generate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"oa_bloom_objective_id": 1}' | jq

# 3. Obtener plan por OA
curl http://localhost:8080/api/learning-plans/by-oa/1 \
  -H "Authorization: Bearer $TOKEN" | jq

# 4. Generar contenido de componente
curl -X POST http://localhost:8080/api/learning-plans/1/components/1/generate-content \
  -H "Authorization: Bearer $TOKEN" | jq
```

---

## 📚 Archivos Relevantes

**Backend:**
- `backend/internal/models/learning_plan.go` - Modelos GORM
- `backend/internal/handlers/learning_plan.go` - Handlers HTTP
- `backend/internal/services/content_generator.go` - Integración OpenAI
- `backend/internal/services/content_prompts.go` - Prompts especializados
- `backend/migrations/000021_create_learning_plans_tables.up.sql` - Schema

**Frontend (componentes existentes):**
- `frontend/src/lib/components/slides/teach/*.svelte` - Componentes de visualización
