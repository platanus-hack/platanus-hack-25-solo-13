<script lang="ts">
  import type { Subject } from '$lib/constants/subjects';
  import { getDomainLevelInfo } from '$lib/constants/subjects';

  interface Props {
    subject: Subject;
    domainLevel: number;
    onClick: () => void;
    comingSoon?: boolean;
  }

  let { subject, domainLevel, onClick, comingSoon = false }: Props = $props();

  const levelInfo = $derived(getDomainLevelInfo(domainLevel));
  const percentage = $derived(domainLevel * 25); // 0->0%, 1->25%, 2->50%, 3->75%, 4->100%
</script>

<button
  onclick={comingSoon ? undefined : onClick}
  disabled={comingSoon}
  class="relative w-full text-center p-8 rounded-2xl bg-canvas-900 border border-slate-700 flex flex-col items-center transition-all duration-200 {comingSoon ? 'opacity-60 cursor-not-allowed' : 'hover:bg-canvas-800 hover:border-slate-600 hover:scale-[1.02]'} group"
>
  <!-- Coming Soon Badge -->
  {#if comingSoon}
    <div class="absolute -top-3 -right-3 px-3 py-1 bg-purple-600 text-white text-xs font-bold rounded-full border-2 border-canvas-950 shadow-lg">
      Próximamente
    </div>
  {/if}

  <!-- Icon -->
  <div class="h-28 w-28 rounded-2xl bg-gradient-to-br from-slate-700 to-slate-800 flex items-center justify-center mb-6 {comingSoon ? '' : 'group-hover:scale-105'} transition-transform shadow-lg border border-slate-600">
    {#if subject.name === 'Lengua y Literatura'}
      <svg class="w-16 h-16 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
      </svg>
    {:else if subject.name === 'Matemáticas'}
      <svg class="w-16 h-16 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
      </svg>
    {:else}
      <span class="text-6xl">{subject.icon}</span>
    {/if}
  </div>

  <!-- Title -->
  <h4 class="text-2xl font-bold text-white mb-3">
    {subject.name}
  </h4>

  <!-- Badge -->
  {#if !comingSoon}
    <span class="text-sm font-bold {levelInfo.badgeColor} px-4 py-1.5 rounded-full {levelInfo.textColor} mb-4">
      {levelInfo.label}
    </span>
  {:else}
    <div class="h-[28px] mb-4"></div>
  {/if}

  <!-- Progress bar -->
  {#if !comingSoon}
    <div class="w-full">
      <div class="flex justify-between text-xs text-slate-400 mb-2">
        <span>Progreso</span>
        <span>{percentage}%</span>
      </div>
      <div class="h-3 w-full bg-black/40 rounded-full overflow-hidden border border-black/50">
        <div
          class="h-full bg-gradient-to-r from-teal-500 to-cyan-500 transition-all duration-500"
          style="width: {percentage}%"
        ></div>
      </div>
    </div>
  {:else}
    <div class="h-[44px] w-full"></div>
  {/if}
</button>
