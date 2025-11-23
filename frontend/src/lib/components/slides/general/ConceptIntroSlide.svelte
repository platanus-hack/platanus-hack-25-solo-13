<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    concepto = "Concepto",
    definicionSimple = "Definición en términos simples",
    definicionTecnica = "Definición técnica con detalles científicos",
    imagen = null,
    terminosClave = [],
    colorTema = "blue",
    materia = "general",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let versionTecnica = $state(false);
  let containerRef = $state(null);
  let hoveredTermino = $state(null);

  // Colores por materia
  const materiaColors = {
    matemáticas: 'cyan',
    lenguaje: 'purple',
    historia: 'amber',
    física: 'blue',
    química: 'green',
    biología: 'emerald',
    general: colorTema
  };

  const currentColor = materiaColors[materia] || colorTema;

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

  function toggleVersion() {
    versionTecnica = !versionTecnica;

    // Animar transición de definición
    const definitionEl = document.querySelector('.definition-text');
    if (definitionEl) {
      gsap.fromTo(definitionEl,
        { opacity: 0, x: -20 },
        { opacity: 1, x: 0, duration: 0.4, ease: 'power2.out' }
      );
    }
  }

  // Resaltar términos clave en el texto
  function highlightTerminos(texto) {
    if (!terminosClave || terminosClave.length === 0) return texto;

    let resultado = texto;
    terminosClave.forEach(termino => {
      const regex = new RegExp(`\\b(${termino.palabra})\\b`, 'gi');
      resultado = resultado.replace(regex, `<mark class="termino-clave" data-termino="${termino.palabra}">$1</mark>`);
    });
    return resultado;
  }

  const definicionMarcada = $derived(
    highlightTerminos(versionTecnica ? definicionTecnica : definicionSimple)
  );
</script>

<div
  bind:this={containerRef}
  class="w-full max-w-5xl mx-auto p-8 bg-slate-950 rounded-2xl border border-slate-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-{currentColor}-500/20 text-{currentColor}-400 border border-{currentColor}-500">
        {materia}
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <!-- Concepto principal -->
    <h2 class="text-4xl font-bold text-white mb-2">
      {concepto}
    </h2>

    <!-- Toggle versión -->
    <button
      onclick={toggleVersion}
      class="
        mt-4 px-4 py-2 rounded-xl text-sm font-semibold
        transition-all duration-300
        {versionTecnica ? 'bg-purple-500/20 text-purple-400 border border-purple-500' : 'bg-canvas-800 text-slate-400 border border-slate-700'}
        hover:scale-105
      "
    >
      {versionTecnica ? '👁️ Ver versión simple' : '🔬 Ver versión técnica'}
    </button>
  </div>

  <!-- Contenido principal -->
  <div class="grid grid-cols-1 {imagen ? 'lg:grid-cols-2' : ''} gap-8 mb-8">
    <!-- Definición -->
    <div>
      <div class="flex items-center gap-2 mb-4">
        <span class="text-2xl">
          {versionTecnica ? '📖' : '💡'}
        </span>
        <h3 class="text-lg font-semibold text-slate-300">
          {versionTecnica ? 'Definición Técnica' : 'Definición Simple'}
        </h3>
      </div>

      <div
        class="definition-text p-6 bg-canvas-900/50 rounded-2xl border border-slate-700 text-slate-200 leading-relaxed"
      >
        {@html definicionMarcada}
      </div>

      <!-- Leyenda de términos clave -->
      {#if terminosClave.length > 0}
        <div class="mt-4 p-4 bg-{currentColor}-500/10 rounded-xl border border-{currentColor}-500/30">
          <p class="text-xs font-semibold text-{currentColor}-400 mb-2 uppercase">
            💡 Términos Clave (pasa el cursor sobre ellos):
          </p>
          <div class="flex flex-wrap gap-2">
            {#each terminosClave as termino}
              <button
                onmouseenter={() => hoveredTermino = termino.palabra}
                onmouseleave={() => hoveredTermino = null}
                class="px-2 py-1 bg-canvas-800 rounded-lg text-xs text-slate-300 hover:bg-{currentColor}-500/20 hover:text-{currentColor}-300 transition-colors relative"
              >
                {termino.palabra}

                {#if hoveredTermino === termino.palabra}
                  <div class="absolute bottom-full left-0 mb-2 w-64 p-3 bg-canvas-800 rounded-xl border border-{currentColor}-500 shadow-lg z-10">
                    <p class="text-xs text-slate-200">{termino.tooltip}</p>
                  </div>
                {/if}
              </button>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Imagen de apoyo -->
    {#if imagen}
      <div class="flex items-center justify-center">
        <div class="relative">
          <img
            src={imagen}
            alt={concepto}
            class="w-full max-w-md rounded-2xl border-2 border-slate-700 shadow-xl"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-{currentColor}-500/10 to-transparent rounded-2xl pointer-events-none"></div>
        </div>
      </div>
    {/if}
  </div>

  <!-- Navegación -->
  {#if showNavigation}
    <div class="flex items-center justify-between pt-6 border-t border-slate-800">
      <button
        onclick={onPrevious}
        disabled={!onPrevious}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-canvas-800 text-slate-300
          border border-slate-700
          transition-all duration-300
          hover:bg-slate-700 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        ← Anterior
      </button>

      <div class="text-center">
        <p class="text-xs text-slate-500">
          Presiona "Siguiente" para continuar
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-{currentColor}-500 to-{currentColor}-600
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-{currentColor}-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        Siguiente →
      </button>
    </div>
  {/if}
</div>

<style>
  :global(.termino-clave) {
    background-color: rgb(234 179 8 / 0.2);
    color: #fde047;
    padding-left: 0.25rem;
    padding-right: 0.25rem;
    border-radius: 0.25rem;
    border-bottom-width: 2px;
    border-color: rgb(234 179 8 / 0.5);
    cursor: help;
  }
</style>
