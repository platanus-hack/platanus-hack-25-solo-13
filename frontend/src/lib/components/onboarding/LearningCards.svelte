<script lang="ts">
  interface Card {
    id: string;
    label: string;
    description: string;
    icon: string;
  }

  interface Props {
    cards: Card[];
    selected: string | string[];
    onSelect: (value: string | string[]) => void;
    multiSelect?: boolean;
  }

  let { cards, selected = $bindable(''), onSelect, multiSelect = false }: Props = $props();

  function handleSelect(cardId: string) {
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

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
  {#each cards as card}
    {@const selected = isSelected(card.id)}
    <button
      type="button"
      class="p-6 rounded-2xl border-2 transition-all duration-300 text-left relative overflow-hidden {
        selected
          ? 'bg-gradient-to-br from-indigo-600/30 to-focus-600/20 border-lumera-500 scale-105 shadow-xl shadow-indigo-500/30'
          : 'bg-canvas-900/40 border-slate-700 hover:border-slate-600 hover:bg-canvas-800/60 hover:scale-102'
      }"
      onclick={() => handleSelect(card.id)}
    >
      {#if selected}
        <div class="absolute top-3 right-3 w-6 h-6 rounded-full bg-indigo-500 flex items-center justify-center text-white text-xs font-bold">
          ✓
        </div>
      {/if}

      <div class="text-4xl mb-3">{card.icon}</div>
      <h3 class="text-lg font-semibold mb-2 {selected ? 'text-indigo-200' : 'text-slate-200'}">
        {card.label}
      </h3>
      <p class="text-sm {selected ? 'text-lumera-300/80' : 'text-slate-400'}">
        {card.description}
      </p>
    </button>
  {/each}
</div>
