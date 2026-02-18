import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import DeleteConfirmation from './DeleteConfirmation.svelte';

describe('DeleteConfirmation component', () => {
  it('should render successfully', () => {
    render(DeleteConfirmation, {});
  });

  it('should fire handleDelete when confirmed', async () => {
    const stub = vi.fn();
    render(DeleteConfirmation, {
      handleDelete: stub,
    });
    const button = page.getByRole('button', { name: 'Confirm Delete' });

    await button.click();

    expect(stub).toHaveBeenCalled();
  });

  it('should not fire handleDelete when cancel and instead fire toggleDelete', async () => {
    const handleDelete = vi.fn();
    const toggleDelete = vi.fn();
    render(DeleteConfirmation, {
      handleDelete,
      toggleDelete,
    });
    const button = page.getByRole('button', { name: 'Cancel' });

    await button.click();

    expect(handleDelete).not.toHaveBeenCalled();
    expect(toggleDelete).toHaveBeenCalled();
  });
});
