// Banco de preguntas diagnósticas para Lengua y Literatura - 1° Medio
// Niveles de dificultad: 1 (básico) - 5 (avanzado)

export interface PreguntaDiagnostica {
  id: string;
  tipo: 'multiple_choice' | 'true_false' | 'fill_blanks';
  dificultad: 1 | 2 | 3 | 4 | 5;
  pregunta: string;
  opciones?: string[]; // Para multiple_choice
  respuestaCorrecta: string | number; // String para fill_blanks, number (índice) para opciones
  explicacion: string;
  bloomLevel?: number;
}

export const PREGUNTAS_DIAGNOSTICO_LENGUA: PreguntaDiagnostica[] = [
  // ===== NIVEL BÁSICO (Dificultad 1-2) =====
  {
    id: 'leng_diag_01',
    tipo: 'multiple_choice',
    dificultad: 1,
    pregunta: '¿Cuál de las siguientes palabras es un sustantivo?',
    opciones: ['Correr', 'Casa', 'Rápido', 'Muy'],
    respuestaCorrecta: 1,
    explicacion: 'Un sustantivo es una palabra que nombra personas, animales, cosas o ideas. "Casa" es un sustantivo.',
    bloomLevel: 1
  },
  {
    id: 'leng_diag_02',
    tipo: 'true_false',
    dificultad: 1,
    pregunta: 'Los adjetivos son palabras que describen o califican a los sustantivos.',
    respuestaCorrecta: 0, // 0 = Verdadero, 1 = Falso
    explicacion: 'Correcto. Los adjetivos acompañan al sustantivo para expresar cualidades o características.',
    bloomLevel: 1
  },
  {
    id: 'leng_diag_03',
    tipo: 'multiple_choice',
    dificultad: 2,
    pregunta: 'Lee el siguiente fragmento: "El sol brillaba intensamente mientras los niños jugaban en el parque". ¿Cuál es el sujeto de la oración?',
    opciones: ['El sol', 'Los niños', 'El parque', 'El día'],
    respuestaCorrecta: 0,
    explicacion: 'El sujeto es "el sol", ya que es quien realiza la acción de brillar.',
    bloomLevel: 2
  },

  // ===== NIVEL INTERMEDIO (Dificultad 2-3) =====
  {
    id: 'leng_diag_04',
    tipo: 'multiple_choice',
    dificultad: 2,
    pregunta: '¿Qué tipo de narrador se utiliza en este fragmento? "Yo caminaba por la calle cuando vi a mi mejor amigo".',
    opciones: [
      'Narrador omnisciente',
      'Narrador protagonista (primera persona)',
      'Narrador testigo',
      'Narrador en tercera persona'
    ],
    respuestaCorrecta: 1,
    explicacion: 'Es un narrador protagonista porque cuenta la historia en primera persona ("yo") y es parte de los acontecimientos.',
    bloomLevel: 2
  },
  {
    id: 'leng_diag_05',
    tipo: 'fill_blanks',
    dificultad: 3,
    pregunta: 'Una ___1___ es una figura literaria que consiste en atribuir características humanas a objetos inanimados o animales.',
    respuestaCorrecta: 'personificación',
    explicacion: 'La personificación (o prosopopeya) es una figura literaria que da cualidades humanas a objetos o seres no humanos.',
    bloomLevel: 2
  },
  {
    id: 'leng_diag_06',
    tipo: 'multiple_choice',
    dificultad: 3,
    pregunta: 'En un texto narrativo, el clímax es:',
    opciones: [
      'La presentación de los personajes',
      'El momento de mayor tensión o conflicto',
      'La conclusión de la historia',
      'El contexto donde ocurre la historia'
    ],
    respuestaCorrecta: 1,
    explicacion: 'El clímax es el punto de mayor tensión en la narración, donde el conflicto alcanza su máxima intensidad.',
    bloomLevel: 3
  },

  // ===== NIVEL AVANZADO (Dificultad 3-4) =====
  {
    id: 'leng_diag_07',
    tipo: 'multiple_choice',
    dificultad: 3,
    pregunta: '¿Cuál es la función principal de un texto argumentativo?',
    opciones: [
      'Narrar una historia ficticia',
      'Describir un objeto o persona',
      'Convencer o persuadir al lector sobre un punto de vista',
      'Explicar cómo hacer algo paso a paso'
    ],
    respuestaCorrecta: 2,
    explicacion: 'Los textos argumentativos buscan convencer al lector mediante razones y evidencias que apoyen una tesis.',
    bloomLevel: 3
  },
  {
    id: 'leng_diag_08',
    tipo: 'true_false',
    dificultad: 4,
    pregunta: 'En una metáfora, se comparan dos elementos utilizando las palabras "como" o "tal como".',
    respuestaCorrecta: 1, // Falso
    explicacion: 'Falso. Esa es la definición de comparación o símil. La metáfora identifica directamente dos elementos sin usar "como".',
    bloomLevel: 2
  },
  {
    id: 'leng_diag_09',
    tipo: 'multiple_choice',
    dificultad: 4,
    pregunta: 'Lee el siguiente verso: "Volverán las oscuras golondrinas / en tu balcón sus nidos a colgar". ¿Qué figura literaria predomina?',
    opciones: [
      'Hipérbole',
      'Personificación',
      'Metáfora',
      'Anáfora'
    ],
    respuestaCorrecta: 1,
    explicacion: 'Hay personificación al atribuir acciones humanas (colgar) a las golondrinas de manera intencional.',
    bloomLevel: 4
  },

  // ===== NIVEL EXPERTO (Dificultad 4-5) =====
  {
    id: 'leng_diag_10',
    tipo: 'multiple_choice',
    dificultad: 4,
    pregunta: 'En el análisis de un texto literario, ¿qué significa el concepto de "intertextualidad"?',
    opciones: [
      'El uso de diferentes tipos de texto en una misma obra',
      'La relación y conexión entre diferentes textos literarios',
      'La estructura interna del texto',
      'El contexto histórico de producción del texto'
    ],
    respuestaCorrecta: 1,
    explicacion: 'La intertextualidad se refiere a las relaciones que un texto establece con otros textos, ya sea mediante citas, alusiones o referencias.',
    bloomLevel: 4
  },
  {
    id: 'leng_diag_11',
    tipo: 'fill_blanks',
    dificultad: 5,
    pregunta: 'El ___1___ es el tiempo verbal que se utiliza para narrar acciones que ya ocurrieron y están completamente finalizadas.',
    respuestaCorrecta: 'pretérito',
    explicacion: 'El pretérito (o pretérito perfecto simple) indica acciones pasadas y acabadas. Ejemplo: "Ella llegó ayer".',
    bloomLevel: 3
  },
  {
    id: 'leng_diag_12',
    tipo: 'multiple_choice',
    dificultad: 5,
    pregunta: 'En un ensayo académico, la estructura básica incluye (en orden):',
    opciones: [
      'Desarrollo - Introducción - Conclusión',
      'Introducción - Desarrollo - Conclusión',
      'Conclusión - Introducción - Desarrollo',
      'Introducción - Conclusión - Desarrollo'
    ],
    respuestaCorrecta: 1,
    explicacion: 'La estructura clásica del ensayo es: Introducción (presentación del tema), Desarrollo (argumentación) y Conclusión (síntesis).',
    bloomLevel: 5
  }
];

/**
 * Calculate diagnostic level based on score
 * @param correctAnswers Number of correct answers
 * @param totalQuestions Total number of questions
 * @returns Level 0-4
 */
export function calcularNivelDiagnostico(correctAnswers: number, totalQuestions: number): {
  nivel: number;
  label: string;
  porcentaje: number;
} {
  const porcentaje = Math.round((correctAnswers / totalQuestions) * 100);

  let nivel = 0;
  let label = 'Sin Dominio';

  if (porcentaje >= 86) {
    nivel = 4;
    label = 'Experto';
  } else if (porcentaje >= 71) {
    nivel = 3;
    label = 'Avanzado';
  } else if (porcentaje >= 51) {
    nivel = 2;
    label = 'Intermedio';
  } else if (porcentaje >= 26) {
    nivel = 1;
    label = 'Básico';
  } else {
    nivel = 0;
    label = 'Sin Dominio';
  }

  return { nivel, label, porcentaje };
}

/**
 * Get feedback message based on level
 */
export function getFeedbackMensaje(nivel: number): string {
  const mensajes: Record<number, string> = {
    0: 'No te preocupes, estamos aquí para ayudarte a fortalecer tus bases en Lengua y Literatura. Comenzaremos con conceptos fundamentales.',
    1: 'Tienes una base inicial en Lengua. Trabajaremos para fortalecer tus conocimientos y desarrollar nuevas habilidades.',
    2: 'Demuestras un buen manejo de Lengua y Literatura. Vamos a profundizar en conceptos más complejos y análisis crítico.',
    3: 'Excelente nivel en Lengua. Estás listo para desafíos más avanzados en análisis literario y producción de textos complejos.',
    4: '¡Sobresaliente! Tienes un dominio experto de Lengua y Literatura. Te propondremos actividades de alto nivel y pensamiento crítico.'
  };

  return mensajes[nivel] || mensajes[0];
}
