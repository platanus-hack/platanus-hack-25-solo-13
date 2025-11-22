import { browser } from '$app/environment';
import { goto } from '$app/navigation';

export const ssr = false;

export async function load() {
  if (browser) {
    const token = localStorage.getItem('auth_token');

    if (!token) {
      goto('/login');
      return {};
    }
  }

  return {};
}
