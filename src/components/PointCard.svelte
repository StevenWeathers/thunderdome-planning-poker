<script>
  import { createEventDispatcher } from "svelte"
  import WarriorIcon from "./WarriorIcon.svelte"

  const dispatch = createEventDispatcher()

  export let point = "1"
  export let active = false
  export let isLocked = true
  export let results = {
      count: 0,
      voters: [],
  }

  $: activeColor = active
    ? "border-green bg-green-lightest text-green-dark"
    : "border-grey-light bg-white"
  $: lockedClass = isLocked
    ? "opacity-25 cursor-not-allowed"
    : "cursor-pointer"

  function voteAction() {
    if (isLocked) {
      return false
    }
    if (!active) {
      dispatch("voted", {
        point
      })
    } else {
      dispatch("voteRetraction")
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
    <span
      class="text-green font-semibold inline-block absolute pin-r pin-t p-2 text-4xl"
      data-testId="pointCardCount"
    >
      {results.count}<WarriorIcon height="24" width="24" />
      <div>
        {#each results.voters as voter}
            {voter}<br />
        {/each}
      </div>
    </span>
  {/if}
  <div
    class="w-full rounded overflow-hidden shadow-md border {activeColor}
    {lockedClass} relative text-3xl lg:text-5xl"
    on:click={voteAction}>
    <div class="py-12 md:py-16 text-center">{point}</div>
  </div>
</div>
