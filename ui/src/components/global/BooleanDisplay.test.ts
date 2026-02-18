import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import BooleanDisplay from './BooleanDisplay.svelte';

describe('BooleanDisplay component', () => {
  it('should render true state successfully', () => {
    render(BooleanDisplay, { props: { boolValue: true } });

    expect(page.getByText('True')).toBeTruthy();
  });

  it('should display false state correctly', () => {
    render(BooleanDisplay, { props: { boolValue: false } });

    expect(page.getByText('False')).toBeTruthy();
  });
});
