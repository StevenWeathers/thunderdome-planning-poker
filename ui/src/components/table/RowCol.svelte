<script lang="ts">
  interface Props {
    type?: string;
    children?: import('svelte').Snippet;
  }

  let { type = 'default', children }: Props = $props();
</script>

{#if type === 'default'}
  <td class="px-4 py-3 whitespace-nowrap font-medium">
    {@render children?.()}
  </td>
{:else if type === 'action'}
  <td class="px-4 py-3 whitespace-nowrap text-right font-medium">
    {@render children?.()}
  </td>
{/if}
