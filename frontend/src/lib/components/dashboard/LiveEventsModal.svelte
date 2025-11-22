<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  interface Event {
    id: number;
    title: string;
    subtitle: string;
    color: string;
  }

  interface Props {
    events: Event[];
    isOpen: boolean;
    onClose: () => void;
  }

  let { events, isOpen, onClose }: Props = $props();
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
      class="bg-canvas-900 rounded-2xl border border-slate-800 max-w-2xl w-full max-h-[80vh] overflow-y-auto shadow-2xl"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="sticky top-0 bg-canvas-900 border-b border-slate-800 p-6 flex items-center justify-between z-10">
        <div class="flex items-center gap-3">
          <div class="h-12 w-12 rounded-xl bg-gradient-to-br from-rose-600/20 to-amber-600/20 border border-rose-500/30 flex items-center justify-center">
            <svg class="w-6 h-6 text-rose-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </div>
          <div>
            <h2 class="text-xl font-bold text-white">Live Events</h2>
            <p class="text-xs text-slate-400 mt-0.5">Eventos especiales activos</p>
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
        {#if events.length > 0}
          <div class="space-y-4">
            {#each events as event}
              <button class="
                w-full text-left p-6 rounded-2xl bg-gradient-to-br {event.color} text-white
                relative overflow-hidden hover:scale-[1.02] transition-transform
                shadow-lg hover:shadow-xl
              ">
                <!-- Decorative background pattern -->
                <div class="absolute inset-0 opacity-10">
                  <div class="absolute top-0 right-0 w-64 h-64 bg-white rounded-full blur-3xl transform translate-x-32 -translate-y-32"></div>
                </div>

                <div class="relative z-10">
                  <!-- Badge -->
                  <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/20 backdrop-blur-sm text-xs font-bold mb-3">
                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10 10-4.477 10-10S17.523 2 12 2zm0 18c-4.411 0-8-3.589-8-8s3.589-8 8-8 8 3.589 8 8-3.589 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z" />
                    </svg>
                    LIMITED TIME
                  </div>

                  <!-- Title -->
                  <h3 class="font-bold text-2xl mb-2">{event.title}</h3>
                  <p class="text-white/90 text-base mb-4">{event.subtitle}</p>

                  <!-- CTA -->
                  <div class="flex items-center gap-2 text-sm font-semibold">
                    <span>Participate now</span>
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                    </svg>
                  </div>
                </div>
              </button>
            {/each}
          </div>

          <!-- Footer info -->
          <div class="mt-6 p-4 bg-canvas-800/40 rounded-xl border border-slate-700">
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-lumera-400 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 24 24">
                <path d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div class="text-xs text-slate-400">
                <p class="font-semibold text-slate-300 mb-1">Sobre los eventos</p>
                <p>Los eventos especiales ofrecen recompensas exclusivas y desafíos únicos. Participa antes de que terminen para ganar XP extra y badges especiales.</p>
              </div>
            </div>
          </div>
        {:else}
          <!-- Empty State -->
          <div class="text-center py-12">
            <div class="text-6xl mb-4">🎉</div>
            <h3 class="text-lg font-semibold text-white mb-2">No hay eventos activos</h3>
            <p class="text-sm text-slate-400">
              Los nuevos eventos aparecerán aquí cuando estén disponibles
            </p>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
