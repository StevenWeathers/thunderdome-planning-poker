<script lang="ts">
  import TextInput from '../forms/TextInput.svelte';

  import type { ColorLegend } from '../../types/storyboard';

  interface Props {
    colorLegend?: Array<ColorLegend>;
  }

  let { colorLegend = $bindable([]) }: Props = $props();

  const labelForColor = (color: string) => `${color.charAt(0).toUpperCase()}${color.slice(1)} legend`;
</script>

<div class="space-y-3">
  {#each colorLegend as color, i}
    <div class="group">
      <label class="sr-only" for="legend-{color.color}">
        {labelForColor(color.color)}
      </label>
      <TextInput
        bind:value={colorLegend[i].legend}
        id="legend-{color.color}"
        name="legend-{color.color}"
        placeholder="Describe what {color.color} means"
      >
        {#snippet startElement()}
          <div class="w-6 h-6 rounded colorcard-{color.color}"></div>
        {/snippet}
      </TextInput>
    </div>
  {/each}
</div>

<style lang="postcss">
  .colorcard-gray {
    @apply bg-gray-400;
  }
  .colorcard-red {
    @apply bg-red-400;
  }
  .colorcard-orange {
    @apply bg-orange-400;
  }
  .colorcard-yellow {
    @apply bg-yellow-400;
  }
  .colorcard-green {
    @apply bg-green-400;
  }
  .colorcard-teal {
    @apply bg-teal-400;
  }
  .colorcard-blue {
    @apply bg-blue-400;
  }
  .colorcard-indigo {
    @apply bg-indigo-400;
  }
  .colorcard-purple {
    @apply bg-purple-400;
  }
  .colorcard-pink {
    @apply bg-pink-400;
  }
</style>
