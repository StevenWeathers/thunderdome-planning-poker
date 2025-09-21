<script lang="ts">
  interface Props {
    type?: string;
    colspan?: number;
    class?: string;
    children?: import('svelte').Snippet;
  }

  let { type = 'default', children, colspan, class: className }: Props = $props();
</script>

{#if type === 'default'}
  <td class="px-4 py-3 whitespace-nowrap font-medium {className}" {colspan}>
    {@render children?.()}
  </td>
{:else if type === 'action'}
  <td class="px-4 py-3 whitespace-nowrap text-right font-medium {className}" {colspan}>
    {@render children?.()}
  </td>
{/if}
