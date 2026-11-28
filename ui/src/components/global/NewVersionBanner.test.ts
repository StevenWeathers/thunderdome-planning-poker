import { afterEach, describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { page } from 'vitest/browser';

import NewVersionBanner from './NewVersionBanner.svelte';
import { newVersion } from '../../stores';

describe('NewVersionBanner', () => {
  afterEach(() => newVersion.set(''));

  it('is hidden while the tab is up to date', async () => {
    render(NewVersionBanner);

    expect(page.getByText('A new version is available.').query()).toBeNull();
    expect(page.getByRole('button', { name: 'Reload' }).query()).toBeNull();
  });

  it('offers a reload once a new version is known', async () => {
    const reload = vi.fn();
    render(NewVersionBanner, { reload });

    newVersion.set('v2');

    await expect.element(page.getByText('A new version is available.')).toBeVisible();
    await page.getByRole('button', { name: 'Reload' }).click();
    expect(reload).toHaveBeenCalledTimes(1);
  });
});
