<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    userName: string;
    userEmail: string;
    userGrade: string;
    level: number;
    xp: number;
    initials: string;
  }

  let { isOpen, onClose, userName, userEmail, userGrade, level, xp, initials }: Props = $props();
  let backdropRef = $state<HTMLDivElement | null>(null);
  let panelRef = $state<HTMLDivElement | null>(null);

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

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  // Animate panel entrance
  $effect(() => {
    if (isOpen && backdropRef && panelRef) {
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
        <div class="inline-flex items-center justify-center h-24 w-24 rounded-full bg-gradient-to-br from-lumera-500 to-focus-600 text-3xl font-bold text-white mb-4 shadow-lg">
          {initials}
        </div>
        <h3 class="text-2xl font-bold text-white mb-1">{userName}</h3>
        <p class="text-sm text-slate-400 mb-2">{userEmail}</p>
        <span class="inline-block px-3 py-1 rounded-full bg-canvas-800 border border-canvas-700 text-sm text-slate-300">
          {userGrade}
        </span>
      </div>

      <!-- Level Progress -->
      <div class="bg-canvas-800/40 rounded-2xl border border-slate-700 p-6">
        <div class="mb-3">
          <span class="text-sm font-semibold text-achievement-400">Level {level}</span>
        </div>

        <!-- XP Bar -->
        <div class="space-y-2">
          <div class="h-3 w-full bg-canvas-900 rounded-full overflow-hidden border border-slate-700">
            <div
              class="h-full bg-gradient-to-r from-lumera-500 via-focus-500 to-achievement-400 transition-all duration-500 rounded-full"
              style="width: {xp}%"
            ></div>
          </div>
          <div class="flex justify-end text-xs text-slate-500">
            <span>1000 XP</span>
          </div>
        </div>
      </div>

      <!-- Action Cards -->
      <div class="grid grid-cols-3 gap-3 pt-4 border-t border-slate-800">
        <button class="aspect-square rounded-xl bg-canvas-800/60 border border-slate-700 text-slate-200 hover:bg-canvas-700 hover:border-slate-600 transition-all flex flex-col items-center justify-center gap-2 p-4">
          <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <span class="text-xs font-medium text-center">Configuración</span>
        </button>

        <button class="aspect-square rounded-xl bg-canvas-800/60 border border-slate-700 text-slate-200 hover:bg-canvas-700 hover:border-slate-600 transition-all flex flex-col items-center justify-center gap-2 p-4">
          <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
          <span class="text-xs font-medium text-center">Mi Aprendizaje</span>
        </button>

        <button class="aspect-square rounded-xl bg-canvas-800/60 border border-slate-700 text-slate-200 hover:bg-canvas-700 hover:border-slate-600 transition-all flex flex-col items-center justify-center gap-2 p-4">
          <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
          </svg>
          <span class="text-xs font-medium text-center">Logros</span>
        </button>
      </div>
    </div>
  </div>
{/if}
