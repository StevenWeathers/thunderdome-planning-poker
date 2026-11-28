import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import { page } from 'vitest/browser';

import LoginForm from './LoginForm.svelte';

vi.mock('../../config', async importOriginal => {
  const original = await importOriginal<typeof import('../../config')>();

  return { ...original, appRoutes: { ...original.appRoutes, login: '/thunderdome/login' } };
});

describe('LoginForm', () => {
  it('keeps the path prefix in the forgot password link', async () => {
    render(LoginForm, { router: { route: vi.fn() } });

    await expect
      .element(page.getByRole('link', { name: 'Forgot Your Password?' }))
      .toHaveAttribute('href', '/thunderdome/login');
  });
});
