<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import SubjectDetailModal from '$lib/components/dashboard/SubjectDetailModal.svelte';
  import SubNavBar from '$lib/components/dashboard/SubNavBar.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';
  import RecentActivityModal from '$lib/components/dashboard/RecentActivityModal.svelte';
  import MissionBoardModal from '$lib/components/dashboard/MissionBoardModal.svelte';
  import CurrentQuestModal from '$lib/components/dashboard/CurrentQuestModal.svelte';
  import LiveEventsModal from '$lib/components/dashboard/LiveEventsModal.svelte';
  import PlayerProfilePanel from '$lib/components/dashboard/PlayerProfilePanel.svelte';

  // State
  let activeTab = $state('daily');
  let userProfile = $state(null);
  let subjects = $state<Subject[]>([]);
  let selectedSubject = $state<Subject | null>(null);
  let isModalOpen = $state(false);
  let selectedDomainLevel = $state(0);
  let isProgressPanelOpen = $state(false);
  let isActivityModalOpen = $state(false);
  let isMissionBoardOpen = $state(false);
  let isCurrentQuestOpen = $state(false);
  let isLiveEventsOpen = $state(false);
  let isPlayerProfileOpen = $state(false);

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

  const missions = {
    daily: [
      { id: 1, subject: 'Math', title: 'Linear Equations Practice', time: '15 min', reward: '50 XP', state: 'start' },
      { id: 2, subject: 'Language', title: 'Reading Comprehension', time: '10 min', reward: '30 XP', state: 'done' },
      { id: 3, subject: 'History', title: 'Chilean Independence', time: '20 min', reward: '60 XP', state: 'start' }
    ],
    weekly: [
      { id: 4, subject: 'Math', title: 'Complete Unit 3: Functions', time: '2 hrs', reward: '500 XP', state: 'progress' },
      { id: 5, subject: 'Physics', title: 'Lab Report: Motion', time: '45 min', reward: '200 XP', state: 'start' }
    ],
    story: [
      { id: 6, subject: 'Campaign', title: 'Chapter 2: The Algebra Realm', time: 'Ongoing', reward: 'Badge', state: 'progress' }
    ],
    side: [
      { id: 7, subject: 'Challenge', title: 'Speed Math: Mental Calculation', time: '5 min', reward: '10 SP', state: 'start' }
    ]
  };

  // Count active missions (not done)
  const activeMissionCount = $derived(
    Object.values(missions).flat().filter(m => m.state !== 'done').length
  );

  const events = [
    { id: 1, title: 'Exam Sprint Week', subtitle: 'Boost your PAES score', color: 'from-indigo-600 to-violet-600' },
    { id: 2, title: 'Math Boss Challenge', subtitle: 'Beat the Function Dragon', color: 'from-rose-600 to-orange-600' }
  ];

  const activities = [
    { id: 1, text: "You unlocked 'Equation Explorer'", time: '2h ago', icon: '🏆' },
    { id: 2, text: 'Teacher left feedback on essay', time: '4h ago', icon: '💬' },
    { id: 3, text: 'New PAES simulator available', time: '1d ago', icon: '🆕' }
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
          subjects = materiasToSubjects(cursoResult.curso.materias);
        } else {
          console.error('Failed to load subjects:', cursoResult.error);
          subjects = [];
        }
      }
    }
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
    <div class="px-6 py-3 flex flex-col md:flex-row md:items-center justify-between gap-4">
      <!-- Profile -->
      <button
        onclick={() => isPlayerProfileOpen = true}
        class="flex items-center gap-3 hover:bg-canvas-800/40 rounded-xl p-2 -m-2 transition-colors group"
      >
        <div class="h-10 w-10 rounded-full bg-gradient-to-br from-lumera-500 to-focus-600 flex items-center justify-center text-sm font-bold group-hover:scale-110 transition-transform">
          {initials}
        </div>
        <div class="text-left">
          <div class="flex items-center gap-2">
            <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
            <span class="px-1.5 py-0.5 rounded bg-canvas-800 border border-canvas-700 text-xs text-slate-400">{student.grade}</span>
          </div>
          <div class="text-xs text-slate-400">{auth.user?.email || 'student@lumera.com'}</div>
        </div>
      </button>

      <!-- XP -->
      <div class="flex-1 max-w-md px-2 hidden md:block">
        <div class="flex justify-between text-xs mb-1">
          <span class="text-achievement-400 font-medium">Level {student.level}</span>
          <span class="text-slate-500">{student.xp}% XP</span>
        </div>
        <div class="h-2 w-full bg-canvas-900 rounded-full overflow-hidden">
          <div class="h-full bg-gradient-to-r from-lumera-500 via-focus-500 to-achievement-400" style="width: {student.xp}%"></div>
        </div>
      </div>

      <!-- Resources -->
      <div class="flex items-center gap-3">
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
    </div>
  </header>

  <!-- Sub Navigation Bar -->
  <SubNavBar
    onQuestClick={() => isCurrentQuestOpen = true}
    onProgressClick={() => isProgressPanelOpen = true}
    onActivityClick={() => isActivityModalOpen = true}
    onMissionsClick={() => isMissionBoardOpen = true}
    onEventsClick={() => isLiveEventsOpen = true}
    onLogoutClick={() => auth.logout()}
    subjectCount={subjects.length}
    activityCount={activities.length}
    missionCount={activeMissionCount}
    eventCount={events.length}
  />

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
  {subjects}
  {userProfile}
  onSubjectClick={openSubjectDetail}
/>

<!-- Recent Activity Modal -->
<RecentActivityModal
  {activities}
  isOpen={isActivityModalOpen}
  onClose={() => isActivityModalOpen = false}
/>

<!-- Mission Board Modal -->
<MissionBoardModal
  {missions}
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
  {events}
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
/>
