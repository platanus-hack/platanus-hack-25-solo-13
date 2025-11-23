<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  /**
   * TrueFalse - Componente de verdadero o falso
   *
   * @example
   * <TrueFalse
   *   statement="La fotosíntesis produce oxígeno."
   *   correctAnswer={true}
   *   explanation="Correcto. Las plantas liberan O₂ como subproducto."
   *   bloomLevel="comprender"
   *   materia="biología"
   *   onAnswer={handleAnswer}
   * />
   */

  // Props
  let {
    /**
     * Afirmación a evaluar
     * @type {string}
     * @required
     */
    statement = "La capital de Chile es Santiago.",

    /**
     * Respuesta correcta
     * @type {boolean}
     * @required
     */
    correctAnswer = true,

    /**
     * Explicación mostrada después de responder
     * @type {string}
     */
    explanation = "",

    // Metadata educativa
    bloomLevel = "recordar",
    materia = "historia",
    oaId = null,

    // Configuración
    showExplanation = true,
    allowMultipleAttempts = false,
    requireJustification = false, // Si requiere que el estudiante justifique su respuesta

    /**
     * Callback cuando el usuario envía respuesta
     * @type {function|null}
     * @param {Object} data - Datos de la respuesta
     * @param {number|null} data.oaId - ID del objetivo de aprendizaje
     * @param {string} data.bloomLevel - Nivel de Bloom
     * @param {string} data.materia - Materia/asignatura
     * @param {boolean} data.userAnswer - Respuesta seleccionada (true/false)
     * @param {boolean} data.isCorrect - Si la respuesta es correcta
     * @param {number} data.attemptCount - Número de intentos
     * @param {string} data.timestamp - ISO timestamp
     * @param {string} [data.justification] - Justificación (si requireJustification=true)
     */
    onAnswer = null,

    /**
     * Callback cuando se completa correctamente
     * @type {function|null}
     * @param {Object} data - Datos de completación
     * @param {number|null} data.oaId - ID del objetivo de aprendizaje
     * @param {string} data.bloomLevel - Nivel de Bloom
     * @param {number} data.score - Siempre 1
     * @param {number} data.attempts - Número de intentos hasta completar
     */
    onComplete = null
  } = $props();

  // Estados locales
  let selectedAnswer = $state(null); // true, false, o null
  let hasSubmitted = $state(false);
  let isCorrect = $state(false);
  let attemptCount = $state(0);
  let justification = $state("");
  let containerRef = $state(null);

  // Colores por nivel de Bloom
  const bloomColors = {
    recordar: { bg: 'bg-red-500/20', border: 'border-red-500', text: 'text-red-400' },
    comprender: { bg: 'bg-orange-500/20', border: 'border-orange-500', text: 'text-orange-400' },
    aplicar: { bg: 'bg-yellow-500/20', border: 'border-yellow-500', text: 'text-yellow-400' },
    analizar: { bg: 'bg-green-500/20', border: 'border-green-500', text: 'text-green-400' },
    evaluar: { bg: 'bg-blue-500/20', border: 'border-blue-500', text: 'text-blue-400' },
    crear: { bg: 'bg-purple-500/20', border: 'border-purple-500', text: 'text-purple-400' }
  };

  // Colores por materia
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

  // Funciones
  function handleSelect(answer) {
    if (!hasSubmitted || allowMultipleAttempts) {
      selectedAnswer = answer;
    }
  }

  function handleSubmit() {
    if (selectedAnswer === null) return;
    if (requireJustification && !justification.trim()) {
      alert("Por favor, justifica tu respuesta.");
      return;
    }

    hasSubmitted = true;
    attemptCount++;
    isCorrect = selectedAnswer === correctAnswer;

    // Animación de feedback
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
        userAnswer: selectedAnswer,
        justification: requireJustification ? justification : null,
        isCorrect: isCorrect,
        attemptCount: attemptCount,
        timestamp: new Date().toISOString()
      });
    }

    // Si es correcto, llamar onComplete
    if (isCorrect && onComplete) {
      setTimeout(() => {
        onComplete({
          oaId: oaId,
          bloomLevel: bloomLevel,
          score: 1,
          attempts: attemptCount
        });
      }, 1500);
    }
  }

  function handleTryAgain() {
    if (allowMultipleAttempts) {
      selectedAnswer = null;
      hasSubmitted = false;
      justification = "";
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
  class="w-full max-w-2xl mx-auto p-6 bg-canvas-900 rounded-2xl border border-slate-700 shadow-2xl"
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

  <!-- Statement -->
  <div class="mb-8">
    <h3 class="text-xl font-semibold text-white mb-2">
      {statement}
    </h3>
    <p class="text-slate-400 text-sm">Selecciona Verdadero o Falso</p>
  </div>

  <!-- True/False Buttons -->
  <div class="flex gap-4 mb-6">
    <!-- Botón Verdadero -->
    <button
      onclick={() => handleSelect(true)}
      disabled={hasSubmitted && !allowMultipleAttempts}
      class="
        flex-1 p-6 rounded-2xl font-bold text-lg
        transition-all duration-300
        border-2
        {selectedAnswer === true && !hasSubmitted ? 'border-green-500 bg-green-500/20' : 'border-slate-700'}
        {hasSubmitted && correctAnswer === true ? 'border-green-500 bg-green-500/20' : ''}
        {hasSubmitted && selectedAnswer === true && !isCorrect ? 'border-red-500 bg-red-500/20' : ''}
        {!hasSubmitted ? 'hover:border-green-500/50 hover:bg-green-500/10 cursor-pointer' : ''}
        {hasSubmitted && !allowMultipleAttempts ? 'cursor-not-allowed opacity-75' : ''}
      "
    >
      <div class="flex flex-col items-center gap-2">
        <span class="text-4xl">{hasSubmitted && correctAnswer === true ? '✓' : '👍'}</span>
        <span class="text-white">Verdadero</span>
      </div>
    </button>

    <!-- Botón Falso -->
    <button
      onclick={() => handleSelect(false)}
      disabled={hasSubmitted && !allowMultipleAttempts}
      class="
        flex-1 p-6 rounded-2xl font-bold text-lg
        transition-all duration-300
        border-2
        {selectedAnswer === false && !hasSubmitted ? 'border-red-500 bg-red-500/20' : 'border-slate-700'}
        {hasSubmitted && correctAnswer === false ? 'border-green-500 bg-green-500/20' : ''}
        {hasSubmitted && selectedAnswer === false && !isCorrect ? 'border-red-500 bg-red-500/20' : ''}
        {!hasSubmitted ? 'hover:border-red-500/50 hover:bg-red-500/10 cursor-pointer' : ''}
        {hasSubmitted && !allowMultipleAttempts ? 'cursor-not-allowed opacity-75' : ''}
      "
    >
      <div class="flex flex-col items-center gap-2">
        <span class="text-4xl">{hasSubmitted && correctAnswer === false ? '✓' : '👎'}</span>
        <span class="text-white">Falso</span>
      </div>
    </button>
  </div>

  <!-- Justificación opcional -->
  {#if requireJustification}
    <div class="mb-6">
      <label for="justification" class="block text-sm font-medium text-slate-300 mb-2">
        Justifica tu respuesta:
      </label>
      <textarea
        id="justification"
        bind:value={justification}
        disabled={hasSubmitted && !allowMultipleAttempts}
        rows="3"
        placeholder="Explica por qué crees que esta afirmación es verdadera o falsa..."
        class="
          w-full p-4 rounded-xl
          bg-canvas-900 border border-slate-700
          text-white placeholder-slate-500
          focus:outline-none focus:ring-2 focus:ring-cyan-500
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      ></textarea>
    </div>
  {/if}

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={selectedAnswer === null || (requireJustification && !justification.trim())}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-[#E1E1E1] hover:bg-[#CCCCCC]
          text-canvas-900
          transition-all duration-300
          shadow-lg hover:shadow-xl
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Enviar Respuesta
      </button>
    {:else if allowMultipleAttempts && !isCorrect}
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
        Intentar de Nuevo
      </button>
    {/if}
  </div>

  <!-- Feedback final con explicación -->
  {#if hasSubmitted}
    <div class="mt-6 p-4 rounded-2xl {isCorrect ? 'bg-green-500/10 border border-green-500/50' : 'bg-red-500/10 border border-red-500/50'}">
      <div class="flex items-start gap-3">
        <span class="text-3xl">
          {isCorrect ? '🎉' : '💡'}
        </span>
        <div class="flex-1">
          <p class="{isCorrect ? 'text-green-400' : 'text-red-400'} font-semibold mb-1">
            {isCorrect ? '¡Correcto!' : 'Respuesta incorrecta'}
          </p>
          <p class="text-slate-400 text-sm mb-2">
            La respuesta correcta es: <strong class="text-white">{correctAnswer ? 'Verdadero' : 'Falso'}</strong>
          </p>

          {#if showExplanation && explanation}
            <div class="mt-3 p-3 bg-canvas-900/50 rounded-lg">
              <p class="text-xs text-slate-400 mb-1 uppercase font-semibold">Explicación:</p>
              <p class="text-slate-300 text-sm">{explanation}</p>
            </div>
          {/if}

          {#if requireJustification && justification}
            <div class="mt-3 p-3 bg-canvas-900/50 rounded-lg">
              <p class="text-xs text-slate-400 mb-1 uppercase font-semibold">Tu justificación:</p>
              <p class="text-slate-300 text-sm italic">"{justification}"</p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
