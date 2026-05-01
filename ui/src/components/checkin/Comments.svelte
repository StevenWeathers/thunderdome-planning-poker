<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import { user } from '../../stores';
  import type { UserDisplay } from '../../types/user';
  import LL from '../../i18n/i18n-svelte';
  import Comment from '../comments/Comment.svelte';
  import { MessageSquareMore } from '@lucide/svelte';
  import CommentsHeader from '../comments/CommentsHeader.svelte';
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
  <button
    onclick={toggleComments}
    class="group inline-flex items-center gap-2 rounded-full border border-gray-200 bg-white px-3 py-1.5 text-sm text-gray-600 transition-colors duration-200 hover:border-blue-200 hover:text-blue-600 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:border-blue-500/50 dark:hover:text-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
    aria-expanded={showComments}
    aria-haspopup="dialog"
    data-testid="checkin-comments-toggle"
  >
    <MessageSquareMore class="h-4 w-4" />
    <span class="font-medium leading-none text-gray-900 dark:text-white">{checkin.comments.length}</span>
    <span class="leading-none">{checkin.comments.length === 1 ? $LL.comment() : $LL.comments()}</span>
    <span class="sr-only">
      {showComments
        ? `Hide ${checkin.comments.length === 1 ? $LL.comment() : $LL.comments()}`
        : `Show ${checkin.comments.length === 1 ? $LL.comment() : $LL.comments()}`}
    </span>
  </button>

  {#if showComments}
    <Modal
      closeModal={toggleComments}
      widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
      ariaLabel={`${$LL.checkIn()} ${$LL.comments()}`}
    >
      <div
        class="mt-6 flex flex-col gap-2 rounded-xl border border-gray-200/50 bg-gray-50 p-4 pt-2 dark:border-gray-600/30 dark:bg-gray-700/30 dark:text-gray-300"
      >
        <CommentsHeader commentsCount={checkin.comments.length} title={$LL.comments()} />
        <div class="flex flex-col gap-3">
          {#each checkin.comments as comment}
            <Comment {comment} {userMap} {isAdmin} handleEdit={handleEditComment} handleDelete={handleDeleteComment} />
          {/each}
        </div>
        {#if checkin.comments.length === 0}
          <CommentEmptyState description="Be the first to share your thoughts on this check-in." />
        {/if}

        <CommentForm onSubmit={handleSubmitComment} />
      </div>
    </Modal>
  {/if}
</div>
