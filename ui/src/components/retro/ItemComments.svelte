<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import UserIcon from '../icons/UserIcon.svelte';
  import { user } from '../../stores';

  export let toggleComments = () => {};
  export let items = [];
  export let selectedItemId = '';
  export let users = [];
  export let isFacilitator = false;
  export let sendSocketEvent = (event: string, value: any) => {};

  const userMap = users.reduce((prev, cur) => {
    prev[cur.id] = cur.name;
    return prev;
  }, {});

  let userComment = '';
  let selectedComment = null;
  let selectedCommentContent = '';

  const toggleCommentEdit = comment => () => {
    selectedComment = comment;
    if (comment !== null) {
      selectedCommentContent = comment.comment;
    }
  };

  function handleCommentSubmit() {
    sendSocketEvent(
      'item_comment_add',
      JSON.stringify({
        item_id: item.id,
        comment: userComment,
      }),
    );
    userComment = '';
  }

  const handleCommentDelete = commentId => () => {
    sendSocketEvent(
      'item_comment_delete',
      JSON.stringify({
        comment_id: commentId,
      }),
    );
  };

  const handleCommentEdit = () => {
    sendSocketEvent(
      'item_comment_edit',
      JSON.stringify({
        comment_id: selectedComment.id,
        comment: selectedCommentContent,
      }),
    );
    toggleCommentEdit(null)();
  };

  $: item =
    items && selectedItemId
      ? items.find(i => i.id === selectedItemId)
      : { id: '', comments: [] };
</script>

<Modal closeModal="{toggleComments}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <div class="mt-12 dark:text-gray-300">
    <h3 class="text-xl pb-2 mb-4 border-b border-gray-600 dark:border-gray-400">
      Comments
    </h3>
    {#each item.comments as comment}
      <div
        class="w-full mb-4 text-gray-700 dark:text-gray-400 border-b border-gray-300 dark:border-gray-700"
        data-commentid="{comment.id}"
      >
        <div class="font-bold">
          <UserIcon class="h-4 w-4" />&nbsp;{userMap[comment.user_id] || '...'}
        </div>
        {#if selectedComment !== null && selectedComment.id === comment.id}
          <div class="w-full my-2">
            <textarea
              class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
              bind:value="{selectedCommentContent}"></textarea>
            <div class="text-right">
              <HollowButton color="blue" onClick="{toggleCommentEdit(null)}">
                {$LL.cancel()}
              </HollowButton>
              <HollowButton
                color="green"
                onClick="{handleCommentEdit}"
                disabled="{selectedCommentContent === ''}"
              >
                {$LL.updateComment()}
              </HollowButton>
            </div>
          </div>
        {:else}
          <div class="py-2">
            {comment.comment}
          </div>
        {/if}
        {#if (comment.user_id === $user.id || comment.user_id === isFacilitator) && !(selectedComment !== null && selectedComment.id === comment.id)}
          <div class="mb-2 text-right">
            <button
              class="text-blue-500 hover:text-blue-300 me-1"
              on:click="{toggleCommentEdit(comment)}"
            >
              {$LL.edit()}
            </button>
            <button
              class="text-red-500"
              on:click="{handleCommentDelete(comment.id)}"
            >
              {$LL.delete()}
            </button>
          </div>
        {/if}
      </div>
    {/each}
    {#if item.comments.length === 0}
      <p class="text-lg dark:text-gray-400">{$LL.noComments()}</p>
    {/if}

    <div class="w-full mt-8">
      <textarea
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
        placeholder="{$LL.writeCommentPlaceholder()}"
        bind:value="{userComment}"></textarea>
      <div class="text-right">
        <HollowButton
          color="teal"
          onClick="{handleCommentSubmit}"
          disabled="{userComment === ''}"
        >
          {$LL.postComment()}
        </HollowButton>
      </div>
    </div>
  </div>
</Modal>
