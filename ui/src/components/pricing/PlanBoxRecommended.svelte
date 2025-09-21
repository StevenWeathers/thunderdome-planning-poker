<script>
  /** @type {{children?: import('svelte').Snippet}} */
  let { children } = $props();
</script>

<div
  class="bg-gradient-to-br from-purple-100 to-indigo-100 dark:from-purple-900 dark:to-indigo-900 rounded-2xl p-8 transform transition duration-500 hover:scale-105 hover:shadow-2xl border-2 border-purple-500 relative"
>
  <div class="absolute top-0 right-0 bg-purple-500 text-white text-xs font-bold px-3 py-1 rounded-bl-lg rounded-tr-lg">
    RECOMMENDED
  </div>
  {@render children?.()}
</div>
