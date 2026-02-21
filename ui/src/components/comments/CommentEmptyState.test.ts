import { describe, it, expect } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CommentEmptyState from './CommentEmptyState.svelte';

describe('CommentEmptyState component', () => {
  it('should render default title and description', () => {
    render(CommentEmptyState, {});

    expect(page.getByText('No comments yet')).toBeTruthy();
    expect(page.getByText('Be the first to share your thoughts.')).toBeTruthy();
  });

  it('should render custom title', () => {
    render(CommentEmptyState, { props: { title: 'Nothing here' } });

    expect(page.getByText('Nothing here')).toBeTruthy();
    expect(page.getByText('Be the first to share your thoughts.')).toBeTruthy();
  });

  it('should render custom description', () => {
    render(CommentEmptyState, { props: { description: 'Start a conversation.' } });

    expect(page.getByText('No comments yet')).toBeTruthy();
    expect(page.getByText('Start a conversation.')).toBeTruthy();
  });

  it('should render custom title and description together', () => {
    render(CommentEmptyState, {
      props: {
        title: 'Empty activity',
        description: 'Share your first update.',
      },
    });

    expect(page.getByText('Empty activity')).toBeTruthy();
    expect(page.getByText('Share your first update.')).toBeTruthy();
  });
});
