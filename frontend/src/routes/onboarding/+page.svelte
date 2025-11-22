<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { createProfile } from '$lib/api/profiles';
  import StepIndicator from '$lib/components/onboarding/StepIndicator.svelte';
  import InterestChips from '$lib/components/onboarding/InterestChips.svelte';
  import LearningCards from '$lib/components/onboarding/LearningCards.svelte';
  import EmojiScale from '$lib/components/onboarding/EmojiScale.svelte';

  // Current step
  let currentStep = $state(1);
  const totalSteps = 4;

  // Form data
  let cursoActual = $state('');
  let interesesSeleccionados = $state([]);
  let profesionSoñada = $state('');
  let formatoPreferido = $state('');
  let tipoActividad = $state([]);
  let canalPreferido = $state('');
  let interesActual = $state('');
  let orientacionMeta = $state('');

  // UI state
  let isLoading = $state(false);
  let errorMessage = $state('');

  const stepTitles = ['Info Básica', 'Tus Intereses', 'Cómo Aprendes', 'Motivación'];

  const cursosDisponibles = [
    '1ro Medio',
    '2do Medio',
    '3ro Medio',
    '4to Medio'
  ];

  const formatosCards = [
    { id: 'visual', label: 'Visual', description: 'Imágenes, diagramas y gráficos', icon: '👁️' },
    { id: 'texto', label: 'Texto', description: 'Lecturas y explicaciones escritas', icon: '📝' },
    { id: 'interactivo', label: 'Interactivo', description: 'Ejercicios prácticos y hands-on', icon: '🎯' }
  ];

  const actividadesCards = [
    { id: 'proyectos', label: 'Proyectos', description: 'Trabajos prácticos y aplicados', icon: '🛠️' },
    { id: 'juegos', label: 'Juegos', description: 'Aprendizaje gamificado', icon: '🎮' },
    { id: 'lecturas', label: 'Lecturas', description: 'Textos y documentos', icon: '📖' },
    { id: 'videos', label: 'Videos', description: 'Contenido audiovisual', icon: '🎥' }
  ];

  const canalesCards = [
    { id: 'video', label: 'Video', description: 'Tutoriales y clases en video', icon: '📹' },
    { id: 'audio', label: 'Audio', description: 'Podcasts y explicaciones de voz', icon: '🎧' },
    { id: 'animaciones', label: 'Animaciones', description: 'Contenido visual animado', icon: '🎬' }
  ];

  const motivacionOptions = [
    { value: 'bajo', emoji: '😐', label: 'Poco' },
    { value: 'medio', emoji: '🙂', label: 'Moderado' },
    { value: 'alto', emoji: '😃', label: 'Muy motivado' }
  ];

  const metaCards = [
    { id: 'maestria', label: 'Dominar el Tema', description: 'Quiero entender profundamente', icon: '🎓' },
    { id: 'desempeño', label: 'Buenas Notas', description: 'Quiero obtener altas calificaciones', icon: '⭐' }
  ];

  onMount(() => {
    // Animate entrance
    gsap.from('.onboarding-container', {
      duration: 0.8,
      y: 50,
      opacity: 0,
      ease: 'power3.out'
    });
  });

  function animateStepTransition(direction: 'forward' | 'backward') {
    const tl = gsap.timeline();
    tl.to('.step-content', {
      duration: 0.3,
      x: direction === 'forward' ? -30 : 30,
      opacity: 0,
      ease: 'power2.in'
    })
    .set('.step-content', { x: direction === 'forward' ? 30 : -30 })
    .to('.step-content', {
      duration: 0.3,
      x: 0,
      opacity: 1,
      ease: 'power2.out'
    });
  }

  function nextStep() {
    if (validateCurrentStep()) {
      animateStepTransition('forward');
      currentStep = Math.min(currentStep + 1, totalSteps);
      errorMessage = '';
    }
  }

  function prevStep() {
    animateStepTransition('backward');
    currentStep = Math.max(currentStep - 1, 1);
    errorMessage = '';
  }

  function validateCurrentStep(): boolean {
    errorMessage = '';

    if (currentStep === 1) {
      if (!cursoActual) {
        errorMessage = 'Por favor selecciona tu curso actual';
        return false;
      }
    }

    if (currentStep === 2) {
      if (interesesSeleccionados.length === 0) {
        errorMessage = 'Por favor selecciona al menos un interés';
        return false;
      }
      if (!profesionSoñada || profesionSoñada.trim().length < 3) {
        errorMessage = 'Por favor ingresa tu profesión soñada';
        return false;
      }
    }

    if (currentStep === 3) {
      if (!formatoPreferido) {
        errorMessage = 'Por favor selecciona un formato preferido';
        return false;
      }
      if (tipoActividad.length === 0) {
        errorMessage = 'Por favor selecciona al menos un tipo de actividad';
        return false;
      }
      if (!canalPreferido) {
        errorMessage = 'Por favor selecciona un canal preferido';
        return false;
      }
    }

    if (currentStep === 4) {
      if (!interesActual) {
        errorMessage = 'Por favor indica tu nivel de motivación';
        return false;
      }
      if (!orientacionMeta) {
        errorMessage = 'Por favor selecciona tu orientación de meta';
        return false;
      }
    }

    return true;
  }

  async function handleSubmit() {
    if (!validateCurrentStep()) {
      return;
    }

    isLoading = true;
    errorMessage = '';

    const profileData = {
      user_id: auth.user?.id,
      curso_actual: cursoActual,
      profile_data: {
        intereses_personales: {
          temas: interesesSeleccionados,
          profesion_soñada: profesionSoñada
        },
        preferencias_aprendizaje: {
          formato_preferido: formatoPreferido,
          tipo_actividad: tipoActividad,
          canal_preferido: canalPreferido
        },
        motivacion: {
          interes_actual: interesActual,
          orientacion_meta: orientacionMeta
        },
        ultima_actualizacion: new Date().toISOString()
      }
    };

    const result = await createProfile(profileData);

    if (result.success) {
      // Success animation and redirect to diagnostic test
      gsap.to('.onboarding-container', {
        duration: 0.5,
        scale: 0.95,
        opacity: 0,
        ease: 'power2.in',
        onComplete: () => {
          goto('/diagnostico');
        }
      });
    } else {
      errorMessage = result.error || 'Error al crear el perfil. Por favor intenta de nuevo.';
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Bienvenido - Lumera App</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950 flex items-center justify-center p-4">
  <div class="onboarding-container bg-canvas-900/40 backdrop-blur-xl rounded-2xl shadow-2xl p-8 max-w-3xl w-full border border-white/10">
    <!-- Header -->
    <div class="text-center mb-8">
      <h1 class="text-4xl font-bold text-white mb-2">¡Bienvenido a Lumera! 🎓</h1>
      <p class="text-slate-400">Cuéntanos un poco sobre ti para personalizar tu experiencia</p>
    </div>

    <!-- Step Indicator -->
    <StepIndicator {currentStep} {totalSteps} {stepTitles} />

    <!-- Error Message -->
    {#if errorMessage}
      <div class="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
        <p class="text-red-400 text-sm">{errorMessage}</p>
      </div>
    {/if}

    <!-- Step Content -->
    <div class="step-content min-h-[400px]">
      <!-- Step 1: Información Básica -->
      {#if currentStep === 1}
        <div class="space-y-6">
          <div>
            <label for="curso" class="block text-lg font-semibold text-white mb-3">
              ¿En qué curso estás? 📚
            </label>
            <select
              id="curso"
              bind:value={cursoActual}
              class="w-full px-6 py-4 bg-canvas-950/50 border border-slate-700 rounded-xl text-white focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all text-lg"
            >
              <option value="">Selecciona tu curso</option>
              {#each cursosDisponibles as curso}
                <option value={curso}>{curso}</option>
              {/each}
            </select>
          </div>
        </div>
      {/if}

      <!-- Step 2: Intereses Personales -->
      {#if currentStep === 2}
        <div class="space-y-6">
          <div>
            <h2 class="text-lg font-semibold text-white mb-4">
              ¿Qué te apasiona? ❤️
            </h2>
            <p class="text-slate-400 text-sm mb-4">Selecciona todos los que quieras</p>
            <InterestChips
              bind:selected={interesesSeleccionados}
              onSelect={(interests) => interesesSeleccionados = interests}
            />
          </div>

          <div>
            <label for="profesion" class="block text-lg font-semibold text-white mb-3">
              ¿Cuál es tu profesión soñada? 💭
            </label>
            <input
              id="profesion"
              type="text"
              bind:value={profesionSoñada}
              class="w-full px-6 py-4 bg-canvas-950/50 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all text-lg"
              placeholder="Ej: Desarrollador de Videojuegos"
            />
          </div>
        </div>
      {/if}

      <!-- Step 3: Preferencias de Aprendizaje -->
      {#if currentStep === 3}
        <div class="space-y-8">
          <div>
            <h2 class="text-lg font-semibold text-white mb-4">
              ¿Cómo prefieres aprender? 🎯
            </h2>
            <LearningCards
              cards={formatosCards}
              bind:selected={formatoPreferido}
              onSelect={(value) => formatoPreferido = value}
            />
          </div>

          <div>
            <h2 class="text-lg font-semibold text-white mb-4">
              ¿Qué tipo de actividades te gustan? 🎨
            </h2>
            <p class="text-slate-400 text-sm mb-4">Puedes seleccionar varias</p>
            <LearningCards
              cards={actividadesCards}
              bind:selected={tipoActividad}
              onSelect={(value) => tipoActividad = value}
              multiSelect={true}
            />
          </div>

          <div>
            <h2 class="text-lg font-semibold text-white mb-4">
              ¿Qué canal prefieres? 📺
            </h2>
            <LearningCards
              cards={canalesCards}
              bind:selected={canalPreferido}
              onSelect={(value) => canalPreferido = value}
            />
          </div>
        </div>
      {/if}

      <!-- Step 4: Motivación -->
      {#if currentStep === 4}
        <div class="space-y-8">
          <div>
            <h2 class="text-lg font-semibold text-white mb-4 text-center">
              ¿Qué tan motivado estás por aprender? 💪
            </h2>
            <EmojiScale
              options={motivacionOptions}
              bind:selected={interesActual}
              onSelect={(value) => interesActual = value}
            />
          </div>

          <div>
            <h2 class="text-lg font-semibold text-white mb-4">
              ¿Qué buscas al estudiar? 🎯
            </h2>
            <LearningCards
              cards={metaCards}
              bind:selected={orientacionMeta}
              onSelect={(value) => orientacionMeta = value}
            />
          </div>
        </div>
      {/if}
    </div>

    <!-- Navigation Buttons -->
    <div class="flex justify-between items-center mt-8 pt-6 border-t border-slate-800">
      <button
        onclick={prevStep}
        disabled={currentStep === 1}
        class="px-6 py-3 rounded-xl bg-canvas-800 text-slate-300 font-semibold hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        ← Anterior
      </button>

      {#if currentStep < totalSteps}
        <button
          onclick={nextStep}
          class="px-6 py-3 rounded-xl bg-lumera-600 hover:bg-lumera-700 text-white font-semibold transition-all shadow-lg hover:shadow-indigo-500/50"
        >
          Siguiente →
        </button>
      {:else}
        <button
          onclick={handleSubmit}
          disabled={isLoading}
          class="px-8 py-3 rounded-xl bg-gradient-to-r from-lumera-600 to-focus-600 hover:from-indigo-700 hover:to-cyan-700 text-white font-bold transition-all shadow-lg hover:shadow-indigo-500/50 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          {#if isLoading}
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Guardando...
          {:else}
            ¡Comenzar! 🚀
          {/if}
        </button>
      {/if}
    </div>
  </div>
</div>
