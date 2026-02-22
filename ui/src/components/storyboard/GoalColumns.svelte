<script lang="ts">
  import { dndzone } from 'svelte-dnd-action';
  import { Grip, User } from '@lucide/svelte';
  import StoryColumn from './StoryColumn.svelte';
  import type { StoryboardGoal, StoryboardColumn } from '../../types/storyboard';
  import type { NotificationService } from '../../types/notifications';
  import type { ColorLegend } from '../../types/storyboard';

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
  }: Props = $props();

  // Calculate column width based on scale (w-40 = 10rem for scale 1)
  const columnWidth = $derived(`${10 * scale}rem`);

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

<div class="flex">
  {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
    <div class="flex-none mx-2" style="width: {columnWidth}" data-testid="goal-personas">
      <div class="w-full mb-2">
        {#each goalColumn.personas as persona}
          <div class="mt-4 dark:text-gray-300 text-right" data-testid="goal-persona">
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
<div
  class="flex"
  data-goalid={goal.id}
  data-goalindex={goalIndex}
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
    <div class:relative={columnOrderEditMode} class:cursor-move={columnOrderEditMode}>
      <StoryColumn
        bind:goals
        {goal}
        {goalIndex}
        {columnOrderEditMode}
        {addStory}
        {toggleColumnEdit}
        {sendSocketEvent}
        {goalColumn}
        {columnIndex}
        {scale}
        {notifications}
        {colorLegend}
        {users}
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
