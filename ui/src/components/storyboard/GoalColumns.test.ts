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

describe('GoalColumns empty goal drop target', () => {
  function renderEmptyGoal(columnOrderEditMode: boolean) {
    const sendSocketEvent = vi.fn();
    const emptyGoal = { ...goal(), id: 'goal-2', columns: [] as ReturnType<typeof goal>['columns'] };
    render(GoalColumns, {
      goals: [emptyGoal],
      goal: emptyGoal,
      goalIndex: 0,
      columnOrderEditMode,
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

  it('shows a drop target with a minimum height in column order edit mode', async () => {
    renderEmptyGoal(true);

    await expect.element(page.getByTestId('column-drop-hint')).toBeVisible();
    // Tailwind is not loaded in component tests, so check the class rather than the rendered height
    const zone = page.getByTestId('goal-columns-dropzone').element();
    expect(zone.classList).toContain('min-h-40');
    expect(zone.children).toHaveLength(0);
  });

  it('shows no drop hint outside column order edit mode', async () => {
    renderEmptyGoal(false);

    expect(page.getByTestId('column-drop-hint').elements()).toHaveLength(0);
  });

  it('moves a column dropped onto an empty goal', async () => {
    const { sendSocketEvent } = renderEmptyGoal(true);
    const column = goal().columns[0];

    page
      .getByTestId('goal-columns-dropzone')
      .element()
      .dispatchEvent(new CustomEvent('finalize', { detail: { items: [column], info: { id: column.id } } }));

    expect(sendSocketEvent).toHaveBeenCalledExactlyOnceWith(
      'move_column',
      JSON.stringify({ goalId: 'goal-2', columnId: 'column-1', placeBefore: '' }),
    );
  });
});
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
