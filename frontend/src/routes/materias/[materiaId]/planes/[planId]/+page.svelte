<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { getPlanById, startLearningPlan, completeLearningPlan, type LearningPlan } from '$lib/api/learningPlans';
  import LessonPlayer from '$lib/components/slides/LessonPlayer.svelte';
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

  // Route params
  let materiaId = $state(0);
  let planId = $state(0);

  // Learning Plan State
  let plan = $state<LearningPlan | null>(null);
  let currentSlideIndex = $state(0);
  let totalSlides = $state(0);

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
  let isLoading = $state(false);
  let errorMessage = $state('');

  // Student data
  const student = $derived({
    name: auth.user?.name || 'Student',
    grade: auth.user?.curso_actual || '2° Medio'
  });

  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  // Progress calculation
  const progreso = $derived(totalSlides > 0
    ? Math.round(((currentSlideIndex + 1) / totalSlides) * 100)
    : 0);

  // Convert learning plan to lesson format
  const lesson = $derived(
    plan ? {
      leccionId: `plan-${plan.id}`,
      titulo: plan.titulo,
      materia: 'learning-plan',
      slides: [...(plan.components || [])]  // Crear copia del array
        .sort((a, b) => a.orden - b.orden)
        .map(component => ({
          orden: component.orden,
          tipo: component.tipo_componente,
          props: component.contenido_props || {},
          componentId: component.id,
          estado: component.estado
        }))
    } : null
  );

  // Load user profile
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

  // Handle avatar change
  function handleAvatarChanged(avatar: CustomizationItem) {
    currentAvatar = avatar;
  }

  // Load learning plan
  async function loadPlan() {
    if (!planId) return;

    isLoading = true;
    errorMessage = '';

    try {
      const response = await fetch(`/api/learning-plans/${planId}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });

      if (!response.ok) {
        throw new Error('Error al cargar el plan');
      }

      const planData = await response.json();

      // Ensure components exist and is an array
      if (!planData.components) {
        planData.components = [];
      }

      plan = planData;
      totalSlides = plan?.components?.length || 0;

      // Mark plan as started if not already started
      if (plan && !plan.fecha_inicio) {
        const startResult = await startLearningPlan(planId);
        if (startResult.success && startResult.plan) {
          // Ensure the returned plan also has components
          if (!startResult.plan.components) {
            startResult.plan.components = planData.components;
          }
          plan = startResult.plan;
        }
      }
    } catch (error) {
      console.error('Error loading plan:', error);
      errorMessage = 'No se pudo cargar el plan de aprendizaje';
    } finally {
      isLoading = false;
    }
  }

  // Handle lesson completion
  async function handlePlanComplete(completionData: any) {
    console.log('Plan completado:', completionData);

    // Mark plan as completed in backend
    if (planId) {
      const completeResult = await completeLearningPlan(planId);
      if (completeResult.success) {
        console.log('Plan marcado como completado en el backend');
      } else {
        console.error('Error al marcar plan como completado:', completeResult.error);
      }
    }

    // Navigate back to objectives
    goto(`/materias/${materiaId}/objetivos`);
  }

  // Handle slide change
  function handleSlideChange(data: { slideIndex: number; slideType: string }) {
    currentSlideIndex = data.slideIndex;
    // TODO: Update component progress in backend
    console.log('Slide changed:', data);
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

  // Initialize on mount
  onMount(async () => {
    loadUserProfile();
    loadEquippedAvatar();

    const unsubscribe = page.subscribe(($page) => {
      materiaId = parseInt($page.params.materiaId || '0');
      planId = parseInt($page.params.planId || '0');
    });

    await loadPlan();

    return () => {
      if (unsubscribe) unsubscribe();
    };
  });
</script>

<svelte:head>
  <title>{plan?.titulo || 'Plan de Aprendizaje'} - Lumera App</title>
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
      <div class="flex-1 max-w-md">
        <div class="flex items-center justify-between gap-2 text-xs mb-1.5">
          {#if plan}
            <span class="text-lumera-400 font-medium truncate">📚 {plan.titulo}</span>
            <span class="text-slate-400 flex-shrink-0">{currentSlideIndex + 1}/{totalSlides}</span>
          {:else}
            <span class="text-slate-400">Cargando...</span>
            <span class="text-slate-400">--</span>
          {/if}
        </div>
        <div class="h-2 w-full bg-canvas-900 rounded-full overflow-hidden border border-canvas-700">
          <div class="h-full bg-gradient-to-r from-lumera-500 to-focus-500 transition-all duration-500" style="width: {progreso}%"></div>
        </div>
      </div>
    {/snippet}
  </AppHeader>

  <!-- Main Content -->
  <main class="relative">
    <!-- Error Message -->
    {#if errorMessage}
      <div class="max-w-4xl mx-auto px-6 pt-6">
        <div class="mb-6 p-4 bg-red-50 border-2 border-red-200 rounded-xl">
          <p class="text-red-700">{errorMessage}</p>
        </div>
      </div>
    {/if}

    <!-- Loading State -->
    {#if isLoading}
      <div class="flex items-center justify-center py-20">
        <div class="text-center">
          <svg class="animate-spin h-12 w-12 mx-auto mb-4 text-lumera-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="text-slate-600">Cargando plan de aprendizaje...</p>
        </div>
      </div>
    {/if}

    <!-- Lesson Player -->
    {#if !isLoading && lesson && lesson.slides && lesson.slides.length > 0}
      <LessonPlayer
        leccion={lesson}
        showProgress={false}
        showHeader={false}
        onComplete={handlePlanComplete}
        onSlideChange={handleSlideChange}
      />
    {:else if !isLoading && plan && (!lesson?.slides || lesson.slides.length === 0)}
      <div class="flex items-center justify-center py-20">
        <div class="text-center">
          <p class="text-slate-400 mb-4">No hay contenido disponible en este plan todavía.</p>
          <button
            onclick={() => goto(`/materias/${materiaId}/objetivos`)}
            class="px-6 py-3 rounded-xl font-semibold bg-[#E1E1E1] hover:bg-[#CCCCCC] text-canvas-900 transition-all duration-300"
          >
            Volver a Objetivos
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
