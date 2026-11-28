import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import TeamHealthChecks from './TeamHealthChecks.svelte';
import { user } from '../../stores';

describe('TeamHealthChecks page', () => {
  it('explains what a health check is and links to the article in a new tab', async () => {
    user.create({ id: 'user-1', name: 'Thor', rank: 'REGISTERED' } as any);
    const xfetch = vi
      .fn()
      .mockImplementation((url: string) =>
        Promise.resolve(
          new Response(
            JSON.stringify(
              url.includes('health-checks')
                ? { data: [], meta: { count: 0 } }
                : { data: { team: { name: 'Squad' }, teamRole: 'MEMBER' } },
            ),
          ),
        ),
      );

    render(TeamHealthChecks, {
      xfetch,
      router: { route: vi.fn() },
      notifications: { danger: vi.fn() } as any,
      teamId: 'team-1',
    });

    await expect.element(page.getByRole('heading', { name: 'What is a Squad Health Check?' })).toBeInTheDocument();
    await expect.element(page.getByText('Quarterly is a good place to start.')).toBeInTheDocument();

    const link = page.getByRole('link', { name: "Read Spotify's article about the model (opens in a new tab)" });
    await expect
      .element(link)
      .toHaveAttribute('href', 'https://engineering.atspotify.com/2014/9/squad-health-check-model');
    await expect.element(link).toHaveAttribute('target', '_blank');
    await expect.element(link).toHaveAttribute('rel', 'noopener noreferrer');
  });
});
