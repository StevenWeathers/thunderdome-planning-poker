<script lang="ts">
  import { onDestroy } from 'svelte';
  import { addMinutesToDate } from '../../dateUtils';

  interface Props {
    retroId?: string;
    timeLimitMin?: number;
    timeStart?: Date;
    onEnded?: () => void;
  }

  let props: Props = $props();

  const defaultTimeStart = new Date();

  let phaseEndTime: Date = $derived.by(() =>
    addMinutesToDate(props.timeStart ?? defaultTimeStart, props.timeLimitMin ?? 0),
  );
  let now: Date = $state(new Date());
  let didDispatchEnd = $state(false);

  let count = $derived(Math.max(0, Math.round((phaseEndTime.getTime() - now.getTime()) / 1000)));
  let h = $derived(Math.floor(count / 3600));
  let m = $derived(Math.floor((count - h * 3600) / 60));
  let s = $derived(count - h * 3600 - m * 60);

  function updateTimer() {
    now = new Date();
  }

  let interval = setInterval(updateTimer, 1000);
  $effect(() => {
    if (count > 0) {
      didDispatchEnd = false;
      return;
    }

    if (!didDispatchEnd) {
      didDispatchEnd = true;
      props.onEnded?.();
    }
  });

  function padValue(value: number, length = 2, char = '0') {
    const { length: currentLength } = value.toString();
    if (currentLength >= length) return value.toString();
    return `${char.repeat(length - currentLength)}${value}`;
  }

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<div class="inline-block me-2 font-semibold dark:text-gray-200 md:text-lg" data-testid="phase-timer">
  {#each Object.entries({ m, s }) as [key, value], i}
    <span class="me-2">{padValue(value)}{key}</span>
  {/each}
</div>
