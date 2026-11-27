<script lang="ts">
  import { user } from '../../stores';
  import type { RetroGroup } from '../../types/retro';
  import RetroFeedbackGroup from './RetroFeedbackGroup.svelte';

  interface Props {
    phase?: string;
    groups?: Array<RetroGroup>;
    handleVote?: any;
    handleVoteSubtract?: any;
    voteLimit?: number;
    columns?: any;
    allowCumulativeVoting?: boolean;
    isFacilitator?: boolean;
    users?: any;
    columnColors?: any;
    sendSocketEvent?: any;
    hideVotesDuringVoting?: boolean;
  }

  let {
    phase = 'vote',
    groups = [] as RetroGroup[],
    handleVote = () => {},
    handleVoteSubtract = () => {},
    voteLimit = 3,
    columns = [],
    allowCumulativeVoting = false,
    isFacilitator = false,
    users = [],
    columnColors = {},
    hideVotesDuringVoting = false,
    sendSocketEvent = (event: string, value: any) => {},
  }: Props = $props();

  // Calculate total votes used by current user across all groups
  const userVotesUsed = $derived(
    groups.reduce((acc, group) => {
      const userVotes = group.votes?.find(v => v.userId === $user.id)?.count || 0;
      return acc + userVotes;
    }, 0),
  );

  const votesRemaining = $derived(Math.max(voteLimit - userVotesUsed, 0));

  // Get user votes for a specific group
  const getUserVotesOnGroup = (group: RetroGroup) => {
    return group.votes?.find(v => v.userId === $user.id)?.count || 0;
  };
</script>

<div
  class="col-span-full flex flex-wrap items-baseline gap-x-3 gap-y-1 px-4 py-3 rounded-lg border bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800 text-blue-700 dark:text-blue-300"
  role="status"
  data-testid="vote-summary"
>
  <p class="font-bold">{votesRemaining} of {voteLimit} votes remaining</p>
  <p class="text-sm">
    {#if votesRemaining === 0}
      To move a vote, remove it from one group first.
    {:else if allowCumulativeVoting}
      Use + to add a vote to a group and − to take one back.
    {:else}
      Vote once per group. Select a voted group's button again to remove your vote.
    {/if}
  </p>
</div>

{#each groups as group, _ (group.id)}
  {#if (group.items ?? []).length > 0}
    {@const userVotesOnThisGroup = getUserVotesOnGroup(group)}

    <RetroFeedbackGroup
      {phase}
      group={{
        ...group,
      }}
      {handleVote}
      {handleVoteSubtract}
      {allowCumulativeVoting}
      {voteLimit}
      {userVotesOnThisGroup}
      {userVotesUsed}
      {users}
      {isFacilitator}
      {sendSocketEvent}
      {columnColors}
      {hideVotesDuringVoting}
    />
  {/if}
{/each}
