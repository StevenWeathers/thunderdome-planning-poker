import { describe, it, expect, vi, beforeEach } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import TeamHealthCheck from './TeamHealthCheck.svelte';
import { user } from '../../stores';
import type { HealthCheckDetail } from '../../types/healthcheck';

const detail = (overrides: Partial<HealthCheckDetail> = {}, closed = false): HealthCheckDetail => ({
  healthCheck: {
    id: 'check-1',
    teamId: 'team-1',
    name: 'Sprint 12',
    areas: [{ name: 'Fun', healthy: '', unhealthy: '' }],
    deadline: '2030-01-01T00:00:00Z',
    closed,
    createdBy: 'user-1',
    responseCount: 3,
    createdAt: '2029-01-01T00:00:00Z',
    updatedAt: '2029-01-01T00:00:00Z',
  },
  canRespond: true,
  canManage: false,
  minResponses: 3,
  myAnswers: [],
  pending: null,
  results: null,
  ...overrides,
});

function renderPage(data: HealthCheckDetail) {
  const xfetch = vi.fn().mockResolvedValue(new Response(JSON.stringify({ data })));
  render(TeamHealthCheck, {
    xfetch,
    router: { route: vi.fn() },
    notifications: { success: vi.fn(), danger: vi.fn() } as any,
    teamId: 'team-1',
    healthCheckId: 'check-1',
  });

  return xfetch;
}

describe('TeamHealthCheck page', () => {
  beforeEach(() => {
    user.create({ id: 'user-1', name: 'Thor', rank: 'REGISTERED' } as any);
  });

  it('lets a team member respond while open, without management actions', async () => {
    const xfetch = renderPage(detail());

    await expect.element(page.getByRole('button', { name: 'Submit Response' })).toBeInTheDocument();
    expect(page.getByTestId('health-check-close').query()).toBeNull();
    expect(page.getByTestId('health-check-delete').query()).toBeNull();
    expect(xfetch).toHaveBeenCalledWith('/api/teams/team-1/health-checks/check-1');
  });

  it('shows managers who is yet to respond', async () => {
    renderPage(detail({ canManage: true, pending: [{ id: 'user-2', name: 'Loki' }] }));

    await expect.element(page.getByTestId('health-check-pending')).toHaveTextContent('Loki');
    await expect.element(page.getByTestId('health-check-close')).toBeInTheDocument();
  });

  it('does not offer the form to admins outside the team', async () => {
    renderPage(detail({ canRespond: false }));

    await expect.element(page.getByText('only team members can respond', { exact: false })).toBeInTheDocument();
    expect(page.getByRole('button', { name: 'Submit Response' }).query()).toBeNull();
  });

  it('shows anonymous results once closed', async () => {
    renderPage(
      detail(
        {
          results: [
            {
              name: 'Fun',
              healthy: '',
              unhealthy: '',
              ratings: { green: 2, yellow: 1, red: 0 },
              trends: { improving: 1, stable: 2, declining: 0 },
              comments: ['More pairing'],
            },
          ],
        },
        true,
      ),
    );

    await expect.element(page.getByTestId('health-check-rating-green')).toHaveTextContent('Good: 2');
    await expect.element(page.getByText('More pairing')).toBeInTheDocument();
    expect(page.getByRole('button', { name: 'Submit Response' }).query()).toBeNull();
  });

  it('explains why results are hidden when too few responded', async () => {
    renderPage(detail({}, true));

    await expect.element(page.getByText('fewer than 3 people responded', { exact: false })).toBeInTheDocument();
  });
});
