<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import { getEquipment, type CustomizationItem } from '$lib/api/customization';
  import SubjectDetailModal from '$lib/components/dashboard/SubjectDetailModal.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';
  import RecentActivityModal from '$lib/components/dashboard/RecentActivityModal.svelte';
  import MissionBoardModal from '$lib/components/dashboard/MissionBoardModal.svelte';
  import CurrentQuestModal from '$lib/components/dashboard/CurrentQuestModal.svelte';
  import LiveEventsModal from '$lib/components/dashboard/LiveEventsModal.svelte';
  import PlayerProfilePanel from '$lib/components/dashboard/PlayerProfilePanel.svelte';
  import AvatarDisplay from '$lib/components/common/AvatarDisplay.svelte';

  // State
  let activeTab = $state('daily');
  let userProfile = $state(null);
  let selectedSubject = $state<Subject | null>(null);
  let isModalOpen = $state(false);
  let selectedDomainLevel = $state(0);
  let isProgressPanelOpen = $state(false);
  let isActivityModalOpen = $state(false);
  let isMissionBoardOpen = $state(false);
  let isCurrentQuestOpen = $state(false);
  let isLiveEventsOpen = $state(false);
  let isPlayerProfileOpen = $state(false);
  let currentAvatar = $state<CustomizationItem | null>(null);

  // Get authenticated user
  const student = {
    name: auth.user?.name || 'Student',
    grade: '2° Medio',
    level: 12,
    xp: 75,
    streak: 5
  };

  // Get initials for avatar
  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  const resources = [
    {
      name: 'Mastery Tokens',
      value: 420,
      iconType: 'currency',
      color: 'text-white'
    },
    {
      name: 'Skill Points',
      value: 15,
      iconType: 'bolt',
      color: 'text-white'
    },
    {
      name: 'PAES Ready',
      value: '63%',
      iconType: 'chart',
      color: 'text-white'
    }
  ];

  // Load user profile and subjects
  async function loadUserProfile() {
    if (auth.user?.id) {
      const result = await getProfile(auth.user.id);
      if (result.success && result.profile) {
        userProfile = result.profile;
        // Get subjects for user's course from backend
        const courseName = result.profile.curso_actual || '1ro Medio';
        const cursoResult = await findCursoByName(courseName);
        if (cursoResult.success && cursoResult.curso?.materias) {
          dashboardStore.updateSubjects(materiasToSubjects(cursoResult.curso.materias));
        } else {
          console.error('Failed to load subjects:', cursoResult.error);
          dashboardStore.updateSubjects([]);
        }
      }
    }
  }

  // Load equipped avatar
  async function loadEquippedAvatar() {
    try {
      console.log('Loading equipped avatar...');
      const equipment = await getEquipment();
      console.log('Equipment response:', equipment);
      currentAvatar = equipment.equipped_avatar || null;
      console.log('Current avatar set to:', currentAvatar);
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
      return 0; // Not evaluated
    }

    // Map subject code to conocimiento_previo key
    // MAT -> matematicas, LYL -> lectura
    const subjectMap: Record<string, string> = {
      'mat': 'matematicas',
      'lyl': 'lectura',
      'lenguaje': 'lectura',
      'lengua': 'lectura'
    };

    const key = subjectMap[subjectId.toLowerCase()] || subjectId.toLowerCase();
    return userProfile.profile_data.conocimiento_previo[key]?.nivel || 0;
  }

  // Open subject detail modal
  function openSubjectDetail(subject: Subject) {
    selectedSubject = subject;
    selectedDomainLevel = getDomainLevel(subject.id);
    isModalOpen = true;
  }

  // Close modal
  function closeModal() {
    isModalOpen = false;
    selectedSubject = null;
  }

  // Load profile on mount
  onMount(() => {
    loadUserProfile();
    loadEquippedAvatar();
  });

  function getSubjectColor(subject) {
    const colors = {
      'Math': 'bg-blue-500/10 text-blue-400',
      'Language': 'bg-emerald-500/10 text-emerald-400',
      'History': 'bg-amber-500/10 text-amber-400',
      'Physics': 'bg-purple-500/10 text-purple-400',
      'Campaign': 'bg-pink-500/10 text-pink-400',
      'Challenge': 'bg-orange-500/10 text-orange-400'
    };
    return colors[subject] || 'bg-canvas-800 text-slate-400';
  }
</script>

<svelte:head>
  <title>Lumera App - Student Dashboard</title>
</svelte:head>

<div class="min-h-screen bg-canvas-950 text-slate-200 font-sans pb-10">
  <!-- Top Player Bar -->
  <header class="sticky top-0 z-50 border-b border-white/5 bg-canvas-950/90 backdrop-blur-md">
    <div class="px-6 py-3 flex items-center justify-between gap-4">
      <!-- Profile -->
      <button
        onclick={() => isPlayerProfileOpen = true}
        class="flex items-center gap-3 hover:bg-canvas-800/40 rounded-xl p-2 -m-2 transition-colors group"
      >
        <div class="group-hover:scale-110 transition-transform">
          <AvatarDisplay
            {currentAvatar}
            {initials}
            size="small"
          />
        </div>
        <div class="text-left">
          <div class="flex items-center gap-2">
            <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
            <span class="px-1.5 py-0.5 rounded bg-canvas-800 border border-canvas-700 text-xs text-slate-400">{student.grade}</span>
          </div>
          <div class="text-xs text-slate-400">{auth.user?.email || 'student@lumera.com'}</div>
        </div>
      </button>

      <!-- Center - Badges -->
      <div class="flex items-center gap-2">
        <!-- Streak Badge -->
        <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-canvas-900/50 border border-slate-800">
          <svg class="w-4 h-4 text-orange-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2.25c-2.429 0-4.817.178-7.152.521C2.87 3.061 1.5 4.795 1.5 6.741v6.018c0 1.946 1.37 3.68 3.348 3.97.877.129 1.761.234 2.652.316V21a.75.75 0 001.28.53l4.184-4.183a.39.39 0 01.266-.112c2.006-.05 3.982-.22 5.922-.506 1.978-.29 3.348-2.023 3.348-3.97V6.741c0-1.947-1.37-3.68-3.348-3.97A49.145 49.145 0 0012 2.25zM8.25 8.625a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25zm2.625 1.125a1.125 1.125 0 112.25 0 1.125 1.125 0 01-2.25 0zm4.875-1.125a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25z" />
          </svg>
          <span class="text-xs font-bold text-white">{student.streak}</span>
        </div>

        {#each resources as res}
          <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-canvas-900/50 border border-slate-800">
            {#if res.iconType === 'currency'}
              <svg class="w-4 h-4 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            {:else if res.iconType === 'bolt'}
              <svg class="w-4 h-4 text-focus-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M11.25 5.337c0-.355-.186-.676-.401-.959a1.647 1.647 0 01-.349-1.003c0-1.036 1.007-1.875 2.25-1.875S15 2.34 15 3.375c0 .369-.128.713-.349 1.003-.215.283-.401.604-.401.959 0 .332.278.598.61.578 1.91-.114 3.79-.342 5.632-.676a.75.75 0 01.878.645 49.17 49.17 0 01.376 5.452.657.657 0 01-.66.664c-.354 0-.675-.186-.958-.401a1.647 1.647 0 00-1.003-.349c-1.035 0-1.875 1.007-1.875 2.25s.84 2.25 1.875 2.25c.369 0 .713-.128 1.003-.349.283-.215.604-.401.959-.401.31 0 .557.262.534.571a48.774 48.774 0 01-.595 4.845.75.75 0 01-.61.61c-1.82.317-3.673.533-5.555.642a.58.58 0 01-.611-.581c0-.355.186-.676.401-.959.221-.29.349-.634.349-1.003 0-1.035-1.007-1.875-2.25-1.875s-2.25.84-2.25 1.875c0 .369.128.713.349 1.003.215.283.401.604.401.959a.641.641 0 01-.658.643 49.118 49.118 0 01-4.708-.36.75.75 0 01-.645-.878c.293-1.614.504-3.257.629-4.924A.53.53 0 005.337 15c-.355 0-.676.186-.959.401-.29.221-.634.349-1.003.349-1.036 0-1.875-1.007-1.875-2.25s.84-2.25 1.875-2.25c.369 0 .713.128 1.003.349.283.215.604.401.959.401a.656.656 0 00.659-.663 47.703 47.703 0 00-.31-4.82.75.75 0 01.83-.832c1.343.155 2.703.254 4.077.294a.64.64 0 00.657-.642z" />
              </svg>
            {:else if res.iconType === 'chart'}
              <svg class="w-4 h-4 text-emerald-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            {/if}
            <span class="text-xs font-bold {res.color}">{res.value}</span>
          </div>
        {/each}
      </div>

      <!-- Right - Navigation Icons -->
      <div class="flex items-center gap-2">
        <!-- Quest Button -->
        <button
          onclick={() => isCurrentQuestOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Current Quest"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
          </svg>
        </button>

        <!-- Missions Button -->
        <button
          onclick={() => isMissionBoardOpen = true}
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Mission Board"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Actividad Reciente"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Live Events"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
          class="relative p-2 rounded-lg hover:bg-canvas-800/60 transition-all duration-200 group"
          title="Tu Progreso"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-slate-200 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          {#if dashboardStore.subjects.length > 0}
            <span class="absolute -top-1 -right-1 text-xs font-bold bg-lumera-600 text-white px-1.5 py-0.5 rounded-full min-w-[20px] text-center">
              {dashboardStore.subjects.length}
            </span>
          {/if}
        </button>

        <!-- Logout Button -->
        <button
          onclick={() => auth.logout()}
          class="relative p-2 rounded-lg hover:bg-red-900/60 transition-all duration-200 group"
          title="Logout"
        >
          <svg class="w-6 h-6 text-slate-400 group-hover:text-red-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>
    </div>
  </header>

  <!-- Main -->
  <main class="px-6 py-8 max-w-7xl mx-auto">
    <!-- Dashboard content area - currently empty, all content in modals -->
  </main>
</div>

<!-- Subject Detail Modal -->
<SubjectDetailModal
  subject={selectedSubject}
  domainLevel={selectedDomainLevel}
  isOpen={isModalOpen}
  onClose={closeModal}
/>

<!-- Progress Panel -->
<ProgressPanel
  isOpen={isProgressPanelOpen}
  onClose={() => isProgressPanelOpen = false}
  subjects={dashboardStore.subjects}
  {userProfile}
  onSubjectClick={openSubjectDetail}
/>

<!-- Recent Activity Modal -->
<RecentActivityModal
  activities={dashboardStore.activities}
  isOpen={isActivityModalOpen}
  onClose={() => isActivityModalOpen = false}
/>

<!-- Mission Board Modal -->
<MissionBoardModal
  missions={dashboardStore.missions}
  {activeTab}
  onTabChange={(tab) => activeTab = tab}
  isOpen={isMissionBoardOpen}
  onClose={() => isMissionBoardOpen = false}
/>

<!-- Current Quest Modal -->
<CurrentQuestModal
  isOpen={isCurrentQuestOpen}
  onClose={() => isCurrentQuestOpen = false}
  onResumeQuest={() => console.log('Resume quest clicked')}
/>

<!-- Live Events Modal -->
<LiveEventsModal
  events={dashboardStore.events}
  isOpen={isLiveEventsOpen}
  onClose={() => isLiveEventsOpen = false}
/>

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
