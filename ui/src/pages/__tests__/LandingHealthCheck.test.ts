import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import Landing from '../Landing.svelte';

vi.mock('../../config', async importOriginal => {
  const config = await importOriginal<typeof import('../../config')>();

  return {
    ...config,
    AppConfig: { ...config.AppConfig, FeatureHealthCheck: true, SubscriptionsEnabled: true },
  };
});

describe('Landing page health checks', () => {
  it('presents the feature and which plans include it', async () => {
    render(Landing, { xfetch: () => Promise.resolve(new Response()) });

    await expect.element(page.getByRole('heading', { name: 'Squad Health Checks' })).toBeInTheDocument();
    await expect.element(page.getByRole('link', { name: 'Team and Organization plans' })).toBeInTheDocument();
  });
});
