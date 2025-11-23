<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth.svelte';
  import {
    startPracticeSession,
    getNextPracticeQuestion,
    submitPracticeAnswer,
    completePracticeSession
  } from '$lib/api/practice';
  import { getOABloomObjectiveId } from '$lib/api/learningPlans';

  // Import question components
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';
  import TrueFalse from '$lib/components/activities/TrueFalse.svelte';
  import FillBlanks from '$lib/components/activities/FillBlanks.svelte';
  import DragDropMatching from '$lib/components/activities/DragDropMatching.svelte';
  import Sequencing from '$lib/components/activities/Sequencing.svelte';
  import ComparisonTable from '$lib/components/activities/ComparisonTable.svelte';
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';
  import { getEquipment } from '$lib/api/customization';

  // Route params
  let materiaId = $state(0);
  let oaId = $state(0);
  let oaInfo = $state(null);
  let materiaInfo = $state(null);

  // Practice Session State
  let session = $state(null);
  let currentQuestion = $state(null);
  let currentQuestionNumber = $state(0);
  let totalQuestions = $state(10);

  // UI State
  let hasAnswered = $state(false);
  let isLoading = $state(false);
  let isCorrect = $state(false);
  let showResults = $state(false);
  let results = $state(null);
  let errorMessage = $state('');
  let currentAvatar = $state(null);
  let slideStartTime = $state(Date.now());

  // Student data
  const student = {
    name: auth.user?.name || 'Student',
    grade: '2° Medio'
  };

  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  const progreso = $derived(currentQuestionNumber > 0 && totalQuestions > 0
    ? Math.round((currentQuestionNumber / totalQuestions) * 100)
    : 0);

  // Fetch OA info
  async function fetchOAById(id) {
    try {
      const response = await fetch(`/api/objetivos-aprendizaje/${id}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });
      if (!response.ok) throw new Error('Error fetching OA');
      return await response.json();
    } catch (error) {
      console.error('Error fetching OA:', error);
      return null;
    }
  }

  // Fetch materia info
  async function fetchMateriaById(id) {
    try {
      const response = await fetch(`/api/materias/${id}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });
      if (!response.ok) throw new Error('Error fetching materia');
      return await response.json();
    } catch (error) {
      console.error('Error fetching materia:', error);
      return null;
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

  // Get user's current Bloom level for this OA
  async function getUserBloomLevel() {
    try {
      const response = await fetch('/api/diagnostic-sessions?estado=completado', {
        headers: { 'Authorization': auth.token ? `Bearer ${auth.token}` : '' }
      });
      if (!response.ok) return 3;
      const sessions = await response.json();
      const session = sessions.find(s => s.materia_id === materiaId);
      if (session?.estrategia?.average_bloom_level) {
        return Math.round(session.estrategia.average_bloom_level);
      }
      return 3;
    } catch (error) {
      console.error('Error loading bloom level:', error);
      return 3;
    }
  }

  // Initialize practice session
  async function initializePractice() {
    if (!oaInfo) return;

    isLoading = true;
    errorMessage = '';

    const bloomLevel = await getUserBloomLevel();
    const objectiveId = await getOABloomObjectiveId(oaId, bloomLevel);

    if (!objectiveId) {
      errorMessage = 'No se pudo encontrar el objetivo de Bloom para este OA';
      isLoading = false;
      return;
    }

    const sessionResult = await startPracticeSession(oaId, objectiveId, 10);

    if (!sessionResult.success) {
      errorMessage = sessionResult.error || 'Error al iniciar la práctica';
      isLoading = false;
      return;
    }

    session = sessionResult.session;
    totalQuestions = session.numero_preguntas;
    await loadNextQuestion();
    isLoading = false;
  }

  // Load next question
  async function loadNextQuestion() {
    if (!session) return;

    isLoading = true;
    errorMessage = '';

    const questionResult = await getNextPracticeQuestion(session.id);

    if (!questionResult.success) {
      errorMessage = questionResult.error || 'No hay más preguntas disponibles';
      isLoading = false;
      if (currentQuestionNumber > 0) {
        await finalizePractice();
      }
      return;
    }

    currentQuestion = questionResult.question;
    currentQuestionNumber = currentQuestion.question_number;
    totalQuestions = currentQuestion.total_questions;
    hasAnswered = false;
    isCorrect = false;
    slideStartTime = Date.now();
    isLoading = false;
  }

  // Handle answer submission
  async function handleAnswer(data) {
    if (hasAnswered || !session || !currentQuestion) return;

    hasAnswered = true;
    isCorrect = data.isCorrect;
    isLoading = true;

    const tiempoSegundos = Math.floor((Date.now() - slideStartTime) / 1000);

    // Determine user answer based on question type
    let userAnswer;
    if (data.selectedOption !== undefined) {
      const letters = ['A', 'B', 'C', 'D', 'E', 'F'];
      userAnswer = { selected: letters[data.selectedOption] };
    } else if (data.userAnswers !== undefined) {
      const blanksMap = {};
      for (const [key, value] of Object.entries(data.userAnswers)) {
        blanksMap[String(key)] = String(value);
      }
      userAnswer = { blanks: blanksMap };
    } else if (data.userAnswer !== undefined) {
      userAnswer = { answer: data.userAnswer };
    } else if (data.userTable !== undefined) {
      userAnswer = { tabla: data.userTable };
    } else if (data.answer !== undefined) {
      userAnswer = data.answer;
    } else {
      userAnswer = { is_correct: data.isCorrect };
    }

    const answerResult = await submitPracticeAnswer(
      session.id,
      currentQuestion.id,
      userAnswer,
      tiempoSegundos
    );

    isLoading = false;

    if (!answerResult.success) {
      errorMessage = answerResult.error || 'Error al enviar la respuesta';
      return;
    }

    setTimeout(() => {
      siguientePregunta();
    }, 1500);
  }

  // Move to next question
  async function siguientePregunta() {
    if (currentQuestionNumber >= totalQuestions) {
      await finalizePractice();
    } else {
      await loadNextQuestion();
    }
  }

  // Complete practice session
  async function finalizePractice() {
    if (!session) return;

    isLoading = true;
    errorMessage = '';

    const completeResult = await completePracticeSession(session.id);

    if (!completeResult.success) {
      errorMessage = completeResult.error || 'Error al completar la práctica';
      isLoading = false;
      return;
    }

    results = completeResult.result;
    showResults = true;
    isLoading = false;
  }

  // Map Bloom level to display
  function getBloomLevelDisplay(bloomLevel) {
    const bloomMap = {
      0: { label: 'Sin Evaluar', color: 'text-slate-400' },
      1: { label: 'Recordar', color: 'text-red-400' },
      2: { label: 'Comprender', color: 'text-orange-400' },
      3: { label: 'Aplicar', color: 'text-yellow-400' },
      4: { label: 'Analizar', color: 'text-green-400' },
      5: { label: 'Evaluar', color: 'text-blue-400' },
      6: { label: 'Crear', color: 'text-purple-400' }
    };
    return bloomMap[bloomLevel] || bloomMap[0];
  }

  function getBloomLevelStars(bloomLevel) {
    const level = Math.floor(bloomLevel);
    if (level === 0) return '—';
    return '⭐'.repeat(Math.min(level, 6));
  }

  // Initialize on mount
  onMount(async () => {
    loadEquippedAvatar();

    const unsubscribe = page.subscribe(($page) => {
      materiaId = parseInt($page.params.materiaId || '0');
      oaId = parseInt($page.params.oaId || '0');
    });

    oaInfo = await fetchOAById(oaId);
    materiaInfo = await fetchMateriaById(materiaId);

    if (!oaInfo || !materiaInfo) {
      errorMessage = 'OA o Materia no encontrados';
      return () => { if (unsubscribe) unsubscribe(); };
    }

    await initializePractice();

    return () => { if (unsubscribe) unsubscribe(); };
  });
</script>

<svelte:head>
  <title>Práctica - {oaInfo?.titulo || 'Lumera App'}</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950">
  <!-- Header -->
  <header class="sticky top-0 z-50 border-b border-white/5 bg-canvas-950/90 backdrop-blur-md">
    <div class="px-6 py-3 flex items-center justify-between gap-4">
      <!-- Left: Profile -->
      <div class="flex items-center gap-3 px-3 py-2">
        <AvatarDisplay
          {currentAvatar}
          {initials}
          size="small"
        />
        <div class="text-left hidden md:block">
          <div class="text-sm font-semibold text-white">{student.name}</div>
          <div class="text-xs text-slate-400">{student.grade}</div>
        </div>
      </div>

      <!-- Center: Progress -->
      <div class="flex-1 max-w-md">
        <div class="flex items-center justify-between text-xs mb-1.5">
          {#if showResults}
            <span class="text-green-400 font-medium">✓ Práctica Completada</span>
            <span class="text-slate-400">{getBloomLevelStars(results?.bloom_level_final || 0)}</span>
          {:else if currentQuestion}
            <span class="text-focus-400 font-medium truncate">🎯 {oaInfo?.titulo || 'Práctica'}</span>
            <span class="text-slate-400 flex-shrink-0">{currentQuestionNumber}/{totalQuestions}</span>
          {:else}
            <span class="text-slate-400">Cargando...</span>
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

      <!-- Right: Exit button -->
      <button
        onclick={() => goto(`/materias/${materiaId}/objetivos`)}
        class="px-4 py-2 rounded-lg bg-canvas-800 hover:bg-canvas-700 text-white transition-all"
      >
        Salir
      </button>
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
          <p class="text-slate-400">Cargando práctica...</p>
        </div>
      </div>
    {/if}

    <!-- Question Display -->
    {#if !showResults && currentQuestion && !isLoading}
      <div class="bg-canvas-900/40 backdrop-blur-xl rounded-2xl border border-white/10 p-8">
        <div class="text-center mb-6">
          <p class="text-slate-400 text-sm">
            Nivel actual:
            <span class="{getBloomLevelDisplay(currentQuestion.current_bloom_level).color} font-semibold">
              {getBloomLevelDisplay(currentQuestion.current_bloom_level).label}
            </span>
          </p>
        </div>

        {#key currentQuestion.id}
          {#if currentQuestion.tipo === 'multiple_choice'}
            <MultipleChoice
              question={currentQuestion.question_data.pregunta || currentQuestion.question_data.question || ''}
              options={Object.values(currentQuestion.question_data.opciones || []).map((text, index) => ({
                id: index,
                text: String(text),
                isCorrect: false
              }))}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              showCorrectAnswer={false}
              showFeedback={false}
              allowMultipleAttempts={false}
              onAnswer={(data) => handleAnswer({ isCorrect: true, selectedOption: data.userAnswer })}
            />
          {:else if currentQuestion.tipo === 'true_false'}
            <TrueFalse
              statement={currentQuestion.question_data.afirmacion || currentQuestion.question_data.pregunta || ''}
              correctAnswer={currentQuestion.question_data.respuesta_correcta || false}
              explanation={currentQuestion.question_data.explicacion || ''}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showExplanation={false}
              onAnswer={(data) => handleAnswer(data)}
            />
          {:else if currentQuestion.tipo === 'fill_blanks'}
            {@const texto = (currentQuestion.question_data.texto || '').replace(/\[BLANK_(\d+)\]/g, '___$1___')}
            {@const blankMatches = [...texto.matchAll(/___(\d+)___/g)]}
            {@const blanks = blankMatches.map(match => ({ id: parseInt(match[1]), answer: '', caseSensitive: false }))}
            <FillBlanks
              text={texto}
              {blanks}
              bloomLevel="recordar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showWordBank={false}
              showHints={false}
              onAnswer={(data) => handleAnswer(data)}
            />
          {:else if currentQuestion.tipo === 'drag_drop_matching'}
            {@const pares = currentQuestion.question_data.columna_izquierda
              ? currentQuestion.question_data.columna_izquierda.map((term, index) => ({
                  id: index + 1,
                  term: term,
                  definition: currentQuestion.question_data.columna_derecha[index] || ''
                }))
              : currentQuestion.question_data.pares || []}
            <DragDropMatching
              title={currentQuestion.question_data.instruccion || 'Relaciona los conceptos'}
              pairs={pares}
              bloomLevel="comprender"
              materia={materiaInfo?.nombre || ''}
              shuffleOptions={true}
              showFeedback={false}
              allowMultipleAttempts={false}
              onAnswer={(data) => {
                const textMatches = {};
                for (const [termId, defId] of Object.entries(data.matches)) {
                  const pair = pares.find(p => p.id === parseInt(termId));
                  const defPair = pares.find(p => p.id === parseInt(defId));
                  if (pair && defPair) {
                    textMatches[pair.term] = defPair.definition;
                  }
                }
                handleAnswer({ isCorrect: data.isCorrect, answer: { matches: textMatches } });
              }}
            />
          {:else if currentQuestion.tipo === 'sequencing'}
            {@const items = currentQuestion.question_data.elementos_desordenados
              ? currentQuestion.question_data.elementos_desordenados.map((content, index) => ({
                  id: index + 1,
                  content: content,
                  correctOrder: index + 1
                }))
              : currentQuestion.question_data.items || []}
            <Sequencing
              title={currentQuestion.question_data.instruccion || 'Ordena la secuencia'}
              {items}
              bloomLevel="comprender"
              materia={materiaInfo?.nombre || ''}
              shuffleItems={true}
              showNumbers={true}
              showHints={false}
              allowMultipleAttempts={false}
              onAnswer={(data) => {
                const sequenceTexts = data.userOrder.map(id => {
                  const item = items.find(i => i.id === id);
                  return item ? item.content : '';
                });
                handleAnswer({ isCorrect: data.isCorrect, answer: { sequence: sequenceTexts } });
              }}
            />
          {:else if currentQuestion.tipo === 'compare_contrast'}
            <ComparisonTable
              title={currentQuestion.question_data.titulo || 'Compara y contrasta'}
              instruction={currentQuestion.question_data.instruccion || 'Completa la tabla comparativa'}
              concepts={currentQuestion.question_data.conceptos || []}
              criteria={currentQuestion.question_data.criterios || []}
              bloomLevel="analizar"
              materia={materiaInfo?.nombre || ''}
              allowMultipleAttempts={false}
              showFeedback={false}
              onAnswer={(data) => handleAnswer(data)}
            />
          {:else}
            <div class="text-center p-8 bg-slate-800/50 rounded-xl">
              <p class="text-slate-400 mb-2">Tipo de pregunta: <strong>{currentQuestion.tipo}</strong></p>
              <p class="text-sm text-slate-500">Este tipo de pregunta se saltará automáticamente</p>
              <button
                onclick={() => siguientePregunta()}
                class="mt-4 px-6 py-3 rounded-xl font-semibold bg-gradient-to-r from-focus-500 to-blue-500 text-white transition-all duration-300 hover:shadow-lg"
              >
                Continuar
              </button>
            </div>
          {/if}
        {/key}
      </div>
    {/if}

    <!-- Results Display -->
    {#if showResults && results}
      <div class="space-y-6">
        <div class="bg-canvas-900/40 backdrop-blur-xl rounded-2xl border border-white/10 p-8 text-center">
          <div class="text-6xl mb-4">
            {#if results.cambio_nivel > 0}
              🎉
            {:else if results.cambio_nivel < 0}
              📚
            {:else}
              🎯
            {/if}
          </div>
          <h2 class="text-3xl font-bold text-white mb-2">¡Práctica Completada!</h2>

          <div class="flex items-center justify-center gap-8 mt-6">
            <div>
              <p class="text-sm text-slate-400">Nivel Inicial</p>
              <p class="text-2xl">{getBloomLevelStars(results.bloom_level_inicial)}</p>
              <p class="text-xs text-slate-500">{getBloomLevelDisplay(results.bloom_level_inicial).label}</p>
            </div>

            <div class="text-4xl {results.cambio_nivel > 0 ? 'text-green-400' : results.cambio_nivel < 0 ? 'text-red-400' : 'text-slate-400'}">
              {results.cambio_nivel > 0 ? '↗' : results.cambio_nivel < 0 ? '↘' : '→'}
            </div>

            <div>
              <p class="text-sm text-slate-400">Nivel Final</p>
              <p class="text-2xl">{getBloomLevelStars(results.bloom_level_final)}</p>
              <p class="text-xs text-slate-500">{getBloomLevelDisplay(results.bloom_level_final).label}</p>
            </div>
          </div>

          <div class="mt-6 pt-6 border-t border-white/10">
            <p class="text-lg text-slate-300">
              {results.resultado.preguntas_correctas} de {results.resultado.preguntas_totales} correctas
            </p>
            <p class="text-sm text-slate-400">
              {results.resultado.porcentaje_aciertos.toFixed(1)}% de aciertos
            </p>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-4 justify-center">
          <button
            onclick={() => goto(`/materias/${materiaId}/objetivos`)}
            class="px-8 py-3 rounded-xl bg-canvas-800 hover:bg-canvas-700 text-white font-semibold transition-all"
          >
            Volver a Objetivos
          </button>
          <button
            onclick={() => window.location.reload()}
            class="px-8 py-3 rounded-xl bg-gradient-to-r from-lumera-600 to-focus-600 hover:from-lumera-500 hover:to-focus-500 text-white font-semibold transition-all"
          >
            Practicar de Nuevo
          </button>
        </div>
      </div>
    {/if}
  </main>
</div>
