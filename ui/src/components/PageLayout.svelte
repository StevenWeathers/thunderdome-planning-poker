<script>
  /** @type {{children?: import('svelte').Snippet}} */
  let { children } = $props();
</script>

<div class="container mx-auto px-4 py-4 md:py-6 lg:py-8">
  {@render children?.()}
</div>
