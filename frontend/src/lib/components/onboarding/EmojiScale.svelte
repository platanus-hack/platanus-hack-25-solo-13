<script lang="ts">
  interface ScaleOption {
    value: string;
    emoji: string;
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
          ? 'bg-lumera-600/20 border-lumera-500 scale-110 shadow-lg shadow-indigo-500/30'
          : 'bg-canvas-900/40 border-slate-700 hover:border-slate-600 hover:bg-canvas-800/60 hover:scale-105'
      }"
      onclick={() => handleSelect(option.value)}
    >
      <div class="text-4xl mb-2 transition-transform duration-300 {isSelected ? 'scale-125' : ''}">
        {option.emoji}
      </div>
      <div class="text-xs font-medium text-center {isSelected ? 'text-lumera-300' : 'text-slate-400'}">
        {option.label}
      </div>
    </button>
  {/each}
</div>
