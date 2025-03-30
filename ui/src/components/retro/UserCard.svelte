<script lang="ts">
  import UserAvatar from '../user/UserAvatar.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user as sessionUser } from '../../stores';
  import { ThumbsUp } from 'lucide-svelte';

  interface Props {
    user?: any;
    votes?: any;
    maxVotes?: number;
    facilitators?: any;
    readyUsers?: any;
    phase?: string;
    handleAddFacilitator?: any;
    handleRemoveFacilitator?: any;
    handleUserReady?: any;
    handleUserUnReady?: any;
  }

  let {
    user = { id: '' },
    votes = [],
    maxVotes = 3,
    facilitators = [],
    readyUsers = [],
    phase = '',
    handleAddFacilitator = () => {},
    handleRemoveFacilitator = () => {},
    handleUserReady = () => {},
    handleUserUnReady = () => {}
  }: Props = $props();

  let voteTally =
    $derived(votes &&
    votes.reduce(
      (p, v) => {
        if (v.userId === user.id) {
          p.votesLeft -= v.count;
          p.userVoteCount += v.count;
        }

        return p;
      },
      { userVoteCount: 0, votesLeft: maxVotes },
    ));

  let userReady = $derived(readyUsers && readyUsers.includes(user.id));
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
            onclick={handleRemoveFacilitator(user.id)}>{$LL.remove()}</button
          >
        {/if}
      </div>
    {:else if facilitators.includes($sessionUser.id)}
      <div>
        <button
          class="text-blue-500 dark:text-sky-400 text-sm"
          onclick={handleAddFacilitator(user.id)}
          >{$LL.makeFacilitator()}</button
        >
      </div>
    {/if}
    {#if phase === 'brainstorm'}
      {#if user.id === $sessionUser.id}
        <div>
          <button
            onclick={!userReady
              ? handleUserReady($sessionUser.id)
              : handleUserUnReady($sessionUser.id)}
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
      {#if voteTally.userVoteCount >= maxVotes}
        <div class="text-lime-500 dark:text-lime-400">
          {$LL.allVotesIn()}
        </div>
      {:else}
        <div class="text-blue-600 dark:text-blue-400">
          {voteTally.votesLeft} votes left
        </div>
      {/if}
    {/if}
  </div>
</div>
