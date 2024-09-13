<script lang="ts">
  import { ThumbsUp } from 'lucide-svelte';

  export let groups = [];
  export let columns = [];

  $: columnColors =
    columns &&
    columns.reduce((p, c) => {
      p[c.name] = c.color;
      return p;
    }, {});
</script>

{#each groups as group, i (group.id)}
  {#if group.items.length > 0}
    <div
      class="border-2 p-2 dark:border-gray-800 rounded flex flex-col flex-wrap"
    >
      <div class="dark:text-gray-200 w-full text-center text-lg mb-2">
        <div class="flex content-center justify-center">
          <div class="inline-block align-middle text-2xl ms-2">
            <ThumbsUp class="w-6 h-6 inline-block" />
            {group.voteCount}
          </div>
        </div>
        {group.name ? group.name : 'Group'}
      </div>
      <div class="flex-1 grow">
        {#each group.items as item, ii (item.id)}
          <div
            class="p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-s-4 dark:text-white"
            class:border-green-400="{columnColors[item.type] === 'green'}"
            class:dark:border-lime-400="{columnColors[item.type] === 'green'}"
            class:border-red-500="{columnColors[item.type] === 'red'}"
            class:border-blue-400="{columnColors[item.type] === 'blue'}"
            class:dark:border-sky-400="{columnColors[item.type] === 'blue'}"
            class:border-yellow-500="{columnColors[item.type] === 'yellow'}"
            class:dark:border-yellow-400="{columnColors[item.type] ===
              'yellow'}"
            class:border-orange-500="{columnColors[item.type] === 'orange'}"
            class:dark:border-orange-400="{columnColors[item.type] ===
              'orange'}"
            class:border-teal-500="{columnColors[item.type] === 'teal'}"
            class:dark:border-teal-400="{columnColors[item.type] === 'teal'}"
            class:border-indigo-500="{columnColors[item.type] === 'purple'}"
            class:dark:border-indigo-400="{columnColors[item.type] ===
              'purple'}"
          >
            {item.content}
          </div>
        {/each}
      </div>
    </div>
  {/if}
{/each}
