import { describe, it, expect, vi } from 'vitest';
import { userEvent } from 'vitest/browser';
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
      additionalDetailsExpanded: true,
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
