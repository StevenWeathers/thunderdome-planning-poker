<script lang="ts">
  import { PencilIcon } from '@lucide/svelte';
  import { LL } from '../../i18n/i18n-svelte';
  import type { StoryboardGoal, StoryboardColumn } from '../../types/storyboard';

  interface Props {
    goal: StoryboardGoal;
    goalColumn: StoryboardColumn;
    columnIndex: number;
    columnOrderEditMode: boolean;
    columnWidth: string;
    toggleColumnEdit: (column: StoryboardColumn) => () => void;
    addStory: (goalId: string, columnId: string) => () => void;
  }

  let {
    goal,
    columnOrderEditMode,
    goalColumn,
    columnIndex,
    columnWidth = '10rem',
    toggleColumnEdit,
    addStory,
  }: Props = $props();

  let columnName = $derived(goalColumn.name || `Column ${columnIndex + 1}`);
</script>

<div class="w-full flex flex-col gap-2 self-stretch justify-between" style="width: {columnWidth}">
  <div class="w-full flex gap-2 items-start leading-tight">
    <span
      class="flex-wrap font-bold flex-grow break-words dark:text-gray-300 {goalColumn.name
        ? ''
        : 'italic text-gray-500 dark:text-gray-600'}"
      title={columnName}
      data-testid="column-name"
    >
      {columnName}
    </span>
    <button
      onclick={toggleColumnEdit(goalColumn)}
      class="flex-none py-1 font-bold text-xl text-gray-600 dark:text-gray-400 hover:text-green-500 dark:hover:text-lime-400"
      class:cursor-not-allowed={columnOrderEditMode}
      disabled={columnOrderEditMode}
      title={$LL.storyboardEditColumn()}
      data-testid="column-edit"
    >
      <PencilIcon />
    </button>
  </div>
  <div class="w-full">
    <div class="flex">
      <button
        onclick={addStory(goal.id, goalColumn.id)}
        class="flex-grow font-bold text-xl py-1 px-2 border-dashed border-2
                border-gray-400 dark:border-gray-600 hover:border-green-500 dark:hover:border-lime-400
                text-gray-600 dark:text-gray-400 hover:text-green-500 dark:hover:text-lime-400"
        class:cursor-not-allowed={columnOrderEditMode}
        disabled={columnOrderEditMode}
        title={$LL.storyboardAddStoryToColumn()}
        data-testid="story-add"
      >
        +
      </button>
    </div>
  </div>
</div>
