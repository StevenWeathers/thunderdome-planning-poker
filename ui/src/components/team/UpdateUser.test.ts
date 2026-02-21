import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import UpdateUser from './UpdateUser.svelte';

describe('UpdateUser component', () => {
  it('should render successfully', () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
    });

    const form = page.getByRole('form', { name: 'teamUpdateUser' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render disabled email input field', () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
    });

    const input = container.querySelector('input[name="userEmail"]') as HTMLInputElement;
    expect(input).toBeTruthy();
    expect(input.getAttribute('id')).toBe('userEmail');
    expect(input.disabled).toBe(true);
    expect(input.value).toBe('user@example.com');
  });

  it('should render role select input with options', () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
    });

    const select = container.querySelector('select[name="userRole"]');
    expect(select).toBeTruthy();
    expect(select?.getAttribute('id')).toBe('userRole');

    const options = Array.from(container.querySelectorAll('option')).map(option => option.getAttribute('value'));
    expect(options).toContain('ADMIN');
    expect(options).toContain('MEMBER');
  });

  it('should disable submit button when role is empty', () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
      role: '',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(true);
  });

  it('should enable submit button when role is provided', () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
      role: 'ADMIN',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(false);
  });

  it('should fire handleUpdate with userId and role when submitted', async () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
      role: 'ADMIN',
    });

    const select = container.querySelector('select[name="userRole"]') as HTMLSelectElement;
    const button = page.getByRole('button', { name: /update/i });

    await userEvent.selectOptions(select, 'MEMBER');
    await userEvent.click(button);

    expect(handleUpdate).toHaveBeenCalledWith('user-1', 'MEMBER');
    expect(handleUpdate).toHaveBeenCalledTimes(1);
  });

  it('should update role value when changed', async () => {
    const toggleUpdate = vi.fn();
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, {
      toggleUpdate,
      handleUpdate,
      userId: 'user-1',
      userEmail: 'user@example.com',
      role: 'ADMIN',
    });

    const select = container.querySelector('select[name="userRole"]') as HTMLSelectElement;

    await userEvent.selectOptions(select, 'MEMBER');

    await expect.element(select).toHaveValue('MEMBER');
  });
});
