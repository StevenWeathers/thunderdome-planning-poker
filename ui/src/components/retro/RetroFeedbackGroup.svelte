<script lang="ts">
  import { Check, Minus, Plus, ThumbsUp, X } from '@lucide/svelte';
  import RetroFeedbackItem from './RetroFeedbackItem.svelte';

  interface Props {
    phase?: string;
    group?: any;
    handleVote: (group: any) => void;
    handleVoteSubtract: (group: any) => void;
    isFacilitator?: boolean;
    users?: any;
    columnColors?: any;
    sendSocketEvent?: any;
    voteLimit?: number;
    userVotesUsed?: number;
    allowCumulativeVoting?: boolean;
    userVotesOnThisGroup?: number;
    hideVotesDuringVoting?: boolean;
  }

  let {
    phase = '',
    group = {
      name: 'Group',
      voteCount: 0,
    },
    handleVote,
    handleVoteSubtract,
    isFacilitator = false,
    users = [],
    columnColors = {},
    sendSocketEvent = (event: string, value: any) => {},
    voteLimit = 3,
    userVotesUsed = 0,
    allowCumulativeVoting = false,
    userVotesOnThisGroup = 0,
    hideVotesDuringVoting = false,
  }: Props = $props();

  // >= because the facilitator can lower the limit after votes are cast
  const voteLimitReached = $derived(userVotesUsed >= voteLimit);
  const hasVoted = $derived(userVotesOnThisGroup > 0);

  const addVoteLabel = $derived(
    voteLimitReached ? `Cannot add vote to ${group.name}, no votes remaining` : `Add a vote to ${group.name}`,
  );
  const removeVoteLabel = $derived(
    hasVoted ? `Remove a vote from ${group.name}` : `Cannot remove vote, you have no votes on ${group.name}`,
  );
  const toggleVoteLabel = $derived(
    hasVoted
      ? `Remove your vote from ${group.name}`
      : voteLimitReached
        ? `Cannot vote for ${group.name}, no votes remaining`
        : `Vote for ${group.name}`,
  );

  const handleToggleVoteClick = () => {
    if (hasVoted) {
      handleVoteSubtract(group.id);
    } else {
      handleVote(group.id);
    }
  };

  const stepButtonClasses =
    'inline-flex items-center justify-center min-w-[44px] min-h-[44px] rounded-lg border-2 transition-colors duration-200 border-blue-300 dark:border-blue-700 text-blue-700 dark:text-blue-300 bg-white dark:bg-gray-800 enabled:hover:bg-blue-100 dark:enabled:hover:bg-blue-900/40 disabled:cursor-not-allowed disabled:border-gray-200 dark:disabled:border-gray-600 disabled:text-gray-400 dark:disabled:text-gray-500';
</script>

<div
  class="group p-4 md:p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg hover:shadow-xl transition-shadow duration-300 flex flex-col text-gray-800 dark:text-white border {phase ===
    'vote' && hasVoted
    ? 'border-blue-500 dark:border-blue-400 ring-1 ring-blue-500 dark:ring-blue-400'
    : 'border-gray-200 dark:border-gray-700'}"
  role="article"
  aria-labelledby="group-title-{group.id || 'default'}"
>
  <!-- Header -->
  <header class="flex items-start justify-between gap-4 mb-4">
    <h2
      id="group-title-{group.id || 'default'}"
      class="text-lg md:text-xl font-bold text-gray-900 dark:text-white flex-1 min-w-0"
      dir="auto"
    >
      {group.name || 'Group'}
    </h2>

    <!-- Total votes from all users -->
    {#if phase !== 'vote' || !hideVotesDuringVoting}
      <div class="flex items-center gap-2 flex-shrink-0" dir="ltr">
        <div class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
          <ThumbsUp class="w-5 h-5 text-green-600 dark:text-green-400" aria-hidden="true" />
        </div>
        <div class="flex items-center gap-1">
          <span
            class="font-bold text-lg md:text-xl text-green-600 dark:text-green-400 tabular-nums"
            aria-label="{group.voteCount} total votes"
          >
            {group.voteCount}
          </span>
          <span class="text-sm text-gray-600 dark:text-gray-400 hidden sm:inline">
            {group.voteCount === 1 ? 'vote' : 'votes'}
          </span>
        </div>
      </div>
    {/if}
  </header>

  <!-- Current user's vote controls -->
  {#if phase === 'vote'}
    <div class="mb-4" dir="ltr">
      {#if allowCumulativeVoting}
        <div
          class="flex items-center justify-between gap-2 p-1 rounded-xl bg-blue-50 dark:bg-blue-900/20"
          role="group"
          aria-label="Your votes on {group.name}"
        >
          <button
            type="button"
            onclick={() => handleVoteSubtract(group.id)}
            disabled={!hasVoted}
            class={stepButtonClasses}
            aria-label={removeVoteLabel}
            title="Remove a vote"
          >
            <Minus class="w-5 h-5" aria-hidden="true" />
          </button>

          <span class="flex flex-col items-center leading-tight text-blue-700 dark:text-blue-300" aria-live="polite">
            <strong class="text-lg tabular-nums" data-testid="user-group-votes">{userVotesOnThisGroup}</strong>
            <span class="text-xs">yours</span>
          </span>

          <button
            type="button"
            onclick={() => handleVote(group.id)}
            disabled={voteLimitReached}
            class={stepButtonClasses}
            aria-label={addVoteLabel}
            title={voteLimitReached ? 'No votes remaining' : 'Add a vote'}
          >
            <Plus class="w-5 h-5" aria-hidden="true" />
          </button>
        </div>
      {:else}
        <button
          type="button"
          onclick={handleToggleVoteClick}
          disabled={!hasVoted && voteLimitReached}
          class="w-full min-h-[44px] px-3 inline-flex flex-wrap items-center justify-center gap-x-3 gap-y-1 rounded-lg border-2 font-semibold transition-colors duration-200 {hasVoted
            ? 'bg-blue-600 border-blue-600 text-white hover:bg-blue-700 hover:border-blue-700'
            : 'bg-white dark:bg-gray-800 border-blue-300 dark:border-blue-700 text-blue-700 dark:text-blue-300 enabled:hover:bg-blue-50 dark:enabled:hover:bg-blue-900/40 disabled:cursor-not-allowed disabled:border-gray-200 dark:disabled:border-gray-600 disabled:text-gray-400 dark:disabled:text-gray-500'}"
          aria-label={toggleVoteLabel}
          title={!hasVoted && voteLimitReached ? 'No votes remaining' : undefined}
        >
          {#if hasVoted}
            <span class="inline-flex items-center gap-1">
              <Check class="w-5 h-5" aria-hidden="true" />
              Voted
            </span>
            <span class="inline-flex items-center gap-1 text-sm font-normal underline">
              <X class="w-4 h-4" aria-hidden="true" />
              Remove
            </span>
          {:else}
            <span class="inline-flex items-center gap-1">
              <Plus class="w-5 h-5" aria-hidden="true" />
              Vote
            </span>
          {/if}
        </button>
      {/if}
    </div>
  {/if}

  <!-- Items Container -->
  <main class="flex-1 space-y-3" role="group" aria-label="Feedback items for {group.name}">
    {#each group.items as item, ii (item.id)}
      <RetroFeedbackItem {item} {phase} {users} {isFacilitator} {sendSocketEvent} {columnColors} />
    {/each}

    {#if group.items?.length === 0}
      <div class="text-center py-8 text-gray-500 dark:text-gray-400">
        <p class="text-sm">No feedback items yet</p>
      </div>
    {/if}
  </main>
</div>

<style>
  /* RTL support */
  :global([dir='rtl']) .group {
    text-align: right;
  }

  /* Enhanced focus styles for accessibility */
  button:focus-visible {
    outline: 2px solid currentColor;
    outline-offset: 2px;
  }

  /* High contrast mode support */
  @media (prefers-contrast: high) {
    .group {
      border-width: 2px;
    }

    button {
      border-width: 2px;
    }
  }

  /* Reduced motion support */
  @media (prefers-reduced-motion: reduce) {
    * {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
    }
  }

  /* Touch device optimizations */
  @media (hover: none) {
    button:hover {
      transform: none;
    }
  }
</style>
