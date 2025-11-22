<script>
  import { onMount } from 'svelte';
  // Componentes Generales
  import ConceptIntroSlide from './general/ConceptIntroSlide.svelte';
  import ComparisonTableSlide from './general/ComparisonTableSlide.svelte';
  import StepByStepProcessSlide from './general/StepByStepProcessSlide.svelte';
  import FormulaExplorerSlide from './general/FormulaExplorerSlide.svelte';
  import PracticePromptSlide from './general/PracticePromptSlide.svelte';
  // Componentes de Enseñanza (TEACH)
  import ReadingStrategySlide from './teach/ReadingStrategySlide.svelte';
  import GrammarConceptSlide from './teach/GrammarConceptSlide.svelte';
  import ConnectorsGuideSlide from './teach/ConnectorsGuideSlide.svelte';
  import VocabularyStrategySlide from './teach/VocabularyStrategySlide.svelte';
  import TextTypesGuideSlide from './teach/TextTypesGuideSlide.svelte';
  import LiteraryDeviceGuideSlide from './teach/LiteraryDeviceGuideSlide.svelte';
  // Componentes de Práctica (PRACTICE)
  import TextAnnotationSlide from './practice/TextAnnotationSlide.svelte';
  import SentenceBuilderSlide from './practice/SentenceBuilderSlide.svelte';
  import VocabularyContextSlide from './practice/VocabularyContextSlide.svelte';
  import TextStructureSlide from './practice/TextStructureSlide.svelte';
  import ConnectorsWorkshopSlide from './practice/ConnectorsWorkshopSlide.svelte';
  import LiteraryDevicesExplorerSlide from './practice/LiteraryDevicesExplorerSlide.svelte';

  // Props
  let {
    leccion = {
      leccionId: "default",
      titulo: "Lección",
      materia: "general",
      slides: []
    },
    onComplete = null,
    onSlideChange = null,
    showProgress = true
  } = $props();

  // Estados locales
  let currentSlideIndex = $state(0);
  let slideStartTime = $state(Date.now());
  let slideInteractions = $state([]);

  const currentSlide = $derived(leccion.slides[currentSlideIndex] || null);
  const isFirstSlide = $derived(currentSlideIndex === 0);
  const isLastSlide = $derived(currentSlideIndex === leccion.slides.length - 1);
  const progressPercentage = $derived(((currentSlideIndex + 1) / leccion.slides.length) * 100);

  // Mapa de componentes
  const componentMap = {
    // Generales
    'ConceptIntroSlide': ConceptIntroSlide,
    'ComparisonTableSlide': ComparisonTableSlide,
    'StepByStepProcessSlide': StepByStepProcessSlide,
    'FormulaExplorerSlide': FormulaExplorerSlide,
    'PracticePromptSlide': PracticePromptSlide,
    // Enseñanza (TEACH)
    'ReadingStrategySlide': ReadingStrategySlide,
    'GrammarConceptSlide': GrammarConceptSlide,
    'ConnectorsGuideSlide': ConnectorsGuideSlide,
    'VocabularyStrategySlide': VocabularyStrategySlide,
    'TextTypesGuideSlide': TextTypesGuideSlide,
    'LiteraryDeviceGuideSlide': LiteraryDeviceGuideSlide,
    // Práctica (PRACTICE)
    'TextAnnotationSlide': TextAnnotationSlide,
    'SentenceBuilderSlide': SentenceBuilderSlide,
    'VocabularyContextSlide': VocabularyContextSlide,
    'TextStructureSlide': TextStructureSlide,
    'ConnectorsWorkshopSlide': ConnectorsWorkshopSlide,
    'LiteraryDevicesExplorerSlide': LiteraryDevicesExplorerSlide
  };

  function handleNext() {
    trackSlideCompletion();

    if (isLastSlide) {
      // Lección completada
      if (onComplete) {
        onComplete({
          leccionId: leccion.leccionId,
          tiempoTotal: calculateTotalTime(),
          slidesCompletados: leccion.slides.length,
          timestamp: new Date().toISOString()
        });
      }
    } else {
      currentSlideIndex++;
      slideStartTime = Date.now();
      slideInteractions = [];

      if (onSlideChange) {
        onSlideChange({
          slideIndex: currentSlideIndex,
          slideType: currentSlide.tipo
        });
      }
    }
  }

  function handlePrevious() {
    trackSlideCompletion();

    if (!isFirstSlide) {
      currentSlideIndex--;
      slideStartTime = Date.now();
      slideInteractions = [];

      if (onSlideChange) {
        onSlideChange({
          slideIndex: currentSlideIndex,
          slideType: currentSlide.tipo
        });
      }
    }
  }

  function trackSlideCompletion() {
    const tiempoEnSlide = (Date.now() - slideStartTime) / 1000;

    // Aquí se enviaría a backend
    console.log('[LessonPlayer] Slide completado:', {
      leccionId: leccion.leccionId,
      slideIndex: currentSlideIndex,
      slideType: currentSlide?.tipo,
      tiempoSegundos: tiempoEnSlide,
      interacciones: slideInteractions.length,
      timestamp: new Date().toISOString()
    });
  }

  function calculateTotalTime() {
    // Calcular tiempo total de la lección
    // En producción, esto vendría del tracking acumulado
    return 0;
  }

  // Inicialización
  onMount(() => {
    slideStartTime = Date.now();
  });
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 p-4 md:p-8">
  <!-- Header de la lección -->
  <div class="max-w-7xl mx-auto mb-8">
    <div class="bg-slate-950 rounded-2xl border border-slate-800 p-6">
      <div class="flex items-center justify-between flex-wrap gap-4">
        <!-- Título -->
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-white mb-1">
            {leccion.titulo}
          </h1>
          <p class="text-sm text-slate-400">
            {leccion.materia} • {leccion.slides.length} slides
          </p>
        </div>

        <!-- Progress -->
        {#if showProgress}
          <div class="flex items-center gap-4">
            <div class="text-right">
              <p class="text-sm font-semibold text-slate-300">
                Slide {currentSlideIndex + 1} de {leccion.slides.length}
              </p>
              <p class="text-xs text-slate-500">
                {progressPercentage.toFixed(0)}% completado
              </p>
            </div>
            <div class="w-32 h-32 relative">
              <!-- Círculo de progreso -->
              <svg class="transform -rotate-90" width="128" height="128">
                <!-- Fondo -->
                <circle
                  cx="64"
                  cy="64"
                  r="56"
                  stroke="rgb(30 41 59)"
                  stroke-width="8"
                  fill="none"
                />
                <!-- Progreso -->
                <circle
                  cx="64"
                  cy="64"
                  r="56"
                  stroke="url(#gradient)"
                  stroke-width="8"
                  fill="none"
                  stroke-dasharray="351.858"
                  stroke-dashoffset={351.858 * (1 - progressPercentage / 100)}
                  class="transition-all duration-500"
                  stroke-linecap="round"
                />
                <defs>
                  <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#3b82f6;stop-opacity:1" />
                    <stop offset="100%" style="stop-color:#8b5cf6;stop-opacity:1" />
                  </linearGradient>
                </defs>
              </svg>
              <div class="absolute inset-0 flex items-center justify-center">
                <span class="text-2xl font-bold text-white">
                  {progressPercentage.toFixed(0)}%
                </span>
              </div>
            </div>
          </div>
        {/if}
      </div>

      <!-- Barra de progreso lineal (mobile) -->
      {#if showProgress}
        <div class="mt-4 md:hidden">
          <div class="w-full h-2 bg-slate-800 rounded-full overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-blue-500 to-purple-500 transition-all duration-500"
              style="width: {progressPercentage}%"
            ></div>
          </div>
        </div>
      {/if}
    </div>
  </div>

  <!-- Slide actual -->
  <div class="max-w-7xl mx-auto">
    {#if currentSlide}
      {#key currentSlide.orden}
        {#if componentMap[currentSlide.tipo]}
          <svelte:component
            this={componentMap[currentSlide.tipo]}
            {...currentSlide.props}
            onNext={handleNext}
            onPrevious={isFirstSlide ? null : handlePrevious}
          />
        {:else}
          <div class="text-center p-12 bg-slate-950 rounded-2xl border border-slate-800">
            <p class="text-red-400 font-semibold mb-2">
              ⚠️ Tipo de slide desconocido: {currentSlide.tipo}
            </p>
            <p class="text-slate-500 text-sm">
              Verifica que el tipo de slide esté correctamente configurado
            </p>
          </div>
        {/if}
      {/key}
    {:else}
      <div class="text-center p-12 bg-slate-950 rounded-2xl border border-slate-800">
        <p class="text-slate-400">
          No hay slides para mostrar en esta lección
        </p>
      </div>
    {/if}
  </div>

  <!-- Footer con navegación de teclado -->
  <div class="max-w-7xl mx-auto mt-8">
    <div class="bg-slate-950 rounded-2xl border border-slate-800 p-4">
      <div class="flex items-center justify-center gap-8 text-xs text-slate-500">
        <div class="flex items-center gap-2">
          <kbd class="px-2 py-1 bg-slate-800 rounded border border-slate-700 font-mono">←</kbd>
          <span>Anterior</span>
        </div>
        <div class="flex items-center gap-2">
          <kbd class="px-2 py-1 bg-slate-800 rounded border border-slate-700 font-mono">→</kbd>
          <span>Siguiente</span>
        </div>
        <div class="flex items-center gap-2">
          <kbd class="px-2 py-1 bg-slate-800 rounded border border-slate-700 font-mono">Esc</kbd>
          <span>Salir</span>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
