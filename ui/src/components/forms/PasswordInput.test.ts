import { describe, it, expect } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import PasswordInput from './PasswordInput.svelte';

describe('PasswordInput component', () => {
  it('should render successfully', () => {
    render(PasswordInput, { name: 'password' });

    const input = page.getByRole('textbox');
    expect(input).toBeTruthy();
  });

  it('should render a password input by default', () => {
    const { container } = render(PasswordInput, { name: 'password' });

    const input = container.querySelector('input[name="password"]') as HTMLInputElement;
    expect(input).toBeTruthy();
    expect(input.getAttribute('type')).toBe('password');
  });

  it('should toggle visibility when button is clicked', async () => {
    const { container } = render(PasswordInput, { name: 'password' });

    const toggle = container.querySelector('button[type="button"]') as HTMLButtonElement;

    let input = container.querySelector('input[name="password"]') as HTMLInputElement;
    expect(input.getAttribute('type')).toBe('password');

    await userEvent.click(toggle);
    input = container.querySelector('input[name="password"]') as HTMLInputElement;
    expect(input.getAttribute('type')).toBe('text');

    await userEvent.click(toggle);
    input = container.querySelector('input[name="password"]') as HTMLInputElement;
    expect(input.getAttribute('type')).toBe('password');
  });

  it('should update value when typing', async () => {
    const { container } = render(PasswordInput, { name: 'password' });

    const input = container.querySelector('input[name="password"]') as HTMLInputElement;
    await userEvent.fill(input, 'secret123');

    await expect.element(input).toHaveValue('secret123');
  });

  it('should pass through attributes to the input', () => {
    const { container } = render(PasswordInput, {
      name: 'password',
      id: 'password',
      placeholder: 'Enter password',
      required: true,
      autocomplete: 'current-password',
    });

    const input = container.querySelector('input[name="password"]') as HTMLInputElement;
    expect(input.getAttribute('id')).toBe('password');
    expect(input.getAttribute('placeholder')).toBe('Enter password');
    expect(input.getAttribute('required')).toBe('');
    expect(input.getAttribute('autocomplete')).toBe('current-password');
  });
});
