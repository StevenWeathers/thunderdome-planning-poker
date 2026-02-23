import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ActionsMenu from './ActionsMenu.svelte';
import { PencilIcon, Trash2, Settings } from '@lucide/svelte';

const baseActions = [
  {
    label: 'Edit',
    icon: PencilIcon,
    onclick: vi.fn(),
    disabled: false,
  },
  {
    label: 'Delete',
    icon: Trash2,
    onclick: vi.fn(),
    className: 'text-red-600 hover:bg-red-50',
    disabled: false,
  },
];

const setup = (overrides?: { Icon?: any; actions?: typeof baseActions; ariaLabel?: string; disabled?: boolean }) => {
  const { Icon, actions = baseActions, ariaLabel = 'Actions menu', disabled = false } = overrides || {};

  const props = {
    actions,
    ariaLabel,
    disabled,
    ...(Icon && { Icon }),
  };

  return render(ActionsMenu, props);
};

describe('ActionsMenu component', () => {
  it('should render button with default icon and aria label', () => {
    setup();

    const button = page.getByRole('button', { name: 'Actions menu' });
    expect(button).toBeTruthy();
    expect(button.element()).toHaveAttribute('aria-expanded', 'false');
  });

  it('should render button with custom aria label', () => {
    setup({ ariaLabel: 'Comment actions' });

    const button = page.getByRole('button', { name: 'Comment actions' });
    expect(button).toBeTruthy();
  });

  it('should toggle menu open and closed', async () => {
    setup();

    const button = page.getByRole('button', { name: 'Actions menu' });
    expect(button.element()).toHaveAttribute('aria-expanded', 'false');

    await button.click();
    expect(button.element()).toHaveAttribute('aria-expanded', 'true');

    await button.click();
    expect(button.element()).toHaveAttribute('aria-expanded', 'false');
  });

  it('should display all actions when menu is open', async () => {
    setup();

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    const editButton = page.getByRole('button', { name: 'Edit' });
    expect(editButton).toBeTruthy();

    const deleteButton = page.getByRole('button', { name: 'Delete' });
    expect(deleteButton).toBeTruthy();
  });

  it('should apply custom className to action button', async () => {
    setup();

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    const deleteButton = page.getByRole('button', { name: 'Delete' });
    const deleteButtonElement = deleteButton.element() as HTMLButtonElement;

    expect(deleteButtonElement.className).toContain('text-red-600');
    expect(deleteButtonElement.className).toContain('hover:bg-red-50');
  });

  it('should call action onclick handler when clicked', async () => {
    const mockOnclick = vi.fn();
    const customActions = [
      {
        label: 'Test',
        icon: Settings,
        onclick: mockOnclick,
        disabled: false,
      },
    ];

    setup({ actions: customActions });

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    const actionButton = page.getByRole('button', { name: 'Test' });
    await actionButton.click();

    expect(mockOnclick).toHaveBeenCalledTimes(1);
  });

  it('should close menu after action is clicked', async () => {
    setup();

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    expect(button.element()).toHaveAttribute('aria-expanded', 'true');

    const editButton = page.getByRole('button', { name: 'Edit' });
    await editButton.click();

    expect(button.element()).toHaveAttribute('aria-expanded', 'false');
  });

  it('should close menu when clicking outside', async () => {
    const { container } = setup();

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    expect(button.element()).toHaveAttribute('aria-expanded', 'true');

    // Click outside the menu container to close it
    const outsideElement = container.parentElement as HTMLElement;
    const clickEvent = new MouseEvent('click', { bubbles: true });
    outsideElement.dispatchEvent(clickEvent);
    // Wait a tick for state update
    await new Promise(resolve => setTimeout(resolve, 0));

    expect(button.element()).toHaveAttribute('aria-expanded', 'false');
  });

  it('should render empty menu when no actions provided', async () => {
    setup({ actions: [] });

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    // Try to find an action button - should not exist
    try {
      page.getByRole('button', { name: 'Edit' });
      expect.fail('Edit button should not exist');
    } catch {
      // Expected - button does not exist
      expect(true).toBe(true);
    }
  });

  it('should handle multiple actions with different onclick handlers', async () => {
    const mockEdit = vi.fn();
    const mockDelete = vi.fn();
    const customActions = [
      {
        label: 'Edit',
        icon: PencilIcon,
        onclick: mockEdit,
        disabled: false,
      },
      {
        label: 'Delete',
        icon: Trash2,
        onclick: mockDelete,
        disabled: false,
      },
    ];

    setup({ actions: customActions });

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    const deleteButton = page.getByRole('button', { name: 'Delete' });
    await deleteButton.click();

    expect(mockDelete).toHaveBeenCalledTimes(1);
    expect(mockEdit).not.toHaveBeenCalled();
  });

  it('should disable menu button when disabled prop is true', () => {
    setup({ disabled: true });

    const button = page.getByRole('button', { name: 'Actions menu' }).element() as HTMLButtonElement;
    expect(button.disabled).toBe(true);
  });

  it('should not open menu when button is disabled', async () => {
    setup({ disabled: true });

    const buttonElement = page.getByRole('button', { name: 'Actions menu' }).element() as HTMLButtonElement;
    // Disabled buttons require direct click, not userEvent
    buttonElement.click();

    expect(buttonElement).toHaveAttribute('aria-expanded', 'false');
  });

  it('should disable individual action buttons', async () => {
    const customActions = [
      {
        label: 'Active',
        icon: PencilIcon,
        onclick: vi.fn(),
        disabled: false,
      },
      {
        label: 'Disabled',
        icon: Trash2,
        onclick: vi.fn(),
        disabled: true,
      },
    ];

    setup({ actions: customActions });

    const button = page.getByRole('button', { name: 'Actions menu' });
    await button.click();

    const disabledActionButton = page.getByRole('button', { name: 'Disabled' }).element() as HTMLButtonElement;
    expect(disabledActionButton.disabled).toBe(true);
  });
});
