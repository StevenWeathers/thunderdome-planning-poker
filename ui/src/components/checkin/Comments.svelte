<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import type { TeamUser } from '../../types/team';
  import Comment from './Comment.svelte';
  import { MessageSquareMore, ChevronDown, ChevronUp, Send } from 'lucide-svelte';

  interface Props {
    checkin?: any;
    userMap?: Map<string, TeamUser>;
    isAdmin?: boolean;
    handleCreate?: any;
    handleEdit?: any;
    handleDelete?: any;
  }

  let {
    checkin = {},
    userMap = new Map<string, TeamUser>(),
    isAdmin = false,
    handleCreate = () => {},
    handleEdit = () => {},
    handleDelete = () => {},
  }: Props = $props();

  let showComments = $state(false);
  let comment = $state('');

  function toggleComments() {
    showComments = !showComments;
  }

  function onSubmit(e) {
    e.preventDefault();

    handleCreate(checkin.id, {
      userId: $user.id,
      comment,
    });
    comment = '';
  }
</script>

<div class="space-y-4">
  <!-- Comments Toggle Button -->
  <button
    onclick={toggleComments}
    class="group flex items-center gap-3 w-full p-3 rounded-xl border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
    aria-expanded={showComments}
    aria-controls="comments-section"
  >
    <div
      class="flex items-center justify-center w-10 h-10 rounded-lg bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 group-hover:bg-blue-100 dark:group-hover:bg-blue-900/40 transition-colors duration-200"
    >
      <MessageSquareMore class="w-5 h-5" />
    </div>

    <div class="flex-1 text-start">
      <div class="flex items-center gap-2">
        <span class="font-medium text-gray-900 dark:text-white">
          {checkin.comments.length}
          {checkin.comments.length === 1 ? 'Comment' : $LL.comments()}
        </span>
        {#if checkin.comments.length > 0}
          <span
            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900/50 text-blue-800 dark:text-blue-200"
          >
            {checkin.comments.length}
          </span>
        {/if}
      </div>
      <p class="text-sm text-gray-500 dark:text-gray-400">
        {showComments
          ? 'Hide conversation'
          : checkin.comments.length > 0
            ? 'View conversation'
            : 'Start a conversation'}
      </p>
    </div>

    <div
      class="flex items-center justify-center w-8 h-8 rounded-lg text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-all duration-200"
    >
      {#if showComments}
        <ChevronUp class="w-5 h-5 transform transition-transform duration-200" />
      {:else}
        <ChevronDown class="w-5 h-5 transform transition-transform duration-200" />
      {/if}
    </div>
  </button>

  <!-- Comments Section with Animation -->
  {#if showComments}
    <div
      id="comments-section"
      class="space-y-4 animate-in slide-in-from-top-2 duration-300"
      role="region"
      aria-label="Comments section"
    >
      <!-- Comments List -->
      {#if checkin.comments.length > 0}
        <div class="space-y-3">
          {#each checkin.comments as comment}
            <div class="animate-in fade-in-50 duration-200">
              <Comment checkinId={checkin.id} {comment} {userMap} {isAdmin} {handleEdit} {handleDelete} />
            </div>
          {/each}
        </div>

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
      {:else}
        <!-- Empty State -->
        <div class="text-center py-8">
          <div
            class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center"
          >
            <MessageSquareMore class="w-8 h-8 text-gray-400 dark:text-gray-500" />
          </div>
          <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">No comments yet</h3>
          <p class="text-gray-500 dark:text-gray-400">Be the first to share your thoughts on this check-in.</p>
        </div>
      {/if}

      <!-- Comment Form -->
      <form onsubmit={onSubmit} name="checkinComment" class="space-y-4">
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
            <div class="absolute bottom-2 end-2 text-xs text-gray-400 dark:text-gray-500">
              {comment.length}/500
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
            Clear
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
    </div>
  {/if}
</div>

<style>
  @keyframes slide-in-from-top {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes fade-in {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  .animate-in {
    animation-fill-mode: both;
  }

  .slide-in-from-top-2 {
    animation: slide-in-from-top 0.3s ease-out;
  }

  .fade-in-50 {
    animation: fade-in 0.2s ease-out;
  }

  /* Custom scrollbar */
  .scrollbar-thin {
    scrollbar-width: thin;
  }

  .scrollbar-thumb-gray-300::-webkit-scrollbar-thumb {
    background-color: rgb(209 213 219);
    border-radius: 9999px;
  }

  .dark .scrollbar-thumb-gray-600::-webkit-scrollbar-thumb {
    background-color: rgb(75 85 99);
  }

  .scrollbar-track-transparent::-webkit-scrollbar-track {
    background-color: transparent;
  }

  ::-webkit-scrollbar {
    width: 6px;
  }
</style>
