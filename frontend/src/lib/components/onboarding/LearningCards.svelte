<script lang="ts">
  interface Card {
    id: string;
    label: string;
    description: string;
    iconPath: string;
    disabled?: boolean;
  }

  interface Props {
    cards: Card[];
    selected: string | string[];
    onSelect: (value: string | string[]) => void;
    multiSelect?: boolean;
    compact?: boolean;
  }

  let { cards, selected = $bindable(''), onSelect, multiSelect = false, compact = false }: Props = $props();

  function handleSelect(cardId: string, isDisabled: boolean) {
    if (isDisabled) return;

    if (multiSelect) {
      const currentSelected = Array.isArray(selected) ? selected : [];
      if (currentSelected.includes(cardId)) {
        const newSelected = currentSelected.filter(id => id !== cardId);
        selected = newSelected;
        onSelect(newSelected);
      } else {
        const newSelected = [...currentSelected, cardId];
        selected = newSelected;
        onSelect(newSelected);
      }
    } else {
      selected = cardId;
      onSelect(cardId);
    }
  }

  function isSelected(cardId: string): boolean {
    if (multiSelect && Array.isArray(selected)) {
      return selected.includes(cardId);
    }
    return selected === cardId;
  }
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 {compact ? 'gap-2' : 'gap-4'}">
  {#each cards as card}
    {@const selected = isSelected(card.id)}
    {@const isDisabled = card.disabled || false}
    <button
      type="button"
      disabled={isDisabled}
      class="{compact ? 'p-3 rounded-lg' : 'p-6 rounded-2xl'} border-2 transition-all duration-300 text-left relative overflow-hidden {
        isDisabled
          ? 'bg-canvas-900/30 border-canvas-800 opacity-60 cursor-not-allowed'
          : selected
          ? 'bg-[#E1E1E1] border-[#E1E1E1] scale-105 shadow-xl'
          : 'bg-canvas-900/60 border-canvas-700 hover:border-canvas-600 hover:bg-canvas-800/60 hover:scale-102'
      }"
      onclick={() => handleSelect(card.id, isDisabled)}
    >
      {#if isDisabled}
        <div class="absolute {compact ? 'top-1.5 right-1.5 px-1.5 py-0.5' : 'top-3 right-3 px-3 py-1'} bg-canvas-800 rounded-full border border-canvas-600">
          <span class="{compact ? 'text-[9px]' : 'text-xs'} font-semibold text-slate-400">Próximamente</span>
        </div>
      {:else if selected}
        <div class="absolute {compact ? 'top-1.5 right-1.5 w-5 h-5' : 'top-3 right-3 w-7 h-7'} rounded-full bg-canvas-900 flex items-center justify-center shadow-md">
          <svg class="{compact ? 'w-3 h-3' : 'w-4 h-4'} text-white" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
        </div>
      {/if}

      <div class="flex items-center justify-center {compact ? 'w-10 h-10 mb-1.5' : 'w-16 h-16 mb-4'} rounded-lg {
        isDisabled ? 'bg-canvas-800/40' :
        selected ? 'bg-canvas-900/10' : 'bg-canvas-800/60'
      }">
        <svg
          class="{compact ? 'w-6 h-6' : 'w-10 h-10'} {
            isDisabled ? 'text-slate-600' :
            selected ? 'text-canvas-900' : 'text-slate-400'
          } transition-colors duration-300"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d={card.iconPath} />
        </svg>
      </div>

      <h3 class="{compact ? 'text-xs' : 'text-lg'} font-bold {compact ? 'mb-0.5' : 'mb-2'} {
        isDisabled ? 'text-slate-500' :
        selected ? 'text-canvas-900' : 'text-slate-200'
      }">
        {card.label}
      </h3>
      <p class="{compact ? 'text-[10px] leading-tight' : 'text-sm'} {
        isDisabled ? 'text-slate-600' :
        selected ? 'text-canvas-800' : 'text-slate-400'
      }">
        {card.description}
      </p>
    </button>
  {/each}
</div>
