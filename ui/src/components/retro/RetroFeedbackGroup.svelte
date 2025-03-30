<script lang="ts">
  import { ThumbsUp } from 'lucide-svelte';
  import RetroFeedbackItem from './RetroFeedbackItem.svelte';

  interface Props {
    phase?: string;
    group?: any;
    handleVoteAction: (group: any) => void;
    voteLimitReached: boolean;
    isFacilitator?: boolean;
    users?: any;
    columnColors?: any;
    sendSocketEvent?: any;
  }

  let {
    phase = '',
    group = {
    name: 'Group',
    voteCount: 0,
    userVoted: false,
  },
    handleVoteAction,
    voteLimitReached,
    isFacilitator = false,
    users = [],
    columnColors = {},
    sendSocketEvent = (event: string, value: any) => {}
  }: Props = $props();
</script>

<div
  class="p-3 bg-white dark:bg-gray-800 rounded-lg shadow-lg flex flex-col flex-wrap text-gray-800 dark:text-white"
>
  <div class="flex items-center justify-between mb-4">
    <h2 class="text-xl font-bold">{group.name ? group.name : 'Group'}</h2>
    <div class="flex items-center space-x-2">
      {#if phase === 'vote'}
        <button
          onclick={() => {
            handleVoteAction(group);
          }}
          disabled="{voteLimitReached && !group.userVoted}"
          class="inline-block leading-none"
          class:text-gray-300="{voteLimitReached && !group.userVoted}"
          class:dark:text-gray-600="{voteLimitReached && !group.userVoted}"
          class:cursor-not-allowed="{voteLimitReached && !group.userVoted}"
          class:hover:text-blue-500="{!(voteLimitReached && !group.userVoted)}"
          class:dark:hover:text-sky-500="{!(
            voteLimitReached && !group.userVoted
          )}"
          class:text-green-500="{group.userVoted}"
          class:dark:text-lime-500="{group.userVoted}"
        >
          <ThumbsUp class="w-5 h-5 inline-block" />
        </button>
      {:else}
        <ThumbsUp
          class="w-5 h-5 inline-block text-green-500 dark:text-lime-500"
        />
      {/if}
      <span class="font-semibold text-green-500 dark:text-lime-500"
        >{group.voteCount}</span
      >
    </div>
  </div>
  <div class="flex-1 grow">
    {#each group.items as item, ii (item.id)}
      <RetroFeedbackItem
        item="{item}"
        phase="{phase}"
        users="{users}"
        isFacilitator="{isFacilitator}"
        sendSocketEvent="{sendSocketEvent}"
        columnColors="{columnColors}"
      />
    {/each}
  </div>
</div>
