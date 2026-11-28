import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import HealthCheckResponseForm from './HealthCheckResponseForm.svelte';

const areas = [
  { name: 'Fun', healthy: 'We love it', unhealthy: 'We do not' },
  { name: 'Speed', healthy: '', unhealthy: '' },
];

describe('HealthCheckResponseForm component', () => {
  it('explains what each rating means', async () => {
    render(HealthCheckResponseForm, { areas, handleSubmit: vi.fn() });

    await expect.element(page.getByRole('heading', { name: 'How to rate' })).toBeInTheDocument();
    await expect.element(page.getByText("It doesn't have to be perfect.", { exact: false })).toBeInTheDocument();
  });

  it('requires a rating for every area before submitting', async () => {
    const handleSubmit = vi.fn();
    render(HealthCheckResponseForm, { areas, handleSubmit });

    const submit = page.getByRole('button', { name: 'Submit Response' });
    await expect.element(submit).toBeDisabled();

    await userEvent.click(page.getByRole('radio', { name: 'Good' }).first());
    await expect.element(submit).toBeDisabled();

    await userEvent.click(page.getByRole('radio', { name: 'Needs work' }).nth(1));
    await userEvent.click(page.getByRole('radio', { name: 'Improving' }).first());
    await userEvent.fill(page.getByLabelText('Speed: Comment (optional)'), 'Too many handoffs');
    await userEvent.click(submit);

    expect(handleSubmit).toHaveBeenCalledWith([
      { area: 'Fun', rating: 'green', trend: 'improving', comment: '' },
      { area: 'Speed', rating: 'red', trend: 'stable', comment: 'Too many handoffs' },
    ]);
  });

  it('starts with the answers already given', async () => {
    render(HealthCheckResponseForm, {
      areas,
      answers: [
        { area: 'Fun', rating: 'yellow', trend: 'declining', comment: 'Long week' },
        { area: 'Speed', rating: 'green', trend: 'stable', comment: '' },
      ],
      handleSubmit: vi.fn(),
    });

    await expect.element(page.getByRole('radio', { name: 'Some concerns' }).first()).toBeChecked();
    await expect.element(page.getByRole('radio', { name: 'Getting worse' }).first()).toBeChecked();
    await expect.element(page.getByLabelText('Fun: Comment (optional)')).toHaveValue('Long week');
    await expect.element(page.getByRole('button', { name: 'Submit Response' })).toBeEnabled();
  });
});
