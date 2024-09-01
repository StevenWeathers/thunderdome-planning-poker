<script lang="ts">
  import { createEventDispatcher, onDestroy } from 'svelte';
  import { addMinutesToDate } from '../../dateUtils';

  export let retroId: string = '';
  export let timeLimitMin: number = 0;
  export let timeStart: Date = new Date();

  const dispatch = createEventDispatcher();

  let phaseEndTime = addMinutesToDate(timeStart, timeLimitMin);
  let now: Date = new Date();

  $: count = Math.round((phaseEndTime - now) / 1000);
  $: h = Math.floor(count / 3600);
  $: m = Math.floor((count - h * 3600) / 60);
  $: s = count - h * 3600 - m * 60;

  function updateTimer() {
    now = new Date();
  }

  let interval = setInterval(updateTimer, 1000);
  $: if (count === 0) {
    clearInterval(interval);
    dispatch('ended');
  }

  function padValue(value, length = 2, char = '0') {
    const { length: currentLength } = value.toString();
    if (currentLength >= length) return value.toString();
    return `${char.repeat(length - currentLength)}${value}`;
  }

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<div
  class="inline-block mr-2 font-semibold dark:text-gray-200 md:text-lg"
  data-testid="phase-timer"
>
  {#each Object.entries({ m, s }) as [key, value], i}
    <span class="mr-2">{padValue(value)}{key}</span>
  {/each}
</div>
