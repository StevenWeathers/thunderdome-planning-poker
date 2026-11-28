import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import StoryForm from './StoryForm.svelte';

describe('StoryForm story points', () => {
  const notifications = { success: vi.fn(), danger: vi.fn(), warning: vi.fn(), info: vi.fn() } as any;
  const story = {
    id: 'story-1',
    name: 'Story',
    content: '',
    color: 'gray',
    points: '',
    closed: false,
    link: '',
    annotations: [],
    sort_order: 'a',
    comments: [],
  };

  const renderForm = (props: { points?: string; pointValues?: string[] }) => {
    const sendSocketEvent = vi.fn();
    const { container } = render(StoryForm, {
      notifications,
      sendSocketEvent,
      story: { ...story, points: props.points ?? '' },
      pointValues: props.pointValues,
    });
    return { sendSocketEvent, field: () => container.querySelector('[name="storyPoints"]') as HTMLSelectElement };
  };

  const optionValues = (select: HTMLSelectElement) => Array.from(select.options).map(option => option.value);

  it('uses a free text input without a point scale', async () => {
    const { field } = renderForm({ points: '4' });

    await vi.waitFor(() => expect(field().tagName).toBe('INPUT'));
    await expect.element(field()).toHaveValue('4');
  });

  it('selects points from the point scale', async () => {
    const { field, sendSocketEvent } = renderForm({ pointValues: ['S', 'M', 'L'] });

    await vi.waitFor(() => expect(field().tagName).toBe('SELECT'));
    expect(optionValues(field())).toEqual(['', 'S', 'M', 'L']);

    await userEvent.selectOptions(field(), 'M');

    expect(sendSocketEvent).toHaveBeenCalledWith(
      'update_story_points',
      JSON.stringify({ storyId: 'story-1', points: 'M' }),
    );
  });

  it('keeps existing points that are not in the point scale', async () => {
    const { field } = renderForm({ points: '4', pointValues: ['S', 'M'] });

    await vi.waitFor(() => expect(optionValues(field())).toEqual(['', '4', 'S', 'M']));
    await expect.element(field()).toHaveValue('4');
  });
});

describe('StoryForm layout', () => {
  const notifications = { success: vi.fn(), danger: vi.fn(), warning: vi.fn(), info: vi.fn() } as any;
  const story = {
    id: 'story-1',
    name: 'Story',
    content: '',
    color: 'gray',
    points: '',
    closed: false,
    link: '',
    annotations: [],
    sort_order: 'a',
    comments: [],
  };

  const renderForm = () => {
    const sendSocketEvent = vi.fn();
    const toggleStoryForm = vi.fn();
    render(StoryForm, { notifications, sendSocketEvent, toggleStoryForm, story });
    return { sendSocketEvent, toggleStoryForm };
  };

  it('shows all fields and the discussion without further clicks', async () => {
    renderForm();

    await expect.element(page.getByLabelText('Story Points')).toBeVisible();
    await expect.element(page.getByLabelText('Story Link')).toBeVisible();
    await expect.element(page.getByRole('button', { name: 'Mark as Closed' })).toBeVisible();
    await expect.element(page.getByRole('button', { name: 'Delete Story' })).toBeVisible();
    await expect.element(page.getByTestId('story-discussion')).toBeVisible();
  });

  it('deletes the story only after confirmation', async () => {
    const { sendSocketEvent, toggleStoryForm } = renderForm();

    await userEvent.click(page.getByRole('button', { name: 'Delete Story' }));
    expect(sendSocketEvent).not.toHaveBeenCalled();

    await userEvent.click(page.getByTestId('confirm-cancel'));
    expect(sendSocketEvent).not.toHaveBeenCalled();

    await userEvent.click(page.getByRole('button', { name: 'Delete Story' }));
    await userEvent.click(page.getByTestId('confirm-confirm'));

    expect(sendSocketEvent).toHaveBeenCalledWith('delete_story', 'story-1');
    expect(toggleStoryForm).toHaveBeenCalled();
  });
});
