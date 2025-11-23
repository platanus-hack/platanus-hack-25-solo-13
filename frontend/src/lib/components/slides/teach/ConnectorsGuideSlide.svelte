<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Guía de Conectores Textuales",
    categorias = [],
    importancia = "",
    consejos = [],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let categoriaSeleccionada = $state(0);

  // Colores por categoría
  const categoriaColors = {
    causal: { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300", badge: "bg-blue-500" },
    adversativo: { bg: "bg-red-500/20", border: "border-red-500", text: "text-red-300", badge: "bg-red-500" },
    consecutivo: { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300", badge: "bg-green-500" },
    aditivo: { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300", badge: "bg-purple-500" },
    temporal: { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300", badge: "bg-orange-500" }
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
  class="w-full max-w-7xl mx-auto p-8 bg-slate-950 rounded-2xl border border-canvas-800 shadow-2xl"
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
        Conectores
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    {#if importancia}
      <p class="text-canvas-300">
        {importancia}
      </p>
    {/if}
  </div>

  <!-- Tabs de categorías -->
  <div class="mb-6 flex flex-wrap gap-2">
    {#each categorias as categoria, index}
      {@const colors = categoriaColors[categoria.tipo] || categoriaColors.causal}
      <button
        onclick={() => categoriaSeleccionada = index}
        class="
          px-5 py-3 rounded-xl font-semibold
          transition-all duration-300
          {categoriaSeleccionada === index
            ? `${colors.badge} text-white scale-105 shadow-lg`
            : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
        "
      >
        {categoria.nombre}
      </button>
    {/each}
  </div>

  <!-- Contenido de categoría seleccionada -->
  {#if categorias[categoriaSeleccionada]}
    {@const categoria = categorias[categoriaSeleccionada]}
    {@const colors = categoriaColors[categoria.tipo] || categoriaColors.causal}

    <div class="space-y-6">
      <!-- Definición y función -->
      <div class="p-6 {colors.bg} rounded-2xl border-2 {colors.border}">
        <h3 class="text-2xl font-bold {colors.text} mb-3">
          {categoria.nombre}
        </h3>
        <p class="text-lg text-canvas-200 mb-4">
          {categoria.definicion}
        </p>
        <div class="p-4 bg-slate-950 rounded-lg border {colors.border}">
          <p class="text-sm text-canvas-300">
            <span class="font-semibold {colors.text}">Función:</span> {categoria.funcion}
          </p>
        </div>
      </div>

      <!-- Lista de conectores con ejemplos -->
      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
        <h4 class="text-lg font-semibold text-white mb-4">
          Conectores de esta categoría
        </h4>

        <div class="space-y-4">
          {#each categoria.conectores as conector}
            <div class="p-4 bg-slate-950 rounded-xl border {colors.border}">
              <div class="flex items-center gap-3 mb-3">
                <span class="px-4 py-2 {colors.badge} rounded-lg text-white font-bold">
                  {conector.palabra}
                </span>
                {#if conector.nivel}
                  <span class="text-xs px-2 py-1 bg-canvas-800 rounded text-canvas-400">
                    Formalidad: {conector.nivel}
                  </span>
                {/if}
              </div>

              <!-- Ejemplos de uso -->
              {#if conector.ejemplos && conector.ejemplos.length > 0}
                <div class="space-y-2">
                  {#each conector.ejemplos as ejemplo}
                    <div class="pl-4 border-l-2 {colors.border}">
                      <p class="text-canvas-200">
                        {ejemplo}
                      </p>
                    </div>
                  {/each}
                </div>
              {/if}

              <!-- Nota de uso (opcional) -->
              {#if conector.nota}
                <p class="text-xs text-canvas-400 italic mt-3">
                  💡 {conector.nota}
                </p>
              {/if}
            </div>
          {/each}
        </div>
      </div>

      <!-- Comparación con otras categorías -->
      {#if categoria.comparacion}
        <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
          <h4 class="text-sm font-semibold uppercase text-purple-400 mb-3">
            📊 Comparación
          </h4>
          <p class="text-canvas-200">
            {categoria.comparacion}
          </p>
        </div>
      {/if}

      <!-- Ejemplo de párrafo -->
      {#if categoria.ejemploParrafo}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-sm font-semibold uppercase text-canvas-400 mb-3">
            📝 Ejemplo de Párrafo
          </h4>
          <p class="text-lg text-canvas-200 leading-relaxed">
            {categoria.ejemploParrafo}
          </p>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Consejos generales -->
  {#if consejos.length > 0}
    <div class="mt-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        💡 Consejos para Usar Conectores Efectivamente
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each consejos as consejo}
          <div class="flex items-start gap-3 p-3 bg-slate-950 rounded-lg">
            <span class="text-xl">✓</span>
            <p class="text-sm text-canvas-300">
              {consejo}
            </p>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Tabla resumen (todas las categorías) -->
  <div class="mt-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
    <h3 class="text-lg font-semibold text-white mb-4">
      📋 Resumen de Todas las Categorías
    </h3>
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-canvas-700">
            <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Categoría</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Función</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Ejemplos</th>
          </tr>
        </thead>
        <tbody>
          {#each categorias as cat}
            {@const colors = categoriaColors[cat.tipo] || categoriaColors.causal}
            <tr class="border-b border-canvas-800 hover:bg-canvas-800/30">
              <td class="py-3 px-4">
                <span class="px-3 py-1 {colors.badge} rounded text-white text-sm font-semibold">
                  {cat.nombre}
                </span>
              </td>
              <td class="py-3 px-4 text-sm text-canvas-300">
                {cat.funcion}
              </td>
              <td class="py-3 px-4 text-sm text-canvas-400">
                {cat.conectores.slice(0, 3).map(c => c.palabra).join(', ')}...
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>

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
          {categorias.length} categorías • {categorias.reduce((acc, cat) => acc + cat.conectores.length, 0)} conectores
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
