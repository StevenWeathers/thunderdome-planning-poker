import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CreateDepartment from './CreateDepartment.svelte';

describe('CreateDepartment component', () => {
  it('should render successfully', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    render(CreateDepartment, {
      toggleCreate,
      handleCreate,
    });

    const form = page.getByRole('form', { name: 'createDepartment' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateDepartment, {
      toggleCreate,
      handleCreate,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render department name input field', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateDepartment, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="departmentName"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('departmentName');
  });

  it('should disable submit button when department name is empty', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateDepartment, {
      toggleCreate,
      handleCreate,
      departmentName: '',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(true);
  });

  it('should enable submit button when department name is provided', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateDepartment, {
      toggleCreate,
      handleCreate,
      departmentName: 'Engineering',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(false);
  });

  it('should fire handleCreate with department name when form is submitted', async () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateDepartment, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="departmentName"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.fill(input, 'Marketing');
    await button.click();

    expect(handleCreate).toHaveBeenCalledWith('Marketing');
    expect(handleCreate).toHaveBeenCalledTimes(1);
  });

  it('should update department name value when input changes', async () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateDepartment, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="departmentName"]') as HTMLInputElement;

    await userEvent.fill(input, 'Operations');

    await expect.element(input).toHaveValue('Operations');
  });
});
