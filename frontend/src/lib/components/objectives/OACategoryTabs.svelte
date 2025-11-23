<script lang="ts">
  import { onMount } from 'svelte';

  // Props
  let {
    categories = [],  // Array of { categoria, count }
    onCategoryChange = null
  } = $props();

  // Local state
  let selectedCategory = $state('Todos');

  // Category icons
  const categoryIcons = {
    'Todos': '📚',
    'Lectura': '📖',
    'Escritura': '✍️',
    'Comunicación Oral': '💬',
    'Investigación': '📝',
    'General': '📚'
  };

  // Calculate total count
  const totalCount = categories.reduce((sum, cat) => sum + cat.count, 0);

  // Build full category list with "Todos" first
  const allCategories = [
    { categoria: 'Todos', count: totalCount },
    ...categories.filter(c => c.count > 0)
  ];

  function handleCategoryClick(categoria: string) {
    selectedCategory = categoria;
    if (onCategoryChange) {
      onCategoryChange(categoria);
    }
  }
</script>

<div class="mb-6">
  <!-- Desktop View -->
  <div class="hidden md:flex items-center gap-2 flex-wrap bg-white rounded-xl p-2 border-2 border-slate-200 shadow-sm">
    {#each allCategories as cat}
      <button
        onclick={() => handleCategoryClick(cat.categoria)}
        class="
          px-4 py-2.5 rounded-lg font-semibold text-sm transition-all duration-300
          {selectedCategory === cat.categoria
            ? 'bg-gradient-to-r from-lumera-600 to-focus-600 text-white shadow-lg scale-105'
            : 'bg-slate-50 text-slate-700 hover:bg-slate-100 hover:scale-102'}
        "
      >
        <span class="mr-2">{categoryIcons[cat.categoria] || '📚'}</span>
        {cat.categoria}
        <span class="ml-2 px-2 py-0.5 rounded-full text-xs font-bold {selectedCategory === cat.categoria ? 'bg-white/20' : 'bg-slate-200 text-slate-600'}">
          {cat.count}
        </span>
      </button>
    {/each}
  </div>

  <!-- Mobile View (Dropdown) -->
  <div class="md:hidden">
    <select
      bind:value={selectedCategory}
      onchange={(e) => handleCategoryClick(e.target.value)}
      class="w-full px-4 py-3 rounded-xl border-2 border-slate-200 bg-white text-slate-900 font-semibold focus:outline-none focus:ring-2 focus:ring-lumera-500 focus:border-lumera-500"
    >
      {#each allCategories as cat}
        <option value={cat.categoria}>
          {categoryIcons[cat.categoria] || '📚'} {cat.categoria} ({cat.count})
        </option>
      {/each}
    </select>
  </div>

  <!-- Active indicator line (desktop) -->
  <div class="hidden md:block mt-0.5 h-0.5 bg-gradient-to-r from-lumera-500 to-focus-500 rounded-full transition-all duration-300"></div>
</div>
