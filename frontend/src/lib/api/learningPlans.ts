import { auth } from '$lib/stores/auth.svelte';

// ============================================================================
// Types
// ============================================================================

export interface LearningPlanComponent {
  id: number;
  learning_plan_id: number;
  orden: number;
  tipo_componente: string;
  objetivo_especifico: string;
  tiempo_estimado_minutos: number;
  estado: 'pendiente' | 'generando' | 'generado' | 'error';
  contenido_props: any | null;
}

export interface LearningPlan {
  id: number;
  user_id: number;
  oa_bloom_objective_id: number;
  titulo: string;
  descripcion: string;
  tiempo_estimado_minutos: number;
  estado: 'generando' | 'generado' | 'error';
  components?: LearningPlanComponent[];
  created_at?: string;
  updated_at?: string;
  // Completion tracking fields
  completado?: boolean;
  fecha_inicio?: string;
  fecha_completado?: string;
  progreso_actual?: number;
  total_slides?: number;
}

export interface OABloomObjective {
  id: number;
  oa_id: number;
  bloom_level_id: number;
  objetivo_especifico: string;
  indicadores_logro: string[];
  tipo_actividad_sugerida: string;
  complejidad_estimada: number;
}

// ============================================================================
// Helper Functions
// ============================================================================

function getAuthHeaders() {
  const token = auth.token;
  return {
    'Content-Type': 'application/json',
    'Authorization': token ? `Bearer ${token}` : ''
  };
}

// ============================================================================
// API Functions
// ============================================================================

/**
 * Get the OA Bloom Objective ID for a specific OA and Bloom level
 * This queries the oa_bloom_objectives table to find the objective ID
 */
export async function getOABloomObjectiveId(oaId: number, bloomLevel: number): Promise<number | null> {
  try {
    // First, get all bloom objectives for this OA
    const response = await fetch(`/api/objetivos-aprendizaje/${oaId}`, {
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      console.error('Failed to fetch OA:', response.statusText);
      return null;
    }

    const oa = await response.json();

    // Find the bloom objective that matches the bloom level
    const bloomObjective = oa.bloom_objectives?.find(
      (obj: OABloomObjective) => obj.bloom_level_id === bloomLevel
    );

    return bloomObjective?.id || null;
  } catch (error) {
    console.error('Error getting OA Bloom Objective ID:', error);
    return null;
  }
}

/**
 * Get learning plan by OA Bloom Objective ID
 * Returns the plan if it exists, null if it doesn't
 */
export async function getPlanByOA(oaBloomObjectiveId: number): Promise<LearningPlan | null> {
  try {
    const response = await fetch(`/api/learning-plans/by-oa/${oaBloomObjectiveId}`, {
      headers: getAuthHeaders()
    });

    if (response.status === 404) {
      return null; // Plan doesn't exist yet - this is normal
    }

    if (!response.ok) {
      // Silently return null for other errors (401, 500, etc)
      return null;
    }

    return await response.json();
  } catch (error) {
    // Silently handle errors - plan doesn't exist or network issue
    return null;
  }
}

/**
 * Generate a new learning plan for a specific OA Bloom Objective
 * This will create the plan structure (without content yet)
 */
export async function generatePlan(oaBloomObjectiveId: number): Promise<{
  success: boolean;
  plan?: LearningPlan;
  error?: string;
}> {
  try {
    const response = await fetch('/api/learning-plans/generate', {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify({
        oa_bloom_objective_id: oaBloomObjectiveId
      })
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      return {
        success: false,
        error: errorData?.error || `HTTP ${response.status}: ${response.statusText}`
      };
    }

    const plan = await response.json();
    return { success: true, plan };
  } catch (error) {
    console.error('Error generating plan:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

/**
 * Generate content for a specific component in a learning plan
 * This lazily loads the OpenAI-generated content when the user views the slide
 */
export async function generateComponentContent(
  planId: number,
  componentId: number
): Promise<{
  success: boolean;
  component?: LearningPlanComponent;
  error?: string;
}> {
  try {
    const response = await fetch(
      `/api/learning-plans/${planId}/components/${componentId}/generate-content`,
      {
        method: 'POST',
        headers: getAuthHeaders()
      }
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      return {
        success: false,
        error: errorData?.error || `HTTP ${response.status}: ${response.statusText}`
      };
    }

    const component = await response.json();
    return { success: true, component };
  } catch (error) {
    console.error('Error generating component content:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

/**
 * Get a learning plan by ID
 * Useful for reloading a plan or viewing it directly via URL
 */
export async function getPlanById(planId: number): Promise<{
  success: boolean;
  plan?: LearningPlan;
  error?: string;
}> {
  try {
    const response = await fetch(`/api/learning-plans/${planId}`, {
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      return {
        success: false,
        error: `HTTP ${response.status}: ${response.statusText}`
      };
    }

    const plan = await response.json();
    return { success: true, plan };
  } catch (error) {
    console.error('Error fetching plan:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

/**
 * Get all OAs for a specific materia
 */
export async function getObjetivosAprendizaje(materiaId: number): Promise<{
  success: boolean;
  oas?: any[];
  error?: string;
}> {
  try {
    const response = await fetch(`/api/objetivos-aprendizaje?materia_id=${materiaId}&activo=true`, {
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      return {
        success: false,
        error: `HTTP ${response.status}: ${response.statusText}`
      };
    }

    const oas = await response.json();
    return { success: true, oas };
  } catch (error) {
    console.error('Error fetching OAs:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

/**
 * Mark a learning plan as started
 * Sets the fecha_inicio timestamp
 */
export async function startLearningPlan(planId: number): Promise<{
  success: boolean;
  plan?: LearningPlan;
  error?: string;
}> {
  try {
    const response = await fetch(`/api/learning-plans/${planId}/start`, {
      method: 'POST',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      return {
        success: false,
        error: errorData?.error || `HTTP ${response.status}: ${response.statusText}`
      };
    }

    const plan = await response.json();
    return { success: true, plan };
  } catch (error) {
    console.error('Error starting plan:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}

/**
 * Mark a learning plan as completed
 * Sets completado = true, fecha_completado timestamp, and progreso_actual = total_slides
 */
export async function completeLearningPlan(planId: number): Promise<{
  success: boolean;
  plan?: LearningPlan;
  error?: string;
}> {
  try {
    const response = await fetch(`/api/learning-plans/${planId}/complete`, {
      method: 'POST',
      headers: getAuthHeaders()
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => null);
      return {
        success: false,
        error: errorData?.error || `HTTP ${response.status}: ${response.statusText}`
      };
    }

    const plan = await response.json();
    return { success: true, plan };
  } catch (error) {
    console.error('Error completing plan:', error);
    return {
      success: false,
      error: error instanceof Error ? error.message : 'Unknown error'
    };
  }
}
