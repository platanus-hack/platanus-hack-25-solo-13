# Activity Components - Guía de Uso

Esta carpeta contiene los componentes de actividades educativas para Lumera, alineados con la taxonomía de Bloom y diseñados para el aprendizaje adaptativo.

## 📚 Componentes Disponibles

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
    { id: 1, text: "Santiago", isCorrect: true },
    { id: 2, text: "Valparaíso", isCorrect: false },
    { id: 3, text: "Concepción", isCorrect: false }
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

**Versión:** 1.0
**Última actualización:** 2025-11-22
**Autores:** Claude (Anthropic) + Lumera Team
