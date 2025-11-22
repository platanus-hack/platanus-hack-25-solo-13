<script>
  import { onMount } from 'svelte';
  import MultipleChoice from '$lib/components/activities/MultipleChoice.svelte';
  import TrueFalse from '$lib/components/activities/TrueFalse.svelte';
  import OpenEndedResponse from '$lib/components/activities/OpenEndedResponse.svelte';
  import FillBlanks from '$lib/components/activities/FillBlanks.svelte';
  import DragDropMatching from '$lib/components/activities/DragDropMatching.svelte';
  import Sequencing from '$lib/components/activities/Sequencing.svelte';
  import CompareContrast from '$lib/components/activities/CompareContrast.svelte';
  import CriteriaEvaluation from '$lib/components/activities/CriteriaEvaluation.svelte';
  import ConceptMapBuilder from '$lib/components/activities/ConceptMapBuilder.svelte';

  // Estado para navegación de componentes
  let selectedComponent = $state('all');
  let showCode = $state(false);

  // Handlers comunes
  function handleAnswer(componentName, data) {
    console.log(`[${componentName}] Respuesta:`, data);
    // TODO: Integrar con backend POST /api/educational/progress
  }

  function handleComplete(componentName, data) {
    console.log(`[${componentName}] Completado:`, data);
    // TODO: Integrar con backend POST /api/educational/complete
  }

  function handleDraft(data) {
    console.log('[OpenEnded] Auto-save:', data);
    // TODO: Guardar borrador en localStorage o backend
  }

  // Componentes disponibles
  const components = [
    { id: 'multiple-choice', name: 'Multiple Choice', icon: '☑️' },
    { id: 'true-false', name: 'True/False', icon: '✓✗' },
    { id: 'open-ended', name: 'Open-Ended', icon: '📝' },
    { id: 'fill-blanks', name: 'Fill Blanks', icon: '___' },
    { id: 'drag-drop', name: 'Drag & Drop', icon: '🔗' },
    { id: 'sequencing', name: 'Sequencing', icon: '🔢' },
    { id: 'compare-contrast', name: 'Compare & Contrast', icon: '⚖️' },
    { id: 'criteria-evaluation', name: 'Criteria Evaluation', icon: '⭐' },
    { id: 'concept-map', name: 'Concept Map', icon: '🗺️' }
  ];
</script>

<svelte:head>
  <title>Components Demo - Lumera</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 p-8">
  <!-- Header -->
  <header class="max-w-7xl mx-auto mb-12">
    <div class="text-center">
      <h1 class="text-5xl font-bold text-white mb-4">
        Componentes Educativos
        <span class="bg-gradient-to-r from-cyan-400 to-blue-500 bg-clip-text text-transparent">
          Lumera
        </span>
      </h1>
      <p class="text-slate-400 text-lg mb-6">
        9 componentes de actividades alineados con la taxonomía de Bloom
      </p>

      <!-- Component Selector -->
      <div class="flex flex-wrap justify-center gap-3 mb-6">
        <button
          onclick={() => selectedComponent = 'all'}
          class="
            px-4 py-2 rounded-xl font-semibold transition-all duration-300
            {selectedComponent === 'all' ? 'bg-gradient-to-r from-focus-500 to-blue-500 text-white shadow-lg shadow-cyan-500/50' : 'bg-canvas-800 text-slate-400 hover:bg-slate-700'}
          "
        >
          Todos
        </button>
        {#each components as component}
          <button
            onclick={() => selectedComponent = component.id}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300
              {selectedComponent === component.id ? 'bg-gradient-to-r from-focus-500 to-blue-500 text-white shadow-lg shadow-cyan-500/50' : 'bg-canvas-800 text-slate-400 hover:bg-slate-700'}
            "
          >
            {component.icon} {component.name}
          </button>
        {/each}
      </div>

      <!-- Info Alert -->
      <div class="max-w-3xl mx-auto p-4 bg-blue-500/10 border border-blue-500/50 rounded-2xl">
        <p class="text-blue-400 text-sm">
          💡 <strong>Mock Data:</strong> Estos componentes usan datos de ejemplo.
          Los callbacks `onAnswer` y `onComplete` están conectados a la consola del navegador.
        </p>
      </div>
    </div>
  </header>

  <!-- Components Showcase -->
  <main class="max-w-7xl mx-auto space-y-16">
    <!-- 1. Multiple Choice -->
    {#if selectedComponent === 'all' || selectedComponent === 'multiple-choice'}
      <section class="scroll-mt-8" id="multiple-choice">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            ☑️ Multiple Choice Question
          </h2>
          <p class="text-slate-400">
            Pregunta de selección múltiple con 2-5 opciones. Ideal para niveles Recordar, Comprender y Aplicar.
          </p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Ejemplo 1: Historia -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Historia</p>
            <MultipleChoice
              question="¿En qué año se firmó la Declaración de Independencia de Chile?"
              options={[
                { id: 1, text: "1810", isCorrect: false },
                { id: 2, text: "1818", isCorrect: true },
                { id: 3, text: "1820", isCorrect: false },
                { id: 4, text: "1823", isCorrect: false }
              ]}
              bloomLevel="recordar"
              materia="historia"
              oaId={101}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('MultipleChoice', data)}
              onComplete={(data) => handleComplete('MultipleChoice', data)}
            />
          </div>

          <!-- Ejemplo 2: Matemáticas -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Matemáticas</p>
            <MultipleChoice
              question="¿Cuál es el resultado de 2³ + 5²?"
              options={[
                { id: 1, text: "13", isCorrect: false },
                { id: 2, text: "23", isCorrect: false },
                { id: 3, text: "33", isCorrect: true },
                { id: 4, text: "43", isCorrect: false }
              ]}
              bloomLevel="aplicar"
              materia="matemáticas"
              oaId={102}
              onAnswer={(data) => handleAnswer('MultipleChoice', data)}
              onComplete={(data) => handleComplete('MultipleChoice', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 2. True/False -->
    {#if selectedComponent === 'all' || selectedComponent === 'true-false'}
      <section class="scroll-mt-8" id="true-false">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            ✓✗ True/False Statement
          </h2>
          <p class="text-slate-400">
            Evalúa afirmaciones como verdaderas o falsas. Rápido y efectivo para diagnósticos.
          </p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Ejemplo 1: Biología -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Biología</p>
            <TrueFalse
              statement="La fotosíntesis produce oxígeno como subproducto."
              correctAnswer={true}
              explanation="Correcto. Durante la fotosíntesis, las plantas liberan O₂ al descomponer H₂O en la fase luminosa."
              bloomLevel="comprender"
              materia="biología"
              oaId={201}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('TrueFalse', data)}
              onComplete={(data) => handleComplete('TrueFalse', data)}
            />
          </div>

          <!-- Ejemplo 2: Historia con justificación -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Historia (con justificación)</p>
            <TrueFalse
              statement="La Guerra del Pacífico ocurrió entre Chile y Argentina."
              correctAnswer={false}
              explanation="Falso. La Guerra del Pacífico (1879-1884) fue entre Chile contra Perú y Bolivia."
              bloomLevel="recordar"
              materia="historia"
              oaId={202}
              requireJustification={true}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('TrueFalse', data)}
              onComplete={(data) => handleComplete('TrueFalse', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 3. Open-Ended Response -->
    {#if selectedComponent === 'all' || selectedComponent === 'open-ended'}
      <section class="scroll-mt-8" id="open-ended">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            📝 Open-Ended Response
          </h2>
          <p class="text-slate-400">
            Respuesta de texto libre para análisis, evaluación y creación. Niveles altos de Bloom.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-6">
          <!-- Ejemplo: Historia -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Historia (con rúbrica y AI feedback)</p>
            <OpenEndedResponse
              prompt="Analiza cómo la Guerra del Pacífico impactó el desarrollo económico de Chile en el siglo XIX y principios del XX."
              minWords={100}
              maxWords={300}
              bloomLevel="analizar"
              materia="historia"
              oaId={301}
              rubric={[
                "Menciona al menos 3 consecuencias económicas específicas",
                "Incluye evidencia histórica con fechas y datos concretos",
                "Analiza causas y efectos de forma clara y coherente",
                "Conecta el impacto económico con el contexto social y político"
              ]}
              enableAiFeedback={true}
              onSubmit={(data) => handleAnswer('OpenEnded', data)}
              onComplete={(data) => handleComplete('OpenEnded', data)}
              onDraft={handleDraft}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 4. Fill in the Blanks -->
    {#if selectedComponent === 'all' || selectedComponent === 'fill-blanks'}
      <section class="scroll-mt-8" id="fill-blanks">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            ___ Fill in the Blanks
          </h2>
          <p class="text-slate-400">
            Completa espacios en blanco en textos. Ideal para vocabulario y conceptos clave.
          </p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Ejemplo 1: Sin banco de palabras -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Sin banco de palabras</p>
            <FillBlanks
              text="La ___1___ es el proceso mediante el cual las plantas convierten la luz solar en ___2___ química, liberando ___3___ como subproducto."
              blanks={[
                { id: 1, answer: "fotosíntesis", caseSensitive: false },
                { id: 2, answer: "energía", caseSensitive: false },
                { id: 3, answer: "oxígeno", caseSensitive: false }
              ]}
              bloomLevel="recordar"
              materia="biología"
              oaId={401}
              allowMultipleAttempts={true}
              showHints={true}
              onAnswer={(data) => handleAnswer('FillBlanks', data)}
              onComplete={(data) => handleComplete('FillBlanks', data)}
            />
          </div>

          <!-- Ejemplo 2: Con banco de palabras -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Con banco de palabras</p>
            <FillBlanks
              text="La capital de Chile es ___1___ y está ubicada en la región ___2___. La segunda ciudad más poblada es ___3___."
              blanks={[
                { id: 1, answer: "Santiago", caseSensitive: false },
                { id: 2, answer: "Metropolitana", caseSensitive: false },
                { id: 3, answer: "Puente Alto", caseSensitive: false }
              ]}
              showWordBank={true}
              wordBank={["Santiago", "Metropolitana", "Puente Alto", "Valparaíso", "Concepción", "Araucanía"]}
              bloomLevel="recordar"
              materia="historia"
              oaId={402}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('FillBlanks', data)}
              onComplete={(data) => handleComplete('FillBlanks', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 5. Drag & Drop Matching -->
    {#if selectedComponent === 'all' || selectedComponent === 'drag-drop'}
      <section class="scroll-mt-8" id="drag-drop">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            🔗 Drag & Drop Matching
          </h2>
          <p class="text-slate-400">
            Relaciona términos con definiciones arrastrando y soltando. Interactivo y visual.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-6">
          <!-- Ejemplo: Biología -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Biología Celular</p>
            <DragDropMatching
              title="Relaciona los procesos celulares con sus definiciones"
              pairs={[
                {
                  id: 1,
                  term: "Fotosíntesis",
                  definition: "Proceso de conversión de luz solar en energía química almacenada en glucosa"
                },
                {
                  id: 2,
                  term: "Respiración Celular",
                  definition: "Proceso de obtención de energía (ATP) a partir de la descomposición de glucosa"
                },
                {
                  id: 3,
                  term: "Mitosis",
                  definition: "División celular que produce dos células hijas genéticamente idénticas"
                },
                {
                  id: 4,
                  term: "Meiosis",
                  definition: "División celular que produce cuatro células hijas con la mitad de cromosomas"
                }
              ]}
              bloomLevel="comprender"
              materia="biología"
              oaId={501}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('DragDrop', data)}
              onComplete={(data) => handleComplete('DragDrop', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 6. Sequencing -->
    {#if selectedComponent === 'all' || selectedComponent === 'sequencing'}
      <section class="scroll-mt-8" id="sequencing">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            🔢 Sequencing/Ordering
          </h2>
          <p class="text-slate-400">
            Ordena elementos en la secuencia correcta. Ideal para cronología y procesos.
          </p>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Ejemplo 1: Historia -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Independencia de Chile</p>
            <Sequencing
              title="Ordena los eventos de la Independencia de Chile cronológicamente"
              items={[
                { id: 1, content: "Primera Junta Nacional de Gobierno (1810)", correctOrder: 1 },
                { id: 2, content: "Batalla de Rancagua - Desastre de Rancagua (1814)", correctOrder: 2 },
                { id: 3, content: "Cruce de los Andes por el Ejército Libertador (1817)", correctOrder: 3 },
                { id: 4, content: "Batalla de Chacabuco (1817)", correctOrder: 4 },
                { id: 5, content: "Declaración de Independencia (1818)", correctOrder: 5 },
                { id: 6, content: "Batalla de Maipú - Victoria definitiva (1818)", correctOrder: 6 }
              ]}
              bloomLevel="comprender"
              materia="historia"
              oaId={601}
              showHints={true}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('Sequencing', data)}
              onComplete={(data) => handleComplete('Sequencing', data)}
            />
          </div>

          <!-- Ejemplo 2: Biología -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Fases de la Mitosis</p>
            <Sequencing
              title="Ordena las fases de la mitosis en el orden correcto"
              items={[
                { id: 1, content: "Interfase - Duplicación del ADN", correctOrder: 1 },
                { id: 2, content: "Profase - Condensación de cromosomas", correctOrder: 2 },
                { id: 3, content: "Metafase - Alineación en el ecuador celular", correctOrder: 3 },
                { id: 4, content: "Anafase - Separación de cromátidas", correctOrder: 4 },
                { id: 5, content: "Telofase - Formación de núcleos hijos", correctOrder: 5 },
                { id: 6, content: "Citocinesis - División del citoplasma", correctOrder: 6 }
              ]}
              bloomLevel="comprender"
              materia="biología"
              oaId={602}
              showHints={true}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('Sequencing', data)}
              onComplete={(data) => handleComplete('Sequencing', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 7. Compare & Contrast -->
    {#if selectedComponent === 'all' || selectedComponent === 'compare-contrast'}
      <section class="scroll-mt-8" id="compare-contrast">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            ⚖️ Compare & Contrast
          </h2>
          <p class="text-slate-400">
            Analiza similitudes y diferencias entre dos conceptos. Nivel Bloom: Analizar.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-6">
          <!-- Ejemplo: Biología -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Células Animales vs Vegetales</p>
            <CompareContrast
              title="Compara las características de células animales y vegetales"
              itemA={{ name: "Célula Animal", color: "cyan" }}
              itemB={{ name: "Célula Vegetal", color: "green" }}
              characteristics={[
                { id: 1, text: "Tiene pared celular", correctColumn: "B" },
                { id: 2, text: "Tiene membrana celular", correctColumn: "both" },
                { id: 3, text: "Tiene cloroplastos", correctColumn: "B" },
                { id: 4, text: "Tiene centriolos bien definidos", correctColumn: "A" },
                { id: 5, text: "Tiene mitocondrias", correctColumn: "both" },
                { id: 6, text: "Tiene vacuola central grande", correctColumn: "B" },
                { id: 7, text: "Tiene núcleo", correctColumn: "both" },
                { id: 8, text: "Forma irregular o redondeada", correctColumn: "A" }
              ]}
              bloomLevel="analizar"
              materia="biología"
              oaId={701}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('CompareContrast', data)}
              onComplete={(data) => handleComplete('CompareContrast', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 8. Criteria Evaluation -->
    {#if selectedComponent === 'all' || selectedComponent === 'criteria-evaluation'}
      <section class="scroll-mt-8" id="criteria-evaluation">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            ⭐ Criteria Evaluation
          </h2>
          <p class="text-slate-400">
            Evalúa la calidad de un argumento o fuente usando criterios específicos. Nivel Bloom: Evaluar.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-6">
          <!-- Ejemplo: Historia -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Evaluación de Fuente Histórica</p>
            <CriteriaEvaluation
              title="Evalúa la calidad de este documento histórico"
              subject="Artículo: 'Consecuencias de la Guerra del Pacífico'"
              description="Un artículo de periódico de 1885 sobre el impacto económico de la guerra."
              content="La victoria en la Guerra del Pacífico (1879-1884) transformó a Chile en una potencia regional. La anexión de las provincias de Tarapacá y Antofagasta proporcionó acceso exclusivo a ricos yacimientos de salitre y cobre. Los ingresos del salitre representaron hasta el 50% del presupuesto nacional en las décadas siguientes, financiando modernización de puertos, ferrocarriles y educación pública."
              criteria={[
                {
                  id: 1,
                  name: "Evidencia histórica",
                  description: "¿Menciona datos, fechas y hechos específicos verificables?",
                  expectedRating: 5,
                  weight: 30
                },
                {
                  id: 2,
                  name: "Objetividad",
                  description: "¿Presenta los hechos sin sesgo nacionalista evidente?",
                  expectedRating: 3,
                  weight: 25
                },
                {
                  id: 3,
                  name: "Contexto temporal",
                  description: "¿Es apropiado considerando que fue escrito solo 1 año después de la guerra?",
                  expectedRating: 4,
                  weight: 20
                },
                {
                  id: 4,
                  name: "Análisis económico",
                  description: "¿Proporciona datos económicos concretos y medibles?",
                  expectedRating: 4,
                  weight: 25
                }
              ]}
              bloomLevel="evaluar"
              materia="historia"
              oaId={801}
              showExpectedRatings={true}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('CriteriaEvaluation', data)}
              onComplete={(data) => handleComplete('CriteriaEvaluation', data)}
            />
          </div>
        </div>
      </section>
    {/if}

    <!-- 9. Concept Map Builder -->
    {#if selectedComponent === 'all' || selectedComponent === 'concept-map'}
      <section class="scroll-mt-8" id="concept-map">
        <div class="mb-6">
          <h2 class="text-3xl font-bold text-white mb-2">
            🗺️ Concept Map Builder
          </h2>
          <p class="text-slate-400">
            Crea un mapa conceptual conectando ideas y relaciones. Nivel Bloom: Crear.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-6">
          <!-- Ejemplo: Biología -->
          <div>
            <p class="text-xs font-semibold text-slate-500 mb-3 uppercase">Ejemplo: Fotosíntesis</p>
            <ConceptMapBuilder
              title="Crea un mapa conceptual"
              topic="Fotosíntesis"
              instructions="Identifica los conceptos clave del proceso de fotosíntesis y cómo se relacionan entre sí"
              requiredConcepts={["Fotosíntesis", "Luz Solar", "Clorofila", "Agua", "CO2", "Oxígeno", "Glucosa"]}
              suggestedConnections={[
                { from: "Luz Solar", to: "Fotosíntesis", label: "inicia" },
                { from: "Clorofila", to: "Fotosíntesis", label: "captura luz para" },
                { from: "Agua", to: "Fotosíntesis", label: "es reactivo de" },
                { from: "CO2", to: "Fotosíntesis", label: "es reactivo de" },
                { from: "Fotosíntesis", to: "Oxígeno", label: "produce" },
                { from: "Fotosíntesis", to: "Glucosa", label: "produce" }
              ]}
              minConcepts={5}
              minConnections={4}
              bloomLevel="crear"
              materia="biología"
              oaId={901}
              allowMultipleAttempts={true}
              onAnswer={(data) => handleAnswer('ConceptMapBuilder', data)}
              onComplete={(data) => handleComplete('ConceptMapBuilder', data)}
            />
          </div>
        </div>
      </section>
    {/if}
  </main>

  <!-- Footer -->
  <footer class="max-w-7xl mx-auto mt-20 pt-8 border-t border-slate-800">
    <div class="text-center">
      <p class="text-slate-500 text-sm mb-4">
        Componentes educativos desarrollados para Lumera - Platanus Hack 25
      </p>
      <div class="flex justify-center gap-6 text-xs text-slate-600">
        <a href="/components-demo" class="hover:text-focus-400 transition-colors">Demo</a>
        <span>|</span>
        <a href="https://github.com" class="hover:text-focus-400 transition-colors">GitHub</a>
        <span>|</span>
        <a href="/docs" class="hover:text-focus-400 transition-colors">Documentación</a>
      </div>
    </div>
  </footer>
</div>

<style>
  /* Estilos adicionales si es necesario */
  :global(body) {
    overflow-x: hidden;
  }
</style>
