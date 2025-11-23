<script lang="ts">
  import { onMount } from 'svelte';

  // Props
  let {
    categories = [],  // Array of { categoria, count }
    onCategoryChange = null
  } = $props();

  // Local state
  let selectedCategory = $state('Todos');

  // Category icons (SVG)
  const categoryIconsSVG = {
    'Todos': `<svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`,
    'Lectura': `<svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`,
    'Escritura': `<svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>`,
    'Comunicación Oral': `<svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>`,
    'Investigación': `<svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>`,
    'General': `<svg class="w-5 h-5 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" /></svg>`
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
          px-4 py-2.5 rounded-lg font-semibold text-sm transition-all duration-300 flex items-center
          {selectedCategory === cat.categoria
            ? 'bg-gradient-to-r from-lumera-600 to-focus-600 text-white shadow-lg scale-105'
            : 'bg-slate-50 text-slate-700 hover:bg-slate-100 hover:scale-102'}
        "
      >
        <span class="mr-2">{@html categoryIconsSVG[cat.categoria] || categoryIconsSVG['General']}</span>
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
          {cat.categoria} ({cat.count})
        </option>
      {/each}
    </select>
  </div>

  <!-- Active indicator line (desktop) -->
  <div class="hidden md:block mt-0.5 h-0.5 bg-gradient-to-r from-lumera-500 to-focus-500 rounded-full transition-all duration-300"></div>
</div>
