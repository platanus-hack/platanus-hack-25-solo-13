/**
 * API client for Practice Sessions
 * Similar to diagnostic but focuses on a single OA with adaptive difficulty
 */

import { auth } from '$lib/stores/auth.svelte';

export interface PracticeSession {
  id: number;
  user_id: number;
  oa_id: number;
  oa_bloom_objective_id: number;
  bloom_level_inicial: number;
  bloom_level_final?: number;
  numero_preguntas: number;
  preguntas_respondidas: number;
  preguntas_correctas: number;
  estado: 'en_progreso' | 'completado';
  estrategia: any;
  resultado?: any;
  started_at: string;
  completed_at?: string;
  oa?: any;
}

export interface PracticeQuestion {
  id: number;
  oa_bloom_objective_id: number;
  tipo: string;
  question_data: any;
  question_number: number;
  total_questions: number;
  current_bloom_level: number;
}

export interface PracticeResult {
  session: PracticeSession;
  bloom_level_inicial: number;
  bloom_level_final: number;
  cambio_nivel: number;
  resultado: {
    porcentaje_aciertos: number;
    preguntas_totales: number;
    preguntas_correctas: number;
    aciertos_por_nivel: Record<number, number>;
    fallos_por_nivel: Record<number, number>;
    patron_respuestas: string[];
  };
}

function getAuthHeaders() {
  const token = auth.token;
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  };
}

/**
 * Start a new practice session for a specific OA
 */
export async function startPracticeSession(
  oaId: number,
  oaBloomObjectiveId: number,
  numeroPreguntas: number = 10
): Promise<{ success: boolean; session?: PracticeSession; error?: string }> {
  try {
    const response = await fetch('/api/practice-sessions', {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({
        oa_id: oaId,
        oa_bloom_objective_id: oaBloomObjectiveId,
        numero_preguntas: numeroPreguntas
      })
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const session = await response.json();
    return { success: true, session };
  } catch (error) {
    console.error('Error starting practice session:', error);
    return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
  }
}

/**
 * Get list of practice sessions with optional filters
 */
export async function getPracticeSessions(filters?: {
  oa_id?: number;
  estado?: 'en_progreso' | 'completado';
}): Promise<{ success: boolean; sessions?: PracticeSession[]; error?: string }> {
  try {
    const params = new URLSearchParams();
    if (filters?.oa_id) params.append('oa_id', filters.oa_id.toString());
    if (filters?.estado) params.append('estado', filters.estado);

    const response = await fetch(`/api/practice-sessions?${params.toString()}`, {
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const sessions = await response.json();
    return { success: true, sessions };
  } catch (error) {
    console.error('Error fetching practice sessions:', error);
    return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
  }
}

/**
 * Get next question in practice session with adaptive difficulty
 */
export async function getNextPracticeQuestion(
  sessionId: number
): Promise<{ success: boolean; question?: PracticeQuestion; error?: string }> {
  try {
    const response = await fetch(`/api/practice-sessions/${sessionId}/next-question`, {
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const question = await response.json();
    return { success: true, question };
  } catch (error) {
    console.error('Error fetching next question:', error);
    return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
  }
}

/**
 * Submit an answer to a practice question
 */
export async function submitPracticeAnswer(
  sessionId: number,
  questionId: number,
  userAnswer: any,
  tiempoSegundos?: number
): Promise<{
  success: boolean;
  is_correct?: boolean;
  score?: number;
  new_bloom_level?: number;
  is_complete?: boolean;
  error?: string;
}> {
  try {
    const response = await fetch(`/api/practice-sessions/${sessionId}/answer`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({
        question_id: questionId,
        user_answer: userAnswer,
        tiempo_segundos: tiempoSegundos
      })
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const result = await response.json();
    return { success: true, ...result };
  } catch (error) {
    console.error('Error submitting answer:', error);
    return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
  }
}

/**
 * Complete practice session and get final results
 */
export async function completePracticeSession(
  sessionId: number
): Promise<{ success: boolean; result?: PracticeResult; error?: string }> {
  try {
    const response = await fetch(`/api/practice-sessions/${sessionId}/complete`, {
      method: 'POST',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const result = await response.json();
    return { success: true, result };
  } catch (error) {
    console.error('Error completing practice session:', error);
    return { success: false, error: error instanceof Error ? error.message : 'Unknown error' };
  }
}
