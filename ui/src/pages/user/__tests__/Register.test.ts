import { describe, it } from 'vitest';
import { render } from 'vitest-browser-svelte';

import Register from '../Register.svelte';

describe('Register Page', () => {
  it('should render successfully', () => {
    render(Register, {
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
