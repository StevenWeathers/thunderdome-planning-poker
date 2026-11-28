<script lang="ts">
  import { onMount } from 'svelte';
  import LL from '../../i18n/i18n-svelte';
  import { formatTimestamp } from '../../dateUtils';
  import UserAvatar from '../user/UserAvatar.svelte';
  import type { StoryboardStory } from '../../types/storyboard';

  interface Props {
    story: StoryboardStory;
    users: any[];
    anchor: HTMLElement;
    onClose: () => void;
    onOpenDiscussion: () => void;
  }

  let { story, users, anchor, onClose, onOpenDiscussion }: Props = $props();

  const gap = 4;

  let popover: HTMLDivElement;
  let openButton: HTMLButtonElement;

  // comments arrive ordered by created date, oldest first
  let comment = $derived(story.comments[story.comments.length - 1]);
  let author = $derived(users.find(u => u.id === comment.user_id));
  let authorName = $derived(author?.name || $LL.unknownWarrior());

  // The popover renders in the top layer so the horizontally scrolling goal can't clip it,
  // which means it has to be placed next to its button by hand
  function position() {
    const button = anchor.getBoundingClientRect();
    const { offsetWidth: width, offsetHeight: height } = popover;
    const rtl = getComputedStyle(anchor).direction === 'rtl';
    const start = rtl ? button.right - width : button.left;
    const below = button.bottom + gap;
    const top = below + height > window.innerHeight ? button.top - height - gap : below;

    popover.style.left = `${Math.max(gap, Math.min(start, window.innerWidth - width - gap))}px`;
    popover.style.top = `${Math.max(gap, top)}px`;
  }

  function handlePointerDown(e: PointerEvent) {
    const target = e.target as Node;
    // the anchor toggles the popover itself
    if (!popover.contains(target) && !anchor.contains(target)) {
      onClose();
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    // on Tab, move focus back first so it continues on from the button that opened the popover
    if (e.key === 'Escape' || (e.key === 'Tab' && popover.contains(e.target as Node))) {
      anchor.focus();
      onClose();
    }
  }

  function handleOpenDiscussion() {
    // so closing the story form returns focus to the card's comment button
    anchor.focus();
    onOpenDiscussion();
  }

  onMount(() => {
    popover.showPopover();
    position();
    openButton.focus({ preventScroll: true });

    document.addEventListener('pointerdown', handlePointerDown, true);
    document.addEventListener('scroll', position, true);

    return () => {
      document.removeEventListener('pointerdown', handlePointerDown, true);
      document.removeEventListener('scroll', position, true);
    };
  });

  // the height changes when the latest comment does
  $effect(() => {
    if (comment) {
      position();
    }
  });
</script>

<svelte:window onkeydown={handleKeydown} onresize={position} />

<div
  bind:this={popover}
  popover="manual"
  role="dialog"
  aria-label={$LL.storyboardStoryLatestComment()}
  class="fixed inset-auto m-0 w-72 max-w-[calc(100vw-0.5rem)] p-3 text-sm bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-200 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600"
  data-testid="story-comment-preview"
>
  <div class="flex items-center gap-2 mb-2">
    <UserAvatar
      warriorId={comment.user_id}
      userName={authorName}
      gravatarHash={author?.gravatarHash || ''}
      avatar={author?.avatar || ''}
      width={24}
      class="rounded-full flex-shrink-0"
    />
    <div class="min-w-0">
      <div class="font-semibold truncate text-gray-900 dark:text-white">{authorName}</div>
      <time class="block text-xs text-gray-500 dark:text-gray-400" datetime={comment.created_date}>
        {formatTimestamp(comment.created_date)}
      </time>
    </div>
  </div>
  <p class="whitespace-pre-wrap break-words line-clamp-6" data-testid="story-comment-preview-text">{comment.comment}</p>
  <button
    bind:this={openButton}
    type="button"
    class="mt-3 text-blue-600 dark:text-cyan-400 hover:underline"
    data-testid="story-comment-preview-open"
    onclick={handleOpenDiscussion}
  >
    {$LL.storyboardStoryOpenDiscussion()} ({story.comments.length})
  </button>
</div>
