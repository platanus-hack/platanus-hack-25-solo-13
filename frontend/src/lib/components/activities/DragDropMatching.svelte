<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    // Data del componente
    title = "Relaciona los conceptos con sus definiciones",
    pairs = [
      { id: 1, term: "Fotosíntesis", definition: "Proceso de conversión de luz solar en energía química" },
      { id: 2, term: "Respiración Celular", definition: "Proceso de obtención de energía de moléculas orgánicas" },
      { id: 3, term: "Mitosis", definition: "División celular que produce dos células hijas idénticas" }
    ],

    // Metadata educativa
    bloomLevel = "comprender",
    materia = "biología",
    oaId = null,

    // Configuración
    shuffleOptions = true,
    showFeedback = true,
    allowMultipleAttempts = true,

    // Callbacks
    onAnswer = null,
    onComplete = null
  } = $props();

  // Estados locales
  let matches = $state({}); // { termId: definitionId }
  let hasSubmitted = $state(false);
  let results = $state({});
  let attemptCount = $state(0);
  let containerRef = $state(null);
  let draggedItem = $state(null);
  let dropTarget = $state(null);

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
  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.comprender;

  // Shufflear definiciones
  const shuffledDefinitions = $derived.by(() => {
    const defs = [...pairs];
    if (shuffleOptions) {
      // Fisher-Yates shuffle
      for (let i = defs.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [defs[i], defs[j]] = [defs[j], defs[i]];
      }
    }
    return defs;
  });

  // Calcular score
  const score = $derived.by(() => {
    if (!hasSubmitted) return 0;
    const correct = Object.values(results).filter(r => r).length;
    return (correct / pairs.length) * 100;
  });

  const allCorrect = $derived(score === 100);

  // Funciones de drag and drop
  function handleDragStart(event, definitionId) {
    if (hasSubmitted && !allowMultipleAttempts) return;
    draggedItem = definitionId;
    event.dataTransfer.effectAllowed = 'move';
  }

  function handleDragOver(event, termId) {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    dropTarget = termId;
  }

  function handleDragLeave() {
    dropTarget = null;
  }

  function handleDrop(event, termId) {
    event.preventDefault();
    if (draggedItem !== null) {
      // Crear nuevo objeto de matches
      const newMatches = { ...matches };

      // Remover cualquier match previo de esta definición
      Object.keys(newMatches).forEach(key => {
        if (newMatches[key] === draggedItem) {
          delete newMatches[key];
        }
      });

      // Asignar nuevo match
      newMatches[termId] = draggedItem;
      matches = newMatches;

      draggedItem = null;
      dropTarget = null;
    }
  }

  function handleDragEnd() {
    draggedItem = null;
    dropTarget = null;
  }

  function handleRemoveMatch(termId) {
    const newMatches = { ...matches };
    delete newMatches[termId];
    matches = newMatches;
  }

  // Verificar respuestas
  function handleSubmit() {
    if (Object.keys(matches).length !== pairs.length) {
      alert("Por favor, relaciona todos los conceptos antes de enviar.");
      return;
    }

    hasSubmitted = true;
    attemptCount++;

    // Verificar cada match
    const newResults = {};
    Object.keys(matches).forEach(termId => {
      const definitionId = matches[termId];
      newResults[termId] = parseInt(termId) === definitionId;
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
        matches: matches,
        results: results,
        score: score,
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
      // Mantener solo los matches correctos
      const newMatches = {};
      Object.keys(matches).forEach(termId => {
        if (results[termId]) {
          newMatches[termId] = matches[termId];
        }
      });
      matches = newMatches;
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
  class="w-full max-w-4xl mx-auto p-6 bg-slate-950 rounded-3xl border border-slate-800 shadow-2xl"
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
  <div class="mb-8">
    <h3 class="text-xl font-semibold text-white mb-2">
      {title}
    </h3>
    <p class="text-slate-400 text-sm">
      Arrastra las definiciones y suéltalas en el concepto correspondiente
    </p>
  </div>

  <!-- Área de matching -->
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
    <!-- Columna izquierda: Términos (drop zones) -->
    <div class="space-y-3">
      <p class="text-xs font-semibold text-slate-400 mb-3 uppercase">Conceptos:</p>
      {#each pairs as pair (pair.id)}
        {@const isMatched = matches[pair.id] !== undefined}
        {@const matchedDefinition = isMatched ? shuffledDefinitions.find(d => d.id === matches[pair.id]) : null}
        {@const isCorrect = hasSubmitted && results[pair.id]}
        {@const isIncorrect = hasSubmitted && !results[pair.id]}
        {@const isDropTarget = dropTarget === pair.id}

        <div
          ondragover={(e) => handleDragOver(e, pair.id)}
          ondragleave={handleDragLeave}
          ondrop={(e) => handleDrop(e, pair.id)}
          class="
            p-4 rounded-2xl border-2 border-dashed
            transition-all duration-300 min-h-[100px]
            {isDropTarget ? 'border-cyan-500 bg-cyan-500/10 scale-105' : 'border-slate-700'}
            {isCorrect ? 'border-green-500 bg-green-500/20' : ''}
            {isIncorrect ? 'border-red-500 bg-red-500/20' : ''}
            {!isMatched && !isDropTarget ? 'bg-slate-900/50' : ''}
          "
        >
          <div class="flex items-start justify-between gap-3">
            <div class="flex-1">
              <p class="text-white font-semibold mb-2">{pair.term}</p>

              {#if isMatched}
                <div class="p-3 bg-slate-800 rounded-xl border border-slate-600 relative">
                  <p class="text-slate-300 text-sm pr-6">{matchedDefinition?.definition}</p>

                  {#if !hasSubmitted || allowMultipleAttempts}
                    <button
                      onclick={() => handleRemoveMatch(pair.id)}
                      class="absolute top-2 right-2 text-slate-500 hover:text-red-400"
                    >
                      ✕
                    </button>
                  {/if}
                </div>
              {:else}
                <p class="text-slate-500 text-sm italic">Arrastra una definición aquí</p>
              {/if}
            </div>

            {#if hasSubmitted}
              <span class="text-2xl">
                {isCorrect ? '✓' : '✗'}
              </span>
            {/if}
          </div>
        </div>
      {/each}
    </div>

    <!-- Columna derecha: Definiciones disponibles (draggable) -->
    <div class="space-y-3">
      <p class="text-xs font-semibold text-slate-400 mb-3 uppercase">Definiciones:</p>
      <div class="space-y-3">
        {#each shuffledDefinitions as def (def.id)}
          {@const isUsed = Object.values(matches).includes(def.id)}
          {@const isDragging = draggedItem === def.id}

          {#if !isUsed || (hasSubmitted && !allowMultipleAttempts)}
            <div
              draggable={!hasSubmitted || allowMultipleAttempts}
              ondragstart={(e) => handleDragStart(e, def.id)}
              ondragend={handleDragEnd}
              class="
                p-4 rounded-2xl border-2
                transition-all duration-300
                {isDragging ? 'opacity-50 scale-95 border-cyan-500' : 'border-slate-700'}
                {!hasSubmitted || allowMultipleAttempts ? 'cursor-move hover:border-cyan-500 hover:bg-slate-900/50' : 'cursor-not-allowed opacity-50'}
                bg-slate-900
              "
            >
              <p class="text-slate-300 text-sm">{def.definition}</p>
            </div>
          {/if}
        {/each}
      </div>

      {#if Object.values(matches).length === pairs.length}
        <div class="p-4 bg-green-500/10 rounded-xl border border-green-500/50">
          <p class="text-green-400 text-sm text-center">
            ✓ Todos los conceptos relacionados
          </p>
        </div>
      {/if}
    </div>
  </div>

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={Object.keys(matches).length !== pairs.length}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-cyan-500 to-blue-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-cyan-500/50
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
          bg-slate-800 text-white
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
  {#if hasSubmitted && showFeedback}
    <div class="mt-6 p-4 rounded-2xl {allCorrect ? 'bg-green-500/10 border border-green-500/50' : 'bg-yellow-500/10 border border-yellow-500/50'}">
      <div class="flex items-center gap-3">
        <span class="text-3xl">
          {allCorrect ? '🎉' : '📝'}
        </span>
        <div class="flex-1">
          <p class="{allCorrect ? 'text-green-400' : 'text-yellow-400'} font-semibold mb-1">
            {allCorrect ? '¡Perfecto! Todas las relaciones son correctas' : `${Object.values(results).filter(r => r).length} de ${pairs.length} correctas`}
          </p>
          <div class="flex items-center gap-3 mt-2">
            <div class="flex-1 bg-slate-900 rounded-full h-2 overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-cyan-500 to-green-500 transition-all duration-500"
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
