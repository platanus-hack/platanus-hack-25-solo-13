<script lang="ts">
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    isOpen = false,
    oaTitle = "Objetivo de Aprendizaje",
    onClose = null
  } = $props();

  let modalRef: HTMLDivElement | null = null;
  let step = $state(0);
  const totalSteps = 4;

  const ingredients = [
    { label: "Preferencias", icon: "user" },
    { label: "Conocimientos", icon: "book" },
    { label: "Brechas", icon: "chart" },
    { label: "IA", icon: "brain" }
  ];

  onMount(() => {
    if (isOpen && modalRef) {
      animateModal();
    }
  });

  $effect(() => {
    if (isOpen && modalRef) {
      animateModal();
    }
  });

  function animateModal() {
    if (!modalRef) return;

    // Reset
    step = 0;
    gsap.set('.ingredient-card', { opacity: 0, y: 30 });
    gsap.set('.bowl-container', { opacity: 0, scale: 0.8 });
    gsap.set('.connecting-line', { scaleX: 0 });

    // Animate backdrop
    gsap.fromTo('.modal-backdrop',
      { opacity: 0 },
      { opacity: 1, duration: 0.3 }
    );

    // Animate modal entrance
    gsap.fromTo('.modal-content',
      { scale: 0.95, opacity: 0 },
      { scale: 1, opacity: 1, duration: 0.4, ease: 'power2.out' }
    );

    // Start gears rotating continuously
    gsap.to('.gear', {
      rotation: 360,
      duration: 4,
      ease: 'none',
      repeat: -1
    });

    gsap.to('.gear-reverse', {
      rotation: -360,
      duration: 3,
      ease: 'none',
      repeat: -1
    });

    // Animate ingredients flowing into bowl
    setTimeout(() => {
      animateRecipeFlow();
    }, 500);
  }

  function animateRecipeFlow() {
    const timeline = gsap.timeline();

    // Show all ingredient cards first
    timeline.to('.ingredient-card', {
      opacity: 1,
      y: 0,
      duration: 0.5,
      stagger: 0.15,
      ease: 'power2.out'
    });

    // Show bowl
    timeline.to('.bowl-container', {
      opacity: 1,
      scale: 1,
      duration: 0.5,
      ease: 'back.out(1.7)'
    }, '-=0.3');

    // Animate each ingredient flowing to the bowl
    ingredients.forEach((ingredient, index) => {
      timeline.call(() => {
        step = index + 1;
        animateIngredientFlow(index);
      }, [], `+=${index === 0 ? 0.5 : 1.2}`);
    });

    // Loop
    timeline.call(() => {
      step = 0;
      animateRecipeFlow();
    }, [], '+=1.5');
  }

  function animateIngredientFlow(index: number) {
    const card = document.querySelector(`.ingredient-card-${index}`);
    const bowl = document.querySelector('.bowl-container');
    const line = document.querySelector(`.connecting-line-${index}`);

    if (!card || !bowl) return;

    // Highlight current ingredient
    gsap.to(card, {
      scale: 1.05,
      duration: 0.3,
      yoyo: true,
      repeat: 1
    });

    // Animate connecting line
    gsap.to(line, {
      scaleX: 1,
      duration: 0.6,
      ease: 'power2.inOut'
    });

    // Create particle flow
    createParticleFlow(index);

    // Pulse bowl
    gsap.fromTo(bowl,
      { scale: 1 },
      {
        scale: 1.08,
        duration: 0.3,
        yoyo: true,
        repeat: 1,
        ease: 'power2.inOut'
      }
    );

    // Hide line after animation
    setTimeout(() => {
      gsap.to(line, {
        scaleX: 0,
        duration: 0.3
      });
    }, 800);
  }

  function createParticleFlow(ingredientIndex: number) {
    const container = document.querySelector('.particles-container');
    if (!container) return;

    for (let i = 0; i < 8; i++) {
      const particle = document.createElement('div');
      particle.className = 'flow-particle';
      container.appendChild(particle);

      const startX = -200 + (ingredientIndex * 100);
      const startY = -150 + (ingredientIndex * 80);

      gsap.fromTo(particle,
        {
          x: startX,
          y: startY,
          opacity: 1,
          scale: 1
        },
        {
          x: 0,
          y: 50,
          opacity: 0,
          scale: 0.3,
          duration: 0.8,
          delay: i * 0.05,
          ease: 'power2.in',
          onComplete: () => particle.remove()
        }
      );
    }
  }

  function handleClose() {
    gsap.to('.modal-content', {
      scale: 0.95,
      opacity: 0,
      duration: 0.3,
      ease: 'power2.in',
      onComplete: () => {
        if (onClose) onClose();
      }
    });

    gsap.to('.modal-backdrop', {
      opacity: 0,
      duration: 0.3
    });
  }
</script>

{#if isOpen}
  <div bind:this={modalRef} class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div class="modal-backdrop absolute inset-0 bg-white/90 backdrop-blur-sm" onclick={handleClose}></div>

    <!-- Modal Content -->
    <div class="modal-content relative bg-white rounded-2xl shadow-2xl max-w-2xl w-full border-2 border-slate-200"
         onclick={(e) => e.stopPropagation()}>

      <!-- Close button -->
      <button
        onclick={handleClose}
        class="absolute top-4 right-4 w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 flex items-center justify-center text-slate-600 hover:text-slate-900 transition-all hover:scale-110 z-20"
        title="Cerrar"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <div class="p-8">
        <!-- Header -->
        <div class="text-center mb-6">
          <div class="flex items-center justify-center gap-2 mb-2">
            <svg class="w-5 h-5 text-lumera-600 gear" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <h2 class="text-xl font-bold text-slate-900">Creando Plan Personalizado</h2>
            <svg class="w-4 h-4 text-lumera-600 gear-reverse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <p class="text-slate-600 text-xs font-medium mb-1">{oaTitle}</p>
          <p class="text-slate-500 text-xs italic">Esto puede tomar 2-3 minutos</p>
        </div>

        <!-- Recipe Visualization -->
        <div class="relative">
          <!-- Particles container -->
          <div class="particles-container absolute inset-0 pointer-events-none"></div>

          <!-- Ingredients Grid (Top) -->
          <div class="grid grid-cols-4 gap-3 mb-8">
            {#each ingredients as ingredient, i}
              <div class="ingredient-card ingredient-card-{i} relative">
                <!-- Connecting line -->
                <div class="connecting-line connecting-line-{i} absolute top-full left-1/2 -translate-x-1/2 w-0.5 h-12 bg-gradient-to-b from-lumera-400 to-transparent origin-top" style="transform-origin: top;"></div>

                <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl p-4 border-2 border-slate-200 hover:border-lumera-300 transition-all text-center">
                  <!-- Icon -->
                  <div class="flex justify-center mb-2">
                    {#if ingredient.icon === 'user'}
                      <svg class="w-6 h-6 text-lumera-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                      </svg>
                    {:else if ingredient.icon === 'book'}
                      <svg class="w-6 h-6 text-lumera-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                      </svg>
                    {:else if ingredient.icon === 'chart'}
                      <svg class="w-6 h-6 text-lumera-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                      </svg>
                    {:else if ingredient.icon === 'brain'}
                      <svg class="w-6 h-6 text-lumera-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                      </svg>
                    {/if}
                  </div>
                  <p class="text-xs font-bold text-slate-700">{ingredient.label}</p>
                </div>
              </div>
            {/each}
          </div>

          <!-- Bowl (Bottom Center) -->
          <div class="bowl-container flex justify-center">
            <div class="relative bg-gradient-to-br from-lumera-50 to-focus-50 rounded-2xl p-6 border-3 border-lumera-400 shadow-lg w-52">
              <!-- Gears decoration -->
              <div class="absolute -top-2 -left-2">
                <svg class="w-6 h-6 text-lumera-500 gear" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>
              <div class="absolute -top-2 -right-2">
                <svg class="w-5 h-5 text-focus-500 gear-reverse" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                </svg>
              </div>

              <div class="text-center">
                <svg class="w-12 h-12 mx-auto mb-2 text-lumera-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <p class="text-base font-bold text-slate-900 mb-0.5">Tu Plan</p>
                <p class="text-xs text-slate-600">Combinando...</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Status text -->
        <div class="mt-6 text-center">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 bg-slate-100 rounded-full">
            <div class="flex gap-1">
              <div class="w-1.5 h-1.5 bg-lumera-500 rounded-full animate-pulse"></div>
              <div class="w-1.5 h-1.5 bg-lumera-500 rounded-full animate-pulse" style="animation-delay: 0.2s;"></div>
              <div class="w-1.5 h-1.5 bg-lumera-500 rounded-full animate-pulse" style="animation-delay: 0.4s;"></div>
            </div>
            <span class="text-xs font-medium text-slate-700">
              {#if step > 0}
                Añadiendo {ingredients[step - 1].label}...
              {:else}
                Preparando ingredientes...
              {/if}
            </span>
          </div>
        </div>

        <!-- Background link -->
        <div class="mt-5 text-center">
          <button
            onclick={handleClose}
            class="group inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-slate-50 hover:bg-slate-100 border border-slate-200 hover:border-lumera-300 text-slate-700 hover:text-lumera-700 transition-all"
          >
            <svg class="w-3.5 h-3.5 group-hover:scale-110 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
            <span class="text-xs font-semibold">Dejar en segundo plano</span>
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .flow-particle {
    position: absolute;
    width: 6px;
    height: 6px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    border-radius: 50%;
    pointer-events: none;
    top: 50%;
    left: 50%;
  }
</style>
