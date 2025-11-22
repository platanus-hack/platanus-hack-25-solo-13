<script>
  import LessonPlayer from '$lib/components/slides/LessonPlayer.svelte';

  // Estado
  let leccionActiva = $state('fotosintesis');

  // Datos de lecciones de ejemplo
  const lecciones = {
    fotosintesis: {
      leccionId: "fotosintesis-intro",
      titulo: "Introducción a la Fotosíntesis",
      materia: "Biología",
      slides: [
        {
          tipo: "ConceptIntroSlide",
          orden: 1,
          props: {
            concepto: "Fotosíntesis",
            definicionSimple: "La fotosíntesis es el proceso mediante el cual las plantas convierten la luz solar en alimento (glucosa), liberando oxígeno como subproducto.",
            definicionTecnica: "6CO₂ + 6H₂O + luz → C₆H₁₂O₆ + 6O₂. Proceso anabólico donde organismos autótrofos sintetizan moléculas orgánicas complejas a partir de sustancias simples usando energía lumínica.",
            imagen: null,
            terminosClave: [
              { palabra: "autótrofos", tooltip: "Organismos que producen su propio alimento, como las plantas" },
              { palabra: "clorofila", tooltip: "Pigmento verde que captura la energía lumínica" },
              { palabra: "glucosa", tooltip: "Azúcar simple (C₆H₁₂O₆) que almacena energía química" }
            ],
            colorTema: "emerald",
            materia: "biología"
          }
        },
        {
          tipo: "StepByStepProcessSlide",
          orden: 2,
          props: {
            titulo: "Las dos fases de la Fotosíntesis",
            pasos: [
              {
                numero: 1,
                titulo: "Fase Luminosa (Reacciones de luz)",
                contenido: "Ocurre en los tilacoides del cloroplasto. La clorofila absorbe luz solar y la convierte en energía química (ATP y NADPH). Se libera oxígeno al descomponer moléculas de agua.",
                ejemplo: "H₂O + luz → O₂ + ATP + NADPH",
                ayudaVisual: null
              },
              {
                numero: 2,
                titulo: "Fase Oscura (Ciclo de Calvin)",
                contenido: "Ocurre en el estroma del cloroplasto. Usa el ATP y NADPH de la fase luminosa para convertir CO₂ en glucosa. No requiere luz directa.",
                ejemplo: "CO₂ + ATP + NADPH → C₆H₁₂O₆",
                ayudaVisual: null
              },
              {
                numero: 3,
                titulo: "Producto Final",
                contenido: "La glucosa producida se usa para crecimiento, almacenamiento de energía (almidón) o respiración celular. El oxígeno se libera a la atmósfera.",
                ejemplo: "1 molécula de glucosa = energía para 38 moléculas de ATP",
                ayudaVisual: null
              }
            ],
            materia: "biología",
            requiereConfirmacion: false,
            mostrarProgreso: true
          }
        },
        {
          tipo: "ComparisonTableSlide",
          orden: 3,
          props: {
            titulo: "Fase Luminosa vs Fase Oscura",
            items: [
              { nombre: "Fase Luminosa", color: "yellow" },
              { nombre: "Fase Oscura", color: "purple" }
            ],
            filas: [
              {
                caracteristica: "Ubicación",
                valores: ["Tilacoides", "Estroma"],
                tipo: "diferencia",
                detalles: "Los tilacoides son membranas internas del cloroplasto, el estroma es el fluido que los rodea"
              },
              {
                caracteristica: "Requiere luz",
                valores: ["Sí, directamente", "No directamente"],
                tipo: "diferencia",
                detalles: "La fase oscura usa productos de la fase luminosa, pero no necesita luz en ese momento"
              },
              {
                caracteristica: "Productos",
                valores: ["ATP, NADPH, O₂", "Glucosa (C₆H₁₂O₆)"],
                tipo: "diferencia"
              },
              {
                caracteristica: "Reactivos principales",
                valores: ["H₂O, luz", "CO₂, ATP, NADPH"],
                tipo: "diferencia"
              },
              {
                caracteristica: "Función",
                valores: ["Capturar energía", "Sintetizar glucosa"],
                tipo: "diferencia"
              }
            ],
            materia: "biología",
            mostrarFiltros: true
          }
        },
        {
          tipo: "PracticePromptSlide",
          orden: 4,
          props: {
            mensaje: "¡Ahora dominas la Fotosíntesis!",
            submensaje: "Refuerza lo aprendido con algunos ejercicios prácticos",
            icono: "🌱",
            previewEjercicios: [
              { tipo: "Verdadero/Falso", cantidad: 3, icono: "✓✗" },
              { tipo: "Selección Múltiple", cantidad: 5, icono: "☑️" },
              { tipo: "Completar Blancos", cantidad: 2, icono: "___" }
            ],
            motivacion: "¡Excelente! Ya comprendes el 70% de Biología Celular",
            botonTexto: "Comenzar Ejercicios",
            colorTema: "emerald",
            materia: "biología"
          }
        }
      ]
    },

    ecuaciones: {
      leccionId: "ecuaciones-segundo-grado",
      titulo: "Ecuaciones de Segundo Grado",
      materia: "Matemáticas",
      slides: [
        {
          tipo: "ConceptIntroSlide",
          orden: 1,
          props: {
            concepto: "Ecuación Cuadrática",
            definicionSimple: "Una ecuación de segundo grado es aquella donde la variable está elevada al cuadrado. Tiene la forma ax² + bx + c = 0.",
            definicionTecnica: "Ecuación polinómica de grado 2 de la forma ax² + bx + c = 0, donde a ≠ 0. Sus soluciones se obtienen mediante la fórmula general: x = (-b ± √(b²-4ac)) / 2a",
            imagen: null,
            terminosClave: [
              { palabra: "coeficientes", tooltip: "Los números a, b, c que multiplican a las variables" },
              { palabra: "discriminante", tooltip: "b²-4ac, determina el número de soluciones reales" },
              { palabra: "raíces", tooltip: "Los valores de x que satisfacen la ecuación" }
            ],
            colorTema: "cyan",
            materia: "matemáticas"
          }
        },
        {
          tipo: "StepByStepProcessSlide",
          orden: 2,
          props: {
            titulo: "Resolver una Ecuación Cuadrática",
            pasos: [
              {
                numero: 1,
                titulo: "Identificar coeficientes a, b, c",
                contenido: "En la forma ax² + bx + c = 0, identifica los valores de a (coeficiente de x²), b (coeficiente de x), y c (término independiente).",
                ejemplo: "2x² + 5x - 3 = 0 → a=2, b=5, c=-3"
              },
              {
                numero: 2,
                titulo: "Calcular el discriminante",
                contenido: "Calcula Δ = b² - 4ac. Este valor determina cuántas soluciones reales tiene la ecuación.",
                ejemplo: "Δ = 5² - 4(2)(-3) = 25 + 24 = 49"
              },
              {
                numero: 3,
                titulo: "Aplicar la fórmula general",
                contenido: "Sustituye en x = (-b ± √Δ) / 2a para obtener las dos soluciones.",
                ejemplo: "x = (-5 ± √49) / 4 = (-5 ± 7) / 4"
              },
              {
                numero: 4,
                titulo: "Calcular ambas soluciones",
                contenido: "Resuelve con + y con - para obtener x₁ y x₂.",
                ejemplo: "x₁ = (-5 + 7) / 4 = 0.5; x₂ = (-5 - 7) / 4 = -3"
              }
            ],
            materia: "matemáticas",
            requiereConfirmacion: true,
            mostrarProgreso: true
          }
        },
        {
          tipo: "FormulaExplorerSlide",
          orden: 3,
          props: {
            titulo: "Explora la Fórmula General",
            formula: "x = (-b ± √(b²-4ac)) / 2a",
            variables: [
              {
                simbolo: "a",
                nombre: "Coeficiente cuadrático",
                unidad: "",
                descripcion: "Coeficiente del término x² (debe ser ≠ 0)"
              },
              {
                simbolo: "b",
                nombre: "Coeficiente lineal",
                unidad: "",
                descripcion: "Coeficiente del término x"
              },
              {
                simbolo: "c",
                nombre: "Término independiente",
                unidad: "",
                descripcion: "Término constante sin variable"
              }
            ],
            ejemploNumerico: {
              a: 1,
              b: -3,
              c: 2,
              sustitucion: "x = (3 ± √(9-8)) / 2 = (3 ± 1) / 2",
              resultado: "x₁ = 2, x₂ = 1",
              unidadResultado: ""
            },
            calculadoraInteractiva: false,
            materia: "matemáticas",
            mostrarUnidades: false
          }
        },
        {
          tipo: "PracticePromptSlide",
          orden: 4,
          props: {
            mensaje: "¡Listo para resolver ecuaciones!",
            submensaje: "Practica con 10 ecuaciones de dificultad progresiva",
            icono: "🎯",
            previewEjercicios: [
              { tipo: "Ecuaciones simples", cantidad: 3, icono: "📐" },
              { tipo: "Ecuaciones intermedias", cantidad: 4, icono: "📊" },
              { tipo: "Problemas aplicados", cantidad: 3, icono: "🧮" }
            ],
            motivacion: "Dominas el 85% de Álgebra - ¡Sigue así!",
            botonTexto: "Practicar Ahora",
            colorTema: "cyan",
            materia: "matemáticas"
          }
        }
      ]
    }
  };

  // Handlers
  function handleLeccionComplete(data) {
    console.log('[LessonsDemo] Lección completada:', data);
    alert(`¡Lección completada! 🎉\n\nSlides: ${data.slidesCompletados}\nTiempo: ${data.tiempoTotal}s`);
  }

  function handleSlideChange(data) {
    console.log('[LessonsDemo] Cambio de slide:', data);
  }
</script>

<svelte:head>
  <title>Lessons Demo - Lumera</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
  <!-- Header con selector de lección -->
  <div class="sticky top-0 z-50 bg-slate-950/95 backdrop-blur-sm border-b border-slate-800">
    <div class="max-w-7xl mx-auto p-4">
      <div class="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 class="text-2xl font-bold text-white">
            Componentes de <span class="bg-gradient-to-r from-cyan-400 to-blue-500 bg-clip-text text-transparent">Lecciones</span>
          </h1>
          <p class="text-sm text-slate-400">
            Slides educativos interactivos
          </p>
        </div>

        <!-- Selector de lección -->
        <div class="flex gap-2">
          <button
            onclick={() => leccionActiva = 'fotosintesis'}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300
              {leccionActiva === 'fotosintesis' ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/50' : 'bg-slate-800 text-slate-400 hover:bg-slate-700'}
            "
          >
            🌱 Fotosíntesis
          </button>
          <button
            onclick={() => leccionActiva = 'ecuaciones'}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300
              {leccionActiva === 'ecuaciones' ? 'bg-cyan-500 text-white shadow-lg shadow-cyan-500/50' : 'bg-slate-800 text-slate-400 hover:bg-slate-700'}
            "
          >
            📐 Ecuaciones
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Player de lección -->
  {#key leccionActiva}
    <LessonPlayer
      leccion={lecciones[leccionActiva]}
      onComplete={handleLeccionComplete}
      onSlideChange={handleSlideChange}
      showProgress={true}
    />
  {/key}
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
