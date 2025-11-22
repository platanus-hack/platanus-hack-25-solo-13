<script>
  /**
   * StudentHome.svelte
   * Implements Svelte 5 Runes syntax ($state, $props, etc.)
   * Uses native DOM attributes (onclick, onkeydown) instead of on: directives.
   */
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // --- State ---
  let activeTab = $state('daily');

  // --- Mock Data ---
  const student = {
    name: 'Sofía Morales',
    grade: '2° Medio',
    level: 12,
    xp: 75, // %
    streak: 5
  };

  const resources = [
    { name: 'Mastery Tokens', value: 420, icon: '🪙', color: 'text-amber-400' },
    { name: 'Skill Points', value: 15, icon: '⚡', color: 'text-cyan-400' },
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

  const collections = [
    { id: 1, title: 'Algebra Concepts', count: '18/40', color: 'from-blue-500 to-cyan-500' },
    { id: 2, title: 'Mechanics Badges', count: '5/12', color: 'from-orange-500 to-red-500' },
    { id: 3, title: 'Literature Trophies', count: '3/8', color: 'from-purple-500 to-pink-500' }
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

  // --- Animations ---
  onMount(() => {
    const tl = gsap.timeline({ defaults: { ease: 'power2.out' } });

    tl.from('.hero-card', { y: 30, opacity: 0, duration: 0.8 })
      .from('.shortcut-tile', { y: 20, opacity: 0, duration: 0.5, stagger: 0.05 }, '-=0.5')
      .from('.mission-card', { x: -20, opacity: 0, duration: 0.4, stagger: 0.1 }, '-=0.5')
      .from('.collection-card', { scale: 0.8, opacity: 0, duration: 0.4, stagger: 0.1 }, '-=0.5');
  });

  function handleKeyDown(event) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      // Handle click event here
    }
  }
</script>

<div class="min-h-screen bg-slate-950 text-slate-200 font-sans selection:bg-indigo-500/30 pb-20 md:pb-10">

  <!-- Top Player Bar -->
  <header class="sticky top-0 z-50 border-b border-white/5 bg-slate-950/80 backdrop-blur-md">
    <div class="px-6 py-3 flex flex-col md:flex-row md:items-center justify-between gap-4">

      <!-- Profile -->
      <div class="flex items-center gap-3">
        <div class="h-10 w-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-sm font-bold shadow-lg shadow-indigo-500/20">
          SM
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h1 class="font-semibold text-white text-sm md:text-base">{student.name}</h1>
            <span class="px-1.5 py-0.5 rounded bg-slate-800 border border-slate-700 text-xs text-slate-400">{student.grade}</span>
          </div>
          <div class="text-xs text-slate-400">Season: PAES Prep</div>
        </div>
      </div>

      <!-- XP & Streak -->
      <div class="flex-1 max-w-md px-2 hidden md:block">
        <div class="flex justify-between text-xs mb-1">
          <span class="text-indigo-300 font-medium">Level {student.level}</span>
          <span class="text-slate-500">{student.xp}% XP</span>
        </div>
        <div class="h-2 w-full bg-slate-900 rounded-full overflow-hidden">
          <div class="h-full bg-gradient-to-r from-indigo-500 via-purple-500 to-cyan-500 w-[75%] shadow-[0_0_10px_rgba(99,102,241,0.5)]"></div>
        </div>
      </div>

      <!-- Resources -->
      <div class="flex items-center gap-3 overflow-x-auto pb-1 md:pb-0 no-scrollbar">
        <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-slate-900/50 border border-slate-800 whitespace-nowrap">
          <span class="text-sm">🔥</span>
          <span class="text-xs font-bold text-orange-400">{student.streak} Days</span>
        </div>
        {#each resources as res}
          <button class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-slate-900/50 border border-slate-800 whitespace-nowrap hover:bg-slate-800/50 transition-colors cursor-help" title={res.name} aria-label={res.name}>
            <span class="text-sm">{res.icon}</span>
            <span class={`text-xs font-bold ${res.color}`}>{res.value}</span>
          </button>
        {/each}
      </div>
    </div>
  </header>

  <!-- Main Layout -->
  <main class="px-6 py-6 md:py-8 grid grid-cols-1 lg:grid-cols-12 gap-8">

    <!-- LEFT COLUMN (Major Interactions) -->
    <div class="lg:col-span-8 space-y-8">

      <!-- Hero Quest CTA -->
      <section class="hero-card relative overflow-hidden rounded-3xl bg-gradient-to-br from-indigo-900/80 via-slate-900 to-slate-900 border border-indigo-500/30 shadow-2xl shadow-indigo-900/20 group">
        <!-- Abstract Background Elements -->
        <div class="absolute top-0 right-0 w-64 h-64 bg-indigo-600/20 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2"></div>
        <div class="absolute bottom-0 left-0 w-48 h-48 bg-cyan-600/10 rounded-full blur-2xl translate-y-1/3 -translate-x-1/4"></div>

        <div class="relative z-10 p-6 md:p-8">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-indigo-500/20 border border-indigo-400/30 text-indigo-300 text-xs font-semibold uppercase tracking-wider mb-4">
            <span class="animate-pulse">●</span> Current Quest
          </div>
          <h2 class="text-3xl md:text-4xl font-bold text-white mb-2 tracking-tight group-hover:text-indigo-100 transition-colors">Continue today's quest</h2>
          <p class="text-slate-400 text-lg mb-6 max-w-lg">Pick up where you left off in <span class="text-cyan-300">Math · Linear Functions</span>. Complete this to unlock the Algebra Badge.</p>

          <div class="flex items-center gap-4">
            <button class="px-6 py-3 rounded-xl bg-white text-indigo-950 font-bold hover:bg-indigo-50 hover:scale-105 active:scale-95 transition-all shadow-lg shadow-white/10 flex items-center gap-2" aria-label="Resume Quest" onclick={() => {}}>
              <span>Resume Quest</span>
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-play"><polygon points="6 3 20 12 6 21 6 3"/></svg>
            </button>
            <button class="px-6 py-3 rounded-xl bg-slate-800/50 text-slate-300 border border-slate-700 hover:bg-slate-800 hover:text-white transition-all" aria-label="View All" onclick={() => {}}>
              View All
            </button>
          </div>
        </div>

        <!-- Progress Indicator -->
        <div class="absolute bottom-6 right-8 hidden md:block text-right">
          <div class="text-2xl font-bold text-white">2/5</div>
          <div class="text-xs text-slate-500 uppercase tracking-widest">Daily Goals</div>
        </div>
      </section>

      <!-- Shortcut Grid -->
      <section>
        <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
          <span class="text-cyan-400">❖</span> Quick Access
        </h3>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-3 xl:grid-cols-6 gap-4">
          {#each shortcuts as item}
            <button class="shortcut-tile flex flex-col items-center justify-center p-4 rounded-2xl bg-slate-900/40 border border-slate-800/60 hover:bg-slate-800 hover:border-indigo-500/30 hover:-translate-y-1 transition-all group text-center h-32" aria-label={item.label} onclick={() => {}}>
              <div class="text-3xl mb-2 group-hover:scale-110 transition-transform">{item.icon}</div>
              <div class="font-medium text-slate-200 text-sm leading-tight">{item.label}</div>
              <div class="text-[10px] text-slate-500 mt-1">{item.desc}</div>
            </button>
          {/each}
        </div>
      </section>

      <!-- Mission Board -->
      <section class="bg-slate-900/20 rounded-3xl border border-slate-800/50 p-6">
        <div class="flex flex-wrap items-center justify-between mb-6 gap-4">
          <h3 class="text-xl font-bold text-white flex items-center gap-2">
            <span>⚔️</span> Mission Board
          </h3>
          <div class="flex p-1 rounded-xl bg-slate-950 border border-slate-800">
            {#each ['daily', 'weekly', 'story', 'side'] as tab}
              <button
                class={`px-4 py-1.5 rounded-lg text-sm font-medium transition-all capitalize ${activeTab === tab ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-900/50' : 'text-slate-500 hover:text-slate-300'}`}
                onclick={() => activeTab = tab}
              >
                {tab}
              </button>
            {/each}
          </div>
        </div>

        <div class="grid md:grid-cols-2 gap-4">
          {#each missions[activeTab] as mission (mission.id)}
            <button class="mission-card group p-4 rounded-2xl bg-slate-900/50 border border-slate-800/80 hover:border-indigo-500/40 transition-colors flex items-center gap-4 w-full text-left" aria-label={mission.title} onclick={() => {}} onkeydown={handleKeyDown}>
              <div class={`h-12 w-12 rounded-xl flex items-center justify-center text-lg font-bold shrink-0 ${
                mission.subject === 'Math' ? 'bg-blue-500/10 text-blue-400' :
                mission.subject === 'Language' ? 'bg-emerald-500/10 text-emerald-400' :
                mission.subject === 'History' ? 'bg-amber-500/10 text-amber-400' :
                'bg-slate-800 text-slate-400'
              }`}>
                {mission.subject[0]}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-0.5">
                  <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">{mission.subject}</span>
                  {#if mission.state === 'done'}
                    <span class="text-[10px] px-1.5 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/20">Done</span>
                  {/if}
                </div>
                <h4 class="font-medium text-slate-200 truncate group-hover:text-indigo-300 transition-colors">{mission.title}</h4>
                <div class="flex items-center gap-3 text-xs text-slate-500 mt-1">
                  <span class="flex items-center gap-1">⏱️ {mission.time}</span>
                  <span class="flex items-center gap-1 text-indigo-400">🎁 {mission.reward}</span>
                </div>
              </div>
              <div class="h-8 w-8 rounded-full bg-slate-900 border border-slate-700 flex items-center justify-center text-slate-400 group-hover:bg-indigo-600 group-hover:text-white group-hover:border-indigo-500 transition-all">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
              </div>
            </button>
          {/each}
        </div>
      </section>
    </div>

    <!-- RIGHT COLUMN (Meta & Social) -->
    <div class="lg:col-span-4 space-y-8">

      <!-- Event Banners (Horizontal Scroll) -->
      <section>
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-white">Live Events</h3>
          <a href="#events" class="text-xs text-indigo-400 hover:text-indigo-300">View All</a>
        </div>
        <div class="flex gap-4 overflow-x-auto pb-4 no-scrollbar snap-x">
          {#each events as event}
            <button class={`snap-center shrink-0 w-64 p-5 rounded-2xl bg-gradient-to-br ${event.color} shadow-lg text-white relative overflow-hidden group cursor-pointer text-left`} onclick={() => {}} aria-label={event.title} onkeydown={handleKeyDown}>
              <div class="relative z-10">
                <div class="text-xs font-bold opacity-80 mb-1 uppercase tracking-wider">Limited Time</div>
                <h4 class="font-bold text-xl mb-1 leading-tight">{event.title}</h4>
                <p class="text-sm opacity-90 mb-4">{event.subtitle}</p>
                <div class="px-3 py-1.5 bg-white/20 hover:bg-white/30 rounded-lg text-xs font-bold backdrop-blur-sm transition-colors inline-block">
                  Join Now
                </div>
              </div>
              <div class="absolute -bottom-4 -right-4 text-8xl opacity-20 rotate-12 group-hover:rotate-0 group-hover:scale-110 transition-transform duration-500">★</div>
            </button>
          {/each}
        </div>
      </section>

      <!-- Collections (Gacha style) -->
      <section>
        <h3 class="text-lg font-semibold text-white mb-4">Collections</h3>
        <div class="space-y-3">
          {#each collections as col}
            <button class="collection-card p-3 rounded-xl bg-slate-900/60 border border-slate-800 flex items-center gap-3 hover:bg-slate-800/60 transition-colors cursor-pointer w-full text-left" aria-label={col.title} onclick={() => {}} onkeydown={handleKeyDown}>
              <div class={`h-12 w-12 rounded-lg bg-gradient-to-br ${col.color} shadow-lg flex items-center justify-center text-xl`}>
                📚
              </div>
              <div class="flex-1">
                <div class="flex justify-between items-center mb-1">
                  <h4 class="text-sm font-medium text-slate-200">{col.title}</h4>
                  <span class="text-xs font-bold text-indigo-400">{col.count}</span>
                </div>
                <div class="h-1.5 w-full bg-slate-950 rounded-full overflow-hidden">
                  <div class={`h-full bg-gradient-to-r ${col.color}`} style={`width: ${parseInt(col.count)/parseInt(col.count.split('/')[1])*100}%`}></div>
                </div>
              </div>
            </button>
          {/each}
        </div>
      </section>

      <!-- Activity Feed -->
      <section class="bg-slate-900/30 rounded-3xl p-6 border border-slate-800/30">
        <h3 class="text-sm font-bold text-slate-500 uppercase tracking-wider mb-4">Recent Activity</h3>
        <div class="space-y-6 relative before:absolute before:left-2.5 before:top-2 before:h-full before:w-px before:bg-slate-800">
          {#each activities as act}
            <div class="relative pl-8">
              <div class="absolute left-0 top-1 h-5 w-5 rounded-full bg-slate-900 border-2 border-slate-700 flex items-center justify-center text-[10px] z-10">
                {act.icon}
              </div>
              <p class="text-sm text-slate-300 leading-snug mb-1">{act.text}</p>
              <span class="text-xs text-slate-600">{act.time}</span>
            </div>
          {/each}
        </div>
      </section>

    </div>
  </main>
</div>

<style>
  /* Utility for hiding scrollbars but allowing scroll */
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }
  .no-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>
