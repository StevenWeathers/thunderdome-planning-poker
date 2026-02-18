<script lang="ts">
  import { dndzone } from 'svelte-dnd-action';
  import StoryCard from './StoryCard.svelte';
  import type { StoryboardGoal, StoryboardColumn, StoryboardStory } from '../../types/storyboard';
  import { Pencil } from '@lucide/svelte';
  import { LL } from '../../i18n/i18n-svelte';

  interface Props {
    goals: StoryboardGoal[];
    goal: StoryboardGoal;
    goalIndex: number;
    goalColumn: StoryboardColumn;
    columnIndex: number;
    columnOrderEditMode: boolean;
    toggleColumnEdit: (column: StoryboardColumn) => () => void;
    addStory: (goalId: string, columnId: string) => () => void;
    toggleStoryForm: (story?: StoryboardStory) => () => void;
    sendSocketEvent: (event: string, data: string) => void;
  }

  let {
    goals = $bindable([]),
    goal,
    goalIndex,
    columnOrderEditMode,
    goalColumn,
    columnIndex,
    toggleColumnEdit,
    addStory,
    toggleStoryForm,
    sendSocketEvent,
  }: Props = $props();

  function handleDndConsider(e) {
    const goalIndex = e.target.dataset.goalindex;
    const columnIndex = e.target.dataset.columnindex;

    goals[goalIndex].columns[columnIndex].stories = e.detail.items;
    goals = goals;
  }

  function handleDndFinalize(e) {
    const goalIndex = e.target.dataset.goalindex;
    const columnIndex = e.target.dataset.columnindex;
    const storyId = e.detail.info.id;

    goals[goalIndex].columns[columnIndex].stories = e.detail.items;
    goals = goals;

    const matchedStory = goals[goalIndex].columns[columnIndex].stories.find((i: StoryboardStory) => i.id === storyId);

    if (matchedStory) {
      const goalId = goals[goalIndex].id;
      const columnId = goals[goalIndex].columns[columnIndex].id;

      // determine what story to place story before in target column
      const matchedStoryIndex = goals[goalIndex].columns[columnIndex].stories.indexOf(matchedStory);
      const sibling = goals[goalIndex].columns[columnIndex].stories[matchedStoryIndex + 1];
      const placeBefore = sibling ? sibling.id : '';

      sendSocketEvent(
        'move_story',
        JSON.stringify({
          storyId,
          goalId,
          columnId,
          placeBefore,
        }),
      );
    }
  }
</script>

<div class="flex-none my-4 mx-2 w-40" data-testid="goal-column">
  <div class="flex-none">
    <div class="w-full mb-2">
      <div class="flex">
        <span class="font-bold flex-grow truncate dark:text-gray-300" title={goalColumn.name} data-testid="column-name">
          {goalColumn.name}
        </span>
        <button
          onclick={toggleColumnEdit(goalColumn)}
          class="flex-none font-bold text-xl
                                border-dashed border-2 border-gray-400 dark:border-gray-600
                                hover:border-green-500 text-gray-600 dark:text-gray-400
                                hover:text-green-500 py-1 px-2"
          class:cursor-not-allowed={columnOrderEditMode}
          disabled={columnOrderEditMode}
          title={$LL.storyboardEditColumn()}
          data-testid="column-edit"
        >
          <Pencil />
        </button>
      </div>
    </div>
    <div class="w-full">
      <div class="flex">
        <button
          onclick={addStory(goal.id, goalColumn.id)}
          class="flex-grow font-bold text-xl py-1
                                px-2 border-dashed border-2
                                border-gray-400 dark:border-gray-600 hover:border-green-500
                                text-gray-600 dark:text-gray-400 hover:text-green-500"
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
  <div
    class="w-full relative"
    data-testid="column-dropzone"
    style="min-height: 160px;"
    data-goalid={goal.id}
    data-columnid={goalColumn.id}
    data-goalIndex={goalIndex}
    data-columnindex={columnIndex}
    use:dndzone={{
      items: goalColumn.stories,
      type: 'story',
      dropTargetStyle: '',
      dropTargetClasses: ['outline', 'outline-2', 'outline-indigo-500', 'dark:outline-yellow-400'],
      dragDisabled: columnOrderEditMode,
    }}
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
  >
    {#each goalColumn.stories as story (story.id)}
      <StoryCard {story} {goalColumn} {goal} {columnOrderEditMode} {toggleStoryForm} />
    {/each}
  </div>
</div>
