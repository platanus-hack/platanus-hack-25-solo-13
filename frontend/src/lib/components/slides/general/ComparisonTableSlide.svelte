<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Comparación de Conceptos",
    items = [
      { nombre: "Concepto A", color: "cyan" },
      { nombre: "Concepto B", color: "purple" }
    ],
    filas = [],
    materia = "general",
    mostrarFiltros = true,
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let expandedRows = $state(new Set());
  let filtroActivo = $state("todos"); // "todos" | "similitudes" | "diferencias"

  // Filtrar filas según filtro activo
  const filasFiltradas = $derived.by(() => {
    if (filtroActivo === "todos") return filas;
    if (filtroActivo === "similitudes") {
      return filas.filter(f => f.tipo === "similitud");
    }
    if (filtroActivo === "diferencias") {
      return filas.filter(f => f.tipo === "diferencia");
    }
    return filas;
  });

  // Toggle expand/collapse de fila
  function toggleRow(filaId) {
    const newSet = new Set(expandedRows);
    if (newSet.has(filaId)) {
      newSet.delete(filaId);
    } else {
      newSet.add(filaId);
    }
    expandedRows = newSet;
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
  class="w-full max-w-6xl mx-auto p-8 bg-canvas-950 rounded-2xl border border-slate-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-blue-500/20 text-blue-400 border border-blue-500">
        {materia}
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-4">
      {titulo}
    </h2>

    <!-- Filtros -->
    {#if mostrarFiltros}
      <div class="flex gap-2">
        <button
          onclick={() => filtroActivo = "todos"}
          class="
            px-4 py-2 rounded-xl text-sm font-semibold transition-all duration-300
            {filtroActivo === 'todos' ? 'bg-blue-500 text-white' : 'bg-canvas-800 text-slate-400 hover:bg-slate-700'}
          "
        >
          Todos ({filas.length})
        </button>
        <button
          onclick={() => filtroActivo = "similitudes"}
          class="
            px-4 py-2 rounded-xl text-sm font-semibold transition-all duration-300
            {filtroActivo === 'similitudes' ? 'bg-green-500 text-white' : 'bg-canvas-800 text-slate-400 hover:bg-slate-700'}
          "
        >
          Similitudes ({filas.filter(f => f.tipo === 'similitud').length})
        </button>
        <button
          onclick={() => filtroActivo = "diferencias"}
          class="
            px-4 py-2 rounded-xl text-sm font-semibold transition-all duration-300
            {filtroActivo === 'diferencias' ? 'bg-orange-500 text-white' : 'bg-canvas-800 text-slate-400 hover:bg-slate-700'}
          "
        >
          Diferencias ({filas.filter(f => f.tipo === 'diferencia').length})
        </button>
      </div>
    {/if}
  </div>

  <!-- Tabla comparativa -->
  <div class="mb-8 overflow-x-auto">
    <table class="w-full border-collapse">
      <!-- Header de columnas -->
      <thead>
        <tr class="border-b-2 border-slate-700">
          <th class="p-4 text-left text-sm font-semibold text-slate-400 uppercase w-1/3">
            Característica
          </th>
          {#each items as item}
            <th class="p-4 text-center text-sm font-semibold uppercase w-1/3 text-{item.color}-400">
              {item.nombre}
            </th>
          {/each}
        </tr>
      </thead>

      <!-- Filas de datos -->
      <tbody>
        {#each filasFiltradas as fila, index (fila.caracteristica)}
          {@const isExpanded = expandedRows.has(fila.caracteristica)}
          {@const tieneDetalles = fila.detalles && fila.detalles.trim().length > 0}

          <!-- Fila principal -->
          <tr
            class="
              border-b border-slate-800 transition-colors duration-200
              {tieneDetalles ? 'cursor-pointer hover:bg-canvas-900/50' : ''}
              {fila.tipo === 'similitud' ? 'bg-green-500/5' : fila.tipo === 'diferencia' ? 'bg-orange-500/5' : ''}
            "
            onclick={() => tieneDetalles && toggleRow(fila.caracteristica)}
          >
            <td class="p-4 text-slate-200 font-medium">
              <div class="flex items-center gap-2">
                <span class="text-lg">
                  {fila.tipo === 'similitud' ? '✓' : fila.tipo === 'diferencia' ? '⚡' : '•'}
                </span>
                <span>{fila.caracteristica}</span>
                {#if tieneDetalles}
                  <span class="text-xs text-slate-500">
                    {isExpanded ? '▼' : '▶'}
                  </span>
                {/if}
              </div>
            </td>
            {#each fila.valores as valor, i}
              <td class="p-4 text-center text-slate-300">
                <div class="inline-block px-3 py-1 bg-canvas-800 rounded-lg">
                  {valor}
                </div>
              </td>
            {/each}
          </tr>

          <!-- Fila expandida con detalles -->
          {#if isExpanded && tieneDetalles}
            <tr class="border-b border-slate-800 bg-canvas-900/30">
              <td colspan={items.length + 1} class="p-4">
                <div class="pl-8 text-sm text-slate-400">
                  <span class="text-blue-400 font-semibold">Detalles:</span> {fila.detalles}
                </div>
              </td>
            </tr>
          {/if}
        {/each}
      </tbody>
    </table>

    {#if filasFiltradas.length === 0}
      <div class="text-center py-12 text-slate-500">
        <p class="text-lg">No hay {filtroActivo} para mostrar</p>
      </div>
    {/if}
  </div>

  <!-- Leyenda -->
  <div class="mb-8 flex gap-4 text-sm text-slate-400">
    <div class="flex items-center gap-2">
      <span class="text-green-400">✓</span>
      <span>Similitud</span>
    </div>
    <div class="flex items-center gap-2">
      <span class="text-orange-400">⚡</span>
      <span>Diferencia</span>
    </div>
    <div class="flex items-center gap-2">
      <span>▶</span>
      <span>Click para expandir detalles</span>
    </div>
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
          Compara las características de ambos conceptos
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-blue-500 to-focus-500
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

<style>
  /* Estilos adicionales si es necesario */
</style>
