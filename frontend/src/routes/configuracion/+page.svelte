<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { updateProfile } from '$lib/api/profiles';
  import StepIndicator from '$lib/components/onboarding/StepIndicator.svelte';
  import InterestChips from '$lib/components/onboarding/InterestChips.svelte';
  import LearningCards from '$lib/components/onboarding/LearningCards.svelte';

  // Current step (only showing steps 1 and 2)
  let currentStep = $state(1);
  const totalSteps = 2;

  // Form data
  let interesesSeleccionados = $state<string[]>([]);
  let profesionSoñada = $state('');
  let formatoPreferido = $state('');
  let tipoActividad = $state<string[]>([]);
  let canalPreferido = $state('');

  // UI state
  let isLoading = $state(false);
  let errorMessage = $state('');
  let isSaving = $state(false);

  const stepTitles = ['Tus Intereses', 'Cómo Aprendes'];

  // Show ALL options including disabled ones
  const formatosCards = [
    { id: 'texto', label: 'Texto', description: 'Lecturas y explicaciones escritas', iconPath: 'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z' },
    { id: 'interactivo', label: 'Interactivo', description: 'Ejercicios prácticos y hands-on', iconPath: 'M15.042 21.672L13.684 16.6m0 0l-2.51 2.225.569-9.47 5.227 7.917-3.286-.672zM12 2.25V4.5m5.834.166l-1.591 1.591M20.25 10.5H18M7.757 14.743l-1.59 1.59M6 10.5H3.75m4.007-4.243l-1.59-1.59' },
    { id: 'visual', label: 'Visual', description: 'Imágenes, diagramas y gráficos', iconPath: 'M2.036 12.322a1.012 1.012 0 010-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178z M15 12a3 3 0 11-6 0 3 3 0 016 0z', disabled: true }
  ];

  const actividadesCards = [
    { id: 'proyectos', label: 'Proyectos', description: 'Trabajos prácticos y aplicados', iconPath: 'M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z' },
    { id: 'lecturas', label: 'Lecturas', description: 'Textos y documentos', iconPath: 'M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25' },
    { id: 'juegos', label: 'Juegos', description: 'Aprendizaje gamificado', iconPath: 'M14.25 6.087c0-.355.186-.676.401-.959.221-.29.349-.634.349-1.003 0-1.036-1.007-1.875-2.25-1.875s-2.25.84-2.25 1.875c0 .369.128.713.349 1.003.215.283.401.604.401.959v0a.64.64 0 01-.657.643 48.39 48.39 0 01-4.163-.3c.186 1.613.293 3.25.315 4.907a.656.656 0 01-.658.663v0c-.355 0-.676-.186-.959-.401a1.647 1.647 0 00-1.003-.349c-1.036 0-1.875 1.007-1.875 2.25s.84 2.25 1.875 2.25c.369 0 .713-.128 1.003-.349.283-.215.604-.401.959-.401v0c.31 0 .555.26.532.57a48.039 48.039 0 01-.642 5.056c1.518.19 3.058.309 4.616.354a.64.64 0 00.657-.643v0c0-.355-.186-.676-.401-.959a1.647 1.647 0 01-.349-1.003c0-1.035 1.008-1.875 2.25-1.875 1.243 0 2.25.84 2.25 1.875 0 .369-.128.713-.349 1.003-.215.283-.4.604-.4.959v0c0 .333.277.599.61.58a48.1 48.1 0 005.427-.63 48.05 48.05 0 00.582-4.717.532.532 0 00-.533-.57v0c-.355 0-.676.186-.959.401-.29.221-.634.349-1.003.349-1.035 0-1.875-1.007-1.875-2.25s.84-2.25 1.875-2.25c.37 0 .713.128 1.003.349.283.215.604.401.96.401v0a.656.656 0 00.658-.663 48.422 48.422 0 00-.37-5.36c-1.886.342-3.81.574-5.766.689a.578.578 0 01-.61-.58v0z', disabled: true }
  ];

  const canalesCards = [
    { id: 'audio', label: 'Audio', description: 'Podcasts y explicaciones de voz', iconPath: 'M12 18.75a6 6 0 006-6v-1.5m-6 7.5a6 6 0 01-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 01-3-3V4.5a3 3 0 116 0v8.25a3 3 0 01-3 3z' },
    { id: 'video', label: 'Video', description: 'Tutoriales y clases en video', iconPath: 'M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z', disabled: true },
    { id: 'animaciones', label: 'Animaciones', description: 'Contenido visual animado', iconPath: 'M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 00-2.456 2.456zM16.894 20.567L16.5 21.75l-.394-1.183a2.25 2.25 0 00-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 001.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 001.423 1.423l1.183.394-1.183.394a2.25 2.25 0 00-1.423 1.423z', disabled: true }
  ];

  // Load current profile data
  onMount(async () => {
    isLoading = true;
    try {
      await auth.loadProfile();
      if (auth.profile?.profile_data) {
        const data = auth.profile.profile_data;

        // Load interests
        if (data.intereses_personales) {
          interesesSeleccionados = data.intereses_personales.temas || [];
          profesionSoñada = data.intereses_personales.profesion_soñada || '';
        }

        // Load learning preferences
        if (data.preferencias_aprendizaje) {
          formatoPreferido = data.preferencias_aprendizaje.formato_preferido || '';
          tipoActividad = data.preferencias_aprendizaje.tipo_actividad || [];
          canalPreferido = data.preferencias_aprendizaje.canal_preferido || '';
        }
      }
    } catch (err) {
      console.error('Error loading profile:', err);
      errorMessage = 'Error al cargar tu perfil';
    } finally {
      isLoading = false;
    }

    // Animate entrance
    gsap.from('.config-container', {
      duration: 0.8,
      y: 50,
      opacity: 0,
      scale: 0.95,
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
      if (interesesSeleccionados.length === 0) {
        errorMessage = 'Por favor selecciona al menos un interés';
        return false;
      }
      if (!profesionSoñada || profesionSoñada.trim().length < 3) {
        errorMessage = 'Por favor ingresa tu profesión soñada';
        return false;
      }
    }

    if (currentStep === 2) {
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

    return true;
  }

  async function handleSave() {
    if (!validateCurrentStep()) {
      return;
    }

    isSaving = true;
    errorMessage = '';

    try {
      if (!auth.user?.id) {
        errorMessage = 'No se encontró el usuario';
        return;
      }

      if (!auth.profile) {
        errorMessage = 'No se encontró el perfil';
        return;
      }

      const updatedProfileData = {
        ...auth.profile.profile_data,
        intereses_personales: {
          temas: interesesSeleccionados,
          profesion_soñada: profesionSoñada
        },
        preferencias_aprendizaje: {
          formato_preferido: formatoPreferido,
          tipo_actividad: tipoActividad,
          canal_preferido: canalPreferido
        },
        ultima_actualizacion: new Date().toISOString()
      };

      const result = await updateProfile(auth.user.id, {
        profile_data: updatedProfileData
      });

      if (result.success) {
        // Reload profile
        await auth.loadProfile();

        // Success animation and redirect
        gsap.to('.config-container', {
          duration: 0.5,
          scale: 0.95,
          opacity: 0,
          ease: 'power2.in',
          onComplete: () => {
            goto('/');
          }
        });
      } else {
        errorMessage = result.error || 'Error al actualizar las preferencias';
      }
    } catch (err) {
      console.error('Error updating preferences:', err);
      errorMessage = 'Error al guardar los cambios';
    } finally {
      isSaving = false;
    }
  }

  function handleCancel() {
    gsap.to('.config-container', {
      duration: 0.4,
      scale: 0.95,
      opacity: 0,
      ease: 'power2.in',
      onComplete: () => {
        goto('/');
      }
    });
  }
</script>

<svelte:head>
  <title>Configuración - Lumera App</title>
</svelte:head>

<div class="relative min-h-screen bg-canvas-950 flex items-center justify-center p-4">
  <!-- Main Content -->
  <div class="relative z-10 w-full max-w-3xl">
    <div class="config-container relative group">
      <!-- Animated glow effect -->
      <div class="absolute -inset-1 bg-gradient-to-r from-lumera-500 via-focus-500 to-purple-500 rounded-3xl blur-lg opacity-25 group-hover:opacity-40 transition-opacity duration-300"></div>

      <div class="relative bg-canvas-800/90 backdrop-blur-xl rounded-3xl shadow-2xl p-6 border-2 border-white/10">
        <!-- Header -->
        <div class="mb-4">
          <div class="flex items-center justify-between mb-2">
            <h1 class="text-2xl font-bold text-white" style="text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);">
              Configuración de Preferencias
            </h1>
            <button
              onclick={handleCancel}
              class="h-9 w-9 rounded-full bg-canvas-900 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
            >
              ✕
            </button>
          </div>
          <p class="text-slate-400 text-sm">Actualiza tus intereses y preferencias de aprendizaje</p>
        </div>

        {#if isLoading}
          <div class="flex items-center justify-center py-20">
            <svg class="animate-spin h-10 w-10 text-lumera-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
        {:else}
          <!-- Step Indicator -->
          <StepIndicator {currentStep} {totalSteps} {stepTitles} />

          <!-- Error Message -->
          {#if errorMessage}
            <div class="mb-6 p-4 bg-red-500/10 border-2 border-red-500/30 rounded-xl backdrop-blur-sm">
              <div class="flex items-center gap-2">
                <svg class="w-5 h-5 text-red-400 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
                </svg>
                <p class="text-red-400 text-sm font-medium">{errorMessage}</p>
              </div>
            </div>
          {/if}

          <!-- Step Content -->
          <div class="step-content min-h-[280px]">
            <!-- Step 1: Intereses Personales -->
            {#if currentStep === 1}
              <div class="space-y-4">
                <div>
                  <h2 class="text-sm font-semibold text-white mb-1.5">
                    ¿Qué te apasiona? <span class="text-slate-500 text-xs font-normal">(selecciona varios)</span>
                  </h2>
                  <InterestChips
                    bind:selected={interesesSeleccionados}
                    onSelect={(interests) => interesesSeleccionados = interests}
                    compact={true}
                  />
                </div>

                <div>
                  <label for="profesion" class="block text-sm font-semibold text-white mb-2">
                    ¿Cuál es tu profesión soñada?
                  </label>
                  <input
                    id="profesion"
                    type="text"
                    bind:value={profesionSoñada}
                    class="w-full px-4 py-2.5 bg-canvas-900/60 border-2 border-canvas-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-lumera-500/20 transition-all text-sm hover:border-canvas-600"
                    placeholder="Ej: Desarrollador de Videojuegos"
                  />
                </div>
              </div>
            {/if}

            <!-- Step 2: Preferencias de Aprendizaje -->
            {#if currentStep === 2}
              <div class="space-y-2.5">
                <div>
                  <h2 class="text-sm font-semibold text-white mb-1.5">
                    ¿Cómo prefieres aprender?
                  </h2>
                  <LearningCards
                    cards={formatosCards}
                    bind:selected={formatoPreferido}
                    onSelect={(value) => formatoPreferido = value}
                    compact={true}
                  />
                </div>

                <div>
                  <h2 class="text-sm font-semibold text-white mb-1.5">
                    ¿Qué tipo de actividades te gustan? <span class="text-slate-500 text-xs font-normal">(varias)</span>
                  </h2>
                  <LearningCards
                    cards={actividadesCards}
                    bind:selected={tipoActividad}
                    onSelect={(value) => tipoActividad = value}
                    multiSelect={true}
                    compact={true}
                  />
                </div>

                <div>
                  <h2 class="text-sm font-semibold text-white mb-1.5">
                    ¿Qué canal prefieres?
                  </h2>
                  <LearningCards
                    cards={canalesCards}
                    bind:selected={canalPreferido}
                    onSelect={(value) => canalPreferido = value}
                    compact={true}
                  />
                </div>
              </div>
            {/if}
          </div>

          <!-- Navigation Buttons -->
          <div class="flex justify-between items-center mt-4 pt-3 border-t border-canvas-700/50">
            <button
              onclick={prevStep}
              disabled={currentStep === 1}
              class="px-4 py-2 text-sm rounded-lg bg-canvas-900/60 text-slate-300 font-semibold hover:bg-canvas-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all border border-canvas-700"
            >
              ← Anterior
            </button>

            <div class="flex gap-2">
              <button
                onclick={handleCancel}
                class="px-4 py-2 text-sm rounded-lg bg-canvas-900/60 text-slate-300 font-semibold hover:bg-canvas-700 transition-all border border-canvas-700"
              >
                Cancelar
              </button>

              {#if currentStep < totalSteps}
                <button
                  onclick={nextStep}
                  class="px-6 py-2 text-sm rounded-lg bg-[#E1E1E1] hover:bg-[#CCCCCC] text-canvas-900 font-bold transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105"
                >
                  Siguiente →
                </button>
              {:else}
                <button
                  onclick={handleSave}
                  disabled={isSaving}
                  class="px-6 py-2 text-sm rounded-lg bg-[#E1E1E1] hover:bg-[#CCCCCC] disabled:bg-[#E1E1E1]/50 text-canvas-900 font-bold transition-all duration-300 shadow-lg hover:shadow-xl disabled:cursor-not-allowed flex items-center gap-2 hover:scale-105 disabled:hover:scale-100"
                >
                  {#if isSaving}
                    <svg class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Guardando...
                  {:else}
                    Guardar Cambios
                  {/if}
                </button>
              {/if}
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
