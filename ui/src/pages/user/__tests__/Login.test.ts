import { describe, it } from 'vitest';
import { render } from 'vitest-browser-svelte';

import Login from '../Login.svelte';

describe('Login Page', () => {
  it('should render successfully', () => {
    render(Login, {
      xfetch: () => Promise.resolve(new Response()),
      notifications: {
        show: () => {},
        success: () => {},
        danger: () => {},
        warning: () => {},
        info: () => {},
        removeToast: () => {},
      },
      router: {},
      battleId: null,
      retroId: null,
      storyboardId: null,
      orgInviteId: null,
      teamInviteId: null,
    });
  });
});
