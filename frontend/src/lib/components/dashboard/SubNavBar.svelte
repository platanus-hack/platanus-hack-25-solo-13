<script lang="ts">
  interface Props {
    onProgressClick: () => void;
    onActivityClick: () => void;
    onMissionsClick: () => void;
    onLogoutClick: () => void;
    subjectCount: number;
    activityCount?: number;
    missionCount?: number;
  }

  let { onProgressClick, onActivityClick, onMissionsClick, onLogoutClick, subjectCount, activityCount = 0, missionCount = 0 }: Props = $props();
</script>

<nav class="sticky top-[73px] z-40 border-b border-white/5 bg-canvas-950/80 backdrop-blur-md sub-nav-bar">
  <div class="px-6 py-2 flex items-center justify-end gap-2">
    <!-- Missions Button -->
    <button
      onclick={onMissionsClick}
      class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
      title="Mission Board"
    >
      <!-- Clipboard/Tasks Icon SVG -->
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
          d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"
        />
      </svg>

      <!-- Badge -->
      {#if missionCount > 0}
        <span class="absolute -top-1 -right-1 text-xs font-bold bg-purple-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
          {missionCount}
        </span>
      {/if}
    </button>

    <!-- Activity/Notifications Button -->
    <button
      onclick={onActivityClick}
      class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
      title="Actividad Reciente"
    >
      <!-- Bell Icon SVG -->
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
          d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
        />
      </svg>

      <!-- Badge -->
      {#if activityCount > 0}
        <span class="absolute -top-1 -right-1 text-xs font-bold bg-rose-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
          {activityCount}
        </span>
      {/if}
    </button>

    <!-- Progress Button -->
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

    <!-- Logout Button -->
    <button
      onclick={onLogoutClick}
      class="relative p-2 rounded-lg hover:bg-red-900/60 transition-all duration-200 group"
      title="Logout"
    >
      <!-- Logout Icon SVG -->
      <svg
        class="w-6 h-6 text-slate-400 group-hover:text-red-400 transition-colors"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
        />
      </svg>
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
