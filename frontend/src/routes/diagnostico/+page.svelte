<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
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

  // State
  let currentQuestionIndex = $state(0);
  let answers = $state<Record<string, boolean>>({}); // preguntaId -> esCorrecta
  let hasAnswered = $state(false);
  let isLoading = $state(false);
  let showResults = $state(false);
  let diagnosticoCompletado = $state(false);

  const preguntas = PREGUNTAS_DIAGNOSTICO_LENGUA;
  const totalPreguntas = preguntas.length;

  // Derived
  const preguntaActual = $derived(preguntas[currentQuestionIndex]);
  const progreso = $derived(Math.round(((currentQuestionIndex + 1) / totalPreguntas) * 100));
  const respuestasCorrectas = $derived(Object.values(answers).filter(Boolean).length);
  const nivelCalculado = $derived(calcularNivelDiagnostico(respuestasCorrectas, totalPreguntas));

  // Verificar autenticación
  onMount(() => {
    if (!auth.checkAuth()) {
      goto('/login');
      return;
    }

    // Animación de entrada
    gsap.from('.diagnostico-container', {
      opacity: 0,
      y: 30,
      duration: 0.6,
      ease: 'power2.out'
    });
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
    if (!auth.user?.id || diagnosticoCompletado) return;

    isLoading = true;

    try {
      const profileDataUpdate = {
        conocimiento_previo: {
          lectura: {
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
  <title>Evaluación Diagnóstica - Lengua y Literatura</title>
</svelte:head>

<div class="min-h-screen bg-canvas-950 text-white p-6">
  <div class="diagnostico-container max-w-4xl mx-auto">
    {#if !showResults}
      <!-- Header con progreso -->
      <div class="mb-8">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h1 class="text-3xl font-bold text-white mb-1">
              📚 Evaluación Diagnóstica
            </h1>
            <p class="text-slate-400">Lengua y Literatura - 1° Medio</p>
          </div>
          <div class="text-right">
            <div class="text-sm text-slate-400 mb-1">
              Pregunta {currentQuestionIndex + 1} de {totalPreguntas}
            </div>
            <div class="text-2xl font-bold text-focus-400">{progreso}%</div>
          </div>
        </div>

        <!-- Barra de progreso -->
        <div class="h-3 w-full bg-canvas-900 rounded-full overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-focus-500 to-blue-500 transition-all duration-500"
            style="width: {progreso}%"
          ></div>
        </div>
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
</div>
