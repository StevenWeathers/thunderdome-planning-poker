<script lang="ts">
  import { EllipsisIcon } from '@lucide/svelte';
  import type { Component } from 'svelte';

  interface Action {
    label: string;
    icon: Component;
    onclick: () => void;
    className?: string;
    testId?: string;
    disabled?: boolean;
  }

  interface Props {
    Icon?: Component;
    actions?: Action[];
    ariaLabel?: string;
    testId?: string;
    disabled?: boolean;
    iconSize?: 'small' | 'medium' | 'large';
  }

  let {
    Icon = EllipsisIcon,
    actions = [],
    ariaLabel = 'Actions menu',
    testId = 'actions-menu-button',
    disabled = false,
    iconSize = 'small',
  }: Props = $props();

  let showActions = $state(false);

  function toggleActions() {
    showActions = !showActions;
  }

  function closeActions() {
    showActions = false;
  }

  function handleAction(action: Action) {
    if (action.disabled) return;
    action.onclick();
    closeActions();
  }

  function clickOutsideActions(node: HTMLElement) {
    function handleClick(event: MouseEvent) {
      if (showActions && !node.contains(event.target as Node)) {
        closeActions();
      }
    }

    document.addEventListener('click', handleClick, true);

    return {
      destroy() {
        document.removeEventListener('click', handleClick, true);
      },
    };
  }
</script>

<div class="relative" use:clickOutsideActions>
  <button
    onclick={toggleActions}
    class="p-1.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-all duration-200 focus:opacity-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
    aria-label={ariaLabel}
    aria-expanded={showActions}
    data-testid={testId}
    {disabled}
  >
    <Icon class={iconSize === 'small' ? 'w-4 h-4' : iconSize === 'large' ? 'w-6 h-6' : 'w-5 h-5'} />
  </button>

  {#if showActions}
    <!-- Actions Dropdown -->
    <div
      class="absolute top-full end-0 rtl:start-0 rtl:end-auto mt-1 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 z-10 min-w-32"
    >
      <div class="py-1">
        {#each actions as action}
          <button
            onclick={() => handleAction(action)}
            class="flex items-center gap-2 w-full px-3 py-2 text-sm {action.className ||
              'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700'} transition-colors duration-150 {action?.disabled
              ? 'opacity-50 cursor-not-allowed'
              : ''}"
            data-testid={action.testId || `action-${action.label.toLowerCase().replace(/\s+/g, '-')}`}
            disabled={action?.disabled}
          >
            <action.icon class="w-4 h-4" />
            {action.label}
          </button>
        {/each}
      </div>
    </div>
  {/if}
</div>
