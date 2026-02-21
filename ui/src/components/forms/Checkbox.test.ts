import { describe, it, expect } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import Checkbox from './Checkbox.svelte';

describe('Checkbox component', () => {
  it('should render successfully', () => {
    render(Checkbox, {
      label: 'Enable alerts',
      id: 'alerts',
      name: 'alerts',
    });

    const checkbox = page.getByRole('checkbox');
    expect(checkbox).toBeTruthy();
  });

  it('should render label text', () => {
    render(Checkbox, {
      label: 'Receive updates',
      id: 'updates',
    });

    expect(page.getByText('Receive updates')).toBeTruthy();
  });

  it('should set input attributes', () => {
    const { container } = render(Checkbox, {
      label: 'Accept terms',
      id: 'terms',
      name: 'terms',
      value: 'yes',
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    expect(input).toBeTruthy();
    expect(input.getAttribute('id')).toBe('terms');
    expect(input.getAttribute('name')).toBe('terms');
    expect(input.getAttribute('value')).toBe('yes');
  });

  it('should reflect checked state', async () => {
    const { container } = render(Checkbox, {
      label: 'Checked option',
      id: 'checked',
      checked: true,
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    await expect.element(input).toBeChecked();
  });

  it('should toggle checked state when clicked', async () => {
    const { container } = render(Checkbox, {
      label: 'Toggle me',
      id: 'toggle',
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    await expect.element(input).not.toBeChecked();

    await userEvent.click(input);
    await expect.element(input).toBeChecked();
  });

  it('should render check icon when checked', () => {
    const { container } = render(Checkbox, {
      label: 'Icon check',
      id: 'icon-check',
      checked: true,
    });

    const icon = container.querySelector('svg');
    expect(icon).toBeTruthy();
  });

  it('should disable input when disabled is true', () => {
    const { container } = render(Checkbox, {
      label: 'Disabled',
      id: 'disabled',
      disabled: true,
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    expect(input.disabled).toBe(true);
  });
});
