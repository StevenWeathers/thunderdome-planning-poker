import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { page } from 'vitest/browser';

import Admin from './Admin.svelte';
import { user } from '../../stores';

vi.mock('../../config', async importOriginal => {
  const config = await importOriginal<typeof import('../../config')>();
  return {
    ...config,
    AppConfig: { ...config.AppConfig, FeatureHealthCheck: true, FeatureRetro: true, FeatureStoryboard: true },
  };
});

describe('Admin dashboard', () => {
  it('shows the health check, template and activity tallies', async () => {
    user.create({ id: 'admin-id', name: 'Admin', rank: 'ADMIN' } as any);
    const stats = {
      unregisteredUserCount: 0,
      registeredUserCount: 0,
      apikeyCount: 0,
      battleCount: 0,
      activeBattleCount: 0,
      activeBattleUserCount: 0,
      planCount: 0,
      organizationCount: 0,
      departmentCount: 0,
      teamCount: 0,
      teamCheckinsCount: 0,
      retroCount: 0,
      activeRetroCount: 0,
      activeRetroUserCount: 0,
      retroItemCount: 0,
      retroActionCount: 0,
      storyboardCount: 0,
      activeStoryboardCount: 0,
      activeStoryboardUserCount: 0,
      storyboardGoalCount: 0,
      storyboardColumnCount: 0,
      storyboardStoryCount: 0,
      storyboardPersonaCount: 0,
      estimationScaleCount: 0,
      publicEstimationScaleCount: 0,
      organizationEstimationScaleCount: 0,
      teamEstimationScaleCount: 0,
      userSubscriptionActiveCount: 0,
      teamSubscriptionActiveCount: 0,
      orgSubscriptionActiveCount: 0,
      retroTemplateCount: 0,
      organizationRetroTemplateCount: 0,
      teamRetroTemplateCount: 0,
      publicRetroTemplateCount: 0,
      colorLegendTemplateCount: 0,
      publicColorLegendTemplateCount: 13,
      organizationColorLegendTemplateCount: 0,
      teamColorLegendTemplateCount: 0,
      personaTemplateCount: 0,
      publicPersonaTemplateCount: 0,
      organizationPersonaTemplateCount: 0,
      teamPersonaTemplateCount: 17,
      retroGroupCount: 19,
      retroTagCount: 0,
      retroItemCommentCount: 0,
      retroActionCommentCount: 0,
      storyboardStoryCommentCount: 23,
      healthCheckTemplateCount: 0,
      publicHealthCheckTemplateCount: 0,
      organizationHealthCheckTemplateCount: 0,
      teamHealthCheckTemplateCount: 29,
      projectCount: 0,
      openSupportTicketCount: 0,
      teamHealthChecksCount: 42,
    };
    const xfetch = vi.fn().mockResolvedValue({ json: () => Promise.resolve({ data: stats }) });

    render(Admin, { xfetch, router: { route: vi.fn() }, notifications: { danger: vi.fn() } as any });

    await expect.element(page.getByText('42', { exact: true })).toBeInTheDocument();
    const labelFor = (count: string) =>
      page.getByText(count, { exact: true }).element().previousElementSibling?.textContent?.trim();
    expect(labelFor('42')).toBe('Health Checks');
    expect(labelFor('13')).toBe('Public Color Legend Templates');
    expect(labelFor('17')).toBe('Team Persona Templates');
    expect(labelFor('19')).toBe('Retro Item Groups');
    expect(labelFor('23')).toBe('Storyboard Story Comments');
    expect(labelFor('29')).toBe('Team Health Check Templates');
  });
});
