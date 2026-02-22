<script lang="ts">
  import { dndzone } from 'svelte-dnd-action';
  import StoryCard from './StoryCard.svelte';
  import StoryForm from './StoryForm.svelte';
  import type { StoryboardGoal, StoryboardColumn, StoryboardStory } from '../../types/storyboard';
  import type { ColorLegend } from '../../types/storyboard';
  import type { NotificationService } from '../../types/notifications';

  interface Props {
    goals: StoryboardGoal[];
    goal: StoryboardGoal;
    goalIndex: number;
    goalColumn: StoryboardColumn;
    columnIndex: number;
    columnOrderEditMode: boolean;
    columnWidth: string;
    scale: number;
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
    goalColumn,
    columnIndex,
    columnWidth = '10rem',
    scale,
    sendSocketEvent,
    notifications,
    colorLegend,
    users,
  }: Props = $props();

  let activeStoryId: string | null = $state(null);
  let storyDiscussionExpanded = $state(false);
  let additionalDetailsExpanded = $state(false);

  const toggleStoryForm =
    (
      story: StoryboardStory | null = null,
      options?: { discussionExpanded?: boolean; additionalDetailsExpanded?: boolean },
    ) =>
    () => {
      if (columnOrderEditMode) {
        return;
      }
      activeStoryId = activeStoryId != null ? null : story?.id || null;
      if (options?.discussionExpanded) {
        storyDiscussionExpanded = true;
      } else {
        storyDiscussionExpanded = false;
      }
      if (options?.additionalDetailsExpanded) {
        additionalDetailsExpanded = true;
      } else {
        additionalDetailsExpanded = false;
      }
    };

  let activeStory = $derived(
    activeStoryId ? goalColumn.stories.find(story => story?.id === activeStoryId) || null : null,
  );

  function handleDndConsider(e: CustomEvent) {
    const goalIndex = Number((e.target as HTMLElement)?.dataset.goalindex);
    const columnIndex = Number((e.target as HTMLElement)?.dataset.columnindex);

    if (!isNaN(goalIndex) && !isNaN(columnIndex)) {
      goals[goalIndex].columns[columnIndex].stories = e.detail.items;
      goals = goals;
    }
  }

  function handleDndFinalize(e: CustomEvent) {
    const goalIndex = Number((e.target as HTMLElement)?.dataset.goalindex);
    const columnIndex = Number((e.target as HTMLElement)?.dataset.columnindex);
    const storyId = e.detail.info.id;

    if (!isNaN(goalIndex) && !isNaN(columnIndex)) {
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
  }
</script>

<div class="w-full h-full flex flex-col" style="width: {columnWidth}" data-testid="goal-column">
  <div
    class="flex-1 relative flex flex-col gap-4"
    data-testid="column-dropzone"
    style="min-height: 160px;"
    data-goalid={goal.id}
    data-columnid={goalColumn.id}
    data-goalIndex={goalIndex}
    data-columnindex={columnIndex}
    use:dndzone={{
      items: goalColumn.stories,
      type: 'story' as const,
      dropTargetStyle: {},
      dropTargetClasses: ['outline', 'outline-2', 'outline-indigo-500', 'dark:outline-yellow-400'],
      dragDisabled: columnOrderEditMode,
    }}
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
  >
    {#each goalColumn.stories as story (story.id)}
      <StoryCard {story} {goalColumn} {goal} {columnOrderEditMode} {toggleStoryForm} {scale} />
    {/each}
  </div>
</div>

{#if activeStory}
  <StoryForm
    toggleStoryForm={toggleStoryForm(null)}
    story={activeStory}
    {sendSocketEvent}
    {notifications}
    {colorLegend}
    {users}
    discussionExpanded={storyDiscussionExpanded}
    {additionalDetailsExpanded}
  />
{/if}
