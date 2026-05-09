import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import ActionComments from './ActionComments.svelte';
import { user } from '../../stores';

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
  comments: [
    {
      id: 'comment-1',
      user_id: 'user-1',
      comment: 'Needs a fallback user lookup',
      created_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
      updated_date: new Date(Date.now() - 5 * 60 * 1000).toISOString(),
      retro_id: 'retro-1',
    },
  ],
};

describe('ActionComments component', () => {
  it('fetches team users when the parent does not provide them', async () => {
    user.create({
      id: 'current-user',
      name: 'Current User',
      createdDate: '2024-01-01T00:00:00.000Z',
      updatedDate: '2024-01-01T00:00:00.000Z',
      lastActive: '2024-01-01T00:00:00.000Z',
      locale: 'en',
      rank: 'warrior',
      subscribed: true,
    });

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

    render(ActionComments, {
      xfetch,
      notifications: createNotifications(),
      actions: [baseAction],
      selectedActionId: 'action-1',
      toggleComments: vi.fn(),
      getRetrosActions: vi.fn(),
    });

    await expect.element(page.getByText('Alex Doe')).toBeInTheDocument();
    expect(xfetch).toHaveBeenCalledWith('/api/teams/team-1/users?limit=1000&offset=0');
  });

  it('does not fetch team users when the parent already provided them', async () => {
    const xfetch = vi.fn() as unknown as ApiClient;

    render(ActionComments, {
      xfetch,
      notifications: createNotifications(),
      actions: [baseAction],
      users: [
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
      selectedActionId: 'action-1',
      toggleComments: vi.fn(),
      getRetrosActions: vi.fn(),
    });

    await expect.element(page.getByText('Alex Doe')).toBeInTheDocument();
    expect(xfetch).not.toHaveBeenCalled();
  });
});
