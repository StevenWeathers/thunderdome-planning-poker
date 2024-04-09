<script lang="ts">
  export let type = 'default';
</script>

{#if type === 'default'}
  <td class="px-4 py-3 whitespace-nowrap font-medium">
    <slot />
  </td>
{:else if type === 'action'}
  <td class="px-4 py-3 whitespace-nowrap text-right font-medium">
    <slot />
  </td>
{/if}
