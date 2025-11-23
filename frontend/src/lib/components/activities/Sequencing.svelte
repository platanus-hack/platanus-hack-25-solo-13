<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    // Data del componente
    title = "Ordena los eventos cronológicamente",
    items = [
      { id: 1, content: "Independencia de Chile (1818)", correctOrder: 1 },
      { id: 2, content: "Batalla de Maipú (1818)", correctOrder: 2 },
      { id: 3, content: "Constitución de 1833", correctOrder: 3 },
      { id: 4, content: "Guerra del Pacífico (1879-1884)", correctOrder: 4 }
    ],

    // Metadata educativa
    bloomLevel = "comprender",
    materia = "historia",
    oaId = null,

    // Configuración
    shuffleItems = true,
    showNumbers = true,
    showHints = false,
    allowMultipleAttempts = true,

    // Callbacks
    onAnswer = null,
    onComplete = null
  } = $props();

  // Estados locales
  let currentOrder = $state([]);
  let hasSubmitted = $state(false);
  let results = $state([]);
  let attemptCount = $state(0);
  let containerRef = $state(null);
  let draggedIndex = $state(null);
  let dragOverIndex = $state(null);

  // Inicializar orden (shuffleado o no)
  $effect(() => {
    const ordered = [...items];
    if (shuffleItems) {
      // Fisher-Yates shuffle
      for (let i = ordered.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [ordered[i], ordered[j]] = [ordered[j], ordered[i]];
      }
    }
    currentOrder = ordered;
  });

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

  // Calcular score
  const score = $derived.by(() => {
    if (!hasSubmitted) return 0;
    const correct = results.filter(r => r).length;
    return (correct / items.length) * 100;
  });

  const allCorrect = $derived(score === 100);

  // Funciones de drag and drop
  function handleDragStart(event, index) {
    if (hasSubmitted && !allowMultipleAttempts) return;
    draggedIndex = index;
    event.dataTransfer.effectAllowed = 'move';
  }

  function handleDragOver(event, index) {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    dragOverIndex = index;
  }

  function handleDragLeave() {
    dragOverIndex = null;
  }

  function handleDrop(event, dropIndex) {
    event.preventDefault();
    if (draggedIndex !== null && draggedIndex !== dropIndex) {
      const newOrder = [...currentOrder];
      const draggedItem = newOrder[draggedIndex];

      // Remover el item de su posición original
      newOrder.splice(draggedIndex, 1);

      // Insertar en la nueva posición
      newOrder.splice(dropIndex, 0, draggedItem);

      currentOrder = newOrder;
    }

    draggedIndex = null;
    dragOverIndex = null;
  }

  function handleDragEnd() {
    draggedIndex = null;
    dragOverIndex = null;
  }

  // Mover items con botones
  function moveUp(index) {
    if (index > 0 && (!hasSubmitted || allowMultipleAttempts)) {
      const newOrder = [...currentOrder];
      [newOrder[index - 1], newOrder[index]] = [newOrder[index], newOrder[index - 1]];
      currentOrder = newOrder;
    }
  }

  function moveDown(index) {
    if (index < currentOrder.length - 1 && (!hasSubmitted || allowMultipleAttempts)) {
      const newOrder = [...currentOrder];
      [newOrder[index], newOrder[index + 1]] = [newOrder[index + 1], newOrder[index]];
      currentOrder = newOrder;
    }
  }

  // Verificar respuestas
  function handleSubmit() {
    hasSubmitted = true;
    attemptCount++;

    // Verificar cada posición
    const newResults = currentOrder.map((item, index) => {
      return item.correctOrder === index + 1;
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
        userOrder: currentOrder.map(item => item.id),
        correctOrder: items.sort((a, b) => a.correctOrder - b.correctOrder).map(item => item.id),
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
      results = [];
    }
  }

  function handleShowCorrect() {
    const correctOrder = [...items].sort((a, b) => a.correctOrder - b.correctOrder);
    currentOrder = correctOrder;
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

  <!-- Título e instrucciones -->
  <div class="mb-8">
    <h3 class="text-xl font-semibold text-white mb-2">
      {title}
    </h3>
    <p class="text-slate-400 text-sm">
      Arrastra los elementos para ordenarlos correctamente o usa las flechas
    </p>
  </div>

  <!-- Items a ordenar -->
  <div class="space-y-3 mb-6">
    {#each currentOrder as item, index (item.id)}
      {@const isDragging = draggedIndex === index}
      {@const isDragOver = dragOverIndex === index}
      {@const isCorrect = hasSubmitted && results[index]}
      {@const isIncorrect = hasSubmitted && !results[index]}

      <div
        draggable={!hasSubmitted || allowMultipleAttempts}
        ondragstart={(e) => handleDragStart(e, index)}
        ondragover={(e) => handleDragOver(e, index)}
        ondragleave={handleDragLeave}
        ondrop={(e) => handleDrop(e, index)}
        ondragend={handleDragEnd}
        class="
          flex items-center gap-4 p-4 rounded-2xl border-2
          transition-all duration-300
          {isDragging ? 'opacity-50 scale-95' : ''}
          {isDragOver ? 'border-cyan-500 bg-cyan-500/10 scale-105' : 'border-slate-700'}
          {isCorrect ? 'border-green-500 bg-green-500/20' : ''}
          {isIncorrect ? 'border-red-500 bg-red-500/20' : ''}
          {!hasSubmitted || allowMultipleAttempts ? 'cursor-move hover:border-slate-600 bg-canvas-900/50' : 'cursor-not-allowed bg-canvas-900/30'}
        "
      >
        <!-- Número de posición -->
        {#if showNumbers}
          <div class="flex-shrink-0 w-10 h-10 rounded-full bg-canvas-800 border border-slate-600 flex items-center justify-center">
            <span class="text-white font-bold">{index + 1}</span>
          </div>
        {/if}

        <!-- Contenido del item -->
        <div class="flex-1">
          <p class="text-white font-medium">{item.content}</p>
        </div>

        <!-- Controles de movimiento -->
        {#if !hasSubmitted || allowMultipleAttempts}
          <div class="flex flex-col gap-1">
            <button
              onclick={() => moveUp(index)}
              disabled={index === 0}
              class="
                p-1 rounded bg-canvas-800 border border-slate-600
                hover:bg-slate-700 hover:border-cyan-500
                disabled:opacity-30 disabled:cursor-not-allowed
                transition-all duration-300
              "
            >
              <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
              </svg>
            </button>
            <button
              onclick={() => moveDown(index)}
              disabled={index === currentOrder.length - 1}
              class="
                p-1 rounded bg-canvas-800 border border-slate-600
                hover:bg-slate-700 hover:border-cyan-500
                disabled:opacity-30 disabled:cursor-not-allowed
                transition-all duration-300
              "
            >
              <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </button>
          </div>
        {/if}

        <!-- Feedback visual -->
        {#if hasSubmitted}
          <div class="flex-shrink-0">
            <span class="text-2xl">
              {isCorrect ? '✓' : '✗'}
            </span>
          </div>
        {/if}
      </div>
    {/each}
  </div>

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-[#E1E1E1] hover:bg-[#CCCCCC]
          text-canvas-900
          transition-all duration-300
          shadow-lg hover:shadow-xl
        "
      >
        Verificar Orden
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
        Intentar de Nuevo
      </button>
      {#if showHints}
        <button
          onclick={handleShowCorrect}
          class="
            px-6 py-3 rounded-xl font-semibold
            bg-yellow-500/20 text-yellow-400
            border border-yellow-500
            transition-all duration-300
            hover:bg-yellow-500/30
          "
        >
          Mostrar Orden Correcto
        </button>
      {/if}
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
            {allCorrect ? '¡Perfecto! El orden es correcto' : `${results.filter(r => r).length} de ${items.length} en posición correcta`}
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

          {#if !allCorrect && showHints}
            <p class="text-slate-400 text-sm mt-3">
              💡 Los elementos con ✗ están en posición incorrecta
            </p>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
