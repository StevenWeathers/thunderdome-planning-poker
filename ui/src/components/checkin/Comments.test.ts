import { describe, it, expect, vi } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import { user } from '../../stores';
import Comments from './Comments.svelte';

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

const baseComments = [
  {
    id: 'comment-1',
    user_id: 'user-1',
    comment: 'First comment',
    created_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
    updated_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
  },
  {
    id: 'comment-2',
    user_id: 'user-2',
    comment: 'Second comment',
    created_date: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
    updated_date: new Date(Date.now() - 2 * 60 * 1000).toISOString(),
  },
];

const baseUserMap = new Map([
  ['user-1', { id: 'user-1', name: 'Alex Doe' }],
  ['user-2', { id: 'user-2', name: 'Sam Sage' }],
]);

const setup = (overrides?: {
  checkin?: { id: string; comments: typeof baseComments };
  userId?: string;
  isAdmin?: boolean;
  handleCreate?: (checkinId: string, data: { userId: string; comment: string }) => void;
  handleEdit?: (checkinId: string, commentId: string, data: { userId: string; comment: string }) => void;
  handleDelete?: (checkinId: string, commentId: string) => void;
}) => {
  const {
    checkin = { id: 'checkin-1', comments: baseComments },
    userId = 'user-1',
    isAdmin = false,
    handleCreate = vi.fn(),
    handleEdit = vi.fn(),
    handleDelete = vi.fn(),
  } = overrides || {};

  user.create({
    ...baseUser,
    id: userId,
  });

  return {
    handleCreate,
    handleEdit,
    handleDelete,
    ...render(Comments, {
      checkin,
      userMap: baseUserMap,
      isAdmin,
      handleCreate,
      handleEdit,
      handleDelete,
    }),
  };
};

describe('Comments component', () => {
  const getToggleButton = (container: HTMLElement) =>
    container.querySelector('button[aria-controls="comments-section"]') as HTMLButtonElement;

  it('should render the toggle button with comment count', () => {
    const { container } = setup({ checkin: { id: 'checkin-1', comments: [baseComments[0]] } });

    const toggleButton = getToggleButton(container);
    expect(toggleButton).toBeTruthy();
    expect(toggleButton.textContent).toMatch(/1\s+comment/i);
  });

  it('should toggle the comments section', async () => {
    const { container } = setup({ checkin: { id: 'checkin-1', comments: [baseComments[0]] } });

    const toggleButton = getToggleButton(container);
    expect(container.querySelector('[role="region"]')).toBeNull();

    await toggleButton.click();

    const region = page.getByRole('region', { name: 'Comments section' });
    expect(region).toBeTruthy();
    expect(page.getByText('Hide conversation')).toBeTruthy();

    await toggleButton.click();

    expect(container.querySelector('[role="region"]')).toBeNull();
    expect(page.getByText('View conversation')).toBeTruthy();
  });

  it('should show empty state when there are no comments', async () => {
    const { container } = setup({ checkin: { id: 'checkin-1', comments: [] } });

    const toggleButton = getToggleButton(container);
    await toggleButton.click();

    const emptyState = page.getByText('Be the first to share your thoughts on this check-in.');
    expect(emptyState).toBeTruthy();
  });

  it('should render comment items when comments exist', async () => {
    const { container } = setup();

    const toggleButton = getToggleButton(container);
    await toggleButton.click();

    const commentItems = container.querySelectorAll('[data-commentid]');
    expect(commentItems.length).toBe(2);
  });

  it('should submit comment through CommentForm', async () => {
    const handleCreate = vi.fn();
    const { container } = setup({ handleCreate });

    const toggleButton = getToggleButton(container);
    await toggleButton.click();

    const textarea = page.getByRole('textbox');
    const submitButton = page.getByRole('button', { name: 'Post Comment' });

    await userEvent.fill(textarea, 'New check-in comment');
    await submitButton.click();

    expect(handleCreate).toHaveBeenCalledWith('checkin-1', {
      userId: 'user-1',
      comment: 'New check-in comment',
    });
    expect(handleCreate).toHaveBeenCalledTimes(1);
  });
});
