import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import GoalColumns from './GoalColumns.svelte';

const story = (id: string, name: string) => ({
  id,
  name,
  color: 'gray',
  points: '',
  link: '',
  comments: [],
  closed: false,
  sort_order: id,
  content: '',
  annotations: [],
});

const goal = () => ({
  id: 'goal-1',
  name: 'Checkout',
  sort_order: 'a1',
  default_story_color: 'gray',
  personas: [],
  columns: [
    {
      id: 'column-1',
      name: 'To Do',
      sort_order: 'a1',
      default_story_color: 'gray',
      personas: [],
      stories: [story('story-1', 'Cart'), story('story-2', 'Payment'), story('story-3', 'Receipt')],
    },
  ],
});

function renderColumns() {
  const sendSocketEvent = vi.fn();
  const testGoal = goal();
  render(GoalColumns, {
    goals: [testGoal],
    goal: testGoal,
    goalIndex: 0,
    columnOrderEditMode: false,
    scale: 1,
    toggleColumnEdit: () => () => {},
    addStory: () => () => {},
    sendSocketEvent,
    notifications: {} as any,
    colorLegend: [],
    users: [],
    personas: [],
  });
  return { sendSocketEvent };
}

async function startSelection() {
  await page.getByTestId('column-actions-menu').click();
  await page.getByTestId('delete-stories-column-action').click();
}

const card = (name: string) => page.getByTestId('column-story').filter({ hasText: name });
const deleteButton = () => page.getByTestId('story-selection-delete');

describe('GoalColumns bulk story delete', () => {
  it('deletes only the selected stories after confirmation', async () => {
    const { sendSocketEvent } = renderColumns();
    await startSelection();
    await expect.element(deleteButton()).toBeDisabled();

    await card('Cart').click();
    await card('Receipt').click();
    await deleteButton().click();
    expect(sendSocketEvent).not.toHaveBeenCalled();
    await page.getByTestId('confirm-confirm').click();

    expect(sendSocketEvent).toHaveBeenCalledExactlyOnceWith(
      'delete_stories',
      JSON.stringify({ storyIds: ['story-1', 'story-3'] }),
    );
    await expect.element(page.getByTestId('story-add')).toBeVisible();
  });

  it('selects all stories and allows deselecting one', async () => {
    renderColumns();
    await startSelection();

    await page.getByTestId('story-selection-all').click();
    await expect.element(deleteButton()).toHaveTextContent('Delete (3)');

    await card('Payment').click();
    await expect.element(deleteButton()).toHaveTextContent('Delete (2)');
    await expect.element(card('Payment')).toHaveAttribute('aria-pressed', 'false');
  });

  it('selects a story when its comments button area is clicked', async () => {
    renderColumns();
    await startSelection();

    // the card contents are inert, so the click lands on the card itself
    await page.getByTestId('story-comments').first().click({ force: true });

    await expect.element(deleteButton()).toHaveTextContent('Delete (1)');
  });

  it('leaves selection mode on Escape without deleting', async () => {
    const { sendSocketEvent } = renderColumns();
    await startSelection();
    await card('Cart').click();

    await userEvent.keyboard('{Escape}');

    await expect.element(page.getByTestId('story-add')).toBeVisible();
    expect(page.getByTestId('story-select-indicator').elements()).toHaveLength(0);
    expect(sendSocketEvent).not.toHaveBeenCalled();
  });
});
