<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import {
    startDiagnosticSession,
    getNextQuestion,
    submitAnswer,
    completeSession,
    getResults,
    type DiagnosticSession,
    type DiagnosticQuestion,
    type DiagnosticResult
  } from '$lib/api/diagnostic';
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';
  import TrueFalse from '$lib/components/activities/TrueFalse.svelte';
  import FillBlanks from '$lib/components/activities/FillBlanks.svelte';
  import DragDropMatching from '$lib/components/activities/DragDropMatching.svelte';
  import CompareContrast from '$lib/components/activities/CompareContrast.svelte';
  import ConceptMapBuilder from '$lib/components/activities/ConceptMapBuilder.svelte';
  import CriteriaEvaluation from '$lib/components/activities/CriteriaEvaluation.svelte';
  import OpenEndedResponse from '$lib/components/activities/OpenEndedResponse.svelte';
  import Sequencing from '$lib/components/activities/Sequencing.svelte';
  import PlayerProfilePanel from '$lib/components/dashboard/PlayerProfilePanel.svelte';
  import RecentActivityModal from '$lib/components/dashboard/RecentActivityModal.svelte';
  import MissionBoardModal from '$lib/components/dashboard/MissionBoardModal.svelte';
  import CurrentQuestModal from '$lib/components/dashboard/CurrentQuestModal.svelte';
  import LiveEventsModal from '$lib/components/dashboard/LiveEventsModal.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import { getEquipment, type CustomizationItem } from '$lib/api/customization';

  // Obtener parámetro de URL
  let materiaSlug = $state('');
  let materiaInfo = $state<{ id: number; nombre: string; codigo: string } | null>(null);

  // Fetch materia by slug/code from backend
  async function fetchMateriaByCodigo(slug: string) {
    try {
      // Normalize slug to codigo (lyl -> LYL, mat -> MAT, etc)
      const codigoMap: Record<string, string> = {
        'lyl': 'LYL',
        'lenguaje': 'LYL',
        'lengua': 'LYL',
        'mat': 'MAT',
        'matematicas': 'MAT'
      };

      const codigo = codigoMap[slug.toLowerCase()] || slug.toUpperCase();

      const response = await fetch('/api/materias', {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });

      if (!response.ok) {
        throw new Error('Error fetching materias');
      }

      const materias = await response.json();
      const materia = materias.find((m: any) => m.codigo === codigo);

      return materia || null;
    } catch (error) {
      console.error('Error fetching materia:', error);
      return null;
    }
  }

  // Diagnostic Session State
  let session = $state<DiagnosticSession | null>(null);
  let currentQuestion = $state<DiagnosticQuestion | null>(null);
  let currentQuestionNumber = $state(0);
  let totalQuestions = $state(5); // Always 5 questions

  // Dashboard State
  let isPlayerProfileOpen = $state(false);
  let isProgressPanelOpen = $state(false);
  let isActivityModalOpen = $state(false);
  let isMissionBoardOpen = $state(false);
  let isCurrentQuestOpen = $state(false);
  let isLiveEventsOpen = $state(false);
  let activeTab = $state('daily');
  let userProfile = $state(null);
  let currentAvatar = $state<CustomizationItem | null>(null);

  // UI State
  let hasAnswered = $state(false);
  let isLoading = $state(false);
  let isCorrect = $state(false);
  let showResults = $state(false);
  let results = $state<DiagnosticResult[]>([]);
  let averageBloomLevel = $state(0);
  let errorMessage = $state('');

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

  // Load user profile and subjects for navigation
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

  // Load equipped avatar
  async function loadEquippedAvatar() {
    try {
      const equipment = await getEquipment();
      currentAvatar = equipment.equipped_avatar || null;
    } catch (err) {
      console.error('Error loading avatar:', err);
    }
  }

  // Handle avatar change from profile panel
  function handleAvatarChanged(avatar: CustomizationItem) {
    currentAvatar = avatar;
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

  // Subject detail modal state
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

  // Progress calculation
  const progreso = $derived(currentQuestionNumber > 0 && totalQuestions > 0
    ? Math.round((currentQuestionNumber / totalQuestions) * 100)
    : 0);

  // Initialize diagnostic
  async function initializeDiagnostic() {
    if (!materiaInfo) return;

    isLoading = true;
    errorMessage = '';

    // Start diagnostic session
    const sessionResult = await startDiagnosticSession(materiaInfo.id);

    if (!sessionResult.success) {
      errorMessage = sessionResult.error || 'Error al iniciar la evaluación';
      isLoading = false;
      return;
    }

    session = sessionResult.session!;

    // Load first question
    await loadNextQuestion();

    isLoading = false;
  }

  // Load next question
  async function loadNextQuestion() {
    if (!session) return;

    isLoading = true;
    errorMessage = '';

    console.log('[Diagnostic] Fetching next question for session:', session.id);
    const questionResult = await getNextQuestion(session.id);
    console.log('[Diagnostic] Question result:', questionResult);

    if (!questionResult.success) {
      errorMessage = questionResult.error || 'No hay más preguntas disponibles';
      isLoading = false;
      // If no more questions, complete the session
      if (currentQuestionNumber > 0) {
        await finalizeDiagnostic();
      }
      return;
    }

    currentQuestion = questionResult.question!;
    currentQuestionNumber = currentQuestion.question_number;
    totalQuestions = currentQuestion.total_questions;
    hasAnswered = false;
    isCorrect = false;
    isLoading = false;

    console.log('[Diagnostic] Current question loaded:', {
      type: currentQuestion.tipo,
      number: currentQuestionNumber,
      total: totalQuestions,
      data: currentQuestion.question_data
    });
  }

  // Handle answer submission
  async function handleAnswer(data: { isCorrect: boolean; selectedOption?: any; answer?: any }) {
    if (hasAnswered || !session || !currentQuestion) {
      console.log('[Diagnostic] Skipping answer - already answered or missing data', {
        hasAnswered, session: !!session, currentQuestion: !!currentQuestion
      });
      return;
    }

    hasAnswered = true;
    isCorrect = data.isCorrect;
    isLoading = true;

    // Determine user answer based on question type
    let userAnswer: any;
    if (data.selectedOption !== undefined) {
      userAnswer = { selected: data.selectedOption }; // Backend expects 'selected' field
    } else if (data.answer !== undefined) {
      userAnswer = data.answer; // Send answer object directly (e.g., { matches: {...} })
    } else {
      userAnswer = { is_correct: data.isCorrect };
    }

    console.log('[Diagnostic] Submitting answer:', {
      sessionId: session.id,
      questionId: currentQuestion.id,
      userAnswer
    });

    // Submit answer to backend
    const answerResult = await submitAnswer(
      session.id,
      currentQuestion.id,
      userAnswer
    );

    console.log('[Diagnostic] Answer result:', answerResult);

    isLoading = false;

    if (!answerResult.success) {
      errorMessage = answerResult.error || 'Error al enviar la respuesta';
      console.error('[Diagnostic] Error submitting answer:', errorMessage);
      return;
    }

    // Show feedback briefly before moving to next question
    setTimeout(() => {
      siguientePregunta();
    }, 1500);
  }

  // Move to next question
  async function siguientePregunta() {
    if (currentQuestionNumber >= totalQuestions) {
      // Finished all questions, complete session
      await finalizeDiagnostic();
    } else {
      // Load next question
      await loadNextQuestion();
    }
  }

  // Complete diagnostic session
  async function finalizeDiagnostic() {
    if (!session) return;

    isLoading = true;
    errorMessage = '';

    const completeResult = await completeSession(session.id);

    if (!completeResult.success) {
      errorMessage = completeResult.error || 'Error al completar la evaluación';
      isLoading = false;
      return;
    }

    averageBloomLevel = completeResult.result!.average_bloom_level;

    // Get detailed results
    const resultsResponse = await getResults(session.id);

    if (resultsResponse.success) {
      results = resultsResponse.results!;
    }

    showResults = true;
    isLoading = false;
  }

  // Transform question data for components
  function transformarPreguntaMultipleChoice(question: DiagnosticQuestion) {
    const questionData = question.question_data;

    // Convert opciones object {A, B, C, D} to array of option objects
    let optionsArray: Array<{id: number, text: string, isCorrect: boolean}> = [];
    if (questionData.opciones) {
      if (Array.isArray(questionData.opciones)) {
        optionsArray = questionData.opciones.map((text, index) => ({
          id: index,
          text: text,
          isCorrect: false // We don't validate locally, backend handles it
        }));
      } else {
        // Convert {A: "...", B: "...", C: "...", D: "..."} to array
        const values = Object.values(questionData.opciones);
        optionsArray = values.map((text: any, index) => ({
          id: index,
          text: text,
          isCorrect: false // We don't validate locally, backend handles it
        }));
      }
    }

    return {
      question: questionData.pregunta || questionData.question || '',
      options: optionsArray,
      showCorrectAnswer: false, // Don't show correct answer indicator (backend validates)
      showFeedback: false // Don't show immediate feedback (wait for backend)
    };
  }

  function transformarPreguntaTrueFalse(question: DiagnosticQuestion) {
    const questionData = question.question_data;
    return {
      pregunta: questionData.pregunta || questionData.question || '',
      respuestaCorrecta: questionData.respuesta_correcta || questionData.correct_answer || false,
      explicacion: questionData.explicacion || questionData.explanation || ''
    };
  }

  function transformarPreguntaFillBlanks(question: DiagnosticQuestion) {
    const questionData = question.question_data;
    return {
      texto: questionData.texto || questionData.text || '',
      blancos: questionData.blancos || questionData.blanks || [],
      respuestasCorrectas: questionData.respuestas_correctas || questionData.correct_answers || [],
      explicacion: questionData.explicacion || questionData.explanation || ''
    };
  }

  function transformarPreguntaDragDropMatching(question: DiagnosticQuestion) {
    const questionData = question.question_data;

    // Transform columna_izquierda/columna_derecha format to pairs array
    let pairs = [];
    if (questionData.columna_izquierda && questionData.columna_derecha) {
      const leftColumn = questionData.columna_izquierda;
      const rightColumn = questionData.columna_derecha;

      pairs = leftColumn.map((term: string, index: number) => ({
        id: index + 1,
        term: term,
        definition: rightColumn[index] || ''
      }));
    } else if (questionData.pares || questionData.pairs) {
      pairs = questionData.pares || questionData.pairs;
    }

    return {
      title: questionData.instruccion || questionData.titulo || questionData.title || 'Relaciona los conceptos',
      pairs: pairs,
      shuffleOptions: true,
      showFeedback: false,
      allowMultipleAttempts: false
    };
  }

  function transformarPreguntaSequencing(question: DiagnosticQuestion) {
    const questionData = question.question_data;

    // Transform elementos_desordenados to items array with correct order
    let items = [];
    if (questionData.elementos_desordenados && Array.isArray(questionData.elementos_desordenados)) {
      // Create items with sequential IDs - backend will validate the actual order
      items = questionData.elementos_desordenados.map((content: string, index: number) => ({
        id: index + 1,
        content: content,
        correctOrder: index + 1 // Placeholder - backend validates actual correct order
      }));
    } else if (questionData.items || questionData.elementos) {
      items = questionData.items || questionData.elementos;
    }

    return {
      title: questionData.instruccion || questionData.titulo || questionData.title || 'Ordena la secuencia',
      items: items,
      shuffleItems: true,
      showNumbers: true,
      showHints: false,
      allowMultipleAttempts: false
    };
  }

  function transformarPreguntaCompareContrast(question: DiagnosticQuestion) {
    const questionData = question.question_data;
    return {
      title: questionData.titulo || questionData.title || 'Compara y contrasta',
      itemA: questionData.itemA || questionData.item_a || { name: 'Item A', color: 'cyan' },
      itemB: questionData.itemB || questionData.item_b || { name: 'Item B', color: 'purple' },
      characteristics: questionData.caracteristicas || questionData.characteristics || [],
      showFeedback: false,
      allowMultipleAttempts: false
    };
  }

  function transformarPreguntaCriteriaEvaluation(question: DiagnosticQuestion) {
    const questionData = question.question_data;
    return {
      title: questionData.titulo || questionData.title || 'Evalúa según criterios',
      subject: questionData.sujeto || questionData.subject || '',
      description: questionData.descripcion || questionData.description || '',
      content: questionData.contenido || questionData.content || null,
      criteria: questionData.criterios || questionData.criteria || [],
      showFeedback: false,
      allowMultipleAttempts: false,
      showExpectedRatings: false
    };
  }

  function transformarPreguntaOpenEnded(question: DiagnosticQuestion) {
    const questionData = question.question_data;
    return {
      prompt: questionData.pregunta || questionData.prompt || '',
      placeholder: questionData.placeholder || 'Escribe tu respuesta aquí...',
      minWords: questionData.min_palabras || questionData.minWords || 0,
      maxWords: questionData.max_palabras || questionData.maxWords || 0,
      showWordCount: true,
      enableAiFeedback: false,
      rubric: questionData.rubrica || questionData.rubric || []
    };
  }

  function transformarPreguntaConceptMap(question: DiagnosticQuestion) {
    const questionData = question.question_data;
    return {
      title: questionData.titulo || questionData.title || 'Crea un mapa conceptual',
      topic: questionData.tema || questionData.topic || '',
      instructions: questionData.instrucciones || questionData.instructions || '',
      requiredConcepts: questionData.conceptos_requeridos || questionData.requiredConcepts || [],
      suggestedConnections: questionData.conexiones_sugeridas || questionData.suggestedConnections || [],
      minConcepts: questionData.min_conceptos || questionData.minConcepts || 3,
      minConnections: questionData.min_conexiones || questionData.minConnections || 2,
      showFeedback: false,
      allowMultipleAttempts: false,
      provideConcepts: false
    };
  }

  // Map Bloom level to display category
  function getBloomLevelDisplay(bloomLevel: number) {
    const bloomMap: Record<number, { label: string; color: string; porcentaje: number }> = {
      0: { label: 'Sin Evaluar', color: 'text-slate-400', porcentaje: 0 },
      1: { label: 'Recordar', color: 'text-red-400', porcentaje: 16 },
      2: { label: 'Comprender', color: 'text-orange-400', porcentaje: 33 },
      3: { label: 'Aplicar', color: 'text-yellow-400', porcentaje: 50 },
      4: { label: 'Analizar', color: 'text-green-400', porcentaje: 66 },
      5: { label: 'Evaluar', color: 'text-blue-400', porcentaje: 83 },
      6: { label: 'Crear', color: 'text-purple-400', porcentaje: 100 }
    };
    return bloomMap[bloomLevel] || bloomMap[0];
  }

  // Initialize on mount
  onMount(async () => {
    loadUserProfile();
    loadEquippedAvatar();

    const unsubscribe = page.subscribe(($page) => {
      materiaSlug = $page.params.materiaId || '';
    });

    // Fetch materia from backend
    materiaInfo = await fetchMateriaByCodigo(materiaSlug);

    if (!materiaInfo) {
      errorMessage = 'Materia no reconocida. Por favor verifica la URL.';
      return () => {
        if (unsubscribe) unsubscribe();
      };
    }

    // Initialize diagnostic session
    await initializeDiagnostic();

    return () => {
      if (unsubscribe) unsubscribe();
    };
  });
</script>

<svelte:head>
  <title>Evaluación Diagnóstica - {materiaInfo?.nombre || 'Lumera App'}</title>
</svelte:head>

<!-- Main Layout -->
<div class="min-h-screen bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950">
  <!-- Header -->
  <header class="sticky top-0 z-50 border-b border-white/5 bg-canvas-950/90 backdrop-blur-md">
    <div class="px-6 py-3 flex items-center justify-between gap-4">
      <!-- Left: Profile -->
      <button
        onclick={() => isPlayerProfileOpen = true}
        class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-canvas-800/40 transition-all group"
      >
        <AvatarDisplay
          {currentAvatar}
          {initials}
          size="small"
        />
        <div class="text-left hidden md:block">
          <div class="text-sm font-semibold text-white group-hover:text-lumera-400 transition-colors">
            {student.name}
          </div>
          <div class="text-xs text-slate-400">{student.grade}</div>
        </div>
      </button>

      <!-- Center: Progress (replaces XP bar during evaluation) -->
      <div class="flex-1 max-w-md">
        <div class="flex items-center justify-between text-xs mb-1.5">
          {#if showResults}
            <span class="text-green-400 font-medium">✓ Evaluación Completada</span>
            <span class="text-slate-400">{getBloomLevelDisplay(averageBloomLevel).label} - {getBloomLevelDisplay(averageBloomLevel).porcentaje}%</span>
          {:else if currentQuestion}
            <span class="text-focus-400 font-medium">📚 {materiaInfo?.nombre}</span>
            <span class="text-slate-400">Pregunta {currentQuestionNumber} de {totalQuestions}</span>
          {:else}
            <span class="text-slate-400">Cargando...</span>
            <span class="text-slate-400">--</span>
          {/if}
        </div>
        <div class="h-2 w-full bg-canvas-900 rounded-full overflow-hidden border border-slate-700">
          {#if showResults}
            <div class="h-full bg-gradient-to-r from-green-500 to-emerald-500 transition-all duration-500" style="width: 100%"></div>
          {:else}
            <div class="h-full bg-gradient-to-r from-focus-500 to-blue-500 transition-all duration-500" style="width: {progreso}%"></div>
          {/if}
        </div>
      </div>

      <!-- Right: Navigation Icons -->
      <div class="flex items-center gap-2">
        <!-- Quest Button -->
        <button
          onclick={() => isCurrentQuestOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Quest Actual"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        </button>

        <!-- Mission Board Button -->
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
    <!-- Error Message -->
    {#if errorMessage}
      <div class="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-xl">
        <p class="text-red-400">{errorMessage}</p>
      </div>
    {/if}

    <!-- Loading State -->
    {#if isLoading && !currentQuestion}
      <div class="flex items-center justify-center py-20">
        <div class="text-center">
          <svg class="animate-spin h-12 w-12 mx-auto mb-4 text-lumera-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="text-slate-400">Cargando evaluación...</p>
        </div>
      </div>
    {/if}

    <!-- Question Display -->
    {#if !showResults && currentQuestion && !isLoading}
      <div class="bg-canvas-900/40 backdrop-blur-xl rounded-2xl border border-white/10 p-8">
        {#key currentQuestion.id}
          {#if currentQuestion.tipo === 'multiple_choice'}
            {@const transformedQuestion = transformarPreguntaMultipleChoice(currentQuestion)}
            <MultipleChoice
              question={transformedQuestion.question}
              options={transformedQuestion.options}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              showCorrectAnswer={false}
              showFeedback={false}
              allowMultipleAttempts={false}
              onAnswer={(data) => handleAnswer({ isCorrect: true, selectedOption: data.userAnswer })}
            />
          {:else if currentQuestion.tipo === 'true_false'}
            {@const transformedQuestion = transformarPreguntaTrueFalse(currentQuestion)}
            <TrueFalse
              statement={transformedQuestion.pregunta}
              correctAnswer={transformedQuestion.respuestaCorrecta}
              explanation={transformedQuestion.explicacion}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showExplanation={false}
              onAnswer={(data) => handleAnswer({ isCorrect: data.isCorrect })}
            />
          {:else if currentQuestion.tipo === 'fill_blanks'}
            {@const transformedQuestion = transformarPreguntaFillBlanks(currentQuestion)}
            <FillBlanks
              text={transformedQuestion.texto}
              blanks={transformedQuestion.blancos}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showWordBank={false}
              showHints={false}
              onAnswer={(data) => handleAnswer({ isCorrect: data.isCorrect })}
            />
          {:else if currentQuestion.tipo === 'drag_drop_matching'}
            {@const transformedQuestion = transformarPreguntaDragDropMatching(currentQuestion)}
            <DragDropMatching
              title={transformedQuestion.title}
              pairs={transformedQuestion.pairs}
              bloomLevel="comprender"
              materia={materiaInfo?.nombre || ''}
              shuffleOptions={transformedQuestion.shuffleOptions}
              showFeedback={transformedQuestion.showFeedback}
              allowMultipleAttempts={transformedQuestion.allowMultipleAttempts}
              onAnswer={(data) => {
                // Transform matches from { termId: definitionId } to { "term": "definition" }
                const textMatches: Record<string, string> = {};
                const pairs = transformedQuestion.pairs;

                for (const [termId, defId] of Object.entries(data.matches)) {
                  const termIdNum = parseInt(termId);
                  const defIdNum = parseInt(defId as string);

                  // Find the term and definition by ID
                  const pair = pairs.find((p: any) => p.id === termIdNum);
                  const defPair = pairs.find((p: any) => p.id === defIdNum);

                  if (pair && defPair) {
                    textMatches[pair.term] = defPair.definition;
                  }
                }

                handleAnswer({ isCorrect: data.isCorrect, answer: { matches: textMatches } });
              }}
            />
          {:else if currentQuestion.tipo === 'sequencing'}
            {@const transformedQuestion = transformarPreguntaSequencing(currentQuestion)}
            <Sequencing
              title={transformedQuestion.title}
              items={transformedQuestion.items}
              bloomLevel="comprender"
              materia={materiaInfo?.nombre || ''}
              shuffleItems={transformedQuestion.shuffleItems}
              showNumbers={transformedQuestion.showNumbers}
              showHints={transformedQuestion.showHints}
              allowMultipleAttempts={transformedQuestion.allowMultipleAttempts}
              onAnswer={(data) => {
                // Transform userOrder from IDs to content texts
                const items = transformedQuestion.items;
                const sequenceTexts = data.userOrder.map((id: number) => {
                  const item = items.find((i: any) => i.id === id);
                  return item ? item.content : '';
                });
                handleAnswer({ isCorrect: data.isCorrect, answer: { sequence: sequenceTexts } });
              }}
            />
          {:else if currentQuestion.tipo === 'compare_contrast'}
            {@const transformedQuestion = transformarPreguntaCompareContrast(currentQuestion)}
            <CompareContrast
              title={transformedQuestion.title}
              itemA={transformedQuestion.itemA}
              itemB={transformedQuestion.itemB}
              characteristics={transformedQuestion.characteristics}
              bloomLevel="analizar"
              materia={materiaInfo?.nombre || ''}
              showFeedback={transformedQuestion.showFeedback}
              allowMultipleAttempts={transformedQuestion.allowMultipleAttempts}
              onAnswer={(data) => handleAnswer({ isCorrect: data.isCorrect })}
            />
          {:else if currentQuestion.tipo === 'criteria_evaluation'}
            {@const transformedQuestion = transformarPreguntaCriteriaEvaluation(currentQuestion)}
            <CriteriaEvaluation
              title={transformedQuestion.title}
              subject={transformedQuestion.subject}
              description={transformedQuestion.description}
              content={transformedQuestion.content}
              criteria={transformedQuestion.criteria}
              bloomLevel="evaluar"
              materia={materiaInfo?.nombre || ''}
              showFeedback={transformedQuestion.showFeedback}
              allowMultipleAttempts={transformedQuestion.allowMultipleAttempts}
              showExpectedRatings={transformedQuestion.showExpectedRatings}
              onAnswer={(data) => handleAnswer({ isCorrect: data.isCorrect })}
            />
          {:else if currentQuestion.tipo === 'open_ended'}
            {@const transformedQuestion = transformarPreguntaOpenEnded(currentQuestion)}
            <OpenEndedResponse
              prompt={transformedQuestion.prompt}
              placeholder={transformedQuestion.placeholder}
              minWords={transformedQuestion.minWords}
              maxWords={transformedQuestion.maxWords}
              bloomLevel="analizar"
              materia={materiaInfo?.nombre || ''}
              showWordCount={transformedQuestion.showWordCount}
              enableAiFeedback={transformedQuestion.enableAiFeedback}
              rubric={transformedQuestion.rubric}
              onSubmit={(data) => handleAnswer({ isCorrect: true, answer: data.response })}
            />
          {:else if currentQuestion.tipo === 'concept_map'}
            {@const transformedQuestion = transformarPreguntaConceptMap(currentQuestion)}
            <ConceptMapBuilder
              title={transformedQuestion.title}
              topic={transformedQuestion.topic}
              instructions={transformedQuestion.instructions}
              requiredConcepts={transformedQuestion.requiredConcepts}
              suggestedConnections={transformedQuestion.suggestedConnections}
              minConcepts={transformedQuestion.minConcepts}
              minConnections={transformedQuestion.minConnections}
              bloomLevel="crear"
              materia={materiaInfo?.nombre || ''}
              showFeedback={transformedQuestion.showFeedback}
              allowMultipleAttempts={transformedQuestion.allowMultipleAttempts}
              provideConcepts={transformedQuestion.provideConcepts}
              onAnswer={(data) => handleAnswer({ isCorrect: data.isCorrect })}
            />
          {:else}
            <div class="text-center p-8 bg-red-500/10 border border-red-500/20 rounded-xl">
              <p class="text-red-400 mb-2">Tipo de pregunta no soportado: <strong>{currentQuestion.tipo}</strong></p>
              <p class="text-sm text-slate-400">Por favor contacta al administrador</p>
            </div>
          {/if}
        {/key}
      </div>
    {/if}

    <!-- Results Display -->
    {#if showResults}
      <div class="space-y-6">
        <!-- Overall Result -->
        <div class="bg-canvas-900/40 backdrop-blur-xl rounded-2xl border border-white/10 p-8 text-center">
          <div class="text-6xl mb-4">
            {averageBloomLevel >= 5 ? '🏆' : averageBloomLevel >= 3 ? '🎯' : '📚'}
          </div>
          <h2 class="text-3xl font-bold text-white mb-2">¡Evaluación Completada!</h2>
          <p class="text-xl text-slate-300 mb-4">
            Nivel Promedio: <span class="{getBloomLevelDisplay(averageBloomLevel).color} font-bold">
              {getBloomLevelDisplay(averageBloomLevel).label}
            </span>
          </p>
          <p class="text-slate-400 mb-6">
            Respondiste {totalQuestions} preguntas de diferentes objetivos de aprendizaje
          </p>
        </div>

        <!-- Results by OA -->
        <div class="space-y-4">
          <h3 class="text-xl font-bold text-white">Resultados por Objetivo de Aprendizaje</h3>

          {#each results as result}
            <div class="bg-canvas-900/40 backdrop-blur-xl rounded-xl border border-white/10 p-6">
              <div class="flex items-center justify-between mb-3">
                <div>
                  <h4 class="font-semibold text-white">{result.oa?.titulo || `OA ${result.oa_id}`}</h4>
                  <p class="text-xs text-slate-400">{result.oa?.codigo}</p>
                </div>
                <div class="text-right">
                  <div class="text-2xl font-bold {getBloomLevelDisplay(result.nivel_bloom_dominado).color}">
                    {getBloomLevelDisplay(result.nivel_bloom_dominado).label}
                  </div>
                  <div class="text-sm text-slate-400">
                    {result.preguntas_correctas}/{result.preguntas_respondidas} correctas
                  </div>
                </div>
              </div>

              <!-- Progress Bar -->
              <div class="h-2 w-full bg-canvas-950 rounded-full overflow-hidden mb-3">
                <div
                  class="h-full bg-gradient-to-r from-lumera-500 to-focus-500 transition-all duration-500"
                  style="width: {result.porcentaje_aciertos}%"
                ></div>
              </div>

              <!-- Recommendation -->
              <p class="text-sm text-slate-300">
                {result.recomendacion}
              </p>
            </div>
          {/each}
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-4 justify-center pt-4">
          <button
            onclick={() => goto('/')}
            class="px-8 py-3 rounded-xl bg-canvas-800 hover:bg-canvas-700 text-white font-semibold transition-all"
          >
            Volver al Dashboard
          </button>
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
  onAvatarChanged={handleAvatarChanged}
/>
