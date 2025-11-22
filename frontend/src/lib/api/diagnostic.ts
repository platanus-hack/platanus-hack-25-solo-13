import { auth } from '$lib/stores/auth.svelte';

export interface DiagnosticSession {
  id: number;
  user_id: number;
  materia_id: number;
  numero_intento: number;
  estado: string;
  estrategia: any;
  preguntas_totales: number;
  preguntas_correctas: number;
  started_at: string;
  completed_at?: string;
}

export interface DiagnosticQuestion {
  id: number;
  oa_bloom_objective_id: number;
  tipo: string;
  question_data: any;
  question_number: number;
  total_questions: number;
  current_bloom_level: number;
}

export interface DiagnosticAnswer {
  is_correct: boolean;
  score: number;
  answer_id: number;
  new_bloom_level: number;
}

export interface DiagnosticResult {
  id: number;
  session_id: number;
  oa_id: number;
  nivel_bloom_dominado: number;
  nivel_bloom_nombre: string;
  preguntas_respondidas: number;
  preguntas_correctas: number;
  porcentaje_aciertos: number;
  recomendacion: string;
  oa?: {
    id: number;
    codigo: string;
    titulo: string;
    descripcion: string;
  };
}

export interface DiagnosticCompleteResponse {
  message: string;
  session: DiagnosticSession;
  average_bloom_level: number;
}

const API_BASE_URL = '/api';

function getAuthHeaders() {
  const token = auth.token;
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  };
}

export async function startDiagnosticSession(materiaId: number) {
  try {
    const response = await fetch(`${API_BASE_URL}/diagnostic-sessions`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({ materia_id: materiaId })
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const session: DiagnosticSession = await response.json();
    return { success: true, session };
  } catch (error) {
    console.error('Error starting diagnostic session:', error);
    return { success: false, error: 'Network error' };
  }
}

export async function getNextQuestion(sessionId: number) {
  try {
    const response = await fetch(`${API_BASE_URL}/diagnostic-sessions/${sessionId}/next-question`, {
      method: 'GET',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const question: DiagnosticQuestion = await response.json();
    return { success: true, question };
  } catch (error) {
    console.error('Error getting next question:', error);
    return { success: false, error: 'Network error' };
  }
}

export async function submitAnswer(sessionId: number, questionId: number, userAnswer: any, tiempoSegundos?: number) {
  try {
    const response = await fetch(`${API_BASE_URL}/diagnostic-sessions/${sessionId}/answer`, {
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

    const result: DiagnosticAnswer = await response.json();
    return { success: true, result };
  } catch (error) {
    console.error('Error submitting answer:', error);
    return { success: false, error: 'Network error' };
  }
}

export async function completeSession(sessionId: number) {
  try {
    const response = await fetch(`${API_BASE_URL}/diagnostic-sessions/${sessionId}/complete`, {
      method: 'POST',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const result: DiagnosticCompleteResponse = await response.json();
    return { success: true, result };
  } catch (error) {
    console.error('Error completing session:', error);
    return { success: false, error: 'Network error' };
  }
}

export async function getResults(sessionId: number) {
  try {
    const response = await fetch(`${API_BASE_URL}/diagnostic-sessions/${sessionId}/results`, {
      method: 'GET',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const results: DiagnosticResult[] = await response.json();
    return { success: true, results };
  } catch (error) {
    console.error('Error getting results:', error);
    return { success: false, error: 'Network error' };
  }
}

export async function getSessionProgress(sessionId: number) {
  try {
    const response = await fetch(`${API_BASE_URL}/diagnostic-sessions/${sessionId}`, {
      method: 'GET',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const error = await response.text();
      return { success: false, error };
    }

    const session: DiagnosticSession = await response.json();
    return { success: true, session };
  } catch (error) {
    console.error('Error getting session progress:', error);
    return { success: false, error: 'Network error' };
  }
}
