<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import LessonPlayer from '$lib/components/slides/LessonPlayer.svelte';
  import type { LearningPlan, LearningPlanComponent } from '$lib/api/learningPlans';
  import { generateComponentContent } from '$lib/api/learningPlans';

  interface Props {
    plan: LearningPlan | null;
    isOpen: boolean;
    onClose: () => void;
    onComplete?: (planId: number) => void;
  }

  let { plan, isOpen, onClose, onComplete }: Props = $props();

  // Convert learning plan to lesson format for LessonPlayer
  const lesson = $derived(
    plan ? {
      leccionId: `plan-${plan.id}`,
      titulo: plan.titulo,
      materia: 'learning-plan',
      slides: [...plan.components]
        .sort((a, b) => a.orden - b.orden)
        .map(component => ({
          tipo: component.tipo_componente,
          props: component.contenido_props || {},
          componentId: component.id,
          estado: component.estado
        }))
    } : null
  );

  // Handle lesson completion
  function handleComplete(completionData: any) {
    if (plan && onComplete) {
      onComplete(plan.id);
    }
  }

  // Close modal
  function handleClose() {
    onClose();
  }

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      handleClose();
    }
  }

  // Animate modal entrance
  $effect(() => {
    if (isOpen && plan) {
      gsap.fromTo(
        '.learning-plan-modal',
        { opacity: 0, scale: 0.95 },
        { opacity: 1, scale: 1, duration: 0.3, ease: 'power2.out' }
      );
    }
  });

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });
</script>

{#if isOpen && plan}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4"
    onclick={handleClose}
  >
    <!-- Modal Container -->
    <div
      class="learning-plan-modal bg-canvas-950 rounded-2xl border border-slate-800 w-full max-w-6xl max-h-[95vh] overflow-hidden shadow-2xl flex flex-col"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-6 border-b border-slate-800 bg-canvas-900/50">
        <div class="flex-1">
          <h2 class="text-2xl font-bold text-white mb-1">{plan.titulo}</h2>
          <p class="text-sm text-slate-400">{plan.descripcion}</p>
        </div>
        <button
          onclick={handleClose}
          class="ml-4 h-10 w-10 rounded-full bg-canvas-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-all"
          title="Cerrar (Esc)"
        >
          ✕
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto relative">
        {#if lesson}
          <!-- Lesson Player -->
          <LessonPlayer
            leccion={lesson}
            showProgress={true}
            onComplete={handleComplete}
          />
        {:else}
          <div class="flex items-center justify-center h-full">
            <p class="text-slate-400">No hay contenido disponible</p>
          </div>
        {/if}
      </div>

      <!-- Footer (Optional - Plan info) -->
      <div class="px-6 py-3 border-t border-slate-800 bg-canvas-900/30 flex items-center justify-between">
        <div class="flex items-center gap-4 text-xs text-slate-400">
          <div class="flex items-center gap-1.5">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>{plan.tiempo_estimado_minutos} min</span>
          </div>
          <div class="flex items-center gap-1.5">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>{plan.components.length} componentes</span>
          </div>
        </div>
        <div class="text-xs text-slate-500">
          Presiona <kbd class="px-2 py-1 bg-canvas-800 rounded border border-slate-700">Esc</kbd> para cerrar
        </div>
      </div>
    </div>
  </div>
{/if}
