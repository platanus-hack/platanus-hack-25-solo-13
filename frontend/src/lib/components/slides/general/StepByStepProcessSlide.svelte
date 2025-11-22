<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    titulo = "Proceso Paso a Paso",
    pasos = [],
    materia = "general",
    requiereConfirmacion = false,
    mostrarProgreso = true,
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let pasoActual = $state(0);
  let pasosConfirmados = $state(new Set());

  const pasoActivo = $derived(pasos[pasoActual] || null);
  const esUltimoPaso = $derived(pasoActual === pasos.length - 1);
  const esPrimerPaso = $derived(pasoActual === 0);
  const progresoPercentage = $derived(((pasoActual + 1) / pasos.length) * 100);

  // Verificar si puede avanzar
  const puedeAvanzar = $derived.by(() => {
    if (!requiereConfirmacion) return true;
    return pasosConfirmados.has(pasoActual);
  });

  function handleAnteriorPaso() {
    if (pasoActual > 0) {
      pasoActual--;
      animarCambioPaso();
    }
  }

  function handleSiguientePaso() {
    if (pasoActual < pasos.length - 1 && puedeAvanzar) {
      pasoActual++;
      animarCambioPaso();
    }
  }

  function animarCambioPaso() {
    const stepEl = document.querySelector('.step-content');
    if (stepEl) {
      gsap.fromTo(stepEl,
        { opacity: 0, x: 20 },
        { opacity: 1, x: 0, duration: 0.4, ease: 'power2.out' }
      );
    }
  }

  function toggleConfirmacion(index) {
    const newSet = new Set(pasosConfirmados);
    if (newSet.has(index)) {
      newSet.delete(index);
    } else {
      newSet.add(index);
    }
    pasosConfirmados = newSet;
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
  class="w-full max-w-5xl mx-auto p-8 bg-slate-950 rounded-2xl border border-slate-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-purple-500/20 text-purple-400 border border-purple-500">
        {materia}
      </span>
      <div class="flex-1 h-px bg-slate-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-6">
      {titulo}
    </h2>

    <!-- Barra de progreso -->
    {#if mostrarProgreso}
      <div class="mb-6">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-semibold text-slate-400">
            Paso {pasoActual + 1} de {pasos.length}
          </span>
          <span class="text-sm font-semibold text-purple-400">
            {progresoPercentage.toFixed(0)}%
          </span>
        </div>
        <div class="w-full h-2 bg-slate-800 rounded-full overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-purple-500 to-pink-500 transition-all duration-500 ease-out"
            style="width: {progresoPercentage}%"
          ></div>
        </div>
      </div>
    {/if}

    <!-- Navegador de pasos -->
    <div class="flex gap-2 overflow-x-auto pb-2">
      {#each pasos as paso, index}
        {@const isPasoActual = index === pasoActual}
        {@const isPasoCompletado = index < pasoActual || (requiereConfirmacion && pasosConfirmados.has(index))}

        <button
          onclick={() => {
            if (!requiereConfirmacion || index <= pasoActual) {
              pasoActual = index;
              animarCambioPaso();
            }
          }}
          class="
            flex items-center gap-2 px-3 py-2 rounded-xl text-sm font-semibold whitespace-nowrap
            transition-all duration-300
            {isPasoActual ? 'bg-purple-500 text-white scale-105' : ''}
            {!isPasoActual && isPasoCompletado ? 'bg-green-500/20 text-green-400 border border-green-500' : ''}
            {!isPasoActual && !isPasoCompletado ? 'bg-slate-800 text-slate-400 border border-slate-700' : ''}
            {!requiereConfirmacion || index <= pasoActual ? 'cursor-pointer hover:scale-105' : 'cursor-not-allowed opacity-50'}
          "
        >
          <span class="flex items-center justify-center w-6 h-6 rounded-full bg-white/20">
            {isPasoCompletado ? '✓' : index + 1}
          </span>
          <span class="hidden md:inline">Paso {index + 1}</span>
        </button>
      {/each}
    </div>
  </div>

  <!-- Contenido del paso actual -->
  {#if pasoActivo}
    <div class="step-content mb-8">
      <div class="p-6 bg-slate-900/50 rounded-2xl border-2 border-purple-500/30">
        <!-- Número de paso y título -->
        <div class="flex items-center gap-4 mb-4">
          <div class="flex items-center justify-center w-12 h-12 rounded-xl bg-purple-500 text-white text-xl font-bold">
            {pasoActual + 1}
          </div>
          <h3 class="text-xl font-bold text-white">
            {pasoActivo.titulo}
          </h3>
        </div>

        <!-- Contenido del paso -->
        <div class="mb-4 text-slate-200 leading-relaxed">
          {pasoActivo.contenido}
        </div>

        <!-- Ejemplo -->
        {#if pasoActivo.ejemplo}
          <div class="p-4 bg-slate-950 rounded-xl border border-slate-700">
            <p class="text-xs font-semibold text-purple-400 mb-2 uppercase">
              💡 Ejemplo:
            </p>
            <p class="text-slate-300 font-mono text-sm">
              {pasoActivo.ejemplo}
            </p>
          </div>
        {/if}

        <!-- Ayuda visual -->
        {#if pasoActivo.ayudaVisual}
          <div class="mt-4">
            <img
              src={pasoActivo.ayudaVisual}
              alt="Ayuda visual paso {pasoActual + 1}"
              class="w-full max-w-md mx-auto rounded-xl border border-slate-700"
            />
          </div>
        {/if}

        <!-- Checkbox de confirmación -->
        {#if requiereConfirmacion}
          <div class="mt-6 pt-4 border-t border-slate-700">
            <label class="flex items-center gap-3 cursor-pointer group">
              <input
                type="checkbox"
                checked={pasosConfirmados.has(pasoActual)}
                onchange={() => toggleConfirmacion(pasoActual)}
                class="w-5 h-5 rounded border-2 border-purple-500 bg-slate-900 checked:bg-purple-500 cursor-pointer"
              />
              <span class="text-sm text-slate-300 group-hover:text-purple-300">
                ✓ Entendí este paso y puedo continuar
              </span>
            </label>
          </div>
        {/if}
      </div>

      <!-- Navegación interna de pasos -->
      <div class="flex gap-3 mt-6">
        <button
          onclick={handleAnteriorPaso}
          disabled={esPrimerPaso}
          class="
            flex-1 px-4 py-3 rounded-xl font-semibold
            bg-slate-800 text-slate-300
            border border-slate-700
            transition-all duration-300
            hover:bg-slate-700
            disabled:opacity-30 disabled:cursor-not-allowed
          "
        >
          ← Paso Anterior
        </button>

        <button
          onclick={handleSiguientePaso}
          disabled={esUltimoPaso || !puedeAvanzar}
          class="
            flex-1 px-4 py-3 rounded-xl font-semibold
            bg-purple-500 text-white
            transition-all duration-300
            hover:bg-purple-600 hover:scale-105
            disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
          "
        >
          {#if requiereConfirmacion && !pasosConfirmados.has(pasoActual)}
            Marca "Entendí" para continuar
          {:else if esUltimoPaso}
            Último paso completado ✓
          {:else}
            Siguiente Paso →
          {/if}
        </button>
      </div>
    </div>
  {/if}

  <!-- Navegación principal del slide -->
  {#if showNavigation}
    <div class="flex items-center justify-between pt-6 border-t border-slate-800">
      <button
        onclick={onPrevious}
        disabled={!onPrevious}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-slate-800 text-slate-300
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
          {esUltimoPaso ? 'Todos los pasos completados' : 'Completa todos los pasos'}
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-purple-500 to-pink-500
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

<style>
  /* Estilos adicionales si es necesario */
</style>
