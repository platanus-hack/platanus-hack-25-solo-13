<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    palabraObjetivo = "",
    pronunciacion = "",
    audioUrl = null,
    etimologia = "",
    definicion = "",
    morfologia = null, // { prefijo, raiz, sufijo }
    contextosEjemplo = [],
    imagenes = [],
    sinonimos = [],
    antonimos = [],
    palabrasFamilia = [],
    materia = "lenguaje",
    onNext = null,
    onPrevious = null,
    showNavigation = true
  } = $props();

  // Estados locales
  let containerRef = $state(null);
  let contextoActual = $state(0);
  let mostrarQuiz = $state(false);
  let respuestaQuiz = $state(null);

  // Navegación del carousel
  function anteriorContexto() {
    contextoActual = (contextoActual - 1 + contextosEjemplo.length) % contextosEjemplo.length;
  }

  function siguienteContexto() {
    contextoActual = (contextoActual + 1) % contextosEjemplo.length;
  }

  // Reproducir audio
  function playAudio() {
    if (audioUrl) {
      const audio = new Audio(audioUrl);
      audio.play();
    }
  }

  // Animación de entrada
  onMount(() => {
    if (containerRef) {
      gsap.from(containerRef, {
        opacity: 0,
        y: 30,
        duration: 0.6,
        ease: 'power2.out'
      });
    }
  });
</script>

<div
  bind:this={containerRef}
  class="w-full max-w-7xl mx-auto p-8 bg-canvas-950 rounded-2xl border border-canvas-800 shadow-2xl"
>
  <!-- Header -->
  <div class="mb-8">
    <div class="flex items-center gap-3 mb-4">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase bg-purple-500/20 text-purple-400 border border-purple-500">
        {materia}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-medium bg-canvas-800 text-canvas-300">
        Vocabulario
      </span>
      <div class="flex-1 h-px bg-canvas-800"></div>
    </div>

    <h2 class="text-3xl font-bold text-white mb-2">
      Explorando: <span class="text-purple-400">{palabraObjetivo}</span>
    </h2>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
    <!-- Columna principal -->
    <div class="lg:col-span-2 space-y-6">
      <!-- Palabra objetivo con pronunciación -->
      <div class="p-6 bg-gradient-to-br from-purple-500/10 to-purple-600/10 rounded-2xl border-2 border-purple-500/30">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-5xl font-bold text-white mb-2">{palabraObjetivo}</h3>
            <div class="flex items-center gap-3">
              <p class="text-lg text-purple-300">{pronunciacion}</p>
              {#if audioUrl}
                <button
                  onclick={playAudio}
                  class="
                    px-3 py-1 rounded-lg text-xs font-semibold
                    bg-purple-500/20 text-purple-400 border border-purple-500
                    transition-all duration-300
                    hover:bg-purple-500/30 hover:scale-105
                  "
                >
                  🔊 Escuchar
                </button>
              {/if}
            </div>
          </div>
        </div>

        <p class="text-xl text-canvas-200 leading-relaxed">
          {definicion}
        </p>
      </div>

      <!-- Etimología y morfología -->
      {#if etimologia || morfologia}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            🌳 Etimología
          </h3>

          {#if etimologia}
            <p class="text-canvas-300 mb-4 italic">
              {etimologia}
            </p>
          {/if}

          {#if morfologia}
            <div class="flex items-center gap-2 mt-4">
              {#if morfologia.prefijo}
                <div class="flex-1 p-3 bg-blue-500/20 rounded-lg border border-blue-500">
                  <p class="text-xs text-blue-400 font-semibold mb-1">Prefijo</p>
                  <p class="text-white font-mono">{morfologia.prefijo}</p>
                </div>
              {/if}

              {#if morfologia.raiz}
                <div class="flex-1 p-3 bg-purple-500/20 rounded-lg border border-purple-500">
                  <p class="text-xs text-purple-400 font-semibold mb-1">Raíz</p>
                  <p class="text-white font-mono font-bold">{morfologia.raiz}</p>
                </div>
              {/if}

              {#if morfologia.sufijo}
                <div class="flex-1 p-3 bg-green-500/20 rounded-lg border border-green-500">
                  <p class="text-xs text-green-400 font-semibold mb-1">Sufijo</p>
                  <p class="text-white font-mono">{morfologia.sufijo}</p>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      {/if}

      <!-- Contextos de uso (carousel) -->
      {#if contextosEjemplo.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h3 class="text-lg font-semibold text-white mb-4 flex items-center gap-2">
            📖 Contextos de Uso
          </h3>

          <div class="relative">
            <!-- Contexto actual -->
            <div class="p-4 bg-canvas-950 rounded-lg border border-canvas-700 min-h-24 flex items-center">
              <p class="text-canvas-200 text-lg italic">
                "{contextosEjemplo[contextoActual]}"
              </p>
            </div>

            <!-- Controles del carousel -->
            {#if contextosEjemplo.length > 1}
              <div class="flex items-center justify-between mt-4">
                <button
                  onclick={anteriorContexto}
                  class="
                    px-4 py-2 rounded-lg font-semibold
                    bg-canvas-800 text-canvas-300
                    border border-canvas-700
                    transition-all duration-300
                    hover:bg-canvas-700 hover:scale-105
                  "
                >
                  ← Anterior
                </button>

                <div class="flex items-center gap-2">
                  {#each contextosEjemplo as _, index}
                    <button
                      onclick={() => contextoActual = index}
                      class="
                        w-2 h-2 rounded-full transition-all duration-300
                        {index === contextoActual ? 'bg-purple-500 w-6' : 'bg-canvas-700'}
                      "
                    ></button>
                  {/each}
                </div>

                <button
                  onclick={siguienteContexto}
                  class="
                    px-4 py-2 rounded-lg font-semibold
                    bg-canvas-800 text-canvas-300
                    border border-canvas-700
                    transition-all duration-300
                    hover:bg-canvas-700 hover:scale-105
                  "
                >
                  Siguiente →
                </button>
              </div>

              <p class="text-center text-xs text-canvas-500 mt-2">
                Contexto {contextoActual + 1} de {contextosEjemplo.length}
              </p>
            {/if}
          </div>
        </div>
      {/if}

      <!-- Familia de palabras -->
      {#if palabrasFamilia.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h3 class="text-lg font-semibold text-white mb-4">👨‍👩‍👧‍👦 Familia de Palabras</h3>

          <div class="flex flex-wrap gap-2">
            {#each palabrasFamilia as palabra}
              <span class="
                px-4 py-2 rounded-lg
                bg-purple-500/20 text-purple-300
                border border-purple-500/50
                font-medium
              ">
                {palabra}
              </span>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Sidebar -->
    <div class="lg:col-span-1 space-y-6">
      <!-- Sinónimos -->
      {#if sinonimos.length > 0}
        <div class="p-6 bg-green-500/10 rounded-2xl border border-green-500/30">
          <h3 class="text-sm font-semibold uppercase text-green-400 mb-3">✓ Sinónimos</h3>
          <ul class="space-y-2">
            {#each sinonimos as sinonimo}
              <li class="flex items-center gap-2">
                <span class="w-2 h-2 bg-green-500 rounded-full"></span>
                <p class="text-canvas-200">{sinonimo}</p>
              </li>
            {/each}
          </ul>
        </div>
      {/if}

      <!-- Antónimos -->
      {#if antonimos.length > 0}
        <div class="p-6 bg-red-500/10 rounded-2xl border border-red-500/30">
          <h3 class="text-sm font-semibold uppercase text-red-400 mb-3">✕ Antónimos</h3>
          <ul class="space-y-2">
            {#each antonimos as antonimo}
              <li class="flex items-center gap-2">
                <span class="w-2 h-2 bg-red-500 rounded-full"></span>
                <p class="text-canvas-200">{antonimo}</p>
              </li>
            {/each}
          </ul>
        </div>
      {/if}

      <!-- Imágenes asociativas -->
      {#if imagenes.length > 0}
        <div class="p-6 bg-canvas-900/50 rounded-2xl border border-canvas-700">
          <h3 class="text-sm font-semibold uppercase text-canvas-400 mb-3">🖼️ Imágenes</h3>
          <div class="grid grid-cols-1 gap-3">
            {#each imagenes as imagen}
              <img
                src={imagen}
                alt="Ilustración de {palabraObjetivo}"
                class="w-full h-32 object-cover rounded-lg border border-canvas-700"
              />
            {/each}
          </div>
        </div>
      {/if}

      <!-- Quiz de uso -->
      <div class="p-6 bg-purple-500/10 rounded-2xl border border-purple-500/30">
        <h3 class="text-sm font-semibold uppercase text-purple-400 mb-3">💡 Autoevaluación</h3>

        {#if !mostrarQuiz}
          <button
            onclick={() => mostrarQuiz = true}
            class="
              w-full px-4 py-3 rounded-lg font-semibold
              bg-purple-500/20 text-purple-300
              border border-purple-500/50
              transition-all duration-300
              hover:bg-purple-500/30 hover:scale-105
            "
          >
            Probar mi Comprensión
          </button>
        {:else}
          <div>
            <p class="text-canvas-300 text-sm mb-4">
              ¿Puedes crear una oración usando la palabra <span class="text-purple-400 font-semibold">{palabraObjetivo}</span>?
            </p>

            <textarea
              bind:value={respuestaQuiz}
              placeholder="Escribe tu oración aquí..."
              class="
                w-full px-4 py-3 bg-canvas-950 border-2 border-canvas-700 rounded-xl
                text-white placeholder-canvas-500
                focus:outline-none focus:border-purple-500
                transition-colors duration-300
                min-h-24 text-sm
              "
            ></textarea>

            <button
              onclick={() => {
                // En producción, aquí se enviaría al backend
                mostrarQuiz = false;
                respuestaQuiz = null;
              }}
              disabled={!respuestaQuiz?.trim()}
              class="
                w-full mt-3 px-4 py-2 rounded-lg font-semibold text-sm
                bg-green-500/20 text-green-300
                border border-green-500/50
                transition-all duration-300
                hover:bg-green-500/30 hover:scale-105
                disabled:opacity-30 disabled:cursor-not-allowed
              "
            >
              ✓ Guardar mi Oración
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>

  <!-- Navegación -->
  {#if showNavigation}
    <div class="flex items-center justify-between pt-6 mt-8 border-t border-canvas-800">
      <button
        onclick={onPrevious}
        disabled={!onPrevious}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-canvas-800 text-canvas-300
          border border-canvas-700
          transition-all duration-300
          hover:bg-canvas-700 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        ← Anterior
      </button>

      <div class="text-center">
        <p class="text-xs text-canvas-500">
          Palabra del vocabulario PAES
        </p>
      </div>

      <button
        onclick={onNext}
        disabled={!onNext}
        class="
          px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-purple-500 to-purple-600
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-purple-500/50 hover:scale-105
          disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100
        "
      >
        Siguiente →
      </button>
    </div>
  {/if}
</div>
