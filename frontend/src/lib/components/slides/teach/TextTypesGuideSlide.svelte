<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Guía de Tipos de Texto",
    tipos = [],
    comparacionTabla = null,
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let tipoSeleccionado = $state(0);

  // Colores por tipo de texto
  const tipoColors = {
    narrativo: { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300", badge: "bg-blue-500" },
    argumentativo: { bg: "bg-red-500/20", border: "border-red-500", text: "text-red-300", badge: "bg-red-500" },
    expositivo: { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300", badge: "bg-green-500" },
    instructivo: { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300", badge: "bg-orange-500" },
    descriptivo: { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300", badge: "bg-purple-500" }
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
      <span class="px-3 py-1 rounded-full text-xs font-medium bg-canvas-800 text-canvas-300">
        Tipos de Texto
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      Comprende las características y estructuras de cada tipo de texto
    </p>
  </div>

  <!-- Tabs de tipos de texto -->
  <div class="mb-6 flex flex-wrap gap-2">
    {#each tipos as tipo, index}
      {@const colors = tipoColors[tipo.tipo] || tipoColors.expositivo}
      <button
        onclick={() => tipoSeleccionado = index}
        class="
          px-5 py-3 rounded-xl font-semibold
          transition-all duration-300
          {tipoSeleccionado === index
            ? `${colors.badge} text-white scale-105 shadow-lg`
            : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
        "
      >
        {tipo.nombre}
      </button>
    {/each}
  </div>

  <!-- Contenido del tipo seleccionado -->
  {#if tipos[tipoSeleccionado]}
    {@const tipo = tipos[tipoSeleccionado]}
    {@const colors = tipoColors[tipo.tipo] || tipoColors.expositivo}

    <div class="space-y-6">
      <!-- Definición y propósito -->
      <div class="p-6 {colors.bg} rounded-2xl border-2 {colors.border}">
        <h3 class="text-2xl font-bold {colors.text} mb-3">
          {tipo.nombre}
        </h3>
        <p class="text-lg text-canvas-200 mb-4">
          {tipo.definicion}
        </p>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="p-4 bg-canvas-950 rounded-lg border {colors.border}">
            <p class="text-xs font-semibold uppercase text-canvas-400 mb-2">Propósito</p>
            <p class="text-canvas-200">{tipo.proposito}</p>
          </div>
          <div class="p-4 bg-canvas-950 rounded-lg border {colors.border}">
            <p class="text-xs font-semibold uppercase text-canvas-400 mb-2">Contexto de uso</p>
            <p class="text-canvas-200">{tipo.contexto}</p>
          </div>
        </div>
      </div>

      <!-- Estructura típica -->
      {#if tipo.estructura && tipo.estructura.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-lg font-semibold text-white mb-4">
            📐 Estructura Típica
          </h4>
          <div class="space-y-3">
            {#each tipo.estructura as seccion, idx}
              <div class="flex items-start gap-4">
                <div class="w-10 h-10 rounded-full {colors.badge} flex items-center justify-center flex-shrink-0">
                  <span class="text-white font-bold">{idx + 1}</span>
                </div>
                <div class="flex-1 p-4 bg-canvas-950 rounded-lg border {colors.border}">
                  <h5 class="text-lg font-bold {colors.text} mb-2">
                    {seccion.nombre}
                  </h5>
                  <p class="text-canvas-200 mb-3">
                    {seccion.descripcion}
                  </p>
                  {#if seccion.ejemplo}
                    <p class="text-sm text-canvas-400 italic">
                      Ejemplo: "{seccion.ejemplo}"
                    </p>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Características del lenguaje -->
      {#if tipo.caracteristicas && tipo.caracteristicas.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-lg font-semibold text-white mb-4">
            ✍️ Características del Lenguaje
          </h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            {#each tipo.caracteristicas as caracteristica}
              <div class="flex items-start gap-2 p-3 bg-canvas-950 rounded-lg">
                <span class="{colors.text}">•</span>
                <p class="text-canvas-200">{caracteristica}</p>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Ejemplo de texto -->
      {#if tipo.ejemploTexto}
        <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
          <h4 class="text-sm font-semibold uppercase text-purple-400 mb-3">
            📝 Ejemplo de Texto {tipo.nombre}
          </h4>
          <p class="text-lg text-canvas-200 leading-relaxed italic">
            {tipo.ejemploTexto}
          </p>
        </div>
      {/if}

      <!-- Conectores típicos -->
      {#if tipo.conectoresTipicos && tipo.conectoresTipicos.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-lg font-semibold text-white mb-4">
            🔗 Conectores Típicos
          </h4>
          <div class="flex flex-wrap gap-2">
            {#each tipo.conectoresTipicos as conector}
              <span class="px-3 py-2 {colors.badge} rounded-lg text-white font-medium">
                {conector}
              </span>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Tabla comparativa (opcional) -->
  {#if comparacionTabla}
    <div class="mt-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        📊 Comparación de Tipos de Texto
      </h3>
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-canvas-700">
              <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Tipo</th>
              <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Propósito</th>
              <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Estructura</th>
              <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Lenguaje</th>
            </tr>
          </thead>
          <tbody>
            {#each tipos as tipo}
              {@const colors = tipoColors[tipo.tipo] || tipoColors.expositivo}
              <tr class="border-b border-canvas-800 hover:bg-canvas-800/30">
                <td class="py-3 px-4">
                  <span class="px-3 py-1 {colors.badge} rounded text-white text-sm font-semibold">
                    {tipo.nombre}
                  </span>
                </td>
                <td class="py-3 px-4 text-sm text-canvas-300">
                  {tipo.proposito}
                </td>
                <td class="py-3 px-4 text-sm text-canvas-400">
                  {tipo.estructura.map(s => s.nombre).join(' → ')}
                </td>
                <td class="py-3 px-4 text-sm text-canvas-400">
                  {tipo.caracteristicas[0]}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
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
          {tipos.length} tipos de texto
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
