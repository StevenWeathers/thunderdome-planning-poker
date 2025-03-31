<script lang="ts">
  interface Props {
    type?: string;
    children?: import('svelte').Snippet;
  }

  let { type = 'default', children }: Props = $props();
</script>

{#if type === 'default'}
  <th scope="col" class="p-4">
    {@render children?.()}
  </th>
{:else if type === 'action'}
  <th scope="col" class="relative px-2 py-2">
    {@render children?.()}
  </th>
{/if}
