<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Estructura del Texto",
    tipoTexto = "argumentativo", // narrativo | argumentativo | expositivo | instructivo
    textoEjemplo = "",
    estructura = {}, // { seccion: { texto, color, descripcion } }
    ejercicioTipo = "solo-visualizar", // identificar-partes | solo-visualizar
    comparacionTipos = false,
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let vistaActual = $state("estructura"); // estructura | texto
  let seccionSeleccionada = $state(null);
  let draggedFragment = $state(null);

  // Definiciones de estructuras por tipo
  const estructurasTipos = {
    narrativo: {
      nombre: "Narrativo",
      color: "cyan",
      secciones: ["inicio", "desarrollo", "climax", "desenlace"],
      descripcion: "Relata una historia con personajes, espacio y tiempo"
    },
    argumentativo: {
      nombre: "Argumentativo",
      color: "purple",
      secciones: ["introduccion", "tesis", "argumentos", "contraargumentos", "conclusion"],
      descripcion: "Defiende una postura con argumentos y evidencia"
    },
    expositivo: {
      nombre: "Expositivo",
      color: "blue",
      secciones: ["introduccion", "desarrollo", "conclusion"],
      descripcion: "Explica o informa sobre un tema de forma objetiva"
    },
    instructivo: {
      nombre: "Instructivo",
      color: "green",
      secciones: ["objetivo", "materiales", "pasos", "resultado"],
      descripcion: "Enseña cómo hacer algo paso a paso"
    }
  };

  const tipoActual = $derived(estructurasTipos[tipoTexto] || estructurasTipos.argumentativo);

  // Colores por sección (genéricos)
  const seccionColors = {
    introduccion: { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300" },
    tesis: { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300" },
    argumentos: { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300" },
    contraargumentos: { bg: "bg-red-500/20", border: "border-red-500", text: "text-red-300" },
    conclusion: { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300" },
    inicio: { bg: "bg-cyan-500/20", border: "border-cyan-500", text: "text-cyan-300" },
    desarrollo: { bg: "bg-teal-500/20", border: "border-teal-500", text: "text-teal-300" },
    climax: { bg: "bg-yellow-500/20", border: "border-yellow-500", text: "text-yellow-300" },
    desenlace: { bg: "bg-pink-500/20", border: "border-pink-500", text: "text-pink-300" },
    objetivo: { bg: "bg-lime-500/20", border: "border-lime-500", text: "text-lime-300" },
    materiales: { bg: "bg-emerald-500/20", border: "border-emerald-500", text: "text-emerald-300" },
    pasos: { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300" },
    resultado: { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300" }
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
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-purple-500/20 text-purple-400 border border-purple-500">
        {materia}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-medium bg-canvas-800 text-canvas-300">
        Estructura de Textos
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      Tipo: <span class="text-purple-400 font-semibold">{tipoActual.nombre}</span> - {tipoActual.descripcion}
    </p>
  </div>

  <!-- Toggle vista -->
  <div class="mb-6 flex items-center justify-center gap-3">
    <button
      onclick={() => vistaActual = "estructura"}
      class="
        px-6 py-3 rounded-xl font-semibold
        transition-all duration-300
        {vistaActual === 'estructura' ? 'bg-purple-500 text-white shadow-lg shadow-purple-500/50' : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
      "
    >
      📊 Vista Estructura
    </button>

    <button
      onclick={() => vistaActual = "texto"}
      class="
        px-6 py-3 rounded-xl font-semibold
        transition-all duration-300
        {vistaActual === 'texto' ? 'bg-purple-500 text-white shadow-lg shadow-purple-500/50' : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
      "
    >
      📖 Vista Texto Completo
    </button>
  </div>

  <!-- Vista Estructura -->
  {#if vistaActual === "estructura"}
    <div class="space-y-4 mb-8">
      {#each Object.entries(estructura) as [seccion, contenido], index}
        {@const colors = seccionColors[seccion] || seccionColors.desarrollo}
        <div
          class="
            p-6 rounded-2xl border-2 transition-all duration-300 cursor-pointer
            {colors.bg} {colors.border}
            {seccionSeleccionada === seccion ? 'scale-105 shadow-2xl' : 'hover:scale-102'}
          "
          onclick={() => seccionSeleccionada = seccionSeleccionada === seccion ? null : seccion}
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <span class="text-3xl font-bold {colors.text}">
                  {index + 1}
                </span>
                <h3 class="text-xl font-bold {colors.text} capitalize">
                  {seccion.replace(/_/g, ' ')}
                </h3>
              </div>

              {#if contenido.descripcion}
                <p class="text-sm text-canvas-400 mb-3 italic">
                  {contenido.descripcion}
                </p>
              {/if}

              {#if seccionSeleccionada === seccion}
                <div class="mt-4 p-4 bg-slate-950/50 rounded-lg border border-canvas-700">
                  <p class="text-canvas-200 leading-relaxed">
                    {contenido.texto}
                  </p>
                </div>
              {:else}
                <p class="text-canvas-300 line-clamp-2">
                  {contenido.texto}
                </p>
              {/if}
            </div>

            <button class="text-canvas-500 hover:text-white transition-colors">
              {seccionSeleccionada === seccion ? '▼' : '▶'}
            </button>
          </div>
        </div>
      {/each}
    </div>

    <!-- Leyenda de estructura -->
    <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4">💡 Estructura del Texto {tipoActual.nombre}</h3>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each tipoActual.secciones as seccion}
          {@const colors = seccionColors[seccion] || seccionColors.desarrollo}
          <div class="flex items-center gap-3 p-3 bg-slate-950 rounded-lg border border-canvas-700">
            <div class="w-4 h-4 rounded-full {colors.border.replace('border', 'bg')}"></div>
            <div>
              <p class="text-sm font-semibold {colors.text} capitalize">
                {seccion.replace(/_/g, ' ')}
              </p>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Vista Texto Completo -->
  {#if vistaActual === "texto"}
    <div class="p-8 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <div class="prose prose-invert max-w-none">
        <div class="space-y-6">
          {#each Object.entries(estructura) as [seccion, contenido]}
            {@const colors = seccionColors[seccion] || seccionColors.desarrollo}
            <div class="relative">
              <div class="absolute left-0 top-0 bottom-0 w-1 {colors.border.replace('border', 'bg')} rounded"></div>
              <div class="pl-6">
                <h4 class="text-sm font-bold uppercase {colors.text} mb-2">
                  {seccion.replace(/_/g, ' ')}
                </h4>
                <p class="text-canvas-200 leading-relaxed">
                  {contenido.texto}
                </p>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  {/if}

  <!-- Comparación de tipos (opcional) -->
  {#if comparacionTipos}
    <div class="mt-8 p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
      <h3 class="text-lg font-semibold text-purple-300 mb-4">🔄 Comparación de Estructuras</h3>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {#each Object.entries(estructurasTipos) as [tipo, config]}
          <div class="p-4 bg-slate-950 rounded-lg border border-canvas-700">
            <h4 class="text-sm font-bold text-white mb-2">{config.nombre}</h4>
            <ul class="space-y-1 text-xs text-canvas-400">
              {#each config.secciones as seccion}
                <li class="flex items-center gap-2">
                  <span class="w-1.5 h-1.5 bg-purple-500 rounded-full"></span>
                  <span class="capitalize">{seccion}</span>
                </li>
              {/each}
            </ul>
          </div>
        {/each}
      </div>

      <p class="text-sm text-purple-400 mt-4 italic">
        💡 Cada tipo de texto tiene una estructura diferente según su propósito comunicativo
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
          Texto {tipoActual.nombre}
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
