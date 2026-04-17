<script lang="ts">
  import { Smile } from '@lucide/svelte';
  import type { EmojiPickerItem } from './emoji-picker';

  interface Props {
    options?: EmojiPickerItem[];
    disabled?: boolean;
    onToggle?: (option: EmojiPickerItem) => void;
    toggleAriaLabel?: string;
    toggleTestId?: string;
    popoverTestId?: string;
  }

  let {
    options = [],
    disabled = false,
    onToggle = () => {},
    toggleAriaLabel = 'Open emoji picker',
    toggleTestId = 'emoji-picker-toggle',
    popoverTestId = 'emoji-picker-popover',
  }: Props = $props();

  let showPicker = $state(false);

  const activeOptions = $derived(options.filter(option => option.count > 0));
  const pickerOptions = $derived(
    options.map(option => ({
      ...option,
      pickerDisabled: disabled || option.pickerDisabled,
    })),
  );

  function togglePicker() {
    if (disabled) {
      return;
    }

    showPicker = !showPicker;
  }

  function closePicker() {
    showPicker = false;
  }

  function clickOutsidePicker(node: HTMLElement) {
    function handleClick(event: MouseEvent) {
      if (showPicker && !node.contains(event.target as Node)) {
        closePicker();
      }
    }

    document.addEventListener('click', handleClick, true);

    return {
      destroy() {
        document.removeEventListener('click', handleClick, true);
      },
    };
  }

  const getReactionClasses = (selected: boolean) =>
    selected
      ? 'inline-flex items-center gap-1 rounded-full border border-blue-400 bg-blue-50 px-2.5 py-1 text-sm text-blue-700 transition-colors duration-200 dark:border-sky-500 dark:bg-sky-500/10 dark:text-sky-300'
      : 'inline-flex items-center gap-1 rounded-full border border-gray-200 bg-white px-2.5 py-1 text-sm text-gray-700 transition-colors duration-200 hover:border-blue-300 hover:text-blue-600 dark:border-gray-600 dark:bg-gray-800/60 dark:text-gray-200 dark:hover:border-sky-500 dark:hover:text-sky-300';

  const getPickerButtonClasses = (optionDisabled: boolean) =>
    optionDisabled
      ? 'flex items-center justify-center rounded-lg px-2 py-2 text-xl opacity-40 cursor-not-allowed'
      : 'flex items-center justify-center rounded-lg px-2 py-2 text-xl hover:bg-gray-100 dark:hover:bg-gray-700/70 transition-colors duration-150';

  function handleToggle(option: EmojiPickerItem) {
    if (disabled) {
      return;
    }

    closePicker();
    onToggle(option);
  }
</script>

<div class="flex flex-wrap items-center gap-2" data-testid="emoji-picker">
  {#each activeOptions as option}
    <button
      type="button"
      class={getReactionClasses(Boolean(option.selected))}
      class:cursor-not-allowed={disabled}
      aria-label="{option.selected ? 'Remove' : 'Add'} {option.label} reaction"
      data-testid="reaction-{option.key}"
      onclick={() => handleToggle(option)}
      {disabled}
    >
      <span aria-hidden="true">{option.value}</span>
      <span class="font-medium">{option.count}</span>
    </button>
  {/each}

  <div class="relative" use:clickOutsidePicker>
    <button
      type="button"
      class="inline-flex items-center justify-center rounded-full border border-gray-200 bg-white p-2 text-gray-600 transition-colors duration-200 hover:border-blue-300 hover:text-blue-600 dark:border-gray-600 dark:bg-gray-800/60 dark:text-gray-200 dark:hover:border-sky-500 dark:hover:text-sky-300"
      class:cursor-not-allowed={disabled}
      aria-label={toggleAriaLabel}
      aria-expanded={showPicker}
      data-testid={toggleTestId}
      onclick={togglePicker}
      {disabled}
    >
      <Smile class="h-4 w-4" />
    </button>

    {#if showPicker}
      <div
        class="absolute left-0 top-full z-10 mt-2 flex gap-1 rounded-2xl border border-gray-200 bg-white p-2 shadow-lg dark:border-gray-600 dark:bg-gray-800"
        data-testid={popoverTestId}
      >
        {#each pickerOptions as option}
          <button
            type="button"
            class={getPickerButtonClasses(Boolean(option.pickerDisabled))}
            aria-label="Add {option.label} reaction"
            data-testid="reaction-picker-option-{option.key}"
            onclick={() => handleToggle(option)}
            disabled={option.pickerDisabled}
          >
            <span aria-hidden="true">{option.value}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
</div>
