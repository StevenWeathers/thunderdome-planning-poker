<script lang="ts">
  import { createBubbler } from 'svelte/legacy';

  const bubble = createBubbler();
  import { ChevronDown } from 'lucide-svelte';

  interface Props {
    value?: string;
    children?: import('svelte').Snippet;
    [key: string]: any
  }

  let { value = $bindable(''), children, ...rest }: Props = $props();
</script>

<div class="relative">
  <select
    bind:value="{value}"
    onchange={bubble('change')}
    class="block appearance-none w-full border border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-2 px-4 pe-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
    {...rest}
  >
    {@render children?.()}
  </select>
  <div
    class="pointer-events-none absolute inset-y-0 end-0 flex
                items-center px-2 text-gray-700 dark:text-gray-400"
  >
    <ChevronDown class="inline-block" />
  </div>
</div>
