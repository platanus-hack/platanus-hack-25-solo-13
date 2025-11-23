<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    materiaName = "Materia",
    totalOAs = 0,
    completedOAs = 0,
    bloomLevel = 0, // Student's diagnostic level (0-6)
    categoryStats = [],  // Array of { categoria, completed, total }
    recommendedOA = null,  // Next recommended OA object
    onStartRecommended = null
  } = $props();

  const progressPercentage = totalOAs > 0 ? Math.round((completedOAs / totalOAs) * 100) : 0;

  // Bloom level labels
  const bloomLevels = [
    { label: 'No evaluado', stars: 0, color: 'text-slate-600', bg: 'bg-slate-100', gradient: 'from-slate-400 to-slate-500' },
    { label: 'Recordar', stars: 1, color: 'text-blue-600', bg: 'bg-blue-50', gradient: 'from-blue-500 to-blue-600' },
    { label: 'Comprender', stars: 2, color: 'text-green-600', bg: 'bg-green-50', gradient: 'from-green-500 to-green-600' },
    { label: 'Aplicar', stars: 3, color: 'text-yellow-600', bg: 'bg-yellow-50', gradient: 'from-yellow-500 to-yellow-600' },
    { label: 'Analizar', stars: 4, color: 'text-orange-600', bg: 'bg-orange-50', gradient: 'from-orange-500 to-orange-600' },
    { label: 'Evaluar', stars: 5, color: 'text-purple-600', bg: 'bg-purple-50', gradient: 'from-purple-500 to-purple-600' },
    { label: 'Crear', stars: 6, color: 'text-pink-600', bg: 'bg-pink-50', gradient: 'from-pink-500 to-pink-600' }
  ];

  const currentBloom = bloomLevels[bloomLevel] || bloomLevels[0];

  // Category icons (SVG)
  const categoryIconsSVG = {
    'Lectura': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`,
    'Escritura': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>`,
    'Comunicación Oral': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>`,
    'Investigación': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>`,
    'General': `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`
  };

  let dashboardRef: HTMLDivElement | null = null;

  onMount(() => {
    if (dashboardRef) {
      gsap.from('.stats-panel', {
        opacity: 0,
        y: -20,
        duration: 0.6,
        ease: 'power2.out'
      });

      // Wait for next tick to ensure category-stat elements are rendered
      setTimeout(() => {
        const categoryElements = dashboardRef?.querySelectorAll('.category-stat');
        if (categoryElements && categoryElements.length > 0) {
          gsap.fromTo(categoryElements,
            {
              opacity: 0,
              scale: 0.9
            },
            {
              opacity: 1,
              scale: 1,
              duration: 0.4,
              stagger: 0.1,
              delay: 0.2,
              ease: 'back.out(1.7)'
            }
          );
        }
      }, 100);

      if (recommendedOA) {
        gsap.from('.recommended-card', {
          opacity: 0,
          scale: 0.95,
          duration: 0.5,
          delay: 0.5,
          ease: 'back.out(1.7)'
        });
      }
    }
  });
</script>

<div bind:this={dashboardRef} class="space-y-5 mb-6">
  <!-- Main Stats Panel -->
  <div class="stats-panel bg-white rounded-2xl border-2 border-slate-200 p-5 shadow-lg">
    <!-- Header - Compact -->
    <div class="flex items-center justify-between mb-4 pb-3 border-b border-slate-200">
      <div>
        <h2 class="text-lg font-bold text-slate-900">
          {materiaName}
        </h2>
        <p class="text-xs text-slate-500">
          Objetivos de Aprendizaje
        </p>
      </div>
      <div>
        <svg class="w-8 h-8 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </div>
    </div>

    <!-- Two-column layout: Progress + Diagnostic Level -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Progress Section -->
      <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl p-4 border border-slate-200">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-lumera-500 to-focus-500 flex items-center justify-center text-white shadow-lg">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="flex-1">
            <p class="text-xs font-medium text-slate-600">Progreso General</p>
            <p class="text-xl font-bold text-slate-900">{completedOAs}<span class="text-sm text-slate-600">/{totalOAs}</span></p>
          </div>
        </div>
        <div class="w-full h-2.5 bg-white rounded-full overflow-hidden shadow-inner">
          <div
            class="h-full bg-gradient-to-r from-lumera-500 to-focus-500 rounded-full transition-all duration-1000"
            style="width: {progressPercentage}%"
          ></div>
        </div>
        <p class="text-right text-xs font-bold text-slate-700 mt-1.5">{progressPercentage}% completado</p>
      </div>

      <!-- Diagnostic Level Section -->
      <div class="bg-gradient-to-br {currentBloom.bg} rounded-xl p-4 border-2 border-{currentBloom.color.split('-')[1]}-300 relative overflow-hidden">
        <!-- Decorative background -->
        <div class="absolute top-0 right-0 opacity-10">
          <svg class="w-16 h-16 {currentBloom.color}" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
        </div>
        <div class="relative z-10">
          <div class="flex items-center gap-2 mb-3">
            <div class="flex gap-0.5">
              {#each Array(Math.min(currentBloom.stars || 0, 3)) as _, i}
                <svg class="w-6 h-6 {currentBloom.color}" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              {/each}
            </div>
            <div class="flex-1">
              <p class="text-xs font-medium text-slate-600">Tu Nivel Diagnóstico</p>
              <p class="{currentBloom.color} text-lg font-bold leading-tight">
                {currentBloom.stars > 0 ? `Nivel ${bloomLevel}` : 'Sin evaluar'}
              </p>
            </div>
          </div>
          <div class="bg-white/60 backdrop-blur-sm rounded-lg px-3 py-1.5 border border-white/40">
            <p class="{currentBloom.color} text-xs font-semibold">
              {currentBloom.label}
            </p>
            <p class="text-xs text-slate-600 mt-0.5">
              {#if currentBloom.stars === 0}
                Completa una evaluación diagnóstica
              {:else if currentBloom.stars <= 2}
                Sigue practicando
              {:else if currentBloom.stars <= 4}
                Muy bien
              {:else}
                Excelente
              {/if}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Category Stats Grid -->
    {#if categoryStats.length > 0}
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mt-4">
        {#each categoryStats as stat}
          <div class="category-stat bg-slate-50 rounded-xl p-3 border border-slate-200 hover:border-lumera-300 hover:shadow-md transition-all">
            <div class="mb-1.5 flex justify-center text-slate-600">
              {@html categoryIconsSVG[stat.categoria] || categoryIconsSVG['General']}
            </div>
            <p class="text-xs font-medium text-slate-600 mb-1 text-center truncate" title={stat.categoria}>
              {stat.categoria}
            </p>
            <p class="text-lg font-bold text-slate-900 text-center">
              {stat.completed}/{stat.total}
            </p>
            <div class="w-full h-1 bg-slate-200 rounded-full mt-1.5 overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-lumera-400 to-focus-400 rounded-full"
                style="width: {stat.total > 0 ? (stat.completed / stat.total) * 100 : 0}%"
              ></div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Recommended OA Card -->
  {#if recommendedOA}
    <div class="recommended-card relative">
      <!-- Animated glow effect -->
      <div class="absolute -inset-1 bg-gradient-to-r from-lumera-400 via-focus-400 to-purple-400 rounded-2xl blur-lg opacity-30 animate-pulse"></div>

      <div class="relative bg-white rounded-2xl border-2 border-lumera-300 p-6 shadow-lg">
        <div class="flex items-start gap-4">
          <div class="flex-shrink-0 w-14 h-14 rounded-xl bg-gradient-to-br from-lumera-500 to-focus-500 flex items-center justify-center text-white shadow-lg">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
            </svg>
          </div>
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <h3 class="text-lg font-bold text-slate-900">Recomendado Hoy</h3>
              <span class="px-2 py-1 rounded-full text-xs font-semibold bg-green-100 text-green-700 border border-green-300">
                Sugerido
              </span>
            </div>
            <p class="text-sm font-semibold text-slate-900 mb-1">
              {recommendedOA.titulo}
            </p>
            <p class="text-xs text-slate-500">
              Basado en tu nivel actual de dominio
            </p>
          </div>
          {#if onStartRecommended}
            <button
              onclick={() => onStartRecommended?.(recommendedOA)}
              class="flex-shrink-0 px-5 py-2.5 bg-gradient-to-r from-lumera-600 to-focus-600 hover:from-lumera-500 hover:to-focus-500 text-white font-semibold rounded-lg transition-all duration-300 hover:scale-105 hover:shadow-xl flex items-center gap-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Comenzar
            </button>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
