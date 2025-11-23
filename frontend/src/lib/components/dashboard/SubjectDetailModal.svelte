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
    class="fixed inset-0 bg-black/60 backdrop-blur-sm z-60 flex items-center justify-center p-4"
    onclick={onClose}
  >
    <!-- Modal -->
    <div
      class="bg-canvas-900 rounded-2xl border border-slate-700 max-w-2xl w-full max-h-[90vh] overflow-y-auto shadow-2xl"
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
            <div class="text-4xl mb-3">📊</div>
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
          <div class="space-y-4">
            <!-- Level summary -->
            <div class="bg-canvas-800 rounded-xl p-4 border border-slate-700">
              <div class="text-sm text-slate-400 mb-1">Tu nivel actual</div>
              <div class="flex items-center gap-3">
                <span class="text-3xl">{levelInfo?.label === 'Recordar' ? '⭐' : levelInfo?.label === 'Comprender' ? '⭐⭐' : levelInfo?.label === 'Aplicar' ? '⭐⭐⭐' : levelInfo?.label === 'Analizar' ? '⭐⭐⭐⭐' : levelInfo?.label === 'Evaluar' ? '⭐⭐⭐⭐⭐' : '⭐⭐⭐⭐⭐⭐'}</span>
                <div>
                  <div class="text-lg font-bold text-white">{levelInfo?.label}</div>
                  <div class="text-xs text-slate-400">Nivel de Bloom dominado</div>
                </div>
              </div>
            </div>

            <!-- Main CTA -->
            <div class="bg-canvas-800 border border-slate-700 rounded-2xl p-6 text-center">
              <div class="text-4xl mb-3">🎯</div>
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
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
