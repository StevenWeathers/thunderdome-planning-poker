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

    expect(handleUpdate).toHaveBeenCalledWith('user-1', 'MEMBER', { title: '', notes: '' });
    expect(handleUpdate).toHaveBeenCalledTimes(1);
  });

  const titleProps = {
    userId: 'user-1',
    userEmail: 'user@example.com',
    role: 'MEMBER',
    titleOptions: ['Backend Developer', 'Product Manager'],
  };

  it('should only offer a text input when there are no title options', () => {
    const { container } = render(UpdateUser, { ...titleProps, titleOptions: [] });

    expect(container.querySelector('select[name="userTitleChoice"]')).toBeNull();
    expect(container.querySelector('input[name="userTitle"]')).toBeTruthy();
  });

  it('should submit a title picked from the options', async () => {
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, { ...titleProps, handleUpdate, notes: 'Knows billing' });
    const select = container.querySelector('select[name="userTitleChoice"]') as HTMLSelectElement;

    expect(container.querySelector('input[name="userTitle"]')).toBeNull();
    await userEvent.selectOptions(select, 'Product Manager');
    await userEvent.click(page.getByRole('button', { name: /update/i }));

    expect(handleUpdate).toHaveBeenCalledWith('user-1', 'MEMBER', {
      title: 'Product Manager',
      notes: 'Knows billing',
    });
  });

  it('should allow a custom title when there are options', async () => {
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, { ...titleProps, handleUpdate });
    const select = container.querySelector('select[name="userTitleChoice"]') as HTMLSelectElement;

    await userEvent.selectOptions(select, '__custom__');
    await userEvent.fill(container.querySelector('input[name="userTitle"]') as HTMLInputElement, 'Release Captain');
    await userEvent.click(page.getByRole('button', { name: /update/i }));

    expect(handleUpdate).toHaveBeenCalledWith('user-1', 'MEMBER', { title: 'Release Captain', notes: '' });
  });

  it('should show an existing custom title as custom', async () => {
    const { container } = render(UpdateUser, { ...titleProps, title: 'Release Captain' });

    const select = container.querySelector('select[name="userTitleChoice"]') as HTMLSelectElement;
    const input = container.querySelector('input[name="userTitle"]') as HTMLInputElement;
    await expect.element(select).toHaveValue('__custom__');
    expect(input.value).toBe('Release Captain');
  });

  it('should leave the title empty when it is inherited', async () => {
    const handleUpdate = vi.fn();
    const { container } = render(UpdateUser, { ...titleProps, handleUpdate, inheritedTitle: 'Designer' });

    const emptyOption = container.querySelector('select[name="userTitleChoice"] option[value=""]');
    expect(emptyOption?.textContent).toContain('Designer');
    await userEvent.click(page.getByRole('button', { name: /update/i }));

    expect(handleUpdate).toHaveBeenCalledWith('user-1', 'MEMBER', { title: '', notes: '' });
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
