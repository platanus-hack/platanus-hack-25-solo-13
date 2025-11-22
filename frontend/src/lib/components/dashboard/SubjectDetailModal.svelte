<script lang="ts">
  import { onMount } from 'svelte';
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

  // Mock data for demonstration
  const mockUnits = [
    { id: 1, name: 'Unidad 1: Números y Álgebra', completed: true },
    { id: 2, name: 'Unidad 2: Geometría', completed: true },
    { id: 3, name: 'Unidad 3: Funciones', completed: false },
    { id: 4, name: 'Unidad 4: Probabilidades', completed: false }
  ];

  const mockStats = {
    exercisesCompleted: 45,
    totalExercises: 120,
    studyTimeHours: 12,
    lastActivity: '2 días atrás'
  };

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
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
      class="bg-canvas-900 rounded-2xl border border-slate-800 max-w-2xl w-full max-h-[90vh] overflow-y-auto shadow-2xl"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="sticky top-0 bg-canvas-900 border-b border-slate-800 p-6 flex items-center justify-between z-10">
        <div class="flex items-center gap-4">
          <div class="h-14 w-14 rounded-xl bg-gradient-to-br {subject.color} flex items-center justify-center text-3xl">
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
          class="h-10 w-10 rounded-full bg-canvas-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-6">
        {#if domainLevel === 0}
          <!-- Not evaluated state -->
          <div class="bg-gradient-to-br from-indigo-600/20 to-focus-600/20 border border-lumera-500/30 rounded-2xl p-6 text-center">
            <div class="text-4xl mb-3">📊</div>
            <h3 class="text-xl font-bold text-white mb-2">Aún no has sido evaluado</h3>
            <p class="text-slate-400 mb-4">
              Completa una evaluación diagnóstica para conocer tu nivel en {subject.name}
            </p>
            <a
              href="/diagnostico"
              class="inline-block px-6 py-3 rounded-xl bg-lumera-600 hover:bg-lumera-700 text-white font-bold transition-all"
            >
              Iniciar Evaluación Diagnóstica
            </a>
          </div>
        {:else}
          <!-- Statistics -->
          <section>
            <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              📊 Estadísticas
            </h3>
            <div class="grid grid-cols-2 gap-4">
              <div class="bg-canvas-800/50 rounded-xl p-4 border border-slate-700">
                <div class="text-2xl font-bold text-white">{mockStats.exercisesCompleted}/{mockStats.totalExercises}</div>
                <div class="text-xs text-slate-400 mt-1">Ejercicios completados</div>
              </div>
              <div class="bg-canvas-800/50 rounded-xl p-4 border border-slate-700">
                <div class="text-2xl font-bold text-white">{mockStats.studyTimeHours}h</div>
                <div class="text-xs text-slate-400 mt-1">Tiempo de estudio</div>
              </div>
            </div>
          </section>

          <!-- Units -->
          <section>
            <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              📚 Unidades del Curso
            </h3>
            <div class="space-y-2">
              {#each mockUnits as unit}
                <div class="bg-canvas-800/50 rounded-xl p-4 border border-slate-700 flex items-center gap-3">
                  <div class="h-6 w-6 rounded-full border-2 {unit.completed ? 'bg-green-500 border-green-500' : 'border-slate-600'} flex items-center justify-center">
                    {#if unit.completed}
                      <span class="text-white text-xs">✓</span>
                    {/if}
                  </div>
                  <span class="text-sm text-slate-300 flex-1">{unit.name}</span>
                  {#if !unit.completed}
                    <span class="text-xs text-lumera-400 hover:text-lumera-300 cursor-pointer">Empezar →</span>
                  {/if}
                </div>
              {/each}
            </div>
          </section>

          <!-- Areas to improve -->
          <section>
            <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              🎯 Áreas a Mejorar
            </h3>
            <div class="bg-amber-500/10 border border-amber-500/20 rounded-xl p-4">
              <ul class="space-y-2 text-sm text-slate-300">
                <li class="flex items-start gap-2">
                  <span class="text-amber-400 mt-0.5">•</span>
                  <span>Reforzar conceptos de funciones lineales</span>
                </li>
                <li class="flex items-start gap-2">
                  <span class="text-amber-400 mt-0.5">•</span>
                  <span>Practicar más ejercicios de geometría analítica</span>
                </li>
              </ul>
            </div>
          </section>

          <!-- Next step -->
          <section>
            <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              🚀 Siguiente Paso Recomendado
            </h3>
            <div class="bg-gradient-to-br from-indigo-600/20 to-focus-600/20 border border-lumera-500/30 rounded-xl p-4">
              <h4 class="font-semibold text-white mb-1">Unidad 3: Funciones</h4>
              <p class="text-sm text-slate-400 mb-3">
                Continúa con el estudio de funciones lineales y cuadráticas
              </p>
              <button class="px-4 py-2 rounded-lg bg-lumera-600 hover:bg-lumera-700 text-white text-sm font-semibold transition-all">
                Empezar a Estudiar
              </button>
            </div>
          </section>
        {/if}
      </div>
    </div>
  </div>
{/if}
