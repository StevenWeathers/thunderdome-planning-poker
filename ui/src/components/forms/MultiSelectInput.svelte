<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    value?: string[];
    children?: Snippet;
    onchange?: (event: Event) => void;
    [key: string]: any;
  }

  let { value = $bindable<string[]>([]), children, onchange, ...restProps }: Props = $props();

  let selectElement: HTMLSelectElement;

  function syncSelectedOptions(selectedValues: string[]) {
    if (!selectElement) {
      return;
    }

    const selectedValueSet = new Set(selectedValues);

    for (const option of Array.from(selectElement.options)) {
      option.selected = selectedValueSet.has(option.value);
    }
  }

  function handleChange(event: Event) {
    value = Array.from(selectElement.selectedOptions, option => option.value);
    onchange?.(event);
  }

  export function focus() {
    selectElement?.focus();
  }

  $effect(() => {
    syncSelectedOptions(value ?? []);
  });
</script>

<select
  multiple
  bind:this={selectElement}
  class="block w-full border border-gray-300 dark:border-gray-700 text-gray-700 dark:text-gray-300 py-2 px-4 rounded leading-tight focus:outline-none focus:border-indigo-500 dark:focus:border-yellow-400 dark:bg-gray-900"
  onchange={handleChange}
  {...restProps}
>
  {@render children?.()}
</select>
