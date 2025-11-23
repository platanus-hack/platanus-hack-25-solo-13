<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';
  import { dashboardStore } from '$lib/stores/dashboard.svelte';
  import { customizationStore } from '$lib/stores/customization.svelte';
  import { getProfile } from '$lib/api/profiles';
  import { findCursoByName } from '$lib/api/courses';
  import { materiasToSubjects, type Subject } from '$lib/constants/subjects';
  import type { CustomizationItem } from '$lib/api/customization';
  import SubjectDetailModal from '$lib/components/dashboard/SubjectDetailModal.svelte';
  import ProgressPanel from '$lib/components/dashboard/ProgressPanel.svelte';
  import RecentActivityModal from '$lib/components/dashboard/RecentActivityModal.svelte';
  import MissionBoardModal from '$lib/components/dashboard/MissionBoardModal.svelte';
  import CurrentQuestModal from '$lib/components/dashboard/CurrentQuestModal.svelte';
  import LiveEventsModal from '$lib/components/dashboard/LiveEventsModal.svelte';
  import PlayerProfilePanel from '$lib/components/dashboard/PlayerProfilePanel.svelte';
  import AppHeader from '$lib/components/common/AppHeader.svelte';

  // State
  let activeTab = $state('daily');
  let userProfile = $state(null);
  let selectedSubject = $state<Subject | null>(null);
  let isModalOpen = $state(false);
  let selectedDomainLevel = $state(0);
  let isProgressPanelOpen = $state(false);
  let isPlayerProfileOpen = $state(false);
  // Keep these for modals even though they're not in header anymore
  let isActivityModalOpen = $state(false);
  let isMissionBoardOpen = $state(false);
  let isCurrentQuestOpen = $state(false);
  let isLiveEventsOpen = $state(false);
  let diagnosticLevels = $state<Record<number, number>>({});
  let dailyRecommendation = $state<any>(null);
  let isLoadingRecommendation = $state(false);

  // Mission Control Center State
  let particles = $state<Array<{ x: number; y: number; vx: number; vy: number; size: number; color: string }>>([]);
  let mouseX = $state(0);
  let mouseY = $state(0);
  let currentTimeOfDay = $state('mañana');
  let currentGreeting = $state('Buenos días');
  let mainContentRef: HTMLElement | null = null;

  // Get authenticated user
  const student = $derived({
    name: auth.user?.name || 'Student',
    grade: '2° Medio',
    level: auth.gamificationStats?.level || 1,
    xp: auth.gamificationStats?.xp || 0,
    streak: auth.gamificationStats?.current_streak || 0,
    coins: auth.gamificationStats?.coins || 100
  });

  // Get initials for avatar
  const initials = auth.user?.name ? auth.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) : 'ST';

  // Calculate progress to next level
  const levelProgress = $derived(() => {
    if (!auth.gamificationStats) return 0;
    const currentLevelStartXP = auth.gamificationStats.xp - auth.gamificationStats.xp_progress;
    const totalXPForLevel = auth.gamificationStats.xp_for_next_level - currentLevelStartXP;
    return Math.round((auth.gamificationStats.xp_progress / totalXPForLevel) * 100);
  });

  const resources = $derived([
    {
      name: 'Mastery Tokens',
      value: student.coins,
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
  ]);

  // Get time-based greeting
  function getTimeOfDay(): { greeting: string; period: string } {
    const hour = new Date().getHours();
    if (hour >= 5 && hour < 12) {
      return { greeting: 'Buenos días', period: 'mañana' };
    } else if (hour >= 12 && hour < 20) {
      return { greeting: 'Buenas tardes', period: 'tarde' };
    } else {
      return { greeting: 'Buenas noches', period: 'noche' };
    }
  }

  // Initialize particles
  function initParticles() {
    const particleColors = ['#3b82f6', '#14b8a6', '#f59e0b', '#8b5cf6'];
    particles = Array.from({ length: 30 }, () => ({
      x: Math.random() * 100,
      y: Math.random() * 100,
      vx: (Math.random() - 0.5) * 0.5,
      vy: (Math.random() - 0.5) * 0.5,
      size: Math.random() * 4 + 2,
      color: particleColors[Math.floor(Math.random() * particleColors.length)]
    }));
  }

  // Update particles position
  function animateParticles() {
    particles = particles.map(p => {
      let newX = p.x + p.vx;
      let newY = p.y + p.vy;

      // Bounce off edges
      if (newX < 0 || newX > 100) p.vx *= -1;
      if (newY < 0 || newY > 100) p.vy *= -1;

      // Mouse interaction - particles move away from mouse
      const dx = newX - mouseX;
      const dy = newY - mouseY;
      const dist = Math.sqrt(dx * dx + dy * dy);

      if (dist < 15) {
        const force = (15 - dist) / 15;
        newX += (dx / dist) * force * 2;
        newY += (dy / dist) * force * 2;
      }

      return {
        ...p,
        x: Math.max(0, Math.min(100, newX)),
        y: Math.max(0, Math.min(100, newY))
      };
    });
  }

  // Handle mouse move for particle interaction
  function handleMouseMove(e: MouseEvent) {
    if (mainContentRef) {
      const rect = mainContentRef.getBoundingClientRect();
      mouseX = ((e.clientX - rect.left) / rect.width) * 100;
      mouseY = ((e.clientY - rect.top) / rect.height) * 100;
    }
  }

  // Load user profile and subjects
  async function loadUserProfile() {
    if (auth.user?.id) {
      // Load profile into auth store for header
      await auth.loadUserProfile();

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

  // Load completed diagnostics
  async function loadDiagnosticLevels() {
    try {
      const token = auth.token;
      const response = await fetch('/api/diagnostic-sessions?estado=completado', {
        headers: {
          'Authorization': token ? `Bearer ${token}` : ''
        }
      });

      if (!response.ok) return;

      const sessions = await response.json();
      const levels: Record<number, number> = {};

      // For each completed session, get the average bloom level
      for (const session of sessions) {
        if (session.estrategia?.average_bloom_level) {
          levels[session.materia_id] = Math.round(session.estrategia.average_bloom_level);
        }
      }

      diagnosticLevels = levels;
      console.log('Diagnostic levels loaded:', diagnosticLevels);
    } catch (err) {
      console.error('Error loading diagnostic levels:', err);
    }
  }

  // Load daily recommendation
  async function loadDailyRecommendation() {
    isLoadingRecommendation = true;
    try {
      const token = auth.token;
      const response = await fetch('/api/recommendations/daily', {
        headers: {
          'Authorization': token ? `Bearer ${token}` : ''
        }
      });

      if (!response.ok) {
        console.error('Failed to load recommendation:', response.statusText);
        return;
      }

      dailyRecommendation = await response.json();
      console.log('Daily recommendation loaded:', dailyRecommendation);
    } catch (err) {
      console.error('Error loading daily recommendation:', err);
    } finally {
      isLoadingRecommendation = false;
    }
  }

  // Handle avatar change from profile panel
  function handleAvatarChanged(avatar: CustomizationItem) {
    customizationStore.setAvatar(avatar);
  }

  // Get domain level for a subject
  function getDomainLevel(materiaId?: number): number {
    if (!materiaId) return 0;
    return diagnosticLevels[materiaId] || 0;
  }

  // Open subject detail modal
  function openSubjectDetail(subject: Subject) {
    selectedSubject = subject;
    selectedDomainLevel = getDomainLevel(subject.materiaId);
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
    customizationStore.loadAvatar();
    loadDiagnosticLevels();
    loadDailyRecommendation();
    auth.loadGamificationStats();

    // Initialize Mission Control Center
    const timeData = getTimeOfDay();
    currentGreeting = timeData.greeting;
    currentTimeOfDay = timeData.period;

    initParticles();

    // Animate particles every 50ms
    const particleInterval = setInterval(animateParticles, 50);

    // Animate elements on mount
    if (mainContentRef) {
      gsap.from('.greeting-text', {
        opacity: 0,
        y: -20,
        duration: 0.8,
        ease: 'power2.out'
      });

      gsap.from('.daily-focus-card', {
        opacity: 0,
        scale: 0.9,
        duration: 0.6,
        delay: 0.2,
        ease: 'back.out(1.7)'
      });

      gsap.from('.floating-stat', {
        opacity: 0,
        y: 20,
        duration: 0.5,
        stagger: 0.1,
        delay: 0.4,
        ease: 'power2.out'
      });
    }

    // Cleanup
    return () => {
      clearInterval(particleInterval);
    };
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

  // Start practice from recommendation
  async function startRecommendedPractice() {
    if (!dailyRecommendation) return;

    const materiaId = dailyRecommendation.oa.materia_id;
    const oaId = dailyRecommendation.oa.id;

    // Navigate to practice page
    window.location.href = `/materias/${materiaId}/practica/${oaId}`;
  }
</script>

<svelte:head>
  <title>Lumera App - Student Dashboard</title>
</svelte:head>

<div class="min-h-screen bg-canvas-950 text-slate-800 font-sans pb-10">
  <!-- App Header -->
  <AppHeader
    currentAvatar={customizationStore.currentAvatar}
    isHomePage={true}
    onProfileClick={() => isPlayerProfileOpen = true}
    onProgressClick={async () => {
      await loadDiagnosticLevels();
      isProgressPanelOpen = true;
    }}
  />

  <!-- Main -->
  <main
    bind:this={mainContentRef}
    onmousemove={handleMouseMove}
    class="relative px-6 py-8 max-w-7xl mx-auto min-h-[calc(100vh-120px)] overflow-hidden"
  >
    <!-- Particle System Background -->
    <div class="absolute inset-0 pointer-events-none">
      {#each particles as particle}
        <div
          class="absolute rounded-full opacity-40 blur-sm"
          style="
            left: {particle.x}%;
            top: {particle.y}%;
            width: {particle.size}px;
            height: {particle.size}px;
            background-color: {particle.color};
            transition: all 0.05s linear;
          "
        ></div>
      {/each}
    </div>

    <!-- Mission Control Center Content -->
    <div class="relative z-10 flex flex-col items-center justify-center min-h-full pt-8 pb-12">

      <!-- Dynamic Greeting -->
      <div class="greeting-text text-center mb-8">
        <div class="inline-block px-8 py-4 rounded-2xl bg-gradient-to-r {
          currentTimeOfDay === 'mañana'
            ? 'from-orange-500/20 via-yellow-500/20 to-blue-500/20'
            : currentTimeOfDay === 'tarde'
            ? 'from-blue-500/20 via-purple-500/20 to-pink-500/20'
            : 'from-purple-500/20 via-blue-900/20 to-slate-800/20'
        } border-2 border-white/10 backdrop-blur-sm">
          <h2 class="text-3xl md:text-4xl font-bold text-white">
            {currentGreeting}, {student.name}
          </h2>
        </div>
      </div>

      <!-- Cards Grid: Mission (2/3) + Subjects (1/3) -->
      <div class="daily-focus-card w-full max-w-5xl mb-8">
        <div class="grid grid-cols-3 gap-6">
          <!-- Daily Focus Card (2/3 width) -->
          <div class="col-span-2">
            <div class="relative group h-full">
              <!-- Animated glow effect -->
              <div class="absolute -inset-1 bg-gradient-to-r from-lumera-500 via-focus-500 to-purple-500 rounded-3xl blur-lg opacity-30 group-hover:opacity-50 transition-opacity duration-300"></div>

              <div class="relative bg-canvas-800/90 backdrop-blur-md rounded-3xl p-8 border-2 border-white/10 h-full">
                <div class="flex items-center gap-3 mb-4">
                  <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-lumera-500 to-focus-500 flex items-center justify-center">
                    <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                  </div>
                  <h3 class="text-2xl font-bold text-white">Tu Misión de Hoy</h3>
                </div>

                {#if isLoadingRecommendation}
                  <div class="mb-6 animate-pulse">
                    <div class="h-6 bg-canvas-700 rounded w-3/4 mb-2"></div>
                    <div class="h-4 bg-canvas-700 rounded w-1/2"></div>
                  </div>
                {:else if dailyRecommendation}
                  <div class="mb-6">
                    <p class="text-xl text-slate-200 mb-2">
                      Completa {dailyRecommendation.numero_preguntas} preguntas de {dailyRecommendation.oa.titulo}
                    </p>
                    <p class="text-sm text-slate-400">
                      {dailyRecommendation.reason}
                    </p>
                    <div class="mt-2 flex items-center gap-2">
                      <span class="px-2 py-1 text-xs font-medium bg-lumera-500/20 text-lumera-300 rounded-lg border border-lumera-500/30">
                        {dailyRecommendation.bloom_level.nombre}
                      </span>
                      <span class="px-2 py-1 text-xs font-medium bg-canvas-700 text-slate-300 rounded-lg">
                        ~{dailyRecommendation.estimated_minutes} min
                      </span>
                    </div>
                  </div>

                  <div class="flex items-center gap-4 mb-6">
                    <div class="flex items-center gap-2 px-4 py-2 bg-canvas-900/60 rounded-lg border border-canvas-700">
                      <svg class="w-5 h-5 text-amber-500" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M10.788 3.21c.448-1.077 1.976-1.077 2.424 0l2.082 5.007 5.404.433c1.164.093 1.636 1.545.749 2.305l-4.117 3.527 1.257 5.273c.271 1.136-.964 2.033-1.96 1.425L12 18.354 7.373 21.18c-.996.608-2.231-.29-1.96-1.425l1.257-5.273-4.117-3.527c-.887-.76-.415-2.212.749-2.305l5.404-.433 2.082-5.006z" />
                      </svg>
                      <span class="text-sm font-semibold text-white">+{dailyRecommendation.xp_reward} XP</span>
                    </div>
                  </div>

                  <button
                    onclick={startRecommendedPractice}
                    class="w-full px-6 py-4 bg-gradient-to-r from-lumera-600 to-focus-600 hover:from-lumera-500 hover:to-focus-500 text-white font-bold rounded-xl transition-all duration-300 hover:scale-105 hover:shadow-xl hover:shadow-lumera-500/30 flex items-center justify-center gap-2"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    COMENZAR AHORA
                  </button>
                {:else}
                  <div class="mb-6">
                    <p class="text-xl text-slate-200 mb-2">
                      Completa 10 preguntas de Comprensión Lectora
                    </p>
                    <p class="text-sm text-slate-400">
                      Recomendado para tu nivel actual
                    </p>
                  </div>

                  <div class="flex items-center gap-4 mb-6">
                    <div class="flex items-center gap-2 px-4 py-2 bg-canvas-900/60 rounded-lg border border-canvas-700">
                      <svg class="w-5 h-5 text-amber-500" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M10.788 3.21c.448-1.077 1.976-1.077 2.424 0l2.082 5.007 5.404.433c1.164.093 1.636 1.545.749 2.305l-4.117 3.527 1.257 5.273c.271 1.136-.964 2.033-1.96 1.425L12 18.354 7.373 21.18c-.996.608-2.231-.29-1.96-1.425l1.257-5.273-4.117-3.527c-.887-.76-.415-2.212.749-2.305l5.404-.433 2.082-5.006z" />
                      </svg>
                      <span class="text-sm font-semibold text-white">+50 XP</span>
                    </div>
                  </div>

                  <button
                    onclick={() => isMissionBoardOpen = true}
                    class="w-full px-6 py-4 bg-gradient-to-r from-lumera-600 to-focus-600 hover:from-lumera-500 hover:to-focus-500 text-white font-bold rounded-xl transition-all duration-300 hover:scale-105 hover:shadow-xl hover:shadow-lumera-500/30 flex items-center justify-center gap-2"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    COMENZAR AHORA
                  </button>
                {/if}
              </div>
            </div>
          </div>

          <!-- Subjects Card (1/3 width) -->
          <div class="col-span-1">
            <div class="relative group h-full">
              <!-- Animated glow effect -->
              <div class="absolute -inset-1 bg-gradient-to-r from-emerald-500 via-teal-500 to-cyan-500 rounded-3xl blur-lg opacity-30 group-hover:opacity-50 transition-opacity duration-300"></div>

              <div class="relative bg-canvas-800/90 backdrop-blur-md rounded-3xl p-6 border-2 border-white/10 h-full flex flex-col">
                <div class="flex items-center gap-3 mb-4">
                  <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-500 flex items-center justify-center">
                    <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                    </svg>
                  </div>
                  <h3 class="text-xl font-bold text-white">Mis Materias</h3>
                </div>

                <div class="flex-1 flex flex-col gap-3">
                  {#if dashboardStore.subjects.length > 0}
                    {#each dashboardStore.subjects as subject}
                      <button
                        onclick={() => window.location.href = `/materias/${subject.materiaId}/objetivos`}
                        class="w-full px-4 py-3 bg-canvas-900/60 hover:bg-canvas-900/80 border border-canvas-700 hover:border-emerald-500/50 rounded-xl transition-all duration-300 hover:scale-105 flex items-center justify-between group"
                      >
                        <div class="flex items-center gap-3">
                          <div class="w-8 h-8 rounded-lg bg-gradient-to-br {subject.color} flex items-center justify-center">
                            {#if subject.id === 'lyl' || subject.id === 'lenguaje' || subject.id === 'lengua'}
                              <!-- Lengua y Literatura - Book icon -->
                              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                              </svg>
                            {:else if subject.id === 'mat' || subject.id === 'matematicas'}
                              <!-- Matemáticas - Calculator icon -->
                              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                              </svg>
                            {:else}
                              <!-- Default emoji for other subjects -->
                              <span class="text-sm">{subject.icon}</span>
                            {/if}
                          </div>
                          <span class="text-sm font-semibold text-white">{subject.name}</span>
                        </div>
                        <svg class="w-4 h-4 text-slate-400 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                        </svg>
                      </button>
                    {/each}
                  {:else}
                    <div class="flex-1 flex items-center justify-center">
                      <p class="text-sm text-slate-400 text-center">Cargando materias...</p>
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Streak & Mini Stats Row -->
      <div class="flex flex-wrap items-center justify-center gap-6 mb-8">
        <!-- Streak Visualizer -->
        <div class="floating-stat">
          <div class="relative group">
            <div class="absolute -inset-1 bg-gradient-to-r from-orange-500 to-red-500 rounded-2xl blur opacity-25 group-hover:opacity-40 transition-opacity"></div>
            <div class="relative px-6 py-4 bg-canvas-800/90 backdrop-blur-md rounded-2xl border-2 border-orange-500/30 flex items-center gap-4">
              <!-- Animated flame icon -->
              <div class="relative">
                <svg class="w-12 h-12 text-orange-500 animate-pulse" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2.25c-2.429 0-4.817.178-7.152.521C2.87 3.061 1.5 4.795 1.5 6.741v6.018c0 1.946 1.37 3.68 3.348 3.97.877.129 1.761.234 2.652.316V21a.75.75 0 001.28.53l4.184-4.183a.39.39 0 01.266-.112c2.006-.05 3.982-.22 5.922-.506 1.978-.29 3.348-2.023 3.348-3.97V6.741c0-1.947-1.37-3.68-3.348-3.97A49.145 49.145 0 0012 2.25zM8.25 8.625a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25zm2.625 1.125a1.125 1.125 0 112.25 0 1.125 1.125 0 01-2.25 0zm4.875-1.125a1.125 1.125 0 100 2.25 1.125 1.125 0 000-2.25z" />
                </svg>
              </div>
              <div class="text-left">
                <p class="text-sm text-slate-400 font-medium">Racha Actual</p>
                <p class="text-3xl font-bold text-white">{student.streak} días</p>
              </div>
            </div>
          </div>
        </div>

        <!-- XP Today -->
        <div class="floating-stat">
          <div class="relative group">
            <div class="absolute -inset-1 bg-gradient-to-r from-purple-500 to-blue-500 rounded-2xl blur opacity-25 group-hover:opacity-40 transition-opacity"></div>
            <div class="relative px-6 py-4 bg-canvas-800/90 backdrop-blur-md rounded-2xl border-2 border-purple-500/30 flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center">
                <svg class="w-7 h-7 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M10.788 3.21c.448-1.077 1.976-1.077 2.424 0l2.082 5.007 5.404.433c1.164.093 1.636 1.545.749 2.305l-4.117 3.527 1.257 5.273c.271 1.136-.964 2.033-1.96 1.425L12 18.354 7.373 21.18c-.996.608-2.231-.29-1.96-1.425l1.257-5.273-4.117-3.527c-.887-.76-.415-2.212.749-2.305l5.404-.433 2.082-5.006z" />
                </svg>
              </div>
              <div class="text-left">
                <p class="text-sm text-slate-400 font-medium">XP Ganado Hoy</p>
                <p class="text-3xl font-bold text-white">+{student.xp}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Level Progress -->
        <div class="floating-stat">
          <div class="relative group">
            <div class="absolute -inset-1 bg-gradient-to-r from-focus-500 to-emerald-500 rounded-2xl blur opacity-25 group-hover:opacity-40 transition-opacity"></div>
            <div class="relative px-6 py-4 bg-canvas-800/90 backdrop-blur-md rounded-2xl border-2 border-focus-500/30 flex items-center gap-4">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-focus-500 to-emerald-500 flex items-center justify-center">
                <svg class="w-7 h-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
              <div class="text-left">
                <p class="text-sm text-slate-400 font-medium">Próximo Nivel</p>
                <p class="text-3xl font-bold text-white">{levelProgress()}%</p>
              </div>
            </div>
          </div>
        </div>
      </div>

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
  subjects={dashboardStore.subjects}
  {userProfile}
  {diagnosticLevels}
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
  gamificationStats={auth.gamificationStats}
  onAvatarChanged={handleAvatarChanged}
/>
