import { describe, it, expect } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import TextInput from './TextInput.svelte';

describe('TextInput component', () => {
  it('should render successfully', () => {
    render(TextInput, { name: 'title' });

    const input = page.getByRole('textbox');
    expect(input).toBeTruthy();
  });

  it('should render with id, name, and placeholder', () => {
    const { container } = render(TextInput, {
      id: 'title',
      name: 'title',
      placeholder: 'Enter title',
    });

    const input = container.querySelector('input[name="title"]') as HTMLInputElement;
    expect(input.getAttribute('id')).toBe('title');
    expect(input.getAttribute('placeholder')).toBe('Enter title');
  });

  it('should default to text type', () => {
    const { container } = render(TextInput, { name: 'default' });

    const input = container.querySelector('input[name="default"]') as HTMLInputElement;
    expect(input.getAttribute('type')).toBe('text');
  });

  it('should render with provided type', () => {
    const { container } = render(TextInput, { name: 'email', type: 'email' });

    const input = container.querySelector('input[name="email"]') as HTMLInputElement;
    expect(input.getAttribute('type')).toBe('email');
  });

  it('should update value when typing', async () => {
    const { container } = render(TextInput, { name: 'title' });

    const input = container.querySelector('input[name="title"]') as HTMLInputElement;
    await userEvent.fill(input, 'Hello');

    await expect.element(input).toHaveValue('Hello');
  });

  it('should apply custom class', () => {
    const { container } = render(TextInput, { name: 'title', class: 'custom-class' });

    const input = container.querySelector('input[name="title"]') as HTMLInputElement;
    expect(input.classList.contains('custom-class')).toBe(true);
  });

  it('should be focusable via component method', () => {
    const { component } = render(TextInput, { name: 'title' });

    component.focus();

    const input = page.getByRole('textbox');
    const element = input.element() as HTMLInputElement;
    expect(document.activeElement).toBe(element);
  });
});
