<script lang="ts">
  interface Props {
    itemIndex?: number;
    children?: import('svelte').Snippet;
  }

  let { itemIndex = 0, children }: Props = $props();
</script>

<tr class="border-b dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700">
  {@render children?.()}
</tr>
