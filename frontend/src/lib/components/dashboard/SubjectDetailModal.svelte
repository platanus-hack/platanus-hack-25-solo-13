<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import type { Subject } from '$lib/constants/subjects';
  import { getDomainLevelInfo } from '$lib/constants/subjects';

  interface Props {
    subject: Subject | null;
    domainLevel: number;
    isOpen: boolean;
    onClose: () => void;
  }

  let { subject, domainLevel, isOpen, onClose }: Props = $props();

  const levelInfo = $derived(subject ? getDomainLevelInfo(domainLevel) : null);

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      onClose();
    }
  }

  function navigateToObjetivos() {
    if (subject?.materiaId) {
      goto(`/materias/${subject.materiaId}/objetivos`);
      onClose();
    }
  }

  onMount(() => {
    if (isOpen) {
      window.addEventListener('keydown', handleKeydown);
      return () => window.removeEventListener('keydown', handleKeydown);
    }
  });
</script>

{#if isOpen && subject}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 top-[72px] bg-black/60 backdrop-blur-sm z-50 flex items-start justify-center pt-24 pb-8 px-4"
    onclick={onClose}
  >
    <!-- Modal -->
    <div
      class="bg-canvas-900 rounded-2xl border border-slate-700 max-w-2xl w-full min-h-[50vh] max-h-[80vh] overflow-y-auto shadow-2xl"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="sticky top-0 bg-canvas-900 border-b border-slate-700 p-6 flex items-center justify-between z-10">
        <div class="flex items-center gap-4">
          <div class="h-14 w-14 rounded-xl bg-gradient-to-br from-slate-700 to-slate-800 flex items-center justify-center text-3xl border border-slate-600 shadow-lg">
            {subject.icon}
          </div>
          <div>
            <h2 class="text-2xl font-bold text-white">{subject.name}</h2>
            <div class="flex items-center gap-2 mt-1">
              <span class="text-xs font-bold {levelInfo?.badgeColor} px-2 py-1 rounded-full {levelInfo?.textColor}">
                {levelInfo?.label}
              </span>
              {#if domainLevel === 0}
                <span class="text-xs text-slate-400">• Completa la evaluación diagnóstica</span>
              {/if}
            </div>
          </div>
        </div>
        <button
          onclick={onClose}
          class="h-10 w-10 rounded-full bg-canvas-800 hover:bg-canvas-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-6">
        {#if domainLevel === 0}
          <!-- Not evaluated state -->
          <div class="bg-canvas-800 border border-slate-700 rounded-2xl p-6 text-center">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-blue-500/20 to-purple-500/20 border border-blue-500/30 flex items-center justify-center">
              <svg class="w-8 h-8 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-2">Aún no has sido evaluado</h3>
            <p class="text-slate-300 mb-4">
              Completa una evaluación diagnóstica para conocer tu nivel en {subject.name}
            </p>
            <a
              href="/diagnostico/{subject.materiaId}"
              class="inline-block px-6 py-3 rounded-xl bg-[#E1E1E1] hover:bg-[#CCCCCC] text-canvas-900 font-bold transition-all shadow-lg hover:shadow-xl"
            >
              Iniciar Evaluación Diagnóstica
            </a>
          </div>
        {:else}
          <!-- Evaluated state - Show Objetivos de Aprendizaje CTA -->
          <div class="bg-canvas-800 border border-slate-700 rounded-2xl p-6 text-center">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gradient-to-br from-lumera-500/20 to-focus-500/20 border border-lumera-500/30 flex items-center justify-center">
              <svg class="w-8 h-8 text-lumera-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-2">Objetivos de Aprendizaje</h3>
            <p class="text-slate-300 mb-5">
              Explora los objetivos de aprendizaje disponibles y genera planes personalizados para mejorar tu dominio en {subject.name}
            </p>
            <button
              onclick={navigateToObjetivos}
              class="px-6 py-3 rounded-xl bg-[#E1E1E1] hover:bg-[#CCCCCC] text-canvas-900 font-bold transition-all shadow-lg hover:shadow-xl"
            >
              Ver Objetivos de Aprendizaje →
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
