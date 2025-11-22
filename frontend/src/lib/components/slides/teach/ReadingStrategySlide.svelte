<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Estrategias de Comprensión Lectora",
    estrategias = [],
    ejemploTexto = "",
    tipsAdicionales = [],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let estrategiaSeleccionada = $state(0);

  // Colores por estrategia
  const estrategiaColors = {
    0: { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300" },
    1: { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300" },
    2: { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300" },
    3: { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300" },
    4: { bg: "bg-pink-500/20", border: "border-pink-500", text: "text-pink-300" }
  };

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
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      Aprende técnicas efectivas para comprender mejor cualquier texto
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Lista de estrategias (sidebar) -->
    <div class="lg:col-span-1 space-y-3">
      <h3 class="text-sm font-semibold uppercase text-canvas-400 mb-4">
        Estrategias Clave
      </h3>

      {#each estrategias as estrategia, index}
        {@const colors = estrategiaColors[index % 5]}
        <button
          onclick={() => estrategiaSeleccionada = index}
          class="
            w-full text-left p-4 rounded-xl transition-all duration-300
            {estrategiaSeleccionada === index
              ? `${colors.bg} ${colors.border} border-2 scale-105`
              : 'bg-canvas-900/50 border border-canvas-700 hover:bg-canvas-800 hover:scale-102'}
          "
        >
          <div class="flex items-start gap-3">
            <span class="text-2xl">{estrategia.icono}</span>
            <div class="flex-1">
              <h4 class="text-lg font-bold {estrategiaSeleccionada === index ? colors.text : 'text-white'} mb-1">
                {estrategia.nombre}
              </h4>
              <p class="text-xs {estrategiaSeleccionada === index ? 'text-canvas-200' : 'text-canvas-500'}">
                {estrategia.resumen}
              </p>
            </div>
          </div>
        </button>
      {/each}
    </div>

    <!-- Contenido de estrategia seleccionada -->
    <div class="lg:col-span-2 space-y-6">
      {#if estrategias[estrategiaSeleccionada]}
        {@const estrategia = estrategias[estrategiaSeleccionada]}
        {@const colors = estrategiaColors[estrategiaSeleccionada % 5]}

        <!-- Explicación detallada -->
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <div class="flex items-center gap-3 mb-4">
            <span class="text-4xl">{estrategia.icono}</span>
            <h3 class="text-2xl font-bold {colors.text}">
              {estrategia.nombre}
            </h3>
          </div>

          <p class="text-lg text-canvas-200 mb-6 leading-relaxed">
            {estrategia.explicacion}
          </p>

          <!-- Pasos específicos -->
          {#if estrategia.pasos && estrategia.pasos.length > 0}
            <div class="mb-6">
              <h4 class="text-sm font-semibold uppercase text-canvas-400 mb-3">
                ¿Cómo aplicarla?
              </h4>
              <div class="space-y-3">
                {#each estrategia.pasos as paso, idx}
                  <div class="flex items-start gap-3 p-3 bg-canvas-950 rounded-lg border {colors.border}">
                    <div class="w-8 h-8 rounded-full {colors.bg} {colors.border} border-2 flex items-center justify-center flex-shrink-0">
                      <span class="text-sm font-bold {colors.text}">
                        {idx + 1}
                      </span>
                    </div>
                    <p class="text-canvas-200 pt-1">
                      {paso}
                    </p>
                  </div>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Ejemplo visual -->
          {#if estrategia.ejemplo}
            <div class="p-4 {colors.bg} rounded-xl border-2 {colors.border}">
              <p class="text-xs font-semibold uppercase {colors.text} mb-2">
                💡 Ejemplo
              </p>
              <p class="text-canvas-200 italic">
                {estrategia.ejemplo}
              </p>
            </div>
          {/if}
        </div>

        <!-- Cuándo usar esta estrategia -->
        {#if estrategia.cuandoUsar}
          <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
            <h4 class="text-sm font-semibold uppercase text-purple-400 mb-3">
              📖 Cuándo usar esta estrategia
            </h4>
            <ul class="space-y-2">
              {#each estrategia.cuandoUsar as situacion}
                <li class="flex items-start gap-2 text-canvas-200">
                  <span class="text-purple-400">•</span>
                  <span>{situacion}</span>
                </li>
              {/each}
            </ul>
          </div>
        {/if}
      {/if}
    </div>
  </div>

  <!-- Tips adicionales -->
  {#if tipsAdicionales.length > 0}
    <div class="mt-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        💎 Tips Adicionales para Mejorar tu Comprensión
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each tipsAdicionales as tip}
          <div class="flex items-start gap-3 p-3 bg-canvas-950 rounded-lg">
            <span class="text-xl">✓</span>
            <p class="text-sm text-canvas-300">
              {tip}
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
          Estrategia {estrategiaSeleccionada + 1} de {estrategias.length}
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
