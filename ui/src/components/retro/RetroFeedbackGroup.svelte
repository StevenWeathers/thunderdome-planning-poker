<script lang="ts">
  import { ThumbsUp, ThumbsDown, Info, AlertCircle } from '@lucide/svelte';
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

  let showTooltip = $state(false);
  let tooltipTimeout: NodeJS.Timeout;
  let voteLimitReached = $derived(() => userVotesUsed === voteLimit);

  const handleVoteUpClick = () => {
    if (allowCumulativeVoting) {
      // Cumulative voting: add vote if not at limit
      if (voteLimitReached()) {
        showTooltip = true;
        clearTimeout(tooltipTimeout);
        tooltipTimeout = setTimeout(() => {
          showTooltip = false;
        }, 3000);
        return;
      }
      handleVote(group.id);
    } else {
      // Non-cumulative voting: add vote if user hasn't voted on this group and not at limit
      if (voteLimitReached() && userVotesOnThisGroup === 0) {
        showTooltip = true;
        clearTimeout(tooltipTimeout);
        tooltipTimeout = setTimeout(() => {
          showTooltip = false;
        }, 3000);
        return;
      }

      // Only allow voting if user hasn't voted on this group yet
      if (userVotesOnThisGroup === 0) {
        handleVote(group.id);
      }
    }
  };

  const handleVoteDownClick = () => {
    // Remove vote if user has votes on this group
    if (userVotesOnThisGroup > 0) {
      handleVoteSubtract(group.id);
    }
  };

  const canVoteUp = () => {
    if (allowCumulativeVoting) {
      return !voteLimitReached();
    } else {
      return userVotesOnThisGroup === 0 && !voteLimitReached();
    }
  };

  const canVoteDown = () => {
    return userVotesOnThisGroup > 0;
  };

  const getVoteUpButtonClasses = () => {
    const baseClasses =
      'relative inline-flex items-center justify-center min-w-[44px] min-h-[44px] rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800';

    if (!canVoteUp()) {
      return `${baseClasses} bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500 cursor-not-allowed border-2 border-gray-200 dark:border-gray-600`;
    }

    // Show active state if user has votes on this group
    if (userVotesOnThisGroup > 0) {
      return `${baseClasses} bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 hover:bg-blue-200 dark:hover:bg-blue-900/50 border-2 border-blue-300 dark:border-blue-700`;
    }

    return `${baseClasses} bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 hover:bg-blue-100 dark:hover:bg-blue-900/40 active:scale-95 border-2 border-blue-200 dark:border-blue-700 hover:border-blue-300 dark:hover:border-blue-600`;
  };

  const getVoteDownButtonClasses = () => {
    const baseClasses =
      'relative inline-flex items-center justify-center min-w-[44px] min-h-[44px] rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800';

    if (!canVoteDown()) {
      return `${baseClasses} bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500 cursor-not-allowed border-2 border-gray-200 dark:border-gray-600`;
    }

    return `${baseClasses} bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900/40 active:scale-95 border-2 border-red-200 dark:border-red-700 hover:border-red-300 dark:hover:border-red-600`;
  };

  const getVoteUpAriaLabel = () => {
    if (allowCumulativeVoting) {
      if (!canVoteUp()) {
        return `Cannot add vote. Vote limit reached (${userVotesUsed}/${voteLimit} votes used)`;
      }
      return `Add a vote to ${group.name}. Current votes: ${group.voteCount}`;
    } else {
      if (userVotesOnThisGroup > 0) {
        return `Cannot add vote. You already voted on ${group.name}`;
      }
      if (!canVoteUp()) {
        return `Cannot vote. Vote limit reached (${userVotesUsed}/${voteLimit} votes used)`;
      }
      return `Vote for ${group.name}. Current votes: ${group.voteCount}`;
    }
  };

  const getVoteDownAriaLabel = () => {
    if (!canVoteDown()) {
      return `Cannot remove vote. You have no votes on ${group.name}`;
    }
    return `Remove a vote from ${group.name}. You have ${userVotesOnThisGroup} vote${userVotesOnThisGroup === 1 ? '' : 's'} on this group`;
  };

  const getTooltipMessage = () => {
    if (allowCumulativeVoting) {
      return `Vote limit reached (${userVotesUsed}/${voteLimit})`;
    } else {
      if (userVotesOnThisGroup > 0) {
        return `You already voted on this group`;
      }
      return `Vote limit reached (${userVotesUsed}/${voteLimit})`;
    }
  };
</script>

<div
  class="group p-4 md:p-6 bg-white dark:bg-gray-800 rounded-xl shadow-lg hover:shadow-xl transition-shadow duration-300 flex flex-col text-gray-800 dark:text-white border border-gray-200 dark:border-gray-700"
  role="article"
  aria-labelledby="group-title-{group.id || 'default'}"
>
  <!-- Header -->
  <header class="flex items-start justify-between gap-4 mb-6">
    <h2
      id="group-title-{group.id || 'default'}"
      class="text-lg md:text-xl font-bold text-gray-900 dark:text-white flex-1 min-w-0"
      dir="auto"
    >
      {group.name || 'Group'}
    </h2>

    <!-- Vote Section -->
    <div class="flex items-center gap-3 flex-shrink-0" dir="ltr">
      {#if phase === 'vote'}
        <!-- Unified Dual Button Interface -->
        <div class="flex items-center gap-2">
          <!-- Vote Down Button -->
          <button
            onclick={handleVoteDownClick}
            disabled={!canVoteDown()}
            class={getVoteDownButtonClasses()}
            aria-label={getVoteDownAriaLabel()}
          >
            <ThumbsDown class="w-5 h-5" aria-hidden="true" />
          </button>

          <!-- Vote Up Button -->
          <div class="relative">
            <button
              onclick={handleVoteUpClick}
              disabled={!canVoteUp()}
              class={getVoteUpButtonClasses()}
              aria-label={getVoteUpAriaLabel()}
              tabindex="0"
            >
              <ThumbsUp class="w-5 h-5" aria-hidden="true" />

              <!-- Vote limit indicator -->
              {#if !canVoteUp() && (voteLimitReached() || (!allowCumulativeVoting && userVotesOnThisGroup > 0))}
                <AlertCircle class="w-3 h-3 absolute -top-1 -end-1 text-red-500 dark:text-red-400" aria-hidden="true" />
              {/if}
            </button>

            <!-- Tooltip for disabled state -->
            {#if showTooltip && !canVoteUp()}
              <div
                class="absolute bottom-full mb-2 start-1/2 transform -translate-x-1/2 px-3 py-2 bg-gray-900 dark:bg-gray-700 text-white text-xs rounded-lg shadow-lg whitespace-nowrap z-10 animate-in fade-in slide-in-from-bottom-2 duration-200"
                role="tooltip"
                aria-live="polite"
              >
                {getTooltipMessage()}
                <div
                  class="absolute top-full start-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900 dark:border-t-gray-700"
                ></div>
              </div>
            {/if}
          </div>
        </div>

        <!-- Vote count with better visual hierarchy -->
        {#if !hideVotesDuringVoting}
          <div class="flex items-center gap-1">
            <span
              class="font-bold text-lg md:text-xl text-green-600 dark:text-green-400 tabular-nums"
              aria-label="{group.voteCount} votes"
            >
              {group.voteCount}
            </span>
            <span class="text-sm text-gray-600 dark:text-gray-400 hidden sm:inline">
              {group.voteCount === 1 ? 'vote' : 'votes'}
            </span>
          </div>
        {/if}
      {:else}
        <!-- Non-voting phase display -->
        <div class="flex items-center gap-3">
          <div class="p-2 bg-green-100 dark:bg-green-900/30 rounded-lg">
            <ThumbsUp class="w-5 h-5 text-green-600 dark:text-green-400" aria-hidden="true" />
          </div>
          <div class="flex items-center gap-1">
            <span
              class="font-bold text-lg md:text-xl text-green-600 dark:text-green-400 tabular-nums"
              aria-label="{group.voteCount} votes"
            >
              {group.voteCount}
            </span>
            <span class="text-sm text-gray-600 dark:text-gray-400 hidden sm:inline">
              {group.voteCount === 1 ? 'vote' : 'votes'}
            </span>
          </div>
        </div>
      {/if}
    </div>
  </header>

  <!-- Vote status info (mobile-friendly) -->
  {#if phase === 'vote' && (voteLimitReached() || userVotesUsed > 0)}
    <div class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
      <div class="flex items-start gap-2">
        <Info class="w-4 h-4 text-blue-600 dark:text-blue-400 mt-0.5 flex-shrink-0" aria-hidden="true" />
        <div class="text-sm text-blue-700 dark:text-blue-300">
          <p>
            You have used <strong>{userVotesUsed}</strong> of <strong>{voteLimit}</strong> votes
            {#if userVotesOnThisGroup > 0}
              <br /><span class="text-blue-600 dark:text-blue-400 font-medium">
                {userVotesOnThisGroup} vote{userVotesOnThisGroup === 1 ? '' : 's'} on this group
              </span>
            {/if}
            {#if voteLimitReached()}
              <br /><span class="text-amber-600 dark:text-amber-400 font-medium">
                Vote limit reached
                {#if userVotesOnThisGroup > 0}
                  (you can still remove votes from this group)
                {/if}
              </span>
            {:else if allowCumulativeVoting}
              <br /><span class="text-green-600 dark:text-green-400 font-medium">
                You can vote multiple times on this group
              </span>
            {:else if userVotesOnThisGroup === 0}
              <br /><span class="text-green-600 dark:text-green-400 font-medium">
                You can vote once on this group
              </span>
            {:else}
              <br /><span class="text-blue-600 dark:text-blue-400 font-medium">
                You have already voted on this group
              </span>
            {/if}
          </p>
        </div>
      </div>
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

  /* Smooth animations */
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes slideInFromBottom {
    from {
      transform: translateY(8px) translateX(-50%);
    }
    to {
      transform: translateY(0) translateX(-50%);
    }
  }

  .animate-in {
    animation: fadeIn 0.2s ease-out;
  }

  .fade-in {
    animation: fadeIn 0.2s ease-out;
  }

  .slide-in-from-bottom-2 {
    animation: slideInFromBottom 0.2s ease-out;
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
