<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { MessageSquareMore, ChevronRight, ChevronDown } from '@lucide/svelte';

  interface Props {
    title?: string;
    commentsCount: number;
    onToggleExpand?: () => void;
    isExpanded?: boolean;
  }

  let { title, commentsCount = 0, onToggleExpand, isExpanded = false }: Props = $props();
</script>

{#if onToggleExpand}
  <button
    type="button"
    onclick={onToggleExpand}
    class="group flex items-center gap-2 w-full text-start rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
    aria-expanded={isExpanded}
  >
    <span
      class="flex items-center justify-center w-8 h-8 rounded-lg text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300 group-hover:bg-gray-100 dark:group-hover:bg-gray-700 transition-colors duration-200"
      aria-hidden="true"
    >
      {#if isExpanded}
        <ChevronDown class="w-4 h-4" />
      {:else}
        <ChevronRight class="w-4 h-4" />
      {/if}
    </span>
    <div
      class="flex items-center justify-center w-10 h-10 rounded-lg bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 group-hover:bg-blue-100 dark:group-hover:bg-blue-900/40 transition-colors duration-200"
    >
      <MessageSquareMore class="w-5 h-5" />
    </div>

    <div class="flex-1 text-start">
      <div class="flex items-center gap-2">
        <span class="font-medium text-gray-900 dark:text-white">
          <span>{commentsCount}</span>&nbsp;
          <span>
            {#if title}
              {title}
            {:else}
              {commentsCount === 1 ? $LL.comment() : $LL.comments()}
            {/if}
          </span>
        </span>
      </div>
    </div>
  </button>
{:else}
  <div class="group flex items-center gap-3 w-full py-2">
    <div
      class="flex items-center justify-center w-10 h-10 rounded-lg bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 group-hover:bg-blue-100 dark:group-hover:bg-blue-900/40 transition-colors duration-200"
    >
      <MessageSquareMore class="w-5 h-5" />
    </div>

    <div class="flex-1 text-start">
      <div class="flex items-center gap-2">
        <span class="font-medium text-gray-900 dark:text-white">
          <span>{commentsCount}</span>&nbsp;
          <span>
            {#if title}
              {title}
            {:else}
              {commentsCount === 1 ? $LL.comment() : $LL.comments()}
            {/if}
          </span>
        </span>
      </div>
    </div>
  </div>
{/if}
