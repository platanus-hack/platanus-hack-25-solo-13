<script lang="ts">
  interface ScaleOption {
    value: string;
    iconPath: string;
    label: string;
  }

  interface Props {
    options: ScaleOption[];
    selected: string;
    onSelect: (value: string) => void;
  }

  let { options, selected = $bindable(''), onSelect }: Props = $props();

  function handleSelect(value: string) {
    selected = value;
    onSelect(value);
  }
</script>

<div class="flex justify-center items-center gap-2 sm:gap-4 flex-wrap">
  {#each options as option}
    {@const isSelected = selected === option.value}
    <button
      type="button"
      class="flex flex-col items-center p-4 rounded-2xl border-2 transition-all duration-300 min-w-[80px] {
        isSelected
          ? 'bg-[#E1E1E1] border-[#E1E1E1] scale-110 shadow-xl'
          : 'bg-canvas-900/40 border-canvas-700 hover:border-canvas-600 hover:bg-canvas-800/60 hover:scale-105'
      }"
      onclick={() => handleSelect(option.value)}
    >
      <div class="flex items-center justify-center w-12 h-12 mb-2 rounded-xl transition-all duration-300 {
        isSelected ? 'bg-canvas-900/10 scale-110' : 'bg-canvas-800/60'
      }">
        <svg
          class="w-8 h-8 transition-colors duration-300 {
            isSelected ? 'text-canvas-900' : 'text-slate-400'
          }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d={option.iconPath} />
        </svg>
      </div>
      <div class="text-xs font-medium text-center {isSelected ? 'text-canvas-900' : 'text-slate-400'}">
        {option.label}
      </div>
    </button>
  {/each}
</div>
