<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import CommentsHeader from '../comments/CommentsHeader.svelte';
  import Comment from '../comments/Comment.svelte';
  import type { RetroAction } from '../../types/retro';
  import type { UserDisplay } from '../../types/user';
  import type { TeamUser } from '../../types/team';
  import CommentEmptyState from '../comments/CommentEmptyState.svelte';
  import CommentForm from '../comments/CommentForm.svelte';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    toggleComments?: () => void;
    getRetrosActions?: () => void;
    actions?: RetroAction[];
    users?: TeamUser[];
    selectedActionId: string;
    isAdmin?: boolean;
  }

  let {
    xfetch,
    notifications,
    toggleComments = () => {},
    getRetrosActions = () => {},
    actions = [] as RetroAction[],
    users = [] as TeamUser[],
    selectedActionId = '',
    isAdmin = false,
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

  let userComment = $state('');

  function handleSubmitComment(commentText: string) {
    const body = {
      comment: commentText,
    };

    xfetch(`/api/retros/${selectedAction?.retroId}/actions/${selectedAction?.id}/comments`, {
      body,
    })
      .then(res => res.json())
      .then(function ({ data }) {
        userComment = '';
        getRetrosActions();
      })
      .catch(function () {
        notifications.danger($LL.retroActionCommentAddError());
      });
  }

  const handleCommentDelete = (commentId: string) => {
    xfetch(`/api/retros/${selectedAction?.retroId}/actions/${selectedAction?.id}/comments/${commentId}`, {
      method: 'DELETE',
    })
      .then(res => res.json())
      .then(function ({ data }) {
        getRetrosActions();
      })
      .catch(function () {
        notifications.danger($LL.retroActionCommentDeleteError());
      });
  };

  const handleCommentEdit = (commentId: string, data: { userId: string; comment: string }) => {
    const body = {
      comment: data.comment,
    };

    xfetch(`/api/retros/${selectedAction?.retroId}/actions/${selectedAction?.id}/comments/${commentId}`, {
      body,
      method: 'PUT',
    })
      .then(res => res.json())
      .then(function ({ data }) {
        getRetrosActions();
      })
      .catch(function () {
        notifications.danger($LL.retroActionCommentAddError());
      });
  };

  let selectedAction = $derived(actions && actions.find((a: RetroAction) => a.id === selectedActionId));
</script>

<Modal closeModal={toggleComments} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2" ariaLabel={$LL.modalRetroActionComments()}>
  <div
    class="mt-6 dark:text-gray-300 flex flex-col gap-2 bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4 pt-2 border border-gray-200/50 dark:border-gray-600/30"
  >
    <CommentsHeader commentsCount={selectedAction?.comments?.length || 0} title={$LL.actionComments()} />
    <div class="flex flex-col gap-3">
      {#each selectedAction?.comments as comment}
        <Comment {comment} {userMap} {isAdmin} handleEdit={handleCommentEdit} handleDelete={handleCommentDelete} />
      {/each}
    </div>
    {#if selectedAction?.comments?.length === 0}
      <CommentEmptyState description="Be the first to share your thoughts on this retro action item." />
    {/if}

    <CommentForm onSubmit={handleSubmitComment} />
  </div>
</Modal>
