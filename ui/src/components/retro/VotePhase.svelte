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
    sendSocketEvent = (event: string, value: any) => {}
  }: Props = $props();

  // Calculate total votes used by current user across all groups
  const userVotesUsed = $derived(
    groups.reduce((acc, group) => {
      const userVotes = group.votes?.find(v => v.userId === $user.id)?.count || 0;
      return acc + userVotes;
    }, 0)
  );

  // Get user votes for a specific group
  const getUserVotesOnGroup = (group: RetroGroup) => {
    return group.votes?.find(v => v.userId === $user.id)?.count || 0;
  };
</script>

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