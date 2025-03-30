<script lang="ts">
  import { createBubbler } from 'svelte/legacy';

  const bubble = createBubbler();
  interface Props {
    type?: string;
    value: string;
    icon: any;
    class?: string;
    [key: string]: any
  }

  let {
    type = 'text',
    value = $bindable(),
    icon,
    class: klass = '',
    ...rest
  }: Props = $props();
  

  let inputElement = $state();

  export function focus() {
    inputElement.focus();
  }

  // works around "svelte(invalid-type)" warning, i.e., can't have a dynamic type AND bind:value...keep an eye on https://github.com/sveltejs/svelte/issues/3921
  const typeWorkaround = node => (node.type = type);
</script>

<div class="relative">
  <input
    use:typeWorkaround
    bind:this="{inputElement}"
    bind:value="{value}"
    onchange={bubble('change')}
    oninput={bubble('input')}
    class="block w-full px-5 py-3 text-lg rounded-lg outline-none transition-all duration-300 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-400 disabled:cursor-not-allowed {klass}"
    {...rest}
  />
  {#if icon}
    {@const SvelteComponent = icon}
    <SvelteComponent
      class="absolute top-3 right-3 text-gray-500 dark:text-gray-400"
      size="{24}"
      tabindex="-1"
    />
  {/if}
</div>
