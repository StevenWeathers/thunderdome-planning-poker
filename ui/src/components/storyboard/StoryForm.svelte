<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import Editor from '../forms/Editor.svelte';
  import StoryColorSelector from './StoryColorSelector.svelte';
  import { Check, Link, RotateCcw, Trash2 } from '@lucide/svelte';
  import { onMount } from 'svelte';
  import Comment from '../comments/Comment.svelte';
  import CommentForm from '../comments/CommentForm.svelte';
  import { sanitizeStoryboardPoints, storyboardPointsMaxLength } from './storyPoints';

  import type { NotificationService } from '../../types/notifications';
  import type { UserDisplay } from '../../types/user';
  import type { StoryboardStory, StoryboardUser } from '../../types/storyboard';
  import CommentsHeader from '../comments/CommentsHeader.svelte';
  import CommentEmptyState from '../comments/CommentEmptyState.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';

  interface Props {
    toggleStoryForm?: any;
    sendSocketEvent?: any;
    notifications: NotificationService;
    story?: StoryboardStory;
    colorLegend?: any;
    users?: StoryboardUser[];
    /** values of the storyboard's point scale, empty for free text points */
    pointValues?: string[];
    /** scroll to the discussion when the form opens */
    showDiscussion?: boolean;
  }

  let {
    toggleStoryForm = () => {},
    sendSocketEvent = () => {},
    notifications,
    story = { id: '' } as StoryboardStory,
    colorLegend = [],
    users = [],
    pointValues = [],
    showDiscussion = false,
  }: Props = $props();

  const isAbsolute = new RegExp('^([a-z]+://|//)', 'i');

  let showDeleteConfirm = $state(false);
  let discussionElement: HTMLElement;
  let focusInput: any = $state();
  let storyPoints = $state('');
  let selectedStoryColor = $state('');

  // keeps points set before the scale was chosen or changed selectable, so they aren't silently lost
  const pointOptions = $derived(
    storyPoints === '' || pointValues.includes(storyPoints) ? pointValues : [storyPoints, ...pointValues],
  );

  $effect(() => {
    selectedStoryColor = story.color;
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
    const points = sanitizeStoryboardPoints((evt.target as HTMLInputElement).value);
    storyPoints = points;
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
      notifications.danger($LL.storyboardStoryLinkAbsoluteError());
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

  onMount(() => {
    focusInput?.focus();
    if (showDiscussion) {
      discussionElement.scrollIntoView({ block: 'start' });
    }
  });

  $effect(() => {
    if (story?.points !== undefined) {
      storyPoints = sanitizeStoryboardPoints(story.points);
    }
  });
</script>

<Modal closeModal={toggleStoryForm} widthClasses="w-full md:w-11/12 max-w-screen-2xl" ariaLabel={$LL.modalStoryboardStory()}>
  <div class="flex flex-col h-[calc(100dvh-8rem)]">
    <!-- right padding keeps clear of the modal close button, which is on the right in both directions -->
    <div class="pb-4 px-1 pr-10 border-b border-gray-200 dark:border-gray-700">
      <label for="storyName" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
        {$LL.storyboardStoryName()} <span class="text-red-500">*</span>
      </label>
      <TextInput
        id="storyName"
        onchange={updateName}
        value={story.name}
        placeholder={$LL.storyboardStoryNamePlaceholder()}
        name="storyName"
        bind:this={focusInput}
      />
    </div>

    <!-- one scrolling area, with the details beside the description on large screens and below it on small ones -->
    <div class="flex-1 min-h-0 overflow-y-auto overscroll-contain dark:[color-scheme:dark] pt-4 px-1 pb-1">
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_22rem] lg:grid-rows-[auto_1fr]">
        <div class="story-description">
          <label for="storyDescription" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
            {$LL.storyboardStoryDescription()}
          </label>
          <div class="border border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden bg-white">
            <Editor
              content={story.content}
              placeholder={$LL.storyboardStoryDescriptionPlaceholder()}
              id="storyDescription"
              handleTextChange={(c: string) => {
                updateContent(c);
              }}
            />
          </div>
        </div>

        <aside
          class="flex flex-col gap-5 self-start rounded-xl p-4 bg-gray-50 dark:bg-gray-700/30 border border-gray-200/50 dark:border-gray-600/30
            lg:col-start-2 lg:row-start-1 lg:row-span-2 lg:sticky lg:top-0 lg:max-h-[calc(100dvh-16rem)] lg:overflow-y-auto"
          aria-label={$LL.storyboardStoryDetails()}
        >
          <div>
            <div class="text-gray-700 dark:text-gray-300 mb-2 text-lg">{$LL.storyboardStoryStatus()}</div>
            {#if !story.closed}
              <HollowButton color="orange" onClick={markClosed}>
                <Check class="w-4 h-4 inline me-2" aria-hidden="true" />{$LL.storyboardStoryMarkClosed()}
              </HollowButton>
            {:else}
              <HollowButton color="green" onClick={markOpen}>
                <RotateCcw class="w-4 h-4 inline me-2" aria-hidden="true" />{$LL.storyboardStoryReopen()}
              </HollowButton>
            {/if}
          </div>

          <div>
            <label for="storyPoints" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
              {$LL.storyboardStoryPoints()}
            </label>
            {#if pointValues.length > 0}
              <SelectInput id="storyPoints" name="storyPoints" value={storyPoints} onchange={updatePoints}>
                <option value="">{$LL.storyboardNoPoints()}</option>
                {#each pointOptions as point (point)}
                  <option value={point}>{point}</option>
                {/each}
              </SelectInput>
            {:else}
              <input
                class="bg-white dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                          rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                          focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="storyPoints"
                bind:value={storyPoints}
                onchange={updatePoints}
                maxlength={storyboardPointsMaxLength}
                placeholder={$LL.storyboardStoryPointsPlaceholder()}
                name="storyPoints"
              />
            {/if}
          </div>

          <div>
            <div class="text-gray-700 dark:text-gray-300 mb-2 text-lg">{$LL.storyboardStoryColor()}</div>
            <StoryColorSelector bind:value={selectedStoryColor} {colorLegend} onSelect={color => changeColor(color)()} />
          </div>

          <div>
            <label for="storyLink" class="block text-gray-700 dark:text-gray-300 mb-2 text-lg">
              {$LL.storyboardStoryLink()}
              <span class="text-gray-500 dark:text-gray-400 font-normal ms-1">{$LL.optional()}</span>
            </label>
            <TextInput
              id="storyLink"
              onchange={updateLink}
              value={story.link}
              placeholder="https://jira.company.com/browse/PROJ-123"
              name="storyLink"
              type="url"
              icon={Link}
              iconPosition="start"
            />
          </div>

          <div class="border-t border-gray-200 dark:border-gray-700 pt-4">
            <div class="text-red-800 dark:text-red-400 mb-2 text-lg">{$LL.storyboardStoryDangerZone()}</div>
            <HollowButton color="red" onClick={() => (showDeleteConfirm = true)} testid="story-delete">
              <Trash2 class="w-4 h-4 inline me-2" aria-hidden="true" />{$LL.storyboardStoryDelete()}
            </HollowButton>
          </div>
        </aside>

        <div
          bind:this={discussionElement}
          class="flex flex-col gap-3 self-start rounded-xl p-4 bg-gray-50 dark:bg-gray-700/30 border border-gray-200/50 dark:border-gray-600/30 lg:col-start-1"
          data-testid="story-discussion"
        >
          <CommentsHeader commentsCount={story.comments ? story.comments.length : 0} />

          {#if story.comments && story.comments.length > 0}
            <div class="flex flex-col gap-3">
              {#each story.comments as comment}
                <Comment {comment} {userMap} handleEdit={handleCommentEdit} handleDelete={handleCommentDelete} />
              {/each}
            </div>
          {:else}
            <CommentEmptyState description={$LL.storyboardStoryNoComments()} />
          {/if}

          <CommentForm onSubmit={handleCommentSubmit} />
        </div>
      </div>
    </div>
  </div>

  {#if showDeleteConfirm}
    <!-- outside the sticky sidebar, whose stacking context would trap the dialog -->
    <DeleteConfirmation
      toggleDelete={() => (showDeleteConfirm = false)}
      handleDelete={handleStoryDelete}
      confirmText={$LL.storyboardStoryDeleteConfirm()}
      confirmBtnText={$LL.delete()}
    />
  {/if}
</Modal>

<style>
  .story-description :global(.ql-editor) {
    min-height: max(12rem, 30vh);
  }
</style>
