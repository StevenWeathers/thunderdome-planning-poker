<script lang="ts">
  import Modal from './Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { TriangleAlert } from 'lucide-svelte';

  interface Props {
    confirmText?: string;
    confirmBtnText?: string;
    permanent?: boolean;
    handleDelete?: any;
    toggleDelete?: any;
  }

  let {
    confirmText = '',
    confirmBtnText = 'Confirm Delete',
    permanent = true,
    handleDelete = () => {},
    toggleDelete = () => {}
  }: Props = $props();
</script>

<Modal closeModal={toggleDelete} ariaLabel={$LL.modalDeleteConfirmation()}>
  <div class="mb-4">
    <div
      class="w-12 h-12 rounded-lg bg-gray-200 dark:bg-gray-700 p-2 flex items-center justify-center mx-auto mb-3.5"
    >
      <TriangleAlert class="w-8 h-8 text-gray-600 dark:text-gray-400" />
    </div>
    <div
      class="text-center text-gray-800 dark:text-gray-100 mb-4 font-semibold text-lg"
    >
      {confirmText}
    </div>
    <div class="text-xl text-center text-red-400 dark:text-[#fdba8c]">
      {#if permanent}
        {$LL.cannotBeUndone()}
      {/if}
    </div>
  </div>
  <div class="flex justify-center items-center space-x-4">
    <button
      type="button"
      onclick={toggleDelete}
      data-testid="confirm-cancel"
      class="py-2 px-3 text-sm font-medium text-gray-500 bg-white rounded-lg border border-gray-200 hover:bg-gray-100 focus:ring-4 focus:outline-none focus:ring-primary-300 hover:text-gray-900 focus:z-10 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-500 dark:hover:text-white dark:hover:bg-gray-600 dark:focus:ring-gray-600"
    >
      {$LL.cancel()}
    </button>
    <button
      onclick={handleDelete}
      data-testid="confirm-confirm"
      class="py-2 px-3 text-sm font-medium text-center text-white bg-red-600 rounded-lg hover:bg-red-700 focus:ring-4 focus:outline-none focus:ring-red-300 dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-900"
    >
      {confirmBtnText}
    </button>
  </div>
</Modal>
