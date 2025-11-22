import { browser } from '$app/environment';
import { goto } from '$app/navigation';

interface User {
  id: number;
  email: string;
  name: string;
  role: string;
  created_at: string;
  updated_at: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

// Create reactive auth state
let authState = $state<AuthState>({
  user: null,
  token: null,
  isAuthenticated: false
});

// Initialize from localStorage on client
if (browser) {
  const savedToken = localStorage.getItem('auth_token');
  const savedUser = localStorage.getItem('auth_user');

  if (savedToken && savedUser) {
    authState.token = savedToken;
    authState.user = JSON.parse(savedUser);
    authState.isAuthenticated = true;
  }
}

export const auth = {
  // Getters (reactive)
  get user() { return authState.user; },
  get token() { return authState.token; },
  get isAuthenticated() { return authState.isAuthenticated; },

  // Register
  async register(email: string, name: string, password: string): Promise<{ success: boolean; error?: string }> {
    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, name, password })
      });

      if (!response.ok) {
        const error = await response.json();
        return { success: false, error: error.error || 'Registration failed' };
      }

      const data = await response.json();
      this.setAuth(data.token, data.user);
      return { success: true };
    } catch (error) {
      return { success: false, error: 'Network error. Please try again.' };
    }
  },

  // Login
  async login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      if (!response.ok) {
        const error = await response.json();
        return { success: false, error: error.error || 'Invalid credentials' };
      }

      const data = await response.json();
      this.setAuth(data.token, data.user);
      return { success: true };
    } catch (error) {
      return { success: false, error: 'Network error. Please try again.' };
    }
  },

  // Logout
  logout() {
    authState.user = null;
    authState.token = null;
    authState.isAuthenticated = false;

    if (browser) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
      goto('/login');
    }
  },

  // Set authentication
  setAuth(token: string, user: User) {
    authState.token = token;
    authState.user = user;
    authState.isAuthenticated = true;

    if (browser) {
      localStorage.setItem('auth_token', token);
      localStorage.setItem('auth_user', JSON.stringify(user));
    }
  },

  // Check if user is authenticated
  checkAuth(): boolean {
    return authState.isAuthenticated && authState.token !== null;
  },

  // Get auth headers for API calls
  getAuthHeaders(): HeadersInit {
    if (authState.token) {
      return {
        'Authorization': `Bearer ${authState.token}`,
        'Content-Type': 'application/json'
      };
    }
    return { 'Content-Type': 'application/json' };
  },

  // Check if user has a profile
  async checkIfHasProfile(): Promise<boolean> {
    if (!authState.user?.id) {
      return false;
    }

    try {
      const response = await fetch(`/api/profiles/${authState.user.id}`, {
        method: 'GET',
        headers: this.getAuthHeaders()
      });

      return response.ok;
    } catch (error) {
      return false;
    }
  }
};
