import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import BecomeFacilitator from './BecomeFacilitator.svelte';

describe('BecomeFacilitator component', () => {
  it('should render successfully', () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const form = page.getByRole('form', { name: 'becomeLeader' });
    expect(form).toBeTruthy();
  });

  it('should render the modal with close functionality', () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render leader code input field', () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const input = container.querySelector('input[name="leaderCode"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('leaderCode');
  });

  it('should render submit button', () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const button = page.getByRole('button', { name: /save/i });
    expect(button).toBeTruthy();
  });

  it('should fire handleBecomeLeader with leader code when form is submitted', async () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const input = container.querySelector('input[name="leaderCode"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    // Enter leader code
    await userEvent.fill(input, 'facilitator-code-123');
    await button.click();

    expect(handleBecomeLeader).toHaveBeenCalledWith('facilitator-code-123');
    expect(handleBecomeLeader).toHaveBeenCalledTimes(1);
  });

  it('should submit form on Enter key when input is focused', async () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const input = container.querySelector('input[name="leaderCode"]') as HTMLInputElement;

    // Enter leader code and submit with Enter key
    await userEvent.fill(input, 'keyboard-submit-test');
    await userEvent.keyboard('{Enter}');

    expect(handleBecomeLeader).toHaveBeenCalledWith('keyboard-submit-test');
    expect(handleBecomeLeader).toHaveBeenCalledTimes(1);
  });

  it('should update leader code value when input changes', async () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const input = container.querySelector('input[name="leaderCode"]') as HTMLInputElement;

    // Change the input value
    await userEvent.fill(input, 'new-leader-code');

    await expect.element(input).toHaveValue('new-leader-code');
  });

  it('should handle multiple form submissions with different codes', async () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const input = container.querySelector('input[name="leaderCode"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    // First submission
    await userEvent.fill(input, 'first-code');
    await button.click();

    expect(handleBecomeLeader).toHaveBeenCalledWith('first-code');

    // Second submission
    await userEvent.clear(input);
    await userEvent.fill(input, 'second-code');
    await button.click();

    expect(handleBecomeLeader).toHaveBeenCalledWith('second-code');
    expect(handleBecomeLeader).toHaveBeenCalledTimes(2);
  });

  it('should allow submission with empty leader code', async () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const button = page.getByRole('button', { name: /save/i });
    await button.click();

    // Form should submit even with empty code (no required attribute)
    expect(handleBecomeLeader).toHaveBeenCalledWith('');
    expect(handleBecomeLeader).toHaveBeenCalledTimes(1);
  });

  it('should render label for leader code input', () => {
    const toggleBecomeLeader = vi.fn();
    const handleBecomeLeader = vi.fn();
    const { container } = render(BecomeFacilitator, {
      toggleBecomeLeader,
      handleBecomeLeader,
    });

    const label = container.querySelector('label[for="leaderCode"]');
    expect(label).toBeTruthy();
  });
});
