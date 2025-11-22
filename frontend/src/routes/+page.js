import { browser } from '$app/environment';
import { goto } from '$app/navigation';

export const ssr = false;

export async function load() {
  if (browser) {
    const token = localStorage.getItem('auth_token');
    const userStr = localStorage.getItem('auth_user');

    // Check authentication
    if (!token) {
      goto('/login');
      return {};
    }

    // Check if user has profile
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        const response = await fetch(`/api/profiles/${user.id}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        // If profile doesn't exist (404), redirect to onboarding
        if (response.status === 404) {
          goto('/onboarding');
          return {};
        }

        // Check if profile has conocimiento_previo (diagnostic completed)
        if (response.ok) {
          const profile = await response.json();
          const hasConocimientoPrevio = profile?.profile_data?.conocimiento_previo?.lectura?.nivel !== undefined;

          // If profile exists but no diagnostic, redirect to diagnostic
          if (!hasConocimientoPrevio) {
            goto('/diagnostico');
            return {};
          }
        }
      } catch (error) {
        console.error('Error checking profile:', error);
      }
    }
  }

  return {};
}
