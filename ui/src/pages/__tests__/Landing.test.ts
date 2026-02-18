import { describe, it } from 'vitest';
import { render } from 'vitest-browser-svelte';

import Landing from '../Landing.svelte';

describe('Landing Page', () => {
  it('should render successfully', () => {
    render(Landing, {
      xfetch: () => Promise.resolve(new Response()),
    });
  });
});
