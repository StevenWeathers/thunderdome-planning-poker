<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import type { RetroAction, RetroActionComment, RetroUser } from '../../types/retro';
  import type { UserDisplay } from '../../types/user';
  import CommentsHeader from '../comments/CommentsHeader.svelte';
  import Comment from '../comments/Comment.svelte';
  import CommentEmptyState from '../comments/CommentEmptyState.svelte';
  import CommentForm from '../comments/CommentForm.svelte';

  interface Props {
    toggleComments?: any;
    item?: RetroAction;
    users?: RetroUser[];
    isFacilitator?: boolean;
    sendSocketEvent?: any;
  }

  let {
    toggleComments = () => {},
    item = {
      id: '',
      comments: [] as RetroActionComment[],
    } as RetroAction,
    users = [] as RetroUser[],
    isFacilitator = false,
    sendSocketEvent = (event: string, value: any) => {},
  }: Props = $props();

  const userMap: Map<string, UserDisplay> = $derived(
    users.reduce((prev, cur) => {
      prev.set(cur.id, {
        id: cur.id,
        name: cur.name,
        avatar: cur.avatar,
        gravatarHash: cur.gravatarHash,
        pictureUrl: '',
      });
      return prev;
    }, new Map<string, UserDisplay>()),
  );

  function handleSubmitComment(commentText: string) {
    sendSocketEvent(
      'item_comment_add',
      JSON.stringify({
        item_id: item.id,
        comment: commentText,
      }),
    );
  }

  const handleCommentDelete = (commentId: string) => {
    sendSocketEvent(
      'item_comment_delete',
      JSON.stringify({
        comment_id: commentId,
      }),
    );
  };

  const handleCommentEdit = (commentId: string, data: { userId: string; comment: string }) => {
    sendSocketEvent(
      'item_comment_edit',
      JSON.stringify({
        comment_id: commentId,
        comment: data.comment,
      }),
    );
  };
</script>

<Modal closeModal={toggleComments} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2" ariaLabel={$LL.modalRetroItemComments()}>
  <div
    class="mt-6 dark:text-gray-300 flex flex-col gap-2 bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4 pt-2 border border-gray-200/50 dark:border-gray-600/30"
  >
    <CommentsHeader commentsCount={item.comments.length} />
    <div class="flex flex-col gap-3">
      {#each item.comments as comment}
        <Comment
          {comment}
          {userMap}
          isAdmin={isFacilitator}
          handleEdit={handleCommentEdit}
          handleDelete={handleCommentDelete}
        />
      {/each}
    </div>
    {#if item.comments.length === 0}
      <CommentEmptyState description="Be the first to share your thoughts on this retro feedback item." />
    {/if}

    <CommentForm onSubmit={handleSubmitComment} />
  </div>
</Modal>
