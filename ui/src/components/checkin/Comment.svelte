<script lang="ts">
  import HollowButton from '../global/HollowButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import SolidButton from '../global/SolidButton.svelte';
  import { User } from 'lucide-svelte';

  interface Props {
    checkinId?: any;
    comment?: any;
    userMap?: any;
    isAdmin?: boolean;
    handleEdit?: any;
    handleDelete?: any;
  }

  let {
    checkinId = {},
    comment = {},
    userMap = {},
    isAdmin = false,
    handleEdit = () => {},
    handleDelete = () => {}
  }: Props = $props();

  let showEdit = $state(false);
  let editcomment = $state(`${comment.comment}`);

  function toggleEdit() {
    showEdit = !showEdit;
  }

  function onSubmit(e) {
    e.preventDefault();

    handleEdit(checkinId, comment.id, {
      userId: $user.id,
      comment: editcomment,
    });

    toggleEdit();
  }
</script>

<div
  class="w-full mb-2 text-gray-700 dark:text-gray-300 border-b border-gray-300 dark:border-gray-700"
  data-commentid="{comment.id}"
>
  <div class="font-bold">
    <User class="h-4 w-4 inline-block" />&nbsp;{userMap[comment.user_id] ||
      '...'}
  </div>
  {#if showEdit}
    <div class="w-full my-2">
      <form onsubmit={onSubmit} name="checkinComment">
        <textarea
          class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
    rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
    focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
          bind:value="{editcomment}"></textarea>
        <div class="text-right">
          <HollowButton color="blue" onClick={toggleEdit}>
            {$LL.cancel()}
          </HollowButton>
          <SolidButton type="submit" disabled={editcomment === ''}>
            {$LL.updateComment()}
          </SolidButton>
        </div>
      </form>
    </div>
  {:else}
    <div class="py-2">
      {comment.comment}
    </div>
  {/if}
  {#if (comment.user_id === $user.id || comment.user_id === isAdmin) && !showEdit}
    <div class="mb-2 text-right">
      <button
        class="text-blue-500 hover:text-blue-300 dark:text-sky-300 dark:hover:text-sky-100 me-1"
        onclick={toggleEdit}
      >
        {$LL.edit()}
      </button>
      <button
        class="text-red-500"
        onclick={handleDelete(checkinId, comment.id)}
      >
        {$LL.delete()}
      </button>
    </div>
  {/if}
</div>
