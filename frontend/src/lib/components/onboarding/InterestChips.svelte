<script lang="ts">
  interface Props {
    selected: string[];
    onSelect: (interests: string[]) => void;
  }

  let { selected = $bindable([]), onSelect }: Props = $props();

  const availableInterests = [
    { id: 'tecnologia', label: 'Tecnología', icon: '💻' },
    { id: 'musica', label: 'Música', icon: '🎵' },
    { id: 'deportes', label: 'Deportes', icon: '⚽' },
    { id: 'arte', label: 'Arte', icon: '🎨' },
    { id: 'ciencia', label: 'Ciencia', icon: '🔬' },
    { id: 'videojuegos', label: 'Videojuegos', icon: '🎮' },
    { id: 'lectura', label: 'Lectura', icon: '📚' },
    { id: 'cine', label: 'Cine', icon: '🎬' },
    { id: 'naturaleza', label: 'Naturaleza', icon: '🌿' },
    { id: 'cocina', label: 'Cocina', icon: '👨‍🍳' },
    { id: 'moda', label: 'Moda', icon: '👗' },
    { id: 'viajes', label: 'Viajes', icon: '✈️' }
  ];

  function toggleInterest(interestId: string) {
    if (selected.includes(interestId)) {
      selected = selected.filter(id => id !== interestId);
    } else {
      selected = [...selected, interestId];
    }
    onSelect(selected);
  }
</script>

<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
  {#each availableInterests as interest}
    {@const isSelected = selected.includes(interest.id)}
    <button
      type="button"
      class="p-4 rounded-xl border-2 transition-all duration-300 text-center {
        isSelected
          ? 'bg-indigo-600/20 border-indigo-500 scale-105 shadow-lg shadow-indigo-500/20'
          : 'bg-slate-900/40 border-slate-700 hover:border-slate-600 hover:bg-slate-800/60'
      }"
      onclick={() => toggleInterest(interest.id)}
    >
      <div class="text-3xl mb-2">{interest.icon}</div>
      <div class="text-sm font-medium {isSelected ? 'text-indigo-300' : 'text-slate-300'}">
        {interest.label}
      </div>
      {#if isSelected}
        <div class="mt-1 text-xs text-indigo-400">✓ Seleccionado</div>
      {/if}
    </button>
  {/each}
</div>
