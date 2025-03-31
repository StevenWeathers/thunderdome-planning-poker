<script lang="ts">
  import { user } from '../../stores';
  import RetroFeedbackGroup from './RetroFeedbackGroup.svelte';

  interface Props {
    phase?: string;
    groups?: any;
    handleVote?: any;
    handleVoteSubtract?: any;
    voteLimitReached?: boolean;
    columns?: any;
    allowCumulativeVoting?: boolean;
    isFacilitator?: boolean;
    users?: any;
    columnColors?: any;
    sendSocketEvent?: any;
  }

  let {
    phase = 'vote',
    groups = [],
    handleVote = () => {},
    handleVoteSubtract = () => {},
    voteLimitReached = false,
    columns = [],
    allowCumulativeVoting = false,
    isFacilitator = false,
    users = [],
    columnColors = {},
    sendSocketEvent = (event: string, value: any) => {}
  }: Props = $props();

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
      group={group}
      handleVoteAction={handleVoteAction}
      voteLimitReached={voteLimitReached}
      users={users}
      isFacilitator={isFacilitator}
      sendSocketEvent={sendSocketEvent}
      columnColors={columnColors}
    />
  {/if}
{/each}
