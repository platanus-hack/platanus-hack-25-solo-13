<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    components = [],
    currentIndex = 0,
    onNavigate = null,
    planTitle = '',
    isCollapsed = $bindable(false)
  } = $props();

  // Component type configuration
  const typeConfig: Record<string, { label: string; icon: string; color: string; bgColor: string }> = {
    'teach': {
      label: 'Aprender',
      icon: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25',
      color: 'text-blue-400',
      bgColor: 'bg-blue-500/10'
    },
    'practice': {
      label: 'Practicar',
      icon: 'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10',
      color: 'text-purple-400',
      bgColor: 'bg-purple-500/10'
    },
    'assess': {
      label: 'Evaluar',
      icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
      color: 'text-amber-400',
      bgColor: 'bg-amber-500/10'
    }
  };

  // State
  let navRef: HTMLDivElement | null = null;

  // Calculate progress by type
  const progressByType = $derived(() => {
    const stats: Record<string, { current: number; total: number }> = {
      teach: { current: 0, total: 0 },
      practice: { current: 0, total: 0 },
      assess: { current: 0, total: 0 }
    };

    components.forEach((comp, idx) => {
      const type = comp.tipo_componente || 'teach';
      if (stats[type]) {
        stats[type].total++;
        if (idx < currentIndex) {
          stats[type].current++;
        }
      }
    });

    return stats;
  });

  // Get status for a component
  function getComponentStatus(index: number) {
    if (index < currentIndex) return 'completed';
    if (index === currentIndex) return 'current';
    return 'upcoming';
  }

  // Check if component is navigable
  function isNavigable(index: number) {
    return index <= currentIndex; // Can navigate to current or completed
  }

  // Handle component click
  function handleComponentClick(index: number) {
    if (isNavigable(index) && onNavigate) {
      onNavigate(index);
    }
  }

  // Toggle collapse
  function toggleCollapse() {
    isCollapsed = !isCollapsed;
  }

  onMount(() => {
    if (navRef) {
      gsap.from(navRef, {
        x: -300,
        opacity: 0,
        duration: 0.5,
        ease: 'power2.out'
      });
    }
  });
</script>

<div
  bind:this={navRef}
  class="fixed left-0 top-16 h-[calc(100vh-4rem)] bg-canvas-900/95 backdrop-blur-sm border-r border-canvas-700 transition-all duration-300 z-40 {
    isCollapsed ? 'w-16' : 'w-80'
  }"
>
  <!-- Toggle Button -->
  <button
    onclick={toggleCollapse}
    class="absolute -right-3 top-6 w-6 h-6 rounded-full bg-canvas-800 border-2 border-canvas-700 flex items-center justify-center hover:bg-canvas-700 transition-colors z-50"
  >
    <svg
      class="w-3 h-3 text-slate-400 transition-transform duration-300 {isCollapsed ? 'rotate-180' : ''}"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
    </svg>
  </button>

  {#if !isCollapsed}
    <div class="h-full flex flex-col">
      <!-- Header -->
      <div class="p-4 border-b border-canvas-800">
        <h3 class="text-sm font-bold text-white mb-3 truncate">{planTitle}</h3>

        <!-- Progress by type -->
        <div class="space-y-2">
          {#each Object.entries(progressByType()) as [type, stats]}
            {@const config = typeConfig[type]}
            {#if stats.total > 0}
              <div class="flex items-center gap-2">
                <div class="w-5 h-5 rounded {config.bgColor} flex items-center justify-center flex-shrink-0">
                  <svg class="w-3 h-3 {config.color}" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d={config.icon} />
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between text-xs mb-1">
                    <span class="text-slate-400 truncate">{config.label}</span>
                    <span class="text-slate-500 flex-shrink-0">{stats.current}/{stats.total}</span>
                  </div>
                  <div class="h-1 bg-canvas-800 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-gradient-to-r from-{config.color.replace('text-', '')}/80 to-{config.color.replace('text-', '')} transition-all duration-500"
                      style="width: {stats.total > 0 ? (stats.current / stats.total * 100) : 0}%"
                    ></div>
                  </div>
                </div>
              </div>
            {/if}
          {/each}
        </div>
      </div>

      <!-- Components List -->
      <div class="flex-1 overflow-y-auto p-4 space-y-2">
        {#each components as component, index}
          {@const status = getComponentStatus(index)}
          {@const config = typeConfig[component.tipo_componente] || typeConfig['teach']}
          {@const navigable = isNavigable(index)}

          <button
            onclick={() => handleComponentClick(index)}
            disabled={!navigable}
            class="w-full text-left relative group {
              navigable ? 'cursor-pointer' : 'cursor-not-allowed opacity-50'
            }"
          >
            <!-- Connecting line (except for last item) -->
            {#if index < components.length - 1}
              <div class="absolute left-[13px] top-7 bottom-0 w-0.5 {
                status === 'completed' ? 'bg-lumera-500/30' : 'bg-canvas-700'
              }"></div>
            {/if}

            <div class="relative flex items-start gap-3 p-2 rounded-lg transition-colors {
              status === 'current'
                ? 'bg-lumera-500/10 border-2 border-lumera-500/50'
                : navigable
                ? 'hover:bg-canvas-800 border-2 border-transparent'
                : 'border-2 border-transparent'
            }">
              <!-- Status indicator -->
              <div class="relative flex-shrink-0 mt-0.5">
                {#if status === 'completed'}
                  <div class="w-7 h-7 rounded-full bg-lumera-500/20 border-2 border-lumera-500 flex items-center justify-center">
                    <svg class="w-4 h-4 text-lumera-400" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                {:else if status === 'current'}
                  <div class="w-7 h-7 rounded-full bg-lumera-500 border-2 border-lumera-400 flex items-center justify-center animate-pulse">
                    <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" />
                    </svg>
                  </div>
                {:else}
                  <div class="w-7 h-7 rounded-full bg-canvas-800 border-2 border-canvas-700 flex items-center justify-center">
                    <svg class="w-3 h-3 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                  </div>
                {/if}
              </div>

              <!-- Content -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-xs font-medium {
                    status === 'current' ? 'text-lumera-400' : status === 'completed' ? 'text-slate-300' : 'text-slate-500'
                  }">Paso {index + 1}</span>
                  <span class="px-2 py-0.5 rounded text-[10px] font-semibold {config.bgColor} {config.color}">
                    {config.label}
                  </span>
                </div>

                {#if component.titulo}
                  <p class="text-xs {
                    status === 'current' ? 'text-white font-medium' : status === 'completed' ? 'text-slate-400' : 'text-slate-600'
                  } line-clamp-2">
                    {component.titulo}
                  </p>
                {/if}
              </div>

              <!-- Current indicator -->
              {#if status === 'current'}
                <div class="absolute -right-1 top-1/2 -translate-y-1/2 w-1 h-8 bg-lumera-500 rounded-l"></div>
              {/if}
            </div>
          </button>
        {/each}
      </div>
    </div>
  {:else}
    <!-- Collapsed view -->
    <div class="h-full flex flex-col items-center py-6 gap-4">
      <div class="text-xs font-bold text-slate-600 writing-mode-vertical transform rotate-180">
        {Math.round((currentIndex / components.length) * 100)}%
      </div>

      <!-- Mini progress indicators -->
      <div class="flex-1 flex flex-col gap-1 overflow-y-auto py-2">
        {#each components as component, index}
          {@const status = getComponentStatus(index)}
          <div
            class="w-2 h-2 rounded-full {
              status === 'completed'
                ? 'bg-lumera-500'
                : status === 'current'
                ? 'bg-lumera-400 animate-pulse'
                : 'bg-canvas-700'
            }"
          ></div>
        {/each}
      </div>
    </div>
  {/if}
</div>
