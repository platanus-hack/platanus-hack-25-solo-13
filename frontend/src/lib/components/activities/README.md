# Activity Components - Guía de Uso

Esta carpeta contiene los componentes de actividades educativas para Lumera, alineados con la taxonomía de Bloom y diseñados para el aprendizaje adaptativo.

## 📚 Componentes Disponibles (9 Total)

### Cobertura por Nivel de Bloom:
- **Recordar:** MultipleChoice, TrueFalse, FillBlanks (3 componentes)
- **Comprender:** MultipleChoice, TrueFalse, FillBlanks, DragDropMatching, Sequencing (5 componentes)
- **Aplicar:** MultipleChoice, DragDropMatching, Sequencing (3 componentes)
- **Analizar:** OpenEndedResponse, CompareContrast (2 componentes)
- **Evaluar:** OpenEndedResponse, CriteriaEvaluation (2 componentes)
- **Crear:** OpenEndedResponse, ConceptMapBuilder (2 componentes)

### 1. MultipleChoice.svelte
**Pregunta de selección múltiple**

- **Bloom Levels:** Recordar, Comprender, Aplicar
- **Uso:** Evaluaciones rápidas, PAES simulator, diagnósticos

**Props:**
```javascript
{
  question: string,              // La pregunta a mostrar
  options: Array<{               // Opciones de respuesta
    id: number,
    text: string,
    isCorrect: boolean
  }>,
  bloomLevel: string,            // recordar | comprender | aplicar | analizar | evaluar | crear
  materia: string,               // matemáticas | lenguaje | historia | física | química | biología
  oaId: number | null,           // ID del objetivo de aprendizaje (backend)
  showFeedback: boolean,         // Mostrar feedback visual (default: true)
  allowMultipleAttempts: boolean, // Permitir reintentos (default: false)
  showCorrectAnswer: boolean,    // Mostrar respuesta correcta (default: true)
  onAnswer: function,            // Callback cuando responde (envía datos a backend)
  onComplete: function           // Callback cuando completa correctamente
}
```

**⚠️ IMPORTANTE: IDs de opciones**

**SIEMPRE comienza los IDs desde 0**, ya que el componente usa índices de array:

```javascript
// ✅ CORRECTO
options={[
  { id: 0, text: "Santiago", isCorrect: true },
  { id: 1, text: "Valparaíso", isCorrect: false },
  { id: 2, text: "Concepción", isCorrect: false }
]}

// ❌ INCORRECTO - ¡No empieces desde 1!
options={[
  { id: 1, text: "Santiago", isCorrect: true },  // Bug: primer elemento no será seleccionable
  { id: 2, text: "Valparaíso", isCorrect: false }
]}
```

**Ejemplo de uso:**
```svelte
<script>
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';

  function handleAnswer(data) {
    console.log('Respuesta:', data);
    // Enviar a backend: POST /api/educational/progress
  }

  function handleComplete(data) {
    console.log('Completado:', data);
    // Actualizar progreso del estudiante
  }
</script>

<MultipleChoice
  question="¿Cuál es la capital de Chile?"
  options={[
    { id: 0, text: "Santiago", isCorrect: true },
    { id: 1, text: "Valparaíso", isCorrect: false },
    { id: 2, text: "Concepción", isCorrect: false }
  ]}
  bloomLevel="recordar"
  materia="historia"
  oaId={123}
  onAnswer={handleAnswer}
  onComplete={handleComplete}
/>
```

---

### 2. TrueFalse.svelte
**Verdadero o Falso**

- **Bloom Levels:** Recordar, Comprender
- **Uso:** Daily missions, warm-ups, diagnósticos rápidos

**Props:**
```javascript
{
  statement: string,             // Afirmación a evaluar
  correctAnswer: boolean,        // true o false
  explanation: string,           // Explicación opcional (se muestra después)
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  showExplanation: boolean,      // Mostrar explicación (default: true)
  allowMultipleAttempts: boolean,
  requireJustification: boolean, // Requiere que justifique (default: false)
  onAnswer: function,
  onComplete: function
}
```

**Ejemplo de uso:**
```svelte
<TrueFalse
  statement="La fotosíntesis produce oxígeno como subproducto."
  correctAnswer={true}
  explanation="Correcto. Durante la fotosíntesis, las plantas liberan O₂ al descomponer H₂O."
  bloomLevel="comprender"
  materia="biología"
  requireJustification={true}
  onAnswer={handleAnswer}
  onComplete={handleComplete}
/>
```

---

### 3. OpenEndedResponse.svelte
**Respuesta Abierta**

- **Bloom Levels:** Analizar, Evaluar, Crear
- **Uso:** Ensayos cortos, reflexiones, análisis crítico

**Props:**
```javascript
{
  prompt: string,                // Pregunta o consigna
  placeholder: string,           // Texto placeholder
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  minWords: number,              // Mínimo de palabras (0 = sin mínimo)
  maxWords: number,              // Máximo de palabras (0 = sin límite)
  showWordCount: boolean,        // Mostrar contador (default: true)
  enableAiFeedback: boolean,     // Habilitar feedback IA (default: false)
  rubric: Array<string>,         // Criterios de evaluación opcionales
  onSubmit: function,            // Cuando envía la respuesta
  onComplete: function,
  onDraft: function              // Auto-save de borrador
}
```

**Ejemplo de uso:**
```svelte
<OpenEndedResponse
  prompt="Analiza cómo la Guerra del Pacífico impactó el desarrollo económico de Chile."
  minWords={100}
  maxWords={300}
  bloomLevel="analizar"
  materia="historia"
  rubric={[
    "Menciona al menos 3 consecuencias económicas",
    "Incluye evidencia histórica específica",
    "Analiza causas y efectos de forma clara"
  ]}
  enableAiFeedback={true}
  onSubmit={handleSubmit}
  onDraft={handleAutosave}
/>
```

---

### 4. FillBlanks.svelte
**Completar Espacios en Blanco**

- **Bloom Levels:** Recordar, Comprender
- **Uso:** Vocabulario, fórmulas matemáticas, conceptos clave

**Props:**
```javascript
{
  text: string,                  // Texto con ___1___, ___2___, etc. para marcar blanks
  blanks: Array<{                // Respuestas correctas
    id: number,
    answer: string,
    caseSensitive: boolean
  }>,
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  showWordBank: boolean,         // Mostrar banco de palabras (default: false)
  wordBank: Array<string>,       // Palabras adicionales (distractores)
  allowMultipleAttempts: boolean,
  showHints: boolean,            // Mostrar respuestas correctas si falla
  onAnswer: function,
  onComplete: function
}
```

**⚠️ FORMATO CRÍTICO: Marcadores de Blanks**

**El formato es OBLIGATORIO: `___N___` (tres underscores + número + tres underscores)**

```javascript
// ✅ CORRECTO
text: "La capital es ___1___ y la moneda es ___2___."
blanks: [
  { id: 1, answer: "Santiago", caseSensitive: false },
  { id: 2, answer: "Peso", caseSensitive: false }
]

// ❌ INCORRECTO - Estos formatos NO funcionarán:
text: "La capital es _____ y la moneda es _____."    // Sin números
text: "La capital es __1__ y la moneda es __2__."   // Solo 2 underscores
text: "La capital es ____1____ y la moneda es ____2____."  // 4 underscores
text: "La capital es [1] y la moneda es [2]."       // No usa underscores
```

**Regex interno:** `/___(\d+)___/g` - Si los marcadores no coinciden, los inputs no aparecerán.

**Ejemplo de uso:**
```svelte
<FillBlanks
  text="La capital de Chile es ___1___ y está ubicada en la región ___2___."
  blanks={[
    { id: 1, answer: "Santiago", caseSensitive: false },
    { id: 2, answer: "Metropolitana", caseSensitive: false }
  ]}
  showWordBank={true}
  wordBank={["Santiago", "Metropolitana", "Valparaíso", "Central"]}
  bloomLevel="recordar"
  materia="historia"
  onAnswer={handleAnswer}
/>
```

---

### 5. DragDropMatching.svelte
**Relacionar Términos (Drag & Drop)**

- **Bloom Levels:** Comprender, Aplicar
- **Uso:** Conectar términos-definiciones, causas-efectos, conceptos

**Props:**
```javascript
{
  title: string,                 // Título de la actividad
  pairs: Array<{                 // Pares a relacionar
    id: number,
    term: string,                // Término (lado izquierdo)
    definition: string           // Definición (lado derecho)
  }>,
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  shuffleOptions: boolean,       // Mezclar definiciones (default: true)
  showFeedback: boolean,
  allowMultipleAttempts: boolean,
  onAnswer: function,
  onComplete: function
}
```

**Ejemplo de uso:**
```svelte
<DragDropMatching
  title="Relaciona los conceptos de Biología con sus definiciones"
  pairs={[
    {
      id: 1,
      term: "Fotosíntesis",
      definition: "Proceso de conversión de luz solar en energía química"
    },
    {
      id: 2,
      term: "Respiración Celular",
      definition: "Proceso de obtención de energía de moléculas orgánicas"
    },
    {
      id: 3,
      term: "Mitosis",
      definition: "División celular que produce dos células hijas idénticas"
    }
  ]}
  bloomLevel="comprender"
  materia="biología"
  onAnswer={handleAnswer}
/>
```

---

### 6. Sequencing.svelte
**Ordenar Secuencia**

- **Bloom Levels:** Comprender, Aplicar
- **Uso:** Ordenar eventos cronológicos, pasos de procesos, procedimientos

**Props:**
```javascript
{
  title: string,                 // Título de la actividad
  items: Array<{                 // Items a ordenar
    id: number,
    content: string,             // Texto del item
    correctOrder: number         // Posición correcta (1, 2, 3...)
  }>,
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  shuffleItems: boolean,         // Mezclar orden inicial (default: true)
  showNumbers: boolean,          // Mostrar números de posición (default: true)
  showHints: boolean,            // Mostrar botón "Ver orden correcto"
  allowMultipleAttempts: boolean,
  onAnswer: function,
  onComplete: function
}
```

**Ejemplo de uso:**
```svelte
<Sequencing
  title="Ordena los eventos de la Independencia de Chile cronológicamente"
  items={[
    { id: 1, content: "Primera Junta Nacional de Gobierno", correctOrder: 1 },
    { id: 2, content: "Batalla de Rancagua", correctOrder: 2 },
    { id: 3, content: "Cruce de los Andes", correctOrder: 3 },
    { id: 4, content: "Batalla de Chacabuco", correctOrder: 4 },
    { id: 5, content: "Declaración de Independencia", correctOrder: 5 }
  ]}
  bloomLevel="comprender"
  materia="historia"
  showHints={true}
  onAnswer={handleAnswer}
/>
```

---

### 7. CompareContrast.svelte
**Comparar y Contrastar Conceptos**

- **Bloom Levels:** Analizar
- **Uso:** Análisis de similitudes y diferencias entre dos conceptos, pensamiento crítico

**Props:**
```javascript
{
  title: string,                 // Título de la actividad
  itemA: {                       // Primer concepto a comparar
    name: string,
    color: string                // Color del badge (cyan, purple, green, etc.)
  },
  itemB: {                       // Segundo concepto a comparar
    name: string,
    color: string
  },
  characteristics: Array<{       // Características para clasificar
    id: number,
    text: string,                // Descripción de la característica
    correctColumn: string        // "A" | "B" | "both"
  }>,
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  showFeedback: boolean,
  allowMultipleAttempts: boolean,
  onAnswer: function,
  onComplete: function
}
```

**Ejemplo de uso:**
```svelte
<CompareContrast
  title="Compara las características de células animales y vegetales"
  itemA={{ name: "Célula Animal", color: "cyan" }}
  itemB={{ name: "Célula Vegetal", color: "green" }}
  characteristics={[
    { id: 1, text: "Tiene pared celular", correctColumn: "B" },
    { id: 2, text: "Tiene membrana celular", correctColumn: "both" },
    { id: 3, text: "Tiene cloroplastos", correctColumn: "B" },
    { id: 4, text: "Tiene centriolos", correctColumn: "A" },
    { id: 5, text: "Tiene núcleo", correctColumn: "both" }
  ]}
  bloomLevel="analizar"
  materia="biología"
  onAnswer={handleAnswer}
/>
```

**Interacción:**
- Estudiante arrastra características a 3 columnas: solo A, solo B, o Ambos
- Feedback visual con colores verde (correcto) y rojo (incorrecto)
- Opción de corregir solo los errores en intentos posteriores

---

### 8. CriteriaEvaluation.svelte
**Evaluación por Criterios**

- **Bloom Levels:** Evaluar
- **Uso:** Evaluar argumentos, fuentes históricas, calidad de trabajos usando rúbricas

**Props:**
```javascript
{
  title: string,                 // Título de la actividad
  subject: string,               // Qué se está evaluando
  description: string,           // Descripción breve
  content: string | null,        // Contenido a evaluar (opcional)
  criteria: Array<{              // Criterios de evaluación
    id: number,
    name: string,                // Nombre del criterio
    description: string,         // Descripción detallada
    expectedRating: number,      // Rating esperado (1-5)
    weight: number               // Peso porcentual en score final
  }>,
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  showFeedback: boolean,
  allowMultipleAttempts: boolean,
  showExpectedRatings: boolean,  // Mostrar ratings esperados después (default: false)
  onAnswer: function,
  onComplete: function
}
```

**Ejemplo de uso:**
```svelte
<CriteriaEvaluation
  title="Evalúa la calidad de este documento histórico"
  subject="Artículo sobre la Guerra del Pacífico"
  description="Un artículo de 1885 sobre el impacto económico"
  content="La victoria en la Guerra del Pacífico (1879-1884) transformó a Chile..."
  criteria={[
    {
      id: 1,
      name: "Evidencia histórica",
      description: "¿Menciona datos y fechas específicas?",
      expectedRating: 5,
      weight: 30
    },
    {
      id: 2,
      name: "Objetividad",
      description: "¿Presenta múltiples perspectivas?",
      expectedRating: 3,
      weight: 25
    }
  ]}
  bloomLevel="evaluar"
  materia="historia"
  showExpectedRatings={true}
  onAnswer={handleAnswer}
/>
```

**Interacción:**
- Escala de 1-5 estrellas por cada criterio
- Sistema de tolerancia: exacto (100%), ±1 (60%), ±2 (30%)
- Feedback diferenciado: perfecto ✓, cercano ~, incorrecto ✗

---

### 9. ConceptMapBuilder.svelte
**Constructor de Mapas Conceptuales**

- **Bloom Levels:** Crear
- **Uso:** Crear representaciones visuales de conceptos y relaciones, síntesis de conocimiento

**Props:**
```javascript
{
  title: string,                 // Título de la actividad
  topic: string,                 // Tema del mapa conceptual
  instructions: string,          // Instrucciones para el estudiante
  requiredConcepts: Array<string>, // Conceptos que deben aparecer
  suggestedConnections: Array<{  // Conexiones esperadas
    from: string,                // Concepto origen
    to: string,                  // Concepto destino
    label: string                // Etiqueta de la relación
  }>,
  minConcepts: number,           // Mínimo de conceptos requeridos
  minConnections: number,        // Mínimo de conexiones requeridas
  bloomLevel: string,
  materia: string,
  oaId: number | null,
  showFeedback: boolean,
  allowMultipleAttempts: boolean,
  provideConcepts: boolean,      // Si true, conceptos vienen predefinidos (default: false)
  onAnswer: function,
  onComplete: function
}
```

**Ejemplo de uso:**
```svelte
<ConceptMapBuilder
  title="Crea un mapa conceptual"
  topic="Fotosíntesis"
  instructions="Identifica conceptos clave y sus relaciones"
  requiredConcepts={["Fotosíntesis", "Luz Solar", "Clorofila", "Oxígeno", "Glucosa"]}
  suggestedConnections={[
    { from: "Luz Solar", to: "Fotosíntesis", label: "inicia" },
    { from: "Fotosíntesis", to: "Oxígeno", label: "produce" }
  ]}
  minConcepts={5}
  minConnections={4}
  bloomLevel="crear"
  materia="biología"
  onAnswer={handleAnswer}
/>
```

**Interacción:**
- Agregar nodos (conceptos) escribiendo texto
- Seleccionar 2 nodos para crear conexión
- Etiquetar la relación entre nodos
- Arrastrar nodos para reorganizar visualmente
- Canvas SVG para dibujar conexiones con flechas

**Validación:**
- Verifica que conceptos requeridos estén presentes (búsqueda fuzzy)
- Verifica que conexiones sugeridas existan
- Score: 50% conceptos + 50% conexiones

---

## ⚡ Uso en Wizards/Secuencias de Preguntas

### ⚠️ PATRÓN OBLIGATORIO: Usar `{#key}` para resetear estado

Cuando uses estos componentes en secuencias (como tests diagnósticos), **SIEMPRE** envuélvelos en `{#key}` para forzar la recreación del componente:

```svelte
<script>
  let currentQuestionIndex = $state(0);
  let questions = [/* ... */];

  const currentQuestion = $derived(questions[currentQuestionIndex]);
</script>

<!-- ✅ CORRECTO - Se resetea el estado al cambiar de pregunta -->
{#key currentQuestion.id}
  <MultipleChoice
    question={currentQuestion.text}
    options={currentQuestion.options}
    onAnswer={handleAnswer}
  />
{/key}

<!-- ❌ INCORRECTO - El estado persiste entre preguntas -->
<MultipleChoice
  question={currentQuestion.text}
  options={currentQuestion.options}
  onAnswer={handleAnswer}
/>
```

**¿Por qué es necesario?**
- Svelte reutiliza componentes del mismo tipo para eficiencia
- Sin `{#key}`, el estado interno (`selectedOption`, `userAnswers`, etc.) persiste
- Resultado: La selección de la pregunta anterior se mantiene visible en la nueva pregunta

### Ejemplo Completo: Wizard de Diagnóstico

```svelte
<script>
  import { MultipleChoice, TrueFalse, FillBlanks } from '$lib/components/activities';

  let currentIndex = $state(0);
  let answers = $state({});

  const questions = [
    {
      id: 'q1',
      tipo: 'multiple_choice',
      pregunta: '¿Cuál es 2+2?',
      opciones: [
        { id: 0, text: '3', isCorrect: false },
        { id: 1, text: '4', isCorrect: true }
      ]
    },
    {
      id: 'q2',
      tipo: 'fill_blanks',
      pregunta: 'La capital es ___1___.',
      blanks: [{ id: 1, answer: 'Santiago', caseSensitive: false }]
    }
  ];

  const current = $derived(questions[currentIndex]);

  function handleAnswer(data) {
    answers[current.id] = data.isCorrect;
  }

  function next() {
    if (currentIndex < questions.length - 1) {
      currentIndex++;
    }
  }
</script>

<!-- CRÍTICO: Usar {#key} para cada tipo de componente -->
{#key current.id}
  {#if current.tipo === 'multiple_choice'}
    <MultipleChoice
      question={current.pregunta}
      options={current.opciones}
      onAnswer={handleAnswer}
    />
  {:else if current.tipo === 'fill_blanks'}
    <FillBlanks
      text={current.pregunta}
      blanks={current.blanks}
      onAnswer={handleAnswer}
    />
  {/if}
{/key}

<button onclick={next}>Siguiente</button>
```

---

## 🎨 Características Comunes

### Estilos y Diseño
Todos los componentes usan:
- **Dark theme gaming:** Fondo slate-950, bordes slate-800
- **Colores por Bloom level:**
  - Recordar: Rojo
  - Comprender: Naranja
  - Aplicar: Amarillo
  - Analizar: Verde
  - Evaluar: Azul
  - Crear: Púrpura

- **Colores por materia:**
  - Matemáticas: Cyan
  - Lenguaje: Purple
  - Historia: Amber
  - Física: Blue
  - Química: Green
  - Biología: Emerald

### Animaciones GSAP
- Entrada con fade-in y slide-up (0.5s)
- Feedback con escala al enviar respuesta
- Transiciones suaves en hover

### Accesibilidad
- Keyboard navigation (donde aplica)
- ARIA labels
- Contraste de colores WCAG-AA
- Disabled states claros

---

## 🔌 Integración con Backend

### Estructura de Callbacks

**onAnswer:**
Enviado cada vez que el estudiante responde (correcto o incorrecto).
```javascript
{
  oaId: number,              // ID del objetivo de aprendizaje
  bloomLevel: string,        // Nivel de Bloom
  materia: string,           // Materia
  userAnswer: any,           // Respuesta del usuario (varía por componente)
  isCorrect: boolean,        // Si la respuesta es correcta
  attemptCount: number,      // Número de intentos
  timestamp: string          // ISO timestamp
}
```

**Endpoint sugerido:** `POST /api/educational/progress`

**onComplete:**
Enviado solo cuando el estudiante completa correctamente la actividad.
```javascript
{
  oaId: number,
  bloomLevel: string,
  score: number,             // 0-100
  attempts: number           // Número de intentos hasta completar
}
```

**Endpoint sugerido:** `POST /api/educational/complete`

---

## 📊 Datos del Backend

### OA Bloom Objectives
Los componentes esperan recibir datos de la tabla `oa_bloom_objectives`:
```javascript
{
  id: number,
  oa_id: number,
  bloom_level: string,
  objective: string,
  success_indicators: Array<string>,
  suggested_activity_type: string,  // "multiple_choice", "true_false", etc.
  estimated_complexity: number      // 1-10
}
```

### Student OA Progress
Cada callback `onAnswer` debería actualizar `student_oa_progress`:
```javascript
{
  student_id: number,
  oa_bloom_id: number,
  current_state: string,     // "no_iniciado" | "en_proceso" | "logrado" | "dominado"
  percentage_completion: number,
  number_of_attempts: number,
  last_activity_date: timestamp
}
```

### Student OA History
Cada callback debería crear un registro en `student_oa_history`:
```javascript
{
  student_id: number,
  oa_bloom_id: number,
  event_type: string,        // "evaluacion" | "practica" | "diagnostico" | "repaso"
  score: number,
  details: jsonb,            // Detalles específicos del componente
  created_at: timestamp
}
```

---

## 🧪 Mock Data para Testing

Ver la página demo `/components-demo` para ejemplos completos con datos de prueba.

---

## 🚀 Próximos Pasos

1. **Integrar con backend de Lumera**
   - Conectar callbacks a endpoints reales
   - Consumir datos de `oa_bloom_objectives`
   - Actualizar `student_oa_progress`

2. **Testing con usuarios**
   - Probar usabilidad en mobile
   - Medir tiempos de completación
   - Recoger feedback de estudiantes

3. **Optimizaciones**
   - Lazy loading de componentes
   - Caching de respuestas
   - Offline support

4. **Siguientes componentes (Fase 2)**
   - Bloom Level Progress Wheel
   - Learning Path Map
   - Achievement Badge Display

---

**Versión:** 2.0
**Última actualización:** 2025-11-22
**Autores:** Claude (Anthropic) + Lumera Team

**Changelog:**
- v2.0: Agregados 3 componentes para niveles superiores de Bloom (CompareContrast, CriteriaEvaluation, ConceptMapBuilder)
- v1.0: Lanzamiento inicial con 6 componentes base
