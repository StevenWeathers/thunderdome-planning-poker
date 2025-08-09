<script lang="ts">
  import { run } from 'svelte/legacy';

  import { createEventDispatcher, onDestroy } from 'svelte';
  import { addMinutesToDate } from '../../dateUtils';

  interface Props {
    retroId?: string;
    timeLimitMin?: number;
    timeStart?: Date;
  }

  let { retroId = '', timeLimitMin = 0, timeStart = new Date() }: Props = $props();

  const dispatch = createEventDispatcher();

  let phaseEndTime = addMinutesToDate(timeStart, timeLimitMin);
  let now: Date = $state(new Date());

  let count = $derived(Math.round((phaseEndTime - now) / 1000));
  let h = $derived(Math.floor(count / 3600));
  let m = $derived(Math.floor((count - h * 3600) / 60));
  let s = $derived(count - h * 3600 - m * 60);

  function updateTimer() {
    now = new Date();
  }

  let interval = setInterval(updateTimer, 1000);
  run(() => {
    if (count === 0) {
      clearInterval(interval);
      dispatch('ended');
    }
  });

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
  class="inline-block me-2 font-semibold dark:text-gray-200 md:text-lg"
  data-testid="phase-timer"
>
  {#each Object.entries({ m, s }) as [key, value], i}
    <span class="me-2">{padValue(value)}{key}</span>
  {/each}
</div>
