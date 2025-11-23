<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import OACard from './OACard.svelte';

  // Props
  let {
    oas = [],
    bloomLevels = {},  // Map of oaId => bloomLevel
    progressData = {},  // Map of oaId => progress percentage
    plansData = {},  // Map of oaId => { hasPlan: boolean, isCompleted: boolean }
    generatingPlanFor = null,  // ID of OA currently generating plan
    onGeneratePlan = null,
    onViewPlan = null,
    onPractice = null
  } = $props();

  let listRef: HTMLDivElement | null = null;

  onMount(() => {
    // Animate cards with stagger effect
    if (listRef) {
      const cards = listRef.querySelectorAll('.oa-card');
      gsap.from(cards, {
        opacity: 0,
        y: 30,
        duration: 0.5,
        stagger: 0.05,
        ease: 'power2.out'
      });
    }
  });
</script>

<div bind:this={listRef} class="space-y-4">
  {#if oas.length === 0}
    <!-- Empty State -->
    <div class="text-center py-16 px-6">
      <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-slate-100 flex items-center justify-center">
        <svg class="w-12 h-12 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>
      <h3 class="text-xl font-semibold text-slate-700 mb-2">
        No se encontraron objetivos de aprendizaje
      </h3>
      <p class="text-slate-500">
        Intenta cambiar los filtros de búsqueda o categoría
      </p>
    </div>
  {:else}
    {#each oas as oa (oa.id)}
      <OACard
        {oa}
        bloomLevel={bloomLevels[oa.id] || 0}
        progress={progressData[oa.id] || 0}
        hasPlan={plansData[oa.id]?.hasPlan || false}
        isCompleted={plansData[oa.id]?.isCompleted || false}
        isGenerating={generatingPlanFor === oa.id}
        {onGeneratePlan}
        {onViewPlan}
        {onPractice}
      />
    {/each}
  {/if}
</div>
