<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Conectores y Coherencia",
    parrafos = [],
    categorias = ["causales", "adversativos", "consecutivos", "aditivos", "temporales"],
    bancoConectores = {
      causales: ["porque", "ya que", "debido a", "puesto que"],
      adversativos: ["sin embargo", "no obstante", "pero", "aunque"],
      consecutivos: ["por lo tanto", "en consecuencia", "así que", "entonces"],
      aditivos: ["además", "asimismo", "también", "igualmente"],
      temporales: ["luego", "después", "mientras", "finalmente", "antes"]
    },
    modoComparacion = false,
    ejercicioCreativo = null,
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let respuestasUsuario = $state({}); // { parrafoIndex: conectorElegido }
  let mostrarBanco = $state(false);
  let categoriaFiltro = $state("todos");
  let feedback = $state(null);
  let respuestaCreativa = $state("");

  // Colores por categoría
  const categoriaColors = {
    causal: { bg: "bg-blue-500/20", border: "border-blue-500", text: "text-blue-300" },
    adversativo: { bg: "bg-red-500/20", border: "border-red-500", text: "text-red-300" },
    consecutivo: { bg: "bg-green-500/20", border: "border-green-500", text: "text-green-300" },
    aditivo: { bg: "bg-purple-500/20", border: "border-purple-500", text: "text-purple-300" },
    temporal: { bg: "bg-orange-500/20", border: "border-orange-500", text: "text-orange-300" }
  };

  // Nombre plural a singular
  const categoriaPlural = {
    causales: "causal",
    adversativos: "adversativo",
    consecutivos: "consecutivo",
    aditivos: "aditivo",
    temporales: "temporal"
  };

  // Conectores filtrados
  const conectoresFiltrados = $derived.by(() => {
    if (categoriaFiltro === "todos") {
      return Object.values(bancoConectores).flat();
    }
    return bancoConectores[categoriaFiltro] || [];
  });

  // Progreso
  const progreso = $derived(
    (Object.keys(respuestasUsuario).length / parrafos.length) * 100
  );

  // Correctas
  const respuestasCorrectas = $derived(
    parrafos.filter((p, i) => respuestasUsuario[i] === p.correcta).length
  );

  function seleccionarConector(parrafoIndex, conector) {
    respuestasUsuario = { ...respuestasUsuario, [parrafoIndex]: conector };
  }

  function verificarRespuestas() {
    const total = parrafos.length;
    const correctas = respuestasCorrectas;

    feedback = {
      tipo: correctas === total ? "success" : "info",
      mensaje: `${correctas} de ${total} conectores correctos`
    };

    setTimeout(() => {
      feedback = null;
    }, 5000);
  }

  function reiniciar() {
    respuestasUsuario = {};
    feedback = null;
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
        Conectores
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>

    <p class="text-canvas-300">
      Los conectores ayudan a dar coherencia y fluidez al texto
    </p>
  </div>

  <!-- Progreso -->
  <div class="mb-6">
    <div class="flex items-center justify-between mb-2">
      <p class="text-sm text-canvas-400">Progreso</p>
      <p class="text-sm font-semibold text-purple-400">
        {Object.keys(respuestasUsuario).length} / {parrafos.length} completados
      </p>
    </div>
    <div class="w-full h-3 bg-canvas-800 rounded-full overflow-hidden">
      <div
        class="h-full bg-gradient-to-r from-purple-500 to-purple-600 transition-all duration-500"
        style="width: {progreso}%"
      ></div>
    </div>
  </div>

  <!-- Párrafos con ejercicios -->
  <div class="mb-8 space-y-6">
    {#each parrafos as parrafo, index}
      {@const tipoSingular = categoriaPlural[parrafo.tipo] || "causal"}
      {@const colors = categoriaColors[tipoSingular] || categoriaColors.causal}
      {@const partes = parrafo.texto.split("___")}

      <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
        <div class="flex items-start gap-3 mb-4">
          <span class="text-2xl font-bold text-purple-400">{index + 1}</span>
          <div class="flex-1">
            <!-- Párrafo con espacio para conector -->
            <div class="text-lg text-canvas-200 leading-relaxed mb-4">
              {partes[0]}
              {#if respuestasUsuario[index]}
                <span class="
                  px-3 py-1 rounded-lg font-semibold
                  {colors.bg} {colors.text} {colors.border} border-2
                ">
                  {respuestasUsuario[index]}
                </span>
              {:else}
                <span class="px-4 py-1 bg-canvas-800 rounded-lg border-2 border-dashed border-purple-500 text-purple-400 font-semibold">
                  ___
                </span>
              {/if}
              {partes[1] || ""}
            </div>

            <!-- Opciones de conectores -->
            <div class="flex flex-wrap gap-2 mb-3">
              {#each parrafo.opcionesConector as opcion}
                <button
                  onclick={() => seleccionarConector(index, opcion)}
                  class="
                    px-4 py-2 rounded-lg font-semibold text-sm
                    transition-all duration-300
                    {respuestasUsuario[index] === opcion
                      ? `${colors.bg} ${colors.text} ${colors.border} border-2 scale-105`
                      : 'bg-canvas-800 text-canvas-300 border border-canvas-700 hover:bg-canvas-700 hover:scale-105'}
                  "
                >
                  {opcion}
                </button>
              {/each}
            </div>

            <!-- Badge del tipo -->
            <div class="flex items-center gap-2">
              <span class="px-2 py-1 rounded text-xs font-semibold {colors.bg} {colors.text} capitalize">
                Conector {tipoSingular}
              </span>

              {#if respuestasUsuario[index]}
                {#if respuestasUsuario[index] === parrafo.correcta}
                  <span class="text-green-400 text-sm">✓ Correcto</span>
                {:else}
                  <span class="text-red-400 text-sm">✗ Incorrecto</span>
                {/if}
              {/if}
            </div>

            <!-- Explicación (si está respondido) -->
            {#if respuestasUsuario[index] && parrafo.explicacion}
              <div class="mt-3 p-3 bg-slate-950 rounded-lg border border-canvas-700">
                <p class="text-xs text-canvas-400 italic">
                  💡 {parrafo.explicacion}
                </p>
              </div>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>

  <!-- Botones de acción -->
  <div class="flex flex-wrap items-center gap-3 mb-8">
    <button
      onclick={verificarRespuestas}
      disabled={Object.keys(respuestasUsuario).length === 0}
      class="
        px-6 py-3 rounded-xl font-semibold
        bg-gradient-to-r from-purple-500 to-purple-600
        text-white
        transition-all duration-300
        hover:shadow-lg hover:shadow-purple-500/50 hover:scale-105
        disabled:opacity-30 disabled:cursor-not-allowed
      "
    >
      ✓ Verificar Respuestas
    </button>

    <button
      onclick={reiniciar}
      class="
        px-6 py-3 rounded-xl font-semibold
        bg-canvas-800 text-canvas-300
        border border-canvas-700
        transition-all duration-300
        hover:bg-canvas-700 hover:scale-105
      "
    >
      ↻ Reiniciar
    </button>

    <button
      onclick={() => mostrarBanco = !mostrarBanco}
      class="
        px-6 py-3 rounded-xl font-semibold text-sm
        bg-canvas-800 text-canvas-300
        border border-canvas-700
        transition-all duration-300
        hover:bg-canvas-700 hover:scale-105
      "
    >
      {mostrarBanco ? '✕ Cerrar' : '📚'} Banco de Conectores
    </button>
  </div>

  <!-- Feedback -->
  {#if feedback}
    <div class="
      mb-6 p-4 rounded-xl border-2
      {feedback.tipo === 'success' ? 'bg-green-500/10 border-green-500' : 'bg-blue-500/10 border-blue-500'}
    ">
      <p class="
        font-semibold
        {feedback.tipo === 'success' ? 'text-green-300' : 'text-blue-300'}
      ">
        {feedback.mensaje}
        {#if feedback.tipo === 'success'}
          🎉 ¡Perfecto! Dominas el uso de conectores.
        {/if}
      </p>
    </div>
  {/if}

  <!-- Banco de conectores (expandible) -->
  {#if mostrarBanco}
    <div class="mb-8 p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
      <h3 class="text-lg font-semibold text-purple-300 mb-4">📚 Banco de Conectores</h3>

      <!-- Filtro por categoría -->
      <div class="flex flex-wrap gap-2 mb-4">
        <button
          onclick={() => categoriaFiltro = "todos"}
          class="
            px-3 py-2 rounded-lg text-sm font-semibold
            transition-all duration-300
            {categoriaFiltro === 'todos' ? 'bg-purple-500 text-white' : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
          "
        >
          Todos
        </button>
        {#each Object.keys(bancoConectores) as cat}
          <button
            onclick={() => categoriaFiltro = cat}
            class="
              px-3 py-2 rounded-lg text-sm font-semibold capitalize
              transition-all duration-300
              {categoriaFiltro === cat ? 'bg-purple-500 text-white' : 'bg-canvas-800 text-canvas-300 hover:bg-canvas-700'}
            "
          >
            {cat}
          </button>
        {/each}
      </div>

      <!-- Lista de conectores -->
      <div class="flex flex-wrap gap-2">
        {#each conectoresFiltrados as conector}
          <span class="px-4 py-2 bg-slate-950 rounded-lg border border-canvas-700 text-canvas-200 font-medium">
            {conector}
          </span>
        {/each}
      </div>
    </div>
  {/if}

  <!-- Ejercicio creativo (opcional) -->
  {#if ejercicioCreativo}
    <div class="mb-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-3">✍️ Ejercicio Creativo</h3>
      <p class="text-canvas-300 mb-4">
        {ejercicioCreativo}
      </p>

      <textarea
        bind:value={respuestaCreativa}
        placeholder="Escribe tu párrafo aquí..."
        class="
          w-full px-4 py-3 bg-slate-950 border-2 border-canvas-700 rounded-xl
          text-white placeholder-canvas-500
          focus:outline-none focus:border-purple-500
          transition-colors duration-300
          min-h-32
        "
      ></textarea>
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
          {respuestasCorrectas} / {parrafos.length} correctas
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
