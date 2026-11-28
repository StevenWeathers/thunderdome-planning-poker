import { afterEach, describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { page } from 'vitest/browser';

import OpenSource from '../OpenSource.svelte';
import PrivacyPolicy from '../PrivacyPolicy.svelte';
import TermsConditions from '../TermsConditions.svelte';

vi.mock('../../../config', async importOriginal => {
  const original = await importOriginal<typeof import('../../../config')>();

  return { ...original, AppConfig: { ...original.AppConfig, PathPrefix: '/thunderdome' } };
});

describe('static content pages under a path prefix', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it.each([
    ['PrivacyPolicy', PrivacyPolicy, 'privacy-policy'],
    ['TermsConditions', TermsConditions, 'terms-and-conditions'],
    ['OpenSource', OpenSource, 'open-source'],
  ])('%s requests its content below the prefix', async (_name, component, file) => {
    const fetchSpy = vi.spyOn(window, 'fetch').mockResolvedValue(new Response('<p>Loaded content</p>'));

    render(component);

    await expect.element(page.getByText('Loaded content')).toBeVisible();
    expect(fetchSpy).toHaveBeenCalledWith(`/thunderdome/content/${file}.html`);
  });
});
