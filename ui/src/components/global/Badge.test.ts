import { describe, it, expect, vi } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import Badge from './Badge.svelte';

describe('Badge component', () => {
  it('should render successfully', () => {
    render(Badge, { props: { label: 'TestBadge' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-gray-100');
  });

  it('should render with gray color', () => {
    render(Badge, { props: { label: 'TestBadge', color: 'gray' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-gray-100');
  });

  it('should render with blue color', () => {
    render(Badge, { props: { label: 'TestBadge', color: 'blue' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-blue-100');
  });

  it('should render with indigo color', () => {
    render(Badge, { props: { label: 'TestBadge', color: 'indigo' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-indigo-100');
  });

  it('should render with red color', () => {
    render(Badge, { props: { label: 'TestBadge', color: 'red' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-red-100');
  });

  it('should render with green color', () => {
    render(Badge, { props: { label: 'TestBadge', color: 'green' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-green-100');
  });

  it('should render with invalid color defaulting to gray', () => {
    // @ts-expect-error Testing default case with invalid color
    render(Badge, { props: { label: 'TestBadge', color: 'invalid' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-gray-100');
  });

  it('should render with custom class', () => {
    render(Badge, { props: { label: 'TestBadge', class: 'custom-class' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('custom-class');
    expect(badge.element().className).toContain('inline-flex');
  });

  it('should render with custom testId', () => {
    render(Badge, { props: { label: 'TestBadge', testId: 'custom-test-id' } });

    const badge = page.getByTestId('custom-test-id');
    expect(badge).toBeTruthy();
  });

  it('should render with custom title', () => {
    render(Badge, { props: { label: 'TestBadge', title: 'Custom Title' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().getAttribute('title')).toBe('Custom Title');
  });

  it('should use label as default title when title not provided', () => {
    render(Badge, { props: { label: 'TestBadge' } });

    const badge = page.getByText('TestBadge');
    expect(badge).toBeTruthy();
    expect(badge.element().getAttribute('title')).toBe('TestBadge');
  });

  it('should render with all optional props provided', () => {
    render(Badge, {
      props: {
        label: 'TestBadge',
        color: 'blue',
        class: 'my-class',
        testId: 'my-id',
        title: 'My Title',
      },
    });

    const badge = page.getByTestId('my-id');
    expect(badge).toBeTruthy();
    expect(badge.element().className).toContain('bg-blue-100');
    expect(badge.element().className).toContain('my-class');
    expect(badge.element().getAttribute('title')).toBe('My Title');
  });
});
