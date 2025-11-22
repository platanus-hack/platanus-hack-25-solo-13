<script>
  import LessonPlayer from '$lib/components/slides/LessonPlayer.svelte';

  // Estado
  let leccionActiva = $state('conectores-completo'); // Lección demo Teach → Practice

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
    },

    // LECCIONES DE LENGUAJE
    'comprension-garcia-marquez': {
      leccionId: "comprension-garcia-marquez",
      titulo: "Comprensión Lectora: Cien Años de Soledad",
      materia: "Lenguaje",
      slides: [
        {
          tipo: "TextAnnotationSlide",
          orden: 1,
          props: {
            titulo: "Análisis: 'Cien Años de Soledad' - Fragmento",
            texto: `Muchos años después, frente al pelotón de fusilamiento, el coronel Aureliano Buendía había de recordar aquella tarde remota en que su padre lo llevó a conocer el hielo. Macondo era entonces una aldea de veinte casas de barro y cañabrava construidas a la orilla de un río de aguas diáfanas que se precipitaban por un lecho de piedras pulidas, blancas y enormes como huevos prehistóricos. El mundo era tan reciente, que muchas cosas carecían de nombre, y para mencionarlas había que señalarlas con el dedo.`,
            tipoLectura: "narrativa",
            preguntasGuia: [
              "¿Qué técnica narrativa usa García Márquez al comenzar in medias res?",
              "¿Qué ambiente describe el autor en el segundo párrafo?",
              "Identifica elementos del realismo mágico en el fragmento"
            ],
            vocabularioDestacado: ["pelotón", "fusilamiento", "remota", "diáfanas"],
            herramientasAnotacion: ["resaltar", "notas"],
            materia: "lenguaje"
          }
        },
        {
          tipo: "VocabularyContextSlide",
          orden: 2,
          props: {
            palabraObjetivo: "diáfano",
            pronunciacion: "/di.ˈa.fa.no/",
            audioUrl: null,
            etimologia: "Del latín 'diaphanus', y este del griego 'diaphanés': transparente",
            definicion: "Que deja pasar la luz casi en su totalidad. También se usa figuradamente para referirse a algo claro, limpio y transparente.",
            morfologia: {
              prefijo: "dia- (a través de)",
              raiz: "phan (aparecer, mostrar)",
              sufijo: "-o (adjetivo)"
            },
            contextosEjemplo: [
              "El río de aguas diáfanas reflejaba el sol de la mañana.",
              "Su explicación fue diáfana y todos la comprendieron.",
              "El cristal diáfano de la ventana permitía ver el jardín."
            ],
            imagenes: [],
            sinonimos: ["transparente", "claro", "límpido", "cristalino"],
            antonimos: ["opaco", "turbio", "oscuro"],
            palabrasFamilia: ["diafanidad", "diáfanamente"],
            materia: "lenguaje"
          }
        },
        {
          tipo: "LiteraryDevicesExplorerSlide",
          orden: 3,
          props: {
            titulo: "Recursos Literarios en García Márquez",
            texto: `Muchos años después, frente al pelotón de fusilamiento, el coronel Aureliano Buendía había de recordar aquella tarde remota en que su padre lo llevó a conocer el hielo.`,
            autor: "Gabriel García Márquez",
            dispositivosLiterarios: [
              {
                tipo: "anafora",
                ejemplos: ["Muchos años después"],
                definicion: "Repetición de palabras al inicio de frases sucesivas",
                efecto: "Crea ritmo y énfasis temporal, estableciendo distancia narrativa",
                color: "purple"
              },
              {
                tipo: "metafora",
                ejemplos: ["tarde remota"],
                definicion: "Comparación implícita entre dos elementos sin usar 'como'",
                efecto: "Evoca nostalgia y distancia temporal de manera poética",
                color: "blue"
              },
              {
                tipo: "simil",
                ejemplos: ["piedras pulidas, blancas y enormes como huevos prehistóricos"],
                definicion: "Comparación explícita usando 'como' o 'cual'",
                efecto: "Crea una imagen visual potente y surrealista, característica del realismo mágico",
                color: "green"
              }
            ],
            preguntasAnalisis: [
              "¿Por qué García Márquez comienza la novela con una anticipación?",
              "¿Qué efecto tiene la imagen de las piedras como 'huevos prehistóricos'?"
            ],
            modoComparacion: false,
            ejercicioCreativo: "Escribe una oración usando metáfora y otra usando símil",
            materia: "lenguaje"
          }
        },
        {
          tipo: "PracticePromptSlide",
          orden: 4,
          props: {
            mensaje: "¡Ahora dominas el análisis literario!",
            submensaje: "Refuerza lo aprendido con ejercicios de comprensión",
            icono: "📚",
            previewEjercicios: [
              { tipo: "Análisis de texto", cantidad: 3, icono: "📖" },
              { tipo: "Vocabulario contextual", cantidad: 5, icono: "📝" },
              { tipo: "Recursos literarios", cantidad: 4, icono: "✨" }
            ],
            motivacion: "¡Excelente! Ya comprendes técnicas narrativas avanzadas",
            botonTexto: "Comenzar Ejercicios",
            colorTema: "purple",
            materia: "lenguaje"
          }
        }
      ]
    },

    'gramatica-oraciones': {
      leccionId: "gramatica-oraciones",
      titulo: "Gramática: Construcción de Oraciones Complejas",
      materia: "Lenguaje",
      slides: [
        {
          tipo: "SentenceBuilderSlide",
          orden: 1,
          props: {
            titulo: "Construye Oraciones Simples y Compuestas",
            modoEjercicio: "guided",
            oracionBase: null,
            palabrasDisponibles: [
              { texto: "El estudiante", tipo: "sujeto", color: "blue" },
              { texto: "chileno", tipo: "adjetivo", color: "purple" },
              { texto: "estudia", tipo: "verbo", color: "red" },
              { texto: "diligentemente", tipo: "adverbio", color: "orange" },
              { texto: "para la PAES", tipo: "complemento", color: "green" },
              { texto: "en la biblioteca", tipo: "complemento", color: "green" }
            ],
            objetivoGramatical: "Construye una oración compleja agregando modificadores al sujeto y al verbo",
            variacionesCorrectas: [
              "El estudiante chileno estudia diligentemente para la PAES",
              "El estudiante estudia diligentemente para la PAES en la biblioteca",
              "El estudiante chileno estudia diligentemente para la PAES en la biblioteca"
            ],
            mostrarScaffolding: true,
            materia: "lenguaje"
          }
        },
        {
          tipo: "ConnectorsWorkshopSlide",
          orden: 2,
          props: {
            titulo: "Conectores: Dando Coherencia al Texto",
            parrafos: [
              {
                texto: "El cambio climático es una realidad científica. ___ debemos actuar con urgencia para mitigar sus efectos.",
                opcionesConector: ["Por lo tanto", "Sin embargo", "Además"],
                correcta: "Por lo tanto",
                tipo: "consecutivos",
                explicacion: "El segundo enunciado es una consecuencia lógica del primero, por eso usamos un conector consecutivo"
              },
              {
                texto: "Chile es líder en energía solar en Sudamérica. ___ muchos países aún dependen de combustibles fósiles.",
                opcionesConector: ["Además", "Sin embargo", "Por eso"],
                correcta: "Sin embargo",
                tipo: "adversativos",
                explicacion: "Se presenta un contraste entre Chile y otros países, requiere un conector adversativo"
              },
              {
                texto: "Las energías renovables reducen emisiones. ___ crean empleos verdes en las comunidades.",
                opcionesConector: ["Además", "Por lo tanto", "Aunque"],
                correcta: "Además",
                tipo: "aditivos",
                explicacion: "Se agrega información complementaria, se usa conector aditivo"
              }
            ],
            categorias: ["causales", "adversativos", "consecutivos", "aditivos"],
            bancoConectores: {
              causales: ["porque", "ya que", "debido a", "puesto que"],
              adversativos: ["sin embargo", "no obstante", "pero", "aunque"],
              consecutivos: ["por lo tanto", "en consecuencia", "así que", "entonces"],
              aditivos: ["además", "asimismo", "también", "igualmente"]
            },
            modoComparacion: false,
            ejercicioCreativo: "Escribe un párrafo sobre tu ciudad usando al menos 3 tipos diferentes de conectores",
            materia: "lenguaje"
          }
        },
        {
          tipo: "VocabularyContextSlide",
          orden: 3,
          props: {
            palabraObjetivo: "resiliente",
            pronunciacion: "/re.si.ˈljen.te/",
            audioUrl: null,
            etimologia: "Del latín 'resilire': saltar hacia atrás, rebotar",
            definicion: "Capacidad de un material, persona o ecosistema para recuperarse ante adversidades o perturbaciones.",
            morfologia: {
              prefijo: "re- (de nuevo)",
              raiz: "sil (saltar)",
              sufijo: "-ente (que realiza la acción)"
            },
            contextosEjemplo: [
              "Los ecosistemas resilientes se recuperan rápidamente después de incendios forestales.",
              "María demostró ser resiliente tras perder su empleo durante la pandemia.",
              "Necesitamos materiales resilientes para construcciones antisísmicas en Chile."
            ],
            imagenes: [],
            sinonimos: ["resistente", "adaptable", "fuerte", "flexible"],
            antonimos: ["frágil", "débil", "vulnerable", "quebradizo"],
            palabrasFamilia: ["resiliencia", "resilir"],
            materia: "lenguaje"
          }
        },
        {
          tipo: "PracticePromptSlide",
          orden: 4,
          props: {
            mensaje: "¡Listo para escribir con coherencia!",
            submensaje: "Practica con ejercicios de redacción y gramática",
            icono: "✍️",
            previewEjercicios: [
              { tipo: "Construcción de oraciones", cantidad: 5, icono: "🔨" },
              { tipo: "Uso de conectores", cantidad: 6, icono: "🔗" },
              { tipo: "Vocabulario en contexto", cantidad: 4, icono: "📚" }
            ],
            motivacion: "Dominas el 80% de Gramática y Redacción - ¡Sigue así!",
            botonTexto: "Practicar Ahora",
            colorTema: "purple",
            materia: "lenguaje"
          }
        }
      ]
    },

    'escritura-argumentativa': {
      leccionId: "escritura-argumentativa",
      titulo: "Escritura Argumentativa: Estructura y Coherencia",
      materia: "Lenguaje",
      slides: [
        {
          tipo: "TextStructureSlide",
          orden: 1,
          props: {
            titulo: "Estructura de un Texto Argumentativo",
            tipoTexto: "argumentativo",
            textoEjemplo: "Texto sobre energías renovables en Chile",
            estructura: {
              introduccion: {
                texto: "Chile enfrenta desafíos importantes en materia energética. El país debe decidir entre continuar con energías tradicionales o transitar hacia fuentes renovables.",
                color: "blue",
                descripcion: "Contextualiza el tema y presenta el problema"
              },
              tesis: {
                texto: "Chile debe aumentar significativamente su inversión en energías renovables para combatir el cambio climático y asegurar la sustentabilidad energética del país.",
                color: "purple",
                descripcion: "Postura clara y defendible del autor"
              },
              argumentos: {
                texto: "Argumento 1: Reduce emisiones de CO₂. Estudios del MIT demuestran que invertir en solar y eólica reduce emisiones en un 40% en 10 años. Argumento 2: Crea empleos verdes. El sector solar creó 50,000 empleos en Chile en 2023, dinamizando la economía regional.",
                color: "green",
                descripcion: "Razones y evidencia que apoyan la tesis"
              },
              contraargumentos: {
                texto: "Algunos argumentan que es muy costoso. Sin embargo, los costos de energía solar han bajado 90% desde 2010 según IRENA, haciéndola competitiva con combustibles fósiles.",
                color: "red",
                descripcion: "Reconoce y refuta objeciones"
              },
              conclusion: {
                texto: "Por lo tanto, la inversión en energías renovables no solo es viable económicamente, sino necesaria para el futuro sustentable de Chile.",
                color: "orange",
                descripcion: "Retoma la tesis y cierra el argumento"
              }
            },
            ejercicioTipo: "solo-visualizar",
            comparacionTipos: true,
            materia: "lenguaje"
          }
        },
        {
          tipo: "TextAnnotationSlide",
          orden: 2,
          props: {
            titulo: "Analiza un Ensayo Argumentativo",
            texto: `Introducción: En la era digital, el uso excesivo de redes sociales ha generado un debate sobre su impacto en adolescentes chilenos.

Tesis: Las redes sociales, si bien ofrecen beneficios de conectividad, generan más perjuicios que beneficios en el desarrollo socioemocional de los jóvenes.

Argumento 1: Diversos estudios demuestran que el uso prolongado de Instagram y TikTok está asociado con mayor ansiedad y depresión en adolescentes. Un estudio de la Universidad de Chile (2023) reveló que el 65% de jóvenes entre 14-18 años reporta sentimientos de inadecuación al compararse con contenido en redes.

Contraargumento y refutación: Aunque las redes permiten mantener contacto con amigos, esta conexión es superficial. La comunicación cara a cara desarrolla habilidades sociales más profundas que el intercambio de mensajes.

Conclusión: Por lo tanto, es urgente educar sobre uso responsable de redes sociales para proteger la salud mental juvenil.`,
            tipoLectura: "argumentativa",
            preguntasGuia: [
              "Identifica la tesis del texto",
              "¿Qué evidencia usa el autor para su argumento principal?",
              "¿Cómo refuta el autor el contraargumento?"
            ],
            vocabularioDestacado: ["inadecuación", "socioemocional", "superficial"],
            herramientasAnotacion: ["resaltar", "notas"],
            materia: "lenguaje"
          }
        },
        {
          tipo: "ConnectorsWorkshopSlide",
          orden: 3,
          props: {
            titulo: "Conectores en Textos Argumentativos",
            parrafos: [
              {
                texto: "La educación online tiene ventajas como flexibilidad horaria. ___ presenta desafíos en la interacción social de los estudiantes.",
                opcionesConector: ["Además", "Sin embargo", "Por consiguiente"],
                correcta: "Sin embargo",
                tipo: "adversativos",
                explicacion: "Se presenta un contraste entre ventajas y desafíos"
              },
              {
                texto: "Los jóvenes pasan más de 5 horas diarias en redes sociales. ___ es necesario implementar programas de alfabetización digital en los colegios.",
                opcionesConector: ["Por lo tanto", "Aunque", "Además"],
                correcta: "Por lo tanto",
                tipo: "consecutivos",
                explicacion: "La segunda oración es una consecuencia lógica de la primera"
              }
            ],
            categorias: ["adversativos", "consecutivos", "causales", "aditivos"],
            bancoConectores: {
              causales: ["porque", "ya que", "debido a"],
              adversativos: ["sin embargo", "no obstante", "aunque"],
              consecutivos: ["por lo tanto", "en consecuencia", "por consiguiente"],
              aditivos: ["además", "asimismo", "también"]
            },
            modoComparacion: false,
            ejercicioCreativo: null,
            materia: "lenguaje"
          }
        },
        {
          tipo: "PracticePromptSlide",
          orden: 4,
          props: {
            mensaje: "¡Ya puedes escribir ensayos argumentativos!",
            submensaje: "Practica redactando tu propio ensayo sobre un tema actual",
            icono: "📝",
            previewEjercicios: [
              { tipo: "Identificar estructuras", cantidad: 4, icono: "🏗️" },
              { tipo: "Redactar argumentos", cantidad: 3, icono: "💭" },
              { tipo: "Uso de evidencia", cantidad: 3, icono: "📊" }
            ],
            motivacion: "¡Excelente progreso! Dominas escritura argumentativa al 85%",
            botonTexto: "Comenzar a Escribir",
            colorTema: "purple",
            materia: "lenguaje"
          }
        }
      ]
    },

    // LECCIÓN DEMO: TEACH → PRACTICE
    'conectores-completo': {
      leccionId: "conectores-completo",
      titulo: "Conectores Textuales: De la Teoría a la Práctica",
      materia: "Lenguaje",
      slides: [
        // SLIDE 1: TEACH - Introducción a conectores
        {
          tipo: "ConnectorsGuideSlide",
          orden: 1,
          props: {
            titulo: "Guía Completa de Conectores Textuales",
            importancia: "Los conectores son palabras o frases que unen ideas y dan coherencia al texto, esenciales para la PAES de Lenguaje",
            categorias: [
              {
                nombre: "Causales",
                tipo: "causal",
                definicion: "Indican la causa o razón de algo",
                funcion: "Expresan relación de causa-efecto entre ideas",
                conectores: [
                  {
                    palabra: "porque",
                    nivel: "básico",
                    ejemplos: [
                      "Estudió toda la noche porque tenía examen al día siguiente.",
                      "No salió porque estaba lloviendo."
                    ],
                    nota: "El conector causal más común en español"
                  },
                  {
                    palabra: "ya que",
                    nivel: "formal",
                    ejemplos: [
                      "Aprobó el curso ya que cumplió con todos los requisitos.",
                      "El proyecto fue cancelado ya que no había presupuesto."
                    ]
                  },
                  {
                    palabra: "debido a que",
                    nivel: "formal",
                    ejemplos: [
                      "La reunión se pospuso debido a que faltó el director."
                    ]
                  }
                ],
                ejemploParrafo: "El reciclaje es fundamental para el medio ambiente porque reduce la cantidad de desechos. Ya que muchas personas no están conscientes de su importancia, es necesario educar desde las escuelas.",
                comparacion: "A diferencia de los consecutivos que muestran resultado, los causales muestran el origen o razón"
              },
              {
                nombre: "Adversativos",
                tipo: "adversativo",
                definicion: "Expresan oposición o contraste entre ideas",
                funcion: "Marcan un cambio de dirección en el argumento",
                conectores: [
                  {
                    palabra: "sin embargo",
                    nivel: "formal",
                    ejemplos: [
                      "Estudió mucho; sin embargo, no aprobó el examen.",
                      "La propuesta era buena; sin embargo, fue rechazada."
                    ],
                    nota: "Usado frecuentemente en textos argumentativos"
                  },
                  {
                    palabra: "pero",
                    nivel: "básico",
                    ejemplos: [
                      "Quería ir al cine, pero no tenía dinero.",
                      "Es inteligente, pero no estudia."
                    ]
                  },
                  {
                    palabra: "no obstante",
                    nivel: "muy formal",
                    ejemplos: [
                      "El clima era adverso; no obstante, continuaron la expedición."
                    ]
                  }
                ],
                ejemploParrafo: "Las energías renovables son más limpias. Sin embargo, su implementación requiere inversión inicial alta. No obstante, a largo plazo resultan más económicas.",
                comparacion: "Los adversativos contradicen la idea anterior, mientras que los aditivos la complementan"
              },
              {
                nombre: "Consecutivos",
                tipo: "consecutivo",
                definicion: "Indican consecuencia o resultado de una acción",
                funcion: "Muestran el efecto o resultado de una causa",
                conectores: [
                  {
                    palabra: "por lo tanto",
                    nivel: "formal",
                    ejemplos: [
                      "El experimento falló; por lo tanto, debemos repetirlo.",
                      "No hay evidencia; por lo tanto, se desestima la acusación."
                    ],
                    nota: "Muy común en textos académicos y científicos"
                  },
                  {
                    palabra: "en consecuencia",
                    nivel: "formal",
                    ejemplos: [
                      "Violó las reglas; en consecuencia, fue expulsado."
                    ]
                  },
                  {
                    palabra: "así que",
                    nivel: "básico",
                    ejemplos: [
                      "Estaba cansado, así que se fue a dormir."
                    ]
                  }
                ],
                ejemploParrafo: "La deforestación aumenta el CO₂ en la atmósfera. Por lo tanto, contribuye al calentamiento global. En consecuencia, debemos proteger los bosques.",
                comparacion: "Los consecutivos muestran el resultado, mientras que los causales muestran la razón"
              },
              {
                nombre: "Aditivos",
                tipo: "aditivo",
                definicion: "Agregan información complementaria",
                funcion: "Suman ideas o argumentos en la misma dirección",
                conectores: [
                  {
                    palabra: "además",
                    nivel: "básico",
                    ejemplos: [
                      "Es inteligente. Además, es muy trabajador.",
                      "Llegó tarde. Además, no trajo los documentos."
                    ],
                    nota: "Uno de los conectores aditivos más versátiles"
                  },
                  {
                    palabra: "asimismo",
                    nivel: "formal",
                    ejemplos: [
                      "La propuesta es viable. Asimismo, es económica."
                    ]
                  },
                  {
                    palabra: "también",
                    nivel: "básico",
                    ejemplos: [
                      "Estudia matemáticas. También estudia física."
                    ]
                  }
                ],
                ejemploParrafo: "La lectura mejora el vocabulario. Además, desarrolla la imaginación. Asimismo, fortalece la concentración.",
                comparacion: "Los aditivos suman ideas del mismo tipo, mientras que los adversativos las contrastan"
              },
              {
                nombre: "Temporales",
                tipo: "temporal",
                definicion: "Indican relaciones de tiempo entre eventos",
                funcion: "Organizan eventos en secuencia cronológica",
                conectores: [
                  {
                    palabra: "luego",
                    nivel: "básico",
                    ejemplos: [
                      "Desayunó. Luego se fue a trabajar.",
                      "Primero estudió teoría. Luego hizo los ejercicios."
                    ]
                  },
                  {
                    palabra: "después",
                    nivel: "básico",
                    ejemplos: [
                      "Terminó la tarea. Después salió a jugar."
                    ]
                  },
                  {
                    palabra: "finalmente",
                    nivel: "básico",
                    ejemplos: [
                      "Preparó los ingredientes. Cocinó el plato. Finalmente, lo sirvió."
                    ],
                    nota: "Indica el último paso en una secuencia"
                  }
                ],
                ejemploParrafo: "Primero, se mezclan los ingredientes secos. Luego, se agregan los líquidos. Después, se hornea durante 30 minutos. Finalmente, se deja enfriar.",
                comparacion: "Los temporales organizan en el tiempo, mientras que los consecutivos organizan por causa-efecto"
              }
            ],
            consejos: [
              "Varía los conectores para evitar repetición",
              "Usa conectores formales en textos académicos",
              "No abuses de los conectores: úsalos solo cuando sea necesario",
              "Asegúrate de que el conector elegido refleje la relación lógica correcta",
              "En la PAES, identifica primero la relación entre ideas antes de elegir el conector"
            ],
            materia: "lenguaje"
          }
        },

        // SLIDE 2: PRACTICE - Aplicar conocimiento de conectores
        {
          tipo: "ConnectorsWorkshopSlide",
          orden: 2,
          props: {
            titulo: "Practica: Elige el Conector Correcto",
            parrafos: [
              {
                texto: "El reciclaje reduce la contaminación. ___ ayuda a conservar los recursos naturales.",
                opcionesConector: ["Por lo tanto", "Además", "Sin embargo"],
                correcta: "Además",
                tipo: "aditivos",
                explicacion: "Se está agregando una ventaja adicional del reciclaje, no una consecuencia ni un contraste"
              },
              {
                texto: "Quería comprar el libro. ___ no tenía suficiente dinero.",
                opcionesConector: ["Sin embargo", "Por lo tanto", "Porque"],
                correcta: "Sin embargo",
                tipo: "adversativos",
                explicacion: "Hay contraste entre el deseo y la imposibilidad de realizarlo"
              },
              {
                texto: "No estudió para el examen. ___ reprobó el curso.",
                opcionesConector: ["Además", "Por lo tanto", "Sin embargo"],
                correcta: "Por lo tanto",
                tipo: "consecutivos",
                explicacion: "Reprobar es la consecuencia lógica de no estudiar"
              },
              {
                texto: "La energía solar es limpia. ___ cada vez es más accesible económicamente.",
                opcionesConector: ["Además", "Sin embargo", "Por lo tanto"],
                correcta: "Además",
                tipo: "aditivos",
                explicacion: "Se suma otro beneficio de la energía solar"
              },
              {
                texto: "Llegó temprano a la reunión. ___ todos lo felicitaron por su puntualidad.",
                opcionesConector: ["Luego", "Porque", "Sin embargo"],
                correcta: "Luego",
                tipo: "temporales",
                explicacion: "Indica la secuencia temporal: primero llegó, después lo felicitaron"
              }
            ],
            ejercicioCreativo: "Escribe un párrafo de 4 oraciones sobre las energías renovables usando al menos 3 conectores diferentes (uno causal, uno adversativo y uno consecutivo)",
            materia: "lenguaje"
          }
        },

        // SLIDE 3: TEACH - Gramática
        {
          tipo: "GrammarConceptSlide",
          orden: 3,
          props: {
            titulo: "Tipos de Sujeto en Español",
            concepto: "Sujeto",
            definicion: "El sujeto es la persona, animal o cosa que realiza la acción del verbo o de quien se dice algo. Responde a las preguntas '¿quién?' o '¿quiénes?'",
            tipos: [
              {
                nombre: "Sujeto Expreso",
                definicion: "El sujeto está explícitamente mencionado en la oración",
                caracteristicas: [
                  "Se menciona directamente en la oración",
                  "Puede estar al inicio, medio o final de la oración",
                  "Concuerda en número y persona con el verbo"
                ],
                ejemplos: [
                  {
                    oracion: "Los estudiantes chilenos se preparan para la PAES",
                    analisis: "Sujeto expreso: 'Los estudiantes chilenos' (3ª persona plural)"
                  },
                  {
                    oracion: "María y Pedro estudian juntos en la biblioteca",
                    analisis: "Sujeto expreso compuesto: 'María y Pedro'"
                  },
                  {
                    oracion: "En la sala de clases explicó el profesor el teorema de Pitágoras",
                    analisis: "Sujeto expreso pospuesto: 'el profesor' (viene después del verbo)"
                  }
                ],
                estructura: "Núcleo + Modificadores (opcionales)"
              },
              {
                nombre: "Sujeto Tácito (Desinencial)",
                definicion: "El sujeto no se menciona explícitamente, pero se entiende por la conjugación del verbo",
                caracteristicas: [
                  "No aparece escrito en la oración",
                  "Se deduce por la desinencia (terminación) del verbo",
                  "Muy común en español, a diferencia del inglés"
                ],
                ejemplos: [
                  {
                    oracion: "Estudiamos para la PAES todos los días",
                    analisis: "Sujeto tácito: 'Nosotros/Nosotras' (deducido de la terminación -amos)"
                  },
                  {
                    oracion: "Llegó tarde a la clase",
                    analisis: "Sujeto tácito: 'Él/Ella/Usted' (deducido de la terminación -ó)"
                  },
                  {
                    oracion: "¿Vas al cine esta noche?",
                    analisis: "Sujeto tácito: 'Tú' (deducido de la terminación -as)"
                  }
                ],
                estructura: "∅ (ausente) + Verbo conjugado"
              },
              {
                nombre: "Sujeto Indeterminado",
                definicion: "No se sabe quién realiza la acción o no interesa especificarlo",
                caracteristicas: [
                  "Se desconoce la identidad del sujeto",
                  "Usa la tercera persona del plural sin sujeto expreso",
                  "También con 'se' + verbo en 3ª persona singular"
                ],
                ejemplos: [
                  {
                    oracion: "Dicen que habrá cambios en la PAES",
                    analisis: "No se sabe quiénes dicen (3ª persona plural sin sujeto)"
                  },
                  {
                    oracion: "Se vende casa",
                    analisis: "Construcción con 'se' impersonal"
                  },
                  {
                    oracion: "Llamaron por teléfono",
                    analisis: "No se especifica quiénes llamaron"
                  }
                ],
                estructura: "Verbo en 3ª persona plural / Se + Verbo"
              }
            ],
            reglas: [
              {
                texto: "El sujeto siempre concuerda con el verbo en número y persona",
                ejemplo: "Los niños juegan (plural) / El niño juega (singular)"
              },
              {
                texto: "El sujeto puede omitirse en español cuando es obvio por el contexto",
                ejemplo: "¿Vas al cine? (tácito: tú) - Sí, voy (tácito: yo)"
              },
              {
                texto: "Para encontrar el sujeto, pregunta al verbo '¿quién?' o '¿quiénes?'",
                ejemplo: "María estudia → ¿Quién estudia? → María (sujeto)"
              }
            ],
            erroresComunes: [
              {
                incorrecto: "Los estudiante estudian",
                correcto: "Los estudiantes estudian",
                explicacion: "El sujeto plural 'los estudiantes' debe concordar con el verbo 'estudian'"
              },
              {
                incorrecto: "Yo y tú vamos al cine",
                correcto: "Tú y yo vamos al cine",
                explicacion: "Por cortesía, la 1ª persona (yo) va al final en sujetos compuestos"
              }
            ],
            materia: "lenguaje"
          }
        },

        // SLIDE 4: PRACTICE - Construcción de oraciones
        {
          tipo: "SentenceBuilderSlide",
          orden: 4,
          props: {
            titulo: "Construye Oraciones con Sujetos Variados",
            modoEjercicio: "guided",
            oracionBase: null,
            palabrasDisponibles: [
              { texto: "Los profesores", tipo: "sujeto", color: "blue" },
              { texto: "de la universidad", tipo: "complemento", color: "green" },
              { texto: "explicaron", tipo: "verbo", color: "red" },
              { texto: "detalladamente", tipo: "adverbio", color: "orange" },
              { texto: "el nuevo sistema", tipo: "complemento", color: "green" },
              { texto: "de evaluación", tipo: "complemento", color: "green" }
            ],
            objetivoGramatical: "Construye una oración con sujeto expreso compuesto por núcleo + modificadores",
            variacionesCorrectas: [
              "Los profesores de la universidad explicaron el nuevo sistema de evaluación",
              "Los profesores de la universidad explicaron detalladamente el nuevo sistema",
              "Los profesores explicaron detalladamente el nuevo sistema de evaluación",
              "Los profesores de la universidad explicaron detalladamente el nuevo sistema de evaluación"
            ],
            mostrarScaffolding: true,
            materia: "lenguaje"
          }
        },

        // SLIDE 5: Transición a evaluación
        {
          tipo: "PracticePromptSlide",
          orden: 5,
          props: {
            mensaje: "¡Dominas conectores y gramática!",
            submensaje: "Ahora estás listo para ejercicios de nivel PAES",
            icono: "🎓",
            previewEjercicios: [
              { tipo: "Conectores textuales", cantidad: 10, icono: "🔗" },
              { tipo: "Análisis sintáctico", cantidad: 8, icono: "📝" },
              { tipo: "Comprensión lectora", cantidad: 5, icono: "📖" }
            ],
            motivacion: "Has completado el ciclo: Aprender → Practicar → Dominar",
            botonTexto: "Ir a Evaluación PAES",
            colorTema: "purple",
            materia: "lenguaje"
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
        <div class="flex flex-wrap gap-2">
          <!-- Matemáticas -->
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

          <!-- Lenguaje -->
          <button
            onclick={() => leccionActiva = 'comprension-garcia-marquez'}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300
              {leccionActiva === 'comprension-garcia-marquez' ? 'bg-purple-500 text-white shadow-lg shadow-purple-500/50' : 'bg-slate-800 text-slate-400 hover:bg-slate-700'}
            "
          >
            📚 Comprensión
          </button>
          <button
            onclick={() => leccionActiva = 'gramatica-oraciones'}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300
              {leccionActiva === 'gramatica-oraciones' ? 'bg-purple-500 text-white shadow-lg shadow-purple-500/50' : 'bg-slate-800 text-slate-400 hover:bg-slate-700'}
            "
          >
            ✍️ Gramática
          </button>
          <button
            onclick={() => leccionActiva = 'escritura-argumentativa'}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300
              {leccionActiva === 'escritura-argumentativa' ? 'bg-purple-500 text-white shadow-lg shadow-purple-500/50' : 'bg-slate-800 text-slate-400 hover:bg-slate-700'}
            "
          >
            📝 Escritura
          </button>

          <!-- NUEVA: Lección Teach → Practice -->
          <button
            onclick={() => leccionActiva = 'conectores-completo'}
            class="
              px-4 py-2 rounded-xl font-semibold transition-all duration-300 border-2
              {leccionActiva === 'conectores-completo' ? 'bg-gradient-to-r from-blue-500 to-purple-500 text-white shadow-lg border-white' : 'bg-slate-800 text-slate-400 hover:bg-slate-700 border-blue-500/30'}
            "
          >
            🎓 Teach→Practice
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
