<script lang="ts">
  interface Props {
    onProgressClick: () => void;
    subjectCount: number;
  }

  let { onProgressClick, subjectCount }: Props = $props();
</script>

<nav class="sticky top-[73px] z-40 border-b border-white/5 bg-canvas-950/80 backdrop-blur-md sub-nav-bar">
  <div class="px-6 py-2 flex items-center justify-end">
    <button
      onclick={onProgressClick}
      class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
      title="Tu Progreso"
    >
      <!-- Chart Icon SVG -->
      <svg
        class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
        />
      </svg>

      <!-- Badge -->
      {#if subjectCount > 0}
        <span class="absolute -top-1 -right-1 text-xs font-bold bg-lumera-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
          {subjectCount}
        </span>
      {/if}
    </button>
  </div>
</nav>

<style>
  /* Ensure smooth animations */
  .sub-nav-bar {
    animation: slideDown 0.4s ease-out;
  }

  @keyframes slideDown {
    from {
      transform: translateY(-100%);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
</style>
