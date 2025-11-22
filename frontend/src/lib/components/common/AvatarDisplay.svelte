<script lang="ts">
  import type { CustomizationItem } from '$lib/api/customization';

  interface Props {
    currentAvatar: CustomizationItem | null;
    initials: string;
    size?: 'small' | 'medium' | 'large';
    editable?: boolean;
    loading?: boolean;
    onEdit?: () => void;
  }

  let {
    currentAvatar,
    initials,
    size = 'medium',
    editable = false,
    loading = false,
    onEdit
  }: Props = $props();

  // Size configurations
  const sizeClasses = {
    small: 'h-10 w-10 text-sm',
    medium: 'h-24 w-24 text-3xl',
    large: 'h-32 w-32 text-4xl'
  };

  const editButtonClasses = {
    small: 'h-6 w-6 bottom-0 right-0',
    medium: 'h-8 w-8 bottom-0 right-0',
    large: 'h-10 w-10 bottom-1 right-1'
  };

  const editIconClasses = {
    small: 'w-3 h-3',
    medium: 'w-4 h-4',
    large: 'w-5 h-5'
  };
</script>

<div class="relative inline-block">
  {#if loading}
    <div class="{sizeClasses[size]} rounded-full bg-canvas-800 animate-pulse"></div>
  {:else if currentAvatar}
    <div class="relative {sizeClasses[size]} rounded-full overflow-hidden bg-gradient-to-br from-lumera-500/20 to-focus-600/20 border-2 border-lumera-500/30 shadow-lg">
      <img
        src="http://localhost:8080{currentAvatar.image_url}"
        alt={currentAvatar.name}
        class="w-full h-full object-cover scale-125"
      />
    </div>
  {:else}
    <div class="inline-flex items-center justify-center {sizeClasses[size]} rounded-full bg-gradient-to-br from-lumera-500 to-focus-600 {sizeClasses[size].split(' ')[2]} font-bold text-white shadow-lg">
      {initials}
    </div>
  {/if}

  {#if editable && onEdit}
    <button
      onclick={onEdit}
      class="absolute {editButtonClasses[size]} rounded-full bg-lumera-500 hover:bg-lumera-600 text-white flex items-center justify-center shadow-lg transition-all hover:scale-110"
      title="Cambiar avatar"
    >
      <svg class="{editIconClasses[size]}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
      </svg>
    </button>
  {/if}
</div>
