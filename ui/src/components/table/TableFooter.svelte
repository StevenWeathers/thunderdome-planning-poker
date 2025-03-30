<script lang="ts">
  import { run } from 'svelte/legacy';

  import { ChevronLeft, ChevronRight } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';

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

  let showingBegin = $derived(num_items === 0 ? 0 : (current - 1) * per_page + 1);
  let showingEnd =
    $derived(current === num_pages
      ? num_items
      : num_items === 0
      ? 0
      : current * per_page);
</script>

<nav
  class="flex flex-col md:flex-row justify-between items-start md:items-center space-y-3 md:space-y-0 p-4"
  aria-label="Table navigation"
>
  <span class="text-sm font-normal text-gray-500 dark:text-gray-400">
    Showing
    <span class="font-semibold text-gray-900 dark:text-white"
      >{showingBegin}
      - {showingEnd}</span
    >
    of
    <span class="font-semibold text-gray-900 dark:text-white">{num_items}</span>
  </span>
  <ul class="inline-flex items-stretch -space-x-px">
    <li>
      <button
        class="flex items-center justify-center h-full py-1.5 px-3 ml-0 text-gray-500 bg-white rounded-l-lg border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
        onclick={() => current > 1 && setCurrent(current - 1)}
        onkeypress={() => current > 1 && setCurrent(current - 1)}
      >
        <span class="sr-only">Previous</span>
        <ChevronLeft class="w-5 h-5 inline-block" />
      </button>
    </li>
    {#each arr_pages as i}
      <li>
        {#if i === current}
          <div
            class="flex items-center justify-center text-sm z-10 py-2 px-3 leading-tight text-primary-600 bg-primary-50 border border-primary-300 hover:bg-primary-100 hover:text-primary-700 dark:border-gray-700 dark:bg-gray-700 dark:text-white"
            aria-current="page"
          >
            {i}
          </div>
        {:else if i === '...'}
          <div
            class="flex items-center justify-center text-sm py-2 px-3 leading-tight text-gray-500 bg-white border border-gray-300 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400"
          >
            {i}
          </div>
        {:else}
          <button
            class="flex items-center justify-center text-sm py-2 px-3 leading-tight text-gray-500 bg-white border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
            onclick={() => setCurrent(i)}
            onkeypress={() => setCurrent(i)}
          >
            {i}
          </button>
        {/if}
      </li>
    {/each}
    <li>
      <button
        class="flex items-center justify-center h-full py-1.5 px-3 leading-tight text-gray-500 bg-white rounded-r-lg border border-gray-300 hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white"
        onclick={() => current < num_pages && setCurrent(current + 1)}
        onkeypress={() => current < num_pages && setCurrent(current + 1)}
      >
        <span class="sr-only">Next</span>
        <ChevronRight class="w-5 h-5 inline-block" />
      </button>
    </li>
  </ul>
</nav>
