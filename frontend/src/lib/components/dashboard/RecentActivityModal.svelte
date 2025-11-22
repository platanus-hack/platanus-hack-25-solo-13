<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  interface Activity {
    id: number;
    text: string;
    time: string;
    icon: string;
  }

  interface Props {
    activities: Activity[];
    isOpen: boolean;
    onClose: () => void;
  }

  let { activities, isOpen, onClose }: Props = $props();
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
      class="bg-canvas-900 rounded-2xl border border-slate-800 max-w-lg w-full max-h-[80vh] overflow-y-auto shadow-2xl"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="sticky top-0 bg-canvas-900 border-b border-slate-800 p-6 flex items-center justify-between z-10">
        <div class="flex items-center gap-3">
          <div class="h-12 w-12 rounded-xl bg-gradient-to-br from-rose-600/20 to-orange-600/20 border border-rose-500/30 flex items-center justify-center text-2xl">
            🔔
          </div>
          <div>
            <h2 class="text-xl font-bold text-white">Actividad Reciente</h2>
            <p class="text-xs text-slate-400 mt-0.5">Tus últimas notificaciones</p>
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
      <div class="p-6">
        {#if activities.length > 0}
          <div class="space-y-4">
            {#each activities as act}
              <div class="flex gap-4 p-4 rounded-xl bg-canvas-800/40 border border-slate-700/50 hover:bg-canvas-800/60 hover:border-slate-600 transition-all">
                <div class="text-2xl flex-shrink-0">{act.icon}</div>
                <div class="flex-1">
                  <p class="text-sm text-slate-200 leading-relaxed">{act.text}</p>
                  <span class="text-xs text-slate-500 mt-1 inline-block">{act.time}</span>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <!-- Empty State -->
          <div class="text-center py-12">
            <div class="text-6xl mb-4">📭</div>
            <h3 class="text-lg font-semibold text-white mb-2">No hay actividad reciente</h3>
            <p class="text-sm text-slate-400">
              Tus notificaciones aparecerán aquí
            </p>
          </div>
        {/if}
      </div>

      <!-- Footer -->
      <div class="border-t border-slate-800 px-6 py-4 bg-canvas-950/50">
        <button class="w-full py-2 text-sm text-slate-400 hover:text-slate-200 transition-colors">
          Ver todas las notificaciones →
        </button>
      </div>
    </div>
  </div>
{/if}
