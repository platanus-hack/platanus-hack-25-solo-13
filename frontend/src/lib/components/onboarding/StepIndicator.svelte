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
    <div class="h-2 bg-slate-800 rounded-full overflow-hidden">
      <div
        class="h-full bg-gradient-to-r from-indigo-500 to-cyan-500 transition-all duration-500 ease-out"
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
          class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold mb-2 transition-all duration-300 {
            isCompleted ? 'bg-gradient-to-br from-indigo-500 to-cyan-500 text-white scale-100' :
            isCurrent ? 'bg-indigo-600 text-white scale-110 shadow-lg shadow-indigo-500/50' :
            'bg-slate-800 text-slate-500 scale-90'
          }"
        >
          {#if isCompleted}
            <span>✓</span>
          {:else}
            <span>{stepNumber}</span>
          {/if}
        </div>

        <span class="text-xs text-center {isCurrent ? 'text-indigo-300 font-semibold' : 'text-slate-500'} hidden sm:block">
          {stepTitle}
        </span>
      </div>
    {/each}
  </div>
</div>
