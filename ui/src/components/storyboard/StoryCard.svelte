<script lang="ts">
  import { Link, MessageSquareMore } from '@lucide/svelte';
  import { SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import type { StoryboardStory, StoryboardColumn, StoryboardGoal } from '../../types/storyboard';

  interface Props {
    story: StoryboardStory;
    goalColumn: StoryboardColumn;
    goal: StoryboardGoal;
    columnOrderEditMode: boolean;
    scale: number;
    toggleStoryForm: (
      story: StoryboardStory,
      options?: { discussionExpanded?: boolean; additionalDetailsExpanded?: boolean },
    ) => () => void;
  }

  let {
    story = {
      id: '',
      name: '',
      color: '',
      points: 0,
      link: '',
      comments: [],
      closed: false,
      sort_order: 'a1',
      content: '',
      annotations: [],
    },
    goalColumn = { id: '', name: '', sort_order: '', stories: [], personas: [] },
    goal = { id: '', name: '', sort_order: '', columns: [], personas: [] },
    columnOrderEditMode = false,
    scale,
    toggleStoryForm = (story: StoryboardStory) => () => {},
  }: Props = $props();

  // Calculate content height based on scale (h-20 = 5rem for scale 1)
  const contentHeight = $derived(`${5 * scale}rem`);

  // Handle click on external link to prevent triggering story form
  function handleExtLinkClick(e: Event) {
    e.stopPropagation();
  }

  // Handle click on comments button to prevent triggering story form without expanding comments
  function handleCommentsClick(e: Event) {
    e.stopPropagation();
    toggleStoryForm(story, { discussionExpanded: true })();
  }

  // Handle click on additional details to prevent triggering story form without expanding additional details
  function handleAdditionalDetailsClick(e: Event) {
    e.stopPropagation();
    toggleStoryForm(story, { additionalDetailsExpanded: true })();
  }
</script>

<div
  class="relative shadow bg-white dark:bg-gray-700 dark:text-white border-s-4 story-{story.color} border my-4"
  style="list-style: none;"
  class:cursor-pointer={!columnOrderEditMode}
  class:cursor-not-allowed={columnOrderEditMode}
  role="button"
  tabindex="0"
  data-goalid={goal.id}
  data-columnid={goalColumn.id}
  data-storyid={story.id}
  data-testid="column-story"
  onclick={toggleStoryForm(story)}
  onkeypress={toggleStoryForm(story)}
>
  <div>
    <div>
      <div
        class="p-2 text-sm overflow-hidden {story.closed ? 'line-through' : ''}"
        style="height: {contentHeight}"
        title={story.name}
        data-testid="story-name"
      >
        {story.name}
      </div>
      <div class="h-10">
        <div class="flex content-center p-2 text-sm">
          <div class="w-1/2 text-gray-600 dark:text-gray-300">
            {#if story.name.length > 0}
              <button
                class="inline-block align-middle hover:text-blue-600 dark:hover:text-cyan-400 transition-all duration-200"
                data-testid="story-comments"
                title="Story has {story.comments.length} {story.comments.length ? 'comments' : 'comment'}"
                onclick={handleCommentsClick}
              >
                {story.comments.length}
                <MessageSquareMore class="inline-block" />
              </button>
            {/if}
          </div>
          <div class="w-1/2 flex space-x-2 justify-end">
            {#if story.link !== ''}
              <a
                href={story.link}
                onclick={handleExtLinkClick}
                target="_blank"
                rel="noopener noreferrer"
                title="Story has external link"
                class="hover:text-blue-600 dark:hover:text-cyan-400 transition-all duration-200"
                ><Link class="inline-block w-4 h-4" /></a
              >
            {/if}
            {#if story.points > 0}
              <button
                class="px-2 bg-gray-300 dark:bg-gray-500 inline-block align-middle rounded-full hover:bg-green-400 dark:hover:bg-lime-400 dark:hover:text-gray-800 transition-all duration-200"
                data-testid="story-points"
                title="Story points"
                onclick={handleAdditionalDetailsClick}
              >
                {story.points}
              </button>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
  {#if (story as any)[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
    <div
      class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible opacity-50 max-w-xs shadow bg-white dark:bg-gray-700 dark:text-white border-s-4
                                story-{story.color} border
                                cursor-pointer"
      style="list-style: none;"
      tabindex="-1"
      data-goalid={goal.id}
      data-columnid={goalColumn.id}
      data-storyid={story.id}
      data-testid="column-story-shadowitem"
    >
      <div>
        <div>
          <div
            class="p-2 text-sm overflow-hidden {story.closed ? 'line-through' : ''}"
            style="height: {contentHeight}"
            title={story.name}
            data-testid="shadow-story-name"
          >
            {story.name}
          </div>
          <div class="h-10">
            <div class="flex content-center p-2 text-sm">
              <div class="w-1/2 text-gray-600 dark:text-gray-300">
                {#if story.comments.length > 0}
                  <span
                    class="inline-block align-middle"
                    data-testid="story-comments"
                    title="Story has {story.comments.length} {story.comments.length ? 'comments' : 'comment'}"
                  >
                    {story.comments.length}
                    <MessageSquareMore class="inline-block" />
                  </span>
                {/if}
              </div>
              <div class="w-1/2 flex space-x-2 justify-end">
                {#if story.link !== ''}
                  <span title="Story has external link"><Link class="inline-block w-4 h-4" /></span>
                {/if}
                {#if story.points > 0}
                  <span
                    class="px-2 bg-gray-300 dark:bg-gray-500 inline-block align-middle rounded-full"
                    data-testid="story-points"
                    title="Story points"
                  >
                    {story.points}
                  </span>
                {/if}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style lang="postcss">
  .story-gray {
    @apply border-gray-400;
  }

  .story-gray:hover {
    @apply border-gray-800;
  }

  .story-red {
    @apply border-red-400;
  }

  .story-red:hover {
    @apply border-red-800;
  }

  .story-orange {
    @apply border-orange-400;
  }

  .story-orange:hover {
    @apply border-orange-800;
  }

  .story-yellow {
    @apply border-yellow-400;
  }

  .story-yellow:hover {
    @apply border-yellow-800;
  }

  .story-green {
    @apply border-green-400;
  }

  .story-green:hover {
    @apply border-green-800;
  }

  .story-teal {
    @apply border-teal-400;
  }

  .story-teal:hover {
    @apply border-teal-800;
  }

  .story-blue {
    @apply border-blue-400;
  }

  .story-blue:hover {
    @apply border-blue-800;
  }

  .story-indigo {
    @apply border-indigo-400;
  }

  .story-indigo:hover {
    @apply border-indigo-800;
  }

  .story-purple {
    @apply border-purple-400;
  }

  .story-purple:hover {
    @apply border-purple-800;
  }

  .story-pink {
    @apply border-pink-400;
  }

  .story-pink:hover {
    @apply border-pink-800;
  }

  .colorcard-gray {
    @apply bg-gray-400;
  }

  .colorcard-red {
    @apply bg-red-400;
  }

  .colorcard-orange {
    @apply bg-orange-400;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
  }

  .colorcard-green {
    @apply bg-green-400;
  }

  .colorcard-teal {
    @apply bg-teal-400;
  }

  .colorcard-blue {
    @apply bg-blue-400;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
  }

  .colorcard-purple {
    @apply bg-purple-400;
  }

  .colorcard-pink {
    @apply bg-pink-400;
  }
</style>
