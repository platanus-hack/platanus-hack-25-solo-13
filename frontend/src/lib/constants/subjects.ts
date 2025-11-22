import type { Materia } from '$lib/api/courses';

export interface Subject {
  id: string; // Código de la materia (MAT, LYL, etc.)
  name: string;
  icon: string;
  color: string; // Tailwind gradient
  description: string;
}

// Mapeo de códigos de materia a íconos (visual only, hardcoded)
const SUBJECT_ICONS: Record<string, string> = {
  'MAT': '📐',
  'MATEMATICAS': '📐',
  'LYL': '📖',
  'LENGUAJE': '📖',
  'LENGUA': '📖',
  'HIS': '🌍',
  'HISTORIA': '🌍',
  'CNA': '🔬',
  'CIENCIAS': '🔬',
  'FIS': '⚛️',
  'FISICA': '⚛️',
  'QUI': '⚗️',
  'QUIMICA': '⚗️',
  'BIO': '🧬',
  'BIOLOGIA': '🧬',
  'ING': '🗣️',
  'INGLES': '🗣️',
  'FIL': '🤔',
  'FILOSOFIA': '🤔',
  'ART': '🎨',
  'ARTES': '🎨',
  'MUS': '🎵',
  'MUSICA': '🎵',
  'EDF': '⚽',
  'ED_FISICA': '⚽'
};

// Mapeo de colores hex a gradientes de Tailwind
function hexToTailwindGradient(hexColor: string): string {
  // Extract base color from hex
  const colorMap: Record<string, string> = {
    '#EF4444': 'from-red-500 to-rose-500',
    '#F87171': 'from-red-400 to-red-500',
    '#DC2626': 'from-red-600 to-red-700',

    '#3B82F6': 'from-blue-500 to-cyan-500',
    '#60A5FA': 'from-blue-400 to-blue-500',
    '#2563EB': 'from-blue-600 to-indigo-600',

    '#10B981': 'from-emerald-500 to-green-500',
    '#34D399': 'from-emerald-400 to-green-400',
    '#059669': 'from-emerald-600 to-green-600',

    '#8B5CF6': 'from-violet-500 to-purple-500',
    '#A78BFA': 'from-violet-400 to-purple-400',
    '#7C3AED': 'from-violet-600 to-purple-600',

    '#F59E0B': 'from-amber-500 to-orange-500',
    '#FBBF24': 'from-amber-400 to-yellow-400',
    '#D97706': 'from-amber-600 to-orange-600',

    '#EC4899': 'from-pink-500 to-rose-500',
    '#F472B6': 'from-pink-400 to-rose-400',
    '#DB2777': 'from-pink-600 to-rose-600',

    '#14B8A6': 'from-teal-500 to-cyan-500',
    '#2DD4BF': 'from-teal-400 to-cyan-400',
    '#0D9488': 'from-teal-600 to-cyan-600',

    '#6366F1': 'from-indigo-500 to-purple-500',
    '#818CF8': 'from-indigo-400 to-purple-400',
    '#4F46E5': 'from-indigo-600 to-purple-600'
  };

  return colorMap[hexColor.toUpperCase()] || 'from-slate-500 to-gray-600';
}

// Get icon for subject by code
function getSubjectIcon(codigo: string, nombre: string): string {
  const key = codigo.toUpperCase();
  if (SUBJECT_ICONS[key]) {
    return SUBJECT_ICONS[key];
  }

  // Try by name
  const nameKey = nombre.toUpperCase().replace(/\s+/g, '_').replace(/[ÁÀ]/g, 'A').replace(/[ÉÈ]/g, 'E').replace(/[ÍÌ]/g, 'I').replace(/[ÓÒ]/g, 'O').replace(/[ÚÙ]/g, 'U');
  return SUBJECT_ICONS[nameKey] || '📚'; // Default icon
}

/**
 * Transform backend Materia to frontend Subject
 */
export function materiaToSubject(materia: Materia): Subject {
  return {
    id: materia.codigo.toLowerCase(),
    name: materia.nombre,
    icon: getSubjectIcon(materia.codigo, materia.nombre),
    color: hexToTailwindGradient(materia.color),
    description: materia.descripcion || ''
  };
}

/**
 * Transform multiple materias to subjects
 */
export function materiasToSubjects(materias: Materia[]): Subject[] {
  return materias.map(materiaToSubject);
}

// ============================================================================
// DOMAIN LEVELS (unchanged)
// ============================================================================

export enum DomainLevel {
  NOT_EVALUATED = 0,
  BASIC = 2,
  INTERMEDIATE = 3,
  ADVANCED = 4
}

export interface DomainLevelInfo {
  level: DomainLevel;
  label: string;
  color: string;
  badgeColor: string;
  textColor: string;
}

export const DOMAIN_LEVELS: Record<number, DomainLevelInfo> = {
  0: {
    level: DomainLevel.NOT_EVALUATED,
    label: 'No Evaluado',
    color: 'from-slate-600 to-slate-700',
    badgeColor: 'bg-slate-600',
    textColor: 'text-slate-300'
  },
  1: {
    level: DomainLevel.NOT_EVALUATED,
    label: 'Sin Dominio',
    color: 'from-red-600 to-rose-600',
    badgeColor: 'bg-red-600',
    textColor: 'text-red-300'
  },
  2: {
    level: DomainLevel.BASIC,
    label: 'Básico',
    color: 'from-yellow-500 to-amber-500',
    badgeColor: 'bg-yellow-500',
    textColor: 'text-yellow-300'
  },
  3: {
    level: DomainLevel.INTERMEDIATE,
    label: 'Intermedio',
    color: 'from-green-500 to-emerald-500',
    badgeColor: 'bg-green-500',
    textColor: 'text-green-300'
  },
  4: {
    level: DomainLevel.ADVANCED,
    label: 'Avanzado',
    color: 'from-blue-500 to-indigo-500',
    badgeColor: 'bg-blue-500',
    textColor: 'text-blue-300'
  }
};

/**
 * Get domain level info
 */
export function getDomainLevelInfo(level: number): DomainLevelInfo {
  return DOMAIN_LEVELS[level] || DOMAIN_LEVELS[0];
}
