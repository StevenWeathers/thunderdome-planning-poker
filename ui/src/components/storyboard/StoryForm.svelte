<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import Editor from '../forms/Editor.svelte';
  import { ChevronRight, ChevronDown } from '@lucide/svelte';
  import { onMount } from 'svelte';
  import Comment from '../comments/Comment.svelte';
  import CommentForm from '../comments/CommentForm.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { UserDisplay } from '../../types/user';
  import type { StoryboardStory, StoryboardUser } from '../../types/storyboard';
  import CommentsHeader from '../comments/CommentsHeader.svelte';
  import CommentEmptyState from '../comments/CommentEmptyState.svelte';

  interface Props {
    toggleStoryForm?: any;
    sendSocketEvent?: any;
    notifications: NotificationService;
    story?: StoryboardStory;
    colorLegend?: any;
    users?: StoryboardUser[];
    discussionExpanded?: boolean;
    additionalDetailsExpanded?: boolean;
  }

  let {
    toggleStoryForm = () => {},
    sendSocketEvent = () => {},
    notifications,
    story = { id: '' } as StoryboardStory,
    colorLegend = [],
    users = [],
    discussionExpanded = false,
    additionalDetailsExpanded = false,
  }: Props = $props();

  const isAbsolute = new RegExp('^([a-z]+://|//)', 'i');

  let additionalDetailsHidden = $state(true);
  let discussionHidden = $state(true);
  let focusInput: any = $state();
  let storyPoints = $state(0);

  $effect(() => {
    discussionHidden = !discussionExpanded;
    additionalDetailsHidden = !additionalDetailsExpanded;
  });

  const userMap: Map<string, UserDisplay> = $derived(
    users.reduce((prev, usr) => {
      prev.set(usr.id, {
        id: usr.id,
        name: usr.name,
        avatar: usr.avatar,
        gravatarHash: usr.gravatarHash,
        pictureUrl: '',
      });
      return prev;
    }, new Map<string, UserDisplay>()),
  );

  function handleStoryDelete() {
    sendSocketEvent('delete_story', story.id);
    toggleStoryForm();
  }

  function markClosed() {
    sendSocketEvent(
      'update_story_closed',
      JSON.stringify({
        storyId: story.id,
        closed: true,
      }),
    );
  }

  function markOpen() {
    sendSocketEvent(
      'update_story_closed',
      JSON.stringify({
        storyId: story.id,
        closed: false,
      }),
    );
  }

  const changeColor = (color: string) => () => {
    sendSocketEvent(
      'update_story_color',
      JSON.stringify({
        storyId: story.id,
        color,
      }),
    );
  };

  const updateName = (evt: Event) => {
    const name = (evt.target as HTMLInputElement).value;
    sendSocketEvent(
      'update_story_name',
      JSON.stringify({
        storyId: story.id,
        name,
      }),
    );
  };

  const updateContent = (content?: string) => {
    const contentToSend = content ?? story.content;
    sendSocketEvent(
      'update_story_content',
      JSON.stringify({
        storyId: story.id,
        content: contentToSend,
      }),
    );
  };

  const updatePoints = (evt: Event) => {
    const points = parseInt((evt.target as HTMLInputElement).value, 10);
    sendSocketEvent(
      'update_story_points',
      JSON.stringify({
        storyId: story.id,
        points,
      }),
    );
  };

  const updateLink = (evt: Event) => {
    const link = (evt.target as HTMLInputElement).value;
    if (link !== '' && !isAbsolute.test(link)) {
      notifications.danger('Link must be an absolute URL');
      return;
    }

    sendSocketEvent(
      'update_story_link',
      JSON.stringify({
        storyId: story.id,
        link,
      }),
    );
  };

  const handleCommentSubmit = (commentText: string) => {
    if (commentText.trim() !== '') {
      sendSocketEvent('add_story_comment', JSON.stringify({ storyId: story.id, comment: commentText }));
    }
  };

  const handleCommentEdit = (commentId: string, data: { userId: string; comment: string }) => {
    sendSocketEvent(
      'edit_story_comment',
      JSON.stringify({
        commentId,
        comment: data.comment,
      }),
    );
  };

  const handleCommentDelete = (commentId: String) => {
    sendSocketEvent('delete_story_comment', JSON.stringify({ commentId }));
  };

  const toggleMoreActions = () => {
    additionalDetailsHidden = !additionalDetailsHidden;
  };

  const toggleDiscussion = () => {
    discussionHidden = !discussionHidden;
  };

  onMount(() => {
    focusInput?.focus();
  });

  $effect(() => {
    if (story?.points !== undefined) {
      storyPoints = story.points;
    }
  });
</script>

<Modal
  closeModal={toggleStoryForm}
  widthClasses="w-full md:w-3/4 xl:2/3 2xl:w-3/5"
  ariaLabel={$LL.modalStoryboardStory()}
>
  <div class="p-6 overflow-y-auto h-full">
    <div class="space-y-6">
      <!-- Story Name (Full Width) -->
      <div>
        <label for="storyName" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
          Story Name <span class="text-red-500">*</span>
        </label>
        <!-- <input 
                  type="text" 
                  id="story-name"
                  placeholder="e.g. User Authentication System"
                  class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                  value=""
              > -->
        <TextInput
          id="storyName"
          onchange={updateName}
          value={story.name}
          placeholder="e.g. User Authentication System"
          name="storyName"
          bind:this={focusInput}
        />
      </div>

      <!-- Story Points and Color -->
      <div>
        <div class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">Story Color</div>
        <div class="flex space-x-2 pt-1">
          {#each colorLegend as color}
            <button
              onclick={changeColor(color.color)}
              class="w-8 h-8 rounded-full colorcard-{color.color}
                  hover:scale-110 transition-transform dark:ring-offset-gray-800 {story.color === color.color
                ? `ring-2 ring-offset-2`
                : ''}"
              title="{color.color}{color.legend !== '' ? ` - ${color.legend}` : ''}"
              ><span class="hidden">change color</span></button
            >
          {/each}
        </div>
      </div>

      <!-- Story Content -->
      <div>
        <label for="story-content" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
          Story Description
        </label>
        <div class="border border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden">
          <!-- <textarea 
                      id="story-content"
                      rows="8"
                      placeholder="As a [user type], I want [functionality] so that [benefit]..."
                      class="w-full px-4 py-3 border-0 dark:bg-gray-700 dark:text-white focus:ring-0 resize-none"
                  ></textarea> -->
          <div class="bg-white">
            <Editor
              content={story.content}
              placeholder="Enter story content"
              id="storyDescription"
              handleTextChange={(c: string) => {
                updateContent(c);
              }}
            />
          </div>
        </div>
      </div>

      <!-- More Actions Section (Collapsed by default) -->
      <div
        class="flex flex-col gap-3 {additionalDetailsHidden
          ? 'border-t border-gray-200 dark:border-gray-700 pt-4 lg:pt-6'
          : 'bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4 border border-gray-200/50 dark:border-gray-600/30'}"
      >
        <button
          type="button"
          id="more-actions-toggle"
          class="group flex items-center gap-2 w-full text-start rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
          onclick={toggleMoreActions}
          aria-expanded={!additionalDetailsHidden}
        >
          <span
            class="flex items-center justify-center w-8 h-8 rounded-lg text-gray-400 dark:text-gray-500 group-hover:text-gray-600 dark:group-hover:text-gray-300 group-hover:bg-gray-100 dark:group-hover:bg-gray-700 transition-colors duration-200"
            aria-hidden="true"
          >
            {#if additionalDetailsHidden}
              <ChevronRight class="w-4 h-4" />
            {:else}
              <ChevronDown class="w-4 h-4" />
            {/if}
          </span>
          <span class="font-medium text-gray-900 dark:text-white">Additional Details</span>
        </button>

        <div id="more-actions-content" class="flex flex-col gap-2 {additionalDetailsHidden ? 'hidden' : ''}">
          <!-- Story Points -->
          <div>
            <label for="storyPoints" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">Story Points</label>
            <input
              class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
              id="storyPoints"
              bind:value={storyPoints}
              onchange={updatePoints}
              placeholder="Enter story points"
              name="storyPoints"
            />
          </div>

          <!-- Story Link -->
          <div>
            <label for="storyLink" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
              Story Link
              <span class="text-gray-500 dark:text-gray-400 font-normal ms-1 text-lg">(Optional)</span>
            </label>
            <div class="relative">
              <TextInput
                id="storyLink"
                onchange={updateLink}
                value={story.link}
                placeholder="https://jira.company.com/browse/PROJ-123"
                name="storyLink"
                type="url"
                class="ps-10"
              />
              <svg
                class="w-5 h-5 text-gray-400 dark:text-gray-500 absolute left-3 top-1/2 transform -translate-y-1/2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                ></path>
              </svg>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Story Status-->
            <div>
              <div class="text-gray-700 dark:text-gray-300 mb-3 text-lg">Story Status</div>
              {#if !story.closed}
                <HollowButton color="orange" onClick={markClosed}>
                  <svg class="w-4 h-4 inline me-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                  </svg> Mark as Closed
                </HollowButton>
              {:else}
                <HollowButton color="green" onClick={markOpen}>Reopen story</HollowButton>
              {/if}
            </div>

            <!-- Delete Story-->
            <div>
              <h4 class="text-red-800 dark:text-red-400 mb-3 flex items-center text-lg">
                <svg class="w-5 h-5 me-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.966-.833-2.736 0L3.478 16.5c-.77.833.192 2.5 1.732 2.5z"
                  ></path>
                </svg>
                Danger Zone
              </h4>
              <div class="flex space-x-3">
                <HollowButton color="red" onClick={handleStoryDelete}>
                  <svg class="w-4 h-4 inline me-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1-1H8a1 1 0 00-1 1v3M4 7h16"
                    ></path>
                  </svg> Delete Story
                </HollowButton>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Discussion Section -->
      <div
        class="flex flex-col gap-3 {discussionHidden
          ? 'border-t border-gray-200 dark:border-gray-700 pt-4 lg:pt-6'
          : 'bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4 border border-gray-200/50 dark:border-gray-600/30'}"
      >
        <CommentsHeader
          commentsCount={story.comments ? story.comments.length : 0}
          onToggleExpand={toggleDiscussion}
          isExpanded={!discussionHidden}
        />

        <div class="flex flex-col gap-3 {discussionHidden ? 'hidden' : ''}">
          {#if story.comments && story.comments.length > 0}
            <div class="flex flex-col gap-3">
              {#each story.comments as comment}
                <Comment {comment} {userMap} handleEdit={handleCommentEdit} handleDelete={handleCommentDelete} />
              {/each}
            </div>
          {:else}
            <CommentEmptyState description="Be the first to share your thoughts on this story." />
          {/if}

          <CommentForm onSubmit={handleCommentSubmit} />
        </div>
      </div>
    </div>
  </div>
</Modal>

<style lang="postcss">
  .colorcard-gray {
    @apply bg-gray-400;
    @apply ring-gray-400;
  }

  .colorcard-gray:hover {
    @apply bg-gray-600;
    @apply ring-gray-600;
  }

  .colorcard-red {
    @apply bg-red-400;
    @apply ring-red-400;
  }

  .colorcard-red:hover {
    @apply bg-red-600;
    @apply ring-red-600;
  }

  .colorcard-orange {
    @apply bg-orange-400;
    @apply ring-orange-400;
  }

  .colorcard-orange:hover {
    @apply bg-orange-600;
    @apply ring-orange-600;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
    @apply ring-yellow-400;
  }

  .colorcard-yellow:hover {
    @apply bg-yellow-600;
    @apply ring-yellow-600;
  }

  .colorcard-green {
    @apply bg-green-400;
    @apply ring-green-400;
  }

  .colorcard-green:hover {
    @apply bg-green-600;
    @apply ring-green-600;
  }

  .colorcard-teal {
    @apply bg-teal-400;
    @apply ring-teal-400;
  }

  .colorcard-teal:hover {
    @apply bg-teal-600;
    @apply ring-teal-600;
  }

  .colorcard-blue {
    @apply bg-blue-400;
    @apply ring-blue-400;
  }

  .colorcard-blue:hover {
    @apply bg-blue-600;
    @apply ring-blue-600;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
    @apply ring-indigo-400;
  }

  .colorcard-indigo:hover {
    @apply bg-indigo-600;
    @apply ring-indigo-600;
  }

  .colorcard-purple {
    @apply bg-purple-400;
    @apply ring-purple-400;
  }

  .colorcard-purple:hover {
    @apply bg-purple-600;
    @apply ring-purple-600;
  }

  .colorcard-pink {
    @apply bg-pink-400;
    @apply ring-pink-400;
  }

  .colorcard-pink:hover {
    @apply bg-pink-600;
    @apply ring-pink-600;
  }
</style>
