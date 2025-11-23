<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  /**
   * MultipleChoice - Componente de selección múltiple
   *
   * ⚠️ IMPORTANTE: Los IDs de opciones DEBEN comenzar desde 0
   *
   * @example
   * // ✅ CORRECTO
   * <MultipleChoice
   *   question="¿Cuál es la capital de Chile?"
   *   options={[
   *     { id: 0, text: "Santiago", isCorrect: true },
   *     { id: 1, text: "Valparaíso", isCorrect: false }
   *   ]}
   * />
   *
   * // ❌ INCORRECTO - No empezar desde 1
   * options={[
   *   { id: 1, text: "Santiago", isCorrect: true },  // Bug: primer elemento no seleccionable
   *   { id: 2, text: "Valparaíso", isCorrect: false }
   * ]}
   */

  // Props
  let {
    /**
     * Texto de la pregunta
     * @type {string}
     * @required
     */
    question = "¿Cuál es la capital de Chile?",

    /**
     * Opciones de respuesta. IDs DEBEN empezar desde 0
     * @type {Array<{id: number, text: string, isCorrect: boolean}>}
     * @required
     */
    options = [
      { id: 0, text: "Santiago", isCorrect: true },
      { id: 1, text: "Valparaíso", isCorrect: false },
      { id: 2, text: "Concepción", isCorrect: false },
      { id: 3, text: "La Serena", isCorrect: false }
    ],

    // Metadata educativa
    bloomLevel = "recordar", // recordar | comprender | aplicar | analizar | evaluar | crear
    materia = "historia",
    oaId = null, // ID del objetivo de aprendizaje (backend integration)

    // Configuración del componente
    showFeedback = true,
    allowMultipleAttempts = false,
    showCorrectAnswer = true,

    /**
     * Callback cuando el usuario envía respuesta
     * @type {function|null}
     * @param {Object} data - Datos de la respuesta
     * @param {number|null} data.oaId - ID del objetivo de aprendizaje
     * @param {string} data.bloomLevel - Nivel de Bloom
     * @param {string} data.materia - Materia/asignatura
     * @param {number} data.userAnswer - ID de la opción seleccionada
     * @param {boolean} data.isCorrect - Si la respuesta es correcta
     * @param {number} data.attemptCount - Número de intentos
     * @param {string} data.timestamp - ISO timestamp
     */
    onAnswer = null,

    /**
     * Callback cuando se completa correctamente
     * @type {function|null}
     * @param {Object} data - Datos de completación
     * @param {number|null} data.oaId - ID del objetivo de aprendizaje
     * @param {string} data.bloomLevel - Nivel de Bloom
     * @param {number} data.score - Siempre 1 para single choice
     * @param {number} data.attempts - Número de intentos hasta completar
     */
    onComplete = null
  } = $props();

  // Estados locales (Svelte 5 runes)
  let selectedOption = $state(null);
  let hasSubmitted = $state(false);
  let isCorrect = $state(false);
  let attemptCount = $state(0);
  let containerRef = $state(null);

  // Colores por nivel de Bloom (alineados con Lumera)
  const bloomColors = {
    recordar: { bg: 'bg-red-500/20', border: 'border-red-500', text: 'text-red-400' },
    comprender: { bg: 'bg-orange-500/20', border: 'border-orange-500', text: 'text-orange-400' },
    aplicar: { bg: 'bg-yellow-500/20', border: 'border-yellow-500', text: 'text-yellow-400' },
    analizar: { bg: 'bg-green-500/20', border: 'border-green-500', text: 'text-green-400' },
    evaluar: { bg: 'bg-blue-500/20', border: 'border-blue-500', text: 'text-blue-400' },
    crear: { bg: 'bg-purple-500/20', border: 'border-purple-500', text: 'text-purple-400' }
  };

  // Colores por materia (alineados con Lumera)
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
  function handleOptionSelect(optionId) {
    if (!hasSubmitted || allowMultipleAttempts) {
      selectedOption = optionId;
    }
  }

  function handleSubmit() {
    if (selectedOption === null || selectedOption === undefined) return;

    hasSubmitted = true;
    attemptCount++;

    const selectedOpt = options.find(opt => opt.id === selectedOption);
    isCorrect = selectedOpt?.isCorrect || false;

    // Animación de feedback
    if (containerRef) {
      gsap.to(containerRef, {
        scale: 0.98,
        duration: 0.1,
        yoyo: true,
        repeat: 1
      });
    }

    // Callback para backend (enviar progreso)
    if (onAnswer) {
      onAnswer({
        oaId: oaId,
        bloomLevel: bloomLevel,
        materia: materia,
        userAnswer: selectedOption,
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
      selectedOption = null;
      hasSubmitted = false;
    }
  }

  // Animación de entrada con GSAP
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
  <!-- Header con Bloom Level Badge -->
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

  <!-- Pregunta -->
  <div class="mb-8">
    <h3 class="text-xl font-semibold text-white mb-2">
      {question}
    </h3>
  </div>

  <!-- Opciones -->
  <div class="space-y-3 mb-6">
    {#each options as option (option.id)}
      {@const isSelected = selectedOption === option.id}
      {@const showAsCorrect = hasSubmitted && option.isCorrect && showCorrectAnswer}
      {@const showAsIncorrect = hasSubmitted && isSelected && !option.isCorrect}

      <button
        onclick={() => handleOptionSelect(option.id)}
        disabled={hasSubmitted && !allowMultipleAttempts}
        class="
          w-full p-4 rounded-2xl text-left transition-all duration-300
          border-2
          {isSelected && !hasSubmitted ? 'border-teal-500 bg-teal-500/10' : 'border-slate-700'}
          {showAsCorrect ? 'border-green-500 bg-green-500/20' : ''}
          {showAsIncorrect ? 'border-red-500 bg-red-500/20' : ''}
          {!hasSubmitted ? 'hover:border-slate-600 hover:bg-canvas-800/50 cursor-pointer' : ''}
          {hasSubmitted && !allowMultipleAttempts ? 'cursor-not-allowed opacity-75' : ''}
        "
      >
        <div class="flex items-center justify-between">
          <span class="text-white font-medium">
            {option.text}
          </span>

          <!-- Iconos de feedback -->
          {#if hasSubmitted && showFeedback}
            {#if showAsCorrect}
              <span class="text-green-400 text-2xl">✓</span>
            {:else if showAsIncorrect}
              <span class="text-red-400 text-2xl">✗</span>
            {/if}
          {/if}
        </div>
      </button>
    {/each}
  </div>

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={selectedOption === null || selectedOption === undefined}
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

  <!-- Feedback final -->
  {#if hasSubmitted && showFeedback}
    <div class="mt-6 p-4 rounded-2xl {isCorrect ? 'bg-green-500/10 border border-green-500/50' : 'bg-red-500/10 border border-red-500/50'}">
      <div class="flex items-center gap-3">
        <span class="text-3xl">
          {isCorrect ? '🎉' : '💡'}
        </span>
        <div>
          <p class="{isCorrect ? 'text-green-400' : 'text-red-400'} font-semibold mb-1">
            {isCorrect ? '¡Correcto!' : 'Respuesta incorrecta'}
          </p>
          <p class="text-slate-400 text-sm">
            {isCorrect
              ? 'Has dominado este concepto. ¡Sigue así!'
              : allowMultipleAttempts
                ? 'No te preocupes, puedes intentarlo de nuevo.'
                : 'Revisa el material y vuelve a intentarlo más tarde.'}
          </p>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
