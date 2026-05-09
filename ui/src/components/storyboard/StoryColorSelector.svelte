<script lang="ts">
  import type { ColorLegend } from '../../types/storyboard';

  interface Props {
    value?: string;
    colorLegend?: ColorLegend[];
    allowNone?: boolean;
    noneLabel?: string;
    onSelect?: (color: string) => void;
  }

  let {
    value = $bindable(''),
    colorLegend = [],
    allowNone = false,
    noneLabel = 'No default',
    onSelect = () => {},
  }: Props = $props();

  function selectColor(color: string) {
    return () => {
      value = color;
      onSelect(color);
    };
  }

  function colorTitle(color: ColorLegend) {
    return color.legend ? `${color.color} - ${color.legend}` : color.color;
  }
</script>

<div class="flex flex-wrap gap-2 pt-1">
  {#if allowNone}
    <button
      type="button"
      onclick={selectColor('')}
      class="flex items-center justify-center w-8 h-8 rounded-full border-2 border-dashed border-gray-400 bg-white hover:scale-110 transition-transform dark:border-gray-500 dark:bg-gray-800 dark:ring-offset-gray-800 {value ===
      ''
        ? 'ring-2 ring-offset-2'
        : ''}"
      title={noneLabel}
      aria-pressed={value === ''}
    >
      <span class="w-2.5 h-2.5 rounded-full bg-gray-400 dark:bg-gray-500"></span>
      <span class="hidden">{noneLabel}</span>
    </button>
  {/if}

  {#each colorLegend as color}
    <button
      type="button"
      onclick={selectColor(color.color)}
      class="w-8 h-8 rounded-full colorcard-{color.color} hover:scale-110 transition-transform dark:ring-offset-gray-800 {value ===
      color.color
        ? 'ring-2 ring-offset-2'
        : ''}"
      title={colorTitle(color)}
      aria-pressed={value === color.color}
    >
      <span class="hidden">{colorTitle(color)}</span>
    </button>
  {/each}
</div>

<style lang="postcss">
  .colorcard-gray {
    @apply bg-gray-400;
    @apply ring-gray-400;
  }

  .colorcard-gray:hover {
    @apply bg-gray-600;
    @apply ring-gray-600;
  }

  .colorcard-red {
    @apply bg-red-400;
    @apply ring-red-400;
  }

  .colorcard-red:hover {
    @apply bg-red-600;
    @apply ring-red-600;
  }

  .colorcard-orange {
    @apply bg-orange-400;
    @apply ring-orange-400;
  }

  .colorcard-orange:hover {
    @apply bg-orange-600;
    @apply ring-orange-600;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
    @apply ring-yellow-400;
  }

  .colorcard-yellow:hover {
    @apply bg-yellow-600;
    @apply ring-yellow-600;
  }

  .colorcard-green {
    @apply bg-green-400;
    @apply ring-green-400;
  }

  .colorcard-green:hover {
    @apply bg-green-600;
    @apply ring-green-600;
  }

  .colorcard-teal {
    @apply bg-teal-400;
    @apply ring-teal-400;
  }

  .colorcard-teal:hover {
    @apply bg-teal-600;
    @apply ring-teal-600;
  }

  .colorcard-blue {
    @apply bg-blue-400;
    @apply ring-blue-400;
  }

  .colorcard-blue:hover {
    @apply bg-blue-600;
    @apply ring-blue-600;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
    @apply ring-indigo-400;
  }

  .colorcard-indigo:hover {
    @apply bg-indigo-600;
    @apply ring-indigo-600;
  }

  .colorcard-purple {
    @apply bg-purple-400;
    @apply ring-purple-400;
  }

  .colorcard-purple:hover {
    @apply bg-purple-600;
    @apply ring-purple-600;
  }

  .colorcard-pink {
    @apply bg-pink-400;
    @apply ring-pink-400;
  }

  .colorcard-pink:hover {
    @apply bg-pink-600;
    @apply ring-pink-600;
  }
</style>
