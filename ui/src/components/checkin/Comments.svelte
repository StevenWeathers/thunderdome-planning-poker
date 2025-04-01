<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import Comment from './Comment.svelte';
  import { MessageSquareMore } from 'lucide-svelte';

  interface Props {
    checkin?: any;
    userMap?: any;
    isAdmin?: boolean;
    handleCreate?: any;
    handleEdit?: any;
    handleDelete?: any;
  }

  let {
    checkin = {},
    userMap = {},
    isAdmin = false,
    handleCreate = () => {},
    handleEdit = () => {},
    handleDelete = () => {}
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

<button class="text-blue-500 dark:text-sky-400" onclick={toggleComments}>
  <MessageSquareMore class="inline-block" />&nbsp;{checkin.comments.length}
  {checkin.comments.length === 1 ? 'Comment' : $LL.comments()}
</button>
{#if showComments}
  <div class="mt-2">
    {#each checkin.comments as comment}
      <Comment
        checkinId={checkin.id}
        comment={comment}
        userMap={userMap}
        isAdmin={isAdmin}
        handleEdit={handleEdit}
        handleDelete={handleDelete}
      />
    {/each}
  </div>
  <div class="text-right mb-2">
    <form onsubmit={onSubmit} name="checkinComment">
      <div class="mb-2 w-full">
        <textarea
          class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
          placeholder={$LL.writeCommentPlaceholder()}
          bind:value="{comment}"></textarea>
      </div>

      <div>
        <div class="text-right">
          <SolidButton type="submit">
            {$LL.postComment()}
          </SolidButton>
        </div>
      </div>
    </form>
  </div>
{/if}
