<script module lang="ts">
  // open modals, topmost last, so that Escape only closes the one on top
  const openModals: symbol[] = [];
  let pageOverflow = '';
</script>

<script lang="ts">
  import { X } from '@lucide/svelte';
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
    ariaDescribedby,
  }: Props = $props();

  const modalId = Symbol();

  const handleKeydown = (e: KeyboardEvent) => {
    // a held key must not close the whole stack
    if (e.key !== 'Escape' || e.repeat || e.isComposing || e.defaultPrevented) return;
    if (openModals[openModals.length - 1] !== modalId) return;

    e.preventDefault();
    // fires the focused field's change event, so an edit in progress is saved as it is with the close button
    if (modalElement.contains(document.activeElement)) {
      (document.activeElement as HTMLElement).blur();
    }
    closeModal();
  };

  let modalElement: HTMLElement;
  let focusTrap: any;

  onMount(() => {
    // keeps the page behind from scrolling, and showing a second scrollbar, while any modal is open
    if (openModals.length === 0) {
      pageOverflow = document.documentElement.style.overflow;
      document.documentElement.style.overflow = 'hidden';
    }
    openModals.push(modalId);
    focusTrap = createFocusTrap(modalElement, {
      escapeDeactivates: false, // Handle escape separately
      returnFocusOnDeactivate: true,
    });
    focusTrap.activate();

    return () => {
      openModals.splice(openModals.indexOf(modalId), 1);
      if (openModals.length === 0) {
        document.documentElement.style.overflow = pageOverflow;
      }
      focusTrap.deactivate();
    };
  });
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="fixed z-50 inset-0 flex items-center z-40 max-h-screen overflow-y-auto overscroll-contain">
  <!-- Background overlay -->
  <div class="fixed inset-0 bg-gray-900 opacity-75" aria-hidden="true"></div>

  <div
    class="relative z-[100] mx-4 md:mx-auto w-full {widthClasses != '' ? widthClasses : 'md:w-2/3 lg:w-3/5 xl:w-1/3'}
        max-h-full"
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
