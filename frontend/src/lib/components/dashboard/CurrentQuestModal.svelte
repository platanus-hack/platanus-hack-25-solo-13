<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    onResumeQuest: () => void;
  }

  let { isOpen, onClose, onResumeQuest }: Props = $props();
  let backdropRef = $state<HTMLDivElement | null>(null);
  let modalRef = $state<HTMLDivElement | null>(null);

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      onClose();
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  // Animate modal entrance
  $effect(() => {
    if (isOpen && backdropRef && modalRef) {
      gsap.fromTo(backdropRef,
        { opacity: 0 },
        { opacity: 1, duration: 0.2, ease: 'power2.out' }
      );
      gsap.fromTo(modalRef,
        { opacity: 0, y: 20, scale: 0.95 },
        { opacity: 1, y: 0, scale: 1, duration: 0.3, ease: 'power2.out', delay: 0.1 }
      );
    }
  });
</script>

{#if isOpen}
  <!-- Backdrop -->
  <div
    bind:this={backdropRef}
    class="fixed inset-0 bg-black/60 backdrop-blur-sm z-60 flex items-center justify-center p-4"
    onclick={onClose}
  >
    <!-- Modal -->
    <div
      bind:this={modalRef}
      class="bg-canvas-900 rounded-2xl border border-slate-800 max-w-2xl w-full shadow-2xl"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="bg-canvas-900 border-b border-slate-800 p-6 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="h-12 w-12 rounded-xl bg-gradient-to-br from-focus-600/20 to-lumera-600/20 border border-focus-500/30 flex items-center justify-center">
            <svg class="w-6 h-6 text-focus-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
            </svg>
          </div>
          <div>
            <h2 class="text-xl font-bold text-white">Current Quest</h2>
            <p class="text-xs text-slate-400 mt-0.5">Tu aventura de aprendizaje actual</p>
          </div>
        </div>
        <button
          onclick={onClose}
          class="h-10 w-10 rounded-full bg-canvas-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="p-8">
        <!-- Quest Hero -->
        <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-lumera-900/90 via-canvas-900 to-canvas-900 border border-lumera-600/40 p-10 mb-6">
          <div class="relative z-10">
            <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-lumera-500/20 border border-lumera-400/30 text-lumera-300 text-xs font-semibold mb-4">
              <span class="animate-pulse">●</span> En Progreso
            </div>
            <h3 class="text-4xl font-display font-bold text-white mb-3 tracking-tight">Continue today's quest</h3>
            <p class="text-slate-400 text-lg mb-6">
              Pick up where you left off in <span class="text-focus-300 font-semibold">Math · Linear Functions</span>.
            </p>

            <!-- Progress Info -->
            <div class="grid grid-cols-3 gap-4 mb-6">
              <div class="bg-canvas-950/50 rounded-xl p-4 border border-slate-800">
                <div class="text-xs text-slate-500 mb-1">Progreso</div>
                <div class="text-2xl font-bold text-white">68%</div>
              </div>
              <div class="bg-canvas-950/50 rounded-xl p-4 border border-slate-800">
                <div class="text-xs text-slate-500 mb-1">Tiempo estimado</div>
                <div class="text-2xl font-bold text-white">15m</div>
              </div>
              <div class="bg-canvas-950/50 rounded-xl p-4 border border-slate-800">
                <div class="text-xs text-slate-500 mb-1">XP Restante</div>
                <div class="text-2xl font-bold text-achievement-400">120</div>
              </div>
            </div>

            <!-- Resume Button -->
            <button
              onclick={() => {
                onResumeQuest();
                onClose();
              }}
              class="w-full px-8 py-4 rounded-xl bg-white text-lumera-950 font-bold text-lg hover:scale-105 transition-all shadow-lg hover:shadow-xl"
            >
              Resume Quest →
            </button>
          </div>
        </div>

        <!-- Quick Stats -->
        <div class="grid grid-cols-2 gap-4">
          <div class="p-4 bg-canvas-800/40 rounded-xl border border-slate-700">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-lg bg-emerald-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-emerald-400" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <div class="text-xs text-slate-500">Ejercicios completados</div>
                <div class="text-lg font-bold text-white">12/18</div>
              </div>
            </div>
          </div>
          <div class="p-4 bg-canvas-800/40 rounded-xl border border-slate-700">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 rounded-lg bg-amber-500/20 flex items-center justify-center">
                <svg class="w-5 h-5 text-amber-400" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z" />
                </svg>
              </div>
              <div>
                <div class="text-xs text-slate-500">Racha actual</div>
                <div class="text-lg font-bold text-white">5 días</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}
