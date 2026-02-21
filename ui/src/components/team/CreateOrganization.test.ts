import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CreateOrganization from './CreateOrganization.svelte';

describe('CreateOrganization component', () => {
  it('should render successfully', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    render(CreateOrganization, {
      toggleCreate,
      handleCreate,
    });

    const form = page.getByRole('form', { name: 'createOrganization' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateOrganization, {
      toggleCreate,
      handleCreate,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render organization name input field', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateOrganization, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="organizationName"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('organizationName');
  });

  it('should disable submit button when organization name is empty', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateOrganization, {
      toggleCreate,
      handleCreate,
      organizationName: '',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(true);
  });

  it('should enable submit button when organization name is provided', () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateOrganization, {
      toggleCreate,
      handleCreate,
      organizationName: 'Acme Org',
    });

    const button = container.querySelector('button[type="submit"]') as HTMLButtonElement;
    expect(button).toBeTruthy();
    expect(button.disabled).toBe(false);
  });

  it('should fire handleCreate with organization name when form is submitted', async () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateOrganization, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="organizationName"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.fill(input, 'Product');
    await button.click();

    expect(handleCreate).toHaveBeenCalledWith('Product');
    expect(handleCreate).toHaveBeenCalledTimes(1);
  });

  it('should update organization name value when input changes', async () => {
    const toggleCreate = vi.fn();
    const handleCreate = vi.fn();
    const { container } = render(CreateOrganization, {
      toggleCreate,
      handleCreate,
    });

    const input = container.querySelector('input[name="organizationName"]') as HTMLInputElement;

    await userEvent.fill(input, 'Operations');

    await expect.element(input).toHaveValue('Operations');
  });
});
