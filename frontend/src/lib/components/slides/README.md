# Educational Slides Components - Guía de Uso

Esta carpeta contiene los componentes de slides educativos para Lumera, diseñados para enseñar conceptos de manera interactiva antes de las actividades de evaluación.

## 📊 Taxonomía de Componentes

Lumera organiza sus componentes educativos en **3 categorías** según su propósito pedagógico:

### 1. 📚 **TEACH** (Enseñanza Expositiva)
**Objetivo:** Explicar conceptos, teoría, definiciones
**Ubicación:** `slides/teach/`
**Interactividad:** Baja-Media (navegación, tooltips, ejemplos expandibles)
**Bloom:** Recordar, Comprender
**6 componentes:** ReadingStrategy, GrammarConcept, ConnectorsGuide, VocabularyStrategy, TextTypesGuide, LiteraryDeviceGuide

### 2. ✏️ **PRACTICE** (Práctica Interactiva)
**Objetivo:** Aplicar conceptos con feedback inmediato, sin calificación
**Ubicación:** `slides/practice/`
**Interactividad:** Muy Alta (drag-drop, highlighting, construcción)
**Bloom:** Aplicar, Analizar
**6 componentes:** TextAnnotation, SentenceBuilder, ConnectorsWorkshop, VocabularyContext, TextStructure, LiteraryDevicesExplorer

### 3. 🎯 **ASSESS** (Evaluación Formal)
**Objetivo:** Medir dominio, registrar progreso, calificar
**Ubicación:** `components/activities/` (fuera de slides)
**Interactividad:** Alta (con validación estricta)
**Bloom:** Todos los niveles
**9 componentes:** MultipleChoice, TrueFalse, DragDropMatching, CriteriaEvaluation, etc.

### 4. 🔧 **GENERAL** (Slides Multiuso)
**Objetivo:** Slides que pueden usarse en cualquier contexto
**Ubicación:** `slides/general/`
**5 componentes:** ConceptIntro, ComparisonTable, StepByStepProcess, FormulaExplorer, PracticePrompt

---

## 🎓 Flujo de Aprendizaje Recomendado

```
📚 TEACH (Explicar)
   ↓
✏️ PRACTICE (Aplicar sin presión)
   ↓
🎯 ASSESS (Evaluar formalmente)
```

**Ejemplo de lección:**
1. ConnectorsGuideSlide (TEACH) → Explica los 5 tipos de conectores
2. ConnectorsWorkshopSlide (PRACTICE) → Practica eligiendo conectores
3. MultipleChoice (ASSESS) → Evaluación formal con puntaje

---

## 📚 Componentes Disponibles (18 Slides + 1 Player)

### Filosofía de Diseño:
- ✅ **Configurables vía JSON/props** (NO WYSIWYG)
- ✅ **Interactivos** con micro-interacciones
- ✅ **Animaciones GSAP** para mejor engagement
- ✅ **Navegación libre** (Anterior/Siguiente)
- ✅ **Tracking de engagement** (tiempo, clicks)
- ✅ **Categorización clara** (TEACH/PRACTICE/GENERAL)

---

## 📖 COMPONENTES GENERALES

Los siguientes componentes pueden usarse en cualquier materia:

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
- 📚 Comprensión Lectora: García Márquez (Lenguaje)
- ✍️ Gramática y Oraciones (Lenguaje)
- 📝 Escritura Argumentativa (Lenguaje)

---

## 📖 COMPONENTES DE LENGUAJE

### 7. TextAnnotationSlide.svelte
**Comprensión Lectora con Anotaciones Interactivas**

Desarrolla habilidades de lectura activa mediante highlighting multi-color y sticky notes.

**Props:**
```javascript
{
  titulo: string,
  texto: string,                    // Texto literario o informativo
  tipoLectura: string,              // "narrativa" | "argumentativa" | "expositiva" | "poética"
  preguntasGuia: Array<string>,
  vocabularioDestacado: Array<string>,
  herramientasAnotacion: Array<string>,  // ["resaltar", "notas", "subrayar"]
  materia: "lenguaje"
}
```

**Interacciones:**
- Seleccionar texto → resaltar con 4 colores (idea principal, evidencia, vocabulario, preguntas)
- Click en texto → agregar sticky note
- Panel de resumen con anotaciones organizadas

---

### 8. SentenceBuilderSlide.svelte
**Construcción Gramatical con Drag & Drop**

Enseña gramática y sintaxis mediante manipulación interactiva de palabras.

**Props:**
```javascript
{
  titulo: string,
  modoEjercicio: string,           // "free" | "guided" | "transformation"
  palabrasDisponibles: Array<{
    texto: string,
    tipo: string,                  // "sujeto" | "verbo" | "objeto" | "adjetivo" | etc.
    color: string
  }>,
  objetivoGramatical: string,
  variacionesCorrectas: Array<string>,
  mostrarScaffolding: boolean
}
```

**Interacciones:**
- Drag & drop de tiles con color coding gramatical
- Validación en tiempo real
- Showcase de múltiples construcciones válidas
- Grammar hints en tooltips

---

### 9. VocabularyContextSlide.svelte
**Vocabulario Contextual y Etimología**

Enseña vocabulario profundo con múltiples contextos y etimología.

**Props:**
```javascript
{
  palabraObjetivo: string,
  pronunciacion: string,           // IPA
  audioUrl: string | null,
  etimologia: string,
  definicion: string,
  morfologia: { prefijo, raiz, sufijo },
  contextosEjemplo: Array<string>,    // 3-5 ejemplos de uso
  sinonimos: Array<string>,
  antonimos: Array<string>,
  palabrasFamilia: Array<string>
}
```

**Interacciones:**
- Carousel de contextos de uso
- Etymology tree visual
- Quiz de autoevaluación
- Audio de pronunciación (opcional)

---

### 10. TextStructureSlide.svelte
**Estructura de Textos**

Enseña la estructura de diferentes tipos de texto (narrativo, argumentativo, expositivo, instructivo).

**Props:**
```javascript
{
  titulo: string,
  tipoTexto: string,               // "narrativo" | "argumentativo" | "expositivo" | "instructivo"
  textoEjemplo: string,
  estructura: {
    [seccion: string]: {
      texto: string,
      color: string,
      descripcion: string
    }
  },
  ejercicioTipo: string,           // "identificar-partes" | "solo-visualizar"
  comparacionTipos: boolean
}
```

**Interacciones:**
- Toggle entre vista estructura vs texto completo
- Click en sección para expandir/contraer
- Comparador de estructuras de diferentes tipos de texto

---

### 11. ConnectorsWorkshopSlide.svelte
**Conectores y Coherencia Textual**

Enseña el uso correcto de conectores (causales, adversativos, consecutivos, aditivos, temporales).

**Props:**
```javascript
{
  titulo: string,
  parrafos: Array<{
    texto: string,                 // Con "___" para espacios
    opcionesConector: Array<string>,
    correcta: string,
    tipo: string,                  // Tipo de conector
    explicacion: string
  }>,
  bancoConectores: {
    [categoria: string]: Array<string>
  },
  ejercicioCreativo: string | null
}
```

**Interacciones:**
- Seleccionar conectores para completar párrafos
- Feedback inmediato con explicación
- Banco de conectores categorizados
- Barra de progreso

---

### 12. LiteraryDevicesExplorerSlide.svelte
**Recursos Literarios Interactivos**

Enseña identificación y análisis de recursos literarios (metáfora, símil, personificación, etc.).

**Props:**
```javascript
{
  titulo: string,
  texto: string,                   // Texto literario
  autor: string,
  dispositivosLiterarios: Array<{
    tipo: string,                  // "metafora" | "simil" | "personificacion" | etc.
    ejemplos: Array<string>,
    definicion: string,
    efecto: string,
    color: string
  }>,
  preguntasAnalisis: Array<string>,
  ejercicioCreativo: string | null
}
```

**Interacciones:**
- Click en dispositivo → resalta todos los ejemplos en el texto
- Author's toolkit: gráfico de frecuencia de recursos
- Preguntas de análisis con textarea
- Ejercicio creativo de escritura

---

## 📚 COMPONENTES DE ENSEÑANZA (TEACH)

Los siguientes componentes están diseñados específicamente para **explicar conceptos** antes de la práctica:

---

### 13. ReadingStrategySlide.svelte
**Estrategias de Comprensión Lectora**

Enseña técnicas efectivas para comprender textos (identificar idea principal, inferencia, evidencia textual).

**Props:**
```javascript
{
  titulo: string,
  estrategias: Array<{
    icono: string,              // Emoji representativo
    nombre: string,             // "Identificar Idea Principal"
    resumen: string,            // Descripción breve
    explicacion: string,        // Explicación detallada
    pasos: Array<string>,       // Cómo aplicar la estrategia
    ejemplo: string,            // Ejemplo visual
    cuandoUsar: Array<string>   // Situaciones apropiadas
  }>,
  ejemploTexto: string,         // Texto de demostración
  tipsAdicionales: Array<string>
}
```

**Características:**
- Navegación por estrategias (sidebar)
- Pasos numerados para aplicación
- Ejemplos contextualizados
- Tips prácticos adicionales

---

### 14. GrammarConceptSlide.svelte
**Conceptos Gramaticales con Ejemplos**

Explica conceptos gramaticales con tabla de tipos, ejemplos y errores comunes.

**Props:**
```javascript
{
  titulo: string,
  concepto: string,             // "Sujeto", "Predicado", etc.
  definicion: string,
  tipos: Array<{
    nombre: string,             // "Sujeto Expreso"
    definicion: string,
    caracteristicas: Array<string>,
    ejemplos: Array<{
      oracion: string,
      analisis: string
    }>,
    estructura: string          // Fórmula/patrón
  }>,
  reglas: Array<{
    texto: string,
    ejemplo: string
  }>,
  erroresComunes: Array<{
    incorrecto: string,
    correcto: string,
    explicacion: string
  }>
}
```

**Características:**
- Tabs para cada tipo gramatical
- Ejemplos con análisis
- Reglas con ejemplos
- Errores comunes para evitar

---

### 15. ConnectorsGuideSlide.svelte
**Guía Completa de Conectores Textuales**

Enseña las 5 categorías de conectores con múltiples ejemplos por cada uno.

**Props:**
```javascript
{
  titulo: string,
  importancia: string,
  categorias: Array<{
    nombre: string,             // "Causales", "Adversativos", etc.
    tipo: string,               // ID interno
    definicion: string,
    funcion: string,
    conectores: Array<{
      palabra: string,          // "porque", "sin embargo"
      nivel: string,            // "básico", "formal"
      ejemplos: Array<string>,
      nota: string              // Observación adicional
    }>,
    ejemploParrafo: string,     // Párrafo que usa esa categoría
    comparacion: string         // vs otras categorías
  }>,
  consejos: Array<string>
}
```

**Características:**
- Tabs por categoría de conector
- Múltiples ejemplos por conector
- Nivel de formalidad
- Tabla resumen comparativa
- Consejos de uso

---

### 16. VocabularyStrategySlide.svelte
**Estrategias para Aprender Vocabulario**

Enseña técnicas de vocabulario: etimología, contexto, familias de palabras.

**Props:**
```javascript
{
  titulo: string,
  palabra: string,              // Palabra ejemplo (opcional)
  estrategias: Array<{
    icono: string,
    nombre: string,             // "Análisis Etimológico"
    descripcion: string,
    pasos: Array<string>,
    ejemplo: string             // Aplicado a la palabra
  }>,
  familiasPalabras: Array<{
    raiz: string,               // "graph"
    significado: string,        // "escribir"
    palabras: Array<{
      palabra: string,          // "grafía"
      definicion: string
    }>
  }>,
  consejos: Array<string>
}
```

**Características:**
- Estrategias expandibles
- Ejemplos etimológicos
- Familias de palabras
- Tips prácticos

---

### 17. TextTypesGuideSlide.svelte
**Guía de Tipos de Texto**

Explica características y estructuras de cada tipo de texto (narrativo, argumentativo, expositivo, instructivo).

**Props:**
```javascript
{
  titulo: string,
  tipos: Array<{
    nombre: string,
    tipo: string,               // ID: "narrativo", "argumentativo"
    definicion: string,
    proposito: string,
    contexto: string,
    estructura: Array<{
      nombre: string,           // "Introducción", "Tesis"
      descripcion: string,
      ejemplo: string
    }>,
    caracteristicas: Array<string>,
    ejemploTexto: string,       // Texto completo de ejemplo
    conectoresTipicos: Array<string>
  }>,
  comparacionTabla: boolean     // Mostrar tabla comparativa
}
```

**Características:**
- Tabs por tipo de texto
- Estructura paso a paso
- Características del lenguaje
- Conectores típicos
- Tabla comparativa

---

### 18. LiteraryDeviceGuideSlide.svelte
**Guía de Recursos Literarios**

Enseña dispositivos literarios en profundidad con múltiples ejemplos de autores reconocidos.

**Props:**
```javascript
{
  titulo: string,
  dispositivos: Array<{
    nombre: string,             // "metáfora", "símil"
    definicion: string,
    efecto: string,             // Efecto en el lector
    ejemplos: Array<{
      texto: string,
      autor: string,
      analisis: string
    }>,
    comoIdentificar: Array<string>,
    variantes: Array<{
      tipo: string,
      descripcion: string,
      ejemplo: string
    }>,
    cuandoUsar: string
  }>,
  comparaciones: Array<{
    dispositivos: Array<string>,
    diferencia: string,
    ejemplos: Array<{
      tipo: string,
      ejemplo: string
    }>
  }>
}
```

**Características:**
- Tabs por dispositivo literario
- Múltiples ejemplos con análisis
- Cómo identificar cada recurso
- Variantes del dispositivo
- Comparaciones entre dispositivos similares
- Tabla resumen

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

**Versión:** 3.0
**Última actualización:** 2025-11-22
**Autores:** Claude (Anthropic) + Lumera Team

**Changelog:**
- v3.0: Reorganización completa con taxonomía TEACH/PRACTICE/GENERAL. Agregados 6 componentes TEACH (ReadingStrategy, GrammarConcept, ConnectorsGuide, VocabularyStrategy, TextTypesGuide, LiteraryDeviceGuide). Lección demo "Conectores Completo" con flujo Teach → Practice. Total: 18 componentes + 1 player
- v2.0: Agregados 6 componentes de Lenguaje PRACTICE (TextAnnotation, SentenceBuilder, VocabularyContext, TextStructure, ConnectorsWorkshop, LiteraryDevicesExplorer) con 3 lecciones demo
- v1.0: Componentes generales + slides de ciencias y matemáticas
