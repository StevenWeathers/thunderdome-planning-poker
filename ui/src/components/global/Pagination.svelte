<script lang="ts">
  import { run } from 'svelte/legacy';

  import { createEventDispatcher } from 'svelte';
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';

  const dispatch = createEventDispatcher();

  interface Props {
    current?: number;
    num_items?: number;
    per_page?: number;
  }

  let { current = $bindable(1), num_items = 120, per_page = 5 }: Props = $props();

  let num_pages;
  run(() => {
    num_pages = Math.ceil(num_items / per_page);
  });

  let arr_pages = $state([]);

  function buildArr(c, n) {
    if (n <= 7) {
      return [...Array(n)].map((_, i) => i + 1);
    } else {
      if (c < 3 || c > n - 2) {
        return [1, 2, 3, '...', n - 2, n - 1, n];
      } else {
        return [1, '...', c - 1, c, c + 1, '...', n];
      }
    }
  }

  function setArrPages() {
    arr_pages = buildArr(current, num_pages);
  }

  run(() => {
    if (current) {
      setArrPages();
    }
  });

  run(() => {
    if (per_page) {
      setArrPages();
      current = 1;
    }
  });

  run(() => {
    if (num_items) {
      num_pages = Math.ceil(num_items / per_page);
      setArrPages();
      current = current || 1;
    }
  });

  function setCurrent(i) {
    if (isNaN(i)) return;
    current = i;
    dispatch('navigate', current);
  }
</script>

<div class="flex text-gray-700 dark:text-gray-400 text-lg">
  <div
    class="h-12 w-12 me-1 flex justify-center items-center {current > 1
      ? 'cursor-pointer'
      : 'text-gray-400 dark:text-gray-700'}"
    role="button"
    tabindex="0"
    onclick={() => current > 1 && setCurrent(current - 1)}
    onkeypress={() => current > 1 && setCurrent(current - 1)}
  >
    <ChevronLeft class="w-6 h-6 inline-block" />
  </div>
  <div class="flex h-12 font-medium">
    {#each arr_pages as i}
      <div
        class="w-12 sm:flex justify-center items-center hidden
                select-none cursor-pointer leading-5 transition duration-150
                ease-in {i == current
          ? `border-t-2 border-indigo-600 dark:border-yellow-400`
          : 'border-t-2 border-slate-100 dark:border-gray-900'}
                "
        role="button"
        tabindex="0"
        onclick={() => setCurrent(i)}
        onkeypress={() => setCurrent(i)}
      >
        {i}
      </div>
    {/each}
    <div
      class="w-12 h-12 sm:hidden flex justify-center select-none
            items-center cursor-pointer leading-5 transition duration-150
            ease-in border-t-2 border-indigo-600"
    >
      {current}
    </div>
  </div>
  <div
    class="h-12 w-12 ms-1 flex justify-center items-center {current < num_pages
      ? 'cursor-pointer'
      : 'text-gray-400 dark:text-gray-700'}"
    role="button"
    tabindex="0"
    onclick={() => current < num_pages && setCurrent(current + 1)}
    onkeypress={() => current < num_pages && setCurrent(current + 1)}
  >
    <ChevronRight class="w-6 h-6 inline-block" />
  </div>
</div>
