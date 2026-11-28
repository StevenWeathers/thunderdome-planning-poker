import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import StoryColumn from './StoryColumn.svelte';

const comment = (id: string, text: string) => ({
  id,
  comment: text,
  story_id: 'story-1',
  user_id: 'user-1',
  created_date: '2026-01-01T10:00:00Z',
  updated_date: '2026-01-01T10:00:00Z',
});

const column = (comments: ReturnType<typeof comment>[]) => ({
  default_story_color: 'gray',
  id: 'column-1',
  name: 'To Do',
  sort_order: 'a1',
  personas: [],
  stories: [
    {
      id: 'story-1',
      name: 'Cart',
      content: '',
      color: 'gray',
      points: '',
      link: '',
      closed: false,
      annotations: [],
      sort_order: 'a1',
      comments,
    },
  ],
});

function renderColumn(comments = [comment('c1', 'First thought'), comment('c2', 'Newest thought')]) {
  const goalColumn = column(comments);
  const goal = {
    id: 'goal-1',
    name: 'Checkout',
    sort_order: 'a1',
    default_story_color: 'gray',
    personas: [],
    columns: [goalColumn],
  };
  return render(StoryColumn, {
    goals: [goal],
    goal,
    goalIndex: 0,
    goalColumn,
    columnIndex: 0,
    columnOrderEditMode: false,
    columnWidth: '10rem',
    scale: 1,
    sendSocketEvent: vi.fn(),
    notifications: {} as any,
    colorLegend: [],
    users: [{ id: 'user-1', name: 'Alex', avatar: 'none', gravatarHash: '' }],
  });
}

const countButton = () => page.getByTestId('story-comments');
const preview = () => page.getByTestId('story-comment-preview');
const storyForm = () => page.getByRole('dialog', { name: /storyboard story/i });

describe('StoryCommentPreview', () => {
  it('shows only the latest comment and its author, without opening the story', async () => {
    renderColumn();

    await userEvent.click(countButton());

    await expect.element(preview()).toBeVisible();
    await expect.element(preview()).toHaveTextContent('Newest thought');
    await expect.element(preview()).toHaveTextContent('Alex');
    await expect.element(preview()).not.toHaveTextContent('First thought');
    await expect.element(countButton()).toHaveAttribute('aria-expanded', 'true');
    expect(storyForm().query()).toBeNull();
  });

  it('closes when the comment count is clicked again', async () => {
    renderColumn();

    await userEvent.click(countButton());
    await userEvent.click(countButton());

    expect(preview().query()).toBeNull();
  });

  it('closes on a click outside', async () => {
    renderColumn();
    await userEvent.click(countButton());

    const outside = document.body.appendChild(document.createElement('button'));
    outside.textContent = 'outside';

    await userEvent.click(outside);

    expect(preview().query()).toBeNull();
    outside.remove();
  });

  it.each(['{Enter}', ' '])('opens from the keyboard with %j and returns focus on Escape', async key => {
    renderColumn();
    (countButton().element() as HTMLElement).focus();

    await userEvent.keyboard(key);
    await expect.element(preview()).toBeVisible();
    await expect.element(page.getByTestId('story-comment-preview-open')).toHaveFocus();
    expect(storyForm().query()).toBeNull();

    await userEvent.keyboard('{Escape}');
    expect(preview().query()).toBeNull();
    await expect.element(countButton()).toHaveFocus();
  });

  it('closes when tabbing away', async () => {
    renderColumn();
    await userEvent.click(countButton());

    await userEvent.tab();

    expect(preview().query()).toBeNull();
  });

  it('opens the discussion from the preview', async () => {
    renderColumn();
    await userEvent.click(countButton());

    await userEvent.click(page.getByRole('button', { name: 'Open discussion (2)' }));

    await expect.element(storyForm()).toBeVisible();
    expect(preview().query()).toBeNull();
  });

  it('opens the story directly when there are no comments', async () => {
    renderColumn([]);

    await userEvent.click(countButton());

    await expect.element(storyForm()).toBeVisible();
    expect(preview().query()).toBeNull();
  });

  it.each([
    [0, '0 comments'],
    [1, '1 comment'],
    [2, '2 comments'],
  ])('labels the button for %i comments', async (count, label) => {
    renderColumn(Array.from({ length: count }, (_, i) => comment(`c${i}`, 'text')));

    await expect.element(countButton()).toHaveAttribute('title', label);
    await expect.element(countButton()).toHaveAccessibleName(label);
  });
});
