<script lang="ts">
  import { PencilIcon, SettingsIcon, TrashIcon, UserIcon } from '@lucide/svelte';
  import { LL } from '../../i18n/i18n-svelte';
  import type { StoryboardGoal, StoryboardColumn } from '../../types/storyboard';
  import ActionsMenu from '../global/ActionsMenu.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import ColumnPersonasModal from './ColumnPersonasModal.svelte';

  interface Props {
    goal: StoryboardGoal;
    goalColumn: StoryboardColumn;
    columnIndex: number;
    columnOrderEditMode: boolean;
    columnWidth: string;
    toggleColumnEdit: (column: StoryboardColumn) => () => void;
    sendSocketEvent: (event: string, data: any) => void;
    addStory: (goalId: string, columnId: string) => () => void;
    personas: any[];
  }

  let {
    goal,
    columnOrderEditMode,
    goalColumn,
    columnIndex,
    columnWidth = '10rem',
    personas = [],
    toggleColumnEdit,
    sendSocketEvent,
    addStory,
  }: Props = $props();

  let showColumnDeleteConfirmation = $state(false);
  let showPersonaManagement = $state(false);
  let columnName = $derived(goalColumn.name || `Column ${columnIndex + 1}`);

  function toggleColumnDeleteConfirmation() {
    showColumnDeleteConfirmation = !showColumnDeleteConfirmation;
  }

  function handleDeleteColumn() {
    sendSocketEvent('delete_column', goalColumn.id);
  }

  function togglePersonaManagement() {
    showPersonaManagement = !showPersonaManagement;
  }

  function handlePersonaAdd(data: { column_id: string; persona_id: string }) {
    sendSocketEvent('column_persona_add', JSON.stringify(data));
  }

  function handlePersonaRemove(data: { column_id: string; persona_id: string }) {
    sendSocketEvent('column_persona_remove', JSON.stringify(data));
  }
</script>

<div class="w-full flex flex-col gap-2 self-stretch justify-between" style="width: {columnWidth}">
  <div class="w-full flex gap-2 items-start leading-tight">
    <span
      class="font-bold flex-grow min-w-0 break-words dark:text-gray-300 {goalColumn.name
        ? ''
        : 'italic text-gray-500 dark:text-gray-600'}"
      title={columnName}
      data-testid="column-name"
    >
      {columnName}
    </span>
    <ActionsMenu
      actions={[
        {
          label: $LL.edit(),
          icon: PencilIcon,
          onclick: toggleColumnEdit(goalColumn),
          testId: 'edit-column-action',
        },
        {
          label: $LL.personas(),
          icon: UserIcon,
          onclick: togglePersonaManagement,
          testId: 'personas-column-action',
        },
        {
          label: $LL.delete(),
          icon: TrashIcon,
          onclick: toggleColumnDeleteConfirmation,
          testId: 'delete-column-action',
        },
      ]}
      ariaLabel="Column actions"
      testId="column-actions-menu"
      Icon={SettingsIcon}
      iconSize="medium"
      disabled={columnOrderEditMode}
    />
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

  {#if showColumnDeleteConfirmation}
    <DeleteConfirmation
      toggleDelete={toggleColumnDeleteConfirmation}
      handleDelete={handleDeleteColumn}
      confirmText={'Are you sure you want to delete this Column?'}
      confirmBtnText={'Delete Column'}
    />
  {/if}

  {#if showPersonaManagement}
    <ColumnPersonasModal
      toggleModal={togglePersonaManagement}
      column={goalColumn}
      {personas}
      onPersonaAdd={handlePersonaAdd}
      onPersonaRemove={handlePersonaRemove}
    />
  {/if}
</div>
