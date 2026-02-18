import { describe, it, expect } from 'vitest';
import { page } from 'vitest/browser';
import { render } from 'vitest-browser-svelte';

import BrowserMock from './BrowserMock.svelte';
import BrowserMockWithChildren from './__tests__/BrowserMockWithChildren.svelte';

describe('BrowserMock component', () => {
  it('should render with default props successfully', () => {
    const { container } = render(BrowserMock);

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow).toBeTruthy();
  });

  it('should render with custom class prop', () => {
    const customClass = 'border-blue-500';
    const { container } = render(BrowserMock, { props: { class: customClass } });

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow?.classList.contains('border-blue-500')).toBe(true);
  });

  it('should render with empty class prop', () => {
    const { container } = render(BrowserMock, { props: { class: '' } });

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow).toBeTruthy();
    expect(browserWindow?.classList.contains('border-gray-300')).toBe(false);
    expect(browserWindow?.classList.contains('border-blue-500')).toBe(false);
  });

  it('should render with null class prop', () => {
    const { container } = render(BrowserMock, { props: { class: null as any } });

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow).toBeTruthy();
  });

  it('should render the three browser control dots', () => {
    const { container } = render(BrowserMock);

    const dots = container.querySelectorAll('.h-3.w-3.rounded-full');
    expect(dots.length).toBe(3);

    // Check for the three colored dots
    expect(dots[0].classList.contains('bg-rose-500')).toBe(true);
    expect(dots[1].classList.contains('bg-amber-300')).toBe(true);
    expect(dots[2].classList.contains('bg-lime-400')).toBe(true);
  });

  it('should render children content', () => {
    const testContent = 'My test child content';
    render(BrowserMockWithChildren, { props: { testContent } });

    const childElement = page.getByTestId('browser-child-content');
    expect(childElement.element().textContent?.trim()).toBe(testContent);
  });

  it('should render children inside the browser mock structure', () => {
    const { container } = render(BrowserMockWithChildren, {
      props: { testContent: 'Nested child content' },
    });

    // Verify the browser window exists
    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow).toBeTruthy();

    // Verify children are rendered inside
    const childElement = container.querySelector('[data-testid="browser-child-content"]');
    expect(childElement).toBeTruthy();
    expect(childElement?.textContent?.trim()).toBe('Nested child content');
  });

  it('should apply default border classes when no class prop provided', () => {
    const { container } = render(BrowserMock);

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow?.classList.contains('border-gray-300')).toBe(true);
    expect(browserWindow?.classList.contains('dark:border-gray-700')).toBe(true);
  });

  it('should apply default border classes when class prop is explicitly undefined', () => {
    const { container } = render(BrowserMock, { props: { class: undefined } });

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow?.classList.contains('border-gray-300')).toBe(true);
    expect(browserWindow?.classList.contains('dark:border-gray-700')).toBe(true);
  });

  it('should render without children when children is explicitly undefined', () => {
    const { container } = render(BrowserMock, { props: { children: undefined } });

    const browserWindow = container.querySelector('.rounded.rounded-b-lg.shadow-xl.border');
    expect(browserWindow).toBeTruthy();

    // Verify no child content is rendered in the inner content area
    const contentAreas = container.querySelectorAll('.rounded-b-lg');
    const innerContentArea = contentAreas[1]; // Second .rounded-b-lg is the inner content area
    expect(innerContentArea?.textContent?.trim()).toBe('');
  });
});
