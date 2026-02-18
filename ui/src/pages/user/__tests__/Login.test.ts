import { describe, it } from 'vitest';
import { render } from 'vitest-browser-svelte';

import Login from '../Login.svelte';

describe('Login Page', () => {
  it('should render successfully', () => {
    render(Login, {
      xfetch: () => {},
      notifications: () => {},
      router: () => {},
      battleId: null,
    });
  });
});
