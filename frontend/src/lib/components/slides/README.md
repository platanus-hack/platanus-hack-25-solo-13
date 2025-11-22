# Educational Slides Components - Guía de Uso

Esta carpeta contiene los componentes de slides educativos para Lumera, diseñados para enseñar conceptos de manera interactiva antes de las actividades de evaluación.

## 📚 Componentes Disponibles (5 Slides + 1 Player)

### Filosofía de Diseño:
- ✅ **Configurables vía JSON/props** (NO WYSIWYG)
- ✅ **Interactivos** con micro-interacciones
- ✅ **Animaciones GSAP** para mejor engagement
- ✅ **Navegación libre** (Anterior/Siguiente)
- ✅ **Tracking de engagement** (tiempo, clicks)

---

## 1. ConceptIntroSlide.svelte
**Introducción de Conceptos**

Presenta un concepto nuevo con definición simple y técnica, términos clave interactivos.

**Props:**
```javascript
{
  concepto: string,              // Nombre del concepto
  definicionSimple: string,      // Versión simple para estudiantes
  definicionTecnica: string,     // Versión técnica/científica
  imagen: string | null,         // URL de imagen de apoyo
  terminosClave: Array<{         // Términos para resaltar
    palabra: string,
    tooltip: string              // Definición al hacer hover
  }>,
  colorTema: string,            // Color del tema (blue, green, etc.)
  materia: string,              // Materia
  onNext: function,             // Callback siguiente
  onPrevious: function,         // Callback anterior
  showNavigation: boolean       // Mostrar botones de navegación
}
```

**Ejemplo:**
```svelte
<ConceptIntroSlide
  concepto="Fotosíntesis"
  definicionSimple="Proceso donde plantas convierten luz en alimento"
  definicionTecnica="6CO₂ + 6H₂O + luz → C₆H₁₂O₆ + 6O₂"
  terminosClave={[
    { palabra: "clorofila", tooltip: "Pigmento verde que captura luz" }
  ]}
  colorTema="emerald"
  materia="biología"
  onNext={handleNext}
/>
```

**Interacciones:**
- Toggle entre versión simple/técnica
- Hover sobre términos clave → tooltip
- Términos resaltados en el texto

---

## 2. ComparisonTableSlide.svelte
**Tabla Comparativa**

Compara 2-3 conceptos lado a lado con filas expandibles.

**Props:**
```javascript
{
  titulo: string,
  items: Array<{                 // Conceptos a comparar
    nombre: string,
    color: string                // Color del concepto
  }>,
  filas: Array<{                 // Características a comparar
    caracteristica: string,
    valores: Array<string>,      // Valores por cada item
    tipo: string,                // "similitud" | "diferencia"
    detalles: string             // Info adicional expandible
  }>,
  materia: string,
  mostrarFiltros: boolean,       // Filtros similitudes/diferencias
  onNext: function,
  onPrevious: function,
  showNavigation: boolean
}
```

**Ejemplo:**
```svelte
<ComparisonTableSlide
  titulo="Mitosis vs Meiosis"
  items={[
    { nombre: "Mitosis", color: "cyan" },
    { nombre: "Meiosis", color: "purple" }
  ]}
  filas={[
    {
      caracteristica: "Número de divisiones",
      valores: ["1 división", "2 divisiones"],
      tipo: "diferencia",
      detalles: "Mitosis produce 2 células, Meiosis produce 4"
    }
  ]}
  materia="biología"
  onNext={handleNext}
/>
```

**Interacciones:**
- Click en fila → expande detalles
- Filtros: Todos / Similitudes / Diferencias
- Hover en filas → highlight

---

## 3. StepByStepProcessSlide.svelte
**Proceso Paso a Paso**

Enseña procesos complejos dividiéndolos en pasos secuenciales.

**Props:**
```javascript
{
  titulo: string,
  pasos: Array<{
    numero: number,
    titulo: string,
    contenido: string,
    ejemplo: string | null,      // Ejemplo numérico/textual
    ayudaVisual: string | null   // URL de imagen
  }>,
  materia: string,
  requiereConfirmacion: boolean, // Checkbox "Entendí este paso"
  mostrarProgreso: boolean,      // Barra de progreso
  onNext: function,
  onPrevious: function,
  showNavigation: boolean
}
```

**Ejemplo:**
```svelte
<StepByStepProcessSlide
  titulo="Resolver ecuación de 2do grado"
  pasos={[
    {
      numero: 1,
      titulo: "Identificar coeficientes",
      contenido: "En ax² + bx + c = 0, identifica a, b, c",
      ejemplo: "2x² + 5x - 3 = 0 → a=2, b=5, c=-3"
    }
  ]}
  requiereConfirmacion={true}
  materia="matemáticas"
  onNext={handleNext}
/>
```

**Interacciones:**
- Navegación paso a paso (Anterior/Siguiente interno)
- Checkbox "Entendí" bloquea avance (opcional)
- Click en indicadores de paso → salta a ese paso
- Progreso visual con barra

---

## 4. FormulaExplorerSlide.svelte
**Explorador de Fórmulas**

Explica fórmulas matemáticas/físicas con definición de variables y ejemplos.

**Props:**
```javascript
{
  titulo: string,
  formula: string,               // Fórmula en texto (E = mc²)
  variables: Array<{
    simbolo: string,
    nombre: string,
    unidad: string,
    descripcion: string
  }>,
  ejemploNumerico: {             // Ejemplo resuelto
    [simbolo]: number,           // Valores de las variables
    sustitucion: string,         // Paso de sustitución
    resultado: string,           // Resultado final
    unidadResultado: string
  },
  calculadoraInteractiva: boolean, // Habilitar inputs interactivos
  materia: string,
  mostrarUnidades: boolean,
  onNext: function,
  onPrevious: function,
  showNavigation: boolean
}
```

**Ejemplo:**
```svelte
<FormulaExplorerSlide
  titulo="Energía Cinética"
  formula="E = ½mv²"
  variables={[
    {
      simbolo: "E",
      nombre: "Energía cinética",
      unidad: "Joules (J)",
      descripcion: "Energía de un objeto en movimiento"
    },
    {
      simbolo: "m",
      nombre: "Masa",
      unidad: "kg",
      descripcion: "Cantidad de materia"
    }
  ]}
  ejemploNumerico={{
    m: 2,
    v: 10,
    resultado: "100 J"
  }}
  calculadoraInteractiva={true}
  materia="física"
  onNext={handleNext}
/>
```

**Interacciones:**
- Hover sobre variables → tooltip con definición
- Calculadora interactiva → inputs para probar valores
- Botón "Ver ejemplo resuelto" → expande solución
- Variables resaltadas en la fórmula

---

## 5. PracticePromptSlide.svelte
**Transición a Práctica**

Slide motivacional que conecta teoría con ejercicios prácticos.

**Props:**
```javascript
{
  mensaje: string,
  submensaje: string,
  icono: string,                 // Emoji grande
  previewEjercicios: Array<{
    tipo: string,
    cantidad: number,
    icono: string
  }>,
  motivacion: string | null,     // Mensaje motivacional opcional
  botonTexto: string,
  colorTema: string,
  materia: string,
  mostrarConfetti: boolean,      // Animación de confetti
  onNext: function,
  onPrevious: function,
  showNavigation: boolean
}
```

**Ejemplo:**
```svelte
<PracticePromptSlide
  mensaje="¡Hora de practicar!"
  submensaje="Resuelve 5 ejercicios de fotosíntesis"
  icono="🌱"
  previewEjercicios={[
    { tipo: "Verdadero/Falso", cantidad: 3, icono: "✓✗" },
    { tipo: "Selección Múltiple", cantidad: 2, icono: "☑️" }
  ]}
  motivacion="¡Vas genial! Ya dominas el 70% de la unidad"
  botonTexto="Comenzar Ejercicios"
  colorTema="emerald"
  materia="biología"
  onNext={handleNextToExercises}
/>
```

**Interacciones:**
- Animación de confetti al entrar
- Botón grande con glow effect
- Preview de tipos de ejercicios

---

## 6. LessonPlayer.svelte
**Reproductor de Lecciones**

Componente contenedor que reproduce una secuencia de slides.

**Props:**
```javascript
{
  leccion: {
    leccionId: string,
    titulo: string,
    materia: string,
    slides: Array<{
      tipo: string,              // Nombre del componente
      orden: number,
      props: object              // Props del slide
    }>
  },
  onComplete: function,          // Callback cuando termina la lección
  onSlideChange: function,       // Callback al cambiar de slide
  showProgress: boolean          // Mostrar progreso circular
}
```

**Ejemplo:**
```svelte
<script>
  const leccion = {
    leccionId: "fotosintesis-intro",
    titulo: "Introducción a la Fotosíntesis",
    materia: "Biología",
    slides: [
      {
        tipo: "ConceptIntroSlide",
        orden: 1,
        props: {
          concepto: "Fotosíntesis",
          definicionSimple: "...",
          // ... más props
        }
      },
      {
        tipo: "StepByStepProcessSlide",
        orden: 2,
        props: { /* ... */ }
      }
    ]
  };

  function handleComplete(data) {
    console.log('Lección completada:', data);
    // { leccionId, tiempoTotal, slidesCompletados, timestamp }
  }
</script>

<LessonPlayer
  {leccion}
  onComplete={handleComplete}
  showProgress={true}
/>
```

**Características:**
- Navegación automática entre slides con `{#key}`
- Progress bar circular y lineal
- Tracking de tiempo por slide
- Navegación con flechas del teclado (planeado)
- Layout responsivo

---

## 🎨 Características Comunes

### Estilos y Diseño
Todos los slides usan:
- **Dark theme gaming:** `bg-slate-950`, `border-slate-800`
- **Animaciones GSAP:** Fade-in, slide-up en entrada
- **Tailwind CSS:** Clases utilitarias responsivas
- **Navegación consistente:** Botones Anterior/Siguiente

### Colores por Materia
```javascript
const materiaColors = {
  matemáticas: 'cyan',
  lenguaje: 'purple',
  historia: 'amber',
  física: 'blue',
  química: 'green',
  biología: 'emerald'
};
```

---

## 🔌 Integración con Backend

### Estructura de Tracking

**onSlideChange:**
```javascript
{
  slideIndex: number,
  slideType: string,
  timestamp: string
}
```

**onComplete (lección terminada):**
```javascript
{
  leccionId: string,
  tiempoTotal: number,          // segundos
  slidesCompletados: number,
  timestamp: string
}
```

**Tracking interno por slide:**
```javascript
{
  leccionId: string,
  slideIndex: number,
  slideType: string,
  tiempoSegundos: number,
  interacciones: number,        // Clicks, hovers, etc.
  timestamp: string
}
```

### Endpoints Sugeridos

- `POST /api/lessons/start` - Iniciar lección
- `POST /api/lessons/slide-progress` - Trackear cada slide
- `POST /api/lessons/complete` - Completar lección

---

## 📊 Estructura de Datos de Lección

```javascript
{
  leccionId: "fotosintesis-intro",
  titulo: "Introducción a la Fotosíntesis",
  descripcion: "Aprende los conceptos básicos...",
  materia: "biología",
  nivel: "2° Medio",
  duracionEstimada: 15,         // minutos
  oaRelacionados: [123, 456],   // IDs de objetivos de aprendizaje

  slides: [
    {
      tipo: "ConceptIntroSlide",
      orden: 1,
      duracionEstimada: 3,      // minutos
      props: {
        concepto: "Fotosíntesis",
        definicionSimple: "...",
        // ... más props
      }
    },
    {
      tipo: "StepByStepProcessSlide",
      orden: 2,
      duracionEstimada: 5,
      props: { /* ... */ }
    },
    {
      tipo: "PracticePromptSlide",
      orden: 3,
      duracionEstimada: 1,
      props: { /* ... */ }
    }
  ],

  actividadesSiguientes: [       // IDs de actividades para practicar
    "activity-mc-101",
    "activity-fb-102"
  ]
}
```

---

## 🚀 Uso en el Sistema Lumera

### Flujo Completo: Lección → Práctica

```svelte
<script>
  let fase = $state('leccion'); // 'leccion' | 'practica'

  function handleLeccionComplete() {
    fase = 'practica';
  }
</script>

{#if fase === 'leccion'}
  <LessonPlayer
    leccion={leccionData}
    onComplete={handleLeccionComplete}
  />
{:else}
  <!-- Componentes de actividades -->
  <MultipleChoice {...} />
  <FillBlanks {...} />
{/if}
```

### Integración con Dashboard

```javascript
// Estado del estudiante
{
  leccionesCompletadas: ["fotosintesis-intro", "ecuaciones-2do-grado"],
  leccionEnCurso: {
    leccionId: "mitosis-meiosis",
    slideActual: 2,
    tiempoAcumulado: 180  // segundos
  }
}
```

---

## 🧪 Testing y Demo

**Ver ejemplos en vivo:**
```bash
# Iniciar frontend
make up

# Abrir en navegador
http://localhost:5173/lessons-demo
```

**Lecciones de ejemplo:**
- 🌱 Fotosíntesis (Biología)
- 📐 Ecuaciones de 2do Grado (Matemáticas)

---

## 📝 Mejoras Futuras

1. **Navegación con teclado** (flechas, ESC)
2. **Modo fullscreen**
3. **Audio narrado** por slide
4. **Subtítulos/transcripciones**
5. **Pausar/reanudar** lección
6. **Bookmarks** en slides específicos
7. **Compartir slide** específico
8. **Modo presentación** para profesores

---

**Versión:** 1.0
**Última actualización:** 2025-11-22
**Autores:** Claude (Anthropic) + Lumera Team
