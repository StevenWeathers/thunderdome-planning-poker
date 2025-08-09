<script lang="ts">
  import { type RetroTemplateFormat } from '../../types/retro';
  import Modal from '../global/Modal.svelte';

  interface Props {
    format: RetroTemplateFormat;
    toggleClose: (template) => void;
  }

  let { format, toggleClose }: Props = $props();
</script>

<Modal closeModal={() => toggleClose(null)} widthClasses="w-full md:w-2/3">
  <div class="mt-10 grid md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each format.columns as column}
      <div
        class="p-2 bg-gray-100 dark:bg-gray-700 dark:text-gray-200 rounded-lg shadow"
      >
        <h2 class="text-xl font-bold mb-2">column</h2>
        <ul class="list-none ps-4">
          <li><strong>name:</strong> {column.name}</li>
          <li><strong>label:</strong> {column.label}</li>
          <li><strong>color:</strong> {column.color}</li>
          <li><strong>icon:</strong> {column.icon}</li>
        </ul>
      </div>
    {/each}
  </div>
</Modal>
