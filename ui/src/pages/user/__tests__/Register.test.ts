import { describe, it } from 'vitest';
import { render } from 'vitest-browser-svelte';

import Register from '../Register.svelte';

describe('Register Page', () => {
  it('should render successfully', () => {
    render(Register, {
      xfetch: () => {},
      notifications: () => {},
      router: () => {},
      battleId: null,
    });
  });
});
