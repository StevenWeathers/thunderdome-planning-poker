import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import OtpForm from './OtpForm.svelte';

describe('OtpForm component', () => {
  it('should render successfully', () => {
    const authMfa = vi.fn();
    render(OtpForm, { authMfa });

    const form = page.getByRole('form', { name: 'authMfa' });
    expect(form).toBeTruthy();
  });

  it('should render label and submit text', () => {
    const authMfa = vi.fn();
    render(OtpForm, { authMfa });

    expect(page.getByText('Authenticator Token')).toBeTruthy();
    expect(page.getByRole('button', { name: 'Login' })).toBeTruthy();
  });

  it('should render otp input with expected attributes', () => {
    const authMfa = vi.fn();
    const { container } = render(OtpForm, { authMfa });

    const input = container.querySelector('input[name="otp"]') as HTMLInputElement;
    expect(input).toBeTruthy();
    expect(input.getAttribute('id')).toBe('otp');
    expect(input.getAttribute('placeholder')).toBe('Enter code');
    expect(input.getAttribute('required')).toBe('');
    expect(input.getAttribute('inputmode')).toBe('numeric');
    expect(input.getAttribute('pattern')).toBe('[0-9]*');
    expect(input.getAttribute('autocomplete')).toBe('one-time-code');
  });

  it('should update value when typing', async () => {
    const authMfa = vi.fn();
    const { container } = render(OtpForm, { authMfa });

    const input = container.querySelector('input[name="otp"]') as HTMLInputElement;
    await userEvent.fill(input, '123456');

    await expect.element(input).toHaveValue('123456');
  });

  it('should call authMfa with token on submit', async () => {
    const authMfa = vi.fn();
    const { container } = render(OtpForm, { authMfa });

    const input = container.querySelector('input[name="otp"]') as HTMLInputElement;
    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;

    await userEvent.fill(input, '987654');
    await button.click();

    expect(authMfa).toHaveBeenCalledWith('987654');
    expect(authMfa).toHaveBeenCalledTimes(1);
  });

  it('should not submit when token is empty due to required attribute', async () => {
    const authMfa = vi.fn();
    const { container } = render(OtpForm, { authMfa });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    await button.click();

    expect(authMfa).not.toHaveBeenCalled();
  });
});
