<script lang="ts">
  import UserAvatar from '../user/UserAvatar.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user as sessionUser } from '../../stores';

  export let user = {};
  export let votes = [];
  export let maxVotes = 3;
  export let facilitators = [];
  export let phase = '';
  export let handleAddFacilitator = () => {};
  export let handleRemoveFacilitator = () => {};

  $: reachedMaxVotes =
    votes && votes.filter(v => v.userId === user.id).length === maxVotes;
</script>

<div
  class="shrink text-center px-2"
  data-testId="userCard"
  data-userName="{user.name}"
>
  <UserAvatar
    warriorId="{user.id}"
    avatar="{user.avatar}"
    gravatarHash="{user.gravatarHash}"
    class="mx-auto mb-2"
  />
  <div
    class="text-l font-bold leading-tight truncate dark:text-white"
    data-testId="userName"
    title="{user.name}"
  >
    {user.name}
    {#if facilitators.includes(user.id)}
      <div class="text-indigo-500 dark:text-violet-400">
        {$LL.facilitator()}
        {#if facilitators.includes($sessionUser.id)}
          <button
            class="text-red-500 text-sm"
            on:click="{handleRemoveFacilitator(user.id)}">{$LL.remove()}</button
          >
        {/if}
      </div>
    {:else if facilitators.includes($sessionUser.id)}
      <div>
        <button
          class="text-blue-500 dark:text-sky-400 text-sm"
          on:click="{handleAddFacilitator(user.id)}"
          >{$LL.makeFacilitator()}</button
        >
      </div>
    {/if}
    {#if phase === 'vote'}
      {#if reachedMaxVotes}
        <div class="text-green-600 dark:text-green-400">
          {$LL.allVotesIn()}
        </div>
      {:else}
        <div class="text-blue-600 dark:text-blue-400">
          {maxVotes - votes.filter(v => v.userId === user.id).length} votes left
        </div>
      {/if}
    {/if}
  </div>
</div>
