import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import EndGameModal from './EndGameModal.svelte';

describe('EndGameModal component', () => {
  const baseProps = {
    toggleModal: vi.fn(),
    handleSubmit: vi.fn(),
    notifications: {} as any,
    xfetch: {} as any,
  };

  it('should render successfully', () => {
    render(EndGameModal, {
      ...baseProps,
    });

    const form = page.getByRole('form', { name: 'createBattle' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const { container } = render(EndGameModal, {
      ...baseProps,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render end game reason select input', () => {
    const { container } = render(EndGameModal, {
      ...baseProps,
    });

    const select = container.querySelector('select[name="endGameReason"]');
    expect(select).toBeTruthy();
    expect(select?.getAttribute('id')).toBe('endGameReason');
  });

  it('should render select options for end game reasons', () => {
    const { container } = render(EndGameModal, {
      ...baseProps,
    });

    const options = Array.from(container.querySelectorAll('option'));
    const values = options.map(option => option.getAttribute('value'));

    expect(values).toContain('Completed');
    expect(values).toContain('Cancelled');
  });

  it('should disable submit when no reason is selected', async () => {
    const handleSubmit = vi.fn();
    const { container } = render(EndGameModal, {
      ...baseProps,
      handleSubmit,
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(true);
    expect(handleSubmit).not.toHaveBeenCalled();
  });

  it('should submit the selected end game reason', async () => {
    const handleSubmit = vi.fn();
    const { container } = render(EndGameModal, {
      ...baseProps,
      handleSubmit,
    });

    const select = container.querySelector('select[name="endGameReason"]') as HTMLSelectElement;
    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;

    await userEvent.selectOptions(select, 'Completed');
    await userEvent.click(button);

    expect(handleSubmit).toHaveBeenCalledWith({ endGameReason: 'Completed' });
    expect(handleSubmit).toHaveBeenCalledTimes(1);
  });

  it('should update selected reason value when changed', async () => {
    const { container } = render(EndGameModal, {
      ...baseProps,
    });

    const select = container.querySelector('select[name="endGameReason"]') as HTMLSelectElement;

    await userEvent.selectOptions(select, 'Cancelled');

    await expect.element(select).toHaveValue('Cancelled');
  });
});
