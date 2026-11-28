<script lang="ts">
  import { dndzone } from 'svelte-dnd-action';
  import { Grip, User } from '@lucide/svelte';
  import StoryColumn from './StoryColumn.svelte';
  import type { StoryboardGoal, StoryboardColumn } from '../../types/storyboard';
  import type { NotificationService } from '../../types/notifications';
  import type { ColorLegend } from '../../types/storyboard';
  import ColumnHeader from './ColumnHeader.svelte';
  import { LL } from '../../i18n/i18n-svelte';

  interface Props {
    goals: StoryboardGoal[];
    goal: StoryboardGoal;
    goalIndex: number;
    columnOrderEditMode: boolean;
    scale: number;
    toggleColumnEdit: (column: StoryboardColumn) => () => void;
    addStory: (goalId: string, columnId: string) => () => void;
    sendSocketEvent: (event: string, data: string) => void;
    notifications: NotificationService;
    colorLegend: ColorLegend[];
    users: any[];
    personas: any[];
  }

  let {
    goals = $bindable([]),
    goal,
    goalIndex,
    columnOrderEditMode,
    scale,
    toggleColumnEdit,
    addStory,
    sendSocketEvent,
    notifications,
    colorLegend,
    users,
    personas,
  }: Props = $props();

  // Calculate column width based on scale (w-40 = 10rem for scale 1)
  const columnWidth = $derived(`${10 * scale}rem`);

  // bulk story delete: one column per goal can be in selection mode at a time
  let selectionColumnId: string | null = $state(null);
  let checkedStoryIds: string[] = $state([]);

  // drop ids of stories that other users have since deleted or moved away
  const selectedStoryIds = $derived(
    goal.columns
      .find(column => column.id === selectionColumnId)
      ?.stories.map(story => story.id)
      .filter(id => checkedStoryIds.includes(id)) ?? [],
  );

  function startSelection(columnId: string) {
    selectionColumnId = columnId;
    checkedStoryIds = [];
  }

  function cancelSelection() {
    selectionColumnId = null;
    checkedStoryIds = [];
  }

  function toggleStorySelection(storyId: string) {
    checkedStoryIds = checkedStoryIds.includes(storyId)
      ? checkedStoryIds.filter(id => id !== storyId)
      : [...checkedStoryIds, storyId];
  }

  $effect(() => {
    if (columnOrderEditMode) {
      cancelSelection();
    }
  });

  function handleDndConsider(e: CustomEvent) {
    const goalIndex = Number((e.target as HTMLElement)?.dataset.goalindex);

    if (!isNaN(goalIndex)) {
      goals[goalIndex].columns = e.detail.items;
      goals = goals;
    }
  }

  function handleDndFinalize(e: CustomEvent) {
    const goalIndex = Number((e.target as HTMLElement)?.dataset.goalindex);
    const columnId = e.detail.info.id;

    if (!isNaN(goalIndex)) {
      goals[goalIndex].columns = e.detail.items;
      goals = goals;

      const matchedColumn = goals[goalIndex].columns.find(column => column.id === columnId);

      if (matchedColumn) {
        const goalId = goals[goalIndex].id;

        // determine what column to place column before in target goal
        const matchedColumnIndex = goals[goalIndex].columns.indexOf(matchedColumn);
        const sibling = goals[goalIndex].columns[matchedColumnIndex + 1];
        const placeBefore = sibling ? sibling.id : '';

        sendSocketEvent(
          'move_column',
          JSON.stringify({
            goalId,
            columnId,
            placeBefore,
          }),
        );
      }
    }
  }
</script>

<div class="flex flex-col gap-4">
  <!-- Column Associated Personas -->
  <div class="flex gap-4 items-stretch">
    {#each goal.columns as goalColumn, _columnIndex (goalColumn.id)}
      <div class="flex flex-shrink-0" style="width: {columnWidth}" data-testid="goal-personas">
        <div class="w-full flex flex-col justify-end gap-2">
          {#each goalColumn.personas as persona}
            <div class="dark:text-gray-300 text-right" data-testid="goal-persona">
              <div class="font-bold" data-testid="persona-name">
                <User class="inline-block h-4 w-4" />
                {persona.name}
              </div>
              <div class="text-sm" data-testid="persona-role">
                {persona.role}
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/each}
  </div>
  <!-- Column Headers -->
  <div class="flex gap-4 items-stretch">
    {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
      <div class="flex" style="width: {columnWidth}">
        <ColumnHeader
          {goal}
          {goalColumn}
          {columnIndex}
          {columnOrderEditMode}
          {toggleColumnEdit}
          {addStory}
          {columnWidth}
          {sendSocketEvent}
          {personas}
          selectionMode={selectionColumnId === goalColumn.id}
          {selectedStoryIds}
          startSelection={() => startSelection(goalColumn.id)}
          selectAllStories={() => (checkedStoryIds = goalColumn.stories.map(story => story.id))}
          {cancelSelection}
        />
      </div>
    {/each}
  </div>
  <div class="relative">
    <!-- a sibling of the drop zone, as every child of the zone is treated as a draggable item -->
    {#if columnOrderEditMode && goal.columns.length === 0}
      <div
        class="absolute inset-0 flex items-center justify-center rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 pointer-events-none"
        data-testid="column-drop-hint"
      >
        {$LL.storyboardDropColumnHere()}
      </div>
    {/if}
    <div
      class="flex gap-4 items-stretch {columnOrderEditMode ? 'min-h-40' : ''}"
      data-goalid={goal.id}
      data-goalindex={goalIndex}
      data-testid="goal-columns-dropzone"
      use:dndzone={{
        items: goal.columns,
        type: 'column' as const,
        dropTargetStyle: {},
        dropTargetClasses: ['outline', 'outline-2', 'outline-indigo-500', 'dark:outline-yellow-400'],
        dragDisabled: !columnOrderEditMode,
      }}
      onconsider={handleDndConsider}
      onfinalize={handleDndFinalize}
    >
      {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
        <div class="self-stretch" class:relative={columnOrderEditMode} class:cursor-move={columnOrderEditMode}>
          <StoryColumn
            bind:goals
            {goal}
            {goalIndex}
            {columnOrderEditMode}
            {sendSocketEvent}
            {goalColumn}
            {columnIndex}
            {scale}
            {columnWidth}
            {notifications}
            {colorLegend}
            {users}
            selectionMode={selectionColumnId === goalColumn.id}
            {selectedStoryIds}
            {toggleStorySelection}
          />

          {#if columnOrderEditMode}
            <div class="absolute top-0 right-0 bottom-0 left-0 z-5 bg-gray-900 opacity-50" aria-hidden="true"></div>
            <div class="absolute top-0 right-0 bottom-0 left-0 z-10 flex flex-col items-center justify-center">
              <Grip class="h-8 w-8 text-white mb-2" />
              <div class="text-white font-bold">Drag to Reorder</div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  </div>
</div>
