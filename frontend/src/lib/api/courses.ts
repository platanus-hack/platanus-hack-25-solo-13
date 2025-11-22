import { auth } from '$lib/stores/auth.svelte';

export interface Materia {
  id: number;
  codigo: string;
  nombre: string;
  descripcion: string;
  color: string; // Hex color
  activo: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface Curso {
  id: number;
  codigo: string;
  nombre: string;
  nivel_educativo: string;
  descripcion: string;
  activo: boolean;
  created_at?: string;
  updated_at?: string;
  materias?: Materia[];
}

/**
 * Get all courses
 */
export async function getCursos(): Promise<{ success: boolean; cursos?: Curso[]; error?: string }> {
  try {
    const response = await fetch('/api/cursos', {
      method: 'GET',
      headers: auth.getAuthHeaders()
    });

    if (!response.ok) {
      return { success: false, error: 'Failed to fetch courses' };
    }

    const cursos = await response.json();
    return { success: true, cursos };
  } catch (error) {
    return { success: false, error: 'Network error' };
  }
}

/**
 * Get a specific course with its subjects
 */
export async function getCurso(cursoId: number): Promise<{ success: boolean; curso?: Curso; error?: string }> {
  try {
    const response = await fetch(`/api/cursos/${cursoId}`, {
      method: 'GET',
      headers: auth.getAuthHeaders()
    });

    if (!response.ok) {
      return { success: false, error: 'Course not found' };
    }

    const curso = await response.json();
    return { success: true, curso };
  } catch (error) {
    return { success: false, error: 'Network error' };
  }
}

/**
 * Get all subjects
 */
export async function getMaterias(): Promise<{ success: boolean; materias?: Materia[]; error?: string }> {
  try {
    const response = await fetch('/api/materias?activo=true', {
      method: 'GET',
      headers: auth.getAuthHeaders()
    });

    if (!response.ok) {
      return { success: false, error: 'Failed to fetch subjects' };
    }

    const materias = await response.json();
    return { success: true, materias };
  } catch (error) {
    return { success: false, error: 'Network error' };
  }
}

/**
 * Find course by name (e.g., "1ro Medio", "Primero Medio")
 */
export async function findCursoByName(courseName: string): Promise<{ success: boolean; curso?: Curso; error?: string }> {
  const result = await getCursos();
  if (!result.success || !result.cursos) {
    return { success: false, error: result.error };
  }

  // Normalize course name for matching
  const normalized = courseName.toLowerCase().trim();

  // Try to find by exact match first
  let curso = result.cursos.find(c => c.nombre.toLowerCase() === normalized);

  // Try partial matches
  if (!curso) {
    if (normalized.includes('1') || normalized.includes('primero')) {
      curso = result.cursos.find(c => c.codigo === '1M' || c.nombre.toLowerCase().includes('primero'));
    } else if (normalized.includes('2') || normalized.includes('segundo')) {
      curso = result.cursos.find(c => c.codigo === '2M' || c.nombre.toLowerCase().includes('segundo'));
    } else if (normalized.includes('3') || normalized.includes('tercero')) {
      curso = result.cursos.find(c => c.codigo === '3M' || c.nombre.toLowerCase().includes('tercero'));
    } else if (normalized.includes('4') || normalized.includes('cuarto')) {
      curso = result.cursos.find(c => c.codigo === '4M' || c.nombre.toLowerCase().includes('cuarto'));
    }
  }

  if (!curso) {
    return { success: false, error: 'Course not found' };
  }

  // Load full course data with subjects
  return getCurso(curso.id);
}
