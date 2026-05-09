import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ActionItemEdit from './ActionItemEdit.svelte';

const baseAction = {
  id: 'action-1',
  retroId: 'retro-1',
  content: 'Ship the modal fix',
  completed: false,
  assignees: [],
  comments: [],
};

describe('ActionItemEdit component', () => {
  it('renders the current action content and completion state', async () => {
    render(ActionItemEdit, {
      toggleEdit: vi.fn(),
      handleEdit: vi.fn(),
      handleDelete: vi.fn(),
      action: baseAction,
    });

    await expect.element(page.getByRole('textbox')).toHaveValue('Ship the modal fix');
    await expect.element(page.getByLabelText('Completed')).not.toBeChecked();
  });

  it('submits updated content and completion state', async () => {
    const handleEdit = vi.fn();

    render(ActionItemEdit, {
      toggleEdit: vi.fn(),
      handleEdit,
      handleDelete: vi.fn(),
      action: baseAction,
    });

    const textbox = page.getByRole('textbox');
    const completed = page.getByLabelText('Completed');
    const saveButton = page.getByRole('button', { name: 'Save' });

    await userEvent.clear(textbox);
    await userEvent.fill(textbox, 'Updated action');
    await completed.click();
    await saveButton.click();

    expect(handleEdit).toHaveBeenCalledWith({
      id: 'action-1',
      retroId: 'retro-1',
      content: 'Updated action',
      completed: true,
    });
  });
});
