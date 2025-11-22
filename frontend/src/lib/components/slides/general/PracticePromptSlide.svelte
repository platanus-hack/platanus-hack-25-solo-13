<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    mensaje = "¡Hora de practicar lo aprendido!",
    submensaje = "Resuelve algunos ejercicios para reforzar tus conocimientos",
    icono = "🎯",
    previewEjercicios = [],
    motivacion = null,
    botonTexto = "Comenzar Práctica",
    colorTema = "green",
    materia = "general",
    mostrarConfetti = true,
    onNext = null,
    onPrevious = null,
    showNavigation = false // Generalmente este slide no tiene navegación estándar
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let confettiContainer = $state(null);

  // Contar total de ejercicios
  const totalEjercicios = $derived(
    previewEjercicios.reduce((sum, preview) => sum + (preview.cantidad || 0), 0)
  );

  // Animación de entrada con confetti
  onMount(() => {
    if (containerRef) {
      gsap.from(containerRef, {
        opacity: 0,
        scale: 0.9,
        duration: 0.6,
        ease: 'back.out(1.7)'
      });
    }

    // Confetti animation (opcional)
    if (mostrarConfetti && confettiContainer) {
      // Crear partículas de confetti
      for (let i = 0; i < 30; i++) {
        const confetti = document.createElement('div');
        confetti.className = 'confetti-piece';
        confetti.style.left = `${Math.random() * 100}%`;
        confetti.style.animationDelay = `${Math.random() * 0.5}s`;
        confetti.style.backgroundColor = ['#3b82f6', '#8b5cf6', '#ec4899', '#10b981'][Math.floor(Math.random() * 4)];
        confettiContainer.appendChild(confetti);
      }
    }
  });
</script>

<div
  bind:this={containerRef}
  class="relative w-full max-w-4xl mx-auto p-8 bg-slate-950 rounded-2xl border border-slate-800 shadow-2xl overflow-hidden"
>
  <!-- Confetti container -->
  {#if mostrarConfetti}
    <div bind:this={confettiContainer} class="confetti-container"></div>
  {/if}

  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-{colorTema}-500/20 text-{colorTema}-400 border border-{colorTema}-500">
        {materia}
      </span>
      <div class="flex-1 h-px bg-slate-800"></div>
    </div>
  </div>

  <!-- Contenido principal -->
  <div class="text-center mb-12">
    <!-- Icono grande -->
    <div class="mb-6 inline-block p-8 bg-gradient-to-br from-{colorTema}-500/20 to-{colorTema}-600/20 rounded-full border-4 border-{colorTema}-500/30 animate-pulse">
      <span class="text-7xl">{icono}</span>
    </div>

    <!-- Mensaje principal -->
    <h2 class="text-4xl font-bold text-white mb-4">
      {mensaje}
    </h2>

    <!-- Submensaje -->
    <p class="text-xl text-slate-300 mb-8 max-w-2xl mx-auto">
      {submensaje}
    </p>

    <!-- Mensaje de motivación -->
    {#if motivacion}
      <div class="inline-block px-6 py-3 bg-purple-500/20 rounded-xl border border-purple-500/50 mb-8">
        <p class="text-purple-300 font-semibold">
          ✨ {motivacion}
        </p>
      </div>
    {/if}
  </div>

  <!-- Preview de ejercicios -->
  {#if previewEjercicios.length > 0}
    <div class="mb-12">
      <h3 class="text-lg font-semibold text-slate-300 mb-4 text-center">
        Qué practicarás:
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-{Math.min(previewEjercicios.length, 3)} gap-4">
        {#each previewEjercicios as preview}
          <div class="p-6 bg-slate-900/50 rounded-2xl border border-slate-700 text-center">
            <div class="text-4xl mb-3">
              {preview.icono || '📝'}
            </div>
            <p class="font-semibold text-white mb-2">
              {preview.tipo}
            </p>
            <p class="text-3xl font-bold text-{colorTema}-400 mb-1">
              {preview.cantidad}
            </p>
            <p class="text-xs text-slate-500 uppercase">
              {preview.cantidad === 1 ? 'ejercicio' : 'ejercicios'}
            </p>
          </div>
        {/each}
      </div>

      {#if totalEjercicios > 0}
        <div class="mt-6 text-center">
          <p class="text-slate-400">
            <span class="text-2xl font-bold text-{colorTema}-400">{totalEjercicios}</span>
            <span class="text-sm"> ejercicios en total</span>
          </p>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Botón de acción -->
  <div class="text-center mb-8">
    <button
      onclick={onNext}
      class="
        group relative inline-flex items-center gap-3
        px-12 py-5 rounded-2xl
        text-xl font-bold text-white
        bg-gradient-to-r from-{colorTema}-500 to-{colorTema}-600
        border-2 border-{colorTema}-400
        transition-all duration-300
        hover:scale-110 hover:shadow-2xl hover:shadow-{colorTema}-500/50
        active:scale-95
      "
    >
      <span>{botonTexto}</span>
      <span class="text-2xl group-hover:translate-x-1 transition-transform duration-300">
        →
      </span>

      <!-- Glow effect -->
      <div class="absolute inset-0 rounded-2xl bg-white/20 opacity-0 group-hover:opacity-100 transition-opacity duration-300 blur-xl -z-10"></div>
    </button>
  </div>

  <!-- Navegación estándar (opcional) -->
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
          Listo para practicar
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-{colorTema}-500 to-{colorTema}-600
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-{colorTema}-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        Siguiente →
      </button>
    </div>
  {/if}
</div>

<style>
  .confetti-container {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    overflow: hidden;
    z-index: 0;
  }

  :global(.confetti-piece) {
    position: absolute;
    width: 10px;
    height: 10px;
    top: -10px;
    opacity: 0;
    animation: confettiFall 3s ease-in forwards;
  }

  @keyframes confettiFall {
    0% {
      transform: translateY(0) rotate(0deg);
      opacity: 1;
    }
    100% {
      transform: translateY(100vh) rotate(720deg);
      opacity: 0;
    }
  }
</style>
