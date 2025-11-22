<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    // Data del componente
    title = "Compara y contrasta los siguientes conceptos",
    itemA = { name: "Concepto A", color: "cyan" },
    itemB = { name: "Concepto B", color: "purple" },
    characteristics = [
      { id: 1, text: "Característica compartida", correctColumn: "both" },
      { id: 2, text: "Característica única de A", correctColumn: "A" },
      { id: 3, text: "Característica única de B", correctColumn: "B" }
    ],

    // Metadata educativa
    bloomLevel = "analizar",
    materia = "biología",
    oaId = null,

    // Configuración
    showFeedback = true,
    allowMultipleAttempts = true,

    // Callbacks
    onAnswer = null,
    onComplete = null
  } = $props();

  // Estados locales
  let placements = $state({}); // { characteristicId: "A" | "B" | "both" }
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
  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.analizar;

  // Características disponibles (no colocadas)
  const availableCharacteristics = $derived.by(() => {
    return characteristics.filter(char => placements[char.id] === undefined);
  });

  // Características por columna
  const columnA = $derived.by(() => {
    return characteristics.filter(char => placements[char.id] === 'A');
  });

  const columnB = $derived.by(() => {
    return characteristics.filter(char => placements[char.id] === 'B');
  });

  const columnBoth = $derived.by(() => {
    return characteristics.filter(char => placements[char.id] === 'both');
  });

  // Calcular score
  const score = $derived.by(() => {
    if (!hasSubmitted) return 0;
    const correct = Object.values(results).filter(r => r).length;
    return (correct / characteristics.length) * 100;
  });

  const allCorrect = $derived(score === 100);
  const allPlaced = $derived(Object.keys(placements).length === characteristics.length);

  // Funciones de drag and drop
  function handleDragStart(event, characteristicId) {
    if (hasSubmitted && !allowMultipleAttempts) return;
    draggedItem = characteristicId;
    event.dataTransfer.effectAllowed = 'move';
  }

  function handleDragOver(event, column) {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
    dropTarget = column;
  }

  function handleDragLeave() {
    dropTarget = null;
  }

  function handleDrop(event, column) {
    event.preventDefault();
    if (draggedItem !== null) {
      placements = { ...placements, [draggedItem]: column };
      draggedItem = null;
      dropTarget = null;
    }
  }

  function handleDragEnd() {
    draggedItem = null;
    dropTarget = null;
  }

  function handleRemoveItem(characteristicId) {
    const newPlacements = { ...placements };
    delete newPlacements[characteristicId];
    placements = newPlacements;
  }

  // Verificar respuestas
  function handleSubmit() {
    if (!allPlaced) {
      alert("Por favor, coloca todas las características antes de enviar.");
      return;
    }

    hasSubmitted = true;
    attemptCount++;

    // Verificar cada placement
    const newResults = {};
    Object.keys(placements).forEach(charId => {
      const characteristic = characteristics.find(c => c.id === parseInt(charId));
      newResults[charId] = placements[charId] === characteristic.correctColumn;
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
        placements: placements,
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
      // Mantener solo los placements correctos
      const newPlacements = {};
      Object.keys(placements).forEach(charId => {
        if (results[charId]) {
          newPlacements[charId] = placements[charId];
        }
      });
      placements = newPlacements;
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
  class="w-full max-w-6xl mx-auto p-6 bg-canvas-950 rounded-2xl border border-slate-800 shadow-2xl"
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
      Arrastra cada característica a la columna correspondiente: solo {itemA.name}, solo {itemB.name}, o Ambos
    </p>
  </div>

  <!-- Área de comparación -->
  <div class="grid grid-cols-1 lg:grid-cols-4 gap-4 mb-6">
    <!-- Columna izquierda: Características disponibles -->
    <div class="lg:col-span-1 space-y-3">
      <div class="sticky top-4">
        <p class="text-xs font-semibold text-slate-400 mb-3 uppercase">Características:</p>
        <div class="space-y-2">
          {#each availableCharacteristics as char (char.id)}
            {@const isDragging = draggedItem === char.id}
            <div
              draggable={!hasSubmitted || allowMultipleAttempts}
              ondragstart={(e) => handleDragStart(e, char.id)}
              ondragend={handleDragEnd}
              class="
                p-3 rounded-xl border-2
                transition-all duration-300
                {isDragging ? 'opacity-50 scale-95 border-green-500' : 'border-slate-700'}
                {!hasSubmitted || allowMultipleAttempts ? 'cursor-move hover:border-green-500 hover:bg-canvas-900/50' : 'cursor-not-allowed opacity-50'}
                bg-canvas-900
              "
            >
              <p class="text-slate-300 text-sm">{char.text}</p>
            </div>
          {/each}

          {#if allPlaced}
            <div class="p-3 bg-green-500/10 rounded-xl border border-green-500/50">
              <p class="text-green-400 text-xs text-center">
                ✓ Todo clasificado
              </p>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Columnas de comparación -->
    <div class="lg:col-span-3 grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Columna A -->
      <div
        ondragover={(e) => handleDragOver(e, 'A')}
        ondragleave={handleDragLeave}
        ondrop={(e) => handleDrop(e, 'A')}
        class="
          p-4 rounded-2xl border-2 border-dashed min-h-[300px]
          transition-all duration-300
          {dropTarget === 'A' ? 'border-cyan-500 bg-cyan-500/10 scale-105' : 'border-slate-700 bg-canvas-900/30'}
        "
      >
        <div class="mb-4 pb-3 border-b border-slate-700">
          <h4 class="text-lg font-bold text-{itemA.color || 'cyan'}-400">
            {itemA.name}
          </h4>
          <p class="text-xs text-slate-500 mt-1">Solo este concepto</p>
        </div>

        <div class="space-y-2">
          {#each columnA as char (char.id)}
            {@const isCorrect = hasSubmitted && results[char.id]}
            {@const isIncorrect = hasSubmitted && !results[char.id]}
            <div class="
              p-3 rounded-xl border
              {isCorrect ? 'bg-green-500/20 border-green-500' : ''}
              {isIncorrect ? 'bg-red-500/20 border-red-500' : ''}
              {!hasSubmitted ? 'bg-canvas-800 border-slate-600' : ''}
              relative
            ">
              <p class="text-slate-300 text-sm pr-8">{char.text}</p>

              {#if hasSubmitted}
                <span class="absolute top-2 right-2 text-lg">
                  {isCorrect ? '✓' : '✗'}
                </span>
              {:else if !hasSubmitted || allowMultipleAttempts}
                <button
                  onclick={() => handleRemoveItem(char.id)}
                  class="absolute top-2 right-2 text-slate-500 hover:text-red-400"
                >
                  ✕
                </button>
              {/if}
            </div>
          {/each}
        </div>
      </div>

      <!-- Columna BOTH -->
      <div
        ondragover={(e) => handleDragOver(e, 'both')}
        ondragleave={handleDragLeave}
        ondrop={(e) => handleDrop(e, 'both')}
        class="
          p-4 rounded-2xl border-2 border-dashed min-h-[300px]
          transition-all duration-300
          {dropTarget === 'both' ? 'border-purple-500 bg-purple-500/10 scale-105' : 'border-slate-700 bg-canvas-900/30'}
        "
      >
        <div class="mb-4 pb-3 border-b border-slate-700">
          <h4 class="text-lg font-bold text-purple-400">
            Ambos
          </h4>
          <p class="text-xs text-slate-500 mt-1">Compartido por ambos</p>
        </div>

        <div class="space-y-2">
          {#each columnBoth as char (char.id)}
            {@const isCorrect = hasSubmitted && results[char.id]}
            {@const isIncorrect = hasSubmitted && !results[char.id]}
            <div class="
              p-3 rounded-xl border
              {isCorrect ? 'bg-green-500/20 border-green-500' : ''}
              {isIncorrect ? 'bg-red-500/20 border-red-500' : ''}
              {!hasSubmitted ? 'bg-canvas-800 border-slate-600' : ''}
              relative
            ">
              <p class="text-slate-300 text-sm pr-8">{char.text}</p>

              {#if hasSubmitted}
                <span class="absolute top-2 right-2 text-lg">
                  {isCorrect ? '✓' : '✗'}
                </span>
              {:else if !hasSubmitted || allowMultipleAttempts}
                <button
                  onclick={() => handleRemoveItem(char.id)}
                  class="absolute top-2 right-2 text-slate-500 hover:text-red-400"
                >
                  ✕
                </button>
              {/if}
            </div>
          {/each}
        </div>
      </div>

      <!-- Columna B -->
      <div
        ondragover={(e) => handleDragOver(e, 'B')}
        ondragleave={handleDragLeave}
        ondrop={(e) => handleDrop(e, 'B')}
        class="
          p-4 rounded-2xl border-2 border-dashed min-h-[300px]
          transition-all duration-300
          {dropTarget === 'B' ? 'border-green-500 bg-green-500/10 scale-105' : 'border-slate-700 bg-canvas-900/30'}
        "
      >
        <div class="mb-4 pb-3 border-b border-slate-700">
          <h4 class="text-lg font-bold text-{itemB.color || 'green'}-400">
            {itemB.name}
          </h4>
          <p class="text-xs text-slate-500 mt-1">Solo este concepto</p>
        </div>

        <div class="space-y-2">
          {#each columnB as char (char.id)}
            {@const isCorrect = hasSubmitted && results[char.id]}
            {@const isIncorrect = hasSubmitted && !results[char.id]}
            <div class="
              p-3 rounded-xl border
              {isCorrect ? 'bg-green-500/20 border-green-500' : ''}
              {isIncorrect ? 'bg-red-500/20 border-red-500' : ''}
              {!hasSubmitted ? 'bg-canvas-800 border-slate-600' : ''}
              relative
            ">
              <p class="text-slate-300 text-sm pr-8">{char.text}</p>

              {#if hasSubmitted}
                <span class="absolute top-2 right-2 text-lg">
                  {isCorrect ? '✓' : '✗'}
                </span>
              {:else if !hasSubmitted || allowMultipleAttempts}
                <button
                  onclick={() => handleRemoveItem(char.id)}
                  class="absolute top-2 right-2 text-slate-500 hover:text-red-400"
                >
                  ✕
                </button>
              {/if}
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={!allPlaced}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-green-500 to-emerald-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-green-500/50
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Verificar Clasificación
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
  {#if hasSubmitted && showFeedback}
    <div class="mt-6 p-4 rounded-2xl {allCorrect ? 'bg-green-500/10 border border-green-500/50' : 'bg-yellow-500/10 border border-yellow-500/50'}">
      <div class="flex items-center gap-3">
        <span class="text-3xl">
          {allCorrect ? '🎉' : '📊'}
        </span>
        <div class="flex-1">
          <p class="{allCorrect ? 'text-green-400' : 'text-yellow-400'} font-semibold mb-1">
            {allCorrect ? '¡Excelente análisis! Todas las clasificaciones son correctas' : `${Object.values(results).filter(r => r).length} de ${characteristics.length} correctas`}
          </p>
          <div class="flex items-center gap-3 mt-2">
            <div class="flex-1 bg-canvas-900 rounded-full h-2 overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-green-500 to-emerald-500 transition-all duration-500"
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
