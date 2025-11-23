<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    oa,
    progress = 0, // 0-100
    bloomLevel = 0, // 0-6 (0 = not evaluated)
    hasPlan = false,
    isCompleted = false, // Whether the learning plan has been completed
    isGenerating = false, // Loading state while generating plan
    categoryIcon = '📚',
    onGeneratePlan = null,
    onViewPlan = null,
    onPractice = null
  } = $props();

  // Local state
  let isExpanded = $state(false);

  // Category icons mapping (using SVG components)
  const categoryIconsSVG = {
    'Lectura': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`,
    'Escritura': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>`,
    'Comunicación Oral': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>`,
    'Investigación': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>`,
    'General': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`
  };

  // Bloom level labels and colors
  const bloomLevels = [
    { label: 'No evaluado', stars: 0, color: 'text-slate-400', bg: 'bg-slate-100', border: 'border-slate-300' },
    { label: 'Recordar', stars: 1, color: 'text-blue-600', bg: 'bg-blue-50', border: 'border-blue-300' },
    { label: 'Comprender', stars: 2, color: 'text-green-600', bg: 'bg-green-50', border: 'border-green-300' },
    { label: 'Aplicar', stars: 3, color: 'text-yellow-600', bg: 'bg-yellow-50', border: 'border-yellow-300' },
    { label: 'Analizar', stars: 4, color: 'text-orange-600', bg: 'bg-orange-50', border: 'border-orange-300' },
    { label: 'Evaluar', stars: 5, color: 'text-purple-600', bg: 'bg-purple-50', border: 'border-purple-300' },
    { label: 'Crear', stars: 6, color: 'text-pink-600', bg: 'bg-pink-50', border: 'border-pink-300' }
  ];

  const currentBloomLevel = bloomLevels[bloomLevel] || bloomLevels[0];

  // Toggle expand/collapse
  function toggleExpand() {
    isExpanded = !isExpanded;
  }

  // Get progress color
  function getProgressColor(percentage: number): string {
    if (percentage === 0) return 'bg-slate-200';
    if (percentage < 30) return 'bg-red-400';
    if (percentage < 70) return 'bg-yellow-400';
    return 'bg-gradient-to-r from-lumera-500 to-focus-500';
  }

  let cardRef: HTMLDivElement | null = null;

  onMount(() => {
    // Entrance animation will be handled by parent component (OAList) with stagger
  });
</script>

<div
  bind:this={cardRef}
  class="oa-card bg-white rounded-xl border-2 border-slate-200 overflow-hidden transition-all duration-300 hover:shadow-lg hover:border-slate-300 group"
>
  <!-- Card Header (Always Visible) -->
  <button
    onclick={toggleExpand}
    class="w-full px-6 py-4 flex items-start gap-4 text-left hover:bg-slate-50 transition-colors"
  >
    <!-- Category Icon -->
    <div class="flex-shrink-0 w-12 h-12 rounded-xl bg-gradient-to-br from-lumera-100 to-focus-100 border-2 {currentBloomLevel.border} flex items-center justify-center text-slate-700">
      {@html categoryIconsSVG[oa.categoria] || categoryIconsSVG['General']}
    </div>

    <!-- Content -->
    <div class="flex-1 min-w-0">
      <!-- Title Row -->
      <div class="flex items-start justify-between gap-3 mb-2">
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 mb-2">
            {#if isCompleted}
              <span class="px-2 py-0.5 rounded-full text-xs font-semibold bg-blue-100 text-blue-700 border border-blue-300 flex items-center gap-1">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
                Estudiado
              </span>
            {:else if hasPlan}
              <span class="px-2 py-0.5 rounded-full text-xs font-semibold bg-green-100 text-green-700 border border-green-300 flex items-center gap-1">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
                Plan Activo
              </span>
            {:else}
              <span class="px-2 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600 border border-slate-300">
                Sin Plan
              </span>
            {/if}
          </div>
          <h3 class="text-base font-semibold text-slate-900 leading-tight">
            {oa.titulo}
          </h3>
        </div>

        <!-- Expand Icon -->
        <div class="flex-shrink-0 w-6 h-6 rounded-full bg-slate-100 flex items-center justify-center text-slate-600 group-hover:bg-slate-200 transition-all {isExpanded ? 'rotate-180' : ''}">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>

      <!-- Progress Bar -->
      <div class="mb-3">
        <div class="w-full h-2 bg-slate-100 rounded-full overflow-hidden">
          <div
            class="{getProgressColor(progress)} h-full transition-all duration-500 rounded-full"
            style="width: {progress}%"
          ></div>
        </div>
        <p class="text-xs text-slate-500 mt-1">{progress}% completado</p>
      </div>

      <!-- Bloom Level Badge -->
      <div class="flex items-center gap-2">
        <div class="flex items-center gap-2 px-3 py-1 rounded-lg {currentBloomLevel.bg} border {currentBloomLevel.border}">
          {#if currentBloomLevel.stars > 0}
            <div class="flex gap-0.5">
              {#each Array(currentBloomLevel.stars) as _, i}
                <svg class="w-4 h-4 {currentBloomLevel.color}" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              {/each}
            </div>
          {/if}
          <span class="{currentBloomLevel.color} text-sm font-semibold">
            {#if currentBloomLevel.stars === 0}
              <span class="text-slate-500">No evaluado</span>
            {:else}
              Nivel {bloomLevel} - {currentBloomLevel.label}
            {/if}
          </span>
        </div>
      </div>
    </div>
  </button>

  <!-- Expanded Content -->
  {#if isExpanded}
    <div class="px-6 pb-6 border-t border-slate-100 pt-4 animate-in slide-in-from-top duration-200">
      <!-- Description -->
      <div class="mb-4">
        <h4 class="text-sm font-semibold text-slate-700 mb-2">Descripción:</h4>
        <div class="text-sm text-slate-600 leading-relaxed whitespace-pre-wrap">
          {oa.descripcion}
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-wrap items-center gap-3">
        {#if hasPlan && onViewPlan}
          <button
            onclick={() => onViewPlan?.(oa)}
            class="px-5 py-2.5 bg-gradient-to-r from-lumera-600 to-focus-600 hover:from-lumera-500 hover:to-focus-500 text-white font-semibold rounded-lg transition-all duration-300 hover:scale-105 hover:shadow-lg flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            {isCompleted ? 'Repasar Material' : 'Ver Plan Existente'}
          </button>
        {/if}

        {#if isCompleted && onPractice}
          <button
            onclick={() => onPractice?.(oa)}
            class="px-5 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-semibold rounded-lg transition-all duration-300 hover:scale-105 hover:shadow-lg flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
            </svg>
            Practicar
          </button>
        {/if}

        {#if !hasPlan && onGeneratePlan}
          <button
            onclick={() => onGeneratePlan?.(oa)}
            disabled={isGenerating}
            class="px-5 py-2.5 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-500 hover:to-emerald-500 text-white font-semibold rounded-lg transition-all duration-300 {isGenerating ? 'opacity-70 cursor-not-allowed' : 'hover:scale-105'} flex items-center gap-2"
          >
            {#if isGenerating}
              <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Generando...
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Generar Plan de Aprendizaje
            {/if}
          </button>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  @keyframes slide-in-from-top {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-in {
    animation: slide-in-from-top 0.2s ease-out;
  }
</style>
