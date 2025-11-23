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
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';
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
  const student = {
    name: auth.user?.name || 'Student',
    grade: '2° Medio',
    level: 12,
    xp: 75,
    streak: 5
  };

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

      plan = await response.json();
      totalSlides = plan?.components?.length || 0;

      // Mark plan as started if not already started
      if (plan && !plan.fecha_inicio) {
        const startResult = await startLearningPlan(planId);
        if (startResult.success && startResult.plan) {
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
  <header class="sticky top-0 z-50 border-b border-canvas-700/20 bg-canvas-800/95 backdrop-blur-md">
    <div class="px-6 py-3 flex items-center justify-between gap-4">
      <!-- Left: Profile -->
      <button
        onclick={() => isPlayerProfileOpen = true}
        class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-canvas-900/60 transition-all group"
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

      <!-- Center: Progress (replaces XP/streak/coins) -->
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

      <!-- Right: Navigation Icons -->
      <div class="flex items-center gap-2">
        <!-- Quest Button -->
        <button
          onclick={() => isCurrentQuestOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Quest Actual"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        </button>

        <!-- Mission Board Button -->
        <button
          onclick={() => isMissionBoardOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Mission Board"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Actividad Reciente"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Live Events"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
          class="relative p-2 rounded-lg hover:bg-canvas-900/80 transition-all duration-200 group"
          title="Tu Progreso"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          {#if dashboardStore.subjects.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-lumera-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.subjects.length}
            </span>
          {/if}
        </button>

        <!-- Close/Exit Button -->
        <button
          onclick={() => goto(`/materias/${materiaId}/objetivos`)}
          class="relative p-2 rounded-lg hover:bg-red-900/60 transition-all duration-200 group"
          title="Salir"
        >
          <svg class="w-6 h-6 text-slate-300 group-hover:text-red-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  </header>

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
    {#if !isLoading && lesson}
      <LessonPlayer
        leccion={lesson}
        showProgress={false}
        showHeader={false}
        onComplete={handlePlanComplete}
        onSlideChange={handleSlideChange}
      />
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
