<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { RetroTemplateColumn, RetroTemplateFormat } from '../../types/retro';
  import { Angry, CircleHelp, Frown, Smile } from 'lucide-svelte';

  export let format: RetroTemplateFormat;

  const dispatch = createEventDispatcher();
  const MIN_COLUMNS = 2;
  const MAX_COLUMNS = 5;

  const colorOptions = [
    { name: 'Red', value: 'red' },
    { name: 'Blue', value: 'blue' },
    { name: 'Green', value: 'green' },
    { name: 'Yellow', value: 'yellow' },
    { name: 'Purple', value: 'purple' },
    { name: 'Teal', value: 'teal' },
    { name: 'Orange', value: 'orange' },
  ];

  const iconOptions = [
    { name: 'Smiley', value: 'smiley', component: Smile },
    { name: 'Frown', value: 'frown', component: Frown },
    { name: 'Question', value: 'question', component: CircleHelp },
    { name: 'Angry', value: 'angry', component: Angry },
  ];

  let newColumn: RetroTemplateColumn = {
    name: '',
    label: '',
    color: '',
    icon: '',
  };

  $: canAddColumn = format && format.columns.length < MAX_COLUMNS;
  $: canRemoveColumn = format && format.columns.length > MIN_COLUMNS;

  function addColumn(event: Event) {
    event.preventDefault();
    if (newColumn.name && newColumn.label && canAddColumn) {
      format.columns = [...format.columns, { ...newColumn }];
      dispatch('update', format);
      resetNewColumn();
    }
  }

  function updateColumn(
    index: number,
    field: keyof RetroTemplateColumn,
    value: string,
  ) {
    format.columns[index][field] = value;
    format.columns = [...format.columns];
    dispatch('update', format);
  }

  function removeColumn(index: number) {
    if (canRemoveColumn) {
      format.columns = format.columns.filter((_, i) => i !== index);
      dispatch('update', format);
    }
  }

  function resetNewColumn() {
    newColumn = {
      name: '',
      label: '',
      color: '',
      icon: '',
    };
  }

  function getIconComponent(iconValue: string) {
    return (
      iconOptions.find(option => option.value === iconValue)?.component || Smile
    );
  }
</script>

<div class="space-y-4">
  <h2 class="text-2xl font-bold dark:text-white">Manage Columns</h2>

  <p class="text-sm text-gray-600 dark:text-gray-400">
    You must have at least {MIN_COLUMNS} columns and can have up to {MAX_COLUMNS}
    columns.
  </p>

  {#each format.columns as column, index}
    <div class="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg shadow">
      <input
        bind:value="{column.name}"
        on:input="{() => updateColumn(index, 'name', column.name)}"
        placeholder="Column name"
        class="w-full p-2 mb-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
      />
      <input
        bind:value="{column.label}"
        on:input="{() => updateColumn(index, 'label', column.label)}"
        placeholder="Column label"
        class="w-full p-2 mb-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
      />
      <select
        bind:value="{column.color}"
        on:change="{() => updateColumn(index, 'color', column.color)}"
        class="w-full p-2 mb-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
      >
        <option value="">Select an optional color</option>
        {#each colorOptions as option}
          <option value="{option.value}">{option.name}</option>
        {/each}
      </select>
      <div class="flex items-center mb-2">
        <select
          bind:value="{column.icon}"
          on:change="{() => updateColumn(index, 'icon', column.icon)}"
          class="flex-grow p-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
        >
          <option value="">Select an optional icon</option>
          {#each iconOptions as option}
            <option value="{option.value}">{option.name}</option>
          {/each}
        </select>
        <div
          class="ml-2 p-2 bg-white dark:bg-gray-700 border rounded dark:border-gray-600"
        >
          {#if column.icon}
            <svelte:component
              this="{getIconComponent(column.icon)}"
              class="w-6 h-6 text-gray-700 dark:text-white"
            />
          {/if}
        </div>
      </div>
      <button
        on:click="{() => removeColumn(index)}"
        class="w-full p-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:bg-red-300 disabled:cursor-not-allowed dark:bg-red-700 dark:hover:bg-red-800 dark:disabled:bg-red-900"
        disabled="{!canRemoveColumn}"
      >
        Remove Column
      </button>
    </div>
  {/each}

  {#if canAddColumn}
    <div class="p-4 bg-gray-100 dark:bg-gray-800 rounded-lg shadow">
      <h3 class="text-xl font-bold mb-2 dark:text-white">Add New Column</h3>
      <input
        bind:value="{newColumn.name}"
        placeholder="Column name"
        class="w-full p-2 mb-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
      />
      <input
        bind:value="{newColumn.label}"
        placeholder="Column label"
        class="w-full p-2 mb-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
      />
      <select
        bind:value="{newColumn.color}"
        class="w-full p-2 mb-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
      >
        <option value="">Select an optional color</option>
        {#each colorOptions as option}
          <option value="{option.value}">{option.name}</option>
        {/each}
      </select>
      <div class="flex items-center mb-2">
        <select
          bind:value="{newColumn.icon}"
          class="flex-grow p-2 border rounded bg-white dark:bg-gray-700 dark:text-white dark:border-gray-600"
        >
          <option value="">Select an optional icon</option>
          {#each iconOptions as option}
            <option value="{option.value}">{option.name}</option>
          {/each}
        </select>
        <div
          class="ml-2 p-2 bg-white dark:bg-gray-700 border rounded dark:border-gray-600"
        >
          {#if newColumn.icon}
            <svelte:component
              this="{getIconComponent(newColumn.icon)}"
              class="w-6 h-6 text-gray-700 dark:text-white"
            />
          {/if}
        </div>
      </div>
      <button
        on:click="{addColumn}"
        class="w-full p-2 bg-blue-500 text-white rounded hover:bg-blue-600 dark:bg-blue-700 dark:hover:bg-blue-800"
      >
        Add Column
      </button>
    </div>
  {:else}
    <p class="text-sm text-red-600 dark:text-red-400">
      Maximum number of columns reached. You cannot add more columns.
    </p>
  {/if}
</div>
