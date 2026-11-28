import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import GoalOrder from './GoalOrder.svelte';

const goal = (id: string, name: string) => ({
  id,
  name,
  sort_order: id,
  default_story_color: null,
  personas: [],
  columns: [],
});

const checkout = goal('goal-1', 'Checkout');
const search = goal('goal-2', 'Search');
const signup = goal('goal-3', 'Signup');

function renderGoalOrder() {
  const sendSocketEvent = vi.fn();
  render(GoalOrder, { goals: [checkout, search, signup], sendSocketEvent });
  return { sendSocketEvent };
}

function drop(items: ReturnType<typeof goal>[], id: string) {
  page
    .getByTestId('goal-order-dropzone')
    .element()
    .dispatchEvent(new CustomEvent('finalize', { detail: { items, info: { id } } }));
}

const names = () =>
  page
    .getByTestId('goal-order-item')
    .elements()
    .map(item => item.textContent?.trim());

describe('GoalOrder', () => {
  it('lists the goals in order', async () => {
    renderGoalOrder();

    await expect.element(page.getByText('Drag goals to change their order')).toBeVisible();
    expect(names()).toEqual(['Checkout', 'Search', 'Signup']);
  });

  it('places a moved goal before the goal that follows it', async () => {
    const { sendSocketEvent } = renderGoalOrder();

    drop([signup, checkout, search], 'goal-3');

    expect(sendSocketEvent).toHaveBeenCalledExactlyOnceWith(
      'move_goal',
      JSON.stringify({ goalId: 'goal-3', placeBefore: 'goal-1' }),
    );
    await expect.poll(names).toEqual(['Signup', 'Checkout', 'Search']);
  });

  it('sends an empty placeBefore for a goal moved to the end', () => {
    const { sendSocketEvent } = renderGoalOrder();

    drop([search, signup, checkout], 'goal-1');

    expect(sendSocketEvent).toHaveBeenCalledExactlyOnceWith(
      'move_goal',
      JSON.stringify({ goalId: 'goal-1', placeBefore: '' }),
    );
  });
});
