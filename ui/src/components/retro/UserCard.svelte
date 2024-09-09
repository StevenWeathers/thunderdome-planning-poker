<script lang="ts">
  import UserAvatar from '../user/UserAvatar.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user as sessionUser } from '../../stores';
  import { ThumbsUp } from 'lucide-svelte';

  export let user = { id: '' };
  export let votes = [];
  export let maxVotes = 3;
  export let facilitators = [];
  export let readyUsers = [];
  export let phase = '';
  export let handleAddFacilitator = () => {};
  export let handleRemoveFacilitator = () => {};
  export let handleUserReady = () => {};
  export let handleUserUnReady = () => {};

  $: reachedMaxVotes =
    votes && votes.filter(v => v.userId === user.id).length === maxVotes;

  $: userReady = readyUsers && readyUsers.includes(user.id);
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
    userName="{user.name}"
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
    {#if phase === 'brainstorm'}
      {#if user.id === $sessionUser.id}
        <div>
          <button
            on:click="{!userReady
              ? handleUserReady($sessionUser.id)
              : handleUserUnReady($sessionUser.id)}"
            class="inline-block pointer text-gray-300 dark:text-gray-500 {userReady
              ? 'text-lime-600 dark:text-lime-400'
              : ''} hover:text-cyan-400 dark:hover:text-cyan-400"
          >
            <ThumbsUp
              class="inline-block w-8 h-8"
              title="Done brainstorming?"
            />
          </button>
        </div>
      {:else}
        <div
          class="inline-block text-gray-300 dark:text-gray-500 {userReady
            ? 'text-lime-600 dark:text-lime-400'
            : ''}"
        >
          <ThumbsUp class="inline-block w-8 h-8" title="Done brainstorming?" />
        </div>
      {/if}
    {/if}
    {#if phase === 'vote'}
      {#if reachedMaxVotes}
        <div class="text-lime-500 dark:text-lime-400">
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
