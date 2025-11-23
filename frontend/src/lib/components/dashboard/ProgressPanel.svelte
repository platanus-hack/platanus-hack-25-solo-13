<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import type { Subject } from '$lib/constants/subjects';
  import SubjectProgressCard from './SubjectProgressCard.svelte';

  interface Props {
    isOpen: boolean;
    onClose: () => void;
    subjects: Subject[];
    userProfile: any;
    onSubjectClick: (subject: Subject) => void;
    diagnosticLevels: Record<number, number>;
  }

  let { isOpen, onClose, subjects, userProfile, onSubjectClick, diagnosticLevels }: Props = $props();

  // Get domain level for a subject
  function getDomainLevel(materiaId?: number): number {
    if (!materiaId) return 0;
    return diagnosticLevels[materiaId] || 0;
  }

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      onClose();
    }
  }

  // Animate panel and reload levels when opened
  $effect(() => {
    if (isOpen) {
      gsap.fromTo(
        '.progress-panel',
        { y: '100%', opacity: 0 },
        { y: 0, opacity: 1, duration: 0.4, ease: 'power3.out' }
      );
      gsap.fromTo(
        '.subject-card-panel',
        { scale: 0.8, opacity: 0 },
        { scale: 1, opacity: 1, duration: 0.4, stagger: 0.05, delay: 0.2, ease: 'power2.out' }
      );
    }
  });

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  // Calculate stats
  const totalSubjects = $derived(subjects.length);
  const evaluatedSubjects = $derived(
    subjects.filter(s => getDomainLevel(s.materiaId) > 0).length
  );
  const advancedSubjects = $derived(
    subjects.filter(s => getDomainLevel(s.materiaId) >= 3).length
  );
</script>

{#if isOpen}
  <div class="progress-panel fixed inset-0 z-40 bg-canvas-950 overflow-y-auto">
    <!-- Header -->
    <div class="sticky top-0 z-10 border-b border-slate-300 bg-canvas-950/90 backdrop-blur-md">
      <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 flex items-center gap-2">
            <span>📊</span>
            <span>Tu Progreso Académico</span>
          </h2>
          {#if userProfile?.curso_actual}
            <p class="text-sm text-slate-600 mt-1">
              {userProfile.curso_actual} • {evaluatedSubjects}/{totalSubjects} materias evaluadas
            </p>
          {/if}
        </div>
        <button
          onclick={onClose}
          class="h-10 w-10 rounded-full bg-canvas-900 hover:bg-canvas-800 border border-slate-700 hover:border-slate-600 flex items-center justify-center text-slate-400 hover:text-white transition-all"
          title="Cerrar (Esc)"
        >
          ✕
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-6 py-8">
      <!-- Summary Stats -->
      {#if evaluatedSubjects > 0}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8 max-w-3xl mx-auto">
          <div class="bg-canvas-900 rounded-xl p-4 border border-slate-700">
            <div class="text-3xl font-bold text-white">{evaluatedSubjects}/{totalSubjects}</div>
            <div class="text-sm text-slate-300 mt-1">Materias Evaluadas</div>
          </div>
          <div class="bg-canvas-900 rounded-xl p-4 border border-slate-700">
            <div class="text-3xl font-bold text-green-400">{advancedSubjects}</div>
            <div class="text-sm text-slate-300 mt-1">Nivel Intermedio o Superior</div>
          </div>
        </div>
      {/if}

      <!-- Subject Cards Grid -->
      <div>
        <h3 class="text-lg font-semibold text-slate-900 mb-6 text-center">Todas las Materias</h3>
        <div class="flex justify-center gap-8 flex-wrap max-w-4xl mx-auto">
          {#each subjects as subject}
            <div class="subject-card-panel w-80">
              <SubjectProgressCard
                {subject}
                domainLevel={getDomainLevel(subject.materiaId)}
                onClick={() => onSubjectClick(subject)}
              />
            </div>
          {/each}
        </div>
      </div>

      <!-- Empty State -->
      {#if subjects.length === 0}
        <div class="text-center py-12">
          <div class="text-6xl mb-4">📚</div>
          <h3 class="text-xl font-semibold text-slate-900 mb-2">No hay materias disponibles</h3>
          <p class="text-slate-600">
            Completa tu perfil para ver tus materias
          </p>
        </div>
      {/if}

      <!-- Footer hint -->
      <div class="mt-12 text-center text-sm text-slate-600">
        <p>Presiona <kbd class="px-2 py-1 bg-canvas-900 text-white rounded border border-slate-700">Esc</kbd> para cerrar</p>
      </div>
    </div>
  </div>
{/if}
