<script>
  /** @type {{header?: import('svelte').Snippet, body?: import('svelte').Snippet<[{ class: string }]>}} */
  let { header, body } = $props();
</script>

<div class="overflow-x-auto">
  <table class="w-full text-sm text-left text-gray-700 dark:text-gray-200">
    <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
      {@render header?.()}
    </thead>
    {@render body?.({ class: '' })}
  </table>
</div>
