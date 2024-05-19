<script lang="ts">
  export let type: string = 'text';
  export let value: any;
  let klass: string = '';
  export { klass as class };

  let inputElement;
  export function focus() {
    inputElement.focus();
  }

  // works around "svelte(invalid-type)" warning, i.e., can't have a dynamic type AND bind:value...keep an eye on https://github.com/sveltejs/svelte/issues/3921
  const typeWorkaround = node => (node.type = type);
</script>

<input
  use:typeWorkaround
  bind:this="{inputElement}"
  bind:value="{value}"
  on:change
  on:input
  class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 {klass}"
  {...$$restProps}
/>
