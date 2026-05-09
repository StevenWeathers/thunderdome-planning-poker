import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ActionItemAssignees from './ActionItemAssignees.svelte';

import type { ApiClient } from '../../types/apiclient';
import type { NotificationService } from '../../types/notifications';

function createNotifications(): NotificationService {
  return {
    show: vi.fn(),
    success: vi.fn(),
    danger: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
    removeToast: vi.fn(),
  };
}

function createJsonResponse(data: unknown): Response {
  return new Response(JSON.stringify(data), {
    status: 200,
    headers: {
      'Content-Type': 'application/json',
    },
  });
}

const baseAction = {
  id: 'action-1',
  retroId: 'retro-1',
  teamId: 'team-1',
  content: 'Ship the modal fix',
  completed: false,
  assignees: [],
  comments: [],
};

describe('ActionItemAssignees component', () => {
  it('fetches team users when assignable users are not provided', async () => {
    const xfetch = vi.fn((url: string) => {
      if (url === '/api/teams/team-1/users?limit=1000&offset=0') {
        return Promise.resolve(
          createJsonResponse({
            data: [
              {
                id: 'user-1',
                name: 'Alex Doe',
                avatar: 'warrior',
                gravatarHash: 'hash',
                email: 'alex@example.com',
                role: 'MEMBER',
                pictureUrl: '',
              },
            ],
          }),
        );
      }

      return Promise.resolve(createJsonResponse({ data: {} }));
    }) as unknown as ApiClient;

    render(ActionItemAssignees, {
      xfetch,
      notifications: createNotifications(),
      toggleAssignees: vi.fn(),
      handleAssigneeAdd: vi.fn(),
      handleAssigneeRemove: vi.fn(() => vi.fn()),
      action: baseAction,
      assignableUsers: [],
    });

    await expect.element(page.getByRole('option', { name: 'Alex Doe' })).toBeInTheDocument();
    expect(xfetch).toHaveBeenCalledWith('/api/teams/team-1/users?limit=1000&offset=0');
  });

  it('does not fetch team users when assignable users are already provided', async () => {
    const xfetch = vi.fn() as unknown as ApiClient;

    render(ActionItemAssignees, {
      xfetch,
      notifications: createNotifications(),
      toggleAssignees: vi.fn(),
      handleAssigneeAdd: vi.fn(),
      handleAssigneeRemove: vi.fn(() => vi.fn()),
      action: baseAction,
      assignableUsers: [
        {
          id: 'user-1',
          name: 'Alex Doe',
        },
      ],
    });

    await expect.element(page.getByRole('option', { name: 'Alex Doe' })).toBeInTheDocument();
    expect(xfetch).not.toHaveBeenCalled();
  });
});
