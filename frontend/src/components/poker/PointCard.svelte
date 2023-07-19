<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import WarriorIcon from '../icons/UserIcon.svelte';
  import LL from '../../i18n/i18n-svelte';

  const dispatch = createEventDispatcher();

  export let point = '1';
  export let active = false;
  export let isLocked = true;
  export let results = {
    count: 0,
    voters: [],
  };
  export let hideVoterIdentity = false;

  let showVoters = false;

  $: activeColor = active
    ? 'border-green-500 bg-green-100 text-green-600 dark:border-lime-500 dark:bg-lime-100 dark:text-lime-700'
    : 'border-gray-300 bg-white dark:bg-gray-600 dark:border-gray-500 dark:text-gray-200';
  $: lockedClass = isLocked
    ? 'opacity-25 cursor-not-allowed'
    : 'cursor-pointer';

  function voteAction() {
    if (isLocked) {
      return false;
    }
    if (!active) {
      dispatch('voted', {
        point,
      });
    } else {
      dispatch('voteRetraction');
    }
  }
</script>

<div
  data-testid="pointCard"
  data-active="{active}"
  data-locked="{isLocked}"
  data-point="{point}"
  class="relative select-none"
>
  {#if results.count}
    <div
      class="text-green-500 dark:text-lime-400 font-semibold inline-block absolute end-0
            top-0 p-2 text-4xl text-right {showVoters ? 'z-20' : 'z-10'}"
      data-testid="pointCardCount"
    >
      {results.count}
      <button
        on:mouseenter="{() => {
          if (!hideVoterIdentity) {
            showVoters = true;
          }
        }}"
        on:mouseleave="{() => {
          if (!hideVoterIdentity) {
            showVoters = false;
          }
        }}"
        title="{$LL.showVoters()}"
        class="text-green-500 dark:text-lime-400 relative leading-none"
      >
        <WarriorIcon class="h-5 w-5" />
        <span
          class="text-right text-sm text-gray-900 font-normal w-48
                    absolute start-0 top-0 mt-0 ms-6 bg-white p-2 rounded
                    shadow-lg {showVoters ? '' : 'hidden'}"
        >
          {#each results.voters as voter}
            {voter}
            <br />
          {/each}
        </span>
      </button>
    </div>
  {/if}
  <div
    class="w-full rounded overflow-hidden shadow-lg border {activeColor}
        {lockedClass} relative text-5xl lg:text-6xl relative z-0 font-rajdhani"
    role="button"
    tabindex="0"
    on:click="{voteAction}"
    on:keypress="{voteAction}"
  >
    <div class="py-12 md:py-16 text-center">{point}</div>
  </div>
</div>
