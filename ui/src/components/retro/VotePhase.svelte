<script lang="ts">
  import { user } from '../../stores';
  import ThumbsUp from '../icons/ThumbsUp.svelte';

  export let groups = [];
  export let handleVote = () => {};
  export let handleVoteSubtract = () => {};
  export let voteLimitReached = false;
  export let columns = [];

  const handleVoteAction = group => {
    const alreadyVoted = group.votes.includes($user.id);

    if (alreadyVoted) {
      handleVoteSubtract(group.id);
    } else {
      handleVote(group.id);
    }
  };

  $: columnColors =
    columns &&
    columns.reduce((p, c) => {
      p[c.name] = c.color;
      return p;
    }, {});
</script>

{#each groups as group, i (group.id)}
  {#if group.items.length > 0}
    <div
      class="border-2 p-2 dark:border-gray-800 rounded flex flex-col flex-wrap"
    >
      <div class="dark:text-gray-200 w-full text-center text-lg mb-2">
        <div class="flex content-center justify-center">
          <button
            on:click="{() => {
              handleVoteAction(group);
            }}"
            disabled="{voteLimitReached && !group.userVoted}"
            class="inline-block align-middle"
            class:text-gray-300="{voteLimitReached && !group.userVoted}"
            class:dark:text-gray-600="{voteLimitReached && !group.userVoted}"
            class:cursor-not-allowed="{voteLimitReached && !group.userVoted}"
            class:hover:text-blue-500="{!(
              voteLimitReached && !group.userVoted
            )}"
            class:dark:hover:text-sky-500="{!(
              voteLimitReached && !group.userVoted
            )}"
            class:text-green-500="{group.userVoted}"
            class:dark:text-lime-500="{group.userVoted}"
          >
            <ThumbsUp class="w-6 h-6 inline-block" />
          </button>
          <div class="inline-block align-middle text-2xl ms-2">
            {group.votes.length}
          </div>
        </div>
        {group.name ? group.name : 'Group'}
      </div>
      <div class="flex-1 grow">
        {#each group.items as item, ii (item.id)}
          <div
            class="p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-s-4 dark:text-white"
            class:border-green-400="{columnColors[item.type] === 'green'}"
            class:dark:border-lime-400="{columnColors[item.type] === 'green'}"
            class:border-red-500="{columnColors[item.type] === 'red'}"
            class:border-blue-400="{columnColors[item.type] === 'blue'}"
            class:dark:border-sky-400="{columnColors[item.type] === 'blue'}"
            class:border-yellow-500="{columnColors[item.type] === 'yellow'}"
            class:dark:border-yellow-400="{columnColors[item.type] ===
              'yellow'}"
            class:border-orange-500="{columnColors[item.type] === 'orange'}"
            class:dark:border-orange-400="{columnColors[item.type] ===
              'orange'}"
            class:border-teal-500="{columnColors[item.type] === 'teal'}"
            class:dark:border-teal-400="{columnColors[item.type] === 'teal'}"
            class:border-indigo-500="{columnColors[item.type] === 'purple'}"
            class:dark:border-indigo-400="{columnColors[item.type] ===
              'purple'}"
          >
            {item.content}
          </div>
        {/each}
      </div>
    </div>
  {/if}
{/each}
