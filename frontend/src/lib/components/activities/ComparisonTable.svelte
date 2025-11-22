<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  /**
   * ComparisonTable - Tabla de comparación conceptos x criterios
   */

  // Props
  let {
    title = "Compara y contrasta",
    instruction = "Completa la tabla comparativa",
    concepts = [], // ["Concepto A", "Concepto B"]
    criteria = [], // ["Criterio 1", "Criterio 2"]

    // Metadata educativa
    bloomLevel = "analizar",
    materia = "general",
    oaId = null,

    // Configuración
    allowMultipleAttempts = false,
    showFeedback = false,

    // Callbacks
    onAnswer = null,
    onComplete = null
  } = $props();

  // Estados locales
  let userTable = $state({}); // { "Concepto A": { "Criterio 1": "respuesta" } }
  let hasSubmitted = $state(false);
  let results = $state({});
  let attemptCount = $state(0);
  let containerRef = $state(null);

  // Inicializar tabla vacía
  $effect(() => {
    if (concepts.length > 0 && criteria.length > 0 && Object.keys(userTable).length === 0) {
      const initialTable = {};
      concepts.forEach(concept => {
        initialTable[concept] = {};
        criteria.forEach(criterion => {
          initialTable[concept][criterion] = '';
        });
      });
      userTable = initialTable;
    }
  });

  // Colores por nivel de Bloom
  const bloomColors = {
    recordar: { bg: 'bg-red-500/20', border: 'border-red-500', text: 'text-red-400' },
    comprender: { bg: 'bg-orange-500/20', border: 'border-orange-500', text: 'text-orange-400' },
    aplicar: { bg: 'bg-yellow-500/20', border: 'border-yellow-500', text: 'text-yellow-400' },
    analizar: { bg: 'bg-green-500/20', border: 'border-green-500', text: 'text-green-400' },
    evaluar: { bg: 'bg-blue-500/20', border: 'border-blue-500', text: 'text-blue-400' },
    crear: { bg: 'bg-purple-500/20', border: 'border-purple-500', text: 'text-purple-400' }
  };

  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.analizar;

  // Verificar si está completa
  const allFilled = $derived.by(() => {
    for (const concept of concepts) {
      for (const criterion of criteria) {
        if (!userTable[concept]?.[criterion]?.trim()) {
          return false;
        }
      }
    }
    return true;
  });

  // Funciones
  function handleSubmit() {
    if (!allFilled) {
      alert("Por favor, completa todas las celdas de la tabla.");
      return;
    }

    hasSubmitted = true;
    attemptCount++;

    // Callback para backend (enviamos la tabla completa)
    if (onAnswer) {
      onAnswer({
        oaId: oaId,
        bloomLevel: bloomLevel,
        materia: materia,
        userTable: userTable,
        isCorrect: true, // El backend validará
        attemptCount: attemptCount,
        timestamp: new Date().toISOString()
      });
    }
  }

  // Animación de entrada
  onMount(() => {
    if (containerRef) {
      gsap.from(containerRef, {
        opacity: 0,
        y: 20,
        duration: 0.5,
        ease: 'power2.out'
      });
    }
  });
</script>

<div
  bind:this={containerRef}
  class="w-full max-w-6xl mx-auto p-6 bg-canvas-950 rounded-2xl border border-slate-800 shadow-2xl"
>
  <!-- Header -->
  <div class="flex items-center justify-between mb-6">
    <div class="flex items-center gap-3">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase {currentBloomStyle.bg} {currentBloomStyle.border} {currentBloomStyle.text} border">
        {bloomLevel}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-semibold bg-slate-500/20 text-slate-400 border border-slate-500">
        {materia}
      </span>
    </div>
    {#if attemptCount > 0}
      <span class="text-xs text-slate-500">Intento {attemptCount}</span>
    {/if}
  </div>

  <!-- Título e instrucciones -->
  <div class="mb-6">
    <h3 class="text-xl font-semibold text-white mb-2">
      {title}
    </h3>
    <p class="text-slate-400 text-sm">
      {instruction}
    </p>
  </div>

  <!-- Tabla comparativa -->
  <div class="overflow-x-auto mb-6">
    <table class="w-full border-collapse">
      <thead>
        <tr>
          <th class="p-4 bg-canvas-900 border border-slate-700 text-left text-white font-semibold">
            Concepto / Criterio
          </th>
          {#each criteria as criterion}
            <th class="p-4 bg-canvas-900 border border-slate-700 text-left text-white font-semibold">
              {criterion}
            </th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each concepts as concept}
          <tr>
            <td class="p-4 bg-canvas-900/50 border border-slate-700 text-white font-semibold">
              {concept}
            </td>
            {#each criteria as criterion}
              <td class="p-2 border border-slate-700">
                <textarea
                  bind:value={userTable[concept][criterion]}
                  disabled={hasSubmitted && !allowMultipleAttempts}
                  placeholder="Escribe aquí..."
                  rows="2"
                  class="
                    w-full p-2 rounded-lg
                    bg-canvas-800 border border-slate-600
                    text-white text-sm
                    placeholder-slate-500
                    focus:outline-none focus:ring-2 focus:ring-cyan-500
                    disabled:opacity-50 disabled:cursor-not-allowed
                    resize-none
                  "
                ></textarea>
              </td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  <!-- Botones de acción -->
  <div class="flex gap-3">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={!allFilled}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-focus-500 to-blue-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-cyan-500/50
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Enviar Tabla
      </button>
    {/if}
  </div>
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
