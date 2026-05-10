import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { page, userEvent } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CreateStoryboard from './CreateStoryboard.svelte';
import { user } from '../../stores';

describe('CreateStoryboard component', () => {
  const notifications = {
    success: vi.fn(),
    danger: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  } as any;

  const router = {
    route: vi.fn(),
  };

  beforeEach(() => {
    user.create({
      id: 'user-1',
      name: 'Alex Doe',
      createdDate: '2024-01-01T00:00:00.000Z',
      updatedDate: '2024-01-01T00:00:00.000Z',
      lastActive: '2024-01-01T00:00:00.000Z',
      locale: 'en',
      rank: 'WARRIOR',
      subscribed: true,
    } as any);
  });

  afterEach(() => {
    user.delete();
    vi.clearAllMocks();
  });

  const buildXfetch = () =>
    vi.fn((endpoint: string, config?: { body?: any }) => {
      if (endpoint === '/api/users/user-1/teams?limit=100') {
        return Promise.resolve({
          json: () => Promise.resolve({ data: [{ id: 'team-1', name: 'Platform Team' }] }),
        } as Response);
      }

      if (endpoint === '/api/teams/team-1') {
        return Promise.resolve({
          json: () =>
            Promise.resolve({
              data: {
                team: {
                  subscribed: true,
                  organization_id: 'org-1',
                },
              },
            }),
        } as Response);
      }

      if (endpoint === '/api/teams/team-1/color-legend-templates') {
        return Promise.resolve({
          json: () =>
            Promise.resolve({
              data: [
                {
                  id: 'team-template',
                  name: 'Team Template',
                  description: 'Template description',
                  teamId: 'team-1',
                  colorLegend: [
                    { color: 'gray', legend: 'Inbox' },
                    { color: 'red', legend: 'Urgent' },
                  ],
                },
              ],
            }),
        } as Response);
      }

      if (endpoint === '/api/organizations/org-1/color-legend-templates') {
        return Promise.resolve({
          json: () => Promise.resolve({ data: [] }),
        } as Response);
      }

      if (endpoint === '/api/teams/team-1/users/user-1/storyboards') {
        return Promise.resolve({
          json: () => Promise.resolve({ data: { id: 'storyboard-1' } }),
        } as Response);
      }

      return Promise.reject(new Error(`Unexpected endpoint: ${endpoint} ${JSON.stringify(config || {})}`));
    });

  it('submits selected template color legend when creating a storyboard', async () => {
    const xfetch = buildXfetch();
    const { container } = render(CreateStoryboard, {
      xfetch,
      notifications,
      router,
    });

    await vi.waitFor(() => {
      expect(container.querySelector('select[name="selectedTeam"]')).toBeTruthy();
    });

    await userEvent.fill(
      container.querySelector('input[name="storyboardName"]') as HTMLInputElement,
      'Discovery Board',
    );
    await userEvent.selectOptions(
      container.querySelector('select[name="selectedTeam"]') as HTMLSelectElement,
      'team-1',
    );

    await vi.waitFor(() => {
      expect(page.getByRole('button', { name: 'Select a color legend template...' })).toBeTruthy();
    });

    await userEvent.click(page.getByRole('button', { name: 'Select a color legend template...' }));
    await userEvent.click(page.getByRole('button', { name: /Team: Team Template Template description/i }));

    await userEvent.click(page.getByRole('button', { name: /create storyboard/i }));

    expect(xfetch).toHaveBeenCalledWith('/api/teams/team-1/users/user-1/storyboards', {
      body: {
        storyboardName: 'Discovery Board',
        joinCode: '',
        facilitatorCode: '',
        colorLegend: [
          { color: 'gray', legend: 'Inbox' },
          { color: 'red', legend: 'Urgent' },
        ],
      },
    });
    expect(router.route).toHaveBeenCalledWith('/storyboard/storyboard-1');
  });
});
