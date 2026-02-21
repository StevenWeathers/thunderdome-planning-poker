<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { MessageSquareMore } from '@lucide/svelte';

  interface Comment {
    comment: string;
    created_date: string;
    id: string;
    updated_date: string;
    user_id: string;
  }

  interface Props {
    title?: string;
    commentsCount: number;
  }

  let { title, commentsCount = 0 }: Props = $props();
</script>

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
