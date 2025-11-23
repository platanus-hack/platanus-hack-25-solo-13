<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { goto } from '$app/navigation';
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';
  import type { CustomizationItem } from '$lib/api/customization';

  interface Props {
    currentAvatar?: CustomizationItem | null;
    onProfileClick?: () => void;
    onQuestClick?: () => void;
    onMissionsClick?: () => void;
    onActivityClick?: () => void;
    onLiveEventsClick?: () => void;
    onProgressClick?: () => void;
    showNavButtons?: boolean;
    isHomePage?: boolean;
    centerContent?: any; // Slot for custom center content
  }

  let {
    currentAvatar = null,
    onProfileClick = () => {},
    onQuestClick = () => {},
    onMissionsClick = () => {},
    onActivityClick = () => {},
    onLiveEventsClick = () => {},
    onProgressClick = () => {},
    showNavButtons = true,
    isHomePage = false,
    centerContent
  }: Props = $props();

  // Get authenticated user info
  const student = {
    name: auth.user?.name || 'Student',
    grade: '2° Medio',
    level: 12,
    xp: 75,
    streak: 5
  };

  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  function handleLogout() {
    auth.logout();
    goto('/login');
  }
</script>

<header class="sticky top-0 z-50 border-b border-canvas-700/20 bg-canvas-800/95 backdrop-blur-md">
  <div class="px-6 py-3 grid grid-cols-3 items-center gap-4">
    <!-- Left Section: Profile + Home Button -->
    <div class="flex items-center gap-3 justify-self-start">
      <!-- Profile -->
      <button
        onclick={onProfileClick}
        class="flex items-center gap-3 hover:bg-canvas-900/60 rounded-xl p-2 -m-2 transition-colors group"
      >
        <div class="group-hover:scale-110 transition-transform">
          <AvatarDisplay
            currentAvatar={currentAvatar}
            {initials}
            size="small"
          />
        </div>
        <div class="text-left">
          <div class="flex items-center gap-2">
            <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
            <span class="px-1.5 py-0.5 rounded bg-canvas-800 border border-canvas-700 text-xs text-slate-400">{student.grade}</span>
          </div>
          <div class="text-xs text-slate-400">{auth.user?.email || 'student@lumera.com'}</div>
        </div>
      </button>

      <!-- Home Button (only show when not on home page) -->
      {#if !isHomePage}
        <button
          onclick={() => goto('/')}
          class="p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Volver al Dashboard"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
        </button>
      {/if}
    </div>

    <!-- Center - Badges or Custom Content -->
    <div class="flex items-center gap-2 justify-self-center">
      {#if centerContent}
        {@render centerContent()}
      {:else}
        <!-- Streak Badge -->
        <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-canvas-900/80 border border-canvas-700">
          <svg class="w-4 h-4 text-orange-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2.25c-2.429 0-4.817.178-7.152.521C2.87 3.061 1.5 4.795 1.5 6.741v6.018c0 1.946 1.37 3.68 3.348 3.97.877.129 1.761.234 2.652.316V21a.75.75 0 001.28.53l4.184-4.183a.39.39 0 01.266-.112c2.006-.05 3.982-.22 5.922-.506 1.978-.29 3.348-2.023 3.348-3.97V6.741c0-1.947-1.37-3.68-3.348-3.97A49.145 49.145 0 0012 2.25zM8.25 8.625a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25zm2.625 1.125a1.125 1.125 0 112.25 0 1.125 1.125 0 01-2.25 0zm4.875-1.125a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25z" />
          </svg>
          <span class="text-xs font-bold text-white">{student.streak}</span>
        </div>

        <!-- Mastery Tokens Badge -->
        <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-canvas-900/80 border border-canvas-700">
          <svg class="w-4 h-4 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="text-xs font-bold text-white">420</span>
        </div>
      {/if}
    </div>

    <!-- Right - Navigation Icons -->
    {#if showNavButtons}
      <div class="flex items-center gap-2 justify-self-end">
        <!-- Quest Button -->
        <button
          onclick={onQuestClick}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Current Quest"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
          </svg>
        </button>

        <!-- Missions Button -->
        <button
          onclick={onMissionsClick}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Mission Board"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          {#if dashboardStore.activeMissionCount > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-purple-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.activeMissionCount}
            </span>
          {/if}
        </button>

        <!-- Activity Button -->
        <button
          onclick={onActivityClick}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Actividad Reciente"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          {#if dashboardStore.activities.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-rose-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.activities.length}
            </span>
          {/if}
        </button>

        <!-- Live Events Button -->
        <button
          onclick={onLiveEventsClick}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Live Events"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
          {#if dashboardStore.events.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-amber-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.events.length}
            </span>
          {/if}
        </button>

        <!-- Progress Button -->
        <button
          onclick={onProgressClick}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Tu Progreso"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          {#if dashboardStore.subjects.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-lumera-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.subjects.length}
            </span>
          {/if}
        </button>

        <!-- Logout Button -->
        <button
          onclick={handleLogout}
          class="p-2 rounded-lg hover:bg-red-500/20 transition-all duration-200 group"
          title="Cerrar Sesión"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-red-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>
    {:else}
      <!-- Empty div to maintain grid structure -->
      <div></div>
    {/if}
  </div>
</header>
