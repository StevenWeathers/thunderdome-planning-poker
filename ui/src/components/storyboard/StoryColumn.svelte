<script lang="ts">
  import { dndzone } from 'svelte-dnd-action';
  import StoryCard from './StoryCard.svelte';
  import StoryCommentPreview from './StoryCommentPreview.svelte';
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
    pointValues?: string[];
    selectionMode?: boolean;
    selectedStoryIds?: string[];
    toggleStorySelection?: (storyId: string) => void;
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
    pointValues = [],
    selectionMode = false,
    selectedStoryIds = [],
    toggleStorySelection = () => {},
  }: Props = $props();

  let activeStoryId: string | null = $state(null);
  let showStoryDiscussion = $state(false);

  const toggleStoryForm =
    (story: StoryboardStory | null = null, options?: { showDiscussion?: boolean }) =>
    () => {
      if (columnOrderEditMode) {
        return;
      }
      activeStoryId = activeStoryId != null ? null : story?.id || null;
      showStoryDiscussion = options?.showDiscussion ?? false;
    };

  let activeStory = $derived(
    activeStoryId ? goalColumn.stories.find(story => story?.id === activeStoryId) || null : null,
  );

  let commentPreview: { storyId: string; anchor: HTMLElement } | null = $state(null);

  let previewStory = $derived(
    goalColumn.stories.find(story => story.id === commentPreview?.storyId && story.comments.length > 0),
  );

  // forget the preview once its story is moved, deleted or left without comments
  $effect(() => {
    if (commentPreview && !previewStory) {
      commentPreview = null;
    }
  });

  function toggleCommentPreview(story: StoryboardStory, anchor: HTMLElement) {
    if (columnOrderEditMode) {
      return;
    }
    commentPreview = commentPreview?.storyId === story.id ? null : { storyId: story.id, anchor };
  }

  function closeCommentPreview() {
    commentPreview = null;
  }

  function openPreviewDiscussion() {
    const story = previewStory;
    closeCommentPreview();
    if (story) {
      toggleStoryForm(story, { showDiscussion: true })();
    }
  }

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
      dragDisabled: columnOrderEditMode || selectionMode,
    }}
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
  >
    {#each goalColumn.stories as story (story.id)}
      <StoryCard
        {story}
        {goalColumn}
        {goal}
        {columnOrderEditMode}
        {toggleStoryForm}
        {scale}
        {selectionMode}
        selected={selectedStoryIds.includes(story.id)}
        {toggleStorySelection}
        {toggleCommentPreview}
        commentPreviewOpen={commentPreview?.storyId === story.id}
      />
    {/each}
  </div>
</div>

<!-- outside the drop zone so clicks in the preview don't reach the card or start a drag -->
{#if commentPreview && previewStory}
  {#key commentPreview.anchor}
    <StoryCommentPreview
      story={previewStory}
      {users}
      anchor={commentPreview.anchor}
      onClose={closeCommentPreview}
      onOpenDiscussion={openPreviewDiscussion}
    />
  {/key}
{/if}

{#if activeStory}
  <StoryForm
    toggleStoryForm={toggleStoryForm(null)}
    story={activeStory}
    {sendSocketEvent}
    {notifications}
    {colorLegend}
    {users}
    {pointValues}
    showDiscussion={showStoryDiscussion}
  />
{/if}
