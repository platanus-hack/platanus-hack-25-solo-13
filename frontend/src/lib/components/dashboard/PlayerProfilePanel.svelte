<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import AvatarSelectorModal from './AvatarSelectorModal.svelte';
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';
  import { getEquipment } from '$lib/api/customization';
  import type { CustomizationItem } from '$lib/api/customization';

  interface GamificationStats {
    level: number;
    xp: number;
    xp_for_next_level: number;
    xp_progress: number;
    coins: number;
    current_streak: number;
    longest_streak: number;
  }

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    userName: string;
    userEmail: string;
    userGrade: string;
    level: number;
    xp: number;
    initials: string;
    gamificationStats?: GamificationStats | null;
    onAvatarChanged?: (avatar: CustomizationItem) => void;
  }

  let { isOpen, onClose, userName, userEmail, userGrade, level, xp, initials, gamificationStats, onAvatarChanged }: Props = $props();

  // Calculate progress percentage for current level
  const progressPercentage = $derived(() => {
    if (!gamificationStats) {
      // Fallback calculation if stats not available
      const currentLevelXP = (level - 1) * (level - 1) * 100;
      const nextLevelXP = level * level * 100;
      const progress = xp - currentLevelXP;
      const total = nextLevelXP - currentLevelXP;
      return Math.min(100, Math.max(0, (progress / total) * 100));
    }

    // Use backend-calculated progress
    const currentLevelXP = gamificationStats.xp - gamificationStats.xp_progress;
    const total = gamificationStats.xp_for_next_level - currentLevelXP;
    return Math.min(100, Math.max(0, (gamificationStats.xp_progress / total) * 100));
  });

  const xpText = $derived(() => {
    if (!gamificationStats) {
      const nextLevelXP = level * level * 100;
      return `${xp} / ${nextLevelXP} XP`;
    }
    return `${gamificationStats.xp} / ${gamificationStats.xp_for_next_level} XP`;
  });
  let backdropRef = $state<HTMLDivElement | null>(null);
  let panelRef = $state<HTMLDivElement | null>(null);
  let avatarSelectorOpen = $state(false);
  let currentAvatar = $state<CustomizationItem | null>(null);
  let avatarLoading = $state(false);

  // Close with animation
  function handleClose() {
    if (backdropRef && panelRef) {
      gsap.to(panelRef, {
        x: -400,
        opacity: 0.8,
        duration: 0.25,
        ease: 'power3.in'
      });
      gsap.to(backdropRef, {
        opacity: 0,
        duration: 0.2,
        ease: 'power2.in',
        onComplete: () => {
          onClose();
        }
      });
    } else {
      onClose();
    }
  }

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      handleClose();
    }
  }

  // Load equipped avatar
  async function loadEquipment() {
    avatarLoading = true;
    try {
      const equipment = await getEquipment();
      currentAvatar = equipment.equipped_avatar || null;
    } catch (err) {
      console.error('Error loading equipment:', err);
    } finally {
      avatarLoading = false;
    }
  }

  // Handle avatar equipped
  function handleAvatarEquipped(avatar: CustomizationItem) {
    currentAvatar = avatar;
    // Notify parent component
    if (onAvatarChanged) {
      onAvatarChanged(avatar);
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  // Animate panel entrance and load avatar
  $effect(() => {
    if (isOpen && backdropRef && panelRef) {
      loadEquipment();

      gsap.fromTo(backdropRef,
        { opacity: 0 },
        { opacity: 1, duration: 0.2, ease: 'power2.out' }
      );
      gsap.fromTo(panelRef,
        { x: -400, opacity: 0.8 },
        { x: 0, opacity: 1, duration: 0.3, ease: 'power3.out' }
      );
    }
  });
</script>

{#if isOpen}
  <!-- Backdrop -->
  <div
    bind:this={backdropRef}
    class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
    onclick={handleClose}
  ></div>

  <!-- Panel -->
  <div
    bind:this={panelRef}
    class="fixed left-0 top-0 h-full w-96 bg-canvas-900 border-r border-slate-800 shadow-2xl z-50 overflow-y-auto"
  >
    <!-- Header -->
    <div class="sticky top-0 bg-canvas-900 border-b border-slate-800 p-6 flex items-center justify-between z-10">
      <h2 class="text-xl font-bold text-white">Mi Perfil</h2>
      <button
        onclick={handleClose}
        class="h-10 w-10 rounded-full bg-canvas-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
      >
        ✕
      </button>
    </div>

    <!-- Content -->
    <div class="p-6 space-y-6">
      <!-- Avatar & Basic Info -->
      <div class="text-center">
        <!-- Avatar with edit button -->
        <div class="mb-4">
          <AvatarDisplay
            {currentAvatar}
            {initials}
            size="medium"
            editable={true}
            loading={avatarLoading}
            onEdit={() => avatarSelectorOpen = true}
          />
        </div>

        <h3 class="text-2xl font-bold text-white mb-2">{userName}</h3>
        <p class="text-sm text-slate-400 mb-3">{userEmail}</p>
        <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-gradient-to-r from-lumera-500/10 to-focus-500/10 border border-lumera-500/30">
          <svg class="w-4 h-4 text-lumera-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
          <span class="text-sm font-semibold text-lumera-300">{userGrade}</span>
        </div>
      </div>

      <!-- Level Progress -->
      <div class="pt-2">
        <div class="mb-3">
          <span class="text-sm font-semibold text-achievement-400">Level {level}</span>
        </div>

        <!-- XP Bar -->
        <div class="space-y-2">
          <div class="h-3 w-full bg-canvas-900 rounded-full overflow-hidden border border-slate-700">
            <div
              class="h-full bg-gradient-to-r from-lumera-500 via-focus-500 to-achievement-400 transition-all duration-500 rounded-full"
              style="width: {progressPercentage()}%"
            ></div>
          </div>
          <div class="flex justify-end text-xs text-slate-500">
            <span>{xpText()}</span>
          </div>
        </div>
      </div>

      <!-- Menu Options -->
      <div class="pt-6 border-t border-slate-800 space-y-2">
        <a
          href="/configuracion"
          class="flex items-center gap-3 px-4 py-3 rounded-xl hover:bg-canvas-800 transition-colors group"
        >
          <div class="w-10 h-10 rounded-lg bg-lumera-500/10 flex items-center justify-center group-hover:bg-lumera-500/20 transition-colors">
            <svg class="w-5 h-5 text-lumera-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <div class="flex-1">
            <h4 class="text-sm font-semibold text-white group-hover:text-lumera-300 transition-colors">Configuración</h4>
            <p class="text-xs text-slate-500">Actualiza tus preferencias</p>
          </div>
          <svg class="w-5 h-5 text-slate-600 group-hover:text-slate-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </a>
      </div>
    </div>
  </div>
{/if}

<!-- Avatar Selector Modal -->
<AvatarSelectorModal
  isOpen={avatarSelectorOpen}
  onClose={() => avatarSelectorOpen = false}
  currentAvatarId={currentAvatar?.id || null}
  onAvatarEquipped={handleAvatarEquipped}
/>
