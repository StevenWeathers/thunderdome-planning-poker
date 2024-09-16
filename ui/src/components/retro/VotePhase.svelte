<script lang="ts">
  import { user } from '../../stores';
  import RetroFeedbackGroup from './RetroFeedbackGroup.svelte';

  export let phase = 'vote';
  export let groups = [];
  export let handleVote = () => {};
  export let handleVoteSubtract = () => {};
  export let voteLimitReached = false;
  export let columns = [];
  export let allowCumulativeVoting: boolean = false;
  export let isFacilitator = false;
  export let users = [];
  export let columnColors: any = {};
  export let sendSocketEvent = (event: string, value: any) => {};

  const handleVoteAction = group => {
    const userVoted = group.votes.find(v => v.userId === $user.id);
    if (
      (userVoted && !allowCumulativeVoting) ||
      (allowCumulativeVoting && voteLimitReached)
    ) {
      handleVoteSubtract(group.id);
    } else {
      handleVote(group.id);
    }
  };
</script>

{#each groups as group, i (group.id)}
  {#if group.items.length > 0}
    <RetroFeedbackGroup
      phase="{phase}"
      group="{group}"
      handleVoteAction="{handleVoteAction}"
      voteLimitReached="{voteLimitReached}"
      users="{users}"
      isFacilitator="{isFacilitator}"
      sendSocketEvent="{sendSocketEvent}"
      columnColors="{columnColors}"
    />
  {/if}
{/each}
