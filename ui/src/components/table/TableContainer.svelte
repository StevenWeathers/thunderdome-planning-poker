<script>
  /** @type {{children?: import('svelte').Snippet}} */
  let { children } = $props();
</script>

<div class="bg-white dark:bg-gray-800 shadow-md sm:rounded-lg overflow-hidden">
  {@render children?.()}
</div>
