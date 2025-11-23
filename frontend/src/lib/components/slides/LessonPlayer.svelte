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
  import ExplainAndExploreSlide from './teach/ExplainAndExploreSlide.svelte';
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
    showProgress = true,
    showHeader = true  // Nueva prop para controlar si se muestra el header
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
    'ExplainAndExploreSlide': ExplainAndExploreSlide,
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

<div class="{showHeader ? 'min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 p-4 md:p-6' : ''}">
  <!-- Header de la lección - Compacto -->
  {#if showHeader}
  <div class="max-w-7xl mx-auto mb-4">
    <div class="bg-slate-950 rounded-xl border border-slate-800 p-4">
      <div class="flex items-center justify-between gap-4">
        <!-- Título y metadata -->
        <div class="flex-1 min-w-0">
          <h1 class="text-xl md:text-2xl font-bold text-white mb-1 truncate">
            {leccion.titulo}
          </h1>
          <div class="flex items-center gap-3 text-xs text-slate-400">
            <span class="flex items-center gap-1">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
              {leccion.materia}
            </span>
            <span class="text-slate-600">•</span>
            <span>{leccion.slides.length} componentes</span>
          </div>
        </div>

        <!-- Progress compacto -->
        {#if showProgress}
          <div class="hidden md:flex items-center gap-3">
            <!-- Stats numéricos -->
            <div class="text-right pr-3 border-r border-slate-700">
              <p class="text-xs text-slate-500">Progreso</p>
              <p class="text-lg font-bold text-white">
                {currentSlideIndex + 1}<span class="text-sm text-slate-500">/{leccion.slides.length}</span>
              </p>
            </div>

            <!-- Círculo de progreso más pequeño -->
            <div class="w-16 h-16 relative flex-shrink-0">
              <svg class="transform -rotate-90" width="64" height="64">
                <circle
                  cx="32"
                  cy="32"
                  r="28"
                  stroke="rgb(30 41 59)"
                  stroke-width="4"
                  fill="none"
                />
                <circle
                  cx="32"
                  cy="32"
                  r="28"
                  stroke="url(#gradient)"
                  stroke-width="4"
                  fill="none"
                  stroke-dasharray="175.929"
                  stroke-dashoffset={175.929 * (1 - progressPercentage / 100)}
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
                <span class="text-sm font-bold text-white">
                  {progressPercentage.toFixed(0)}%
                </span>
              </div>
            </div>
          </div>
        {/if}
      </div>

      <!-- Barra de progreso lineal (todos los tamaños) -->
      {#if showProgress}
        <div class="mt-3">
          <div class="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-blue-500 to-purple-500 transition-all duration-500"
              style="width: {progressPercentage}%"
            ></div>
          </div>
        </div>
      {/if}
    </div>
  </div>
  {/if}

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
  {#if showHeader}
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
  {/if}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
