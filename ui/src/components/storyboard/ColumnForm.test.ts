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

  const personas = [
    { id: 'persona-1', name: 'Sam', role: 'Designer' },
    { id: 'persona-2', name: 'Alex', role: 'Developer' },
  ];

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

  it('should call delete column handler when delete is clicked', async () => {
    const deleteHandler = vi.fn();
    const deleteColumn = vi.fn(() => deleteHandler);

    render(ColumnForm, {
      column: baseColumn,
      deleteColumn,
    });

    const button = page.getByRole('button', { name: /delete column/i });
    await userEvent.click(button);

    expect(deleteColumn).toHaveBeenCalledWith('column-1');
    expect(deleteHandler).toHaveBeenCalledTimes(1);
  });

  it('should disable add persona button when no persona is selected', () => {
    render(ColumnForm, {
      column: baseColumn,
      personas,
    });

    const button = page.getByRole('button', { name: /add persona/i });
    expect((button.element() as HTMLButtonElement).disabled).toBe(true);
  });

  it('should add a persona when selected', async () => {
    const handlePersonaAdd = vi.fn();
    const { container } = render(ColumnForm, {
      column: baseColumn,
      personas,
      handlePersonaAdd,
    });

    const select = container.querySelector('select[name="persona"]') as HTMLSelectElement;
    const button = page.getByRole('button', { name: /add persona/i });

    await userEvent.selectOptions(select, 'persona-2');
    await userEvent.click(button);

    expect(handlePersonaAdd).toHaveBeenCalledWith({
      column_id: 'column-1',
      persona_id: 'persona-2',
    });
    expect(handlePersonaAdd).toHaveBeenCalledTimes(1);
  });

  it('should render personas list when column has personas', () => {
    const { container } = render(ColumnForm, {
      column: {
        ...baseColumn,
        personas: [{ id: 'persona-1', name: 'Sam', role: 'Designer' }],
      },
      personas,
    });

    const personaText = container.querySelector('div.grid')?.textContent || '';
    expect(personaText).toContain('Sam (Designer)');
  });

  it('should call persona remove handler when remove is clicked', async () => {
    const removeHandler = vi.fn();
    const handlePersonaRemove = vi.fn(() => removeHandler);

    const { container } = render(ColumnForm, {
      column: {
        ...baseColumn,
        personas: [{ id: 'persona-1', name: 'Sam', role: 'Designer' }],
      },
      personas,
      handlePersonaRemove,
    });

    const removeButton = container.querySelector('div.grid button') as HTMLButtonElement;
    await userEvent.click(removeButton);

    expect(handlePersonaRemove).toHaveBeenCalledWith({
      column_id: 'column-1',
      persona_id: 'persona-1',
    });
    expect(removeHandler).toHaveBeenCalledTimes(1);
  });
});
