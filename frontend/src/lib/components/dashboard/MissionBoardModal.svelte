<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  interface Mission {
    id: number;
    subject: string;
    title: string;
    time: string;
    reward: string;
    state: string;
  }

  interface Missions {
    daily: Mission[];
    weekly: Mission[];
    story: Mission[];
    side: Mission[];
  }

  interface Props {
    missions: Missions;
    activeTab: string;
    onTabChange: (tab: string) => void;
    isOpen: boolean;
    onClose: () => void;
  }

  let { missions, activeTab, onTabChange, isOpen, onClose }: Props = $props();
  let backdropRef = $state<HTMLDivElement | null>(null);
  let modalRef = $state<HTMLDivElement | null>(null);

  // Close on Escape key
  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      onClose();
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleKeydown);
    return () => window.removeEventListener('keydown', handleKeydown);
  });

  // Animate modal entrance
  $effect(() => {
    if (isOpen && backdropRef && modalRef) {
      gsap.fromTo(backdropRef,
        { opacity: 0 },
        { opacity: 1, duration: 0.2, ease: 'power2.out' }
      );
      gsap.fromTo(modalRef,
        { opacity: 0, y: 20, scale: 0.95 },
        { opacity: 1, y: 0, scale: 1, duration: 0.3, ease: 'power2.out', delay: 0.1 }
      );
    }
  });

  function getSubjectColor(subject: string) {
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

  // Mission type metadata
  const missionTypes = [
    { id: 'daily', label: 'Daily', icon: '☀️', desc: 'Misiones diarias' },
    { id: 'weekly', label: 'Weekly', icon: '📅', desc: 'Misiones semanales' },
    { id: 'story', label: 'Story', icon: '📖', desc: 'Modo campaña' },
    { id: 'side', label: 'Side', icon: '🎯', desc: 'Desafíos extras' }
  ];

  // Get mission count for a type
  function getMissionCount(type: string): number {
    return missions[type]?.filter(m => m.state !== 'done').length || 0;
  }
</script>

{#if isOpen}
  <!-- Backdrop -->
  <div
    bind:this={backdropRef}
    class="fixed inset-0 bg-black/60 backdrop-blur-sm z-60 flex items-center justify-center p-4"
    onclick={onClose}
  >
    <!-- Modal -->
    <div
      bind:this={modalRef}
      class="bg-canvas-900 rounded-2xl border border-slate-800 max-w-5xl w-full h-[80vh] overflow-hidden shadow-2xl flex flex-col"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- Header -->
      <div class="flex-shrink-0 bg-canvas-900 border-b border-slate-800 p-6 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="h-12 w-12 rounded-xl bg-gradient-to-br from-lumera-600/20 to-purple-600/20 border border-lumera-500/30 flex items-center justify-center text-2xl">
            ⚔️
          </div>
          <div>
            <h2 class="text-xl font-bold text-white">Mission Board</h2>
            <p class="text-xs text-slate-400 mt-0.5">Tus misiones activas</p>
          </div>
        </div>
        <button
          onclick={onClose}
          class="h-10 w-10 rounded-full bg-canvas-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors"
        >
          ✕
        </button>
      </div>

      <!-- Content: Two Column Layout -->
      <div class="grid grid-cols-12 flex-1 overflow-hidden">
        <!-- Left Column: Mission Types -->
        <div class="col-span-4 border-r border-slate-800 p-4 space-y-2">
          {#each missionTypes as type}
            <button
              onclick={() => onTabChange(type.id)}
              class="
                w-full p-4 rounded-xl text-left transition-all duration-200
                {activeTab === type.id
                  ? 'bg-lumera-600 border-lumera-500 shadow-lg shadow-lumera-600/20'
                  : 'bg-canvas-900/50 border-canvas-800/80 hover:bg-canvas-800 hover:border-canvas-700'
                }
                border
              "
            >
              <div class="flex items-center gap-3 mb-2">
                <span class="text-2xl">{type.icon}</span>
                <div class="flex-1">
                  <h3 class="font-bold {activeTab === type.id ? 'text-white' : 'text-slate-200'}">{type.label}</h3>
                  <p class="text-xs {activeTab === type.id ? 'text-lumera-200' : 'text-slate-500'}">{type.desc}</p>
                </div>
                {#if getMissionCount(type.id) > 0}
                  <span class="
                    px-2 py-1 rounded-full text-xs font-bold
                    {activeTab === type.id
                      ? 'bg-white text-lumera-600'
                      : 'bg-purple-600 text-white'
                    }
                  ">
                    {getMissionCount(type.id)}
                  </span>
                {/if}
              </div>
            </button>
          {/each}
        </div>

        <!-- Right Column: Missions List -->
        <div class="col-span-8 overflow-y-auto p-4">
          {#if missions[activeTab] && missions[activeTab].length > 0}
            <div class="space-y-3">
              {#each missions[activeTab] as mission}
                <button class="
                  w-full p-4 rounded-xl bg-canvas-900/50 border border-canvas-800/80
                  hover:border-lumera-500/40 hover:bg-canvas-800/60
                  transition-all flex items-center gap-4 text-left
                ">
                  <div class="h-12 w-12 rounded-lg {getSubjectColor(mission.subject)} flex items-center justify-center text-lg font-bold flex-shrink-0">
                    {mission.subject[0]}
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-xs font-bold text-slate-500 uppercase">{mission.subject}</span>
                      {#if mission.state === 'done'}
                        <span class="text-[10px] px-1.5 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/20">Done</span>
                      {:else if mission.state === 'progress'}
                        <span class="text-[10px] px-1.5 py-0.5 rounded bg-lumera-500/20 text-lumera-400 border border-lumera-500/20">In Progress</span>
                      {/if}
                    </div>
                    <h4 class="font-medium text-slate-200 mb-1">{mission.title}</h4>
                    <div class="flex items-center gap-3 text-xs text-slate-500">
                      <span>⏱️ {mission.time}</span>
                      <span class="text-achievement-400">🎁 {mission.reward}</span>
                    </div>
                  </div>
                  <div class="flex-shrink-0">
                    <svg class="w-5 h-5 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </button>
              {/each}
            </div>
          {:else}
            <!-- Empty State -->
            <div class="flex items-center justify-center h-full">
              <div class="text-center">
                <div class="text-6xl mb-4">📋</div>
                <h3 class="text-lg font-semibold text-white mb-2">No hay misiones {activeTab}</h3>
                <p class="text-sm text-slate-400">
                  Completa tus misiones actuales para desbloquear nuevas
                </p>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
{/if}
