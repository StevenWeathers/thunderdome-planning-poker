<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { PencilIcon, Trash2, Save, X, EllipsisIcon } from '@lucide/svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import type { UserDisplay } from '../../types/user';
  import UserAvatar from '../user/UserAvatar.svelte';

  interface Comment {
    comment: string;
    created_date: string;
    id: string;
    updated_date: string;
    user_id: string;
  }

  interface Props {
    comment?: Comment;
    userMap?: Map<string, UserDisplay>;
    isAdmin?: boolean;
    handleEdit?: (commentId: string, data: { userId: string; comment: string }) => void;
    handleDelete?: (commentId: string) => void;
  }

  let {
    comment = {} as Comment,
    userMap = new Map<string, UserDisplay>(),
    isAdmin = false,
    handleEdit = () => {},
    handleDelete = () => {},
  }: Props = $props();

  let showEdit = $state(false);
  let editcomment = $state('');
  let showActions = $state(false);

  $effect(() => {
    editcomment = comment.comment;
  });

  function toggleEdit() {
    showEdit = !showEdit;
    showActions = false;
    if (showEdit) {
      editcomment = comment.comment;
    }
  }

  function toggleActions() {
    showActions = !showActions;
  }

  function onSubmit(e: Event) {
    e.preventDefault();

    handleEdit(comment.id, {
      userId: $user.id,
      comment: editcomment,
    });

    toggleEdit();
  }

  function formatDate(dateString: string) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);

    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;
    return date.toLocaleDateString();
  }

  let canEdit = $derived(comment.user_id === $user.id || isAdmin);

  let showDeleteConfirm = $state(false);
</script>

<article
  class="animate-in fade-in-50 duration-200 group relative bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-600 p-4 hover:border-gray-300 dark:hover:border-gray-500 transition-all duration-200"
  data-commentid={comment.id}
  aria-label="Comment by {userMap.get(comment?.user_id)?.name || 'Unknown user'}"
>
  <!-- Comment Header -->
  <header class="flex items-start justify-between gap-3 mb-3">
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <!-- User Avatar -->
      <div class="flex-shrink-0">
        <UserAvatar
          warriorId={comment.user_id}
          pictureUrl={userMap.get(comment?.user_id)?.pictureUrl || ''}
          userName={userMap.get(comment?.user_id)?.name || 'Unknown'}
          gravatarHash={userMap.get(comment?.user_id)?.gravatarHash || ''}
          avatar={userMap.get(comment?.user_id)?.avatar || ''}
        />
      </div>

      <!-- User Info -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <h4 class="font-semibold text-gray-900 dark:text-white text-sm truncate">
            {userMap.get(comment.user_id)?.name || 'Loading...'}
          </h4>
          <span class="text-xs text-gray-500 dark:text-gray-400">•</span>
          <time class="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap">
            {formatDate(comment.created_date)}
          </time>
        </div>
      </div>
    </div>

    <!-- Actions Menu -->
    {#if canEdit && !showEdit}
      <div class="relative">
        <button
          onclick={toggleActions}
          class="opacity-0 group-hover:opacity-100 p-1.5 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-all duration-200 focus:opacity-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
          aria-label="Comment actions"
          aria-expanded={showActions}
        >
          <EllipsisIcon class="w-4 h-4" />
        </button>

        {#if showActions}
          <!-- Actions Dropdown -->
          <div
            class="absolute top-full end-0 rtl:start-0 rtl:end-auto mt-1 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 z-10 min-w-32"
          >
            <div class="py-1">
              <button
                onclick={toggleEdit}
                class="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-150"
              >
                <PencilIcon class="w-4 h-4" />
                {$LL.edit()}
              </button>
              <button
                onclick={() => {
                  showDeleteConfirm = true;
                  showActions = false;
                }}
                class="flex items-center gap-2 w-full px-3 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors duration-150"
              >
                <Trash2 class="w-4 h-4" />
                {$LL.delete()}
              </button>
            </div>
          </div>
        {/if}
      </div>
    {/if}
  </header>

  <!-- Comment Content -->
  <div class="space-y-3">
    {#if showEdit}
      <!-- Edit Mode -->
      <form onsubmit={onSubmit} name="editComment" class="space-y-4">
        <div>
          <label for="edit-comment-{comment.id}" class="sr-only"> Edit comment </label>
          <textarea
            id="edit-comment-{comment.id}"
            bind:value={editcomment}
            rows="3"
            class="block w-full resize-none rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-3 py-2.5 text-gray-900 dark:text-white placeholder:text-gray-500 dark:placeholder:text-gray-400 focus:border-blue-500 dark:focus:border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:focus:ring-blue-400/20 transition-all duration-200"
            placeholder="Edit your comment..."
          ></textarea>
        </div>

        <!-- Edit Actions -->
        <div class="flex items-center justify-end gap-2">
          <button
            type="button"
            onclick={toggleEdit}
            class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
          >
            <X class="w-4 h-4" />
            {$LL.cancel()}
          </button>

          <button
            type="submit"
            disabled={!editcomment.trim() || editcomment === comment.comment}
            class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 text-white text-sm font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-blue-600 dark:disabled:hover:bg-blue-500"
          >
            <Save class="w-4 h-4" />
            {$LL.updateComment()}
          </button>
        </div>
      </form>
    {:else}
      <!-- View Mode -->
      <div class="prose prose-sm max-w-none text-gray-800 dark:text-gray-200 leading-relaxed">
        <p class="whitespace-pre-wrap break-words mb-0">
          {comment.comment}
        </p>
      </div>
    {/if}
  </div>

  <!-- Updated Indicator -->
  {#if comment.updated_date && comment.updated_date !== comment.created_date}
    <footer class="mt-3 pt-2 border-t border-gray-100 dark:border-gray-700">
      <p class="text-xs text-gray-500 dark:text-gray-400">
        Edited {formatDate(comment.updated_date)}
      </p>
    </footer>
  {/if}

  <!-- Delete Confirmation Modal -->
  {#if showDeleteConfirm}
    <DeleteConfirmation
      toggleDelete={() => (showDeleteConfirm = false)}
      handleDelete={() => {
        handleDelete(comment.id);
        showDeleteConfirm = false;
      }}
      confirmText={'Are you sure you want to delete this comment?'}
      confirmBtnText={$LL.delete()}
      permanent={true}
    />
  {/if}
</article>

<!-- Click outside to close dropdown -->
{#if showActions}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="fixed inset-0 z-0" onclick={() => (showActions = false)} tabindex="-1"></div>
{/if}

<style>
  .prose p {
    margin-bottom: 0;
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

  .fade-in-50 {
    animation: fade-in 0.2s ease-out;
  }
</style>
