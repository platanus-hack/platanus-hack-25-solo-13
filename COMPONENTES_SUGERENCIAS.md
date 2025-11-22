# Tipos de Componentes Educativos para Lumera

**Fecha de creación:** 2025-11-22
**Contexto:** Lumera es una plataforma de aprendizaje adaptativo para estudiantes chilenos (1° a 4° Medio) preparándose para la PAES, con gamificación estilo RPG y alineación con taxonomía de Bloom.

---

## 📚 A. COMPONENTES DE ACTIVIDAD DE APRENDIZAJE (Core Learning)

Componentes que los estudiantes usan directamente para aprender y practicar:

### 1. ✅ Multiple Choice Question (MCQ) - IMPLEMENTADO
- **Bloom levels:** Recordar, Comprender, Aplicar
- **Variantes:** 2-5 opciones, con/sin imágenes, feedback inmediato
- **Uso:** Evaluaciones rápidas, PAES simulator
- **Interactividad:** Alta
- **Tiempo estimado:** 1-3 minutos
- **Archivo:** `MultipleChoice.svelte`

### 2. ✅ True/False Statement - IMPLEMENTADO
- **Bloom levels:** Recordar, Comprender
- **Variantes:** Simple, con justificación opcional
- **Uso:** Daily missions, warm-ups, diagnósticos rápidos
- **Interactividad:** Media
- **Tiempo estimado:** 30 segundos - 1 minuto
- **Archivo:** `TrueFalse.svelte`

### 3. ✅ Open-Ended Text Response - IMPLEMENTADO
- **Bloom levels:** Analizar, Evaluar, Crear
- **Variantes:** Con/sin límite de palabras, puede incluir AI feedback
- **Uso:** Ensayos cortos, reflexiones, análisis crítico
- **Interactividad:** Alta
- **Tiempo estimado:** 5-15 minutos
- **Archivo:** `OpenEndedResponse.svelte`

### 4. ✅ Fill in the Blanks / Cloze Test - IMPLEMENTADO
- **Bloom levels:** Recordar, Comprender
- **Variantes:** Palabras únicas, frases, selección de banco de palabras
- **Uso:** Vocabulario, fórmulas matemáticas, conceptos clave
- **Interactividad:** Media-Alta
- **Tiempo estimado:** 2-5 minutos
- **Archivo:** `FillBlanks.svelte`

### 5. ✅ Drag & Drop Matching - IMPLEMENTADO
- **Bloom levels:** Comprender, Aplicar
- **Variantes:** Términos-definiciones, causas-efectos, imágenes-textos
- **Uso:** Historia (eventos-fechas), Ciencias (conceptos-aplicaciones)
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 3-7 minutos
- **Archivo:** `DragDropMatching.svelte`

### 6. ✅ Sequencing/Ordering - IMPLEMENTADO
- **Bloom levels:** Comprender, Aplicar
- **Variantes:** Ordenar pasos, eventos cronológicos, procesos lógicos
- **Uso:** Historia (cronología), procedimientos científicos, algoritmos
- **Interactividad:** Alta
- **Tiempo estimado:** 2-5 minutos
- **Archivo:** `Sequencing.svelte`

### 6b. ✅ Compare & Contrast - IMPLEMENTADO (NUEVO)
- **Bloom levels:** Analizar
- **Variantes:** Tabla de 3 columnas (A | Ambos | B), drag & drop características
- **Uso:** Análisis comparativo, pensamiento crítico, ciencias/historia
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 5-8 minutos
- **Archivo:** `CompareContrast.svelte`

### 6c. ✅ Criteria Evaluation - IMPLEMENTADO (NUEVO)
- **Bloom levels:** Evaluar
- **Variantes:** Rúbrica interactiva con escala de estrellas (1-5), evaluación ponderada
- **Uso:** Evaluar argumentos, fuentes históricas, calidad de trabajos
- **Interactividad:** Alta
- **Tiempo estimado:** 5-10 minutos
- **Archivo:** `CriteriaEvaluation.svelte`

### 7. Interactive Diagram/Labeling
- **Bloom levels:** Recordar, Comprender, Aplicar
- **Variantes:** Anatomía, mapas geográficos, diagramas moleculares
- **Uso:** Biología, Geografía, Química, Física
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 3-8 minutos

### 8. Math Equation Solver
- **Bloom levels:** Aplicar, Analizar
- **Variantes:** Editor de ecuaciones con paso a paso, verificador
- **Uso:** Matemáticas (álgebra, cálculo, trigonometría)
- **Interactividad:** Alta
- **Tiempo estimado:** 5-10 minutos

### 9. Flashcard Stack
- **Bloom levels:** Recordar
- **Variantes:** Spaced repetition, flip animations, auto-avance
- **Uso:** Vocabulario, fórmulas, fechas históricas, definiciones
- **Interactividad:** Media
- **Tiempo estimado:** 5-10 minutos (stacks de 10-20 cards)

### 10. Code/Formula Editor
- **Bloom levels:** Aplicar, Crear
- **Variantes:** Syntax highlighting, autocomplete
- **Uso:** Matemáticas avanzadas, física, lógica computacional
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 10-20 minutos

---

## 📊 B. COMPONENTES DE SLIDES EDUCATIVOS (Lesson Slides) - ✅ IMPLEMENTADO

Slides interactivos configurables vía JSON para enseñar conceptos (NO WYSIWYG):

### 11. ✅ ConceptIntroSlide - IMPLEMENTADO
- **Bloom levels:** Recordar, Comprender
- **Variantes:** Toggle simple/técnica, términos clave con tooltips, imagen de apoyo
- **Uso:** Introducir conceptos nuevos, definiciones, terminología
- **Interactividad:** Alta
- **Tiempo estimado:** 2-4 minutos
- **Archivo:** `ConceptIntroSlide.svelte`

### 12. ✅ ComparisonTableSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Analizar
- **Variantes:** 2-3 conceptos, filas expandibles, filtros similitudes/diferencias
- **Uso:** Comparar conceptos (Mitosis vs Meiosis, Capitalismo vs Socialismo)
- **Interactividad:** Alta
- **Tiempo estimado:** 3-5 minutos
- **Archivo:** `ComparisonTableSlide.svelte`

### 13. ✅ StepByStepProcessSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Aplicar
- **Variantes:** Navegación paso a paso, checkboxes de confirmación, progreso visual
- **Uso:** Procesos secuenciales (resolver ecuación, método científico, experimentos)
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 5-8 minutos
- **Archivo:** `StepByStepProcessSlide.svelte`

### 14. ✅ FormulaExplorerSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Aplicar
- **Variantes:** Variables con tooltips, calculadora interactiva, ejemplo resuelto
- **Uso:** Fórmulas matemáticas/físicas (E=mc², ecuaciones, leyes)
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 4-6 minutos
- **Archivo:** `FormulaExplorerSlide.svelte`

### 15. ✅ PracticePromptSlide - IMPLEMENTADO
- **Bloom levels:** N/A (transición)
- **Variantes:** Preview de ejercicios, mensaje motivacional, confetti animado
- **Uso:** Transición entre teoría y práctica, motivar estudiante
- **Interactividad:** Media
- **Tiempo estimado:** 1-2 minutos
- **Archivo:** `PracticePromptSlide.svelte`

### 16. ✅ LessonPlayer - IMPLEMENTADO (Contenedor)
- **Bloom levels:** N/A (sistema)
- **Variantes:** Progress circular/lineal, tracking de tiempo, navegación por teclado
- **Uso:** Reproductor de secuencias de slides configurables
- **Interactividad:** Sistema
- **Tiempo estimado:** N/A
- **Archivo:** `LessonPlayer.svelte`

### B.1 📖 SLIDES DE LENGUAJE (Nuevos) - ✅ IMPLEMENTADO

Componentes especializados para comprensión lectora, gramática, vocabulario y literatura:

### 17. ✅ TextAnnotationSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Analizar
- **Variantes:** 4 colores de highlighting (idea principal, evidencia, vocabulario, preguntas), sticky notes
- **Uso:** Comprensión lectora activa, análisis de textos narrativos/argumentativos/expositivos
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 5-10 minutos
- **Archivo:** `TextAnnotationSlide.svelte`

### 18. ✅ SentenceBuilderSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Aplicar
- **Variantes:** 3 modos (libre, guiado, transformación), drag & drop palabras color-coded
- **Uso:** Gramática, sintaxis, estructura de oraciones
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 4-7 minutos
- **Archivo:** `SentenceBuilderSlide.svelte`

### 19. ✅ VocabularyContextSlide - IMPLEMENTADO
- **Bloom levels:** Recordar, Comprender, Aplicar
- **Variantes:** Etimología (prefijo + raíz + sufijo), 3-5 contextos con carousel, sinónimos/antónimos
- **Uso:** Vocabulario académico, etimología, familias de palabras
- **Interactividad:** Alta
- **Tiempo estimado:** 3-5 minutos
- **Archivo:** `VocabularyContextSlide.svelte`

### 20. ✅ TextStructureSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Analizar
- **Variantes:** 4 tipos (narrativo, argumentativo, expositivo, instructivo), toggle estructura/texto
- **Uso:** Tipos de texto PAES, estructura argumentativa, análisis de géneros
- **Interactividad:** Alta
- **Tiempo estimado:** 4-6 minutos
- **Archivo:** `TextStructureSlide.svelte`

### 21. ✅ ConnectorsWorkshopSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Aplicar
- **Variantes:** 5 categorías (causales, adversativos, consecutivos, aditivos, temporales), fill-in-the-blank
- **Uso:** Coherencia textual, conectores lógicos, escritura académica
- **Interactividad:** Alta
- **Tiempo estimado:** 5-8 minutos
- **Archivo:** `ConnectorsWorkshopSlide.svelte`

### 22. ✅ LiteraryDevicesExplorerSlide - IMPLEMENTADO
- **Bloom levels:** Comprender, Analizar, Evaluar
- **Variantes:** 8 dispositivos (metáfora, símil, personificación, etc.), frecuency chart, análisis crítico
- **Uso:** Literatura, análisis literario, recursos estilísticos
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 6-10 minutos
- **Archivo:** `LiteraryDevicesExplorerSlide.svelte`

**Demo disponible:** `http://localhost:5173/lessons-demo`

---

## 📖 C. COMPONENTES DE CONTENIDO ENRIQUECIDO (Content Delivery)

Para presentar información de manera engaging:

### 17. Reading Comprehension Passage
- **Bloom levels:** Comprender, Analizar
- **Variantes:** Texto + preguntas integradas (MCQ, open-ended)
- **Uso:** Lenguaje, Historia, comprensión lectora PAES
- **Interactividad:** Media-Alta
- **Tiempo estimado:** 10-20 minutos

### 18. Interactive Video Player
- **Bloom levels:** Todos
- **Variantes:** Video con pausas programadas para preguntas, anotaciones
- **Uso:** Clases grabadas, tutoriales, demostraciones científicas
- **Interactividad:** Alta
- **Tiempo estimado:** 5-30 minutos

### 19. Audio Listening Exercise
- **Bloom levels:** Comprender, Analizar
- **Variantes:** Audio + transcripción opcional + preguntas
- **Uso:** Lenguaje (comprensión auditiva), Historia (testimonios)
- **Interactividad:** Media
- **Tiempo estimado:** 5-15 minutos

### 20. Interactive Timeline
- **Bloom levels:** Recordar, Comprender
- **Variantes:** Línea de tiempo clickeable con eventos expandibles
- **Uso:** Historia, Literatura (biografías), procesos históricos
- **Interactividad:** Alta
- **Tiempo estimado:** 5-10 minutos

### 21. ✅ Concept Map/Mind Map - IMPLEMENTADO
- **Bloom levels:** Comprender, Analizar, Crear
- **Variantes:** Constructor interactivo con nodos draggables, conexiones con etiquetas, canvas SVG
- **Uso:** Todas las materias (resúmenes, conexiones conceptuales, síntesis de conocimiento)
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 10-20 minutos
- **Archivo:** `ConceptMapBuilder.svelte`

### 22. Infographic Display
- **Bloom levels:** Comprender
- **Variantes:** Información visual con hotspots interactivos
- **Uso:** Estadísticas, datos científicos, geografía
- **Interactividad:** Media-Alta
- **Tiempo estimado:** 3-7 minutos

---

## ✅ C. COMPONENTES DE EVALUACIÓN Y FEEDBACK (Assessment)

Para medir progreso y dar retroalimentación:

### 17. Mini Quiz (Bundle)
- **Bloom levels:** Mixto
- **Variantes:** Set de 5-10 preguntas con timer, sin timer
- **Uso:** Daily missions, diagnósticos de unidad
- **Interactividad:** Alta
- **Tiempo estimado:** 5-15 minutos

### 18. PAES Practice Section
- **Bloom levels:** Todos
- **Variantes:** Simulación de sección de PAES (timer, formato oficial)
- **Uso:** PAES Simulator, práctica completa
- **Interactividad:** Alta
- **Tiempo estimado:** 30-90 minutos

### 19. Answer Review Panel
- **Bloom levels:** N/A (post-assessment)
- **Variantes:** Respuestas correctas/incorrectas con explicaciones detalladas
- **Uso:** Después de quizzes/exams, revisión de errores
- **Interactividad:** Media
- **Tiempo estimado:** 5-15 minutos

### 20. Progress Check-in
- **Bloom levels:** Evaluar (metacognición)
- **Variantes:** Auto-evaluación de confianza (emoji scale, slider)
- **Uso:** Inicio/fin de unidades, reflexión personal
- **Interactividad:** Baja-Media
- **Tiempo estimado:** 1-2 minutos

### 21. Peer Review Widget
- **Bloom levels:** Evaluar
- **Variantes:** Rúbricas, comentarios, calificación por pares
- **Uso:** Proyectos colaborativos, ensayos
- **Interactividad:** Alta
- **Tiempo estimado:** 10-20 minutos

---

## 📊 D. COMPONENTES DE PROGRESO Y GAMIFICACIÓN (Progress Tracking)

Para motivar y visualizar avances:

### 22. Bloom Level Progress Wheel
- **Bloom levels:** N/A (visualización)
- **Variantes:** Círculo de 6 secciones (colores Bloom), animado
- **Uso:** Dashboard por OA, visualización de mastery
- **Interactividad:** Baja (visual)
- **Tiempo estimado:** N/A (siempre visible)

### 23. Subject Mastery Dashboard
- **Bloom levels:** N/A (visualización)
- **Variantes:** Grid de materias con % completado, colores por materia
- **Uso:** Home dashboard, perfil de estudiante
- **Interactividad:** Baja (clickeable)
- **Tiempo estimado:** N/A

### 24. Learning Path Map / Skill Tree
- **Bloom levels:** N/A (navegación)
- **Variantes:** Mapa visual de OAs (bloqueados/desbloqueados), estilo RPG
- **Uso:** "Mi Currículum", navegación de unidades
- **Interactividad:** Media (navegación)
- **Tiempo estimado:** N/A

### 25. XP Progress Bar
- **Bloom levels:** N/A (gamificación)
- **Variantes:** Barra animada con niveles, iconos de logros
- **Uso:** Header de dashboard (ya existe en Lumera)
- **Interactividad:** Baja (visual)
- **Tiempo estimado:** N/A

### 26. Streak Calendar
- **Bloom levels:** N/A (gamificación)
- **Variantes:** Calendario visual con días completados, estilo GitHub
- **Uso:** Dashboard, perfil, motivación diaria
- **Interactividad:** Baja (visual)
- **Tiempo estimado:** N/A

### 27. Achievement Badge Display
- **Bloom levels:** N/A (gamificación)
- **Variantes:** Showcase de badges/trophies ganados, categorías
- **Uso:** Colecciones, perfil público
- **Interactividad:** Media (clickeable para detalles)
- **Tiempo estimado:** N/A

### 28. Leaderboard
- **Bloom levels:** N/A (social)
- **Variantes:** Top estudiantes (XP, PAES score), filtros por materia
- **Uso:** Social features, competencia amistosa
- **Interactividad:** Baja (scroll)
- **Tiempo estimado:** N/A

### 29. PAES Readiness Meter
- **Bloom levels:** N/A (visualización)
- **Variantes:** Gauge/medidor predictivo de score PAES
- **Uso:** Dashboard principal, motivación
- **Interactividad:** Baja (visual)
- **Tiempo estimado:** N/A

---

## 🤖 E. COMPONENTES ADAPTATIVOS E IA (AI-Powered)

Componentes que usan el perfil del estudiante:

### 30. Personalized Recommendation Card
- **Bloom levels:** N/A (recomendación)
- **Variantes:** "Basado en tu perfil, te recomendamos..." con razón
- **Uso:** Home feed, sugerencias de OAs según weak spots
- **Interactividad:** Media (clickeable)
- **Tiempo estimado:** N/A

### 31. Difficulty Adjuster
- **Bloom levels:** N/A (herramienta)
- **Variantes:** Botón "Hacerlo más fácil/difícil", adapta complejidad
- **Uso:** Dentro de actividades, ajuste en tiempo real
- **Interactividad:** Alta
- **Tiempo estimado:** N/A

### 32. Weak Spot Identifier
- **Bloom levels:** N/A (análisis)
- **Variantes:** Alerta visual de OAs que necesitan refuerzo
- **Uso:** Dashboard, adaptive gym, diagnósticos
- **Interactividad:** Media (clickeable)
- **Tiempo estimado:** N/A

### 33. Learning Style Matcher Badge
- **Bloom levels:** N/A (metadata)
- **Variantes:** Icono "Esta actividad es visual ✓" (match con perfil)
- **Uso:** Mission cards, filtros de actividades
- **Interactividad:** Baja (visual)
- **Tiempo estimado:** N/A

### 34. AI Tutor Chatbot
- **Bloom levels:** Comprender, Analizar
- **Variantes:** Asistente conversacional para dudas, explicaciones adaptativas
- **Uso:** Help modal, dudas durante actividades
- **Interactividad:** Muy Alta
- **Tiempo estimado:** Variable (1-10 minutos)

### 35. Explanation Generator
- **Bloom levels:** N/A (herramienta)
- **Variantes:** Botón "Explícame de otra forma" (regenera con IA)
- **Uso:** Feedback panels, respuestas incorrectas
- **Interactividad:** Alta
- **Tiempo estimado:** N/A

---

## 👥 F. COMPONENTES SOCIALES Y COLABORATIVOS (Social Learning)

Para fomentar estudio en grupo:

### 36. Study Group Card
- **Bloom levels:** N/A (social)
- **Variantes:** Card con info de grupo + botón "Unirse"
- **Uso:** Friends tab, creación de grupos
- **Interactividad:** Media (clickeable)
- **Tiempo estimado:** N/A

### 37. Challenge Card
- **Bloom levels:** Variable
- **Variantes:** Desafío time-limited con recompensas, competitivo
- **Uso:** Live Events, competencias semanales
- **Interactividad:** Alta
- **Tiempo estimado:** Variable (5-30 minutos)

### 38. Forum Discussion Thread
- **Bloom levels:** Analizar, Evaluar
- **Variantes:** Pregunta + respuestas estilo foro, votación
- **Uso:** Classroom tab, discusiones de clase
- **Interactividad:** Alta
- **Tiempo estimado:** Variable

### 39. Teacher Feedback Widget
- **Bloom levels:** N/A (feedback)
- **Variantes:** Comentarios del profesor en actividades, calificaciones
- **Uso:** Activity review, retroalimentación personalizada
- **Interactividad:** Baja (lectura)
- **Tiempo estimado:** Variable

### 40. Study Buddy Match
- **Bloom levels:** N/A (social)
- **Variantes:** Sugiere compañeros con intereses similares
- **Uso:** Social features, formación de grupos
- **Interactividad:** Media (aceptar/rechazar)
- **Tiempo estimado:** N/A

---

## 🎨 G. COMPONENTES DE PROYECTOS Y CREACIÓN (Higher-Order Thinking)

Para niveles altos de Bloom (Crear):

### 41. Essay Builder
- **Bloom levels:** Crear
- **Variantes:** Editor con estructura guiada (intro, desarrollo, conclusión)
- **Uso:** Lenguaje, ensayos argumentativos
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 30-60 minutos

### 42. Presentation Maker
- **Bloom levels:** Crear
- **Variantes:** Slides simples para proyectos, plantillas
- **Uso:** Todos los temas, presentaciones orales
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 30-90 minutos

### 43. Experiment Designer
- **Bloom levels:** Crear, Evaluar
- **Variantes:** Plantilla para diseñar experimentos científicos
- **Uso:** Ciencias (Física, Química, Biología)
- **Interactividad:** Muy Alta
- **Tiempo estimado:** 20-40 minutos

### 44. Creative Portfolio
- **Bloom levels:** N/A (showcase)
- **Variantes:** Galería de trabajos del estudiante, categorías
- **Uso:** Perfil personal, muestra de progreso
- **Interactividad:** Media (navegación)
- **Tiempo estimado:** N/A

---

## 🎯 PRIORIZACIÓN RECOMENDADA

### **Fase 1: Core Learning (MVP)** ✅ COMPLETADA
1. ✅ Multiple Choice Question
2. ✅ True/False Statement
3. ✅ Open-Ended Text Response
4. ✅ Fill in the Blanks
5. ✅ Drag & Drop Matching
6. ✅ Sequencing/Ordering
7. ✅ Compare & Contrast (NUEVO)
8. ✅ Criteria Evaluation (NUEVO)
9. ✅ Concept Map Builder (NUEVO)

**Justificación:** Estos 9 componentes cubren los 6 niveles completos de Bloom (Recordar → Crear). Los primeros 6 son esenciales para cualquier actividad de aprendizaje. Los 3 adicionales (Compare & Contrast, Criteria Evaluation, Concept Map Builder) completan la cobertura de niveles superiores de pensamiento (Analizar, Evaluar, Crear) con componentes especializados.

### **Fase 2: Gamification & Progress**
7. Bloom Level Progress Wheel
8. Learning Path Map
9. Achievement Badge Display
10. Streak Calendar
11. PAES Readiness Meter

**Justificación:** Aumentan la motivación y engagement mediante visualización de progreso y mecánicas de juego.

### **Fase 3: Adaptive & Advanced**
12. Personalized Recommendation Card
13. PAES Practice Section
14. AI Tutor Chatbot
15. Weak Spot Identifier

**Justificación:** Aprovechan el perfil adaptativo del estudiante y preparan para el examen PAES.

### **Fase 4: Social & Collaborative**
16. Study Group Card
17. Forum Discussion Thread
18. Challenge Card

**Justificación:** Fomentan aprendizaje social y retención mediante comunidad.

### **Fase 5: Higher-Order Thinking**
19. Essay Builder
20. Math Equation Solver
21. Interactive Diagram/Labeling
22. Experiment Designer

**Justificación:** Cubren niveles altos de Bloom (Evaluar, Crear) para aprendizaje profundo.

---

## 📐 ESTRUCTURA PROPUESTA DE COMPONENTE

Cada componente debería incluir esta metadata:

```javascript
{
  // Identificación
  id: "multiple-choice-v1",
  name: "Multiple Choice Question",
  displayName: "Pregunta de Selección Múltiple",

  // Categoría
  category: "learning-activity", // learning-activity | content-delivery | assessment | progress | adaptive | social | creation

  // Alineación educativa
  bloomLevels: ["recordar", "comprender", "aplicar"],
  learningFormats: ["visual", "text", "interactive"],

  // Características
  estimatedTime: 2, // minutos
  difficulty: "1-10", // escala de complejidad
  materias: ["matemáticas", "lenguaje", "historia"], // todas o específicas
  interactivity: "high", // low | medium | high | very-high

  // Características técnicas
  adaptiveSupport: true, // puede adaptar dificultad
  paesRelevant: true, // útil para PAES
  requiresInternet: false, // funciona offline
  aiPowered: false, // usa IA para generar/evaluar

  // Integración backend
  endpoints: ["/api/educational/progress", "/api/educational/history"],
  dataModel: "oa_bloom_objectives",

  // UX
  mobileOptimized: true,
  keyboardNavigation: true,
  accessibility: "WCAG-AA"
}
```

---

## 🔧 CONSIDERACIONES TÉCNICAS

### **Stack Tecnológico (Lumera)**
- **Frontend:** Svelte 5 (runes: `$state`, `$derived`, `$effect`)
- **Estilos:** Tailwind CSS + GSAP animations
- **Backend:** Go + Chi + GORM + PostgreSQL
- **Despliegue:** Docker Compose

### **Patrón de Componente Svelte**

```svelte
<script>
  // Props
  let {
    data,
    bloomLevel,
    materia,
    onComplete,
    onAnswer
  } = $props();

  // Estado local
  let userAnswer = $state(null);
  let isCorrect = $state(false);
  let showFeedback = $state(false);

  // Efectos reactivos
  $effect(() => {
    // Validación, animaciones, etc.
  });

  // Funciones
  function handleSubmit() {
    // Lógica de validación
    // Emitir evento onAnswer con datos para backend
    onAnswer?.({
      oa_id: data.oa_id,
      bloom_level: bloomLevel,
      user_answer: userAnswer,
      is_correct: isCorrect,
      timestamp: new Date().toISOString()
    });
  }
</script>

<!-- Template con Tailwind + GSAP -->
<div class="component-container">
  <!-- UI del componente -->
</div>

<style>
  /* Estilos adicionales si es necesario */
</style>
```

### **Integración con Backend Lumera**

Los componentes deben:
1. Consumir datos de `oa_bloom_objectives` (objetivos de aprendizaje)
2. Trackear progreso en `student_oa_progress` (estados: no_iniciado, en_proceso, logrado, dominado)
3. Registrar eventos en `student_oa_history` (tipo: evaluación, práctica, diagnóstico, repaso)
4. Adaptar según `student_profile.profile_data` (preferencias de aprendizaje)

---

## 📚 PRÓXIMOS PASOS

1. ✅ **Documentar todas las sugerencias** (este archivo)
2. ✅ **Implementar Fase 1** (9 componentes core con cobertura completa de Bloom)
3. ✅ **Crear página demo** para probar componentes (`/components-demo`)
4. ⏳ **Integrar con backend** Lumera (endpoints `/api/educational/progress` y `/api/educational/complete`)
5. ⏳ **Testing con usuarios** reales
6. ⏳ **Iterar y mejorar** según feedback
7. ⏳ **Fase 2:** Implementar componentes de gamificación y progreso

---

## 📊 ESTADO ACTUAL DE IMPLEMENTACIÓN

**Componentes Implementados:** 21/50 (42%)
- ✅ A. Core Learning: 9/10 componentes
  - Multiple Choice, True/False, Open-Ended, Fill Blanks, Drag & Drop, Sequencing
  - Compare & Contrast, Criteria Evaluation, Concept Map Builder
- ✅ B. Lesson Slides: 12/12 componentes (100% COMPLETADO)
  - **Generales:** ConceptIntroSlide, ComparisonTableSlide, StepByStepProcessSlide, FormulaExplorerSlide, PracticePromptSlide, LessonPlayer
  - **Lenguaje:** TextAnnotationSlide, SentenceBuilderSlide, VocabularyContextSlide, TextStructureSlide, ConnectorsWorkshopSlide, LiteraryDevicesExplorerSlide
- ⏳ C. Content Delivery: 0/6 componentes
- ⏳ D. Assessment & Feedback: 0/5 componentes
- ⏳ E. Progress & Gamification: 0/8 componentes
- ⏳ F. AI-Powered: 0/6 componentes
- ⏳ G. Social & Collaborative: 0/5 componentes
- ⏳ H. Higher-Order Thinking: 0/4 componentes

**Cobertura de Bloom:**
- ✅ Recordar: 3 componentes especializados
- ✅ Comprender: 5 componentes especializados
- ✅ Aplicar: 3 componentes especializados
- ✅ Analizar: 2 componentes (OpenEnded + CompareContrast)
- ✅ Evaluar: 2 componentes (OpenEnded + CriteriaEvaluation)
- ✅ Crear: 2 componentes (OpenEnded + ConceptMapBuilder)

---

**Versión:** 3.1
**Última actualización:** 2025-11-22
**Autor:** Claude (Anthropic) + Johnny (Lumera Team)

**Changelog:**
- v3.1: Implementados 6 componentes de slides especializados en Lenguaje (TextAnnotation, SentenceBuilder, VocabularyContext, TextStructure, ConnectorsWorkshop, LiteraryDevicesExplorer) con 3 lecciones demo chilenas (García Márquez, gramática, escritura argumentativa). Cobertura Lenguaje: ~80-85%
- v3.0: Implementados 6 componentes de slides educativos generales (ConceptIntro, ComparisonTable, StepByStep, FormulaExplorer, PracticePrompt, LessonPlayer) con demo en `/lessons-demo`
- v2.0: Fase 1 completada con 9 componentes de actividades (agregados CompareContrast, CriteriaEvaluation, ConceptMapBuilder para completar cobertura de Bloom)
- v1.0: Documentación inicial de 44 componentes sugeridos
