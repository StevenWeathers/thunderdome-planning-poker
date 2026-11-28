import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import Pricing from '../subscription/Pricing.svelte';

vi.mock('../../config', async importOriginal => {
  const config = await importOriginal<typeof import('../../config')>();
  const plan = { Enabled: true, MonthPrice: '1', YearPrice: '10', MonthCheckoutLink: '', YearCheckoutLink: '' };

  return {
    ...config,
    AppConfig: {
      ...config.AppConfig,
      FeatureHealthCheck: true,
      Subscription: { Individual: plan, Team: plan, Organization: plan },
    },
  };
});

describe('Pricing page', () => {
  it('lists the premium features', async () => {
    render(Pricing);

    for (const feature of [
      'Custom estimation scales for Planning Poker and Story Mapping',
      'Custom tags for retrospective brainstorm items',
      'Custom story map persona and color legend templates',
      'Squad Health Checks with anonymous responses and custom templates',
      'Share custom templates, tags and estimation scales across all your teams',
    ]) {
      await expect.element(page.getByText(feature)).toBeInTheDocument();
    }
  });
});
