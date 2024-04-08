<script lang="ts">
  export let type = 'default';
</script>

{#if type === 'default'}
  <th scope="col" class="p-4">
    <slot />
  </th>
{:else if type === 'action'}
  <th scope="col" class="relative px-2 py-2">
    <slot />
  </th>
{/if}
