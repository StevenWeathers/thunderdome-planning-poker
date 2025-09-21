<script lang="ts">
  interface Props {
    class?: string;
    children?: import('svelte').Snippet;
  }

  let { class: klass = 'border-gray-300 dark:border-gray-700', children }: Props = $props();
</script>

<div class="rounded rounded-b-lg shadow-xl border {klass}">
  <div class="rounded-t bg-gray-200 dark:bg-gray-800 px-4 py-2.5">
    <div class="flex flex-row space-x-2">
      <div class="h-3 w-3 rounded-full bg-rose-500"></div>
      <div class="h-3 w-3 rounded-full bg-amber-300"></div>
      <div class="h-3 w-3 rounded-full bg-lime-400"></div>
    </div>
  </div>
  <div class="rounded-b-lg">
    {@render children?.()}
  </div>
</div>
