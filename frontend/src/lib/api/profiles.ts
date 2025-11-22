import { auth } from '$lib/stores/auth.svelte';

export interface ProfileData {
  conocimiento_previo?: {
    [key: string]: {
      nivel: number;
      fuente: string;
    };
  };
  perfil_cognitivo?: {
    memoria_trabajo?: number;
    razonamiento_inductivo?: number;
    estilo_cognitivo?: string;
    carga_cognitiva_tolerada?: string;
  };
  preferencias_aprendizaje?: {
    formato_preferido?: string;
    tipo_actividad?: string[];
    canal_preferido?: string;
  };
  motivacion?: {
    intrinseca?: number;
    extrinseca?: number;
    interes_actual?: string;
    orientacion_meta?: string;
  };
  autoeficacia?: {
    general?: number;
    confianza_resolutiva?: number;
  };
  autonomia?: {
    nivel?: string;
    gestiona_tiempo?: boolean;
    estrategias?: string[];
  };
  intereses_personales?: {
    temas?: string[];
    profesion_soñada?: string;
  };
  ultima_actualizacion?: string;
}

export interface StudentProfile {
  id: number;
  user_id: number;
  edad?: number;
  curso_actual?: string;
  profile_data: ProfileData;
  created_at: string;
  updated_at: string;
}

export interface CreateProfileRequest {
  user_id: number;
  edad?: number;
  curso_actual?: string;
  profile_data: ProfileData;
}

export interface UpdateProfileRequest {
  edad?: number;
  curso_actual?: string;
  profile_data?: ProfileData;
}

/**
 * Get student profile by user ID
 */
export async function getProfile(userId: number): Promise<{ success: boolean; profile?: StudentProfile; error?: string }> {
  try {
    const response = await fetch(`/api/profiles/${userId}`, {
      method: 'GET',
      headers: auth.getAuthHeaders()
    });

    if (response.status === 404) {
      return { success: false, error: 'Profile not found' };
    }

    if (!response.ok) {
      const error = await response.json();
      return { success: false, error: error.error || 'Failed to get profile' };
    }

    const profile = await response.json();
    return { success: true, profile };
  } catch (error) {
    return { success: false, error: 'Network error. Please try again.' };
  }
}

/**
 * Create new student profile
 */
export async function createProfile(data: CreateProfileRequest): Promise<{ success: boolean; profile?: StudentProfile; error?: string }> {
  try {
    const response = await fetch('/api/profiles', {
      method: 'POST',
      headers: auth.getAuthHeaders(),
      body: JSON.stringify(data)
    });

    if (response.status === 409) {
      return { success: false, error: 'Profile already exists' };
    }

    if (!response.ok) {
      const error = await response.json();
      return { success: false, error: error.error || 'Failed to create profile' };
    }

    const profile = await response.json();
    return { success: true, profile };
  } catch (error) {
    return { success: false, error: 'Network error. Please try again.' };
  }
}

/**
 * Update student profile (partial update)
 * IMPORTANT: profile_data is completely replaced, not merged
 */
export async function updateProfile(userId: number, data: UpdateProfileRequest): Promise<{ success: boolean; profile?: StudentProfile; error?: string }> {
  try {
    const response = await fetch(`/api/profiles/${userId}`, {
      method: 'PATCH',
      headers: auth.getAuthHeaders(),
      body: JSON.stringify(data)
    });

    if (response.status === 404) {
      return { success: false, error: 'Profile not found' };
    }

    if (!response.ok) {
      const error = await response.json();
      return { success: false, error: error.error || 'Failed to update profile' };
    }

    const profile = await response.json();
    return { success: true, profile };
  } catch (error) {
    return { success: false, error: 'Network error. Please try again.' };
  }
}

/**
 * Check if user has a profile
 */
export async function hasProfile(userId: number): Promise<boolean> {
  const result = await getProfile(userId);
  return result.success && result.profile !== undefined;
}

/**
 * Update profile_data with merge (preserves existing data)
 * This is a safer alternative to updateProfile when you only want to update specific fields
 */
export async function mergeProfileData(userId: number, newProfileData: Partial<ProfileData>): Promise<{ success: boolean; profile?: StudentProfile; error?: string }> {
  try {
    // First get the current profile
    const currentProfile = await getProfile(userId);
    if (!currentProfile.success || !currentProfile.profile) {
      return { success: false, error: currentProfile.error || 'Profile not found' };
    }

    // Deep merge profile_data
    const mergedProfileData = {
      ...currentProfile.profile.profile_data,
      ...newProfileData
    };

    // Merge nested objects specifically
    if (newProfileData.conocimiento_previo) {
      mergedProfileData.conocimiento_previo = {
        ...currentProfile.profile.profile_data.conocimiento_previo,
        ...newProfileData.conocimiento_previo
      };
    }

    if (newProfileData.preferencias_aprendizaje) {
      mergedProfileData.preferencias_aprendizaje = {
        ...currentProfile.profile.profile_data.preferencias_aprendizaje,
        ...newProfileData.preferencias_aprendizaje
      };
    }

    if (newProfileData.motivacion) {
      mergedProfileData.motivacion = {
        ...currentProfile.profile.profile_data.motivacion,
        ...newProfileData.motivacion
      };
    }

    if (newProfileData.intereses_personales) {
      mergedProfileData.intereses_personales = {
        ...currentProfile.profile.profile_data.intereses_personales,
        ...newProfileData.intereses_personales
      };
    }

    // Update with merged data
    return await updateProfile(userId, { profile_data: mergedProfileData });
  } catch (error) {
    return { success: false, error: 'Error merging profile data' };
  }
}
