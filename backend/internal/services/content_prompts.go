package services

import "fmt"

// buildPlanStructurePrompt construye el prompt para generar la estructura del plan
func buildPlanStructurePrompt(ctx OAContext) string {
	return fmt.Sprintf(`CONTEXTO EDUCATIVO:
- Materia: %s
- Curso: %s
- Objetivo de Aprendizaje (OA): %s
- Descripción del OA: %s
- Nivel de Bloom: %s (Nivel %d)
- Descripción del nivel de Bloom: %s
- Objetivo Específico: %s
- Indicadores de Logro: %v

TAREA: Diseñar un plan de aprendizaje personalizado

Debes generar un plan de aprendizaje que guíe al estudiante hacia el dominio del objetivo de aprendizaje usando SCAFFOLDING PEDAGÓGICO (andamiaje).

COMPONENTE DISPONIBLE:
ExplainAndExploreSlide - Componente flexible con bloques de contenido variables

Este componente puede incluir cualquier combinación de:
- Bloques de texto explicativo (explicaciones profundas y fundacionales)
- Bloques de ejemplos (con análisis)
- Bloques de definiciones (términos clave)
- Bloques de notas destacadas (tips, advertencias, información adicional)
- Bloques de ejercicios (práctica guiada)
- Bloques de resumen (síntesis de conceptos)
- Bloques de comparación (tablas comparativas)

INSTRUCCIONES:
1. Analiza el objetivo de aprendizaje y su nivel de Bloom
2. IMPORTANTE: Diseña una progresión pedagógica que comience desde FUNDAMENTOS (niveles Bloom 1-2) y construya gradualmente hacia el nivel objetivo
3. Divide el aprendizaje en componentes ExplainAndExploreSlide secuenciales - crea TANTOS componentes como sean necesarios para cubrir el tema en profundidad
4. Cada tema, concepto o habilidad importante merece su propio componente dedicado - NO intentes comprimir múltiples conceptos complejos en un solo componente
5. Los primeros componentes deben enfocarse en enseñar BASES (conceptos fundamentales, definiciones, ejemplos simples)
6. Los componentes posteriores pueden aumentar complejidad gradualmente, dedicando tiempo suficiente a cada nivel de profundización
7. Estima el tiempo en minutos para cada componente (considerar que pueden tener mucho contenido - 10-15 min por componente es razonable)

IMPORTANTE - PROFUNDIDAD POR NIVEL DE BLOOM:
- TODOS los componentes deben ser tipo "ExplainAndExploreSlide"
- Niveles Bloom 1-2 (Recordar/Comprender): Plan más directo, enfocado en fundamentos
- Niveles Bloom 3-4 (Aplicar/Analizar): Plan más extenso que incluya múltiples ejemplos y casos de aplicación
- Niveles Bloom 5-6 (Evaluar/Crear): Plan completo y detallado con múltiples componentes que exploren diferentes aspectos, perspectivas y aplicaciones avanzadas
- Cada componente debe tener un objetivo específico claro que construya sobre el anterior
- SCAFFOLDING: Los primeros componentes enseñan fundamentos, los componentes intermedios desarrollan, los últimos profundizan y aplican
- NO asumas conocimiento previo en el primer componente
- Prioriza CALIDAD sobre brevedad - es mejor un plan completo que uno superficial

FORMATO DE RESPUESTA (JSON):
{
  "titulo": "Título atractivo del plan de aprendizaje",
  "descripcion": "Descripción breve de lo que el estudiante aprenderá, empezando desde fundamentos",
  "componentes": [
    {
      "tipo": "ExplainAndExploreSlide",
      "objetivo_especifico": "Qué aprenderá el estudiante con este componente específico",
      "tiempo_estimado_minutos": 15
    }
  ]
}

Responde ÚNICAMENTE con el JSON, sin texto adicional.`,
		ctx.MateriaNombre,
		ctx.CursoNombre,
		ctx.OATitulo,
		ctx.OADescripcion,
		ctx.BloomLevelNombre,
		ctx.BloomLevelNumero,
		ctx.BloomDescripcion,
		ctx.ObjetivoEspecifico,
		ctx.IndicadoresLogro,
	)
}

// buildComponentPrompt construye el prompt para ExplainAndExploreSlide con bloques tipados
func buildComponentPrompt(componentType string, ctx OAContext, componentObjective string) string {
	return fmt.Sprintf(`CONTEXTO EDUCATIVO:
- Materia: %s (%s)
- Curso: %s
- Objetivo de Aprendizaje (OA): %s
- Descripción del OA: %s
- Nivel de Bloom: %s (Nivel %d) - %s
- Objetivo Específico de ESTE componente: %s

INDICADORES DE LOGRO:
%v

TAREA: Generar contenido educativo usando BLOQUES FLEXIBLES

Debes crear contenido pedagógico que alterne entre explicaciones, ejemplos, definiciones y práctica según sea necesario.

TIPOS DE BLOQUES DISPONIBLES:

1. BLOQUE TEXTO (tipo: "texto")
   Uso: Explicaciones, introducciones, desarrollo de conceptos
   Campos: { tipo: "texto", contenido: "texto del párrafo" }
   Cuándo usar: Para explicaciones profundas, contexto, transiciones

2. BLOQUE EJEMPLO (tipo: "ejemplo")
   Uso: Ilustrar conceptos con casos concretos
   Campos: { tipo: "ejemplo", titulo: "...", contenido: "...", analisis: "..." }
   Cuándo usar: Después de explicar un concepto

3. BLOQUE DEFINICIÓN (tipo: "definicion")
   Uso: Términos clave que el estudiante debe conocer
   Campos: { tipo: "definicion", termino: "...", texto: "..." }
   Cuándo usar: Para vocabulario técnico o conceptos fundamentales

4. BLOQUE NOTA (tipo: "nota")
   Uso: Destacar información importante, tips, advertencias
   Campos: { tipo: "nota", estilo: "info"|"warning"|"tip", texto: "..." }
   Cuándo usar: Para enfatizar puntos clave

5. BLOQUE EJERCICIO (tipo: "ejercicio")
   Uso: Práctica guiada o actividades
   Campos: { tipo: "ejercicio", instruccion: "...", ejemplo: "..." (opcional) }
   Cuándo usar: Para aplicar lo aprendido

6. BLOQUE RESUMEN (tipo: "resumen")
   Uso: Síntesis de conceptos clave
   Campos: { tipo: "resumen", puntos: ["...", "...", "..."] }
   Cuándo usar: Al final o después de secciones extensas

7. BLOQUE COMPARACIÓN (tipo: "comparacion")
   Uso: Tablas comparativas entre conceptos
   Campos: { tipo: "comparacion", items: [{aspecto: "...", opcion1: "...", opcion2: "..."}] }
   Cuándo usar: Para contrastar conceptos similares

INSTRUCCIONES PEDAGÓGICAS:
1. COMIENZA SIMPLE: Asume que el estudiante puede no conocer el tema. Usa bloques "definicion" y "texto" para fundamentos
2. ALTERNA TEORÍA Y PRÁCTICA: Después de explicar (bloques "texto"), proporciona ejemplos (bloques "ejemplo")
3. PROFUNDIDAD VARIABLE: Los bloques "texto" pueden ser largos (2-4 párrafos) si se necesita explicar bien
4. USA MÚLTIPLES BLOQUES: No te limites. Genera contenido ABUNDANTE - cada componente debe ser rico en información, ejemplos y práctica
5. SECUENCIA LÓGICA: Sigue un flujo natural de aprendizaje
6. CANTIDAD: Apunta a crear al menos 6-10 bloques por componente para cubrir el tema en profundidad

SCAFFOLDING (ANDAMIAJE) - ADAPTAR CANTIDAD Y COMPLEJIDAD:
- Si el nivel Bloom es bajo (1-2): Enfócate en definiciones, explicaciones claras y muchos ejemplos. Mínimo 6-8 bloques
- Si el nivel Bloom es medio (3-4): Incluye múltiples ejercicios de aplicación, comparaciones y casos de uso. Mínimo 8-10 bloques
- Si el nivel Bloom es alto (5-6): Incluye ejercicios complejos, análisis críticos, múltiples perspectivas y aplicaciones avanzadas. Mínimo 10-12 bloques
- Recuerda: Es mejor un componente completo y profundo que uno superficial

FORMATO DE RESPUESTA (JSON):
{
  "titulo": "Título descriptivo del componente",
  "bloques": [
    {
      "tipo": "texto",
      "contenido": "Párrafo explicativo. Puede ser largo y detallado."
    },
    {
      "tipo": "definicion",
      "termino": "Término importante",
      "texto": "Definición clara del término"
    },
    {
      "tipo": "ejemplo",
      "titulo": "Nombre del ejemplo",
      "contenido": "Contenido del ejemplo",
      "analisis": "Explicación de qué ilustra este ejemplo"
    },
    {
      "tipo": "nota",
      "estilo": "tip",
      "texto": "Consejo útil para el estudiante"
    },
    {
      "tipo": "ejercicio",
      "instruccion": "Qué debe hacer el estudiante",
      "ejemplo": "Ejemplo de cómo hacerlo (opcional)"
    },
    {
      "tipo": "resumen",
      "puntos": [
        "Punto clave 1",
        "Punto clave 2",
        "Punto clave 3"
      ]
    },
    {
      "tipo": "comparacion",
      "items": [
        {
          "aspecto": "Característica a comparar",
          "opcion1": "Valor en opción 1",
          "opcion2": "Valor en opción 2"
        }
      ]
    }
  ]
}

IMPORTANTE:
- El array "bloques" puede tener tantos elementos como necesites
- Alterna tipos de bloques para mantener interés
- Usa bloques "texto" largos cuando necesites explicar fundamentos
- SIEMPRE valida que el JSON esté bien formado
- NO uses tipos de bloques que no estén en la lista
- El objetivo específico es: %s

Responde ÚNICAMENTE con el JSON, sin texto adicional.`,
		ctx.MateriaNombre,
		ctx.MateriaDescripcion,
		ctx.CursoNombre,
		ctx.OATitulo,
		ctx.OADescripcion,
		ctx.BloomLevelNombre,
		ctx.BloomLevelNumero,
		ctx.BloomDescripcion,
		componentObjective,
		ctx.IndicadoresLogro,
		componentObjective,
	)
}
