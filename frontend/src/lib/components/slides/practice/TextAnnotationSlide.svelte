<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Análisis de Texto",
    texto = "",
    tipoLectura = "narrativa", // narrativa | argumentativa | expositiva | poética
    preguntasGuia = [],
    vocabularioDestacado = [],
    herramientasAnotacion = ["resaltar", "notas"],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let textoRef = $state(null);
  let selectedText = $state("");
  let selectionRange = $state(null);
  let annotations = $state([]); // { tipo, texto, posicion, nota }
  let highlightColor = $state("yellow"); // Color actual de resaltado
  let showSummaryPanel = $state(false);
  let activeNote = $state(null); // Para editar nota
  let noteText = $state("");

  // Colores de resaltado
  const highlightColors = {
    yellow: { name: "Idea Principal", bg: "bg-yellow-400/30", border: "border-yellow-500" },
    green: { name: "Evidencia", bg: "bg-green-400/30", border: "border-green-500" },
    blue: { name: "Vocabulario", bg: "bg-blue-400/30", border: "border-blue-500" },
    red: { name: "Preguntas", bg: "bg-red-400/30", border: "border-red-500" }
  };

  // Anotaciones agrupadas por tipo
  const annotationsByType = $derived.by(() => {
    const grouped = {};
    annotations.forEach(ann => {
      if (!grouped[ann.tipo]) grouped[ann.tipo] = [];
      grouped[ann.tipo].push(ann);
    });
    return grouped;
  });

  function handleTextSelection() {
    const selection = window.getSelection();
    const text = selection.toString().trim();

    if (text && textoRef && textoRef.contains(selection.anchorNode)) {
      selectedText = text;
      selectionRange = {
        start: selection.anchorOffset,
        end: selection.focusOffset,
        node: selection.anchorNode
      };
    } else {
      selectedText = "";
      selectionRange = null;
    }
  }

  function addHighlight(color) {
    if (!selectedText) return;

    const newAnnotation = {
      id: Date.now(),
      tipo: color,
      texto: selectedText,
      nota: null,
      timestamp: new Date().toISOString()
    };

    annotations = [...annotations, newAnnotation];
    selectedText = "";
    window.getSelection().removeAllRanges();
  }

  function addNote() {
    if (!selectedText) return;

    activeNote = {
      id: Date.now(),
      tipo: highlightColor,
      texto: selectedText
    };
    noteText = "";
  }

  function saveNote() {
    if (!activeNote || !noteText.trim()) return;

    const newAnnotation = {
      ...activeNote,
      nota: noteText.trim(),
      timestamp: new Date().toISOString()
    };

    annotations = [...annotations, newAnnotation];
    activeNote = null;
    noteText = "";
    selectedText = "";
    window.getSelection().removeAllRanges();
  }

  function deleteAnnotation(id) {
    annotations = annotations.filter(ann => ann.id !== id);
  }

  // Renderizar texto con highlights
  const textoConHighlights = $derived.by(() => {
    if (annotations.length === 0) return texto;

    // Para simplificar, mostrar texto original
    // En producción, usar HTML parsing para aplicar highlights
    return texto;
  });

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

    // Event listener para selección de texto
    document.addEventListener('mouseup', handleTextSelection);

    return () => {
      document.removeEventListener('mouseup', handleTextSelection);
    };
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
        {tipoLectura}
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>
  </div>

  <!-- Toolbar de anotación -->
  {#if herramientasAnotacion.includes("resaltar")}
    <div class="mb-6 p-4 bg-canvas-900/50 rounded-xl border border-canvas-700">
      <h3 class="text-sm font-semibold text-canvas-300 mb-3">Herramientas de Anotación</h3>

      <div class="flex flex-wrap items-center gap-3">
        <!-- Color selector -->
        <div class="flex items-center gap-2">
          {#each Object.entries(highlightColors) as [color, config]}
            <button
              onclick={() => { highlightColor = color; addHighlight(color); }}
              disabled={!selectedText}
              class="
                px-3 py-2 rounded-lg text-xs font-semibold
                transition-all duration-300
                {highlightColor === color ? `${config.bg} ${config.border} border-2` : 'bg-canvas-800 border border-canvas-700 text-canvas-400'}
                hover:scale-105 disabled:opacity-30 disabled:cursor-not-allowed
              "
              title={config.name}
            >
              {config.name}
            </button>
          {/each}
        </div>

        <!-- Nota button -->
        {#if herramientasAnotacion.includes("notas")}
          <div class="h-6 w-px bg-canvas-700"></div>
          <button
            onclick={addNote}
            disabled={!selectedText}
            class="
              px-4 py-2 rounded-lg text-xs font-semibold
              bg-purple-500/20 text-purple-400 border border-purple-500
              transition-all duration-300
              hover:bg-purple-500/30 hover:scale-105
              disabled:opacity-30 disabled:cursor-not-allowed
            "
          >
            📝 Agregar Nota
          </button>
        {/if}

        <!-- Summary toggle -->
        <div class="ml-auto">
          <button
            onclick={() => showSummaryPanel = !showSummaryPanel}
            class="
              px-4 py-2 rounded-lg text-xs font-semibold
              bg-canvas-800 text-canvas-300 border border-canvas-700
              transition-all duration-300
              hover:bg-canvas-700 hover:scale-105
            "
          >
            {showSummaryPanel ? '✕ Cerrar' : '📋'} Resumen ({annotations.length})
          </button>
        </div>
      </div>

      {#if selectedText}
        <div class="mt-3 p-2 bg-slate-950 rounded-lg border border-canvas-700">
          <p class="text-xs text-canvas-400">
            Texto seleccionado: <span class="text-white font-semibold">"{selectedText.substring(0, 50)}{selectedText.length > 50 ? '...' : ''}"</span>
          </p>
        </div>
      {/if}
    </div>
  {/if}

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Texto principal -->
    <div class="lg:col-span-2">
      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
        <div
          bind:this={textoRef}
          class="prose prose-invert max-w-none text-canvas-200 leading-relaxed select-text"
        >
          {textoConHighlights}
        </div>
      </div>

      <!-- Preguntas guía -->
      {#if preguntasGuia.length > 0}
        <div class="mt-6 p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
          <h3 class="text-lg font-semibold text-purple-300 mb-3">📚 Preguntas Guía</h3>
          <ul class="space-y-2">
            {#each preguntasGuia as pregunta}
              <li class="flex items-start gap-2">
                <span class="text-purple-400 mt-1">•</span>
                <p class="text-canvas-300">{pregunta}</p>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>

    <!-- Panel de resumen (sidebar) -->
    {#if showSummaryPanel}
      <div class="lg:col-span-1">
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700 sticky top-4">
          <h3 class="text-lg font-semibold text-white mb-4">📋 Mis Anotaciones</h3>

          {#if annotations.length === 0}
            <p class="text-sm text-canvas-400 italic">
              Aún no has hecho anotaciones. Selecciona texto y usa las herramientas arriba.
            </p>
          {:else}
            <div class="space-y-4 max-h-96 overflow-y-auto">
              {#each Object.entries(annotationsByType) as [tipo, anns]}
                <div>
                  <h4 class="text-xs font-bold uppercase text-canvas-400 mb-2">
                    {highlightColors[tipo].name} ({anns.length})
                  </h4>
                  <div class="space-y-2">
                    {#each anns as ann}
                      <div class="p-3 bg-slate-950 rounded-lg border {highlightColors[tipo].border}">
                        <p class="text-sm text-white font-medium mb-1">
                          "{ann.texto.substring(0, 40)}{ann.texto.length > 40 ? '...' : ''}"
                        </p>
                        {#if ann.nota}
                          <p class="text-xs text-canvas-400 italic mt-2">
                            📝 {ann.nota}
                          </p>
                        {/if}
                        <button
                          onclick={() => deleteAnnotation(ann.id)}
                          class="text-xs text-red-400 hover:text-red-300 mt-2"
                        >
                          Eliminar
                        </button>
                      </div>
                    {/each}
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>

  <!-- Modal para agregar nota -->
  {#if activeNote}
    <div class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4">
      <div class="bg-slate-950 rounded-2xl border border-canvas-700 p-6 max-w-lg w-full">
        <h3 class="text-xl font-bold text-white mb-4">📝 Agregar Nota</h3>

        <div class="mb-4 p-3 bg-canvas-900 rounded-lg border border-canvas-700">
          <p class="text-sm text-canvas-400 mb-1">Texto seleccionado:</p>
          <p class="text-white font-medium">"{activeNote.texto}"</p>
        </div>

        <textarea
          bind:value={noteText}
          placeholder="Escribe tu nota o pregunta aquí..."
          class="
            w-full px-4 py-3 bg-canvas-900 border-2 border-canvas-700 rounded-xl
            text-white placeholder-canvas-500
            focus:outline-none focus:border-purple-500
            transition-colors duration-300
            min-h-32
          "
        ></textarea>

        <div class="flex items-center gap-3 mt-4">
          <button
            onclick={saveNote}
            disabled={!noteText.trim()}
            class="
              flex-1 px-6 py-3 rounded-xl font-semibold
              bg-gradient-to-r from-purple-500 to-purple-600
              text-white
              transition-all duration-300
              hover:shadow-lg hover:shadow-purple-500/50 hover:scale-105
              disabled:opacity-30 disabled:cursor-not-allowed
            "
          >
            Guardar Nota
          </button>
          <button
            onclick={() => { activeNote = null; noteText = ""; }}
            class="
              px-6 py-3 rounded-xl font-semibold
              bg-canvas-800 text-canvas-300
              border border-canvas-700
              transition-all duration-300
              hover:bg-canvas-700 hover:scale-105
            "
          >
            Cancelar
          </button>
        </div>
      </div>
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
          {annotations.length} anotaciones realizadas
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
  .prose {
    font-size: 1.125rem;
    line-height: 1.75;
  }
</style>
