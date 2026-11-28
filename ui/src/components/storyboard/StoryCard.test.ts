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

  it('shows the comment count on an unnamed story that has comments', async () => {
    const comments = [
      { id: 'c1', comment: 'Hi', story_id: 'story-1', user_id: 'u', created_date: '', updated_date: '' },
    ];
    render(StoryCard, { story: { ...story, name: '', comments }, scale: 1 });

    await expect.element(page.getByTestId('story-comments')).toHaveAccessibleName('1 comment');
  });

  it('hides the comment count on an unnamed story without comments', () => {
    render(StoryCard, { story: { ...story, name: '' }, scale: 1 });

    expect(page.getByTestId('story-comments').query()).toBeNull();
  });

  it.each(['{Enter}', ' '])('acts once when %j is pressed on the comments button', async key => {
    const { toggleStoryForm } = renderCard();
    (page.getByTestId('story-comments').element() as HTMLElement).focus();

    await userEvent.keyboard(key);

    expect(toggleStoryForm).toHaveBeenCalledExactlyOnceWith(story, { showDiscussion: true });
  });
});
