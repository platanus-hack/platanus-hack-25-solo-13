<script lang="ts">
  // Props
  let {
    onSearchChange = null,
    onSortChange = null
  } = $props();

  // Local state
  let searchQuery = $state('');
  let sortBy = $state('titulo');  // titulo, progreso

  function handleSearchInput(e: Event) {
    const target = e.target as HTMLInputElement;
    searchQuery = target.value;
    if (onSearchChange) {
      onSearchChange(searchQuery);
    }
  }

  function handleSortChange(newSort: string) {
    sortBy = newSort;
    if (onSortChange) {
      onSortChange(sortBy);
    }
  }
</script>

<div class="flex flex-col md:flex-row gap-4 mb-6">
  <!-- Search Input -->
  <div class="flex-1 relative">
    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
      <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>
    <input
      type="text"
      value={searchQuery}
      oninput={handleSearchInput}
      placeholder="Buscar por título o descripción..."
      class="w-full pl-12 pr-4 py-3 rounded-xl border-2 border-slate-200 bg-white text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-lumera-500 focus:border-lumera-500 transition-all"
    />
    {#if searchQuery}
      <button
        onclick={() => {
          searchQuery = '';
          if (onSearchChange) onSearchChange('');
        }}
        class="absolute inset-y-0 right-0 pr-4 flex items-center text-slate-400 hover:text-slate-600"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    {/if}
  </div>

  <!-- Sort Dropdown -->
  <div class="flex items-center gap-2">
    <span class="text-sm font-medium text-slate-600 whitespace-nowrap">Ordenar por:</span>
    <select
      bind:value={sortBy}
      onchange={() => handleSortChange(sortBy)}
      class="px-4 py-3 rounded-xl border-2 border-slate-200 bg-white text-slate-900 font-medium focus:outline-none focus:ring-2 focus:ring-lumera-500 focus:border-lumera-500"
    >
      <option value="titulo">Título (A-Z)</option>
      <option value="progreso">Progreso (%)</option>
    </select>
  </div>
</div>
