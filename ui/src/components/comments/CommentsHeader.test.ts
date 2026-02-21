import { describe, it, expect } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import CommentsHeader from './CommentsHeader.svelte';

describe('CommentsHeader component', () => {
  it('should render the comment count', () => {
    render(CommentsHeader, { props: { commentsCount: 3 } });

    expect(page.getByText('3')).toBeTruthy();
  });

  it('should render singular label when count is one', () => {
    render(CommentsHeader, { props: { commentsCount: 1 } });

    expect(page.getByText('1')).toBeTruthy();
    expect(page.getByText('Comment')).toBeTruthy();
  });

  it('should render plural label when count is not one', () => {
    render(CommentsHeader, { props: { commentsCount: 2 } });

    expect(page.getByText('2')).toBeTruthy();
    expect(page.getByText('Comments')).toBeTruthy();
  });

  it('should render custom title when provided', () => {
    render(CommentsHeader, { props: { commentsCount: 5, title: 'Activity Feed' } });

    expect(page.getByText('5')).toBeTruthy();
    expect(page.getByText('Activity Feed')).toBeTruthy();
  });
});
