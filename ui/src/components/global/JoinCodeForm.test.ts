import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import JoinCodeForm from './JoinCodeForm.svelte';

describe('JoinCodeForm component', () => {
  it('should render successfully', () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const form = page.getByRole('form', { name: 'joinPasscodeForm' });
    expect(form).toBeTruthy();
  });

  it('should render with default submit button text', () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const button = page.getByRole('button', { name: 'Join' });
    expect(button).toBeTruthy();
  });

  it('should render with custom submit button text', () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
      submitText: 'Enter Game',
    });

    const button = page.getByRole('button', { name: 'Enter Game' });
    expect(button).toBeTruthy();
  });

  it('should render password input field', () => {
    const handleSubmit = vi.fn();
    const { container } = render(JoinCodeForm, {
      handleSubmit,
    });

    const input = container.querySelector('input[name="passCode"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('type')).toBe('password');
    expect(input?.getAttribute('required')).toBe('');
  });

  it('should fire handleSubmit with passcode when form is submitted', async () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const input = page.getByRole('textbox', { name: /passCode/i });
    const button = page.getByRole('button', { name: 'Join' });

    // Enter passcode
    await userEvent.fill(input, 'test-passcode-123');
    await button.click();

    expect(handleSubmit).toHaveBeenCalledWith('test-passcode-123');
    expect(handleSubmit).toHaveBeenCalledTimes(1);
  });

  it('should not submit form when passcode is empty due to required attribute', async () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const button = page.getByRole('button', { name: 'Join' });
    await button.click();

    // Form should not submit because input is required
    expect(handleSubmit).not.toHaveBeenCalled();
  });

  it('should submit form on Enter key when input is focused', async () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const input = page.getByRole('textbox', { name: /passCode/i });

    // Enter passcode and submit with Enter key
    await userEvent.fill(input, 'keyboard-submit-test');
    await userEvent.keyboard('{Enter}');

    expect(handleSubmit).toHaveBeenCalledWith('keyboard-submit-test');
    expect(handleSubmit).toHaveBeenCalledTimes(1);
  });

  it('should update passcode value when input changes', async () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const input = page.getByRole('textbox', { name: /passCode/i });

    // Change the input value
    await userEvent.fill(input, 'new-passcode');

    await expect.element(input).toHaveValue('new-passcode');
  });

  it('should handle multiple form submissions with different passcodes', async () => {
    const handleSubmit = vi.fn();
    render(JoinCodeForm, {
      handleSubmit,
    });

    const input = page.getByRole('textbox', { name: /passCode/i });
    const button = page.getByRole('button', { name: 'Join' });

    // First submission
    await userEvent.fill(input, 'first-code');
    await button.click();

    expect(handleSubmit).toHaveBeenCalledWith('first-code');

    // Second submission
    await userEvent.clear(input);
    await userEvent.fill(input, 'second-code');
    await button.click();

    expect(handleSubmit).toHaveBeenCalledWith('second-code');
    expect(handleSubmit).toHaveBeenCalledTimes(2);
  });
});
