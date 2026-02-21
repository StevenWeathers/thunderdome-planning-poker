<script lang="ts">
  import { user } from '../../stores';
  import type { UserDisplay } from '../../types/user';
  import LL from '../../i18n/i18n-svelte';
  import Comment from '../comments/Comment.svelte';
  import { MessageSquareMore, ChevronDown, ChevronUp } from '@lucide/svelte';
  import CommentEmptyState from '../comments/CommentEmptyState.svelte';
  import CommentForm from '../comments/CommentForm.svelte';

  interface Props {
    checkin?: any;
    userMap?: Map<string, UserDisplay>;
    isAdmin?: boolean;
    handleCreate?: any;
    handleEdit?: any;
    handleDelete?: any;
  }

  let {
    checkin = {},
    userMap = new Map<string, UserDisplay>(),
    isAdmin = false,
    handleCreate = () => {},
    handleEdit = () => {},
    handleDelete = () => {},
  }: Props = $props();

  let showComments = $state(false);

  function toggleComments() {
    showComments = !showComments;
  }

  function handleSubmitComment(commentText: string) {
    handleCreate(checkin.id, {
      userId: $user.id,
      comment: commentText,
    });
  }

  function handleEditComment(commentId: string, data: { userId: string; comment: string }) {
    handleEdit(checkin.id, commentId, data);
  }

  function handleDeleteComment(commentId: string) {
    handleDelete(checkin.id, commentId);
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
          <span>{checkin.comments.length}</span>&nbsp;
          <span>{checkin.comments.length === 1 ? $LL.comment() : $LL.comments()}</span>
        </span>
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
        <div class="flex flex-col gap-3">
          {#each checkin.comments as comment}
            <Comment {comment} {userMap} {isAdmin} handleEdit={handleEditComment} handleDelete={handleDeleteComment} />
          {/each}
        </div>
      {:else}
        <CommentEmptyState description="Be the first to share your thoughts on this check-in." />
      {/if}

      <!-- Comment Form -->
      <CommentForm onSubmit={handleSubmitComment} />
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

  .animate-in {
    animation-fill-mode: both;
  }

  .slide-in-from-top-2 {
    animation: slide-in-from-top 0.3s ease-out;
  }
</style>
