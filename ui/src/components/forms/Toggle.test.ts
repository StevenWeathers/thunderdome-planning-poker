import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import Toggle from './Toggle.svelte';

describe('Toggle component', () => {
  it('should render successfully', () => {
    render(Toggle, {
      label: 'Enable notifications',
      id: 'notifications',
      name: 'notifications',
    });

    const checkbox = page.getByRole('checkbox');
    expect(checkbox).toBeTruthy();
  });

  it('should render label text', () => {
    render(Toggle, {
      label: 'Dark mode',
      id: 'darkmode',
    });

    expect(page.getByText('Dark mode')).toBeTruthy();
  });

  it('should set input attributes', () => {
    const { container } = render(Toggle, {
      label: 'Accept',
      id: 'accept',
      name: 'accept',
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    expect(input).toBeTruthy();
    expect(input.getAttribute('id')).toBe('accept');
    expect(input.getAttribute('name')).toBe('accept');
  });

  it('should reflect checked state', async () => {
    const { container } = render(Toggle, {
      label: 'Checked',
      id: 'checked',
      checked: true,
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    await expect.element(input).toBeChecked();
  });

  it('should call change handler when toggled', async () => {
    const changeHandler = vi.fn();
    const { container } = render(Toggle, {
      label: 'Toggle',
      id: 'toggle',
      changeHandler,
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    await userEvent.click(input);

    expect(changeHandler).toHaveBeenCalledTimes(1);
  });

  it('should toggle checked state when clicked', async () => {
    const { container } = render(Toggle, {
      label: 'Toggle state',
      id: 'toggle-state',
    });

    const input = container.querySelector('input[type="checkbox"]') as HTMLInputElement;
    await expect.element(input).not.toBeChecked();

    await userEvent.click(input);
    await expect.element(input).toBeChecked();
  });
});
