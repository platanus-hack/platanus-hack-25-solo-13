<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Recursos Literarios",
    texto = "",
    autor = "",
    dispositivosLiterarios = [],
    preguntasAnalisis = [],
    modoComparacion = false,
    ejercicioCreativo = null,
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let filtroDispositivo = $state(null); // Filtrar por tipo de dispositivo
  let respuestasAnalisis = $state({});
  let respuestaCreativa = $state("");

  // Colores por tipo de dispositivo
  const dispositivoColors = {
    metafora: { bg: "bg-purple-500/30", border: "border-purple-500", text: "text-purple-300" },
    simil: { bg: "bg-blue-500/30", border: "border-blue-500", text: "text-blue-300" },
    personificacion: { bg: "bg-green-500/30", border: "border-green-500", text: "text-green-300" },
    hiperbole: { bg: "bg-orange-500/30", border: "border-orange-500", text: "text-orange-300" },
    repeticion: { bg: "bg-yellow-500/30", border: "border-yellow-500", text: "text-yellow-300" },
    aliteracion: { bg: "bg-pink-500/30", border: "border-pink-500", text: "text-pink-300" },
    anafora: { bg: "bg-teal-500/30", border: "border-teal-500", text: "text-teal-300" },
    ironia: { bg: "bg-red-500/30", border: "border-red-500", text: "text-red-300" }
  };

  // Contar frecuencia de cada dispositivo
  const frecuenciaDispositivos = $derived.by(() => {
    const freq = {};
    dispositivosLiterarios.forEach(disp => {
      freq[disp.tipo] = (freq[disp.tipo] || 0) + disp.ejemplos.length;
    });
    return freq;
  });

  // Dispositivo más usado
  const dispositivoMasUsado = $derived(
    Object.entries(frecuenciaDispositivos).sort((a, b) => b[1] - a[1])[0]
  );

  // Resaltar texto con dispositivos
  function highlightText() {
    if (!texto) return texto;

    // Simplificación: en producción se usaría HTML parsing
    return texto;
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
  class="w-full max-w-7xl mx-auto p-8 bg-slate-950 rounded-2xl border border-canvas-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
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

    {#if autor}
      <p class="text-canvas-300">
        Autor: <span class="text-purple-400 font-semibold">{autor}</span>
      </p>
    {/if}
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Texto principal -->
    <div class="lg:col-span-2">
      <!-- Texto literario -->
      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700 mb-6">
        <div class="prose prose-invert max-w-none">
          <p class="text-xl text-canvas-200 leading-relaxed whitespace-pre-line italic">
            {texto}
          </p>
        </div>
      </div>

      <!-- Dispositivos encontrados -->
      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700 mb-6">
        <h3 class="text-lg font-semibold text-white mb-4">📚 Recursos Literarios Identificados</h3>

        <div class="space-y-4">
          {#each dispositivosLiterarios as dispositivo}
            {@const colors = dispositivoColors[dispositivo.tipo] || dispositivoColors.metafora}

            <div
              class="
                p-4 rounded-xl border-2 transition-all duration-300 cursor-pointer
                {filtroDispositivo === dispositivo.tipo
                  ? `${colors.bg} ${colors.border} scale-105`
                  : 'bg-slate-950 border-canvas-700 hover:scale-102'}
              "
              onclick={() => filtroDispositivo = filtroDispositivo === dispositivo.tipo ? null : dispositivo.tipo}
            >
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1">
                  <h4 class="text-lg font-bold {colors.text} capitalize mb-2">
                    {dispositivo.tipo.replace(/_/g, ' ')}
                  </h4>

                  <p class="text-sm text-canvas-400 mb-3">
                    {dispositivo.definicion}
                  </p>

                  <!-- Ejemplos del texto -->
                  {#if filtroDispositivo === dispositivo.tipo}
                    <div class="mt-3 space-y-2">
                      <p class="text-xs font-semibold uppercase text-canvas-500">Ejemplos en el texto:</p>
                      {#each dispositivo.ejemplos as ejemplo}
                        <div class="p-2 bg-slate-950/50 rounded-lg border {colors.border}">
                          <p class="text-sm text-white italic">"{ejemplo}"</p>
                        </div>
                      {/each}
                    </div>
                  {/if}

                  <!-- Efecto/propósito -->
                  <div class="mt-3 p-3 bg-purple-500/10 rounded-lg border border-purple-500/30">
                    <p class="text-xs {colors.text}">
                      💡 <span class="font-semibold">Efecto:</span> {dispositivo.efecto}
                    </p>
                  </div>
                </div>

                <!-- Badge de conteo -->
                <div class="flex flex-col items-center gap-2">
                  <div class="w-12 h-12 rounded-full {colors.bg} border-2 {colors.border} flex items-center justify-center">
                    <span class="text-xl font-bold {colors.text}">
                      {dispositivo.ejemplos.length}
                    </span>
                  </div>
                  <p class="text-xs text-canvas-500">
                    {dispositivo.ejemplos.length === 1 ? 'vez' : 'veces'}
                  </p>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <!-- Preguntas de análisis -->
      {#if preguntasAnalisis.length > 0}
        <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
          <h3 class="text-lg font-semibold text-purple-300 mb-4">🤔 Análisis Crítico</h3>

          <div class="space-y-4">
            {#each preguntasAnalisis as pregunta, index}
              <div>
                <p class="text-canvas-200 font-medium mb-2">
                  {index + 1}. {pregunta}
                </p>

                <textarea
                  bind:value={respuestasAnalisis[index]}
                  placeholder="Escribe tu análisis aquí..."
                  class="
                    w-full px-4 py-3 bg-slate-950 border-2 border-canvas-700 rounded-xl
                    text-white placeholder-canvas-500
                    focus:outline-none focus:border-purple-500
                    transition-colors duration-300
                    min-h-20 text-sm
                  "
                ></textarea>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Sidebar -->
    <div class="lg:col-span-1 space-y-6">
      <!-- Author's Toolkit -->
      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
        <h3 class="text-sm font-semibold uppercase text-canvas-400 mb-4">🎨 Kit del Autor</h3>

        {#if dispositivoMasUsado}
          <div class="mb-4 p-3 bg-purple-500/10 rounded-lg border border-purple-500/30">
            <p class="text-xs text-purple-400 mb-1">Más usado:</p>
            <p class="text-lg font-bold text-white capitalize">
              {dispositivoMasUsado[0]}
            </p>
            <p class="text-sm text-canvas-300">
              {dispositivoMasUsado[1]} {dispositivoMasUsado[1] === 1 ? 'vez' : 'veces'}
            </p>
          </div>
        {/if}

        <!-- Gráfico de barras simple -->
        <div class="space-y-2">
          {#each Object.entries(frecuenciaDispositivos) as [tipo, cantidad]}
            {@const colors = dispositivoColors[tipo] || dispositivoColors.metafora}
            {@const maxCantidad = Math.max(...Object.values(frecuenciaDispositivos))}
            {@const porcentaje = (cantidad / maxCantidad) * 100}

            <div>
              <div class="flex items-center justify-between mb-1">
                <p class="text-xs text-canvas-400 capitalize">{tipo}</p>
                <p class="text-xs font-bold {colors.text}">{cantidad}</p>
              </div>
              <div class="w-full h-2 bg-canvas-800 rounded-full overflow-hidden">
                <div
                  class="h-full {colors.border.replace('border', 'bg')} transition-all duration-500"
                  style="width: {porcentaje}%"
                ></div>
              </div>
            </div>
          {/each}
        </div>
      </div>

      <!-- Leyenda de dispositivos -->
      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
        <h3 class="text-sm font-semibold uppercase text-canvas-400 mb-3">📖 Leyenda</h3>

        <div class="space-y-2">
          {#each dispositivosLiterarios as disp}
            {@const colors = dispositivoColors[disp.tipo] || dispositivoColors.metafora}
            <button
              onclick={() => filtroDispositivo = filtroDispositivo === disp.tipo ? null : disp.tipo}
              class="
                w-full text-left px-3 py-2 rounded-lg transition-all duration-300
                {filtroDispositivo === disp.tipo ? `${colors.bg} ${colors.border} border-2` : 'bg-slate-950 border border-canvas-700 hover:bg-canvas-800'}
              "
            >
              <div class="flex items-center gap-2">
                <div class="w-3 h-3 rounded-full {colors.border.replace('border', 'bg')}"></div>
                <span class="text-sm font-medium text-white capitalize">{disp.tipo}</span>
              </div>
            </button>
          {/each}
        </div>
      </div>

      <!-- Ejercicio creativo -->
      {#if ejercicioCreativo}
        <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
          <h3 class="text-sm font-semibold uppercase text-purple-400 mb-3">✍️ Crea</h3>

          <p class="text-sm text-canvas-300 mb-4">
            {ejercicioCreativo}
          </p>

          <textarea
            bind:value={respuestaCreativa}
            placeholder="Escribe tu creación aquí..."
            class="
              w-full px-4 py-3 bg-slate-950 border-2 border-canvas-700 rounded-xl
              text-white placeholder-canvas-500
              focus:outline-none focus:border-purple-500
              transition-colors duration-300
              min-h-32 text-sm
            "
          ></textarea>

          <button
            disabled={!respuestaCreativa.trim()}
            class="
              w-full mt-3 px-4 py-2 rounded-lg font-semibold text-sm
              bg-purple-500/20 text-purple-300
              border border-purple-500/50
              transition-all duration-300
              hover:bg-purple-500/30 hover:scale-105
              disabled:opacity-30 disabled:cursor-not-allowed
            "
          >
            💾 Guardar Creación
          </button>
        </div>
      {/if}
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
          {dispositivosLiterarios.length} recursos identificados
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
