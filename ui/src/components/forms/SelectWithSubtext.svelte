<script lang="ts">
  import { run } from 'svelte/legacy';

  import { Check, ChevronDown } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';

  interface Props {
    items?: any;
    nameField?: string;
    descriptionField?: string;
    label?: string;
    selectedItemId: any;
    itemType?: string;
  }

  let {
    items = [],
    nameField = 'name',
    descriptionField = 'description',
    label = 'Select an item',
    selectedItemId,
    itemType = ''
  }: Props = $props();

  const dispatch = createEventDispatcher();

  let isOpen = $state(false);
  let selectedItemIdx = $state(-1);

  function toggleDropdown(e) {
    e.preventDefault();
    isOpen = !isOpen;
  }

  function handleChange(itemIdx) {
    selectedItemIdx = itemIdx;
    isOpen = false;
    dispatch('change', items[itemIdx]);
  }

  run(() => {
    if (selectedItemId) {
      selectedItemIdx = items.findIndex(item => item.id === selectedItemId);
    }
  });
</script>

<div class="relative">
  <button
    type="button"
    onclick={toggleDropdown}
    class="w-full flex justify-between px-4 py-2 text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
  >
    <span class="flex">
      <Check class="w-6 h-6 mr-2 inline-block" />
      <span data-testid="{itemType}_item-selected"
        >{selectedItemIdx > -1
          ? items[selectedItemIdx][nameField]
          : label}</span
      >
    </span>
    <ChevronDown class="w-6 h-6 ml-2 inline-block" />
  </button>

  {#if isOpen}
    <div
      class="absolute z-10 w-full mt-2 bg-white dark:bg-gray-800 rounded-md shadow-lg dark:border dark:border-gray-700"
    >
      <div class="py-1">
        {#each items as item, idx}
          <button
            type="button"
            onclick={() => handleChange(idx)}
            class="flex content-center w-full px-4 py-2 text-gray-700 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-200 dark:hover:text-gray-100 dark:hover:bg-gray-700"
          >
            <Check
              class="flex-none w-6 h-6 mr-2 {selectedItemIdx === idx
                ? 'text-indigo-600 dark:text-indigo-400'
                : 'invisible'}"
            />
            <div class="flex-grow text-left">
              <div class="font-medium" data-testid="{itemType}_item-name">
                {item[nameField]}
              </div>
              <div
                class="text-sm text-gray-500 dark:text-gray-400"
                data-testid="{itemType}_item-description"
              >
                {item[descriptionField]}
              </div>
            </div>
          </button>
        {/each}
      </div>
    </div>
  {/if}
</div>
