import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import { user } from '../../stores';
import ItemComments from './ItemComments.svelte';

const baseUser = {
  id: 'user-1',
  name: 'Alex Doe',
  createdDate: '2024-01-01T00:00:00.000Z',
  updatedDate: '2024-01-01T00:00:00.000Z',
  lastActive: '2024-01-01T00:00:00.000Z',
  locale: 'en',
  rank: 'warrior',
  subscribed: true,
};

const baseUsers = [
  {
    id: 'user-1',
    name: 'Alex Doe',
    avatar: '',
    gravatarHash: '',
    active: true,
  },
  {
    id: 'user-2',
    name: 'Sam Sage',
    avatar: '',
    gravatarHash: '',
    active: true,
  },
];

const baseComments = [
  {
    id: 'comment-1',
    user_id: 'user-1',
    retro_id: 'retro-1',
    comment: 'First comment',
    created_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    updated_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
  },
  {
    id: 'comment-2',
    user_id: 'user-2',
    retro_id: 'retro-1',
    comment: 'Second comment',
    created_date: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
    updated_date: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
  },
];

const setup = (overrides?: {
  item?: {
    id: string;
    comments: typeof baseComments;
    completed: boolean;
    content: string;
    assignees: typeof baseUsers;
    retroId: string;
  };
  users?: typeof baseUsers;
  isFacilitator?: boolean;
  sendSocketEvent?: (event: string, value: any) => void;
}) => {
  const {
    item = {
      id: 'item-1',
      comments: baseComments,
      completed: false,
      content: 'Test Action Item',
      assignees: baseUsers,
      retroId: 'retro-1',
    },
    users = baseUsers,
    isFacilitator = false,
    sendSocketEvent = vi.fn(),
  } = overrides || {};

  user.create({
    ...baseUser,
    id: baseUser.id,
  });

  return {
    sendSocketEvent,
    ...render(ItemComments, {
      item,
      users,
      isFacilitator,
      sendSocketEvent,
    }),
  };
};

describe('ItemComments component', () => {
  it('should render the modal and comments header', () => {
    setup();

    const modal = page.getByRole('dialog');
    expect(modal).toBeTruthy();

    const header = page.getByText(/2\s+comments/i);
    expect(header).toBeTruthy();
  });

  it('should show empty state when there are no comments', () => {
    setup({
      item: {
        id: 'item-1',
        comments: [],
        completed: false,
        content: 'Test Action Item',
        assignees: baseUsers,
        retroId: 'retro-1',
      },
    });

    const emptyState = page.getByText('Be the first to share your thoughts on this retro feedback item.');
    expect(emptyState).toBeTruthy();
  });

  it('should render comment items when comments exist', () => {
    const { container } = setup();

    const commentItems = container.querySelectorAll('[data-commentid]');
    expect(commentItems.length).toBe(2);
  });

  it('should submit comment and send socket event', async () => {
    const sendSocketEvent = vi.fn();
    setup({ sendSocketEvent });

    const textarea = page.getByRole('textbox');
    const submitButton = page.getByRole('button', { name: 'Post Comment' });

    await userEvent.fill(textarea, 'New retro comment');
    await submitButton.click();

    expect(sendSocketEvent).toHaveBeenCalledWith(
      'item_comment_add',
      JSON.stringify({
        item_id: 'item-1',
        comment: 'New retro comment',
      }),
    );
    expect(sendSocketEvent).toHaveBeenCalledTimes(1);
  });

  it('should send edit socket event from comment actions', async () => {
    const sendSocketEvent = vi.fn();
    const { container } = setup({ sendSocketEvent, isFacilitator: true, users: baseUsers });

    const actionsButtons = container.querySelectorAll('button[aria-label="Comment actions"]');
    await (actionsButtons[0] as HTMLButtonElement).click();

    const editButton = page.getByRole('button', { name: /edit/i });
    await editButton.click();

    const textarea = page.getByLabelText(/edit comment/i);
    await userEvent.clear(textarea);
    await userEvent.fill(textarea, 'Updated comment');

    const updateButton = page.getByRole('button', { name: /update comment/i });
    await updateButton.click();

    expect(sendSocketEvent).toHaveBeenCalledWith(
      'item_comment_edit',
      JSON.stringify({
        comment_id: 'comment-1',
        comment: 'Updated comment',
      }),
    );
  });

  it('should send delete socket event from comment actions', async () => {
    const sendSocketEvent = vi.fn();
    const { container } = setup({ sendSocketEvent, isFacilitator: true, users: baseUsers });

    const actionsButtons = container.querySelectorAll('button[aria-label="Comment actions"]');
    await (actionsButtons[0] as HTMLButtonElement).click();

    const deleteButton = page.getByRole('button', { name: /delete/i });
    await deleteButton.click();

    const confirmButton = page.getByTestId('confirm-confirm');
    await confirmButton.click();

    expect(sendSocketEvent).toHaveBeenCalledWith(
      'item_comment_delete',
      JSON.stringify({
        comment_id: 'comment-1',
      }),
    );
  });
});
