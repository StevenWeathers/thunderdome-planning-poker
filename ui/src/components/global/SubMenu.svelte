<script lang="ts">
  import SolidButton from './SolidButton.svelte';

  interface Props {
    class?: string;
    relativeClass?: string;
    children?: import('svelte').Snippet<[{ toggleSubmenu: () => void }]>;
    button?: import('svelte').Snippet<[{ toggleSubmenu: () => void }]>;
    label?: string;
    icon?: any;
    testId?: string;
  }

  let { class: klass = 'w-56', relativeClass, children, button, label, icon, testId = 'settings' }: Props = $props();
  let showSubmenu = $state(false);

  const toggleSubmenu = () => {
    showSubmenu = !showSubmenu;
  };

  function clickOutsideSubmenu(node: HTMLElement) {
    function handleClick(event: MouseEvent) {
      if (showSubmenu && !node.contains(event.target as Node)) {
        toggleSubmenu();
      }
    }

    document.addEventListener('click', handleClick, true);

    return {
      destroy() {
        document.removeEventListener('click', handleClick, true);
      }
    };
  };
</script>

<div class="inline-block relative {relativeClass}" use:clickOutsideSubmenu>
  {#if button}
    {@render button({ toggleSubmenu })}
  {:else}
    <SolidButton onClick={toggleSubmenu} color="blue" testid={testId} options={{"aria-haspopup": true, "aria-label": label}}>
      {#if icon}{@const SvelteComponent = icon}
      <SvelteComponent class="inline-block h-3.5 w-3.5 me-1.5 -ms-1" />{/if}{label}
    </SolidButton>
  {/if}

  {#if showSubmenu}
    <ul class="absolute right-0 p-2 mt-2 space-y-2 text-gray-600 bg-white border border-gray-100 rounded-md shadow-md dark:border-gray-700 dark:text-gray-300 dark:bg-gray-700 {klass}"
        aria-label="submenu">
        {@render children?.({ toggleSubmenu })}
    </ul>
  {/if}
</div>