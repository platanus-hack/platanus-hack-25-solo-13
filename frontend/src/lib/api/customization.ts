import { auth } from '$lib/stores/auth.svelte';

const API_BASE = 'http://localhost:8080/api';

export interface CustomizationItem {
  id: number;
  name: string;
  type: string;
  rarity: string;
  tier: number;
  image_url: string;
  description: string;
  base_coins_cost: number;
  is_default: boolean;
  is_owned: boolean;
  can_purchase: boolean;
  is_equipped: boolean;
  status: 'owned' | 'locked' | 'can_purchase';
}

export interface UserEquipment {
  user_id: number;
  equipped_avatar_id: number | null;
  equipped_frame_id: number | null;
  equipped_avatar?: CustomizationItem;
  equipped_frame?: CustomizationItem;
}

export interface UserInventoryItem {
  id: number;
  user_id: number;
  item_id: number;
  item: CustomizationItem;
  unlocked_at: string;
}

/**
 * Get user's customization inventory (owned items)
 */
export async function getInventory(): Promise<UserInventoryItem[]> {
  const response = await fetch(`${API_BASE}/customization/inventory`, {
    headers: {
      'Authorization': `Bearer ${auth.token}`
    }
  });

  if (!response.ok) {
    throw new Error('Failed to fetch inventory');
  }

  return response.json();
}

/**
 * Get user's equipped items
 */
export async function getEquipment(): Promise<UserEquipment> {
  const response = await fetch(`${API_BASE}/customization/equipment`, {
    headers: {
      'Authorization': `Bearer ${auth.token}`
    }
  });

  if (!response.ok) {
    throw new Error('Failed to fetch equipment');
  }

  return response.json();
}

/**
 * Get full customization catalog with ownership status
 */
export async function getCatalog(): Promise<CustomizationItem[]> {
  const response = await fetch(`${API_BASE}/customization/catalog`, {
    headers: {
      'Authorization': `Bearer ${auth.token}`
    }
  });

  if (!response.ok) {
    throw new Error('Failed to fetch catalog');
  }

  return response.json();
}

/**
 * Equip an avatar or frame
 */
export async function equipItem(itemId: number, slot: 'avatar' | 'frame'): Promise<UserEquipment> {
  const response = await fetch(`${API_BASE}/customization/equip`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${auth.token}`
    },
    body: JSON.stringify({
      item_id: itemId,
      slot: slot
    })
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to equip item');
  }

  return response.json();
}

/**
 * Purchase an item with coins
 */
export async function purchaseItem(itemId: number): Promise<UserInventoryItem> {
  const response = await fetch(`${API_BASE}/customization/purchase`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${auth.token}`
    },
    body: JSON.stringify({
      item_id: itemId
    })
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to purchase item');
  }

  return response.json();
}
