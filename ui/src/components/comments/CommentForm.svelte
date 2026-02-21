<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { Send } from '@lucide/svelte';

  interface Props {
    onSubmit: (comment: string) => void;
  }

  let { onSubmit = () => {} }: Props = $props();

  let comment = $state('');
  function handleSubmit(e: Event) {
    e.preventDefault();
    if (comment.trim()) {
      onSubmit(comment);
      comment = '';
    }
  }
</script>

<!-- Divider -->
<div class="relative">
  <div class="absolute inset-0 flex items-center" aria-hidden="true">
    <div class="w-full border-t border-gray-200 dark:border-gray-600"></div>
  </div>
  <div class="relative flex justify-center">
    <span class="px-3 text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700 rounded-full">
      Add your thoughts
    </span>
  </div>
</div>

<form onsubmit={handleSubmit} name="checkinComment" class="space-y-4">
  <div class="relative">
    <label for="comment-input" class="sr-only">
      {$LL.writeCommentPlaceholder()}
    </label>
    <textarea
      id="comment-input"
      bind:value={comment}
      placeholder={$LL.writeCommentPlaceholder()}
      rows="3"
      class="block w-full resize-none rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-4 py-3 text-gray-900 dark:text-white placeholder:text-gray-500 dark:placeholder:text-gray-400 focus:border-blue-500 dark:focus:border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:focus:ring-blue-400/20 transition-all duration-200"
    ></textarea>

    <!-- Character counter (optional) -->
    {#if comment.length > 0}
      <div class="absolute bottom-2 end-2 text-xs text-gray-400 dark:text-gray-500" data-testid="comment-counter">
        {comment.length}
      </div>
    {/if}
  </div>

  <!-- Form Actions -->
  <div class="flex items-center justify-end gap-2">
    <button
      type="button"
      onclick={() => (comment = '')}
      class="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 transition-colors duration-200 disabled:opacity-50"
      disabled={!comment.trim()}
    >
      {$LL.clear()}
    </button>

    <button
      type="submit"
      disabled={!comment.trim()}
      class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 text-white font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-blue-600 dark:disabled:hover:bg-blue-500"
    >
      <Send class="w-4 h-4" />
      {$LL.postComment()}
    </button>
  </div>
</form>
