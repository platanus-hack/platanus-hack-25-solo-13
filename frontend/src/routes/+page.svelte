<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import SubjectDetailModal from '$lib/components/dashboard/SubjectDetailModal.svelte';
  import SubNavBar from '$lib/components/dashboard/SubNavBar.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';

  // State
  let activeTab = $state('daily');
  let userProfile = $state(null);
  let subjects = $state<Subject[]>([]);
  let selectedSubject = $state<Subject | null>(null);
  let isModalOpen = $state(false);
  let selectedDomainLevel = $state(0);
  let isProgressPanelOpen = $state(false);

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
    { name: 'Mastery Tokens', value: 420, icon: '🪙', color: 'text-amber-400' },
    { name: 'Skill Points', value: 15, icon: '⚡', color: 'text-focus-400' },
    { name: 'PAES Ready', value: '63%', icon: '📈', color: 'text-emerald-400' }
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

  const shortcuts = [
    { id: 1, label: 'My Curriculum', desc: 'Map of units', icon: '🗺️' },
    { id: 2, label: 'PAES Sim', desc: 'Full practice exams', icon: '🎓' },
    { id: 3, label: 'Adaptive Gym', desc: 'Train weak spots', icon: '🏋️' },
    { id: 4, label: 'Classroom', desc: 'Teacher updates', icon: '👨‍🏫' },
    { id: 5, label: 'Friends', desc: 'Study groups', icon: '👥' },
    { id: 6, label: 'Settings', desc: 'Profile & App', icon: '⚙️' }
  ];

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

  // Animations
  onMount(() => {
    loadUserProfile();

    const tl = gsap.timeline({ defaults: { ease: 'power2.out' } });
    tl.from('.hero-card', { y: 30, opacity: 0, duration: 0.8 })
      .from('.shortcut-tile', { y: 20, opacity: 0, duration: 0.5, stagger: 0.05 }, '-=0.5')
      .from('.mission-card', { x: -20, opacity: 0, duration: 0.4, stagger: 0.1 }, '-=0.5')
      .from('.subject-card', { scale: 0.8, opacity: 0, duration: 0.4, stagger: 0.1 }, '-=0.5');
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
      <div class="flex items-center gap-3">
        <div class="h-10 w-10 rounded-full bg-gradient-to-br from-lumera-500 to-focus-600 flex items-center justify-center text-sm font-bold">
          {initials}
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
            <span class="px-1.5 py-0.5 rounded bg-canvas-800 border border-canvas-700 text-xs text-slate-400">{student.grade}</span>
          </div>
          <div class="text-xs text-slate-400">{auth.user?.email || 'student@lumera.com'}</div>
        </div>
      </div>

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
        <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-canvas-900/50 border border-slate-800">
          <span class="text-sm">🔥</span>
          <span class="text-xs font-bold text-orange-400">{student.streak} Days</span>
        </div>
        {#each resources as res}
          <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-canvas-900/50 border border-slate-800">
            <span class="text-sm">{res.icon}</span>
            <span class="text-xs font-bold {res.color}">{res.value}</span>
          </div>
        {/each}
        <button
          onclick={() => auth.logout()}
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-red-900/50 border border-red-800 hover:bg-red-800/50 transition-colors"
          title="Logout"
        >
          <span class="text-sm">🚪</span>
          <span class="text-xs font-bold text-red-400">Logout</span>
        </button>
      </div>
    </div>
  </header>

  <!-- Sub Navigation Bar -->
  <SubNavBar
    onProgressClick={() => isProgressPanelOpen = true}
    subjectCount={subjects.length}
  />

  <!-- Main -->
  <main class="px-6 py-8 grid grid-cols-1 lg:grid-cols-12 gap-8">
    <!-- LEFT -->
    <div class="lg:col-span-8 space-y-8">
      <!-- Hero Quest -->
      <section class="hero-card relative overflow-hidden rounded-2xl bg-gradient-to-br from-lumera-900/90 via-canvas-900 to-canvas-900 border border-lumera-600/40 p-10">
        <div class="relative z-10">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-lumera-500/20 border border-lumera-400/30 text-lumera-300 text-xs font-semibold mb-4">
            <span class="animate-pulse">●</span> Current Quest
          </div>
          <h2 class="text-5xl font-display font-bold text-white mb-2 tracking-tight">Continue today's quest</h2>
          <p class="text-slate-400 text-lg mb-6">Pick up where you left off in <span class="text-focus-300">Math · Linear Functions</span>.</p>
          <button class="px-6 py-3 rounded-lg bg-white text-lumera-950 font-bold hover:scale-105 transition-all">Resume Quest</button>
        </div>
      </section>

      <!-- Shortcuts -->
      <section>
        <h3 class="text-lg font-semibold text-white mb-4">Quick Access</h3>
        <div class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-4">
          {#each shortcuts as item}
            <button class="shortcut-tile p-4 rounded-2xl bg-canvas-900/40 border border-slate-800/60 hover:bg-canvas-800 transition-all text-center">
              <div class="text-3xl mb-2">{item.icon}</div>
              <div class="font-medium text-slate-200 text-sm">{item.label}</div>
              <div class="text-[10px] text-slate-500 mt-1">{item.desc}</div>
            </button>
          {/each}
        </div>
      </section>

      <!-- Missions -->
      <section class="bg-canvas-900/20 rounded-2xl border border-canvas-800/50 p-6">
        <div class="flex flex-wrap items-center justify-between mb-6 gap-4">
          <h3 class="text-xl font-display font-bold text-white">⚔️ Mission Board</h3>
          <div class="flex p-1 rounded-xl bg-canvas-950 border border-canvas-800">
            {#each ['daily', 'weekly', 'story', 'side'] as tab}
              <button
                class="px-4 py-1.5 rounded-lg text-sm font-medium transition-all capitalize {activeTab === tab ? 'bg-lumera-600 text-white' : 'text-slate-500 hover:text-slate-300'}"
                onclick={() => activeTab = tab}
              >
                {tab}
              </button>
            {/each}
          </div>
        </div>

        <div class="grid md:grid-cols-2 gap-4">
          {#each missions[activeTab] as mission}
            <button class="mission-card p-4 rounded-xl bg-canvas-900/50 border border-canvas-800/80 hover:border-lumera-500/40 transition-colors flex items-center gap-4 w-full text-left">
              <div class="h-12 w-12 rounded-lg {getSubjectColor(mission.subject)} flex items-center justify-center text-lg font-bold">
                {mission.subject[0]}
              </div>
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-0.5">
                  <span class="text-xs font-bold text-slate-500 uppercase">{mission.subject}</span>
                  {#if mission.state === 'done'}
                    <span class="text-[10px] px-1.5 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/20">Done</span>
                  {/if}
                </div>
                <h4 class="font-medium text-slate-200">{mission.title}</h4>
                <div class="flex items-center gap-3 text-xs text-slate-500 mt-1">
                  <span>⏱️ {mission.time}</span>
                  <span class="text-achievement-400">🎁 {mission.reward}</span>
                </div>
              </div>
            </button>
          {/each}
        </div>
      </section>
    </div>

    <!-- RIGHT -->
    <div class="lg:col-span-4 space-y-8">
      <!-- Events -->
      <section>
        <h3 class="text-lg font-semibold text-white mb-4">Live Events</h3>
        <div class="space-y-4">
          {#each events as event}
            <button class="w-full text-left p-5 rounded-2xl bg-gradient-to-br {event.color} text-white relative overflow-hidden">
              <div class="relative z-10">
                <div class="text-xs font-bold opacity-80 mb-1">LIMITED TIME</div>
                <h4 class="font-bold text-xl mb-1">{event.title}</h4>
                <p class="text-sm opacity-90">{event.subtitle}</p>
              </div>
            </button>
          {/each}
        </div>
      </section>

      <!-- Activity -->
      <section class="bg-canvas-900/30 rounded-2xl p-6 border border-slate-800/30">
        <h3 class="text-sm font-bold text-slate-500 uppercase tracking-wider mb-4">Recent Activity</h3>
        <div class="space-y-4">
          {#each activities as act}
            <div class="flex gap-3">
              <div class="text-lg">{act.icon}</div>
              <div>
                <p class="text-sm text-slate-300">{act.text}</p>
                <span class="text-xs text-slate-600">{act.time}</span>
              </div>
            </div>
          {/each}
        </div>
      </section>
    </div>
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
