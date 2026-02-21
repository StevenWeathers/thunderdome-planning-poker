import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import { user } from '../../stores';
import Comment from './Comment.svelte';

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

const baseComment = {
  id: 'comment-1',
  user_id: 'user-1',
  comment: 'Hello there',
  created_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
  updated_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
};

const baseUserMap = new Map([
  ['user-1', { id: 'user-1', name: 'Alex Doe' }],
  ['user-2', { id: 'user-2', name: 'Sam Sage' }],
]);

const setup = (overrides?: {
  comment?: typeof baseComment;
  userId?: string;
  isAdmin?: boolean;
  handleEdit?: (commentId: string, data: { userId: string; comment: string }) => void;
  handleDelete?: (commentId: string) => void;
}) => {
  const {
    comment = baseComment,
    userId = 'user-1',
    isAdmin = false,
    handleEdit = vi.fn(),
    handleDelete = vi.fn(),
  } = overrides || {};

  user.create({
    ...baseUser,
    id: userId,
  });

  return {
    handleEdit,
    handleDelete,
    ...render(Comment, {
      comment,
      userMap: baseUserMap,
      isAdmin,
      handleEdit,
      handleDelete,
    }),
  };
};

describe('Comment component', () => {
  it('should render comment content with user name and timestamp', () => {
    setup();

    const article = page.getByRole('article', { name: /comment by alex doe/i });
    expect(article).toBeTruthy();

    const text = page.getByText('Hello there');
    expect(text).toBeTruthy();

    const time = page.getByText(/m ago/i);
    expect(time).toBeTruthy();
  });

  it('should not show actions when user cannot edit', () => {
    const { container } = setup({
      userId: 'user-2',
      comment: {
        ...baseComment,
        user_id: 'user-1',
      },
      isAdmin: false,
    });

    const actionsButton = container.querySelector('button[aria-label="Comment actions"]');
    expect(actionsButton).toBeNull();
  });

  it('should open edit mode and submit updates', async () => {
    const handleEdit = vi.fn();
    setup({ handleEdit });

    const actionsButton = page.getByRole('button', { name: 'Comment actions' });
    await actionsButton.click();

    const editButton = page.getByRole('button', { name: /edit/i });
    await editButton.click();

    const form = page.getByRole('form', { name: 'editComment' });
    expect(form).toBeTruthy();

    const textarea = page.getByLabelText(/edit comment/i);
    await userEvent.clear(textarea);
    await userEvent.fill(textarea, 'Updated comment text');

    const submitButton = page.getByRole('button', { name: /update comment/i });
    await submitButton.click();

    expect(handleEdit).toHaveBeenCalledWith('comment-1', {
      userId: 'user-1',
      comment: 'Updated comment text',
    });
    expect(handleEdit).toHaveBeenCalledTimes(1);
  });

  it('should disable submit when comment is unchanged', async () => {
    setup();

    const actionsButton = page.getByRole('button', { name: 'Comment actions' });
    await actionsButton.click();

    const editButton = page.getByRole('button', { name: /edit/i });
    await editButton.click();

    const submitButton = page.getByRole('button', { name: /update comment/i }).element() as HTMLButtonElement;
    expect(submitButton.disabled).toBe(true);
  });

  it('should confirm delete and fire handleDelete', async () => {
    const handleDelete = vi.fn();
    setup({ handleDelete });

    const actionsButton = page.getByRole('button', { name: 'Comment actions' });
    await actionsButton.click();

    const deleteButton = page.getByRole('button', { name: /delete/i });
    await deleteButton.click();

    const confirmButton = page.getByTestId('confirm-confirm');
    await confirmButton.click();

    expect(handleDelete).toHaveBeenCalledWith('comment-1');
    expect(handleDelete).toHaveBeenCalledTimes(1);
  });
});
