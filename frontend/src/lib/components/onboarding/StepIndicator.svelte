<script lang="ts">
  interface Props {
    currentStep: number;
    totalSteps: number;
    stepTitles?: string[];
  }

  let { currentStep, totalSteps, stepTitles = [] }: Props = $props();

  const percentage = $derived((currentStep / totalSteps) * 100);
</script>

<div class="w-full mb-8">
  <!-- Progress bar -->
  <div class="relative mb-4">
    <div class="h-3 w-full bg-black/40 rounded-full overflow-hidden border border-black/50">
      <div
        class="h-full bg-gradient-to-r from-teal-500 to-cyan-500 transition-all duration-500 ease-out"
        style="width: {percentage}%"
      ></div>
    </div>
  </div>

  <!-- Step indicators -->
  <div class="flex justify-between items-center">
    {#each Array(totalSteps) as _, index}
      {@const stepNumber = index + 1}
      {@const isCompleted = stepNumber < currentStep}
      {@const isCurrent = stepNumber === currentStep}
      {@const stepTitle = stepTitles[index] || `Paso ${stepNumber}`}

      <div class="flex flex-col items-center flex-1">
        <div
          class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold mb-2 transition-all duration-300 {
            isCompleted ? 'bg-gradient-to-br from-slate-700 to-slate-800 text-white scale-100 border border-slate-600 shadow-lg' :
            isCurrent ? 'bg-[#E1E1E1] text-canvas-900 scale-110 border-2 border-[#E1E1E1] shadow-lg shadow-[#E1E1E1]/30' :
            'bg-black/40 text-slate-500 scale-90 border border-black/50'
          }"
        >
          {#if isCompleted}
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          {:else}
            <span>{stepNumber}</span>
          {/if}
        </div>

        <span class="text-xs text-center {isCurrent ? 'text-[#E1E1E1] font-semibold' : 'text-slate-500'} hidden sm:block">
          {stepTitle}
        </span>
      </div>
    {/each}
  </div>
</div>
