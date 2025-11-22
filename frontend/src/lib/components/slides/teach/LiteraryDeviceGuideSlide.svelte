<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Guía de Recursos Literarios",
    dispositivos = [],
    comparaciones = [],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let dispositivoSeleccionado = $state(0);

  // Colores por dispositivo
  const dispositivoColors = [
    { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300", badge: "bg-purple-500" },
    { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300", badge: "bg-blue-500" },
    { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300", badge: "bg-green-500" },
    { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300", badge: "bg-orange-500" },
    { bg: "bg-pink-500/20", border: "border-pink-500", text: "text-pink-300", badge: "bg-pink-500" }
  ];

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
        Literatura
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      Aprende los recursos estilísticos que enriquecen los textos literarios
    </p>
  </div>

  <!-- Tabs de dispositivos -->
  <div class="mb-6 flex flex-wrap gap-2">
    {#each dispositivos as dispositivo, index}
      {@const colors = dispositivoColors[index % 5]}
      <button
        onclick={() => dispositivoSeleccionado = index}
        class="
          px-5 py-3 rounded-xl font-semibold capitalize
          transition-all duration-300
          {dispositivoSeleccionado === index
            ? `${colors.badge} text-white scale-105 shadow-lg`
            : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
        "
      >
        {dispositivo.nombre}
      </button>
    {/each}
  </div>

  <!-- Contenido del dispositivo seleccionado -->
  {#if dispositivos[dispositivoSeleccionado]}
    {@const dispositivo = dispositivos[dispositivoSeleccionado]}
    {@const colors = dispositivoColors[dispositivoSeleccionado % 5]}

    <div class="space-y-6">
      <!-- Definición y propósito -->
      <div class="p-6 {colors.bg} rounded-2xl border-2 {colors.border}">
        <h3 class="text-2xl font-bold {colors.text} mb-3 capitalize">
          {dispositivo.nombre}
        </h3>
        <p class="text-lg text-canvas-200 mb-4">
          {dispositivo.definicion}
        </p>
        <div class="p-4 bg-canvas-950 rounded-lg border {colors.border}">
          <p class="text-sm text-canvas-300">
            <span class="font-semibold {colors.text}">Efecto en el lector:</span> {dispositivo.efecto}
          </p>
        </div>
      </div>

      <!-- Ejemplos variados -->
      {#if dispositivo.ejemplos && dispositivo.ejemplos.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-lg font-semibold text-white mb-4">
            📖 Ejemplos
          </h4>
          <div class="space-y-4">
            {#each dispositivo.ejemplos as ejemplo, idx}
              <div class="p-4 bg-canvas-950 rounded-xl border {colors.border}">
                <div class="flex items-start gap-3 mb-3">
                  <div class="w-8 h-8 rounded-full {colors.badge} flex items-center justify-center flex-shrink-0">
                    <span class="text-white font-bold">{idx + 1}</span>
                  </div>
                  <p class="text-lg text-white italic flex-1">
                    "{ejemplo.texto}"
                  </p>
                </div>
                {#if ejemplo.autor}
                  <p class="text-sm text-canvas-400 mb-2">
                    — {ejemplo.autor}
                  </p>
                {/if}
                {#if ejemplo.analisis}
                  <div class="mt-3 p-3 bg-canvas-900/50 rounded-lg">
                    <p class="text-sm text-canvas-300">
                      💡 {ejemplo.analisis}
                    </p>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Cómo identificarlo -->
      {#if dispositivo.comoIdentificar && dispositivo.comoIdentificar.length > 0}
        <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
          <h4 class="text-sm font-semibold uppercase text-purple-400 mb-3">
            🔍 Cómo Identificar {dispositivo.nombre}
          </h4>
          <ul class="space-y-2">
            {#each dispositivo.comoIdentificar as pista}
              <li class="flex items-start gap-2 text-canvas-200">
                <span class="text-purple-400">•</span>
                <span>{pista}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/if}

      <!-- Variantes o tipos -->
      {#if dispositivo.variantes && dispositivo.variantes.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-lg font-semibold text-white mb-4">
            🎨 Variantes de {dispositivo.nombre}
          </h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            {#each dispositivo.variantes as variante}
              <div class="p-4 bg-canvas-950 rounded-lg border {colors.border}">
                <h5 class="font-bold {colors.text} mb-2">{variante.tipo}</h5>
                <p class="text-sm text-canvas-200 mb-2">{variante.descripcion}</p>
                {#if variante.ejemplo}
                  <p class="text-xs text-canvas-400 italic">
                    Ejemplo: "{variante.ejemplo}"
                  </p>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- Cuándo usar este dispositivo -->
      {#if dispositivo.cuandoUsar}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h4 class="text-sm font-semibold uppercase text-canvas-400 mb-3">
            ✍️ Cuándo usar {dispositivo.nombre}
          </h4>
          <p class="text-canvas-200">
            {dispositivo.cuandoUsar}
          </p>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Comparaciones entre dispositivos -->
  {#if comparaciones.length > 0}
    <div class="mt-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        ⚖️ Comparaciones Entre Dispositivos
      </h3>
      <div class="space-y-4">
        {#each comparaciones as comparacion}
          <div class="p-4 bg-canvas-950 rounded-xl border border-canvas-700">
            <h4 class="font-bold text-purple-300 mb-3">
              {comparacion.dispositivos.join(' vs ')}
            </h4>
            <p class="text-canvas-200 mb-3">
              {comparacion.diferencia}
            </p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              {#each comparacion.ejemplos as ejemplo}
                <div class="p-3 bg-canvas-900/50 rounded-lg">
                  <p class="text-sm font-semibold text-canvas-300 mb-1">
                    {ejemplo.tipo}:
                  </p>
                  <p class="text-sm text-canvas-400 italic">
                    "{ejemplo.ejemplo}"
                  </p>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Tabla resumen -->
  <div class="mt-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
    <h3 class="text-lg font-semibold text-white mb-4">
      📋 Tabla Resumen
    </h3>
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-canvas-700">
            <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Dispositivo</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Definición</th>
            <th class="text-left py-3 px-4 text-sm font-semibold text-canvas-300">Efecto</th>
          </tr>
        </thead>
        <tbody>
          {#each dispositivos as disp, idx}
            {@const colors = dispositivoColors[idx % 5]}
            <tr class="border-b border-canvas-800 hover:bg-canvas-800/30">
              <td class="py-3 px-4">
                <span class="px-3 py-1 {colors.badge} rounded text-white text-sm font-semibold capitalize">
                  {disp.nombre}
                </span>
              </td>
              <td class="py-3 px-4 text-sm text-canvas-300">
                {disp.definicion.split('.')[0]}.
              </td>
              <td class="py-3 px-4 text-sm text-canvas-400">
                {disp.efecto.split('.')[0]}.
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
          {dispositivos.length} dispositivos literarios
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
