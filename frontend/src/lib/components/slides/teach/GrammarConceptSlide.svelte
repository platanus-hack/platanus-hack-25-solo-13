<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Concepto Gramatical",
    concepto = "",
    definicion = "",
    tipos = [],
    reglas = [],
    erroresComunes = [],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let tipoSeleccionado = $state(0);

  // Colores por tipo
  const tipoColors = [
    { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300", badgeBg: "bg-blue-500" },
    { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300", badgeBg: "bg-green-500" },
    { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300", badgeBg: "bg-purple-500" },
    { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300", badgeBg: "bg-orange-500" },
    { bg: "bg-pink-500/20", border: "border-pink-500", text: "text-pink-300", badgeBg: "bg-pink-500" }
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
        Gramática
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      {concepto}
    </p>
  </div>

  <!-- Definición principal -->
  <div class="mb-8 p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
    <h3 class="text-sm font-semibold uppercase text-purple-400 mb-3">
      📖 Definición
    </h3>
    <p class="text-lg text-canvas-200 leading-relaxed">
      {definicion}
    </p>
  </div>

  <!-- Tabla de tipos con ejemplos -->
  {#if tipos.length > 0}
    <div class="mb-8">
      <h3 class="text-lg font-semibold text-white mb-4">
        Tipos de {concepto}
      </h3>

      <!-- Tabs para tipos -->
      <div class="flex flex-wrap gap-2 mb-6">
        {#each tipos as tipo, index}
          {@const colors = tipoColors[index % 5]}
          <button
            onclick={() => tipoSeleccionado = index}
            class="
              px-4 py-2 rounded-lg font-semibold text-sm
              transition-all duration-300
              {tipoSeleccionado === index
                ? `${colors.badgeBg} text-white scale-105`
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
        {@const colors = tipoColors[tipoSeleccionado % 5]}

        <div class="p-6 {colors.bg} rounded-2xl border-2 {colors.border}">
          <!-- Definición del tipo -->
          <h4 class="text-xl font-bold {colors.text} mb-3">
            {tipo.nombre}
          </h4>
          <p class="text-canvas-200 mb-6">
            {tipo.definicion}
          </p>

          <!-- Características -->
          {#if tipo.caracteristicas && tipo.caracteristicas.length > 0}
            <div class="mb-6">
              <p class="text-xs font-semibold uppercase text-canvas-400 mb-3">
                Características:
              </p>
              <ul class="space-y-2">
                {#each tipo.caracteristicas as caracteristica}
                  <li class="flex items-start gap-2 text-canvas-200">
                    <span class="{colors.text}">•</span>
                    <span>{caracteristica}</span>
                  </li>
                {/each}
              </ul>
            </div>
          {/if}

          <!-- Ejemplos -->
          {#if tipo.ejemplos && tipo.ejemplos.length > 0}
            <div class="space-y-3">
              <p class="text-xs font-semibold uppercase text-canvas-400">
                Ejemplos:
              </p>
              {#each tipo.ejemplos as ejemplo}
                <div class="p-4 bg-slate-950 rounded-lg border {colors.border}">
                  <p class="text-lg text-white mb-2 font-medium">
                    "{ejemplo.oracion}"
                  </p>
                  {#if ejemplo.analisis}
                    <p class="text-sm text-canvas-400 italic">
                      → {ejemplo.analisis}
                    </p>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}

          <!-- Fórmula/Estructura (opcional) -->
          {#if tipo.estructura}
            <div class="mt-6 p-4 bg-canvas-900/50 rounded-lg border border-canvas-700">
              <p class="text-xs font-semibold uppercase text-canvas-400 mb-2">
                Estructura:
              </p>
              <p class="text-lg font-mono {colors.text}">
                {tipo.estructura}
              </p>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  {/if}

  <!-- Reglas importantes -->
  {#if reglas.length > 0}
    <div class="mb-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">
        📏 Reglas Importantes
      </h3>
      <div class="space-y-3">
        {#each reglas as regla, index}
          <div class="flex items-start gap-3 p-4 bg-slate-950 rounded-lg">
            <div class="w-8 h-8 rounded-full bg-blue-500/20 border-2 border-blue-500 flex items-center justify-center flex-shrink-0">
              <span class="text-sm font-bold text-blue-300">
                {index + 1}
              </span>
            </div>
            <div class="flex-1">
              <p class="text-canvas-200">
                {regla.texto}
              </p>
              {#if regla.ejemplo}
                <p class="text-sm text-canvas-400 italic mt-2">
                  Ejemplo: "{regla.ejemplo}"
                </p>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Errores comunes -->
  {#if erroresComunes.length > 0}
    <div class="mb-8 p-6 bg-red-500/10 rounded-2xl border border-red-500/30">
      <h3 class="text-lg font-semibold text-red-300 mb-4">
        ⚠️ Errores Comunes a Evitar
      </h3>
      <div class="space-y-4">
        {#each erroresComunes as error}
          <div class="p-4 bg-slate-950 rounded-lg border border-red-500/30">
            <div class="flex items-start gap-3 mb-2">
              <span class="text-red-400">✗</span>
              <p class="text-canvas-200 font-medium">
                {error.incorrecto}
              </p>
            </div>
            <div class="flex items-start gap-3">
              <span class="text-green-400">✓</span>
              <p class="text-canvas-200 font-medium">
                {error.correcto}
              </p>
            </div>
            {#if error.explicacion}
              <p class="text-sm text-canvas-400 italic mt-3 ml-6">
                💡 {error.explicacion}
              </p>
            {/if}
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
          {tipos.length} tipos • {reglas.length} reglas
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
