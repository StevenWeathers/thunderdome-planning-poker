import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import StoryCard from './StoryCard.svelte';

describe('StoryCard', () => {
  const story = {
    id: 'story-1',
    name: 'Story',
    content: '',
    color: 'gray',
    points: '3',
    closed: false,
    link: '',
    annotations: [],
    sort_order: 'a',
    comments: [],
  };

  const renderCard = () => {
    const open = vi.fn();
    const toggleStoryForm = vi.fn(() => open);
    render(StoryCard, { story, scale: 1, toggleStoryForm });
    return { open, toggleStoryForm };
  };

  it('opens the story form once from the points badge', async () => {
    const { open, toggleStoryForm } = renderCard();

    await userEvent.click(page.getByTestId('story-points'));

    expect(toggleStoryForm).toHaveBeenCalledExactlyOnceWith(story);
    expect(open).toHaveBeenCalledOnce();
  });

  it('opens the story form at the discussion from the comments button', async () => {
    const { toggleStoryForm } = renderCard();

    await userEvent.click(page.getByTestId('story-comments'));

    expect(toggleStoryForm).toHaveBeenCalledExactlyOnceWith(story, { showDiscussion: true });
  });
});
