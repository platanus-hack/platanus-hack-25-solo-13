<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Estrategias para Aprender Vocabulario",
    palabra = "",
    estrategias = [],
    familiasPalabras = [],
    consejos = [],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let estrategiaActiva = $state(0);

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
  class="w-full max-w-7xl mx-auto p-8 bg-canvas-950 rounded-2xl border border-canvas-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-blue-500/20 text-blue-400 border border-blue-500">
        📚 Enseñanza
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-purple-500/20 text-purple-400 border border-purple-500">
        {materia}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-medium bg-canvas-800 text-canvas-300">
        Vocabulario
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      Aprende técnicas efectivas para expandir tu vocabulario académico
    </p>
  </div>

  <!-- Palabra ejemplo (si se proporciona) -->
  {#if palabra}
    <div class="mb-8 p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
      <p class="text-sm text-purple-400 mb-2">Palabra de ejemplo:</p>
      <h3 class="text-3xl font-bold text-white">
        {palabra}
      </h3>
    </div>
  {/if}

  <!-- Estrategias de aprendizaje -->
  <div class="mb-8">
    <h3 class="text-lg font-semibold text-white mb-4">
      🎯 Estrategias de Aprendizaje
    </h3>

    <div class="space-y-4">
      {#each estrategias as estrategia, index}
        {@const colors = [
          { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300" },
          { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300" },
          { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300" },
          { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300" }
        ][index % 4]}

        <button
          onclick={() => estrategiaActiva = index}
          class="
            w-full text-left p-6 rounded-2xl border-2 transition-all duration-300
            {estrategiaActiva === index
              ? `${colors.bg} ${colors.border} scale-102`
              : 'bg-canvas-900/50 border-canvas-700 hover:bg-canvas-800'}
          "
        >
          <div class="flex items-start gap-4">
            <span class="text-3xl">{estrategia.icono}</span>
            <div class="flex-1">
              <h4 class="text-xl font-bold {estrategiaActiva === index ? colors.text : 'text-white'} mb-2">
                {estrategia.nombre}
              </h4>
              <p class="text-canvas-200 mb-4">
                {estrategia.descripcion}
              </p>

              {#if estrategiaActiva === index}
                <!-- Pasos detallados -->
                {#if estrategia.pasos && estrategia.pasos.length > 0}
                  <div class="mt-4 space-y-2">
                    <p class="text-xs font-semibold uppercase text-canvas-400">Cómo aplicarla:</p>
                    {#each estrategia.pasos as paso, idx}
                      <div class="flex items-start gap-2 p-3 bg-canvas-950 rounded-lg">
                        <span class="text-sm font-bold {colors.text}">{idx + 1}.</span>
                        <p class="text-sm text-canvas-200">{paso}</p>
                      </div>
                    {/each}
                  </div>
                {/if}

                <!-- Ejemplo aplicado -->
                {#if estrategia.ejemplo}
                  <div class="mt-4 p-4 bg-canvas-950 rounded-lg border {colors.border}">
                    <p class="text-xs font-semibold uppercase text-canvas-400 mb-2">
                      💡 Ejemplo con "{palabra || 'la palabra'}"
                    </p>
                    <p class="text-canvas-200">
                      {estrategia.ejemplo}
                    </p>
                  </div>
                {/if}
              {/if}
            </div>
          </div>
        </button>
      {/each}
    </div>
  </div>

  <!-- Análisis etimológico (diagrama) -->
  {#if familiasPalabras.length > 0}
    <div class="mb-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        🌳 Familias de Palabras
      </h3>
      <p class="text-sm text-canvas-400 mb-6">
        Las palabras relacionadas ayudan a recordar y comprender mejor el vocabulario
      </p>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each familiasPalabras as familia}
          <div class="p-4 bg-canvas-950 rounded-xl border border-canvas-700">
            <h4 class="text-lg font-bold text-purple-300 mb-3">
              {familia.raiz}
            </h4>
            <p class="text-xs text-canvas-400 mb-3">
              Significado: {familia.significado}
            </p>
            <div class="space-y-2">
              {#each familia.palabras as palabra}
                <div class="flex items-center gap-2 text-sm">
                  <span class="w-2 h-2 rounded-full bg-purple-500"></span>
                  <span class="text-white font-medium">{palabra.palabra}</span>
                  <span class="text-canvas-400">→ {palabra.definicion}</span>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Consejos prácticos -->
  {#if consejos.length > 0}
    <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        💎 Consejos Prácticos
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        {#each consejos as consejo}
          <div class="flex items-start gap-3 p-3 bg-canvas-950 rounded-lg">
            <span class="text-green-400">✓</span>
            <p class="text-sm text-canvas-300">
              {consejo}
            </p>
          </div>
        {/each}
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
          {estrategias.length} estrategias
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-blue-500 to-purple-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-blue-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        Siguiente →
      </button>
    </div>
  {/if}
</div>
