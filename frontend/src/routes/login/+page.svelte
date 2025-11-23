<script lang="ts">
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

  // Particle system
  let particles = $state<Array<{ x: number; y: number; vx: number; vy: number; size: number; color: string }>>([]);
  let mouseX = $state(0);
  let mouseY = $state(0);
  let mainRef: HTMLElement | null = null;

  // Initialize particles
  function initParticles() {
    const particleColors = ['#2563eb', '#14b8a6', '#f59e0b', '#8b5cf6'];
    particles = Array.from({ length: 25 }, () => ({
      x: Math.random() * 100,
      y: Math.random() * 100,
      vx: (Math.random() - 0.5) * 0.3,
      vy: (Math.random() - 0.5) * 0.3,
      size: Math.random() * 3 + 2,
      color: particleColors[Math.floor(Math.random() * particleColors.length)]
    }));
  }

  // Animate particles
  function animateParticles() {
    particles = particles.map(p => {
      let newX = p.x + p.vx;
      let newY = p.y + p.vy;

      // Bounce off edges
      if (newX < 0 || newX > 100) p.vx *= -1;
      if (newY < 0 || newY > 100) p.vy *= -1;

      // Mouse interaction
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

  // Handle mouse move
  function handleMouseMove(e: MouseEvent) {
    if (mainRef) {
      const rect = mainRef.getBoundingClientRect();
      mouseX = ((e.clientX - rect.left) / rect.width) * 100;
      mouseY = ((e.clientY - rect.top) / rect.height) * 100;
    }
  }

  // Check if already authenticated
  onMount(() => {
    if (auth.checkAuth()) {
      goto('/');
      return;
    }

    // Initialize particles
    initParticles();
    const particleInterval = setInterval(animateParticles, 50);

    // Animate card entrance
    gsap.from('.auth-card', {
      duration: 0.8,
      y: 50,
      opacity: 0,
      scale: 0.95,
      ease: 'power3.out'
    });

    gsap.from('.logo-section', {
      duration: 1,
      y: -30,
      opacity: 0,
      delay: 0.2,
      ease: 'power2.out'
    });

    return () => {
      clearInterval(particleInterval);
    };
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

<main
  bind:this={mainRef}
  onmousemove={handleMouseMove}
  class="relative min-h-screen bg-canvas-950 flex items-center justify-center p-4 overflow-hidden"
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

  <!-- Main Content -->
  <div class="relative z-10 w-full max-w-md">
    <!-- Logo/Title -->
    <div class="logo-section text-center mb-8">
      <div class="inline-block px-8 py-4 rounded-2xl bg-gradient-to-r from-lumera-500/20 via-focus-500/20 to-purple-500/20 border-2 border-white/10 backdrop-blur-sm mb-4">
        <h1 class="text-5xl font-bold text-white" style="text-shadow: 0 2px 8px rgba(0, 0, 0, 0.5), 0 0 20px rgba(37, 99, 235, 0.3);">Lumera</h1>
      </div>
      <p class="text-slate-400 text-lg">Tu plataforma de aprendizaje</p>
    </div>

    <!-- Auth Card -->
    <div class="auth-card relative group">
      <!-- Animated glow effect -->
      <div class="absolute -inset-1 bg-gradient-to-r from-lumera-500 via-focus-500 to-purple-500 rounded-3xl blur-lg opacity-25 group-hover:opacity-40 transition-opacity duration-300"></div>

      <div class="relative bg-canvas-800/90 backdrop-blur-xl rounded-3xl shadow-2xl p-8 border-2 border-white/10">

        <!-- Tabs -->
        <div class="flex p-1 mb-6 rounded-xl bg-canvas-900/60 border border-canvas-700">
          <button
            class="flex-1 px-4 py-3 rounded-lg text-sm font-semibold transition-all duration-300 {activeTab === 'login' ? 'bg-[#E1E1E1] text-canvas-900 shadow-lg' : 'text-slate-400 hover:text-white'}"
            onclick={() => { activeTab = 'login'; clearError(); }}
          >
            Iniciar Sesión
          </button>
          <button
            class="flex-1 px-4 py-3 rounded-lg text-sm font-semibold transition-all duration-300 {activeTab === 'register' ? 'bg-[#E1E1E1] text-canvas-900 shadow-lg' : 'text-slate-400 hover:text-white'}"
            onclick={() => { activeTab = 'register'; clearError(); }}
          >
            Crear Cuenta
          </button>
        </div>

        <!-- Error Message -->
        {#if errorMessage}
          <div class="mb-6 p-4 bg-red-500/10 border-2 border-red-500/30 rounded-xl backdrop-blur-sm">
            <div class="flex items-center gap-2">
              <svg class="w-5 h-5 text-red-400 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
              </svg>
              <p class="text-red-400 text-sm font-medium">{errorMessage}</p>
            </div>
          </div>
        {/if}

        <!-- Login Form -->
        {#if activeTab === 'login'}
          <form onsubmit={(e) => { e.preventDefault(); handleLogin(); }} class="space-y-5">
            <div>
              <label for="login-email" class="block text-sm font-semibold text-slate-300 mb-2">Correo Electrónico</label>
              <input
                id="login-email"
                type="email"
                bind:value={loginEmail}
                oninput={clearError}
                disabled={isLoading}
                class="w-full px-4 py-3.5 bg-canvas-900/60 border-2 border-canvas-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-lumera-500/20 transition-all disabled:opacity-50 hover:border-canvas-600"
                placeholder="tu@email.com"
                required
              />
            </div>

            <div>
              <label for="login-password" class="block text-sm font-semibold text-slate-300 mb-2">Contraseña</label>
              <input
                id="login-password"
                type="password"
                bind:value={loginPassword}
                oninput={clearError}
                disabled={isLoading}
                class="w-full px-4 py-3.5 bg-canvas-900/60 border-2 border-canvas-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-lumera-500/20 transition-all disabled:opacity-50 hover:border-canvas-600"
                placeholder="••••••••"
                required
              />
            </div>

            <button
              type="submit"
              disabled={isLoading}
              class="w-full mt-6 px-6 py-4 bg-[#E1E1E1] hover:bg-[#CCCCCC] disabled:bg-[#E1E1E1]/50 text-canvas-900 font-bold rounded-xl transition-all duration-300 shadow-lg hover:shadow-xl disabled:cursor-not-allowed flex items-center justify-center gap-2 hover:scale-105 disabled:hover:scale-100"
            >
              {#if isLoading}
                <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Iniciando sesión...
              {:else}
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
                </svg>
                INICIAR SESIÓN
              {/if}
            </button>
          </form>
        {/if}

        <!-- Register Form -->
        {#if activeTab === 'register'}
          <form onsubmit={(e) => { e.preventDefault(); handleRegister(); }} class="space-y-5">
            <div>
              <label for="register-name" class="block text-sm font-semibold text-slate-300 mb-2">Nombre Completo</label>
              <input
                id="register-name"
                type="text"
                bind:value={registerName}
                oninput={clearError}
                disabled={isLoading}
                class="w-full px-4 py-3.5 bg-canvas-900/60 border-2 border-canvas-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-lumera-500/20 transition-all disabled:opacity-50 hover:border-canvas-600"
                placeholder="Tu Nombre"
                required
              />
            </div>

            <div>
              <label for="register-email" class="block text-sm font-semibold text-slate-300 mb-2">Correo Electrónico</label>
              <input
                id="register-email"
                type="email"
                bind:value={registerEmail}
                oninput={clearError}
                disabled={isLoading}
                class="w-full px-4 py-3.5 bg-canvas-900/60 border-2 border-canvas-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-lumera-500/20 transition-all disabled:opacity-50 hover:border-canvas-600"
                placeholder="tu@email.com"
                required
              />
            </div>

            <div>
              <label for="register-password" class="block text-sm font-semibold text-slate-300 mb-2">Contraseña</label>
              <input
                id="register-password"
                type="password"
                bind:value={registerPassword}
                oninput={clearError}
                disabled={isLoading}
                class="w-full px-4 py-3.5 bg-canvas-900/60 border-2 border-canvas-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-lumera-500 focus:ring-2 focus:ring-lumera-500/20 transition-all disabled:opacity-50 hover:border-canvas-600"
                placeholder="••••••••"
                minlength="6"
                required
              />
              <p class="mt-2 text-xs text-slate-500">Mínimo 6 caracteres</p>
            </div>

            <button
              type="submit"
              disabled={isLoading}
              class="w-full mt-6 px-6 py-4 bg-[#E1E1E1] hover:bg-[#CCCCCC] disabled:bg-[#E1E1E1]/50 text-canvas-900 font-bold rounded-xl transition-all duration-300 shadow-lg hover:shadow-xl disabled:cursor-not-allowed flex items-center justify-center gap-2 hover:scale-105 disabled:hover:scale-100"
            >
              {#if isLoading}
                <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Creando cuenta...
              {:else}
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                </svg>
                CREAR CUENTA
              {/if}
            </button>
          </form>
        {/if}

        <!-- Footer -->
        <div class="mt-8 pt-6 border-t border-canvas-700/50 text-center">
          <p class="text-slate-500 text-sm">
            Platanus Hack 25 · Plataforma de Aprendizaje
          </p>
        </div>
      </div>
    </div>
  </div>
</main>
