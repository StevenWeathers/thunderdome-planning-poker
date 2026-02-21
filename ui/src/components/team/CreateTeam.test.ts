import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CreateTeam from './CreateTeam.svelte';

describe('CreateTeam component', () => {
  it('should render successfully', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    render(CreateTeam, {
      toggleCreate,
      handleCreate,
    });

    const form = page.getByRole('form', { name: 'createTeam' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateTeam, {
      toggleCreate,
      handleCreate,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render team name input field', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateTeam, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="teamName"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('teamName');
  });

  it('should disable submit button when team name is empty', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateTeam, {
      toggleCreate,
      handleCreate,
      teamName: '',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(true);
  });

  it('should enable submit button when team name is provided', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateTeam, {
      toggleCreate,
      handleCreate,
      teamName: 'Alpha Team',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(false);
  });

  it('should fire handleCreate with team name when form is submitted', async () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateTeam, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="teamName"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.fill(input, 'Product Team');
    await button.click();

    expect(handleCreate).toHaveBeenCalledWith('Product Team');
    expect(handleCreate).toHaveBeenCalledTimes(1);
  });

  it('should update team name value when input changes', async () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateTeam, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="teamName"]') as HTMLInputElement;

    await userEvent.fill(input, 'Operations');

    await expect.element(input).toHaveValue('Operations');
  });
});
