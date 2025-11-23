<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { goto } from '$app/navigation';
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';
  import type { CustomizationItem } from '$lib/api/customization';

  interface Props {
    currentAvatar?: CustomizationItem | null;
    onProfileClick?: () => void;
    onProgressClick?: () => void;
    showNavButtons?: boolean;
    isHomePage?: boolean;
    centerContent?: any; // Slot for custom center content
  }

  let {
    currentAvatar = null,
    onProfileClick = () => {},
    onProgressClick = () => {},
    showNavButtons = true,
    isHomePage = false,
    centerContent
  }: Props = $props();

  // Get authenticated user info
  const student = $derived({
    name: auth.user?.name || 'Student',
    grade: auth.user?.curso_actual || '2° Medio',
    level: auth.gamificationStats?.level || 1,
    xp: auth.gamificationStats?.xp || 0,
    streak: auth.gamificationStats?.current_streak || 0,
    coins: auth.gamificationStats?.coins || 0
  });

  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  function handleLogout() {
    auth.logout();
    goto('/login');
  }

  // Load gamification stats on mount
  onMount(() => {
    auth.loadGamificationStats();
  });
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
          <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
          <div class="text-xs text-slate-400">{student.grade}</div>
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
    <div class="flex items-center justify-self-center {centerContent ? 'w-full' : 'gap-2'}">
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
          <span class="text-xs font-bold text-white">{student.coins}</span>
        </div>
      {/if}
    </div>

    <!-- Right - Navigation Icons -->
    {#if showNavButtons}
      <div class="flex items-center gap-2 justify-self-end">
        <!-- Progress Button -->
        <button
          onclick={onProgressClick}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Tu Progreso"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
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
