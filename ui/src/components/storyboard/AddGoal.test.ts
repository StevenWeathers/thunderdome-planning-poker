import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import AddGoal from './AddGoal.svelte';

describe('AddGoal component', () => {
  it('should render successfully', () => {
    const handleGoalAdd = vi.fn();
    const toggleAddGoal = vi.fn();
    render(AddGoal, {
      handleGoalAdd,
      toggleAddGoal,
    });

    const form = page.getByRole('form', { name: 'addGoal' });
    expect(form).toBeTruthy();
  });

  it('should render goal name input field', () => {
    const handleGoalAdd = vi.fn();
    const toggleAddGoal = vi.fn();
    const { container } = render(AddGoal, {
      handleGoalAdd,
      toggleAddGoal,
    });

    const input = container.querySelector('input[name="goalName"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('goalName');
  });

  it('should fire handleGoalAdd with goal name when form is submitted', async () => {
    const handleGoalAdd = vi.fn();
    const handleGoalRevision = vi.fn();
    const toggleAddGoal = vi.fn();
    const { container } = render(AddGoal, {
      handleGoalAdd,
      handleGoalRevision,
      toggleAddGoal,
    });

    const input = container.querySelector('input[name="goalName"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.fill(input, 'New Goal');
    await button.click();

    expect(handleGoalAdd).toHaveBeenCalledWith('New Goal');
    expect(handleGoalAdd).toHaveBeenCalledTimes(1);
    expect(handleGoalRevision).not.toHaveBeenCalled();
    expect(toggleAddGoal).toHaveBeenCalledTimes(1);
  });

  it('should fire handleGoalRevision when goalId is provided', async () => {
    const handleGoalAdd = vi.fn();
    const handleGoalRevision = vi.fn();
    const toggleAddGoal = vi.fn();
    const { container } = render(AddGoal, {
      handleGoalAdd,
      handleGoalRevision,
      toggleAddGoal,
      goalId: 'goal-1',
      goalName: 'Existing Goal',
    });

    const input = container.querySelector('input[name="goalName"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.clear(input);
    await userEvent.fill(input, 'Revised Goal');
    await button.click();

    expect(handleGoalRevision).toHaveBeenCalledWith({
      goalId: 'goal-1',
      name: 'Revised Goal',
    });
    expect(handleGoalRevision).toHaveBeenCalledTimes(1);
    expect(handleGoalAdd).not.toHaveBeenCalled();
    expect(toggleAddGoal).toHaveBeenCalledTimes(1);
  });

  it('should update goal name value when input changes', async () => {
    const handleGoalAdd = vi.fn();
    const toggleAddGoal = vi.fn();
    const { container } = render(AddGoal, {
      handleGoalAdd,
      toggleAddGoal,
    });

    const input = container.querySelector('input[name="goalName"]') as HTMLInputElement;

    await userEvent.fill(input, 'Updated Goal');

    await expect.element(input).toHaveValue('Updated Goal');
  });

  it('should allow submission with empty goal name', async () => {
    const handleGoalAdd = vi.fn();
    const toggleAddGoal = vi.fn();
    render(AddGoal, {
      handleGoalAdd,
      toggleAddGoal,
    });

    const button = page.getByRole('button', { name: /save/i });
    await button.click();

    expect(handleGoalAdd).toHaveBeenCalledWith('');
    expect(handleGoalAdd).toHaveBeenCalledTimes(1);
    expect(toggleAddGoal).toHaveBeenCalledTimes(1);
  });
});
