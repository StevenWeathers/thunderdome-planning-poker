<script lang="ts">
  interface Props {
    class?: string;
    children?: import('svelte').Snippet;
  }

  let { class: klass = 'w-56', children }: Props = $props();
  
</script>

<ul class="absolute right-0 p-2 mt-2 space-y-2 text-gray-600 bg-white border border-gray-100 rounded-md shadow-md dark:border-gray-700 dark:text-gray-300 dark:bg-gray-700 {klass}"
    aria-label="submenu">
    {@render children?.()}
</ul>