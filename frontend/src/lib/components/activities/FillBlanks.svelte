<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  /**
   * FillBlanks - Componente de completar espacios en blanco
   *
   * ⚠️ FORMATO CRÍTICO: El texto debe usar ___N___ (tres underscores + número + tres underscores)
   *
   * @example
   * // ✅ CORRECTO
   * <FillBlanks
   *   text="La capital es ___1___ y la moneda es ___2___."
   *   blanks={[
   *     { id: 1, answer: "Santiago", caseSensitive: false },
   *     { id: 2, answer: "Peso", caseSensitive: false }
   *   ]}
   * />
   *
   * // ❌ INCORRECTO - NO usar estos formatos:
   * text="La capital es _____ y la moneda es _____."    // Sin números
   * text="La capital es __1__ y la moneda es __2__."   // Solo 2 underscores
   */

  // Props
  let {
    /**
     * Texto con marcadores de blanks en formato ___N___
     * Regex interna: /___(\d+)___/g
     * @type {string}
     * @required
     */
    text = "La capital de Chile es ___1___ y está ubicada en la región ___2___.",

    /**
     * Definición de respuestas correctas para cada blank
     * @type {Array<{id: number, answer: string, caseSensitive: boolean}>}
     * @required
     */
    blanks = [
      { id: 1, answer: "Santiago", caseSensitive: false },
      { id: 2, answer: "Metropolitana", caseSensitive: false }
    ],

    // Metadata educativa
    bloomLevel = "recordar",
    materia = "historia",
    oaId = null,

    // Configuración
    showWordBank = false, // Si true, muestra un banco de palabras
    wordBank = [], // Palabras adicionales para confundir
    allowMultipleAttempts = true,
    showHints = false,

    /**
     * Callback cuando el usuario envía respuesta
     * @type {function|null}
     * @param {Object} data - Datos de la respuesta
     * @param {number|null} data.oaId - ID del objetivo de aprendizaje
     * @param {string} data.bloomLevel - Nivel de Bloom
     * @param {string} data.materia - Materia/asignatura
     * @param {Record<number, string>} data.userAnswers - Todas las respuestas del usuario
     * @param {Record<number, boolean>} data.results - Resultado por cada blank (true = correcto)
     * @param {number} data.score - Porcentaje de acierto (0-100)
     * @param {boolean} data.isCorrect - Si todas las respuestas son correctas
     * @param {number} data.attemptCount - Número de intentos
     * @param {string} data.timestamp - ISO timestamp
     */
    onAnswer = null,

    /**
     * Callback cuando se completa correctamente (todas correctas)
     * @type {function|null}
     * @param {Object} data - Datos de completación
     * @param {number|null} data.oaId - ID del objetivo de aprendizaje
     * @param {string} data.bloomLevel - Nivel de Bloom
     * @param {number} data.score - Siempre 100
     * @param {number} data.attempts - Número de intentos hasta completar
     */
    onComplete = null
  } = $props();

  // Estados locales
  let hasSubmitted = $state(false);
  let results = $state({});
  let attemptCount = $state(0);
  let containerRef = $state(null);

  // Inicializar respuestas vacías
  let userAnswers = $state(
    blanks.reduce((acc, blank) => {
      acc[blank.id] = "";
      return acc;
    }, {})
  );

  // Colores por nivel de Bloom
  const bloomColors = {
    recordar: { bg: 'bg-red-500/20', border: 'border-red-500', text: 'text-red-400' },
    comprender: { bg: 'bg-orange-500/20', border: 'border-orange-500', text: 'text-orange-400' },
    aplicar: { bg: 'bg-yellow-500/20', border: 'border-yellow-500', text: 'text-yellow-400' },
    analizar: { bg: 'bg-green-500/20', border: 'border-green-500', text: 'text-green-400' },
    evaluar: { bg: 'bg-blue-500/20', border: 'border-blue-500', text: 'text-blue-400' },
    crear: { bg: 'bg-purple-500/20', border: 'border-purple-500', text: 'text-purple-400' }
  };

  const materiaColors = {
    matemáticas: 'cyan',
    lenguaje: 'purple',
    historia: 'amber',
    física: 'blue',
    química: 'green',
    biología: 'emerald',
    default: 'slate'
  };

  const currentMateriaColor = materiaColors[materia] || materiaColors.default;
  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.recordar;

  // Parsear texto con blanks
  const textParts = $derived.by(() => {
    const parts = [];
    let remainingText = text;
    const regex = /___(\d+)___/g;
    let lastIndex = 0;
    let match;

    while ((match = regex.exec(text)) !== null) {
      // Texto antes del blank
      if (match.index > lastIndex) {
        parts.push({
          type: 'text',
          content: remainingText.substring(lastIndex, match.index)
        });
      }

      // El blank
      parts.push({
        type: 'blank',
        id: parseInt(match[1])
      });

      lastIndex = match.index + match[0].length;
    }

    // Texto después del último blank
    if (lastIndex < text.length) {
      parts.push({
        type: 'text',
        content: text.substring(lastIndex)
      });
    }

    return parts;
  });

  // Calcular score
  const score = $derived.by(() => {
    if (!hasSubmitted) return 0;
    const correct = Object.values(results).filter(r => r).length;
    return (correct / blanks.length) * 100;
  });

  const allCorrect = $derived(score === 100);

  // Funciones
  function checkAnswer(blankId, userAnswer) {
    const blank = blanks.find(b => b.id === blankId);
    if (!blank) return false;

    const correctAnswer = blank.answer;
    const userAns = userAnswer.trim();

    if (blank.caseSensitive) {
      return userAns === correctAnswer;
    } else {
      return userAns.toLowerCase() === correctAnswer.toLowerCase();
    }
  }

  function handleSubmit() {
    hasSubmitted = true;
    attemptCount++;

    // Verificar cada respuesta
    const newResults = {};
    blanks.forEach(blank => {
      newResults[blank.id] = checkAnswer(blank.id, userAnswers[blank.id]);
    });
    results = newResults;

    // Animación
    if (containerRef) {
      gsap.to(containerRef, {
        scale: 0.98,
        duration: 0.1,
        yoyo: true,
        repeat: 1
      });
    }

    // Callback para backend
    if (onAnswer) {
      onAnswer({
        oaId: oaId,
        bloomLevel: bloomLevel,
        materia: materia,
        userAnswers: userAnswers,
        results: results,
        score: score,
        isCorrect: allCorrect,
        attemptCount: attemptCount,
        timestamp: new Date().toISOString()
      });
    }

    // Si todo está correcto, llamar onComplete
    if (allCorrect && onComplete) {
      setTimeout(() => {
        onComplete({
          oaId: oaId,
          bloomLevel: bloomLevel,
          score: 100,
          attempts: attemptCount
        });
      }, 1500);
    }
  }

  function handleTryAgain() {
    if (allowMultipleAttempts) {
      hasSubmitted = false;
      // Mantener las respuestas incorrectas vacías, las correctas las dejamos
      const newAnswers = {};
      blanks.forEach(blank => {
        newAnswers[blank.id] = results[blank.id] ? userAnswers[blank.id] : "";
      });
      userAnswers = newAnswers;
      results = {};
    }
  }

  function handleWordBankClick(word) {
    // Encontrar el primer blank vacío
    const emptyBlank = blanks.find(blank => !userAnswers[blank.id] || userAnswers[blank.id].trim() === "");
    if (emptyBlank) {
      userAnswers = { ...userAnswers, [emptyBlank.id]: word };
    }
  }

  // Animación de entrada
  onMount(() => {
    if (containerRef) {
      gsap.from(containerRef, {
        opacity: 0,
        y: 20,
        duration: 0.5,
        ease: 'power2.out'
      });
    }
  });
</script>

<div
  bind:this={containerRef}
  class="w-full max-w-3xl mx-auto p-6 bg-canvas-900 rounded-2xl border border-slate-700 shadow-2xl"
>
  <!-- Header -->
  <div class="flex items-center justify-between mb-6">
    <div class="flex items-center gap-3">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase {currentBloomStyle.bg} {currentBloomStyle.border} {currentBloomStyle.text} border">
        {bloomLevel}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-semibold bg-{currentMateriaColor}-500/20 text-{currentMateriaColor}-400 border border-{currentMateriaColor}-500">
        {materia}
      </span>
    </div>
    {#if attemptCount > 0}
      <span class="text-xs text-slate-500">Intento {attemptCount}</span>
    {/if}
  </div>

  <!-- Instrucciones -->
  <div class="mb-6">
    <h3 class="text-lg font-semibold text-white mb-2">
      Completa los espacios en blanco
    </h3>
    <p class="text-slate-400 text-sm">
      {showWordBank ? 'Escribe o selecciona palabras del banco' : 'Escribe las palabras correctas en cada espacio'}
    </p>
  </div>

  <!-- Banco de palabras (opcional) -->
  {#if showWordBank && wordBank.length > 0}
    <div class="mb-6 p-4 bg-canvas-900/50 rounded-xl border border-slate-700">
      <p class="text-xs font-semibold text-slate-400 mb-3 uppercase">Banco de palabras:</p>
      <div class="flex flex-wrap gap-2">
        {#each wordBank as word}
          <button
            onclick={() => handleWordBankClick(word)}
            disabled={hasSubmitted && !allowMultipleAttempts}
            class="
              px-4 py-2 rounded-lg
              bg-canvas-800 text-white border border-slate-600
              hover:bg-slate-700 hover:border-cyan-500
              transition-all duration-300
              disabled:opacity-50 disabled:cursor-not-allowed
            "
          >
            {word}
          </button>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Texto con blanks -->
  <div class="mb-6 p-6 bg-canvas-900/50 rounded-xl border border-slate-700">
    <div class="text-lg text-white leading-relaxed flex flex-wrap items-center gap-2">
      {#each textParts as part}
        {#if part.type === 'text'}
          <span>{part.content}</span>
        {:else if part.type === 'blank'}
          {@const isCorrect = hasSubmitted && results[part.id]}
          {@const isIncorrect = hasSubmitted && !results[part.id]}

          <span class="inline-flex items-center">
            <input
              type="text"
              bind:value={userAnswers[part.id]}
              disabled={hasSubmitted && !allowMultipleAttempts}
              placeholder="..."
              class="
                inline-block px-3 py-1 min-w-[120px]
                bg-canvas-800 border-2 rounded-lg
                text-white text-center
                focus:outline-none focus:ring-2 focus:ring-cyan-500
                disabled:cursor-not-allowed
                transition-all duration-300
                {isCorrect ? 'border-green-500 bg-green-500/20' : ''}
                {isIncorrect ? 'border-red-500 bg-red-500/20' : ''}
                {!hasSubmitted ? 'border-slate-600' : ''}
              "
            />
            {#if hasSubmitted}
              <span class="ml-2 text-2xl">
                {isCorrect ? '✓' : '✗'}
              </span>
            {/if}
          </span>
        {/if}
      {/each}
    </div>
  </div>

  <!-- Respuestas correctas (si falló) -->
  {#if hasSubmitted && !allCorrect && showHints}
    <div class="mb-6 p-4 bg-yellow-500/10 rounded-xl border border-yellow-500/50">
      <p class="text-xs font-semibold text-yellow-400 mb-2 uppercase">Respuestas correctas:</p>
      <div class="space-y-1">
        {#each blanks as blank}
          {#if !results[blank.id]}
            <p class="text-sm text-slate-300">
              Espacio {blank.id}: <strong class="text-white">{blank.answer}</strong>
            </p>
          {/if}
        {/each}
      </div>
    </div>
  {/if}

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={Object.values(userAnswers).some(ans => !ans || ans.trim() === "")}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-[#E1E1E1] hover:bg-[#CCCCCC]
          text-canvas-900
          transition-all duration-300
          shadow-lg hover:shadow-xl
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Verificar Respuestas
      </button>
    {:else if allowMultipleAttempts && !allCorrect}
      <button
        onclick={handleTryAgain}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-canvas-800 text-white
          border border-slate-700
          transition-all duration-300
          hover:bg-slate-700
        "
      >
        Corregir Errores
      </button>
    {/if}
  </div>

  <!-- Feedback final -->
  {#if hasSubmitted}
    <div class="mt-6 p-4 rounded-2xl {allCorrect ? 'bg-green-500/10 border border-green-500/50' : 'bg-yellow-500/10 border border-yellow-500/50'}">
      <div class="flex items-center gap-3">
        <span class="text-3xl">
          {allCorrect ? '🎉' : '📝'}
        </span>
        <div class="flex-1">
          <p class="{allCorrect ? 'text-green-400' : 'text-yellow-400'} font-semibold mb-1">
            {allCorrect ? '¡Perfecto! Todas las respuestas son correctas' : `${Object.values(results).filter(r => r).length} de ${blanks.length} correctas`}
          </p>
          <div class="flex items-center gap-3 mt-2">
            <div class="flex-1 bg-black/40 rounded-full h-2 overflow-hidden border border-black/50">
              <div
                class="h-full bg-gradient-to-r from-teal-500 to-cyan-500 transition-all duration-500"
                style="width: {score}%"
              ></div>
            </div>
            <span class="text-sm font-semibold text-white">
              {score.toFixed(0)}%
            </span>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
