<script>
  import { createEventDispatcher } from "svelte";
  import WarriorIcon from "./icons/WarriorIcon.svelte";

  const dispatch = createEventDispatcher();

  export let point = "1";
  export let active = false;
  export let isLocked = true;
  export let results = {
    count: 0,
    voters: []
  };

  let showVoters = false;

  $: activeColor = active
    ? "border-green bg-green-lightest text-green-dark"
    : "border-grey-light bg-white";
  $: lockedClass = isLocked
    ? "opacity-25 cursor-not-allowed"
    : "cursor-pointer";

  function voteAction() {
    if (isLocked) {
      return false;
    }
    if (!active) {
      dispatch("voted", {
        point
      });
    } else {
      dispatch("voteRetraction");
    }
  }
</script>

<div
  data-testId="pointCard"
  data-active={active}
  data-locked={isLocked}
  data-point={point}
  class="relative">
  {#if results.count}
    <div
      class="text-green font-semibold inline-block absolute pin-r pin-t p-2 text-4xl text-right {showVoters ? 'z-20' : 'z-10'}"
      data-testId="pointCardCount">
      {results.count}<button
        on:mouseenter={() => (showVoters = true)}
        on:mouseleave={() => (showVoters = false)}
        title="Show Voters"
        class="text-green relative">
        <WarriorIcon height="24" width="24" /><span
        class="text-right text-sm text-black font-normal w-48 absolute pin-l pin-t mt-0 ml-6 bg-white p-2 rounded shadow-md {showVoters ? '' : 'hidden'}">
        {#each results.voters as voter}
          {voter}
          <br />
        {/each}
      </span>
      </button>
    </div>
  {/if}
  <div
    class="w-full rounded overflow-hidden shadow-md border {activeColor}
    {lockedClass} relative text-3xl lg:text-5xl relative z-0"
    on:click={voteAction}>
    <div class="py-12 md:py-16 text-center">{point}</div>
  </div>
</div>
