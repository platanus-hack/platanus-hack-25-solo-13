<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  let healthStatus = $state({
    status: 'loading...',
    database: 'loading...',
    timestamp: null,
    error: null
  });

  async function checkHealth() {
    try {
      const response = await fetch('/api/health');
      const data = await response.json();
      healthStatus = {
        status: data.status,
        database: data.database,
        timestamp: new Date(data.timestamp).toLocaleString(),
        error: null
      };
    } catch (error) {
      healthStatus = {
        status: 'error',
        database: 'error',
        timestamp: null,
        error: error.message
      };
    }
  }

  onMount(() => {
    checkHealth();
    gsap.from('.health-card', {
      duration: 0.8,
      y: 50,
      opacity: 0,
      ease: 'power3.out'
    });
  });
</script>

<svelte:head>
  <title>System Health - Lumera App</title>
</svelte:head>

<main class="min-h-screen bg-gradient-to-br from-gray-900 via-purple-900 to-gray-900 flex items-center justify-center p-4">
  <div class="health-card bg-white/10 backdrop-blur-lg rounded-2xl shadow-2xl p-8 max-w-md w-full border border-white/20">
    <h1 class="text-4xl font-bold text-white mb-2">System Health</h1>
    <p class="text-purple-300 mb-8">Backend & Database Status</p>

    <div class="space-y-4">
      <div class="bg-black/30 rounded-lg p-4 border border-white/10">
        <div class="flex items-center justify-between mb-2">
          <span class="text-gray-300 font-medium">API Status</span>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 rounded-full {healthStatus.status === 'ok' ? 'bg-green-500' : healthStatus.status === 'error' ? 'bg-red-500' : 'bg-yellow-500'}"></div>
            <span class="text-white font-semibold">{healthStatus.status}</span>
          </div>
        </div>
      </div>

      <div class="bg-black/30 rounded-lg p-4 border border-white/10">
        <div class="flex items-center justify-between mb-2">
          <span class="text-gray-300 font-medium">Database</span>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 rounded-full {healthStatus.database === 'connected' ? 'bg-green-500' : 'bg-red-500'}"></div>
            <span class="text-white font-semibold">{healthStatus.database}</span>
          </div>
        </div>
      </div>

      {#if healthStatus.timestamp}
        <div class="bg-black/30 rounded-lg p-4 border border-white/10">
          <span class="text-gray-300 text-sm">Last check: {healthStatus.timestamp}</span>
        </div>
      {/if}

      {#if healthStatus.error}
        <div class="bg-red-500/20 rounded-lg p-4 border border-red-500/50">
          <p class="text-red-200 text-sm">{healthStatus.error}</p>
        </div>
      {/if}
    </div>

    <button
      onclick={checkHealth}
      class="mt-8 w-full bg-purple-600 hover:bg-purple-700 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200 shadow-lg hover:shadow-purple-500/50"
    >
      Refresh Status
    </button>

    <div class="mt-6 text-center">
      <p class="text-gray-400 text-sm">
        Go + SvelteKit + PostgreSQL + Docker
      </p>
    </div>

    <div class="mt-4 text-center">
      <a href="/" class="text-indigo-400 hover:text-indigo-300 text-sm">← Back to Dashboard</a>
    </div>
  </div>
</main>
