<script lang="ts">
  import { X } from 'lucide-svelte';
  import { createFocusTrap } from 'focus-trap';
  import { onMount } from 'svelte';

  interface Props {
    closeModal?: any;
    widthClasses?: string;
    children?: import('svelte').Snippet;
    ariaLabel?: string;
    ariaLabelledby?: string;
    ariaDescribedby?: string;
  }

  let { 
    closeModal = () => {}, 
    widthClasses = '', 
    children,
    ariaLabel,
    ariaLabelledby,
    ariaDescribedby
  }: Props = $props();

  const handle_keydown = (e) => {
    if (e.key === 'Escape') return closeModal();
  };

  let modalElement: any;
  let focusTrap: any;
  
  onMount(() => {
    focusTrap = createFocusTrap(modalElement, {
      escapeDeactivates: false, // Handle escape separately
      returnFocusOnDeactivate: true
    });
    focusTrap.activate();
    
    return () => focusTrap.deactivate();
  });
</script>

<svelte:window on:keydown|once={handle_keydown} />

<div
  class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll"
>
  <!-- Background overlay -->
  <div class="fixed inset-0 bg-gray-900 opacity-75" aria-hidden="true"></div>

  <div
    class="relative mx-4 md:mx-auto w-full {widthClasses != ''
      ? widthClasses
      : 'md:w-2/3 lg:w-3/5 xl:w-1/3'}
        z-50 max-h-full"
  >
    <div class="py-8">
      <div
        class="relative shadow-xl bg-white dark:bg-gray-800 rounded-lg p-4 xl:p-6 max-h-full"
        role="dialog"
        aria-modal="true"
        aria-label={ariaLabel}
        aria-labelledby={ariaLabelledby}
        aria-describedby={ariaDescribedby}
        bind:this={modalElement}
      >
        <button
          onclick={closeModal}
          aria-label="Close modal"
          type="button"
          class="text-gray-400 absolute top-2.5 right-2.5 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ms-auto inline-flex items-center dark:hover:bg-gray-600 dark:hover:text-white"
        >
          <X class="w-5 h-5" aria-hidden="true" />
          <span class="sr-only">Close modal</span>
        </button>
        <div>
          {@render children?.()}
        </div>
      </div>
    </div>
  </div>
</div>
