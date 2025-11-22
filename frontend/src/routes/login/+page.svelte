<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import gsap from 'gsap';
  import { auth } from '$lib/stores/auth.svelte';

  let activeTab = $state('login');
  let isLoading = $state(false);
  let errorMessage = $state('');

  // Login form
  let loginEmail = $state('');
  let loginPassword = $state('');

  // Register form
  let registerEmail = $state('');
  let registerName = $state('');
  let registerPassword = $state('');

  // Check if already authenticated
  onMount(() => {
    if (auth.checkAuth()) {
      goto('/');
      return;
    }

    gsap.from('.auth-card', {
      duration: 0.8,
      y: 50,
      opacity: 0,
      ease: 'power3.out'
    });
  });

  async function handleLogin() {
    if (!loginEmail || !loginPassword) {
      errorMessage = 'Please fill in all fields';
      return;
    }

    isLoading = true;
    errorMessage = '';

    const result = await auth.login(loginEmail, loginPassword);

    if (result.success) {
      // Check if user has profile
      const hasProfile = await auth.checkIfHasProfile();
      if (hasProfile) {
        goto('/');
      } else {
        goto('/onboarding');
      }
    } else {
      errorMessage = result.error || 'Login failed';
      isLoading = false;
    }
  }

  async function handleRegister() {
    if (!registerEmail || !registerName || !registerPassword) {
      errorMessage = 'Please fill in all fields';
      return;
    }

    if (registerPassword.length < 6) {
      errorMessage = 'Password must be at least 6 characters';
      return;
    }

    isLoading = true;
    errorMessage = '';

    const result = await auth.register(registerEmail, registerName, registerPassword);

    if (result.success) {
      // New users don't have profile, redirect to onboarding
      goto('/onboarding');
    } else {
      errorMessage = result.error || 'Registration failed';
      isLoading = false;
    }
  }

  function clearError() {
    errorMessage = '';
  }
</script>

<svelte:head>
  <title>Login - Lumera App</title>
</svelte:head>

<main class="min-h-screen bg-gradient-to-br from-slate-950 via-indigo-950 to-slate-950 flex items-center justify-center p-4">
  <div class="auth-card bg-canvas-900/40 backdrop-blur-xl rounded-2xl shadow-2xl p-8 max-w-md w-full border border-white/10">
    <!-- Logo/Title -->
    <div class="text-center mb-8">
      <h1 class="text-4xl font-bold text-white mb-2">Lumera App</h1>
      <p class="text-slate-400">Student Learning Platform</p>
    </div>

    <!-- Tabs -->
    <div class="flex p-1 mb-6 rounded-xl bg-canvas-950/50 border border-slate-800">
      <button
        class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium transition-all {activeTab === 'login' ? 'bg-lumera-600 text-white shadow-lg shadow-indigo-900/50' : 'text-slate-400 hover:text-white'}"
        onclick={() => { activeTab = 'login'; clearError(); }}
      >
        Login
      </button>
      <button
        class="flex-1 px-4 py-2.5 rounded-lg text-sm font-medium transition-all {activeTab === 'register' ? 'bg-lumera-600 text-white shadow-lg shadow-indigo-900/50' : 'text-slate-400 hover:text-white'}"
        onclick={() => { activeTab = 'register'; clearError(); }}
      >
        Register
      </button>
    </div>

    <!-- Error Message -->
    {#if errorMessage}
      <div class="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
        <p class="text-red-400 text-sm">{errorMessage}</p>
      </div>
    {/if}

    <!-- Login Form -->
    {#if activeTab === 'login'}
      <form onsubmit={(e) => { e.preventDefault(); handleLogin(); }} class="space-y-4">
        <div>
          <label for="login-email" class="block text-sm font-medium text-slate-300 mb-2">Email</label>
          <input
            id="login-email"
            type="email"
            bind:value={loginEmail}
            oninput={clearError}
            disabled={isLoading}
            class="w-full px-4 py-3 bg-canvas-950/50 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all disabled:opacity-50"
            placeholder="your@email.com"
            required
          />
        </div>

        <div>
          <label for="login-password" class="block text-sm font-medium text-slate-300 mb-2">Password</label>
          <input
            id="login-password"
            type="password"
            bind:value={loginPassword}
            oninput={clearError}
            disabled={isLoading}
            class="w-full px-4 py-3 bg-canvas-950/50 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all disabled:opacity-50"
            placeholder="••••••••"
            required
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
          class="w-full mt-6 px-6 py-3 bg-lumera-600 hover:bg-lumera-700 disabled:bg-lumera-600/50 text-white font-semibold rounded-lg transition-all shadow-lg hover:shadow-indigo-500/50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {#if isLoading}
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Logging in...
          {:else}
            Login
          {/if}
        </button>
      </form>
    {/if}

    <!-- Register Form -->
    {#if activeTab === 'register'}
      <form onsubmit={(e) => { e.preventDefault(); handleRegister(); }} class="space-y-4">
        <div>
          <label for="register-name" class="block text-sm font-medium text-slate-300 mb-2">Name</label>
          <input
            id="register-name"
            type="text"
            bind:value={registerName}
            oninput={clearError}
            disabled={isLoading}
            class="w-full px-4 py-3 bg-canvas-950/50 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all disabled:opacity-50"
            placeholder="Your Name"
            required
          />
        </div>

        <div>
          <label for="register-email" class="block text-sm font-medium text-slate-300 mb-2">Email</label>
          <input
            id="register-email"
            type="email"
            bind:value={registerEmail}
            oninput={clearError}
            disabled={isLoading}
            class="w-full px-4 py-3 bg-canvas-950/50 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all disabled:opacity-50"
            placeholder="your@email.com"
            required
          />
        </div>

        <div>
          <label for="register-password" class="block text-sm font-medium text-slate-300 mb-2">Password</label>
          <input
            id="register-password"
            type="password"
            bind:value={registerPassword}
            oninput={clearError}
            disabled={isLoading}
            class="w-full px-4 py-3 bg-canvas-950/50 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-indigo-500/20 transition-all disabled:opacity-50"
            placeholder="••••••••"
            minlength="6"
            required
          />
          <p class="mt-1 text-xs text-slate-500">Minimum 6 characters</p>
        </div>

        <button
          type="submit"
          disabled={isLoading}
          class="w-full mt-6 px-6 py-3 bg-lumera-600 hover:bg-lumera-700 disabled:bg-lumera-600/50 text-white font-semibold rounded-lg transition-all shadow-lg hover:shadow-indigo-500/50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {#if isLoading}
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating account...
          {:else}
            Create Account
          {/if}
        </button>
      </form>
    {/if}

    <!-- Footer -->
    <div class="mt-8 text-center">
      <p class="text-slate-500 text-sm">
        Platanus Hack 25 · Student Learning Platform
      </p>
    </div>
  </div>
</main>
