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
  import ComparisonTable from '$lib/components/activities/ComparisonTable.svelte';
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
  import AppHeader from '$lib/components/common/AppHeader.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import { getEquipment, type CustomizationItem } from '$lib/api/customization';

  // Obtener parámetro de URL
  let materiaId = $state(0);
  let materiaInfo = $state<{ id: number; nombre: string; codigo: string } | null>(null);

  // Fetch materia by ID from backend
  async function fetchMateriaById(id: number) {
    try {
      const response = await fetch(`/api/materias/${id}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });

      if (!response.ok) {
        throw new Error('Error fetching materia');
      }

      return await response.json();
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
  // Keep these for modals even though they're not in header anymore
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
  let showFeedback = $state(false);
  let feedbackMessage = $state('');

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
  // Helper to get auth headers
  function getAuthHeaders() {
    const token = auth.token;
    return {
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    };
  }

  // Check if user already has a completed diagnostic for this materia
  async function checkExistingDiagnostic() {
    if (!materiaInfo) return null;

    try {
      const response = await fetch(`/api/diagnostic-sessions?materia_id=${materiaInfo.id}&estado=completado`, {
        headers: getAuthHeaders()
      });

      if (!response.ok) return null;

      const sessions = await response.json();
      // Return the most recent completed session
      return sessions.length > 0 ? sessions[0] : null;
    } catch (error) {
      console.error('Error checking existing diagnostic:', error);
      return null;
    }
  }

  // Load existing diagnostic results
  async function loadExistingDiagnostic(existingSession: any) {
    session = existingSession;
    averageBloomLevel = existingSession.estrategia?.average_bloom_level || 0;
    totalQuestions = existingSession.preguntas_totales;

    // Get detailed results
    const resultsResponse = await getResults(existingSession.id);

    if (resultsResponse.success) {
      results = resultsResponse.results!;
    }

    showResults = true;
    isLoading = false;
  }

  async function initializeDiagnostic() {
    if (!materiaInfo) return;

    isLoading = true;
    errorMessage = '';

    // Check if user wants to start a new diagnostic (ignore existing one)
    const urlParams = new URLSearchParams(window.location.search);
    const isNewDiagnostic = urlParams.get('new') === 'true';

    // Check if user already completed this diagnostic (only if not forcing new)
    if (!isNewDiagnostic) {
      const existingSession = await checkExistingDiagnostic();

      if (existingSession) {
        // Load existing results
        await loadExistingDiagnostic(existingSession);
        return;
      }
    }

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

    // Apply the first question immediately (no feedback to wait for)
    applyNextQuestion();

    isLoading = false;
  }

  // Load next question
  let nextQuestionData = $state<DiagnosticQuestion | null>(null);

  async function loadNextQuestion() {
    if (!session) return;

    // DON'T set isLoading here - we want to show feedback overlay
    errorMessage = '';

    console.log('[Diagnostic] Fetching next question for session:', session.id);
    const questionResult = await getNextQuestion(session.id);
    console.log('[Diagnostic] Question result:', questionResult);

    if (!questionResult.success) {
      errorMessage = questionResult.error || 'No hay más preguntas disponibles';
      // If no more questions, complete the session
      if (currentQuestionNumber > 0) {
        await finalizeDiagnostic();
      }
      return;
    }

    // Store next question but don't update currentQuestion yet
    // This will be updated after feedback is shown
    nextQuestionData = questionResult.question!;

    console.log('[Diagnostic] Next question loaded, waiting for feedback to finish');
  }

  function applyNextQuestion() {
    if (!nextQuestionData) return;

    currentQuestion = nextQuestionData;
    currentQuestionNumber = currentQuestion.question_number;
    totalQuestions = currentQuestion.total_questions;
    hasAnswered = false;
    isCorrect = false;
    isLoading = false;
    nextQuestionData = null;

    console.log('[Diagnostic] Applied next question:', {
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

    console.log('[Diagnostic] Data received from component:', JSON.stringify(data, null, 2));

    // Determine user answer based on question type
    let userAnswer: any;
    if (data.selectedOption !== undefined) {
      // MultipleChoice sends index (0,1,2,3), convert to letter (A,B,C,D)
      const letters = ['A', 'B', 'C', 'D', 'E', 'F'];
      userAnswer = { selected: letters[data.selectedOption] };
      console.log('[Diagnostic] MultipleChoice transformed:', { index: data.selectedOption, letter: letters[data.selectedOption] });
    } else if (data.userAnswers !== undefined) {
      // FillBlanks component sends userAnswers, backend expects blanks
      // Convert to string keys for backend
      const blanksMap: Record<string, string> = {};
      for (const [key, value] of Object.entries(data.userAnswers)) {
        blanksMap[String(key)] = String(value);
      }
      userAnswer = { blanks: blanksMap };
      console.log('[Diagnostic] FillBlanks transformed:', JSON.stringify({ userAnswers: data.userAnswers, blanksMap, userAnswer }, null, 2));
    } else if (data.userAnswer !== undefined) {
      // TrueFalse component sends userAnswer (boolean), backend expects answer field
      userAnswer = { answer: data.userAnswer };
    } else if (data.userTable !== undefined) {
      // ComparisonTable component sends userTable
      userAnswer = { tabla: data.userTable };
    } else if (data.answer !== undefined) {
      userAnswer = data.answer; // Send answer object directly (e.g., { matches: {...} })
    } else {
      userAnswer = { is_correct: data.isCorrect };
    }

    console.log('[Diagnostic] Submitting answer:', {
      sessionId: session.id,
      questionId: currentQuestion.id,
      userAnswer,
      userAnswerJSON: JSON.stringify(userAnswer)
    });

    // Submit answer to backend
    const answerResult = await submitAnswer(
      session.id,
      currentQuestion.id,
      userAnswer
    );

    console.log('[Diagnostic] Answer result:', answerResult);

    if (!answerResult.success) {
      isLoading = false;
      errorMessage = answerResult.error || 'Error al enviar la respuesta';
      console.error('[Diagnostic] Error submitting answer:', errorMessage);
      return;
    }

    // Show success feedback BEFORE turning off loading
    hasAnswered = true;
    isCorrect = answerResult.result?.is_correct || false;
    feedbackMessage = isCorrect ? '¡Correcto! ✓' : 'Respuesta enviada';

    console.log('[Diagnostic] Setting showFeedback = true');
    showFeedback = true;

    console.log('[Diagnostic] Setting isLoading = false');
    isLoading = false;

    console.log('[Diagnostic] Feedback state:', { showFeedback, isLoading, feedbackMessage });

    // Load next question in background while showing feedback
    if (currentQuestionNumber >= totalQuestions) {
      console.log('[Diagnostic] Last question, will finalize after feedback');
      // Last question - prepare to complete
      setTimeout(async () => {
        console.log('[Diagnostic] Hiding feedback and finalizing');
        showFeedback = false;
        hasAnswered = false;
        await finalizeDiagnostic();
      }, 1500);
    } else {
      console.log('[Diagnostic] Loading next question in background');
      // Load next question in background
      loadNextQuestion();

      // Show feedback briefly, then apply next question
      setTimeout(() => {
        console.log('[Diagnostic] Hiding feedback and applying next question');
        showFeedback = false;
        applyNextQuestion();
      }, 1500);
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
      pregunta: questionData.afirmacion || questionData.pregunta || questionData.statement || '',
      respuestaCorrecta: questionData.respuesta_correcta || questionData.correct_answer || false,
      explicacion: questionData.explicacion || questionData.explanation || ''
    };
  }

  function transformarPreguntaFillBlanks(question: DiagnosticQuestion) {
    const questionData = question.question_data;

    // Transform [BLANK_X] format to ___X___ format
    let texto = questionData.texto || questionData.text || '';
    texto = texto.replace(/\[BLANK_(\d+)\]/g, '___$1___');

    // Create blanks array from the text (extract blank numbers)
    const blankMatches = [...texto.matchAll(/___(\d+)___/g)];
    const blanks = blankMatches.map(match => ({
      id: parseInt(match[1]),
      answer: '', // Will be validated by backend
      caseSensitive: false
    }));

    return {
      text: texto,
      blanks: blanks,
      showWordBank: false,
      showHints: false,
      allowMultipleAttempts: false
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
      instruction: questionData.instruccion || questionData.instruction || 'Completa la tabla comparativa',
      concepts: questionData.conceptos || questionData.concepts || [],
      criteria: questionData.criterios || questionData.criteria || [],
      allowMultipleAttempts: false,
      showFeedback: false
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

  function getBloomLevelStars(bloomLevel: number) {
    const level = Math.floor(bloomLevel);
    if (level === 0) return '—';
    return '⭐'.repeat(Math.min(level, 6));
  }

  // Initialize on mount
  onMount(async () => {
    loadUserProfile();
    loadEquippedAvatar();

    const unsubscribe = page.subscribe(($page) => {
      materiaId = parseInt($page.params.materiaId || '0');
    });

    // Fetch materia from backend
    materiaInfo = await fetchMateriaById(materiaId);

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
<div class="min-h-screen bg-canvas-950">
  <!-- Header -->
  <AppHeader
    {currentAvatar}
    onProfileClick={() => isPlayerProfileOpen = true}
    onProgressClick={() => isProgressPanelOpen = true}
    showNavButtons={true}
    isHomePage={false}
  >
    {#snippet centerContent()}
      <!-- Progress Bar for Diagnostic in Header -->
      <div class="flex-1 max-w-md">
        <div class="flex items-center justify-between text-xs mb-1.5">
          {#if showResults}
            <span class="text-green-400 font-medium">✓ Evaluación Completada</span>
            <span class="text-slate-400">{getBloomLevelStars(averageBloomLevel)} - {getBloomLevelDisplay(averageBloomLevel).porcentaje}%</span>
          {:else if currentQuestion}
            <span class="text-cyan-400 font-medium">📚 {materiaInfo?.nombre}</span>
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
            <div class="h-full bg-gradient-to-r from-teal-500 to-cyan-500 transition-all duration-500" style="width: {progreso}%"></div>
          {/if}
        </div>
      </div>
    {/snippet}
  </AppHeader>

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
    {#if !showResults && currentQuestion}
      <div class="relative bg-canvas-900 backdrop-blur-xl rounded-2xl border border-slate-700 p-8">
        <!-- Feedback Overlay -->
        {#if showFeedback}
          <div class="absolute inset-0 bg-red-500 rounded-2xl flex items-center justify-center z-[9999]">
            <div class="text-center px-8 py-6 bg-canvas-900 rounded-xl border-2 {isCorrect ? 'border-green-500' : 'border-lumera-500'} shadow-2xl animate-in zoom-in-95 duration-300">
              <div class="text-6xl mb-4">
                {#if isCorrect}
                  <svg class="w-20 h-20 mx-auto text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                {:else}
                  <svg class="w-20 h-20 mx-auto text-lumera-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                {/if}
              </div>
              <p class="text-2xl font-bold {isCorrect ? 'text-green-400' : 'text-white'} mb-2">
                {feedbackMessage}
              </p>
              <p class="text-sm text-slate-400">
                Continuando a la siguiente pregunta...
              </p>
            </div>
          </div>
        {/if}

        {#key currentQuestion.id}
          {#if isLoading}
            <div class="flex items-center justify-center py-20">
              <div class="text-center">
                <svg class="animate-spin h-12 w-12 mx-auto mb-4 text-lumera-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <p class="text-slate-400">Enviando respuesta...</p>
              </div>
            </div>
          {:else if currentQuestion.tipo === 'multiple_choice'}
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
              onAnswer={(data) => handleAnswer(data)}
            />
          {:else if currentQuestion.tipo === 'fill_blanks'}
            {@const transformedQuestion = transformarPreguntaFillBlanks(currentQuestion)}
            <FillBlanks
              text={transformedQuestion.text}
              blanks={transformedQuestion.blanks}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showWordBank={false}
              showHints={false}
              onAnswer={(data) => handleAnswer(data)}
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
            <ComparisonTable
              title={transformedQuestion.title}
              instruction={transformedQuestion.instruction}
              concepts={transformedQuestion.concepts}
              criteria={transformedQuestion.criteria}
              bloomLevel="analizar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showFeedback={false}
              onAnswer={(data) => handleAnswer(data)}
            />
          {:else if currentQuestion.tipo === 'open_ended' || currentQuestion.tipo === 'concept_map'}
            <!-- Tipo de pregunta no soportado en diagnóstico - se salta automáticamente -->
            <div class="bg-yellow-500/10 border border-yellow-500/50 rounded-xl p-6 text-center">
              <p class="text-yellow-400 mb-4">⚠️ Tipo de pregunta no disponible en modo diagnóstico</p>
              <p class="text-slate-400 text-sm mb-4">Este tipo requiere evaluación manual. Se saltará automáticamente.</p>
              <button
                onclick={() => siguientePregunta()}
                class="px-6 py-3 rounded-xl font-semibold bg-[#E1E1E1] hover:bg-[#CCCCCC] text-canvas-900 transition-all duration-300 shadow-lg hover:shadow-xl"
              >
                Continuar
              </button>
            </div>
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
        <div class="bg-canvas-900 backdrop-blur-xl rounded-2xl border border-slate-700 p-8 text-center">
          <div class="text-6xl mb-4">
            {averageBloomLevel >= 5 ? '🏆' : averageBloomLevel >= 3 ? '🎯' : '📚'}
          </div>
          <h2 class="text-3xl font-bold text-white mb-2">¡Evaluación Completada!</h2>
          <p class="text-xl text-slate-300 mb-4">
            Nivel Alcanzado: <span class="text-4xl">
              {getBloomLevelStars(averageBloomLevel)}
            </span>
          </p>
          <p class="text-sm text-slate-400">
            {getBloomLevelDisplay(averageBloomLevel).label}
          </p>
          <p class="text-slate-400 mb-6">
            Respondiste {totalQuestions} preguntas de diferentes objetivos de aprendizaje
          </p>
        </div>

        <!-- Results by OA -->
        <div class="space-y-4">
          <h3 class="text-xl font-bold text-slate-900">Resultados por Objetivo de Aprendizaje</h3>

          {#each results as result}
            <div class="bg-canvas-800 backdrop-blur-xl rounded-xl border border-slate-700 p-6">
              <div class="flex items-center justify-between mb-3">
                <div>
                  <h4 class="font-semibold text-white">{result.oa?.titulo || `OA ${result.oa_id}`}</h4>
                  <p class="text-xs text-slate-400">{result.oa?.codigo}</p>
                </div>
                <div class="text-right">
                  <div class="text-3xl">
                    {getBloomLevelStars(result.nivel_bloom_dominado)}
                  </div>
                  <div class="text-xs text-slate-500 mt-1">
                    {getBloomLevelDisplay(result.nivel_bloom_dominado).label}
                  </div>
                  <div class="text-sm text-slate-400 mt-1">
                    {result.preguntas_correctas}/{result.preguntas_respondidas} correctas
                  </div>
                </div>
              </div>

              <!-- Progress Bar -->
              <div class="h-2 w-full bg-black/40 rounded-full overflow-hidden mb-3 border border-black/50">
                <div
                  class="h-full bg-gradient-to-r from-teal-500 to-cyan-500 transition-all duration-500"
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
            class="px-8 py-3 rounded-xl bg-[#E1E1E1] hover:bg-[#CCCCCC] text-canvas-900 font-semibold transition-all shadow-lg hover:shadow-xl"
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
