import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import BecomeFacilitator from './BecomeFacilitator.svelte';

describe('BecomeFacilitator component', () => {
  it('should render successfully', () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const form = page.getByRole('form', { name: 'becomeFacilitator' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render facilitator code input field', () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const input = container.querySelector('input[name="facilitatorCode"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('facilitatorCode');
  });

  it('should render submit button', () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const button = page.getByRole('button', { name: /save/i });
    expect(button).toBeTruthy();
  });

  it('should fire handleBecomeFacilitator with code when form is submitted', async () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const input = container.querySelector('input[name="facilitatorCode"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.fill(input, 'facilitator-123');
    await button.click();

    expect(handleBecomeFacilitator).toHaveBeenCalledWith('facilitator-123');
    expect(handleBecomeFacilitator).toHaveBeenCalledTimes(1);
  });

  it('should submit form on Enter key when input is focused', async () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const input = container.querySelector('input[name="facilitatorCode"]') as HTMLInputElement;

    await userEvent.fill(input, 'keyboard-submit');
    await userEvent.keyboard('{Enter}');

    expect(handleBecomeFacilitator).toHaveBeenCalledWith('keyboard-submit');
    expect(handleBecomeFacilitator).toHaveBeenCalledTimes(1);
  });

  it('should update facilitator code value when input changes', async () => {
    const toggleBecomeFacilitator = vi.fn();
    const handleBecomeFacilitator = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeFacilitator,
      handleBecomeFacilitator,
    });

    const input = container.querySelector('input[name="facilitatorCode"]') as HTMLInputElement;

    await userEvent.fill(input, 'new-code');

    await expect.element(input).toHaveValue('new-code');
  });
});
