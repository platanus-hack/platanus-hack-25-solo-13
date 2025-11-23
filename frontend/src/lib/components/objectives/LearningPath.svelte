<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import OACard from './OACard.svelte';

  // Props
  let {
    oas = [],
    bloomLevels = {},
    progressData = {},
    plansData = {},
    generatingPlanFor = null,
    diagnosticLevel = 0,
    onGeneratePlan = null,
    onViewPlan = null,
    onPractice = null
  } = $props();

  // Category configuration with colors and icons
  const categoryConfig: Record<string, { color: string; icon: string; description: string; order: number }> = {
    'Lectura': {
      color: 'from-blue-500 to-blue-600',
      icon: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25',
      description: 'Comprensión y análisis de textos',
      order: 1
    },
    'Escritura': {
      color: 'from-purple-500 to-purple-600',
      icon: 'M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10',
      description: 'Producción y creación de textos',
      order: 2
    },
    'Comunicación Oral': {
      color: 'from-green-500 to-green-600',
      icon: 'M12 18.75a6 6 0 006-6v-1.5m-6 7.5a6 6 0 01-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 01-3-3V4.5a3 3 0 116 0v8.25a3 3 0 01-3 3z',
      description: 'Expresión y diálogo oral',
      order: 3
    },
    'Investigación': {
      color: 'from-amber-500 to-amber-600',
      icon: 'M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z',
      description: 'Investigación y síntesis',
      order: 4
    },
    'General': {
      color: 'from-slate-500 to-slate-600',
      icon: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25',
      description: 'Objetivos generales',
      order: 5
    }
  };

  // Group OAs by category and sort by orden
  const groupedOAs = $derived(() => {
    const groups: Record<string, any[]> = {};

    oas.forEach(oa => {
      const category = oa.categoria || 'General';
      if (!groups[category]) {
        groups[category] = [];
      }
      groups[category].push(oa);
    });

    // Sort OAs within each category by orden field
    Object.keys(groups).forEach(category => {
      groups[category].sort((a, b) => (a.orden || 0) - (b.orden || 0));
    });

    return groups;
  });

  // Get ordered categories
  const orderedCategories = $derived(() => {
    return Object.keys(groupedOAs())
      .sort((a, b) => {
        const orderA = categoryConfig[a]?.order || 999;
        const orderB = categoryConfig[b]?.order || 999;
        return orderA - orderB;
      });
  });

  // Get stats for each category
  const categoryStats = $derived(() => {
    const stats: Record<string, { total: number; completed: number; inProgress: number }> = {};

    orderedCategories().forEach(category => {
      const oasInCategory = groupedOAs()[category] || [];
      const completed = oasInCategory.filter(oa => plansData[oa.id]?.isCompleted).length;
      const inProgress = oasInCategory.filter(oa =>
        plansData[oa.id]?.hasPlan && !plansData[oa.id]?.isCompleted
      ).length;

      stats[category] = {
        total: oasInCategory.length,
        completed,
        inProgress
      };
    });

    return stats;
  });

  let pathRef: HTMLDivElement | null = null;

  onMount(() => {
    // Animate path sections
    if (pathRef) {
      const sections = pathRef.querySelectorAll('.path-section');
      gsap.from(sections, {
        opacity: 0,
        x: -50,
        duration: 0.6,
        stagger: 0.1,
        ease: 'power2.out'
      });
    }
  });
</script>

<div bind:this={pathRef} class="relative">
  <!-- Path Container -->
  <div class="space-y-8">
    {#each orderedCategories() as category, index}
      {@const oasInCategory = groupedOAs()[category] || []}
      {@const stats = categoryStats()[category]}
      {@const config = categoryConfig[category] || categoryConfig['General']}

      {#if oasInCategory.length > 0}
        <div class="path-section relative">
          <!-- Category Header with connecting line -->
          <div class="flex items-center gap-4 mb-6">
            <!-- Category Badge -->
            <div class="relative flex-shrink-0">
              <!-- Connecting line to previous category (except first) -->
              {#if index > 0}
                <div class="absolute bottom-full left-1/2 -translate-x-1/2 w-0.5 h-8 bg-gradient-to-b from-canvas-700 to-transparent"></div>
              {/if}

              <div class="relative group">
                <!-- Animated glow -->
                <div class="absolute -inset-1 bg-gradient-to-br {config.color} rounded-2xl blur-md opacity-40 group-hover:opacity-60 transition-opacity"></div>

                <!-- Badge content -->
                <div class="relative bg-canvas-800 rounded-2xl p-4 border-2 border-canvas-700">
                  <div class="flex items-center gap-3">
                    <!-- Icon -->
                    <div class="w-12 h-12 rounded-xl bg-gradient-to-br {config.color} flex items-center justify-center">
                      <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d={config.icon} />
                      </svg>
                    </div>

                    <!-- Info -->
                    <div>
                      <div class="text-xs text-slate-400 font-medium">{config.description}</div>
                      <div class="text-lg font-bold text-white">{category}</div>
                      <div class="text-xs text-slate-500 mt-0.5">
                        {stats.completed}/{stats.total} completados
                        {#if stats.inProgress > 0}
                          · {stats.inProgress} en progreso
                        {/if}
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Connecting line to next category (except last) -->
              {#if index < orderedCategories().length - 1}
                <div class="absolute top-full left-1/2 -translate-x-1/2 w-0.5 h-8 bg-gradient-to-b from-canvas-700 to-transparent"></div>
              {/if}
            </div>

            <!-- Progress bar for this category -->
            <div class="flex-1">
              <div class="h-2 bg-canvas-900 rounded-full overflow-hidden border border-canvas-700">
                <div
                  class="h-full bg-gradient-to-r {config.color} transition-all duration-500"
                  style="width: {stats.total > 0 ? (stats.completed / stats.total * 100) : 0}%"
                ></div>
              </div>
            </div>
          </div>

          <!-- OAs in this category (ordered by orden field) -->
          <div class="ml-20 space-y-3">
            {#each oasInCategory as oa, oaIndex (oa.id)}
              <div class="relative">
                <!-- Step number indicator -->
                <div class="absolute -left-10 top-1/2 -translate-y-1/2 w-6 h-6 rounded-full bg-canvas-700 border-2 border-canvas-600 flex items-center justify-center">
                  <span class="text-xs font-bold text-slate-400">{oaIndex + 1}</span>
                </div>

                <OACard
                  {oa}
                  bloomLevel={bloomLevels[oa.id] || diagnosticLevel || 0}
                  progress={progressData[oa.id] || 0}
                  hasPlan={plansData[oa.id]?.hasPlan || false}
                  isCompleted={plansData[oa.id]?.isCompleted || false}
                  isGenerating={generatingPlanFor === oa.id}
                  {onGeneratePlan}
                  {onViewPlan}
                  {onPractice}
                />
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/each}
  </div>

  <!-- Empty state if no OAs -->
  {#if oas.length === 0}
    <div class="text-center py-16 px-6">
      <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-slate-800 flex items-center justify-center">
        <svg class="w-12 h-12 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>
      <h3 class="text-xl font-semibold text-slate-300 mb-2">
        No se encontraron objetivos de aprendizaje
      </h3>
      <p class="text-slate-500">
        Intenta cambiar los filtros de búsqueda o categoría
      </p>
    </div>
  {/if}
</div>
