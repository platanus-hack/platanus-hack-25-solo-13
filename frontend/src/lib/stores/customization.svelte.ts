import { browser } from '$app/environment';
import { getEquipment, type CustomizationItem } from '$lib/api/customization';

interface CustomizationState {
  currentAvatar: CustomizationItem | null;
  isLoading: boolean;
}

// Create reactive customization state
let customizationState = $state<CustomizationState>({
  currentAvatar: null,
  isLoading: false
});

export const customizationStore = {
  // Getters (reactive)
  get currentAvatar() { return customizationState.currentAvatar; },
  get isLoading() { return customizationState.isLoading; },

  // Load user's equipped avatar
  async loadAvatar() {
    if (!browser) return;

    customizationState.isLoading = true;
    try {
      const equipment = await getEquipment();
      customizationState.currentAvatar = equipment.equipped_avatar || null;
    } catch (err) {
      console.error('Error loading avatar:', err);
      customizationState.currentAvatar = null;
    } finally {
      customizationState.isLoading = false;
    }
  },

  // Update avatar (when user changes it)
  setAvatar(avatar: CustomizationItem | null) {
    customizationState.currentAvatar = avatar;
  },

  // Clear avatar (on logout)
  clear() {
    customizationState.currentAvatar = null;
    customizationState.isLoading = false;
  }
};
