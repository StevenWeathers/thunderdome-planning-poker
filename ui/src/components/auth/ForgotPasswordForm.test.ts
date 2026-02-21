import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ForgotPasswordForm from './ForgotPasswordForm.svelte';

describe('ForgotPasswordForm component', () => {
  it('should render successfully', () => {
    render(ForgotPasswordForm, {
      toggleForgotPassword: vi.fn(),
      onSubmit: vi.fn(),
    });

    const form = page.getByRole('form', { name: 'resetPassword' });
    expect(form).toBeTruthy();
  });

  it('should render header and helper text', () => {
    render(ForgotPasswordForm, {
      toggleForgotPassword: vi.fn(),
      onSubmit: vi.fn(),
    });

    expect(page.getByText('Forgot Your Password?')).toBeTruthy();
    expect(
      page.getByText("Don't fret! Just enter your email and we will send you instructions to reset your password."),
    ).toBeTruthy();
  });

  it('should render input with expected attributes', () => {
    const { container } = render(ForgotPasswordForm, {
      toggleForgotPassword: vi.fn(),
      onSubmit: vi.fn(),
    });

    const input = container.querySelector('input[name="yourResetEmail"]') as HTMLInputElement;
    expect(input).toBeTruthy();
    expect(input.getAttribute('id')).toBe('yourResetEmail');
    expect(input.getAttribute('type')).toBe('email');
    expect(input.getAttribute('placeholder')).toBe('Enter your email');
    expect(input.getAttribute('required')).toBe('');
  });

  it('should update email value when typing', async () => {
    const { container } = render(ForgotPasswordForm, {
      toggleForgotPassword: vi.fn(),
      onSubmit: vi.fn(),
    });

    const input = container.querySelector('input[name="yourResetEmail"]') as HTMLInputElement;
    await userEvent.fill(input, 'user@example.com');

    await expect.element(input).toHaveValue('user@example.com');
  });

  it('should call onSubmit with email on submit', async () => {
    const onSubmit = vi.fn();
    const { container } = render(ForgotPasswordForm, {
      toggleForgotPassword: vi.fn(),
      onSubmit,
    });

    const input = container.querySelector('input[name="yourResetEmail"]') as HTMLInputElement;
    const submitButton = container.querySelector('button[type="submit"]') as HTMLButtonElement;

    await userEvent.fill(input, 'reset@example.com');
    await submitButton.click();

    expect(onSubmit).toHaveBeenCalledWith('reset@example.com');
    expect(onSubmit).toHaveBeenCalledTimes(1);
  });

  it('should not submit when email is empty due to required attribute', async () => {
    const onSubmit = vi.fn();
    const { container } = render(ForgotPasswordForm, {
      toggleForgotPassword: vi.fn(),
      onSubmit,
    });

    const submitButton = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    await submitButton.click();

    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('should call toggleForgotPassword when return button is clicked', async () => {
    const toggleForgotPassword = vi.fn();
    const { container } = render(ForgotPasswordForm, {
      toggleForgotPassword,
      onSubmit: vi.fn(),
    });

    const returnButton = container.querySelector('button[type="button"]') as HTMLButtonElement;
    await returnButton.click();

    expect(toggleForgotPassword).toHaveBeenCalledTimes(1);
  });
});
