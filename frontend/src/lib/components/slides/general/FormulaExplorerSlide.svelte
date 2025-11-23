<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Explorador de Fórmula",
    formula = "E = mc²",
    variables = [],
    ejemploNumerico = null,
    calculadoraInteractiva = true,
    materia = "física",
    mostrarUnidades = true,
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let hoveredVariable = $state(null);
  let valoresCalculadora = $state({});
  let mostrarEjemplo = $state(false);

  // Inicializar valores de la calculadora
  $effect(() => {
    const valores = {};
    variables.forEach(v => {
      if (ejemploNumerico && ejemploNumerico[v.simbolo] !== undefined) {
        valores[v.simbolo] = ejemploNumerico[v.simbolo];
      } else {
        valores[v.simbolo] = 0;
      }
    });
    valoresCalculadora = valores;
  });

  // Calcular resultado (eval seguro - solo para demo, en producción usar parser)
  const resultadoCalculado = $derived.by(() => {
    if (!calculadoraInteractiva) return null;

    try {
      // Reemplazar símbolos con valores
      let formulaEvaluable = formula;

      // Convertir símbolos matemáticos comunes
      formulaEvaluable = formulaEvaluable.replace(/²/g, '**2');
      formulaEvaluable = formulaEvaluable.replace(/³/g, '**3');
      formulaEvaluable = formulaEvaluable.replace(/×/g, '*');
      formulaEvaluable = formulaEvaluable.replace(/÷/g, '/');

      // Reemplazar variables con valores
      variables.forEach(v => {
        const regex = new RegExp(`\\b${v.simbolo}\\b`, 'g');
        formulaEvaluable = formulaEvaluable.replace(regex, valoresCalculadora[v.simbolo] || 0);
      });

      // Evaluar (solo para demo - usar mathjs en producción)
      const resultado = eval(formulaEvaluable);
      return isNaN(resultado) ? 0 : resultado;
    } catch (e) {
      return 0;
    }
  });

  // Resaltar variables en la fórmula
  function highlightFormula(formulaTexto) {
    let resultado = formulaTexto;
    variables.forEach(v => {
      const regex = new RegExp(`\\b(${v.simbolo})\\b`, 'g');
      resultado = resultado.replace(regex, `<mark class="formula-variable" data-variable="${v.simbolo}">$1</mark>`);
    });
    return resultado;
  }

  const formulaMarcada = $derived(highlightFormula(formula));

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
  class="w-full max-w-6xl mx-auto p-8 bg-slate-950 rounded-2xl border border-canvas-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-focus-500/20 text-focus-400 border border-focus-500">
        {materia}
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      {titulo}
    </h2>
  </div>

  <!-- Fórmula principal -->
  <div class="mb-8 text-center">
    <div class="inline-block p-8 bg-gradient-to-br from-focus-500/10 to-lumera-500/10 rounded-2xl border-2 border-focus-500/30">
      <p class="text-5xl font-bold text-white formula-display">
        {@html formulaMarcada}
      </p>
    </div>
  </div>

  <!-- Variables y sus definiciones -->
  <div class="mb-8">
    <h3 class="text-lg font-semibold text-canvas-300 mb-4">
      Variables de la fórmula:
    </h3>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {#each variables as variable}
        <div
          class="
            p-4 bg-canvas-900/50 rounded-xl border-2 transition-all duration-300 cursor-help
            {hoveredVariable === variable.simbolo ? 'border-focus-500 bg-focus-500/5 scale-105' : 'border-canvas-700'}
          "
          onmouseenter={() => hoveredVariable = variable.simbolo}
          onmouseleave={() => hoveredVariable = null}
        >
          <!-- Símbolo -->
          <div class="flex items-center gap-3 mb-2">
            <span class="text-3xl font-bold text-focus-400">
              {variable.simbolo}
            </span>
            <div class="flex-1">
              <p class="text-sm font-semibold text-white">
                {variable.nombre}
              </p>
              {#if mostrarUnidades && variable.unidad}
                <p class="text-xs text-canvas-400">
                  ({variable.unidad})
                </p>
              {/if}
            </div>
          </div>

          <!-- Descripción -->
          <p class="text-sm text-canvas-400">
            {variable.descripcion}
          </p>
        </div>
      {/each}
    </div>
  </div>

  <!-- Calculadora interactiva -->
  {#if calculadoraInteractiva}
    <div class="mb-8 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
      <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
        <span>🧮</span>
        <span>Calculadora Interactiva</span>
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Inputs de variables -->
        <div class="space-y-4">
          {#each variables as variable}
            <div>
              <label class="block text-sm font-semibold text-canvas-300 mb-2">
                {variable.simbolo} - {variable.nombre}
                {#if mostrarUnidades && variable.unidad}
                  <span class="text-xs text-canvas-500">({variable.unidad})</span>
                {/if}
              </label>
              <input
                type="number"
                bind:value={valoresCalculadora[variable.simbolo]}
                step="0.01"
                class="
                  w-full px-4 py-3 bg-slate-950 border-2 border-canvas-700 rounded-xl
                  text-white font-mono text-lg
                  focus:outline-none focus:border-focus-500
                  transition-colors duration-300
                "
              />
            </div>
          {/each}
        </div>

        <!-- Resultado -->
        <div class="flex items-center justify-center">
          <div class="text-center p-8 bg-gradient-to-br from-focus-500/10 to-lumera-500/10 rounded-2xl border-2 border-focus-500/30 w-full">
            <p class="text-sm font-semibold text-focus-400 mb-2 uppercase">
              Resultado:
            </p>
            <p class="text-5xl font-bold text-white mb-2">
              {resultadoCalculado !== null ? resultadoCalculado.toFixed(2) : '0'}
            </p>
            {#if ejemploNumerico?.unidadResultado}
              <p class="text-sm text-canvas-400">
                {ejemploNumerico.unidadResultado}
              </p>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {/if}

  <!-- Ejemplo numérico resuelto -->
  {#if ejemploNumerico}
    <div class="mb-8">
      <button
        onclick={() => mostrarEjemplo = !mostrarEjemplo}
        class="
          w-full px-6 py-4 bg-purple-500/20 rounded-xl border border-purple-500/50
          text-purple-300 font-semibold
          transition-all duration-300
          hover:bg-purple-500/30 hover:scale-105
        "
      >
        {mostrarEjemplo ? '▼' : '▶'} Ver Ejemplo Resuelto
      </button>

      {#if mostrarEjemplo}
        <div class="mt-4 p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h3 class="text-lg font-semibold text-white mb-4">
            Ejemplo:
          </h3>

          <div class="space-y-3">
            <!-- Datos -->
            <div>
              <p class="text-sm font-semibold text-purple-400 mb-2">Datos:</p>
              <div class="pl-4 space-y-1">
                {#each variables as variable}
                  {#if ejemploNumerico[variable.simbolo] !== undefined}
                    <p class="text-canvas-300 font-mono">
                      {variable.simbolo} = {ejemploNumerico[variable.simbolo]}
                      {#if variable.unidad}
                        {variable.unidad}
                      {/if}
                    </p>
                  {/if}
                {/each}
              </div>
            </div>

            <!-- Sustitución -->
            <div>
              <p class="text-sm font-semibold text-purple-400 mb-2">Sustitución:</p>
              <p class="text-canvas-300 font-mono pl-4">
                {ejemploNumerico.sustitucion || formula}
              </p>
            </div>

            <!-- Resultado -->
            <div>
              <p class="text-sm font-semibold text-purple-400 mb-2">Resultado:</p>
              <p class="text-xl font-bold text-focus-400 pl-4">
                {ejemploNumerico.resultado}
                {#if ejemploNumerico.unidadResultado}
                  {ejemploNumerico.unidadResultado}
                {/if}
              </p>
            </div>
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Navegación -->
  {#if showNavigation}
    <div class="flex items-center justify-between pt-6 border-t border-canvas-800">
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
          Experimenta con diferentes valores
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-focus-500 to-lumera-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-focus-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        Siguiente →
      </button>
    </div>
  {/if}
</div>

<style>
  :global(.formula-variable) {
    color: var(--color-focus-400);
    background-color: rgb(from var(--color-focus-500) r g b / 0.2);
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    border-bottom-width: 2px;
    border-color: var(--color-focus-500);
    cursor: help;
  }

  :global(.formula-display) {
    line-height: 1.625;
  }
</style>
