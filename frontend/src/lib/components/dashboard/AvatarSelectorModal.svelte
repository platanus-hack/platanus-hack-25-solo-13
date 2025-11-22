<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import type { CustomizationItem } from '$lib/api/customization';
  import { getInventory, equipItem } from '$lib/api/customization';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    currentAvatarId: number | null;
    onAvatarEquipped: (avatar: CustomizationItem) => void;
  }

  let { isOpen, onClose, currentAvatarId, onAvatarEquipped }: Props = $props();

  let backdropRef = $state<HTMLDivElement | null>(null);
  let modalRef = $state<HTMLDivElement | null>(null);
  let avatars = $state<CustomizationItem[]>([]);
  let loading = $state(false);
  let error = $state<string | null>(null);
  let selectedAvatarId = $state(currentAvatarId);

  // Load owned avatars
  async function loadAvatars() {
    loading = true;
    error = null;
    try {
      const inventory = await getInventory();
      // Filter only avatars
      avatars = inventory
        .filter(item => item.item.type === 'avatar')
        .map(item => item.item);

      console.log('Loaded avatars:', avatars);
    } catch (err) {
      console.error('Error loading avatars:', err);
      error = 'No se pudieron cargar los avatares';
    } finally {
      loading = false;
    }
  }

  // Equip selected avatar
  async function handleEquip() {
    if (!selectedAvatarId) return;

    loading = true;
    error = null;
    try {
      await equipItem(selectedAvatarId, 'avatar');
      const equippedAvatar = avatars.find(a => a.id === selectedAvatarId);
      if (equippedAvatar) {
        onAvatarEquipped(equippedAvatar);
      }
      handleClose();
    } catch (err) {
      console.error('Error equipping avatar:', err);
      error = 'No se pudo equipar el avatar';
    } finally {
      loading = false;
    }
  }

  // Close with animation
  function handleClose() {
    if (backdropRef && modalRef) {
      gsap.to(modalRef, {
        scale: 0.95,
        opacity: 0,
        duration: 0.2,
        ease: 'power2.in'
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

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  // Animate modal entrance
  $effect(() => {
    if (isOpen && backdropRef && modalRef) {
      loadAvatars();
      selectedAvatarId = currentAvatarId;

      gsap.fromTo(backdropRef,
        { opacity: 0 },
        { opacity: 1, duration: 0.2, ease: 'power2.out' }
      );
      gsap.fromTo(modalRef,
        { scale: 0.95, opacity: 0 },
        { scale: 1, opacity: 1, duration: 0.3, ease: 'back.out(1.2)' }
      );
    }
  });

  // Get rarity color
  function getRarityColor(rarity: string): string {
    const colors: Record<string, string> = {
      common: 'from-slate-600 to-slate-700',
      rare: 'from-blue-600 to-blue-700',
      epic: 'from-purple-600 to-purple-700',
      legendary: 'from-amber-500 to-orange-600'
    };
    return colors[rarity] || colors.common;
  }

  // Get tier stars
  function getTierStars(tier: number): string {
    return '⭐'.repeat(tier);
  }
</script>

{#if isOpen}
  <!-- Backdrop -->
  <div
    bind:this={backdropRef}
    class="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    onclick={handleClose}
  ></div>

  <!-- Modal -->
  <div
    bind:this={modalRef}
    class="fixed inset-0 z-50 flex items-center justify-center p-4 pointer-events-none"
  >
    <div
      class="bg-canvas-900 rounded-2xl shadow-2xl border border-slate-700 w-full max-w-3xl max-h-[80vh] overflow-hidden pointer-events-auto"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="bg-gradient-to-r from-lumera-600 to-focus-600 p-6 flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-white">Selecciona tu Avatar</h2>
          <p class="text-sm text-white/80 mt-1">Elige un avatar de tu colección</p>
        </div>
        <button
          onclick={handleClose}
          class="h-10 w-10 rounded-full bg-white/20 hover:bg-white/30 flex items-center justify-center text-white transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 overflow-y-auto max-h-[calc(80vh-200px)]">
        {#if loading}
          <div class="flex items-center justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-4 border-lumera-500 border-t-transparent"></div>
          </div>
        {:else if error}
          <div class="bg-red-500/10 border border-red-500/30 rounded-lg p-4 text-red-400 text-center">
            {error}
          </div>
        {:else if avatars.length === 0}
          <div class="text-center py-12">
            <div class="text-6xl mb-4">🎭</div>
            <p class="text-slate-400 text-lg">No tienes avatares desbloqueados aún</p>
            <p class="text-slate-500 text-sm mt-2">Completa objetivos de aprendizaje para desbloquear avatares</p>
          </div>
        {:else}
          <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {#each avatars as avatar}
              <button
                onclick={() => selectedAvatarId = avatar.id}
                class="group relative aspect-square rounded-xl overflow-hidden border-2 transition-all {
                  selectedAvatarId === avatar.id
                    ? 'border-lumera-500 shadow-lg shadow-lumera-500/30 scale-105'
                    : 'border-slate-700 hover:border-slate-600 hover:scale-102'
                }"
              >
                <!-- Rarity gradient background -->
                <div class="absolute inset-0 bg-gradient-to-br {getRarityColor(avatar.rarity)} opacity-20"></div>

                <!-- Avatar image -->
                <div class="relative w-full h-full p-3">
                  <img
                    src="http://localhost:8080{avatar.image_url}"
                    alt={avatar.name}
                    class="w-full h-full object-contain"
                  />
                </div>

                <!-- Selected checkmark -->
                {#if selectedAvatarId === avatar.id}
                  <div class="absolute top-2 right-2 bg-lumera-500 rounded-full p-1.5 shadow-lg">
                    <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                {/if}

                <!-- Currently equipped badge -->
                {#if currentAvatarId === avatar.id}
                  <div class="absolute top-2 left-2 bg-green-500 text-white text-xs px-2 py-1 rounded-full shadow-lg font-medium">
                    Equipado
                  </div>
                {/if}

                <!-- Tier stars -->
                <div class="absolute bottom-2 left-2 text-xs">
                  {getTierStars(avatar.tier)}
                </div>

                <!-- Hover info -->
                <div class="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/80 to-transparent p-3 translate-y-full group-hover:translate-y-0 transition-transform">
                  <p class="text-white font-medium text-sm truncate">{avatar.name}</p>
                  <p class="text-white/70 text-xs capitalize">{avatar.rarity}</p>
                </div>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Footer -->
      <div class="border-t border-slate-800 p-6 bg-canvas-950 flex items-center justify-between">
        <div class="text-sm text-slate-400">
          {#if selectedAvatarId}
            {avatars.find(a => a.id === selectedAvatarId)?.name || 'Seleccionado'}
          {:else}
            Selecciona un avatar
          {/if}
        </div>
        <div class="flex gap-3">
          <button
            onclick={handleClose}
            class="px-6 py-2.5 rounded-lg bg-canvas-800 hover:bg-canvas-700 text-white font-medium transition-colors"
          >
            Cancelar
          </button>
          <button
            onclick={handleEquip}
            disabled={!selectedAvatarId || selectedAvatarId === currentAvatarId || loading}
            class="px-6 py-2.5 rounded-lg bg-gradient-to-r from-lumera-500 to-focus-500 hover:from-lumera-600 hover:to-focus-600 disabled:from-slate-700 disabled:to-slate-700 disabled:cursor-not-allowed text-white font-medium transition-colors shadow-lg shadow-lumera-500/20"
          >
            {loading ? 'Equipando...' : 'Equipar Avatar'}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
