<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import StoryColorSelector from './StoryColorSelector.svelte';
  import LL from '../../i18n/i18n-svelte';
  import type { ColorLegend, StoryboardColumn } from '../../types/storyboard';
  import { onMount } from 'svelte';

  interface Props {
    toggleColumnEdit?: any;
    handleColumnRevision?: any;
    column?: StoryboardColumn | any;
    goalId?: string;
    colorLegend?: ColorLegend[];
  }

  let {
    toggleColumnEdit = () => {},
    handleColumnRevision = () => {},
    goalId = '',
    colorLegend = [],
    column = $bindable({
      id: '',
      name: '',
      default_story_color: null,
      personas: [],
      stories: [],
      sort_order: '',
    }),
  }: Props = $props();

  let focusInput: any = $state();
  let columnName = $state(column.name);
  let selectedDefaultStoryColor = $state(column.default_story_color ?? '');

  $effect(() => {
    columnName = column.name;
    selectedDefaultStoryColor = column.default_story_color ?? '';
  });

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const defaultStoryColor = selectedDefaultStoryColor === '' ? null : selectedDefaultStoryColor;

    const c = {
      id: column.id,
      name: columnName,
      defaultStoryColor,
    };

    handleColumnRevision(c);
    toggleColumnEdit();
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleColumnEdit} ariaLabel={$LL.modalStoryboardColumnSettings()}>
  <form onsubmit={handleSubmit} name="addColumn">
    <div class="mb-4">
      <label class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2" for="columnName">
        Column Name
      </label>
      <TextInput
        id="columnName"
        bind:value={columnName}
        placeholder="Enter a column name"
        name="columnName"
        bind:this={focusInput}
      />
    </div>
    <div class="mb-4">
      <label class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2" for="columnDefaultStoryColor">
        Default Story Color
      </label>
      <div id="columnDefaultStoryColor">
        <StoryColorSelector
          bind:value={selectedDefaultStoryColor}
          {colorLegend}
          allowNone={true}
          noneLabel="No default"
        />
      </div>
    </div>
    <div class="flex justify-end gap-2">
      <SolidButton type="submit">Save</SolidButton>
    </div>
  </form>
</Modal>
