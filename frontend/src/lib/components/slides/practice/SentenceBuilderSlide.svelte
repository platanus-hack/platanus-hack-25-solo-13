<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Construcción de Oraciones",
    modoEjercicio = "free", // free | guided | transformation
    oracionBase = null,
    palabrasDisponibles = [],
    objetivoGramatical = "",
    variacionesCorrectas = [],
    mostrarScaffolding = true,
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let palabrasBanco = $state([...palabrasDisponibles]); // Palabras disponibles
  let oracionConstruida = $state([]); // Palabras en la oración
  let draggedWord = $state(null);
  let dropZone = $state(null);
  let mostrarVariaciones = $state(false);
  let feedback = $state(null); // { tipo: "success" | "info", mensaje }

  // Color mapping por tipo gramatical
  const tipoColors = {
    sujeto: { bg: "bg-blue-500", text: "text-white", border: "border-blue-600" },
    verbo: { bg: "bg-red-500", text: "text-white", border: "border-red-600" },
    objeto: { bg: "bg-green-500", text: "text-white", border: "border-green-600" },
    adjetivo: { bg: "bg-purple-500", text: "text-white", border: "border-purple-600" },
    adverbio: { bg: "bg-orange-500", text: "text-white", border: "border-orange-600" },
    complemento: { bg: "bg-teal-500", text: "text-white", border: "border-teal-600" }
  };

  // Oración como string
  const oracionTexto = $derived(
    oracionConstruida.map(p => p.texto).join(' ')
  );

  // Verificar si es válida
  const esOracionValida = $derived.by(() => {
    if (oracionConstruida.length === 0) return false;
    if (variacionesCorrectas.length === 0) return true; // Modo libre

    const oracionNormalizada = oracionTexto.toLowerCase().trim();
    return variacionesCorrectas.some(variacion =>
      variacion.toLowerCase().trim() === oracionNormalizada
    );
  });

  // Drag & Drop handlers
  function handleDragStart(event, palabra) {
    draggedWord = palabra;
    event.dataTransfer.effectAllowed = 'move';
  }

  function handleDragEnd() {
    draggedWord = null;
    dropZone = null;
  }

  function handleDragOver(event, zone) {
    event.preventDefault();
    dropZone = zone;
  }

  function handleDragLeave() {
    dropZone = null;
  }

  function handleDrop(event, zone) {
    event.preventDefault();
    if (!draggedWord) return;

    if (zone === 'oracion') {
      // Mover del banco a la oración
      oracionConstruida = [...oracionConstruida, draggedWord];
      palabrasBanco = palabrasBanco.filter(p => p !== draggedWord);
    } else if (zone === 'banco') {
      // Devolver a banco
      if (oracionConstruida.includes(draggedWord)) {
        oracionConstruida = oracionConstruida.filter(p => p !== draggedWord);
        palabrasBanco = [...palabrasBanco, draggedWord];
      }
    }

    draggedWord = null;
    dropZone = null;
  }

  function removerPalabra(palabra) {
    oracionConstruida = oracionConstruida.filter(p => p !== palabra);
    palabrasBanco = [...palabrasBanco, palabra];
  }

  function verificarOracion() {
    if (esOracionValida) {
      feedback = {
        tipo: "success",
        mensaje: "¡Excelente! Tu oración es gramaticalmente correcta."
      };
    } else {
      feedback = {
        tipo: "info",
        mensaje: `Tu oración: "${oracionTexto}". Revisa si cumple el objetivo gramatical.`
      };
    }

    // Limpiar feedback después de 5 segundos
    setTimeout(() => {
      feedback = null;
    }, 5000);
  }

  function reiniciar() {
    oracionConstruida = [];
    palabrasBanco = [...palabrasDisponibles];
    feedback = null;
  }

  // Animación de entrada
  onMount(() => {
    if (containerRef) {
      gsap.from(containerRef, {
        opacity: 0,
        y: 30,
        duration: 0.6,
        ease: 'power2.out'
      });
    }
  });
</script>

<div
  bind:this={containerRef}
  class="w-full max-w-7xl mx-auto p-8 bg-slate-950 rounded-2xl border border-canvas-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-purple-500/20 text-purple-400 border border-purple-500">
        {materia}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-medium bg-canvas-800 text-canvas-300">
        Gramática
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    {#if objetivoGramatical}
      <p class="text-canvas-300">
        🎯 {objetivoGramatical}
      </p>
    {/if}
  </div>

  <!-- Oración base (modo transformation) -->
  {#if modoEjercicio === "transformation" && oracionBase}
    <div class="mb-6 p-4 bg-purple-500/10 rounded-xl border border-purple-500/30">
      <h3 class="text-sm font-semibold text-purple-300 mb-2">Oración base:</h3>
      <p class="text-xl text-white font-medium">"{oracionBase}"</p>
    </div>
  {/if}

  <!-- Zona de construcción de oración -->
  <div class="mb-6">
    <h3 class="text-lg font-semibold text-white mb-3">Tu Oración:</h3>

    <div
      ondragover={(e) => handleDragOver(e, 'oracion')}
      ondragleave={handleDragLeave}
      ondrop={(e) => handleDrop(e, 'oracion')}
      class="
        min-h-32 p-6 rounded-2xl border-2 border-dashed transition-all duration-300
        {dropZone === 'oracion' ? 'border-purple-500 bg-purple-500/10' : 'border-canvas-700 bg-canvas-900/50'}
      "
    >
      {#if oracionConstruida.length === 0}
        <p class="text-center text-canvas-500 italic">
          Arrastra palabras aquí para construir tu oración...
        </p>
      {:else}
        <div class="flex flex-wrap items-center gap-2">
          {#each oracionConstruida as palabra, index}
            <button
              onclick={() => removerPalabra(palabra)}
              class="
                group relative px-4 py-2 rounded-lg font-semibold
                {tipoColors[palabra.tipo]?.bg || 'bg-canvas-700'}
                {tipoColors[palabra.tipo]?.text || 'text-white'}
                border-2 {tipoColors[palabra.tipo]?.border || 'border-canvas-600'}
                transition-all duration-300
                hover:scale-110 hover:shadow-lg
              "
              title="Click para remover"
            >
              {palabra.texto}
              <span class="absolute -top-2 -right-2 w-5 h-5 bg-red-500 rounded-full text-xs flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                ✕
              </span>
            </button>
          {/each}
        </div>

        <!-- Oración como texto -->
        <div class="mt-4 p-4 bg-slate-950 rounded-lg border border-canvas-700">
          <p class="text-xl text-white font-medium">
            "{oracionTexto}"
          </p>
        </div>
      {/if}
    </div>

    <!-- Botones de acción -->
    <div class="flex items-center gap-3 mt-4">
      <button
        onclick={verificarOracion}
        disabled={oracionConstruida.length === 0}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-purple-500 to-purple-600
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-purple-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed
        "
      >
        ✓ Verificar Oración
      </button>

      <button
        onclick={reiniciar}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-canvas-800 text-canvas-300
          border border-canvas-700
          transition-all duration-300
          hover:bg-canvas-700 hover:scale-105
        "
      >
        ↻ Reiniciar
      </button>

      {#if variacionesCorrectas.length > 0}
        <button
          onclick={() => mostrarVariaciones = !mostrarVariaciones}
          class="
            px-6 py-3 rounded-xl font-semibold text-sm
            bg-canvas-800 text-canvas-300
            border border-canvas-700
            transition-all duration-300
            hover:bg-canvas-700 hover:scale-105
          "
        >
          {mostrarVariaciones ? '✕ Cerrar' : '💡'} Ver Variaciones Correctas
        </button>
      {/if}
    </div>

    <!-- Feedback -->
    {#if feedback}
      <div class="
        mt-4 p-4 rounded-xl border-2
        {feedback.tipo === 'success' ? 'bg-green-500/10 border-green-500' : 'bg-blue-500/10 border-blue-500'}
      ">
        <p class="
          text-sm font-semibold
          {feedback.tipo === 'success' ? 'text-green-300' : 'text-blue-300'}
        ">
          {feedback.mensaje}
        </p>
      </div>
    {/if}
  </div>

  <!-- Banco de palabras -->
  <div class="mb-6">
    <h3 class="text-lg font-semibold text-white mb-3">Palabras Disponibles:</h3>

    <div
      ondragover={(e) => handleDragOver(e, 'banco')}
      ondragleave={handleDragLeave}
      ondrop={(e) => handleDrop(e, 'banco')}
      class="
        p-6 rounded-2xl border-2 transition-all duration-300
        {dropZone === 'banco' ? 'border-purple-500 bg-purple-500/10' : 'border-canvas-700 bg-canvas-900/50'}
      "
    >
      {#if palabrasBanco.length === 0}
        <p class="text-center text-canvas-500 italic">
          Todas las palabras han sido usadas
        </p>
      {:else}
        <div class="flex flex-wrap gap-3">
          {#each palabrasBanco as palabra}
            <div
              draggable="true"
              ondragstart={(e) => handleDragStart(e, palabra)}
              ondragend={handleDragEnd}
              class="
                group relative px-4 py-3 rounded-lg font-semibold cursor-move
                {tipoColors[palabra.tipo]?.bg || 'bg-canvas-700'}
                {tipoColors[palabra.tipo]?.text || 'text-white'}
                border-2 {tipoColors[palabra.tipo]?.border || 'border-canvas-600'}
                transition-all duration-300
                hover:scale-110 hover:shadow-lg
                {draggedWord === palabra ? 'opacity-50' : 'opacity-100'}
              "
              title={palabra.tipo}
            >
              {palabra.texto}
              <span class="absolute -bottom-6 left-0 right-0 text-xs text-canvas-500 text-center opacity-0 group-hover:opacity-100 transition-opacity">
                {palabra.tipo}
              </span>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- Leyenda de colores -->
    <div class="mt-4 p-4 bg-canvas-900/50 rounded-xl border border-canvas-700">
      <h4 class="text-xs font-semibold uppercase text-canvas-400 mb-2">Tipos Gramaticales:</h4>
      <div class="flex flex-wrap gap-3">
        {#each Object.entries(tipoColors) as [tipo, colors]}
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 rounded {colors.bg} border {colors.border}"></div>
            <span class="text-xs text-canvas-300 capitalize">{tipo}</span>
          </div>
        {/each}
      </div>
    </div>
  </div>

  <!-- Variaciones correctas (panel expandible) -->
  {#if mostrarVariaciones && variacionesCorrectas.length > 0}
    <div class="mb-6 p-6 bg-green-500/10 rounded-2xl border border-green-500/30">
      <h3 class="text-lg font-semibold text-green-300 mb-3">✓ Variaciones Correctas:</h3>
      <ul class="space-y-2">
        {#each variacionesCorrectas as variacion, index}
          <li class="flex items-start gap-2">
            <span class="text-green-400 mt-1">{index + 1}.</span>
            <p class="text-white font-medium">"{variacion}"</p>
          </li>
        {/each}
      </ul>
      <p class="text-sm text-green-400 mt-4 italic">
        💡 Todas estas construcciones son gramaticalmente válidas
      </p>
    </div>
  {/if}

  <!-- Navegación -->
  {#if showNavigation}
    <div class="flex items-center justify-between pt-6 mt-8 border-t border-canvas-800">
      <button
        onclick={onPrevious}
        disabled={!onPrevious}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-canvas-800 text-canvas-300
          border border-canvas-700
          transition-all duration-300
          hover:bg-canvas-700 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        ← Anterior
      </button>

      <div class="text-center">
        <p class="text-xs text-canvas-500">
          {oracionConstruida.length} palabras en tu oración
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-purple-500 to-purple-600
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-purple-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        Siguiente →
      </button>
    </div>
  {/if}
</div>

<style>
  /* Estilo para drag cursor */
  [draggable="true"] {
    user-select: none;
  }
</style>
