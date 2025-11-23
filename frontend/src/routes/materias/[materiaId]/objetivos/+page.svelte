<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { customizationStore } from '$lib/stores/customization.svelte';
  import { getObjetivosAprendizaje, getOABloomObjectiveId, getPlanByOA, generatePlan } from '$lib/api/learningPlans';
  import { getDomainLevelInfo } from '$lib/constants/subjects';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import { findCursoByName } from '$lib/api/courses';
  import type { CustomizationItem } from '$lib/api/customization';
  import { getProfile } from '$lib/api/profiles';
  import AppHeader from '$lib/components/common/AppHeader.svelte';
  import PlayerProfilePanel from '$lib/components/dashboard/PlayerProfilePanel.svelte';
  import CurrentQuestModal from '$lib/components/dashboard/CurrentQuestModal.svelte';
  import MissionBoardModal from '$lib/components/dashboard/MissionBoardModal.svelte';
  import RecentActivityModal from '$lib/components/dashboard/RecentActivityModal.svelte';
  import LiveEventsModal from '$lib/components/dashboard/LiveEventsModal.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';

  // New components
  import OAStatsDashboard from '$lib/components/objectives/OAStatsDashboard.svelte';
  import OACategoryTabs from '$lib/components/objectives/OACategoryTabs.svelte';
  import OASearchBar from '$lib/components/objectives/OASearchBar.svelte';
  import OAList from '$lib/components/objectives/OAList.svelte';

  // State
  let materiaId = $state(0);
  let materiaInfo = $state<any>(null);
  let oas = $state<any[]>([]);
  let filteredOAs = $state<any[]>([]);
  let searchQuery = $state('');
  let sortBy = $state('titulo');
  let selectedCategory = $state('Todos');
  let isLoading = $state(true);
  let errorMessage = $state('');
  let diagnosticLevel = $state(0);
  let generatingPlanFor = $state<number | null>(null);
  let userProfile = $state(null);

  // Modal states for header navigation
  let isPlayerProfileOpen = $state(false);
  let isCurrentQuestOpen = $state(false);
  let isMissionBoardOpen = $state(false);
  let isActivityModalOpen = $state(false);
  let isLiveEventsOpen = $state(false);
  let isProgressPanelOpen = $state(false);
  let diagnosticLevels = $state<Record<number, number>>({});
  let subjects = $state<any[]>([]);
  let activeTab = $state('daily');

  // Student info for modals
  const student = {
    name: auth.user?.name || 'Student',
    grade: '2° Medio',
    level: 12,
    xp: 75,
    streak: 5
  };

  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  const levelInfo = $derived(getDomainLevelInfo(diagnosticLevel));

  // Derived data for stats dashboard
  const categoryStats = $derived(() => {
    const stats = new Map();

    oas.forEach(oa => {
      const cat = oa.categoria || 'General';
      if (!stats.has(cat)) {
        stats.set(cat, { categoria: cat, completed: 0, total: 0 });
      }
      const catStat = stats.get(cat);
      catStat.total++;
      // TODO: Add actual completion logic based on progress
      // For now, using placeholder logic
      catStat.completed += Math.random() > 0.7 ? 1 : 0;
    });

    return Array.from(stats.values());
  });

  const categoryCounts = $derived(() => {
    const counts = new Map();

    oas.forEach(oa => {
      const cat = oa.categoria || 'General';
      counts.set(cat, (counts.get(cat) || 0) + 1);
    });

    return Array.from(counts.entries()).map(([categoria, count]) => ({ categoria, count }));
  });

  // Get completed OAs count (placeholder - TODO: implement real progress tracking)
  const completedOAs = $derived(Math.floor(oas.length * 0.65)); // 65% completion placeholder

  // Get recommended OA based on diagnostic level
  const recommendedOA = $derived(() => {
    if (oas.length === 0) return null;

    // Find first OA with completion < 100%
    // For now, return first OA as placeholder
    return oas[0] || null;
  });

  // Get materia ID from URL
  $effect(() => {
    const unsubscribe = page.subscribe(($page) => {
      materiaId = parseInt($page.params.materiaId || '0');
    });

    return unsubscribe;
  });

  // Load materia info
  async function loadMateria() {
    try {
      const response = await fetch(`/api/materias/${materiaId}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });

      if (!response.ok) {
        throw new Error('Error al cargar materia');
      }

      materiaInfo = await response.json();
    } catch (error) {
      console.error('Error loading materia:', error);
      errorMessage = 'No se pudo cargar la información de la materia';
    }
  }

  // Load user profile
  async function loadUserProfile() {
    if (!auth.user?.id) return;
    const result = await getProfile(auth.user.id);
    if (result.success && result.profile) {
      userProfile = result.profile;
    }
  }

  // Load diagnostic level
  async function loadDiagnosticLevel() {
    try {
      const response = await fetch('/api/diagnostic-sessions?estado=completado', {
        headers: {
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });

      if (!response.ok) return;

      const sessions = await response.json();

      // Find session for this materia
      const session = sessions.find((s: any) => s.materia_id === materiaId);

      if (session?.estrategia?.average_bloom_level) {
        diagnosticLevel = Math.round(session.estrategia.average_bloom_level);
      }
    } catch (error) {
      console.error('Error loading diagnostic level:', error);
    }
  }

  // Load diagnostic levels for all materias (for progress panel)
  async function loadDiagnosticLevels() {
    try {
      const response = await fetch('/api/diagnostic-sessions?estado=completado', {
        headers: {
          'Authorization': auth.token ? `Bearer ${auth.token}` : ''
        }
      });

      if (!response.ok) return;

      const sessions = await response.json();

      // Build levels map
      const levels: Record<number, number> = {};
      sessions.forEach((session: any) => {
        if (session.estrategia?.average_bloom_level) {
          levels[session.materia_id] = Math.round(session.estrategia.average_bloom_level);
        }
      });

      diagnosticLevels = levels;

      // Load subjects if not loaded
      if (subjects.length === 0 && userProfile) {
        const curso = await findCursoByName(userProfile.curso_actual);
        if (curso.success && curso.curso?.materias) {
          subjects = materiasToSubjects(curso.curso.materias);
        }
      }
    } catch (error) {
      console.error('Error loading diagnostic levels:', error);
    }
  }

  // Load OAs for materia
  async function loadOAs() {
    isLoading = true;
    const result = await getObjetivosAprendizaje(materiaId);

    if (result.success && result.oas) {
      oas = result.oas;
      // filteredOAs will be updated by the $effect() reactively
    } else {
      errorMessage = result.error || 'No se pudieron cargar los objetivos de aprendizaje';
    }

    isLoading = false;
  }

  // Load plans data - check which OAs have learning plans
  async function loadPlansData() {
    if (diagnosticLevel === 0 || oas.length === 0) {
      plansData = {};
      return;
    }

    const data: Record<number, { hasPlan: boolean; isCompleted: boolean }> = {};

    // Check plans for each OA in parallel
    await Promise.all(
      oas.map(async (oa) => {
        try {
          // Get the OA Bloom Objective ID for this OA and user's level
          const objectiveId = await getOABloomObjectiveId(oa.id, diagnosticLevel);

          if (objectiveId) {
            // Check if a plan exists
            const plan = await getPlanByOA(objectiveId);
            data[oa.id] = {
              hasPlan: plan !== null,
              isCompleted: plan?.completado || false
            };
          } else {
            data[oa.id] = { hasPlan: false, isCompleted: false };
          }
        } catch (error) {
          console.error(`Error checking plan for OA ${oa.id}:`, error);
          data[oa.id] = { hasPlan: false, isCompleted: false };
        }
      })
    );

    plansData = data;
  }

  // Filter and sort OAs
  $effect(() => {
    let filtered = oas;

    // Filter by category
    if (selectedCategory !== 'Todos') {
      filtered = filtered.filter(oa => (oa.categoria || 'General') === selectedCategory);
    }

    // Filter by search query
    if (searchQuery.trim() !== '') {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(oa =>
        oa.titulo.toLowerCase().includes(query) ||
        oa.descripcion.toLowerCase().includes(query)
      );
    }

    // Sort
    filtered = [...filtered].sort((a, b) => {
      switch (sortBy) {
        case 'titulo':
          return a.titulo.localeCompare(b.titulo);
        case 'progreso':
          // TODO: Sort by actual progress
          const progressA = progressData()[a.id] || 0;
          const progressB = progressData()[b.id] || 0;
          return progressB - progressA; // Descending order (highest progress first)
        default:
          return 0;
      }
    });

    filteredOAs = filtered;
  });

  // Generate learning plan for an OA
  async function handleGeneratePlan(oa: any) {
    if (diagnosticLevel === 0) {
      alert('Primero debes completar la evaluación diagnóstica de esta materia');
      goto(`/diagnostico/${materiaId}`);
      return;
    }

    generatingPlanFor = oa.id;

    try {
      // Get OA Bloom Objective ID based on user's level
      const objectiveId = await getOABloomObjectiveId(oa.id, diagnosticLevel);

      if (!objectiveId) {
        throw new Error('No se pudo encontrar el objetivo de Bloom para este OA');
      }

      // Check if plan already exists
      let plan = await getPlanByOA(objectiveId);

      // If not, generate it
      if (!plan) {
        const generateResult = await generatePlan(objectiveId);

        if (!generateResult.success) {
          throw new Error(generateResult.error || 'Error al generar el plan');
        }

        plan = generateResult.plan!;

        // Update plansData to reflect the new plan
        plansData = { ...plansData, [oa.id]: { hasPlan: true, isCompleted: false } };
      }

      // Navigate to plan page
      goto(`/materias/${materiaId}/planes/${plan.id}`);
    } catch (error) {
      console.error('Error generating plan:', error);
      alert(error instanceof Error ? error.message : 'Error al generar el plan de aprendizaje');
    } finally {
      generatingPlanFor = null;
    }
  }

  // View existing plan
  async function handleViewPlan(oa: any) {
    // Same logic as generate, but should have plan already
    handleGeneratePlan(oa);
  }

  // Practice OA (only available after completing the plan)
  async function handlePractice(oa: any) {
    // Navigate to practice page for this OA
    console.log('Iniciando práctica para OA:', oa.titulo);
    goto(`/materias/${materiaId}/practica/${oa.id}`);
  }

  // Start recommended OA
  function handleStartRecommended(oa: any) {
    handleGeneratePlan(oa);
  }

  // Handle avatar change
  function handleAvatarChanged(avatar: CustomizationItem) {
    customizationStore.setAvatar(avatar);
  }

  // Placeholder data for progress (TODO: implement real progress tracking)
  const progressData = $derived(() => {
    const data: Record<number, number> = {};
    oas.forEach(oa => {
      // Placeholder: random progress between 0-100
      data[oa.id] = Math.floor(Math.random() * 100);
    });
    return data;
  });

  // Real plans data - check which OAs have learning plans
  let plansData = $state<Record<number, { hasPlan: boolean; isCompleted: boolean }>>({});

  // Bloom levels for each OA (use diagnostic level as placeholder)
  const bloomLevels = $derived(() => {
    const data: Record<number, number> = {};
    oas.forEach(oa => {
      data[oa.id] = diagnosticLevel;
    });
    return data;
  });

  // Initialize
  onMount(async () => {
    loadUserProfile(); // Load profile in parallel
    customizationStore.loadAvatar(); // Load avatar from global store
    if (materiaId) {
      await loadMateria();
      await loadDiagnosticLevel();
      await loadOAs();
      // Load plans data after OAs and diagnostic level are loaded
      await loadPlansData();
    } else {
      errorMessage = 'ID de materia inválido';
    }
  });
</script>

<svelte:head>
  <title>{materiaInfo?.nombre || 'Objetivos de Aprendizaje'} - Lumera App</title>
</svelte:head>

<div class="min-h-screen bg-canvas-950 text-slate-900">
  <!-- App Header -->
  <AppHeader
    currentAvatar={customizationStore.currentAvatar}
    onProfileClick={() => isPlayerProfileOpen = true}
    onQuestClick={() => isCurrentQuestOpen = true}
    onMissionsClick={() => isMissionBoardOpen = true}
    onActivityClick={() => isActivityModalOpen = true}
    onLiveEventsClick={() => isLiveEventsOpen = true}
    onProgressClick={async () => {
      await loadDiagnosticLevels();
      isProgressPanelOpen = true;
    }}
  />

  <!-- Main Content -->
  <main class="max-w-7xl mx-auto px-6 py-8">
    <!-- Breadcrumb -->
    <div class="flex items-center gap-2 text-sm text-slate-600 mb-6">
      <a href="/" class="hover:text-lumera-600 transition-colors">Dashboard</a>
      <span>›</span>
      <span class="text-slate-900 font-medium">{materiaInfo?.nombre || 'Cargando...'}</span>
      <span>›</span>
      <span class="text-lumera-600">Objetivos de Aprendizaje</span>
    </div>

    {#if errorMessage}
      <div class="bg-red-50 border-2 border-red-200 rounded-xl p-4 text-red-700 mb-6">
        {errorMessage}
      </div>
    {/if}

    {#if isLoading}
      <div class="flex items-center justify-center py-20">
        <div class="text-center">
          <svg class="animate-spin h-12 w-12 mx-auto mb-4 text-lumera-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="text-slate-600">Cargando objetivos de aprendizaje...</p>
        </div>
      </div>
    {:else}
      <!-- Stats Dashboard -->
      <OAStatsDashboard
        materiaName={materiaInfo?.nombre || 'Materia'}
        totalOAs={oas.length}
        completedOAs={completedOAs}
        bloomLevel={diagnosticLevel}
        categoryStats={categoryStats()}
        recommendedOA={recommendedOA()}
        onStartRecommended={handleStartRecommended}
      />

      <!-- Category Tabs -->
      <OACategoryTabs
        categories={categoryCounts()}
        onCategoryChange={(cat) => selectedCategory = cat}
      />

      <!-- Search and Sort -->
      <OASearchBar
        onSearchChange={(query) => searchQuery = query}
        onSortChange={(sort) => sortBy = sort}
      />

      <!-- Results Count -->
      <div class="mb-4 text-sm text-slate-600">
        Mostrando {filteredOAs.length} de {oas.length} objetivos
        {#if selectedCategory !== 'Todos'}
          <span class="font-semibold">en {selectedCategory}</span>
        {/if}
      </div>

      <!-- OAs List -->
      <OAList
        oas={filteredOAs}
        bloomLevels={bloomLevels()}
        progressData={progressData()}
        plansData={plansData}
        {generatingPlanFor}
        onGeneratePlan={handleGeneratePlan}
        onViewPlan={handleViewPlan}
        onPractice={handlePractice}
      />
    {/if}
  </main>
</div>

<!-- Player Profile Panel -->
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

<!-- Current Quest Modal -->
<CurrentQuestModal
  isOpen={isCurrentQuestOpen}
  onClose={() => isCurrentQuestOpen = false}
  onResumeQuest={() => console.log('Resume quest clicked')}
/>

<!-- Mission Board Modal -->
<MissionBoardModal
  missions={dashboardStore.missions}
  {activeTab}
  onTabChange={(tab) => activeTab = tab}
  isOpen={isMissionBoardOpen}
  onClose={() => isMissionBoardOpen = false}
/>

<!-- Recent Activity Modal -->
<RecentActivityModal
  activities={dashboardStore.activities}
  isOpen={isActivityModalOpen}
  onClose={() => isActivityModalOpen = false}
/>

<!-- Live Events Modal -->
<LiveEventsModal
  events={dashboardStore.events}
  isOpen={isLiveEventsOpen}
  onClose={() => isLiveEventsOpen = false}
/>

<!-- Progress Panel -->
<ProgressPanel
  isOpen={isProgressPanelOpen}
  onClose={() => isProgressPanelOpen = false}
  {subjects}
  {userProfile}
  {diagnosticLevels}
  onSubjectClick={() => {}}
/>
