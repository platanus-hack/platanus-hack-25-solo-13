# Activity Components - Ejemplos de Uso

Ejemplos prácticos de implementación de los componentes educativos, incluyendo casos comunes y mejores prácticas.

---

## Caso 1: Wizard de Evaluación Diagnóstica

### Problema Resuelto
Crear una secuencia de preguntas mixtas (multiple choice, true/false, fill blanks) donde el usuario avanza pregunta por pregunta.

### Solución Completa

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';
  import TrueFalse from '$lib/components/activities/TrueFalse.svelte';
  import FillBlanks from '$lib/components/activities/FillBlanks.svelte';

  // State
  let currentQuestionIndex = $state(0);
  let answers = $state<Record<string, boolean>>({});
  let hasAnswered = $state(false);
  let showResults = $state(false);

  // Banco de preguntas
  const preguntas = [
    {
      id: 'q1',
      tipo: 'multiple_choice',
      pregunta: '¿Cuál es la capital de Chile?',
      opciones: [
        { id: 0, text: 'Santiago', isCorrect: true },
        { id: 1, text: 'Valparaíso', isCorrect: false },
        { id: 2, text: 'Concepción', isCorrect: false }
      ]
    },
    {
      id: 'q2',
      tipo: 'true_false',
      pregunta: 'Los adjetivos califican a los sustantivos.',
      respuestaCorrecta: true,
      explicacion: 'Correcto. Los adjetivos describen o califican sustantivos.'
    },
    {
      id: 'q3',
      tipo: 'fill_blanks',
      pregunta: 'La capital de Francia es ___1___ y su moneda es el ___2___.',
      blanks: [
        { id: 1, answer: 'París', caseSensitive: false },
        { id: 2, answer: 'Euro', caseSensitive: false }
      ]
    }
  ];

  const totalPreguntas = preguntas.length;
  const preguntaActual = $derived(preguntas[currentQuestionIndex]);
  const progreso = $derived(Math.round(((currentQuestionIndex + 1) / totalPreguntas) * 100));

  // Handlers
  function handleAnswer(data: { isCorrect: boolean }) {
    answers[preguntaActual.id] = data.isCorrect;
    hasAnswered = true;
  }

  function siguientePregunta() {
    if (!hasAnswered) return;

    if (currentQuestionIndex < totalPreguntas - 1) {
      currentQuestionIndex++;
      hasAnswered = false;
    } else {
      showResults = true;
    }
  }

  function preguntaAnterior() {
    if (currentQuestionIndex > 0) {
      currentQuestionIndex--;
      hasAnswered = true; // Ya fue respondida
    }
  }

  // Cálculo de resultados
  const respuestasCorrectas = $derived(Object.values(answers).filter(Boolean).length);
  const porcentaje = $derived(Math.round((respuestasCorrectas / totalPreguntas) * 100));
</script>

<div class="max-w-4xl mx-auto p-6">
  {#if !showResults}
    <!-- Barra de progreso -->
    <div class="mb-8">
      <div class="flex justify-between mb-2">
        <span class="text-slate-400">Pregunta {currentQuestionIndex + 1} de {totalPreguntas}</span>
        <span class="text-cyan-400 font-bold">{progreso}%</span>
      </div>
      <div class="h-3 w-full bg-slate-900 rounded-full overflow-hidden">
        <div
          class="h-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all duration-500"
          style="width: {progreso}%"
        ></div>
      </div>
    </div>

    <!-- ⚠️ CRÍTICO: Usar {#key} para forzar reset de estado -->
    <div class="mb-6">
      {#key preguntaActual.id}
        {#if preguntaActual.tipo === 'multiple_choice'}
          <MultipleChoice
            question={preguntaActual.pregunta}
            options={preguntaActual.opciones}
            bloomLevel="comprender"
            materia="historia"
            allowMultipleAttempts={false}
            onAnswer={handleAnswer}
          />
        {:else if preguntaActual.tipo === 'true_false'}
          <TrueFalse
            statement={preguntaActual.pregunta}
            correctAnswer={preguntaActual.respuestaCorrecta}
            explanation={preguntaActual.explicacion}
            bloomLevel="comprender"
            materia="lenguaje"
            allowMultipleAttempts={false}
            onAnswer={handleAnswer}
          />
        {:else if preguntaActual.tipo === 'fill_blanks'}
          <FillBlanks
            text={preguntaActual.pregunta}
            blanks={preguntaActual.blanks}
            bloomLevel="recordar"
            materia="historia"
            allowMultipleAttempts={false}
            onAnswer={handleAnswer}
          />
        {/if}
      {/key}
    </div>

    <!-- Navegación -->
    <div class="flex justify-between">
      <button
        onclick={preguntaAnterior}
        disabled={currentQuestionIndex === 0}
        class="px-6 py-3 rounded-xl bg-slate-800 hover:bg-slate-700 disabled:opacity-50"
      >
        ← Anterior
      </button>

      <button
        onclick={siguientePregunta}
        disabled={!hasAnswered}
        class="px-6 py-3 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-500 disabled:opacity-50"
      >
        {currentQuestionIndex < totalPreguntas - 1 ? 'Siguiente →' : 'Ver Resultados'}
      </button>
    </div>
  {:else}
    <!-- Pantalla de resultados -->
    <div class="text-center">
      <h2 class="text-4xl font-bold mb-4">Resultados</h2>
      <p class="text-2xl text-cyan-400">{porcentaje}% de acierto</p>
      <p class="text-slate-400">{respuestasCorrectas} de {totalPreguntas} correctas</p>

      <button
        onclick={() => goto('/')}
        class="mt-8 px-8 py-4 rounded-xl bg-gradient-to-r from-cyan-500 to-blue-500"
      >
        Continuar
      </button>
    </div>
  {/if}
</div>
```

### Puntos Clave
- ✅ Usar `{#key preguntaActual.id}` para resetear componentes
- ✅ IDs de opciones empiezan desde 0
- ✅ Formato `___N___` en FillBlanks
- ✅ Estado `hasAnswered` controla navegación

---

## Caso 2: Actividad Individual con Reintentos

### Problema Resuelto
Permitir que un estudiante practique una sola actividad con múltiples intentos hasta acertar.

### Solución

```svelte
<script>
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';

  let attemptCount = $state(0);
  let isCompleted = $state(false);

  function handleAnswer(data) {
    attemptCount = data.attemptCount;

    // Registrar progreso en backend
    console.log('Intento #', attemptCount, '- Correcto:', data.isCorrect);
  }

  function handleComplete(data) {
    isCompleted = true;
    console.log('¡Completado en', data.attempts, 'intentos!');

    // Actualizar progreso del estudiante
    // await updateProgress(data);
  }
</script>

<div class="max-w-2xl mx-auto p-6">
  {#if !isCompleted}
    <MultipleChoice
      question="¿Qué es la fotosíntesis?"
      options={[
        { id: 0, text: 'Proceso de respiración celular', isCorrect: false },
        { id: 1, text: 'Conversión de luz solar en energía química', isCorrect: true },
        { id: 2, text: 'Proceso de división celular', isCorrect: false }
      ]}
      bloomLevel="comprender"
      materia="biología"
      oaId={42}
      allowMultipleAttempts={true}
      showCorrectAnswer={true}
      onAnswer={handleAnswer}
      onComplete={handleComplete}
    />

    {#if attemptCount > 0}
      <div class="mt-4 text-center text-slate-400">
        Intentos: {attemptCount}
      </div>
    {/if}
  {:else}
    <div class="text-center p-8 bg-green-500/10 border border-green-500/50 rounded-2xl">
      <span class="text-6xl mb-4">🎉</span>
      <h3 class="text-2xl font-bold text-green-400">¡Actividad Completada!</h3>
      <p class="text-slate-400 mt-2">Lo lograste en {attemptCount} intentos</p>
    </div>
  {/if}
</div>
```

### Puntos Clave
- ✅ `allowMultipleAttempts={true}` permite reintentar
- ✅ `onComplete` solo se llama cuando es correcta
- ✅ `attemptCount` está disponible en ambos callbacks

---

## Caso 3: FillBlanks con Banco de Palabras

### Problema Resuelto
Crear una actividad de completar espacios con palabras predefinidas (incluye distractores).

### Solución

```svelte
<script>
  import FillBlanks from '$lib/components/activities/FillBlanks.svelte';

  function handleAnswer(data) {
    console.log('Score:', data.score + '%');
    console.log('Respuestas:', data.userAnswers);
    console.log('Resultados:', data.results);
  }

  function handleComplete(data) {
    console.log('¡Completado! Intentos:', data.attempts);
  }
</script>

<FillBlanks
  text="El ___1___ es un planeta que orbita alrededor del ___2___, el cual proporciona ___3___ y calor."
  blanks={[
    { id: 1, answer: 'Tierra', caseSensitive: false },
    { id: 2, answer: 'Sol', caseSensitive: false },
    { id: 3, answer: 'luz', caseSensitive: false }
  ]}
  showWordBank={true}
  wordBank={['Tierra', 'Sol', 'luz', 'Luna', 'agua', 'energía']}
  bloomLevel="recordar"
  materia="física"
  allowMultipleAttempts={true}
  showHints={false}
  onAnswer={handleAnswer}
  onComplete={handleComplete}
/>
```

### Puntos Clave
- ✅ Formato `___1___`, `___2___`, `___3___` (tres underscores + número + tres underscores)
- ✅ `wordBank` incluye respuestas correctas + distractores
- ✅ `caseSensitive: false` para respuestas flexibles
- ✅ `showHints={false}` oculta respuestas correctas al fallar

---

## Caso 4: Integración con Backend

### Problema Resuelto
Enviar progreso del estudiante a la API y actualizar su perfil.

### Solución

```svelte
<script lang="ts">
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';
  import { auth } from '$lib/stores/auth.svelte';

  let oaBloomId = 156; // ID del objetivo Bloom

  async function handleAnswer(data) {
    // Registrar cada intento en el historial
    try {
      const response = await fetch('/api/educational/history', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${auth.token}`
        },
        body: JSON.stringify({
          student_id: auth.user.id,
          oa_bloom_id: oaBloomId,
          event_type: 'practica',
          score: data.isCorrect ? 100 : 0,
          details: {
            bloom_level: data.bloomLevel,
            materia: data.materia,
            attempt_count: data.attemptCount,
            user_answer: data.userAnswer,
            timestamp: data.timestamp
          }
        })
      });

      if (!response.ok) {
        console.error('Error al registrar progreso');
      }
    } catch (error) {
      console.error('Error de red:', error);
    }
  }

  async function handleComplete(data) {
    // Actualizar progreso del estudiante cuando completa correctamente
    try {
      const response = await fetch(`/api/educational/progress/${auth.user.id}/${oaBloomId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${auth.token}`
        },
        body: JSON.stringify({
          current_state: 'logrado',
          percentage_completion: 100,
          number_of_attempts: data.attempts
        })
      });

      if (response.ok) {
        console.log('✅ Progreso actualizado');
      }
    } catch (error) {
      console.error('Error al actualizar progreso:', error);
    }
  }
</script>

<MultipleChoice
  question="¿Cuál es la fórmula del teorema de Pitágoras?"
  options={[
    { id: 0, text: 'a² + b² = c²', isCorrect: true },
    { id: 1, text: 'a + b = c', isCorrect: false },
    { id: 2, text: 'a² - b² = c²', isCorrect: false }
  ]}
  bloomLevel="recordar"
  materia="matemáticas"
  oaId={oaBloomId}
  onAnswer={handleAnswer}
  onComplete={handleComplete}
/>
```

### Estructura de Datos Esperada

**POST `/api/educational/history`**
```json
{
  "student_id": 6,
  "oa_bloom_id": 156,
  "event_type": "practica",
  "score": 100,
  "details": {
    "bloom_level": "recordar",
    "materia": "matemáticas",
    "attempt_count": 2,
    "user_answer": 0,
    "timestamp": "2025-11-22T10:30:00Z"
  }
}
```

**PUT `/api/educational/progress/{student_id}/{oa_bloom_id}`**
```json
{
  "current_state": "logrado",
  "percentage_completion": 100,
  "number_of_attempts": 2
}
```

### Puntos Clave
- ✅ `onAnswer` registra cada intento (correcto o no)
- ✅ `onComplete` actualiza estado a "logrado"
- ✅ Incluir `Authorization` header con JWT token
- ✅ Manejar errores de red adecuadamente

---

## Caso 5: Componentes Anidados en Dashboard

### Problema Resuelto
Mostrar actividades dentro de una modal o tarjeta expandible.

### Solución

```svelte
<script>
  import { fly } from 'svelte/transition';
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';

  let showActivity = $state(false);
  let currentActivity = $state(null);

  const actividades = [
    {
      id: 1,
      titulo: 'Fotosíntesis',
      tipo: 'multiple_choice',
      data: {
        question: '¿Qué produce la fotosíntesis?',
        options: [
          { id: 0, text: 'Oxígeno', isCorrect: true },
          { id: 1, text: 'Dióxido de carbono', isCorrect: false }
        ]
      }
    }
  ];

  function openActivity(actividad) {
    currentActivity = actividad;
    showActivity = true;
  }

  function closeActivity() {
    showActivity = false;
    currentActivity = null;
  }

  function handleComplete(data) {
    console.log('Actividad completada:', data);
    setTimeout(closeActivity, 2000); // Cerrar después de 2 segundos
  }
</script>

<!-- Lista de actividades -->
<div class="grid grid-cols-3 gap-4">
  {#each actividades as actividad}
    <button
      onclick={() => openActivity(actividad)}
      class="p-6 bg-slate-900 rounded-2xl hover:bg-slate-800 border border-slate-700"
    >
      <h3 class="font-bold text-white">{actividad.titulo}</h3>
      <p class="text-sm text-slate-400">Bloom: Comprender</p>
    </button>
  {/each}
</div>

<!-- Modal con actividad -->
{#if showActivity && currentActivity}
  <div
    class="fixed inset-0 bg-black/80 flex items-center justify-center p-6 z-50"
    onclick={closeActivity}
  >
    <div
      class="max-w-3xl w-full"
      onclick={(e) => e.stopPropagation()}
      transition:fly={{ y: 20, duration: 300 }}
    >
      {#if currentActivity.tipo === 'multiple_choice'}
        <MultipleChoice
          {...currentActivity.data}
          bloomLevel="comprender"
          materia="biología"
          onComplete={handleComplete}
        />
      {/if}

      <button
        onclick={closeActivity}
        class="mt-4 w-full px-6 py-3 rounded-xl bg-slate-800 text-white"
      >
        Cerrar
      </button>
    </div>
  </div>
{/if}
```

### Puntos Clave
- ✅ Usar `{...data}` spread para pasar props dinámicamente
- ✅ `stopPropagation()` evita cerrar al clickear dentro
- ✅ Transiciones de Svelte para animaciones suaves
- ✅ Auto-cerrar después de completar

---

## Caso 6: Progreso Adaptativo

### Problema Resuelto
Ajustar dificultad de actividades según el desempeño del estudiante.

### Solución

```svelte
<script lang="ts">
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';

  let nivelDificultad = $state(1); // 1 = fácil, 2 = medio, 3 = difícil
  let racha = $state(0); // Respuestas correctas consecutivas

  const preguntasPorNivel = {
    1: [
      {
        id: 'facil_1',
        question: '¿Cuánto es 2 + 2?',
        options: [
          { id: 0, text: '3', isCorrect: false },
          { id: 1, text: '4', isCorrect: true },
          { id: 2, text: '5', isCorrect: false }
        ]
      }
    ],
    2: [
      {
        id: 'medio_1',
        question: '¿Cuánto es 15 × 3?',
        options: [
          { id: 0, text: '40', isCorrect: false },
          { id: 1, text: '45', isCorrect: true },
          { id: 2, text: '50', isCorrect: false }
        ]
      }
    ],
    3: [
      {
        id: 'dificil_1',
        question: 'Si x² - 5x + 6 = 0, ¿cuál es el valor de x?',
        options: [
          { id: 0, text: 'x = 2 o x = 3', isCorrect: true },
          { id: 1, text: 'x = 1 o x = 6', isCorrect: false },
          { id: 2, text: 'x = 5 o x = 1', isCorrect: false }
        ]
      }
    ]
  };

  let currentQuestionIndex = $state(0);
  const preguntaActual = $derived(preguntasPorNivel[nivelDificultad][currentQuestionIndex]);

  function handleAnswer(data) {
    if (data.isCorrect) {
      racha++;

      // Subir de nivel después de 3 aciertos consecutivos
      if (racha >= 3 && nivelDificultad < 3) {
        nivelDificultad++;
        racha = 0;
        currentQuestionIndex = 0;
        console.log('🎉 ¡Subiste al nivel', nivelDificultad + '!');
      } else {
        // Siguiente pregunta del mismo nivel
        currentQuestionIndex = (currentQuestionIndex + 1) % preguntasPorNivel[nivelDificultad].length;
      }
    } else {
      racha = 0;

      // Bajar de nivel después de 2 errores
      if (nivelDificultad > 1) {
        nivelDificultad--;
        currentQuestionIndex = 0;
        console.log('📉 Bajaste al nivel', nivelDificultad);
      }
    }
  }
</script>

<div class="max-w-2xl mx-auto p-6">
  <!-- Indicadores -->
  <div class="flex justify-between mb-6">
    <div class="text-sm">
      <span class="text-slate-400">Nivel:</span>
      <span class="text-cyan-400 font-bold ml-2">{nivelDificultad}</span>
    </div>
    <div class="text-sm">
      <span class="text-slate-400">Racha:</span>
      <span class="text-green-400 font-bold ml-2">{racha}</span>
    </div>
  </div>

  <!-- Actividad -->
  {#key preguntaActual.id}
    <MultipleChoice
      {...preguntaActual}
      bloomLevel="aplicar"
      materia="matemáticas"
      allowMultipleAttempts={false}
      onAnswer={handleAnswer}
    />
  {/key}
</div>
```

### Puntos Clave
- ✅ Adaptar dificultad según racha de aciertos
- ✅ Resetear índice al cambiar de nivel
- ✅ Usar `{#key}` para forzar reset entre preguntas
- ✅ Feedback visual del nivel actual

---

## Troubleshooting Common Issues

### Problema: "No puedo seleccionar la primera opción"

```javascript
// ❌ INCORRECTO
options={[
  { id: 1, text: "Santiago", isCorrect: true },  // ¡Empieza desde 1!
  { id: 2, text: "Valparaíso", isCorrect: false }
]}

// ✅ CORRECTO
options={[
  { id: 0, text: "Santiago", isCorrect: true },  // Empieza desde 0
  { id: 1, text: "Valparaíso", isCorrect: false }
]}
```

### Problema: "Los inputs de FillBlanks no aparecen"

```javascript
// ❌ INCORRECTO
text: "La capital es _____ y la moneda es _____."  // Sin números

// ✅ CORRECTO
text: "La capital es ___1___ y la moneda es ___2___."  // Con números
```

### Problema: "Las respuestas persisten entre preguntas"

```svelte
<!-- ❌ INCORRECTO -->
<MultipleChoice
  question={current.text}
  options={current.options}
/>

<!-- ✅ CORRECTO -->
{#key current.id}
  <MultipleChoice
    question={current.text}
    options={current.options}
  />
{/key}
```

---

**Última actualización:** 2025-11-22
**Versión:** 1.0
