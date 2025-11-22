<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { mergeProfileData } from '$lib/api/profiles';
  import {
    PREGUNTAS_DIAGNOSTICO_LENGUA,
    calcularNivelDiagnostico,
    getFeedbackMensaje,
    type PreguntaDiagnostica
  } from '$lib/data/diagnostico-lengua';
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';
  import TrueFalse from '$lib/components/activities/TrueFalse.svelte';
  import FillBlanks from '$lib/components/activities/FillBlanks.svelte';
  import PlayerProfilePanel from '$lib/components/dashboard/PlayerProfilePanel.svelte';
  import RecentActivityModal from '$lib/components/dashboard/RecentActivityModal.svelte';
  import MissionBoardModal from '$lib/components/dashboard/MissionBoardModal.svelte';
  import CurrentQuestModal from '$lib/components/dashboard/CurrentQuestModal.svelte';
  import LiveEventsModal from '$lib/components/dashboard/LiveEventsModal.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';

  // Mapeo de IDs de materias a claves de conocimiento_previo
  const MATERIA_MAP: Record<string, { key: string; nombre: string; disponible: boolean }> = {
    'lyl': { key: 'lectura', nombre: 'Lengua y Literatura', disponible: true },
    'lenguaje': { key: 'lectura', nombre: 'Lengua y Literatura', disponible: true },
    'lengua': { key: 'lectura', nombre: 'Lengua y Literatura', disponible: true },
    'mat': { key: 'matematicas', nombre: 'Matemáticas', disponible: false },
    'matematicas': { key: 'matematicas', nombre: 'Matemáticas', disponible: false }
  };

  // Obtener parámetro de URL
  let materiaId = $state('');
  let materiaInfo = $state<{ key: string; nombre: string; disponible: boolean } | null>(null);
  let preguntas = $state<PreguntaDiagnostica[]>([]);

  // Dashboard State
  let isPlayerProfileOpen = $state(false);
  let isProgressPanelOpen = $state(false);
  let isActivityModalOpen = $state(false);
  let isMissionBoardOpen = $state(false);
  let isCurrentQuestOpen = $state(false);
  let isLiveEventsOpen = $state(false);
  let activeTab = $state('daily');
  let userProfile = $state(null);

  // Student data
  const student = {
    name: auth.user?.name || 'Student',
    grade: '2° Medio',
    level: 12,
    xp: 75,
    streak: 5
  };

  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  const resources = [
    { name: 'Mastery Tokens', value: 420, iconType: 'currency', color: 'text-white' },
    { name: 'Skill Points', value: 15, iconType: 'bolt', color: 'text-white' },
    { name: 'PAES Ready', value: '63%', iconType: 'chart', color: 'text-white' }
  ];

  // Load user profile and subjects
  async function loadUserProfile() {
    if (auth.user?.id) {
      const result = await getProfile(auth.user.id);
      if (result.success && result.profile) {
        userProfile = result.profile;
        const courseName = result.profile.curso_actual || '1ro Medio';
        const cursoResult = await findCursoByName(courseName);
        if (cursoResult.success && cursoResult.curso?.materias) {
          dashboardStore.updateSubjects(materiasToSubjects(cursoResult.curso.materias));
        }
      }
    }
  }

  // Get domain level for a subject
  function getDomainLevel(subjectId: string): number {
    if (!userProfile?.profile_data?.conocimiento_previo) {
      return 0;
    }
    const subjectMap: Record<string, string> = {
      'mat': 'matematicas',
      'lyl': 'lectura',
      'lenguaje': 'lectura',
      'lengua': 'lectura'
    };
    const key = subjectMap[subjectId.toLowerCase()] || subjectId.toLowerCase();
    return userProfile.profile_data.conocimiento_previo[key]?.nivel || 0;
  }

  // Open subject detail modal
  let selectedSubject = $state<Subject | null>(null);
  let isModalOpen = $state(false);
  let selectedDomainLevel = $state(0);

  function openSubjectDetail(subject: Subject) {
    selectedSubject = subject;
    selectedDomainLevel = getDomainLevel(subject.id);
    isModalOpen = true;
  }

  function closeModal() {
    isModalOpen = false;
    selectedSubject = null;
  }

  // Evaluation State
  let currentQuestionIndex = $state(0);
  let answers = $state<Record<string, boolean>>({}); // preguntaId -> esCorrecta
  let hasAnswered = $state(false);
  let isLoading = $state(false);
  let showResults = $state(false);
  let diagnosticoCompletado = $state(false);

  // Derived
  const totalPreguntas = $derived(preguntas.length);
  const preguntaActual = $derived(preguntas[currentQuestionIndex]);
  const progreso = $derived(totalPreguntas > 0 ? Math.round(((currentQuestionIndex + 1) / totalPreguntas) * 100) : 0);
  const respuestasCorrectas = $derived(Object.values(answers).filter(Boolean).length);
  const nivelCalculado = $derived(calcularNivelDiagnostico(respuestasCorrectas, totalPreguntas));

  // Verificar autenticación e inicializar
  onMount(() => {
    if (!auth.checkAuth()) {
      goto('/login');
      return;
    }

    // Load user profile and subjects for navigation
    loadUserProfile();

    // Obtener materiaId de la URL
    const unsubscribe = page.subscribe(($page) => {
      materiaId = $page.params.materiaId || '';

      // Normalizar y validar materia
      const normalizedId = materiaId.toLowerCase();
      materiaInfo = MATERIA_MAP[normalizedId] || null;

      if (!materiaInfo) {
        // Materia no reconocida
        alert('Materia no reconocida');
        goto('/');
        return;
      }

      if (!materiaInfo.disponible) {
        // Materia no tiene evaluación disponible aún
        alert(`La evaluación diagnóstica de ${materiaInfo.nombre} estará disponible próximamente.`);
        goto('/');
        return;
      }

      // Cargar banco de preguntas según materia
      if (materiaInfo.key === 'lectura') {
        preguntas = PREGUNTAS_DIAGNOSTICO_LENGUA;
      } else if (materiaInfo.key === 'matematicas') {
        // TODO: Cargar PREGUNTAS_DIAGNOSTICO_MATEMATICAS cuando esté disponible
        preguntas = [];
      }
    });

    return () => unsubscribe();
  });

  // Manejar respuesta de activity component
  function handleAnswer(data: { isCorrect: boolean }) {
    answers[preguntaActual.id] = data.isCorrect;
    hasAnswered = true;
  }

  // Siguiente pregunta
  function siguientePregunta() {
    if (!hasAnswered) return;

    if (currentQuestionIndex < totalPreguntas - 1) {
      currentQuestionIndex++;
      hasAnswered = false;

      // Animación al cambiar pregunta
      gsap.fromTo(
        '.question-container',
        { x: 50, opacity: 0 },
        { x: 0, opacity: 1, duration: 0.4, ease: 'power2.out' }
      );
    } else {
      // Última pregunta - mostrar resultados
      mostrarResultados();
    }
  }

  // Pregunta anterior
  function preguntaAnterior() {
    if (currentQuestionIndex > 0) {
      currentQuestionIndex--;
      hasAnswered = true; // Ya fue respondida
    }
  }

  // Mostrar resultados
  function mostrarResultados() {
    showResults = true;

    gsap.from('.results-container', {
      opacity: 0,
      scale: 0.9,
      duration: 0.6,
      ease: 'back.out(1.7)'
    });
  }

  // Guardar nivel en perfil
  async function guardarNivelEnPerfil() {
    if (!auth.user?.id || diagnosticoCompletado || !materiaInfo) return;

    isLoading = true;

    try {
      const profileDataUpdate = {
        conocimiento_previo: {
          [materiaInfo.key]: {
            nivel: nivelCalculado.nivel,
            fuente: 'diagnostico_inicial',
            fecha_evaluacion: new Date().toISOString(),
            porcentaje_acierto: nivelCalculado.porcentaje
          }
        }
      };

      const result = await mergeProfileData(auth.user.id, profileDataUpdate);

      if (result.success) {
        diagnosticoCompletado = true;
        setTimeout(() => {
          goto('/');
        }, 3000);
      } else {
        alert('Error al guardar el resultado: ' + (result.error || 'Error desconocido'));
      }
    } catch (error) {
      console.error('Error guardando nivel:', error);
      alert('Error al guardar el resultado');
    } finally {
      isLoading = false;
    }
  }

  // Transformar pregunta a formato de componente
  function transformarPreguntaMultipleChoice(pregunta: PreguntaDiagnostica) {
    if (pregunta.tipo !== 'multiple_choice') return null;

    return pregunta.opciones!.map((text, index) => ({
      id: index,
      text,
      isCorrect: index === pregunta.respuestaCorrecta
    }));
  }
</script>

<svelte:head>
  <title>Evaluación Diagnóstica - {materiaInfo?.nombre || 'Cargando...'}</title>
</svelte:head>

<div class="min-h-screen bg-canvas-950 text-slate-200 font-sans pb-10">
  <!-- Top Player Bar -->
  <header class="sticky top-0 z-50 border-b border-white/5 bg-canvas-950/90 backdrop-blur-md">
    <div class="px-6 py-3 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <!-- Profile -->
      <button
        onclick={() => isPlayerProfileOpen = true}
        class="flex items-center gap-3 hover:bg-canvas-800/40 rounded-xl p-2 -m-2 transition-colors group"
      >
        <div class="h-10 w-10 rounded-full bg-gradient-to-br from-lumera-500 to-focus-600 flex items-center justify-center text-sm font-bold group-hover:scale-110 transition-transform">
          {initials}
        </div>
        <div class="text-left">
          <div class="flex items-center gap-2">
            <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
            <span class="px-1.5 py-0.5 rounded bg-canvas-800 border border-canvas-700 text-xs text-slate-400">{student.grade}</span>
          </div>
          <div class="text-xs text-slate-400">{auth.user?.email || 'student@lumera.com'}</div>
        </div>
      </button>

      <!-- Evaluation Progress (replaces XP bar) -->
      <div class="flex-1 max-w-md px-2 hidden md:block">
        {#if materiaInfo && totalPreguntas > 0}
          {#if showResults}
            <div class="flex justify-between text-xs mb-1">
              <span class="text-green-400 font-medium">✓ Evaluación Completada</span>
              <span class="text-slate-400">{nivelCalculado.label} - {nivelCalculado.porcentaje}%</span>
            </div>
            <div class="h-2 w-full bg-canvas-900 rounded-full overflow-hidden">
              <div class="h-full bg-gradient-to-r from-green-500 to-emerald-500" style="width: 100%"></div>
            </div>
          {:else}
            <div class="flex justify-between text-xs mb-1">
              <span class="text-focus-400 font-medium">📚 {materiaInfo.nombre}</span>
              <span class="text-slate-400">Pregunta {currentQuestionIndex + 1} de {totalPreguntas}</span>
            </div>
            <div class="h-2 w-full bg-canvas-900 rounded-full overflow-hidden">
              <div class="h-full bg-gradient-to-r from-focus-500 to-blue-500 transition-all duration-500" style="width: {progreso}%"></div>
            </div>
          {/if}
        {:else}
          <div class="flex justify-between text-xs mb-1">
            <span class="text-slate-400 font-medium">Cargando evaluación...</span>
          </div>
          <div class="h-2 w-full bg-canvas-900 rounded-full overflow-hidden">
            <div class="h-full bg-slate-700" style="width: 0%"></div>
          </div>
        {/if}
      </div>

      <!-- Right - Navigation Icons -->
      <div class="flex items-center gap-2">
        <!-- Quest Button -->
        <button
          onclick={() => isCurrentQuestOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Current Quest"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
          </svg>
        </button>

        <!-- Missions Button -->
        <button
          onclick={() => isMissionBoardOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Mission Board"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          {#if dashboardStore.activeMissionCount > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-purple-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.activeMissionCount}
            </span>
          {/if}
        </button>

        <!-- Activity Button -->
        <button
          onclick={() => isActivityModalOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Actividad Reciente"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          {#if dashboardStore.activities.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-rose-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.activities.length}
            </span>
          {/if}
        </button>

        <!-- Events Button -->
        <button
          onclick={() => isLiveEventsOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Live Events"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
          {#if dashboardStore.events.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-amber-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.events.length}
            </span>
          {/if}
        </button>

        <!-- Progress Button -->
        <button
          onclick={() => isProgressPanelOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Tu Progreso"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          {#if dashboardStore.subjects.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-lumera-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.subjects.length}
            </span>
          {/if}
        </button>

        <!-- Logout Button -->
        <button
          onclick={() => auth.logout()}
          class="relative p-2 rounded-lg hover:bg-red-900/60 transition-all duration-200 group"
          title="Logout"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-red-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>
    </div>
  </header>

  <!-- Main Content -->
  <main class="px-6 py-8 max-w-4xl mx-auto">
    {#if materiaInfo && preguntas.length > 0}
      <div class="diagnostico-container">
      {#if !showResults}
        <!-- Título de la evaluación -->
        <div class="mb-6">
          <h1 class="text-2xl font-bold text-white">Evaluación Diagnóstica</h1>
          <p class="text-slate-400 text-sm mt-1">Responde todas las preguntas para conocer tu nivel inicial</p>
        </div>

        <!-- Pregunta actual -->
        <div class="question-container mb-6">
          {#key preguntaActual.id}
            {#if preguntaActual.tipo === 'multiple_choice'}
              <MultipleChoice
                question={preguntaActual.pregunta}
                options={transformarPreguntaMultipleChoice(preguntaActual) || []}
                bloomLevel="comprender"
                materia="lenguaje"
                showFeedback={true}
                allowMultipleAttempts={false}
                showCorrectAnswer={true}
                onAnswer={handleAnswer}
              />
            {:else if preguntaActual.tipo === 'true_false'}
              <TrueFalse
                statement={preguntaActual.pregunta}
                correctAnswer={preguntaActual.respuestaCorrecta === 0}
                explanation={preguntaActual.explicacion}
                bloomLevel="comprender"
                materia="lenguaje"
                showExplanation={true}
                allowMultipleAttempts={false}
                onAnswer={handleAnswer}
              />
            {:else if preguntaActual.tipo === 'fill_blanks'}
              <FillBlanks
                text={preguntaActual.pregunta}
                blanks={[{
                  id: 1,
                  answer: String(preguntaActual.respuestaCorrecta),
                  caseSensitive: false
                }]}
                bloomLevel="comprender"
                materia="lenguaje"
                allowMultipleAttempts={false}
                onAnswer={handleAnswer}
              />
            {/if}
          {/key}
        </div>

        <!-- Navegación -->
        <div class="flex justify-between items-center">
          <button
            onclick={preguntaAnterior}
            disabled={currentQuestionIndex === 0}
            class="px-6 py-3 rounded-xl bg-canvas-800 hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed border border-slate-700 transition-all"
          >
            ← Anterior
          </button>

          <div class="text-center">
            <p class="text-sm text-slate-400">
              {hasAnswered ? 'Respuesta registrada' : 'Selecciona una respuesta'}
            </p>
          </div>

          <button
            onclick={siguientePregunta}
            disabled={!hasAnswered}
            class="px-6 py-3 rounded-xl bg-gradient-to-r from-focus-500 to-blue-500 hover:shadow-lg hover:shadow-cyan-500/50 disabled:opacity-50 disabled:cursor-not-allowed font-semibold transition-all"
          >
            {currentQuestionIndex < totalPreguntas - 1 ? 'Siguiente →' : 'Ver Resultados'}
          </button>
        </div>
      {:else}
        <!-- Pantalla de resultados -->
        <div class="results-container">
          <!-- Resultado general -->
          <div class="text-center mb-8">
            <div class="inline-block p-8 rounded-2xl bg-gradient-to-br from-cyan-600/20 to-blue-600/20 border border-cyan-500/50 mb-6">
              <div class="text-6xl mb-4">
                {nivelCalculado.nivel === 0 ? '📖' :
                 nivelCalculado.nivel === 1 ? '📕' :
                 nivelCalculado.nivel === 2 ? '📗' :
                 nivelCalculado.nivel === 3 ? '📘' : '📙'}
              </div>
              <h2 class="text-4xl font-bold text-white mb-2">{nivelCalculado.label}</h2>
              <p class="text-2xl text-focus-400">Nivel {nivelCalculado.nivel}</p>
            </div>

            <div class="max-w-2xl mx-auto space-y-4">
              <div class="p-6 bg-canvas-900/60 rounded-2xl border border-slate-800">
                <div class="grid grid-cols-3 gap-4 text-center">
                  <div>
                    <div class="text-3xl font-bold text-white">{respuestasCorrectas}</div>
                    <div class="text-sm text-slate-400">Correctas</div>
                  </div>
                  <div>
                    <div class="text-3xl font-bold text-white">{totalPreguntas - respuestasCorrectas}</div>
                    <div class="text-sm text-slate-400">Incorrectas</div>
                  </div>
                  <div>
                    <div class="text-3xl font-bold text-focus-400">{nivelCalculado.porcentaje}%</div>
                    <div class="text-sm text-slate-400">Acierto</div>
                  </div>
                </div>
              </div>

              <div class="p-6 bg-canvas-900/60 rounded-2xl border border-slate-800 text-left">
                <h3 class="text-lg font-semibold text-white mb-3 flex items-center gap-2">
                  <span>💬</span> Retroalimentación
                </h3>
                <p class="text-slate-300 leading-relaxed">
                  {getFeedbackMensaje(nivelCalculado.nivel)}
                </p>
              </div>
            </div>
          </div>

          <!-- Acciones -->
          <div class="flex justify-center gap-4">
            {#if !diagnosticoCompletado}
              <button
                onclick={guardarNivelEnPerfil}
                disabled={isLoading}
                class="px-8 py-4 rounded-xl bg-gradient-to-r from-focus-500 to-blue-500 hover:shadow-lg hover:shadow-cyan-500/50 font-semibold text-lg disabled:opacity-50 disabled:cursor-not-allowed transition-all"
              >
                {isLoading ? 'Guardando...' : 'Guardar y Continuar'}
              </button>
            {:else}
              <div class="text-center">
                <div class="text-green-400 text-2xl mb-2">✓ Nivel guardado</div>
                <p class="text-slate-400">Redirigiendo al dashboard...</p>
              </div>
            {/if}
          </div>
        </div>
      {/if}
      </div>
    {:else}
      <div class="flex items-center justify-center min-h-[60vh]">
        <div class="text-center">
          <div class="text-6xl mb-4">⏳</div>
          <p class="text-slate-400">Cargando evaluación...</p>
        </div>
      </div>
    {/if}
  </main>
</div>

<!-- Modals -->
<RecentActivityModal
  activities={dashboardStore.activities}
  isOpen={isActivityModalOpen}
  onClose={() => isActivityModalOpen = false}
/>

<MissionBoardModal
  missions={dashboardStore.missions}
  {activeTab}
  onTabChange={(tab) => activeTab = tab}
  isOpen={isMissionBoardOpen}
  onClose={() => isMissionBoardOpen = false}
/>

<CurrentQuestModal
  isOpen={isCurrentQuestOpen}
  onClose={() => isCurrentQuestOpen = false}
  onResumeQuest={() => console.log('Resume quest clicked')}
/>

<LiveEventsModal
  events={dashboardStore.events}
  isOpen={isLiveEventsOpen}
  onClose={() => isLiveEventsOpen = false}
/>

<ProgressPanel
  isOpen={isProgressPanelOpen}
  onClose={() => isProgressPanelOpen = false}
  subjects={dashboardStore.subjects}
  {userProfile}
  onSubjectClick={openSubjectDetail}
/>

<PlayerProfilePanel
  isOpen={isPlayerProfileOpen}
  onClose={() => isPlayerProfileOpen = false}
  userName={student.name}
  userEmail={auth.user?.email || 'student@lumera.com'}
  userGrade={student.grade}
  level={student.level}
  xp={student.xp}
  {initials}
/>

<!-- Subject Detail Modal -->
{#if selectedSubject}
  <div class="fixed inset-0 z-70" style="display: {isModalOpen ? 'block' : 'none'}">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm" onclick={closeModal}></div>
    <div class="fixed inset-0 flex items-center justify-center p-4 pointer-events-none">
      <div class="pointer-events-auto">
        <!-- Subject detail content would go here if needed during evaluation -->
      </div>
    </div>
  </div>
{/if}
