<script lang="ts">
  import type { Subject } from '$lib/constants/subjects';
  import { getDomainLevelInfo } from '$lib/constants/subjects';

  interface Props {
    subject: Subject;
    domainLevel: number;
    onClick: () => void;
  }

  let { subject, domainLevel, onClick }: Props = $props();

  const levelInfo = $derived(getDomainLevelInfo(domainLevel));
  const percentage = $derived(domainLevel * 25); // 0->0%, 1->25%, 2->50%, 3->75%, 4->100%
</script>

<button
  onclick={onClick}
  class="w-full text-center p-8 rounded-2xl bg-canvas-900 border border-slate-700 flex flex-col items-center hover:bg-canvas-800 hover:border-slate-600 transition-all duration-200 hover:scale-[1.02] group"
>
  <!-- Icon -->
  <div class="h-28 w-28 rounded-2xl bg-gradient-to-br from-slate-700 to-slate-800 flex items-center justify-center text-6xl mb-6 group-hover:scale-105 transition-transform shadow-lg border border-slate-600">
    {subject.icon}
  </div>

  <!-- Title -->
  <h4 class="text-2xl font-bold text-white mb-3">
    {subject.name}
  </h4>

  <!-- Badge -->
  <span class="text-sm font-bold {levelInfo.badgeColor} px-4 py-1.5 rounded-full {levelInfo.textColor} mb-4">
    {levelInfo.label}
  </span>

  <!-- Progress bar -->
  <div class="w-full mb-4">
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

  <!-- Description -->
  <p class="text-sm text-slate-400 leading-relaxed">
    {subject.description}
  </p>
</button>
