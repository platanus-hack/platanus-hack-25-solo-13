<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    // Data del componente
    prompt = "Analiza cómo la Guerra del Pacífico impactó el desarrollo económico de Chile en el siglo XIX.",
    placeholder = "Escribe tu respuesta aquí...",

    // Metadata educativa
    bloomLevel = "analizar",
    materia = "historia",
    oaId = null,

    // Configuración
    minWords = 0,
    maxWords = 0, // 0 = sin límite
    showWordCount = true,
    enableAiFeedback = false, // Para integración futura con IA
    rubric = [], // Array de criterios de evaluación (opcional)

    // Callbacks
    onSubmit = null,
    onComplete = null,
    onDraft = null // Callback para auto-save
  } = $props();

  // Estados locales
  let response = $state("");
  let hasSubmitted = $state(false);
  let wordCount = $state(0);
  let characterCount = $state(0);
  let containerRef = $state(null);
  let aiFeedback = $state(null);
  let isGeneratingFeedback = $state(false);

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
  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.analizar;

  // Reactivamente calcular conteo de palabras
  $effect(() => {
    const trimmed = response.trim();
    wordCount = trimmed ? trimmed.split(/\s+/).length : 0;
    characterCount = response.length;

    // Auto-save draft cada 3 segundos
    if (onDraft && response.length > 0 && !hasSubmitted) {
      const timer = setTimeout(() => {
        onDraft({
          oaId: oaId,
          draft: response,
          wordCount: wordCount
        });
      }, 3000);

      return () => clearTimeout(timer);
    }
  });

  // Validación de límites de palabras
  const isUnderMinWords = $derived(minWords > 0 && wordCount < minWords);
  const isOverMaxWords = $derived(maxWords > 0 && wordCount > maxWords);
  const canSubmit = $derived(!isUnderMinWords && !isOverMaxWords && response.trim().length > 0);

  // Funciones
  async function handleSubmit() {
    if (!canSubmit) return;

    hasSubmitted = true;

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
    if (onSubmit) {
      onSubmit({
        oaId: oaId,
        bloomLevel: bloomLevel,
        materia: materia,
        response: response,
        wordCount: wordCount,
        characterCount: characterCount,
        timestamp: new Date().toISOString()
      });
    }

    // Si AI feedback está habilitado, generar feedback
    if (enableAiFeedback) {
      await generateAiFeedback();
    }

    // Llamar onComplete
    if (onComplete) {
      setTimeout(() => {
        onComplete({
          oaId: oaId,
          bloomLevel: bloomLevel,
          wordCount: wordCount
        });
      }, 1000);
    }
  }

  async function generateAiFeedback() {
    isGeneratingFeedback = true;

    // Simulación de llamada a IA (integración futura con backend)
    // TODO: Integrar con /api/educational/ai-feedback
    await new Promise(resolve => setTimeout(resolve, 2000));

    aiFeedback = {
      strengths: [
        "Buen uso de evidencia histórica",
        "Análisis coherente de causas y efectos"
      ],
      improvements: [
        "Podrías expandir el análisis del impacto económico",
        "Considera incluir más fuentes primarias"
      ],
      score: 85
    };

    isGeneratingFeedback = false;
  }

  function handleEdit() {
    hasSubmitted = false;
    aiFeedback = null;
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
  class="w-full max-w-3xl mx-auto p-6 bg-slate-950 rounded-3xl border border-slate-800 shadow-2xl"
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
    {#if hasSubmitted}
      <span class="px-3 py-1 rounded-full text-xs font-semibold bg-green-500/20 text-green-400 border border-green-500">
        ✓ Enviado
      </span>
    {/if}
  </div>

  <!-- Prompt -->
  <div class="mb-6">
    <h3 class="text-xl font-semibold text-white mb-3">
      {prompt}
    </h3>

    <!-- Rúbrica (si existe) -->
    {#if rubric && rubric.length > 0}
      <details class="mt-4 p-4 bg-slate-900/50 rounded-xl border border-slate-700">
        <summary class="text-sm font-semibold text-cyan-400 cursor-pointer">
          Ver criterios de evaluación
        </summary>
        <ul class="mt-3 space-y-2">
          {#each rubric as criterion}
            <li class="flex items-start gap-2">
              <span class="text-cyan-400 mt-1">•</span>
              <span class="text-slate-300 text-sm">{criterion}</span>
            </li>
          {/each}
        </ul>
      </details>
    {/if}
  </div>

  <!-- Textarea para respuesta -->
  <div class="mb-4">
    <textarea
      bind:value={response}
      disabled={hasSubmitted}
      rows="12"
      {placeholder}
      class="
        w-full p-4 rounded-xl
        bg-slate-900 border-2
        {isOverMaxWords ? 'border-red-500' : isUnderMinWords ? 'border-yellow-500' : 'border-slate-700'}
        text-white placeholder-slate-500
        focus:outline-none focus:ring-2 focus:ring-cyan-500
        disabled:opacity-50 disabled:cursor-not-allowed
        transition-all duration-300
      "
    ></textarea>
  </div>

  <!-- Word Count & Status -->
  {#if showWordCount}
    <div class="flex items-center justify-between mb-6 text-sm">
      <div class="flex items-center gap-4">
        <span class="text-slate-400">
          <strong class="{isOverMaxWords ? 'text-red-400' : isUnderMinWords ? 'text-yellow-400' : 'text-white'}">{wordCount}</strong>
          palabras
        </span>
        <span class="text-slate-500">
          {characterCount} caracteres
        </span>
      </div>

      <div class="flex items-center gap-2">
        {#if minWords > 0}
          <span class="text-xs text-slate-500">
            Mín: {minWords} palabras
          </span>
        {/if}
        {#if maxWords > 0}
          <span class="text-xs text-slate-500">
            Máx: {maxWords} palabras
          </span>
        {/if}
      </div>
    </div>
  {/if}

  <!-- Mensajes de validación -->
  {#if !hasSubmitted}
    {#if isUnderMinWords}
      <div class="mb-4 p-3 bg-yellow-500/10 border border-yellow-500/50 rounded-xl">
        <p class="text-yellow-400 text-sm">
          Necesitas al menos {minWords - wordCount} palabras más.
        </p>
      </div>
    {/if}

    {#if isOverMaxWords}
      <div class="mb-4 p-3 bg-red-500/10 border border-red-500/50 rounded-xl">
        <p class="text-red-400 text-sm">
          Has excedido el límite por {wordCount - maxWords} palabras.
        </p>
      </div>
    {/if}
  {/if}

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={!canSubmit}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-cyan-500 to-blue-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-cyan-500/50
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Enviar Respuesta
      </button>
    {:else}
      <button
        onclick={handleEdit}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-slate-800 text-white
          border border-slate-700
          transition-all duration-300
          hover:bg-slate-700
        "
      >
        Editar Respuesta
      </button>
    {/if}
  </div>

  <!-- AI Feedback (si está habilitado) -->
  {#if hasSubmitted && enableAiFeedback}
    <div class="mt-6">
      {#if isGeneratingFeedback}
        <div class="p-6 bg-slate-900/50 rounded-2xl border border-slate-700">
          <div class="flex items-center gap-3">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-cyan-500"></div>
            <span class="text-slate-400">Generando retroalimentación con IA...</span>
          </div>
        </div>
      {:else if aiFeedback}
        <div class="p-6 bg-gradient-to-br from-cyan-500/10 to-blue-500/10 rounded-2xl border border-cyan-500/50">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-2xl">🤖</span>
            <h4 class="text-lg font-semibold text-white">Retroalimentación IA</h4>
          </div>

          <div class="mb-4">
            <p class="text-sm font-semibold text-green-400 mb-2">Fortalezas:</p>
            <ul class="space-y-1">
              {#each aiFeedback.strengths as strength}
                <li class="text-sm text-slate-300 flex items-start gap-2">
                  <span class="text-green-400">✓</span>
                  {strength}
                </li>
              {/each}
            </ul>
          </div>

          <div>
            <p class="text-sm font-semibold text-yellow-400 mb-2">Áreas de mejora:</p>
            <ul class="space-y-1">
              {#each aiFeedback.improvements as improvement}
                <li class="text-sm text-slate-300 flex items-start gap-2">
                  <span class="text-yellow-400">→</span>
                  {improvement}
                </li>
              {/each}
            </ul>
          </div>

          <div class="mt-4 pt-4 border-t border-cyan-500/30">
            <p class="text-sm text-slate-400">
              Puntuación estimada: <strong class="text-cyan-400 text-lg">{aiFeedback.score}/100</strong>
            </p>
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Confirmación de envío -->
  {#if hasSubmitted && !enableAiFeedback}
    <div class="mt-6 p-4 rounded-2xl bg-green-500/10 border border-green-500/50">
      <div class="flex items-center gap-3">
        <span class="text-3xl">🎉</span>
        <div>
          <p class="text-green-400 font-semibold mb-1">
            ¡Respuesta enviada con éxito!
          </p>
          <p class="text-slate-400 text-sm">
            Tu respuesta ha sido guardada y será revisada por el profesor.
          </p>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
