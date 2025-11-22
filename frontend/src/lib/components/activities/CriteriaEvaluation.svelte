<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    // Data del componente
    title = "Evalúa la calidad usando los criterios proporcionados",
    subject = "Argumento o Fuente a Evaluar",
    description = "",
    content = null, // Texto/contenido a evaluar (opcional)
    criteria = [
      {
        id: 1,
        name: "Criterio de ejemplo",
        description: "Descripción del criterio",
        expectedRating: 3, // 1-5
        weight: 100
      }
    ],

    // Metadata educativa
    bloomLevel = "evaluar",
    materia = "historia",
    oaId = null,

    // Configuración
    showFeedback = true,
    allowMultipleAttempts = true,
    showExpectedRatings = false, // Mostrar ratings esperados después de submit

    // Callbacks
    onAnswer = null,
    onComplete = null
  } = $props();

  // Estados locales
  let ratings = $state({}); // { criteriaId: rating (1-5) }
  let hasSubmitted = $state(false);
  let results = $state({});
  let attemptCount = $state(0);
  let containerRef = $state(null);
  let hoveredStar = $state({}); // { criteriaId: starIndex }

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
  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.evaluar;

  // Verificar si todos los criterios están evaluados
  const allRated = $derived(
    criteria.every(criterion => ratings[criterion.id] !== undefined)
  );

  // Calcular score
  const score = $derived.by(() => {
    if (!hasSubmitted) return 0;

    let totalWeightedScore = 0;
    let totalWeight = 0;

    criteria.forEach(criterion => {
      const userRating = ratings[criterion.id] || 0;
      const expectedRating = criterion.expectedRating;
      const weight = criterion.weight || 100 / criteria.length;

      // Tolerancia de ±1 punto da score parcial
      const difference = Math.abs(userRating - expectedRating);
      let criterionScore = 0;

      if (difference === 0) {
        criterionScore = 100; // Perfecto
      } else if (difference === 1) {
        criterionScore = 60; // Cercano
      } else if (difference === 2) {
        criterionScore = 30; // Moderado
      } // difference > 2 = 0 puntos

      totalWeightedScore += (criterionScore * weight);
      totalWeight += weight;
    });

    return totalWeight > 0 ? totalWeightedScore / totalWeight : 0;
  });

  const allCorrect = $derived(score >= 80); // 80% o más se considera éxito

  // Funciones de interacción
  function handleRating(criterionId, rating) {
    if (hasSubmitted && !allowMultipleAttempts) return;
    ratings = { ...ratings, [criterionId]: rating };
  }

  function handleMouseEnter(criterionId, starIndex) {
    hoveredStar = { ...hoveredStar, [criterionId]: starIndex };
  }

  function handleMouseLeave(criterionId) {
    const newHovered = { ...hoveredStar };
    delete newHovered[criterionId];
    hoveredStar = newHovered;
  }

  // Verificar respuestas
  function handleSubmit() {
    if (!allRated) {
      alert("Por favor, evalúa todos los criterios antes de enviar.");
      return;
    }

    hasSubmitted = true;
    attemptCount++;

    // Verificar cada criterio
    const newResults = {};
    criteria.forEach(criterion => {
      const userRating = ratings[criterion.id];
      const expectedRating = criterion.expectedRating;
      const difference = Math.abs(userRating - expectedRating);

      // Perfecto, cercano, o incorrecto
      newResults[criterion.id] = {
        rating: userRating,
        expected: expectedRating,
        difference: difference,
        status: difference === 0 ? 'perfect' : difference === 1 ? 'close' : 'incorrect'
      };
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
        ratings: ratings,
        results: results,
        score: score,
        attemptCount: attemptCount,
        timestamp: new Date().toISOString()
      });
    }

    // Si está bien, llamar onComplete
    if (allCorrect && onComplete) {
      setTimeout(() => {
        onComplete({
          oaId: oaId,
          bloomLevel: bloomLevel,
          score: score,
          attempts: attemptCount
        });
      }, 1500);
    }
  }

  function handleTryAgain() {
    if (allowMultipleAttempts) {
      hasSubmitted = false;
      // Mantener solo los ratings perfectos
      const newRatings = {};
      Object.keys(ratings).forEach(criterionId => {
        if (results[criterionId]?.status === 'perfect') {
          newRatings[criterionId] = ratings[criterionId];
        }
      });
      ratings = newRatings;
      results = {};
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
  class="w-full max-w-5xl mx-auto p-6 bg-canvas-950 rounded-2xl border border-slate-800 shadow-2xl"
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

  <!-- Título e instrucciones -->
  <div class="mb-6">
    <h3 class="text-xl font-semibold text-white mb-2">
      {title}
    </h3>
    {#if description}
      <p class="text-slate-400 text-sm mb-4">
        {description}
      </p>
    {/if}
  </div>

  <!-- Sujeto a evaluar -->
  <div class="mb-6 p-4 bg-canvas-900/50 rounded-2xl border border-slate-700">
    <div class="flex items-start gap-3">
      <div class="p-2 bg-blue-500/20 rounded-lg">
        <span class="text-2xl">📄</span>
      </div>
      <div class="flex-1">
        <h4 class="text-lg font-semibold text-blue-400 mb-2">
          {subject}
        </h4>
        {#if content}
          <div class="text-slate-300 text-sm leading-relaxed max-h-48 overflow-y-auto p-3 bg-canvas-950 rounded-lg border border-slate-700">
            {content}
          </div>
        {/if}
      </div>
    </div>
  </div>

  <!-- Rúbrica de evaluación -->
  <div class="mb-6">
    <h4 class="text-sm font-semibold text-slate-400 mb-4 uppercase">Criterios de Evaluación:</h4>
    <div class="space-y-4">
      {#each criteria as criterion (criterion.id)}
        {@const userRating = ratings[criterion.id]}
        {@const result = results[criterion.id]}
        {@const hovered = hoveredStar[criterion.id]}

        <div class="p-4 bg-canvas-900/50 rounded-2xl border border-slate-700 transition-all duration-300
          {result?.status === 'perfect' ? 'border-green-500 bg-green-500/5' : ''}
          {result?.status === 'close' ? 'border-yellow-500 bg-yellow-500/5' : ''}
          {result?.status === 'incorrect' ? 'border-red-500 bg-red-500/5' : ''}
        ">
          <div class="flex items-start justify-between gap-4 mb-3">
            <div class="flex-1">
              <h5 class="text-white font-semibold mb-1">
                {criterion.name}
              </h5>
              <p class="text-slate-400 text-sm">
                {criterion.description}
              </p>
            </div>

            {#if hasSubmitted && result}
              <div class="flex items-center gap-2">
                {#if result.status === 'perfect'}
                  <span class="text-green-400 text-2xl">✓</span>
                {:else if result.status === 'close'}
                  <span class="text-yellow-400 text-2xl">~</span>
                {:else}
                  <span class="text-red-400 text-2xl">✗</span>
                {/if}
              </div>
            {/if}
          </div>

          <!-- Escala de estrellas -->
          <div class="flex items-center gap-2">
            <span class="text-xs text-slate-500 w-16">Calificación:</span>
            <div class="flex gap-1">
              {#each [1, 2, 3, 4, 5] as star}
                {@const isActive = userRating >= star || (hovered !== undefined && hovered >= star)}
                <button
                  onclick={() => handleRating(criterion.id, star)}
                  onmouseenter={() => handleMouseEnter(criterion.id, star)}
                  onmouseleave={() => handleMouseLeave(criterion.id)}
                  disabled={hasSubmitted && !allowMultipleAttempts}
                  class="text-2xl transition-all duration-200 transform hover:scale-110
                    {isActive ? 'text-blue-400' : 'text-slate-700'}
                    {hasSubmitted && !allowMultipleAttempts ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'}
                  "
                >
                  {isActive ? '★' : '☆'}
                </button>
              {/each}
            </div>
            <span class="text-sm text-slate-400 ml-2">
              {userRating || 0}/5
            </span>
          </div>

          <!-- Mostrar rating esperado si configurado -->
          {#if hasSubmitted && showExpectedRatings && result}
            <div class="mt-3 pt-3 border-t border-slate-700">
              <div class="flex items-center gap-2 text-xs">
                <span class="text-slate-500">Evaluación esperada:</span>
                <div class="flex gap-1">
                  {#each [1, 2, 3, 4, 5] as star}
                    <span class="text-sm {criterion.expectedRating >= star ? 'text-slate-400' : 'text-slate-800'}">
                      {criterion.expectedRating >= star ? '★' : '☆'}
                    </span>
                  {/each}
                </div>
                <span class="text-slate-400">({criterion.expectedRating}/5)</span>
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </div>

  <!-- Progreso -->
  {#if !hasSubmitted}
    <div class="mb-6 p-3 bg-canvas-900/50 rounded-xl border border-slate-700">
      <div class="flex items-center justify-between text-sm">
        <span class="text-slate-400">
          Criterios evaluados: {Object.keys(ratings).length}/{criteria.length}
        </span>
        {#if allRated}
          <span class="text-green-400 font-semibold">✓ Listo para enviar</span>
        {/if}
      </div>
    </div>
  {/if}

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={!allRated}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-blue-500 to-indigo-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-blue-500/50
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Enviar Evaluación
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
        Revisar Evaluación
      </button>
    {/if}
  </div>

  <!-- Feedback final -->
  {#if hasSubmitted && showFeedback}
    <div class="mt-6 p-4 rounded-2xl {allCorrect ? 'bg-blue-500/10 border border-blue-500/50' : 'bg-yellow-500/10 border border-yellow-500/50'}">
      <div class="flex items-center gap-3">
        <span class="text-3xl">
          {allCorrect ? '🎯' : '🤔'}
        </span>
        <div class="flex-1">
          <p class="{allCorrect ? 'text-blue-400' : 'text-yellow-400'} font-semibold mb-1">
            {allCorrect
              ? '¡Excelente evaluación! Tu juicio crítico está muy alineado'
              : 'Revisa algunos criterios. Considera diferentes perspectivas'}
          </p>
          <p class="text-slate-400 text-sm mt-1">
            {Object.values(results).filter(r => r.status === 'perfect').length} criterios perfectos,
            {Object.values(results).filter(r => r.status === 'close').length} cercanos,
            {Object.values(results).filter(r => r.status === 'incorrect').length} por mejorar
          </p>
          <div class="flex items-center gap-3 mt-2">
            <div class="flex-1 bg-canvas-900 rounded-full h-2 overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-blue-500 to-indigo-500 transition-all duration-500"
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
