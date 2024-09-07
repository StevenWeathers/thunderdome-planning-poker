<script lang="ts">
  import LL from '../../i18n/i18n-svelte';

  export let retro = {
    items: [],
    actionItems: [],
  };
</script>

<div class="flex flex-grow p-4 dark:text-white">
  <div class="px-4">
    {#each retro.template.format.columns as column (column.id)}
      <div class="mb-4">
        <h2 class="text-3xl font-rajdhani">{column.label}</h2>
        <ul class="ps-12 text-lg">
          {#each retro.items as item (item.id)}
            {#if item.column_id === column.id}
              <li>
                &bull; {item.content}
              </li>
            {/if}
          {/each}
        </ul>
      </div>
    {/each}
    <div class="mb-4">
      <h2 class="text-3xl font-rajdhani">{$LL.actionItems()}</h2>
      <ul class="ps-12 text-lg">
        {#each retro.actionItems as item (item.id)}
          <li>&bull; {item.content}</li>
        {/each}
      </ul>
    </div>
  </div>
</div>
