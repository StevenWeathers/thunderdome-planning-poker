<script lang="ts">
  import type { Component, Snippet } from 'svelte';
  
  interface Props {
    type?: string;
    value?: any;
    icon?: Component;
    iconPosition?: 'start' | 'end';
    startElement?: Snippet;
    endElement?: Snippet;
    class?: string;
    [key: string]: any; // For rest props
  }

  let {
    type = 'text',
    value = $bindable(),
    icon,
    iconPosition = 'end',
    startElement,
    endElement,
    class: klass = '',
    ...restProps
  }: Props = $props();

  let inputElement: HTMLInputElement;

  export function focus() {
    inputElement?.focus();
  }

  // Works around "svelte(invalid-type)" warning, i.e., can't have a dynamic type AND bind:value
  // Keep an eye on https://github.com/sveltejs/svelte/issues/3921
  const typeWorkaround = (node: HTMLInputElement) => {
    node.type = type;
  };

  // Calculate padding based on icon position and elements
  const getInputPadding = () => {
    const hasStart = startElement || (icon && iconPosition === 'start');
    const hasEnd = endElement || (icon && iconPosition === 'end');
    
    if (hasStart && hasEnd) {
      return 'ps-12 pe-12';
    } else if (hasStart) {
      return 'ps-12 pe-5';
    } else if (hasEnd) {
      return 'ps-5 pe-12';
    } else {
      return 'px-5';
    }
  };
</script>

<div class="relative">
  <input
    use:typeWorkaround
    bind:this={inputElement}
    bind:value={value}
    class="block w-full {getInputPadding()} py-3 text-lg rounded-lg outline-none transition-all duration-300 bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-400 disabled:cursor-not-allowed {klass}"
    {...restProps}
  />
  
  <!-- Start side elements -->
  {#if startElement}
    <div class="absolute start-3 top-1/2 transform -translate-y-1/2 pointer-events-none">
      {@render startElement()}
    </div>
  {:else if icon && iconPosition === 'start'}
    {@const SvelteComponent = icon}
    <SvelteComponent
      class="absolute top-1/2 start-3 transform -translate-y-1/2 text-gray-500 dark:text-gray-400"
      size={24}
      tabindex={-1}
    />
  {/if}
  
  <!-- End side elements -->
  {#if endElement}
    <div class="absolute end-3 top-1/2 transform -translate-y-1/2 pointer-events-none">
      {@render endElement()}
    </div>
  {:else if icon && iconPosition === 'end'}
    {@const SvelteComponent = icon}
    <SvelteComponent
      class="absolute top-1/2 end-3 transform -translate-y-1/2 text-gray-500 dark:text-gray-400"
      size={24}
      tabindex={-1}
    />
  {/if}
</div>