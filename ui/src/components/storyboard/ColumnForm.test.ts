import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ColumnForm from './ColumnForm.svelte';

describe('ColumnForm component', () => {
  const baseColumn = {
    id: 'column-1',
    name: 'To Do',
    personas: [],
  };

  it('should render successfully', () => {
    render(ColumnForm, {
      column: baseColumn,
    });

    const form = page.getByRole('form', { name: 'addColumn' });
    expect(form).toBeTruthy();
  });

  it('should render the modal dialog', () => {
    const { container } = render(ColumnForm, {
      column: baseColumn,
    });

    const modal = container.querySelector('[role="dialog"]');
    expect(modal).toBeTruthy();
  });

  it('should render column name input field', () => {
    const { container } = render(ColumnForm, {
      column: baseColumn,
    });

    const input = container.querySelector('input[name="columnName"]');
    expect(input).toBeTruthy();
    expect(input?.getAttribute('id')).toBe('columnName');
  });

  it('should submit column updates and close the modal', async () => {
    const handleColumnRevision = vi.fn();
    const toggleColumnEdit = vi.fn();
    const { container } = render(ColumnForm, {
      column: { ...baseColumn },
      handleColumnRevision,
      toggleColumnEdit,
    });

    const input = container.querySelector('input[name="columnName"]') as HTMLInputElement;
    const button = page.getByRole('button', { name: /save/i });

    await userEvent.clear(input);
    await userEvent.fill(input, 'In Progress');
    await userEvent.click(button);

    expect(handleColumnRevision).toHaveBeenCalledWith({
      id: 'column-1',
      name: 'In Progress',
    });
    expect(handleColumnRevision).toHaveBeenCalledTimes(1);
    expect(toggleColumnEdit).toHaveBeenCalledTimes(1);
  });
});
